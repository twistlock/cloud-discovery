package gcp

import (
	"encoding/json"
	"fmt"
	"github.com/twistlock/cloud-discovery/internal/shared"
	"google.golang.org/api/cloudfunctions/v1"
	"io/ioutil"
	"net/http"
)

type functionsClient struct {
	opt Options
}

// NewFunctionsClient creates a new google functions client
func NewFunctionsClient(opt Options) *functionsClient {
	return &functionsClient{opt: opt}
}

// Discover retrieves a list of functions for the settings provided to the client
func (s *functionsClient) Discover() (*shared.CloudDiscoveryResult, error) {
	client, projectID, err := client(s.opt.ServiceAccount)
	if err != nil {
		return nil, err
	}
	var res []shared.CloudAsset
	var nextToken string
	truncated := true
	// https://cloud.google.com/functions/docs/reference/rest/v1/projects.locations.functions/list
	url := fmt.Sprintf("https://cloudfunctions.googleapis.com/v1/projects/%s/locations/%s/functions", projectID, s.opt.Region)
	for truncated {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		// Set next page token in case there are more results to be queried
		if nextToken != "" {
			q := req.URL.Query()
			q.Add("pageToken", nextToken)
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var functions cloudfunctions.ListFunctionsResponse
		if err := json.Unmarshal(body, &functions); err != nil {
			return nil, err
		}

		for _, f := range functions.Functions {
			res = append(res, shared.CloudAsset{ID: f.Name, Data: f})
		}
		if functions.NextPageToken == "" {
			truncated = false
		}
		nextToken = functions.NextPageToken
	}
	return &shared.CloudDiscoveryResult{
		Region: s.opt.Region,
		Type:   "Functions",
		Assets: res,
	}, nil
}
