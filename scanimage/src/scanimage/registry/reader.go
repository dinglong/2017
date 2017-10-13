package registry

import (
	"context"
	"fmt"
	"net/http"

	"github.com/docker/distribution"
	_ "github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/client"
)

var (
	MediaTypeManifest = "application/vnd.docker.distribution.manifest.v2+json"
)

func ReadManifest(registryEndpoint, repository, tag string) distribution.Manifest {
	ctx, cancelFunc := context.WithCancel(context.TODO())
	defer cancelFunc()

	named, err := reference.ParseNamed(repository + ":" + tag)
	if err != nil {
		fmt.Printf("ParseNamed error: %v\n", err)
		return nil
	}

	// ctx context.Context, name reference.Named, baseURL string, transport http.RoundTripper
	repo, err := client.NewRepository(ctx, named, registryEndpoint, http.DefaultTransport)
	if err != nil {
		fmt.Printf("ParseNamed error: %v\n", err)
		return nil
	}
	fmt.Printf("repo %v\n", repo)

	mfs, err := repo.Manifests(ctx, nil)
	if err != nil {
		fmt.Printf("Manifests error: %v\n", err)
		return nil
	}
	fmt.Printf("mfs: %v\n", mfs)

	mf, err := mfs.Get(ctx, "", distribution.WithTag(tag))
	if err != nil {
		fmt.Printf("Manifests error: %v\n", err)
		return nil
	}

	fmt.Printf("mf: %v\n", mf)

	return mf
}
