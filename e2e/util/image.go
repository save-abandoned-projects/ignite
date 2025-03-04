package util

import (
	"os/exec"

	"github.com/save-abandoned-projects/ignite/pkg/runtime"
	"github.com/save-abandoned-projects/ignite/pkg/runtime/containerd"
)

// RmiDocker removes an image from docker content store.
func RmiDocker(img string) {
	_, _ = exec.Command(
		"docker",
		"rmi", img,
	).CombinedOutput()
}

// RmiContainerd removes an image from containerd content store.
func RmiContainerd(img string) {
	socketAddr, _ := containerd.StatContainerdSocket()
	_, _ = exec.Command(
		"ctr", "-a", socketAddr,
		"-n", "firecracker",
		"image", "rm", img,
	).CombinedOutput()
}

// rmiCompletely removes a given image completely, from ignite image store and
// runtime image store.
func RmiCompletely(img string, cmd *Command, rt runtime.Name) {
	// Remote from ignite content store.
	_, _ = cmd.New().
		With("image", "rm", img).
		Cmd.CombinedOutput()

	// Remove from runtime content store.
	switch rt {
	case runtime.RuntimeContainerd:
		RmiContainerd(img)
	case runtime.RuntimeDocker:
		RmiDocker(img)
	}
}
