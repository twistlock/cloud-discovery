package azure

import (
	"encoding/base64"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/twistlock/cloud-discovery/internal/shared"
)

// Discover discovers all GCR assets
func Discover(serviceAccount string, emitFn func(result shared.CloudDiscoveryResult)) {
	sa, err := base64.RawStdEncoding.DecodeString(serviceAccount)
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	var opt Options
	if err := json.Unmarshal([]byte(sa), &opt); err != nil {
		log.Errorf(err.Error())
		return
	}
	if err := DiscoverFunctions(opt, emitFn); err != nil {
		log.Debugf(err.Error())
	}
	if err := DiscoverACR(opt, emitFn); err != nil {
		log.Debugf(err.Error())
	}
}
