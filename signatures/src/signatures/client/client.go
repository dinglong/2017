package client

import (
	"log"

	"github.com/docker/notary"
	"github.com/docker/notary/client"
	"github.com/docker/notary/trustpinning"
	"github.com/docker/notary/tuf/data"
	"net/http"
)

var (
	notaryCachePath = "/root/notary"
	trustPin        trustpinning.TrustPinConfig
	mockRetriever   notary.PassRetriever
)

// GetTargets is a help function called by API to fetch signature information of a given repository.
// Per docker's convention the repository should contain the information of endpoint, i.e. it should look
// like "192.168.0.1/library/ubuntu", instead of "library/ubuntu" (fqRepo for fully-qualified repo)
func GetTargets(notaryEndpoint string, fqRepo string) ([]*client.TargetWithRole, error) {
	gun := data.GUN(fqRepo)
	tr := &http.Transport{}

	notaryRepo, err := client.NewFileCachedNotaryRepository(notaryCachePath, gun, notaryEndpoint, tr, mockRetriever, trustPin)
	if err != nil {
		return nil, err
	}

	targets, err := notaryRepo.ListTargets(data.CanonicalTargetsRole)
	if _, ok := err.(client.ErrRepositoryNotExist); ok {
		log.Printf("Repository not exist, repo: %s, error: %v, returning empty signature", fqRepo, err)
		return nil, err
	} else if err != nil {
		return nil, err
	}

	return targets, nil
}
