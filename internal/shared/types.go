package shared

// CloudDiscoveryResults contains the result of the cloud discovery
type CloudDiscoveryResults struct {
	Results []CloudDiscoveryResult
}

type CloudDiscoveryResult struct {
	Region string
	Type string
	Assets []CloudAsset
}

type CloudAsset struct {
	ID string
}

type Discoverer interface {
	Discover() (*CloudDiscoveryResult, error)
}