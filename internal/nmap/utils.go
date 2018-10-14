package nmap

import (
	"crypto/tls"
	"net/http"
	"time"
)

// insecureClient retuns an HTTP client without TLS validations
func insecureClient() *http.Client {
	return &http.Client{Timeout: time.Second * 2, Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
}

// portMapper is a generic port mapper for known apps
type portMapper struct {
	app   string
	ports []int
}

func NewPortMapper(app string, ports ...int) *portMapper                        { return &portMapper{app: app, ports: ports} }
func (p *portMapper) App() string                                               { return p.app }
func (p *portMapper) HasApp(port int, respHeader http.Header, body string) bool { return false }
func (p *portMapper) KnownPorts() []int                                         { return p.ports }
func (p *portMapper) Insecure(addr string) (bool, string)                       { return false, "" }
