package gcp

// Options are options for getting cloud data for GCP services
type Options struct {
	ServiceAccount string // ServiceAccount is the secret key used for authentication (in JSON format)
	Region         string // Region is the region to query
}
