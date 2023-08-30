package lookup

import (
	api "github.com/save-abandoned-projects/ignite/pkg/apis/ignite"
	"github.com/save-abandoned-projects/ignite/pkg/client"
	"github.com/save-abandoned-projects/libgitops/pkg/filter"
	"k8s.io/apimachinery/pkg/types"
)

func ImageUIDForVM(vm *api.VM, c *client.Client) (types.UID, error) {
	image, err := c.Images().Find(filter.NameFilter{Name: vm.Spec.Image.OCI.String()})
	if err != nil {
		return "", err
	}

	return image.GetUID(), nil
}

func KernelUIDForVM(vm *api.VM, c *client.Client) (types.UID, error) {
	kernel, err := c.Kernels().Find(filter.NameFilter{Name: vm.Spec.Kernel.OCI.String()})
	if err != nil {
		return "", err
	}

	return kernel.GetUID(), nil
}
