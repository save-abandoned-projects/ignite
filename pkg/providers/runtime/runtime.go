package runtime

import (
	"fmt"

	"github.com/save-abandoned-projects/ignite/pkg/providers"
	containerdprovider "github.com/save-abandoned-projects/ignite/pkg/providers/containerd"
	dockerprovider "github.com/save-abandoned-projects/ignite/pkg/providers/docker"
	"github.com/save-abandoned-projects/ignite/pkg/runtime"
)

func SetRuntime() error {
	switch providers.RuntimeName {
	case runtime.RuntimeDocker:
		return dockerprovider.SetDockerRuntime() // Use the Docker runtime
	case runtime.RuntimeContainerd:
		return containerdprovider.SetContainerdRuntime() // Use the containerd runtime
	}

	return fmt.Errorf("unknown runtime %q", providers.RuntimeName)
}
