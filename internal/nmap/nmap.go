package nmap

import (
	"fmt"
	"github.com/globalsign/mgo"
	log "github.com/sirupsen/logrus"
	"github.com/twistlock/cloud-discovery/internal/shared"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func Nmap(subnet string, minPort, maxPort int, nmapWriter io.Writer, emitFn func(result shared.CloudNmapResult)) error {
	log.Debugf("Scanning %s", subnet)
	dir, err := ioutil.TempDir("", "nmap")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	resultPath := filepath.Join(dir, "nmap")
	// https://nmap.org/book/nping-man-ou	tput-options.html
	cmd := exec.Command(
		"nmap", subnet,
		"-v", "1",
		"-sT",
		"--max-retries", "0",
		"-p", fmt.Sprintf("%d-%d", minPort, maxPort),
		"-oX", resultPath,
		"--max-scan-delay", "3")
	cmd.Stdout = nmapWriter
	cmd.Stderr = nmapWriter
	if err := cmd.Run(); err != nil {
		return err
	}
	out, err := ioutil.ReadFile(resultPath)
	if err != nil {
		return err
	}
	nmap, err := Parse(out)
	if err != nil {
		return err
	}

	const (
		mongo          = "mongod"
		dockerRegistry = "docker registry"
		mysql          = "mysql"
	)

	for _, host := range nmap.Hosts {
		for _, port := range host.Ports {
			if len(host.Addresses) == 0 {
				continue
			}
			if host.Addresses[0].Addr == "0.0.0.0" {
				continue
			}
			client := http.Client{Timeout: time.Second * 2}
			addr := fmt.Sprintf("%s:%d", host.Addresses[0].Addr, port.PortId)
			service := port.Service.Name
			log.Debugf("Checking port %v %v %v", host.Addresses[0], port.PortId, port.Protocol)
			if service == "unknown" || (port.PortId >= 5000 && port.PortId <= 6000 && port.Protocol == "tcp") {
				resp, err := client.Get(fmt.Sprintf("http://%s/v2/_catalog", addr))
				if err == nil {
					out, _ := ioutil.ReadAll(resp.Body)
					resp.Body.Close()
					respBody := strings.ToLower(string(out))
					if svc := func() string {
						if strings.Contains(respBody, "mongo") {
							return mongo
						}
						if strings.Contains(respBody, "packets") && strings.Contains(respBody, "order") {
							return mysql

						}
						for h, _ := range resp.Header {
							if strings.Contains(strings.ToLower(h), "docker") {
								return dockerRegistry
							}
						}
						return ""
					}(); svc != "" {
						service = svc
					}
				}
			}
			result := shared.CloudNmapResult{
				Host: host.Addresses[0].Addr,
				Port: port.PortId,
				App:  service,
			}
			switch service {
			case mongo:
				conn, err := mgo.DialWithTimeout(fmt.Sprintf("mongodb://%s", addr), time.Second*1)
				if err == nil {
					_, err := conn.DatabaseNames()
					if err == nil {
						result.Insecure = true
					}
				}
			case dockerRegistry:
				resp, err := client.Get(fmt.Sprintf("http://%s/v2/_catalog", addr))
				if err == nil {
					if resp.StatusCode == http.StatusOK {
						result.Insecure = true
					}
					resp.Body.Close()
				}
			}

			emitFn(result)
		}
	}
	return nil
}
