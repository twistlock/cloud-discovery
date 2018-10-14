package nmap

import "net/http"

// Mapper detects apps and checks if they are securely configured
type Mapper interface {
	// App returns the app associated with the mapper
	App() string
	// KnownPorts returns the app known ports
	KnownPorts() []int
	// HasApp returns true if the response is associated with the targeted app
	HasApp(port int, header http.Header, body string) bool
	// Insecure returns whether the app is insecure and the reason for that
	Insecure(addr string) (bool, string)
}
