package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/client"
	"github.com/docker/docker/image"
	"github.com/docker/docker/image/v1"
	"github.com/docker/docker/layer"
	digest "github.com/opencontainers/go-digest"
)

var (
	repo_endpoint = "http://192.168.1.50:5000"

	repo_name = "postgres"
	repo_tag  = "1.0.0"

	image_root = "E:\\data\\temp\\test"
)

type manifestItem struct {
	Config       string
	RepoTags     []string
	Layers       []string
	Parent       image.ID                                 `json:",omitempty"`
	LayerSources map[layer.DiffID]distribution.Descriptor `json:",omitempty"`
}

func main() {
	ctx, cancelFunc := context.WithCancel(context.TODO())
	defer cancelFunc()

	named, err := reference.WithName(repo_name)
	if err != nil {
		log.Fatalf("ParseNamed error: %v\n", err)
	}

	// ctx context.Context, name reference.Named, baseURL string, transport http.RoundTripper
	repo, err := client.NewRepository(ctx, named, repo_endpoint, http.DefaultTransport)
	if err != nil {
		log.Fatalf("ParseNamed error: %v\n", err)
	}

	mfs, err := repo.Manifests(ctx, nil)
	if err != nil {
		log.Fatalf("Manifests error: %v\n", err)
	}

	mf, err := mfs.Get(ctx, "", distribution.WithTag(repo_tag))
	if err != nil {
		log.Fatalf("Manifests error: %v\n", err)
	}

	var imageConfig []byte

	// 读config并写入文件
	for _, ref := range mf.References() {
		if ref.MediaType == schema2.MediaTypeImageConfig {
			imageConfig, err = repo.Blobs(ctx).Get(ctx, ref.Digest)
			if err != nil {
				log.Fatalf("Read config error %v\n", err)
			}
			ioutil.WriteFile(filepath.Join(image_root, ref.Digest.Hex()+".json"), imageConfig, 0644)
		}
	}

	// 解析config
	img, err := image.NewFromJSON(imageConfig)
	if err != nil {
		log.Fatalf("Image NewFromJSON error: %v\n", err)
	}

	var manifest []manifestItem
	var layers []string
	var parent digest.Digest

	for i := range img.RootFS.DiffIDs {
		v1Img := image.V1Image{}
		if i == len(img.RootFS.DiffIDs)-1 {
			v1Img = img.V1Image
		}
		rootFS := *img.RootFS
		rootFS.DiffIDs = rootFS.DiffIDs[:i+1]
		v1ID, err := v1.CreateID(v1Img, rootFS.ChainID(), parent)
		if err != nil {
			log.Fatalf("CreateID error: %v\n", err)
		}

		v1Img.ID = v1ID.Hex()
		if parent != "" {
			v1Img.Parent = parent.Hex()
		}
		log.Printf("layer id: %s\n", v1Img.ID)

		outDir := filepath.Join(image_root, v1Img.ID)
		if err := os.Mkdir(outDir, 0755); err != nil {
			log.Fatalf("Mkdir error: %v\n", err)
		}

		if err := ioutil.WriteFile(filepath.Join(outDir, "VERSION"), []byte("1.0"), 0644); err != nil {
			log.Fatalf("WriteFile error: %v\n", err)
		}

		v1ImgConfig, err := json.Marshal(v1Img)
		if err != nil {
			log.Fatalf("WriteFile error: %v\n", err)
		}

		if err := ioutil.WriteFile(filepath.Join(outDir, "json"), v1ImgConfig, 0644); err != nil {
			log.Fatalf("WriteFile error: %v\n", err)
		}

		if s2mf, ok := mf.(*schema2.DeserializedManifest); ok {
			downloadLayer(ctx, repo, s2mf.Layers[i].Digest, outDir)
		}

		layers = append(layers, filepath.Join(v1Img.ID, "layer.tar"))
		parent = v1ID
	}

	var repoTags []string
	repoTags = append(repoTags, named.Name())

	manifest = append(manifest, manifestItem{
		Config:   digest.FromBytes(imageConfig).String() + ".json",
		RepoTags: repoTags,
		Layers:   layers,
	})

	manifestFileName := filepath.Join(image_root, "manifest.json")
	f, err := os.OpenFile(manifestFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("OpenFile manifest error: %v\n", err)
	}

	if err := json.NewEncoder(f).Encode(manifest); err != nil {
		f.Close()
		log.Fatalf("Encode manifest error: %v\n", err)
	}
	f.Close()
}

func downloadLayer(ctx context.Context, repo distribution.Repository, digest digest.Digest, outDir string) {
	layer, err := repo.Blobs(ctx).Get(ctx, digest)
	if err != nil {
		log.Fatalf("Read config error %v\n", err)
	}

	gr, err := gzip.NewReader(bytes.NewReader(layer))
	if err != nil {
		log.Fatalf("Open file error %v\n", err)
	}

	f, err := os.OpenFile(filepath.Join(outDir, "layer.tar"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Open file error %v\n", err)
	}

	_, err = io.Copy(f, gr)
	if err != nil {
		log.Fatalf("Write file error %v\n", err)
	}
}
