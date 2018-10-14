package nmap

import (
	"fmt"
	"net/http"
	"strings"
)

type registryMapper struct{}

func NewRegistryMapper() *registryMapper { return &registryMapper{} }

func (m *registryMapper) App() string {
	return "docker registry"
}

func (m *registryMapper) HasApp(port int, respHeader http.Header, body string) bool {
	for h, _ := range respHeader {
		if strings.Contains(strings.ToLower(h), "docker") {
			return true
		}
	}
	return false
}

func (m *registryMapper) KnownPorts() []int {
	// https://docs.mongodb.com/manual/reference/default-mongodb-port/
	return []int{5000}
}

func (m *registryMapper) Insecure(addr string) (bool, string) {
	resp, err := insecureClient().Get(fmt.Sprintf("http://%s/v2/_catalog", addr))
	if err == nil {
		resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			return true, "Missing authorization"
		}
	}
	return false, ""
}
