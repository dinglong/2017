package main

import (
	"context"
	"fmt"
	"github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/client"
	"net/http"
)

var (
	MediaTypeManifest = "application/vnd.docker.distribution.manifest.v2+json"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.TODO())
	defer cancelFunc()

	// ctx context.Context, baseURL string, transport http.RoundTripper
	registry, err := client.NewRegistry(ctx, "http://192.168.1.61:5000", http.DefaultTransport)
	if err != nil {
		fmt.Printf("new registry error: %v\n", err)
		return
	}

	///////////////////////////////////////////////////////////////////////////////////////////////////////

	repos := make([]string, 5, 5)
	size, err := registry.Repositories(ctx, repos, "")
	if err != nil {
		fmt.Printf("Repositories error: %v\n", err)
	}

	for i := 0; i < size; i++ {
		fmt.Printf("repos: %s\n", repos[i])
	}

	///////////////////////////////////////////////////////////////////////////////////////////////////////

	named, err := reference.ParseNamed("firstimage:1.0.0")
	if err != nil {
		fmt.Printf("ParseNamed error: %v\n", err)
		return
	}

	// ctx context.Context, name reference.Named, baseURL string, transport http.RoundTripper
	repo, err := client.NewRepository(ctx, named, "http://192.168.1.61:5000", http.DefaultTransport)
	if err != nil {
		fmt.Printf("ParseNamed error: %v\n", err)
		return
	}
	fmt.Printf("repo %v\n", repo)

	mfs, err := repo.Manifests(ctx, nil)
	if err != nil {
		fmt.Printf("Manifests error: %v\n", err)
		return
	}
	fmt.Printf("mfs: %v\n", mfs)

	manifestListFunc := func(b []byte) (distribution.Manifest, distribution.Descriptor, error) {
		fmt.Printf("date: %s\n", string(b))
		return nil, distribution.Descriptor{}, err
	}
	distribution.RegisterManifestSchema(MediaTypeManifest, manifestListFunc)

	mf, err := mfs.Get(ctx, "", distribution.WithTagOption{"1.0.0"})
	if err != nil {
		fmt.Printf("Manifests error: %v\n", err)
		return
	}

	fmt.Printf("mf: %v\n", mf)
}
