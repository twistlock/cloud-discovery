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
	ProviderGCP Provider = "gcp"
)

// Credentials holds authentication data for a specific provider
type Credentials struct {
	Provider Provider `json:"provider"` // Provider is the authentication provider (AWS/Azure/GCP)
	ID       string   `json:"id"`       // ID is the access key id used to access the provider data
	Secret   string   `json:"secret"`   // Secret is the access key secret
}

// CloudDiscoveryRequest repsents a request to scan a cloud provider using a set of credentials
type CloudDiscoveryRequest struct {
	Credentials []Credentials `json:"credentials"`
}

type CloudNmapRequest struct {
	Subnet     string `json:"subnet"`  // Subnet is the subnet to scan
	AutoDetect bool   `json:"auto"`    // AutoDetect indicates subnet should be auto-detected
	Verbose    bool   `json:"verbose"` // Verbose indicates whether to output debug data from nmap
}

// CloudNmapResult is a single host-port result of a cloud discovery nmap query
type CloudNmapResult struct {
	Host     string `json:"host"`     // Host is the target host IP
	Port     int    `json:"port"`     // Port is the target host port
	App      string `json:"app"`      // App is the name of the detected app
	Insecure bool   `json:"insecure"` // Insecure indicates whether the app has secure
	Reason   string `json:"reason"`   // Reason provides detailed reasoning when the app is insecure
}
