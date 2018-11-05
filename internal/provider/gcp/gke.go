package gcp

import (
	"encoding/json"
	"fmt"
	"github.com/twistlock/cloud-discovery/internal/shared"
	"google.golang.org/api/container/v1"
	"io/ioutil"
	"net/http"
)

// DiscoverGKE retrieves a list of all GKE data
func DiscoverGKE(opt Options, emitFn func(result shared.CloudDiscoveryResult)) error {
	client, projectID, err := client(opt.ServiceAccount)
	if err != nil {
		return err
	}
	var nextToken string

	// Get all available zones
	// https://cloud.google.com/kubernetes-engine/docs/reference/rest/
	resp, err := client.Get(fmt.Sprintf("https://container.googleapis.com/v1beta1/projects/%s/locations", projectID))
	if err != nil {
		return err
	}
	var locations struct {
		Locations []struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"locations"`
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &locations); err != nil {
		return err
	}

	for _, location := range locations.Locations {
		if location.Type != "ZONE" {
			continue
		}
		// https://cloud.google.com/kubernetes-engine/docs/reference/rest/
		url := fmt.Sprintf("https://container.googleapis.com/v1beta1/projects/%s/locations/%s/clusters", projectID, location.Name)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		// Set next page token in case there are more results to be queried
		if nextToken != "" {
			q := req.URL.Query()
			q.Add("pageToken", nextToken)
		}

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var clusters container.ListClustersResponse
		if err := json.Unmarshal(body, &clusters); err != nil {
			return err
		}

		res := shared.CloudDiscoveryResult{
			Region: location.Name,
			Type:   "GKE",
		}
		for _, f := range clusters.Clusters {
			res.Assets = append(res.Assets, shared.CloudAsset{ID: f.Name, Data: f})
		}
		emitFn(res)
	}
	return nil
}
