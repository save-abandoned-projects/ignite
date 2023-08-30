package cni

import (
	"github.com/save-abandoned-projects/ignite/pkg/network/cni"
	"github.com/save-abandoned-projects/ignite/pkg/providers"
	log "github.com/sirupsen/logrus"
)

func SetCNINetworkPlugin() (err error) {
	log.Trace("Initializing the CNI provider...")
	providers.NetworkPlugin, err = cni.GetCNINetworkPlugin(providers.Runtime)
	return
}
