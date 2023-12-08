package operations

import (
	"fmt"
	"os"

	api "github.com/save-abandoned-projects/ignite/pkg/apis/ignite"
	meta "github.com/save-abandoned-projects/ignite/pkg/apis/meta/v1alpha1"
	"github.com/save-abandoned-projects/ignite/pkg/client"
	"github.com/save-abandoned-projects/ignite/pkg/dmlegacy"
	"github.com/save-abandoned-projects/ignite/pkg/logs"
	"github.com/save-abandoned-projects/ignite/pkg/providers"
	"github.com/save-abandoned-projects/ignite/pkg/runtime"
	log "github.com/sirupsen/logrus"
)

const (
	signalSIGQUIT = "SIGQUIT"
)

// DeleteVM removes the specified VM from the Client and performs a cleanup
func DeleteVM(c *client.Client, vm *api.VM) error {
	if err := CleanupVM(vm); err != nil {
		return err
	}

	return c.VMs().Delete(vm.GetUID())
}

// CleanupVM removes the resources of the given VM
func CleanupVM(vm *api.VM) error {
	// Runtime information is available only when the VM is running.
	if vm.Running() {
		// Inspect the container before trying to stop it, and it gets auto-removed
		inspectResult, _ := providers.Runtime.InspectContainer(vm.PrefixedID())

		// If the VM is running, try to kill it first, so we don't leave dangling containers. Otherwise, try to cleanup VM networking.
		if err := StopVM(vm, true, true); err != nil {
			if vm.Running() {
				return err
			}
		}

		// Remove the VM container if it exists
		RemoveVMContainer(inspectResult)
	}

	// if vm is not running, just remove the container
	providers.Runtime.RemoveContainer(vm.PrefixedID())

	// After removing the VM container, if the Snapshot Device is still there, clean up
	if _, err := os.Stat(vm.SnapshotDev()); err == nil {
		// try remove it again with DeactivateSnapshot
		if err := dmlegacy.DeactivateSnapshot(vm); err != nil {
			return err
		}
	}

	if logs.Quiet {
		fmt.Println(vm.GetUID())
	} else {
		log.Infof("Removed %s with name %q and ID %q", vm.GetObjectKind().GroupVersionKind().Kind, vm.GetName(), vm.GetUID())
	}

	return nil
}

func RemoveVMContainer(result *runtime.ContainerInspectResult) {
	if result == nil {
		return // If given no result, don't attempt removal
	}

	// Remove the VM container. If the container has been/is being automatically removed
	// between InspectContainer and this call, RemoveContainer will throw an error. Currently
	// we have no real way to inspect and remove immediately without having a potential race
	// condition, so ignore the error for now. TODO: Robust conditional remove support
	providers.Runtime.RemoveContainer(result.ID)
}

// StopVM removes networking of the given VM and stops or kills it
func StopVM(vm *api.VM, kill, silent bool) error {
	var err error
	container := vm.PrefixedID()
	action := "stop"

	if !vm.Running() && !logs.Quiet {
		log.Warnf("VM %q is not running but trying to cleanup networking for stopped container\n", vm.GetUID())
	}

	// Remove VM networking
	if err = removeNetworking(vm.Status.Runtime.ID, vm.Spec.Network.Ports...); err != nil {
		log.Warnf("Failed to cleanup networking for stopped container %s %q: %v", vm.Kind, vm.GetUID(), err)

		return err
	}

	if vm.Running() {
		// Stop or kill the VM container
		if kill {
			action = "kill"
			err = providers.Runtime.KillContainer(container, signalSIGQUIT) // TODO: common constant for SIGQUIT
		} else {
			err = providers.Runtime.StopContainer(container, nil)
		}

		if err != nil {
			return fmt.Errorf("failed to %s container for %s %q: %v", action, vm.Kind, vm.GetUID(), err)
		}

		vm.Status.Running = false
		vm.Status.Network.IPAddresses = nil
		err := providers.Client.VMs().Update(vm)
		if err != nil {
			log.Errorf("Update status %s with name %q and ID %q  failed!", vm.Kind, vm.GetName(), vm.GetUID())
		}

		if silent {
			return nil
		}

		if logs.Quiet {
			fmt.Println(vm.GetUID())
		} else {
			log.Infof("Stopped %s with name %q and ID %q", vm.Kind, vm.GetName(), vm.GetUID())
		}
	}

	return nil
}

func removeNetworking(containerID string, portmappings ...meta.PortMapping) error {
	log.Infof("Removing the container with ID %q from the %q network", containerID, providers.NetworkPlugin.Name())
	return providers.NetworkPlugin.RemoveContainerNetwork(containerID, portmappings...)
}
