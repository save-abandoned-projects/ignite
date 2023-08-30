package containerd

import (
	"github.com/save-abandoned-projects/ignite/pkg/providers"
	containerdruntime "github.com/save-abandoned-projects/ignite/pkg/runtime/containerd"
	log "github.com/sirupsen/logrus"
)

func SetContainerdRuntime() (err error) {
	log.Trace("Initializing the containerd runtime provider...")
	providers.Runtime, err = containerdruntime.GetContainerdClient()
	return
}
