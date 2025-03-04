package docker

import (
	"fmt"

	"github.com/save-abandoned-projects/ignite/pkg/network"
	dockernetwork "github.com/save-abandoned-projects/ignite/pkg/network/docker"
	"github.com/save-abandoned-projects/ignite/pkg/providers"
	"github.com/save-abandoned-projects/ignite/pkg/runtime"
	log "github.com/sirupsen/logrus"
)

func SetDockerNetwork() error {
	log.Trace("Initializing the Docker network provider...")
	if providers.Runtime.Name() != runtime.RuntimeDocker {
		return fmt.Errorf("the %q network plugin can only be used with the %q runtime", network.PluginDockerBridge, runtime.RuntimeDocker)
	}

	providers.NetworkPlugin = dockernetwork.GetDockerNetworkPlugin(providers.Runtime)
	return nil
}
