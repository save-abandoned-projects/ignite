package run

import (
	api "github.com/save-abandoned-projects/ignite/pkg/apis/ignite"
	"github.com/save-abandoned-projects/ignite/pkg/providers"
	"github.com/save-abandoned-projects/libgitops/pkg/filter"
	"k8s.io/apimachinery/pkg/types"
)

func getVMForMatch(vmMatch string) (*api.VM, error) {
	vm, err := providers.Client.VMs().Find(filter.NameFilter{Name: vmMatch})
	// If it can't find the vm by name, use uid instead
	if vm == nil && err == nil {
		return providers.Client.VMs().Find(filter.UIDFilter{UID: types.UID(vmMatch)})
	}
	return vm, err
}

func getVMsForMatches(vmMatches []string) ([]*api.VM, error) {
	allVMs := make([]*api.VM, 0, len(vmMatches))
	for _, match := range vmMatches {
		vm, err := getVMForMatch(match)
		if err != nil {
			return nil, err
		}
		if vm != nil {
			allVMs = append(allVMs, vm)
		}
	}
	return allVMs, nil
}

func getAllVMs() ([]*api.VM, error) {
	return providers.Client.VMs().FindAll(nil)
}
