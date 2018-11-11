package nmap

import (
	"fmt"
	gonmap "github.com/lair-framework/go-nmap"
	log "github.com/sirupsen/logrus"
	"github.com/twistlock/cloud-discovery/internal/shared"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func Nmap(wr io.Writer, subnet string, verbose bool) {
	tw := shared.NewTabWriter(wr)
	var nmapWriter io.Writer
	if verbose {
		nmapWriter = wr
	} else {
		nmapWriter = os.Stdout
	}
	fmt.Fprintf(tw, "\nHost\tPort\tApp\tInsecure\tReason\t\n")
	if err := nmap(subnet, 30, 30000, nmapWriter, func(result shared.CloudNmapResult) {
		fmt.Fprintf(tw, "%s\t%d\t%s\t%t\t%s\t\n", result.Host, result.Port, result.App, result.Insecure, result.Reason)
	}); err != nil {
		log.Error(err)
	}
	tw.Flush()
}

func nmap(subnet string, minPort, maxPort int, nmapWriter io.Writer, emitFn func(result shared.CloudNmapResult)) error {
	log.Debugf("Scanning subnet %s", subnet)
	dir, err := ioutil.TempDir("", "nmap")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	resultPath := filepath.Join(dir, "nmap")
	// https://nmap.org/book/nping-man-output-options.html
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
	nmap, err := gonmap.Parse(out)
	if err != nil {
		return err
	}
	mappers := []Mapper{
		NewMongoMapper(),
		NewMysqlMapper(),
		NewRegistryMapper(),
		NewPortMapper("kube-apiserver", 6443),
		NewPortMapper("kube-router", 20244),
		NewPortMapper("kube-proxy", 10256),
		NewPortMapper("kubelet", 10255),
		NewKubeletMapper()}
	client := http.Client{Timeout: time.Second * 2}
	for _, host := range nmap.Hosts {
		for _, targetPort := range host.Ports {
			if len(host.Addresses) == 0 {
				continue
			}
			if host.Addresses[0].Addr == "0.0.0.0" {
				continue
			}
			addr := fmt.Sprintf("%s:%d", host.Addresses[0].Addr, targetPort.PortId)
			app := targetPort.Service.Name
			log.Debugf("Checking target port %v %v %v", host.Addresses[0], targetPort.PortId, targetPort.Protocol)
			if app == "unknown" || app == "" || (targetPort.PortId >= 5000 && targetPort.PortId <= 30000 && targetPort.Protocol == "tcp") {
				resp, err := client.Get(fmt.Sprintf("http://%s/v2/_catalog", addr))
				found := false
				if err == nil {
					out, _ := ioutil.ReadAll(resp.Body)
					resp.Body.Close()
					body := string(out)
					for _, mapper := range mappers {
						if mapper.HasApp(targetPort.PortId, resp.Header, body) {
							app = mapper.App()
							found = true
							break
						}
					}
				}
				if !found {
					// Fallback to known port
					for _, mapper := range mappers {
						for _, port := range mapper.KnownPorts() {
							if port == targetPort.PortId {
								app = mapper.App()
								break
							}
						}
					}
				}
			}
			result := shared.CloudNmapResult{
				Host: host.Addresses[0].Addr,
				Port: targetPort.PortId,
				App:  app,
			}
			for _, mapper := range mappers {
				if mapper.App() == app {
					insecure, reason := mapper.Insecure(addr)
					if insecure { // Multiple mappers can apply to the same app, take precedence based on insecure findings
						result.Insecure = insecure
						result.Reason = reason
						break
					}
				}
			}
			emitFn(result)
		}
	}
	return nil
}
