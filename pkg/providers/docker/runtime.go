package docker

import (
	"github.com/save-abandoned-projects/ignite/pkg/providers"
	dockerruntime "github.com/save-abandoned-projects/ignite/pkg/runtime/docker"
	log "github.com/sirupsen/logrus"
)

func SetDockerRuntime() (err error) {
	log.Trace("Initializing the Docker runtime provider...")
	providers.Runtime, err = dockerruntime.GetDockerClient()
	return
}
