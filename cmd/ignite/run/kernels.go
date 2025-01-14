package run

import (
	"os"

	api "github.com/save-abandoned-projects/ignite/pkg/apis/ignite"
	"github.com/save-abandoned-projects/ignite/pkg/providers"
	"github.com/save-abandoned-projects/ignite/pkg/util"
)

type KernelsOptions struct {
	allKernels []*api.Kernel
}

func NewKernelsOptions() (ko *KernelsOptions, err error) {
	ko = &KernelsOptions{}
	ko.allKernels, err = providers.Client.Kernels().FindAll(nil)
	// If the storage is uninitialized, avoid failure and continue with empty
	// kernel list.
	if err != nil && os.IsNotExist(err) {
		err = nil
	}
	return
}

func Kernels(ko *KernelsOptions) error {
	o := util.NewOutput()
	defer o.Flush()

	o.Write("KERNEL ID", "NAME", "CREATED", "SIZE", "VERSION")
	for _, kernel := range ko.allKernels {
		o.Write(kernel.GetUID(), kernel.GetName(), kernel.GetCreationTimestamp(), kernel.Status.OCISource.Size.String(), kernel.Status.Version)
	}

	return nil
}
