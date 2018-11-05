package gcp

import (
	"encoding/json"
	"fmt"
	"github.com/twistlock/cloud-discovery/internal/shared"
	"io/ioutil"
	"net/http"
)

type gcrClient struct {
	opt Options
}

// NewGCRClient creates a new google container registry client
func NewGCRClient(opt Options) *gcrClient {
	return &gcrClient{opt: opt}
}

// Discover retrieves all repositories for the settings provided to the client
func (s *gcrClient) Discover() (*shared.CloudDiscoveryResult, error) {
	client, _, err := client(s.opt.ServiceAccount)
	if err != nil {
		return nil, err
	}
	// Use catalog API
	// https://docs.docker.com/registry/spec/api/#catalog
	catalogUrl := fmt.Sprintf("https://%s/v2/_catalog", s.opt.Region)
	result := &shared.CloudDiscoveryResult{
		Region: s.opt.Region,
		Type:   "GCR",
	}
	b, err := client.Get(catalogUrl)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(b.Body)
	b.Body.Close()
	if err != nil {
		return nil, err
	}
	if b.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to to query %s %s %s", catalogUrl, b.Status, string(body))
	}
	var response struct {
		Repositories []string `json:"repositories"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	result.Assets = append(result.Assets)
	for _, repo := range response.Repositories {
		result.Assets = append(result.Assets, shared.CloudAsset{
			ID: repo,
		})
	}
	return result, nil
}
