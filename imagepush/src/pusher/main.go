package main

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/docker/distribution"
	"github.com/docker/distribution/digest"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/client"
	"io/ioutil"
)

var (
	// image_dir = "E:\\data\\temp\\alpine~3.5.2"
	image_dir      = "E:\\data\\temp\\postgres~9.6~alpine"
	image_manifest = "manifest.json"

	// resp_name = "allx"
	resp_name = "postgres"
	resp_tag  = "1.0.0"

	baseUrl = "http://192.168.1.50:5000"
)

type manifestItem struct {
	Config   string
	RepoTags []string
	Layers   []string
}

func main() {
	ctx, cancelFunc := context.WithCancel(context.TODO())
	defer cancelFunc()

	// new repository
	named, err := reference.ParseNamed(resp_name + ":" + resp_tag)
	if err != nil {
		log.Fatalf("ParseNamed error: %v\n", err)
	}
	repo, err := client.NewRepository(ctx, named, baseUrl, http.DefaultTransport)

	// upload image
	manifestFile, err := os.Open(filepath.Join(image_dir, image_manifest))
	if err != nil {
		log.Fatalf("Open manifest error, %v\n", err)
	}
	defer manifestFile.Close()

	var manifest []manifestItem
	if err := json.NewDecoder(manifestFile).Decode(&manifest); err != nil {
		log.Fatalf("Decode manifest error, %v\n", err)
	}

	bs := repo.Blobs(ctx)
	for _, m := range manifest {
		var descriptors []distribution.Descriptor

		// upload layer
		for _, l := range m.Layers {
			createOpts := []distribution.BlobCreateOption{}
			bw, err := bs.Create(ctx, createOpts...)
			if err != nil {
				log.Fatalf("Create blob write error, %v\n", err)
			}

			reader, err := os.Open(filepath.Join(image_dir, l))
			if err != nil {
				log.Fatalf("Open layer error, %v\n", err)
			}

			dt := digest.Canonical.New()
			tee := io.TeeReader(bufio.NewReader(reader), dt.Hash())

			nn, err := bw.ReadFrom(tee)
			if err != nil {
				log.Fatalf("Upload layer error, %v\n", err)
			}
			// TODO defer layerUpload.Close()
			reader.Close()

			pushDigest := dt.Digest()
			if _, err := bw.Commit(ctx, distribution.Descriptor{Digest: pushDigest}); err != nil {
				log.Fatalf("Commit layer error, %v\n", err)
			}

			desc := distribution.Descriptor{
				Digest:    pushDigest,
				MediaType: schema2.MediaTypeLayer,
				Size:      nn,
			}
			descriptors = append(descriptors, desc)

			log.Printf("Uploaded layer %s (%s), %d bytes", l, dt.Digest(), nn)
		}

		// upload manifest
		imgConfig, err := ioutil.ReadFile(filepath.Join(image_dir, m.Config))
		if err != nil {
			log.Fatalf("Read config error, %v\n", err)
		}

		builder := schema2.NewManifestBuilder(repo.Blobs(ctx), schema2.MediaTypeImageConfig, imgConfig)
		for _, d := range descriptors {
			if err := builder.AppendReference(d); err != nil {
				log.Fatalf("AppendReference error, %v\n", err)
			}
		}

		manifest, err := builder.Build(ctx)
		if err != nil {
			log.Fatalf("Build manifest error, %v\n", err)
		}

		manSvc, err := repo.Manifests(ctx)
		if err != nil {
			log.Fatalf("Manifest service error, %v\n", err)
		}

		putOptions := []distribution.ManifestServiceOption{distribution.WithTag(resp_tag)}
		if _, err = manSvc.Put(ctx, manifest, putOptions...); err != nil {
			log.Fatalf("Put manifest error, %v\n", err)
		}
		log.Printf("Put manifest over\n")
	}
}
