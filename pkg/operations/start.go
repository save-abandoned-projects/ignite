package operations

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	api "github.com/save-abandoned-projects/ignite/pkg/apis/ignite"
	meta "github.com/save-abandoned-projects/ignite/pkg/apis/meta/v1alpha1"
	"github.com/save-abandoned-projects/ignite/pkg/constants"
	"github.com/save-abandoned-projects/ignite/pkg/dmlegacy"
	"github.com/save-abandoned-projects/ignite/pkg/logs"
	"github.com/save-abandoned-projects/ignite/pkg/operations/lookup"
	"github.com/save-abandoned-projects/ignite/pkg/providers"
	"github.com/save-abandoned-projects/ignite/pkg/runtime"
	"github.com/save-abandoned-projects/ignite/pkg/util"
	log "github.com/sirupsen/logrus"
)

// VMChannels can be used to get signals for different stages of VM lifecycle
type VMChannels struct {
	SpawnFinished chan error
}

func StartVM(vm *api.VM, debug bool) error {

	vmChans, err := StartVMNonBlocking(vm, debug)
	if err != nil {
		return err
	}

	if err := <-vmChans.SpawnFinished; err != nil {
		return err
	}

	return nil
}

func StartVMNonBlocking(vm *api.VM, debug bool) (*VMChannels, error) {
	// Inspect the VM container and remove it if it exists
	inspectResult, _ := providers.Runtime.InspectContainer(vm.PrefixedID())
	RemoveVMContainer(inspectResult)

	// Make sure we always initialize all channels
	vmChans := &VMChannels{
		SpawnFinished: make(chan error),
	}

	// Setup the snapshot overlay filesystem
	snapshotDevPath, err := dmlegacy.ActivateSnapshot(vm)
	if err != nil {
		return vmChans, err
	}

	kernelUID, err := lookup.KernelUIDForVM(vm, providers.Client)
	if err != nil {
		return vmChans, err
	}

	vmDir := filepath.Join(constants.VM_DIR, string(vm.GetUID()))
	kernelDir := filepath.Join(constants.KERNEL_DIR, string(kernelUID))

	// Verify that the image containing ignite-spawn is pulled
	// TODO: Integrate automatic pulling into pkg/runtime
	if err := verifyPulled(vm.Spec.Sandbox.OCI); err != nil {
		return vmChans, err
	}

	config := &runtime.ContainerConfig{
		Cmd: []string{
			fmt.Sprintf("--log-level=%s", logs.Logger.Level.String()),
			string(vm.GetUID()),
		},
		Labels: map[string]string{"ignite.name": vm.GetName()},
		Binds: []*runtime.Bind{
			{
				HostPath:      vmDir,
				ContainerPath: vmDir,
			},
			{
				// Mount the metadata.json file specifically into the container, to a well-known place for ignite-spawn to access
				HostPath:      path.Join(vmDir, constants.METADATA),
				ContainerPath: constants.IGNITE_SPAWN_VM_FILE_PATH,
			},
			{
				// Mount the vmlinux file specifically into the container, to a well-known place for ignite-spawn to access
				HostPath:      path.Join(kernelDir, constants.KERNEL_FILE),
				ContainerPath: constants.IGNITE_SPAWN_VMLINUX_FILE_PATH,
			},
		},
		CapAdds: []string{
			"SYS_ADMIN", // Needed to run "dmsetup remove" inside the container
			"NET_ADMIN", // Needed for removing the IP from the container's interface
		},
		Devices: []*runtime.Bind{
			runtime.BindBoth("/dev/mapper/control"), // This enables containerized Ignite to remove its own dm snapshot
			runtime.BindBoth("/dev/net/tun"),        // Needed for creating TAP adapters
			runtime.BindBoth("/dev/kvm"),            // Pass through virtualization support
			runtime.BindBoth(snapshotDevPath),       // The block device to boot from
		},
		StopTimeout:  constants.STOP_TIMEOUT + constants.IGNITE_TIMEOUT,
		PortBindings: vm.Spec.Network.Ports, // Add the port mappings to Docker
	}

	var envVars []string
	for k, v := range vm.GetObjectMeta().GetAnnotations() {
		if strings.HasPrefix(k, constants.IGNITE_SANDBOX_ENV_VAR) {
			k := strings.TrimPrefix(k, constants.IGNITE_SANDBOX_ENV_VAR)
			envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
		}
	}
	config.EnvVars = envVars

	// Add the volumes to the container devices
	for _, volume := range vm.Spec.Storage.Volumes {
		if volume.BlockDevice == nil {
			continue // Skip all non block device volumes for now
		}

		config.Devices = append(config.Devices, &runtime.Bind{
			HostPath:      volume.BlockDevice.Path,
			ContainerPath: path.Join(constants.IGNITE_SPAWN_VOLUME_DIR, volume.Name),
		})
	}

	// Prepare the networking for the container, for the given network plugin
	if err := providers.NetworkPlugin.PrepareContainerSpec(config); err != nil {
		return vmChans, err
	}

	// If we're not debugging, remove the container post-run
	if !debug {
		config.AutoRemove = true
	}

	// Run the VM container in Docker
	containerID, err := providers.Runtime.RunContainer(vm.Spec.Sandbox.OCI, config, vm.PrefixedID(), string(vm.GetUID()))
	if err != nil {
		return vmChans, fmt.Errorf("failed to start container for VM %q: %v", vm.GetUID(), err)
	}

	// Set the container ID for the VM
	vm.Status.Runtime.ID = containerID
	vm.Status.Runtime.Name = providers.RuntimeName

	// Set up the networking
	result, err := providers.NetworkPlugin.SetupContainerNetwork(containerID, vm.Spec.Network.Ports...)
	if err != nil {
		return vmChans, err
	}

	if !logs.Quiet {
		log.Infof("Networking is handled by %q", providers.NetworkPlugin.Name())
		log.Infof("Started Firecracker VM %q in a container with ID %q", vm.GetUID(), containerID)
	}

	// Append non-loopback runtime IP addresses of the VM to its state
	for _, addr := range result.Addresses {
		if !addr.IP.IsLoopback() {
			vm.Status.Network.IPAddresses = append(vm.Status.Network.IPAddresses, addr.IP)
		}
	}
	vm.Status.Network.Plugin = providers.NetworkPluginName

	// write the API object in a non-running state before we wait for spawn's network logic and firecracker
	if err := providers.Client.VMs().Update(vm); err != nil {
		return vmChans, err
	}

	// TODO: This is temporary until we have proper communication to the container
	// It's best to perform any imperative changes to the VM object pointer before this go-routine starts
	go waitForSpawn(vm, vmChans)

	return vmChans, nil
}

// verifyPulled pulls the ignite-spawn image if it's not present
func verifyPulled(image meta.OCIImageRef) error {
	if _, err := providers.Runtime.InspectImage(image); err != nil {
		log.Infof("Pulling image %q...", image)
		if err = providers.Runtime.PullImage(image); err != nil {
			return err
		}

		// Verify the image was pulled
		if _, err = providers.Runtime.InspectImage(image); err != nil {
			return err
		}
	}

	return nil
}

// TODO: This check for the Prometheus socket file is temporary
// until we get a proper ignite <-> ignite-spawn communication channel
func waitForSpawn(vm *api.VM, vmChans *VMChannels) {
	const checkInterval = 100 * time.Millisecond

	timer := time.Now()
	for time.Since(timer) < constants.IGNITE_SPAWN_TIMEOUT {
		time.Sleep(checkInterval)

		if util.FileExists(path.Join(vm.ObjectPath(), constants.PROMETHEUS_SOCKET)) {
			// Before we write the VM, we should REALLY REALLY re-fetch the API object from storage
			vm, err := providers.Client.VMs().Get(vm.GetUID())
			if err != nil {
				vmChans.SpawnFinished <- err
			}

			// Set the VM's status to running
			vm.Status.Running = true

			// Set the start time for the VM
			startTime := metav1.Time{Time: time.Now().UTC()}
			vm.Status.StartTime = startTime

			// Write the state changes, send any errors through the channel
			vmChans.SpawnFinished <- providers.Client.VMs().Update(vm)
			return
		}
	}

	vmChans.SpawnFinished <- fmt.Errorf("timeout waiting for ignite-spawn startup")
}
