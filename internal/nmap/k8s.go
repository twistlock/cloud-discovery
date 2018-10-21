package nmap

import (
	"fmt"
	"net/http"
)

type kubeletMapper struct{}

func NewKubeletMapper() *kubeletMapper { return &kubeletMapper{} }

func (m *kubeletMapper) App() string {
	return "kubelet"
}

func (m *kubeletMapper) HasApp(port int, respHeader http.Header, body string) bool {
	return false
}

func (m *kubeletMapper) KnownPorts() []int {
	return []int{10250}
}

func (m *kubeletMapper) Insecure(addr string) (bool, string) {
	resp, err := insecureClient().Get(fmt.Sprintf("https://%s", addr))
	if err == nil {
		resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			return true, "missing authorization for metrics API"
		}
	}
	return false, ""
}
