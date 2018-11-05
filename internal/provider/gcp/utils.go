package gcp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/twistlock/cloud-discovery/internal/shared"
	"golang.org/x/oauth2/google"
	"net/http"
)

// Discover discovers all GCR assets
func Discover(serviceAccount string, emitFn func(result shared.CloudDiscoveryResult)) {
	sa, err := base64.RawStdEncoding.DecodeString(serviceAccount)
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	var discoverers []shared.Discoverer
	opt := Options{
		ServiceAccount: string(sa),
	}

	for _, region := range gcrRegions {
		opt.Region = region
		discoverers = append(discoverers, NewGCRClient(opt))
	}
	for _, region := range functionRegions {
		opt.Region = region
		discoverers = append(discoverers, NewFunctionsClient(opt))
	}
	if err := DiscoverGKE(opt, emitFn); err != nil {
		log.Debugf(err.Error())
	}
	for _, discoverer := range discoverers {
		result, err := discoverer.Discover()
		if err != nil {
			log.Debugf(err.Error())
		} else if len(result.Assets) > 0 {
			emitFn(*result)
		}
	}
}

// functionRegions are regions for GCP cloud functions
// See https://cloud.google.com/functions/docs/locations
var functionRegions = []string{
	"europe-west1",
	"asia-northeast1",
	"us-central1",
	"us-east1",
}

var gcrRegions = []string{
	"gcr.io",
	"us.gcr.io",
	"eu.gcr.io",
	"asia.gcr.io",
}

func client(sa string) (client *http.Client, projectID string, err error) {
	// Use service key to create an authentication token
	conf, err := google.JWTConfigFromJSON([]byte(sa), "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, "", err
	}
	var serviceAccount struct {
		ProjectID string `json:"project_id"`
	}
	if err := json.Unmarshal([]byte(sa), &serviceAccount); err != nil {
		return nil, "", err
	}
	return conf.Client(context.Background()), serviceAccount.ProjectID, nil
}
