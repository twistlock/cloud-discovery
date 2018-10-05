package shared

import "time"

// CloudDiscoveryResults contains the result of the cloud discovery
type CloudDiscoveryResults struct {
	Modified time.Time
	Results  []CloudDiscoveryResult
}

// CloudDiscoveryResult is the result of scanning a specific cloud provider region and type
type CloudDiscoveryResult struct {
	Region string       `json:"region"`  // Region is the provider region/projects
	Type   string       `json:"type"`    // Type is the provider result type
	Assets []CloudAsset `json:"assests"` // Assets are the list of assets discovered in the cloud scan
}

// CloudAsset represents an entity (cluster/container/etc...) for a specific cloud provider
type CloudAsset struct {
	ID   string      `json:"id"`   // ID is the assest ID
	Data interface{} `json:"data"` // Data is expanded customized asset data
}

type Discoverer interface {
	Discover() (*CloudDiscoveryResult, error)
}

// Provider is the cloud provider
type Provider string

const (
	ProviderAWS Provider = "aws"
)

// Credentials holds authentication data for a specific provider
type Credentials struct {
	Provider Provider `json:"provider"` // Provider is the authentication provider (AWS/Azure/GCP)
	ID       string   `json:"id"`       // ID is the access key id used to access the provider data
	Secret   string   `json:"secret"`   // Secret is the access key secret
}

type CloudDiscoveryRequest struct {
	Credentials []Credentials `json:"credentials"`
}
