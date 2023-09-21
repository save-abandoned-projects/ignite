package network

import (
	"fmt"

	"github.com/save-abandoned-projects/ignite/pkg/network"
	"github.com/save-abandoned-projects/ignite/pkg/providers"
	cniprovider "github.com/save-abandoned-projects/ignite/pkg/providers/cni"
	dockerprovider "github.com/save-abandoned-projects/ignite/pkg/providers/docker"
)

func SetNetworkPlugin() error {
	switch providers.NetworkPluginName {
	case network.PluginDockerBridge:
		return dockerprovider.SetDockerNetwork() // Use the Docker bridge network
	case network.PluginCNI:
		return cniprovider.SetCNINetworkPlugin() // Use the CNI Network plugin
	}

	return fmt.Errorf("unknown network plugin %q", providers.NetworkPluginName)
}
