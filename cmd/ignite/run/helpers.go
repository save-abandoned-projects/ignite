package run

import (
	api "github.com/save-abandoned-projects/ignite/pkg/apis/ignite"
	"github.com/save-abandoned-projects/ignite/pkg/providers"
	"github.com/save-abandoned-projects/libgitops/pkg/filter"
)

// TODO: This
func getVMForMatch(vmMatch string) (*api.VM, error) {
	return providers.Client.VMs().Find(filter.NameFilter{Name: vmMatch})
}

// TODO: This
func getVMsForMatches(vmMatches []string) ([]*api.VM, error) {
	allVMs := make([]*api.VM, 0, len(vmMatches))
	for _, match := range vmMatches {
		vm, err := getVMForMatch(match)
		if err != nil {
			return nil, err
		}
		allVMs = append(allVMs, vm)
	}
	return allVMs, nil
}

func getAllVMs() ([]*api.VM, error) {
	return providers.Client.VMs().FindAll(filter.ListOptions{})
}
