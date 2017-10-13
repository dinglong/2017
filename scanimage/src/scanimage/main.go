package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"scanimage/registry"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/schema2"
)

type LayerEnvelope struct {
	Layer *Layer `json:"Layer,omitempty"`
	Error *Error `json:"Error,omitempty"`
}

type Error struct {
	Message string `json:"Message,omitempty"`
}

type Layer struct {
	Name             string            `json:"Name,omitempty"`
	NamespaceName    string            `json:"NamespaceName,omitempty"`
	Path             string            `json:"Path,omitempty"`
	Headers          map[string]string `json:"Headers,omitempty"`
	ParentName       string            `json:"ParentName,omitempty"`
	Format           string            `json:"Format,omitempty"`
	IndexedByVersion int               `json:"IndexedByVersion,omitempty"`
	Features         []Feature         `json:"Features,omitempty"`
}

type Feature struct {
	Name            string          `json:"Name,omitempty"`
	NamespaceName   string          `json:"NamespaceName,omitempty"`
	VersionFormat   string          `json:"VersionFormat,omitempty"`
	Version         string          `json:"Version,omitempty"`
	Vulnerabilities []Vulnerability `json:"Vulnerabilities,omitempty"`
	AddedBy         string          `json:"AddedBy,omitempty"`
}

type Vulnerability struct {
	Name          string                 `json:"Name,omitempty"`
	NamespaceName string                 `json:"NamespaceName,omitempty"`
	Description   string                 `json:"Description,omitempty"`
	Link          string                 `json:"Link,omitempty"`
	Severity      string                 `json:"Severity,omitempty"`
	Metadata      map[string]interface{} `json:"Metadata,omitempty"`
	FixedBy       string                 `json:"FixedBy,omitempty"`
	FixedIn       []Feature              `json:"FixedIn,omitempty"`
}

var client = &http.Client{}

// ScanLayer calls Clair's API to scan a layer.
func ScanLayer(clairEndpoint string, l Layer) error {
	layer := LayerEnvelope{
		Layer: &l,
		Error: nil,
	}

	data, err := json.Marshal(layer)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, clairEndpoint+"/v1/layers", bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set(http.CanonicalHeaderKey("Content-Type"), "application/json")
	_, err = send(req, http.StatusCreated)

	return err
}

func GetResult(clairEndpoint string, layerName string) (*LayerEnvelope, error) {
	req, err := http.NewRequest(http.MethodGet, clairEndpoint+"/v1/layers/"+layerName+"?features&vulnerabilities", nil)
	if err != nil {
		return nil, err
	}
	b, err := send(req, http.StatusOK)
	if err != nil {
		return nil, err
	}
	var res LayerEnvelope
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func send(req *http.Request, expectedStatus int) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != expectedStatus {
		return nil, fmt.Errorf("Unexpected status code: %d, text: %s", resp.StatusCode, string(b))
	}
	return b, nil
}

func prepareLayers(repository string, descriptors []distribution.Descriptor) []Layer {
	ls := make([]Layer, 0)
	tokenHeader := map[string]string{"Connection": "close"}

	// form the chain by using the digests of all parent layers in the image, such that if another image is built on top of this image the layer name can be re-used.
	shaChain := ""
	for _, d := range descriptors {
		if d.MediaType == schema2.MediaTypeImageConfig {
			continue
		}

		shaChain += string(d.Digest) + "-"
		l := Layer{
			Name:    fmt.Sprintf("%x", sha256.Sum256([]byte(shaChain))),
			Headers: tokenHeader,
			Format:  "Docker",
			Path:    fmt.Sprintf("%s/v2/%s/blobs/%s", "http://registry:5000", repository, string(d.Digest)),
		}
		ls = append(ls, l)
	}

	return ls
}

func main() {
	mf := registry.ReadManifest("http://192.168.1.50:5000", "test", "1.0.0")
	ls := prepareLayers("test", mf.References())

	for _, l := range ls {
		fmt.Printf("layer [%v]\n", l)
		err := ScanLayer("http://192.168.1.50:6060", l)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		le, err := GetResult("http://192.168.1.50:6060", l.Name)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		} else {
			s, err := json.Marshal(le)
			if err != nil {
				fmt.Printf("err: %v\n", err)
			} else {
				fmt.Printf("%s\n", s)
			}
		}
	}
}
