package containerd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	meta "github.com/save-abandoned-projects/ignite/pkg/apis/meta/v1alpha1"
	"github.com/save-abandoned-projects/ignite/pkg/constants"
	"github.com/save-abandoned-projects/ignite/pkg/runtime"

	v2shim "github.com/containerd/containerd/runtime/v2/shim"
	"gotest.tools/assert"
)

var client runtime.Interface

func init() {
	var clienterr error
	client, clienterr = GetContainerdClient()
	if clienterr != nil {
		panic(clienterr)
	}
}

var imageName, _ = meta.NewOCIImageRef("docker.io/library/busybox:latest")

func TestPull(t *testing.T) {
	err := client.PullImage(imageName)
	if err != nil {
		t.Errorf("Error Pulling image: %s", err)
	}
}

func TestInspect(t *testing.T) {
	result, err := client.InspectImage(imageName)
	t.Log(result)
	if err != nil {
		t.Errorf("Error Inspecting image: %s", err)
	}
}

/*func TestExport(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	fmt.Println(tempDir)

	tarCmd := exec.Command("tar", "-x", "-C", tempDir)
	reader, _, err = client.ExportImage(imageName)
	if err != nil {
		t.Fatal("export err:", err)
	}

	tarCmd.Stdin = reader
	if err := tarCmd.Start(); err != nil {
		t.Fatal(err)
	}

	if err := tarCmd.Wait(); err != nil {
		t.Fatal(err)
	}

	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
	t.Log("done", tempDir)
}*/

func TestRunRemove(t *testing.T) {
	cName := "ignite-test-foo2"
	cID := "test-foo2"
	vmDir := filepath.Join(constants.VM_DIR, cID)

	// TODO: refactor client RunContainer() to take in generic stateDir
	//       remove dependency on VM constants for runtime client
	//       this specific dir is currently required to support resolvconf
	//       ideally, we could pass any tempdir with any permissions here
	assert.NilError(t, os.MkdirAll(vmDir, constants.DATA_DIR_PERM))

	cmds := []string{"/bin/sh", "-c", "echo hello"}
	cfg := getContainerConfig(cmds, vmDir)

	taskID, err := client.RunContainer(imageName, cfg, cName, cID)
	if err != nil {
		t.Errorf("Error Running Container /w TaskID %q: %s", taskID, err)
	} else {
		t.Logf("TaskID: %q", taskID)
	}

	// TODO: this works around a race where the task is not yet stopped
	//       do this better -- wait on taskID returned above?
	time.Sleep(time.Second / 4)

	err = client.RemoveContainer(cName)
	if err != nil {
		t.Errorf("Error Removing Container: %s", err)
	}

	// just in case the process is hung -- cleanup
	containerCleanup(client, cName)
}

func TestInspectContainer(t *testing.T) {
	cName := "ignite-test-foo3"
	cID := "test-foo3"
	vmDir := filepath.Join(constants.VM_DIR, cID)
	assert.NilError(t, os.MkdirAll(vmDir, constants.DATA_DIR_PERM))

	cmds := []string{"/bin/sh", "-c", "sleep 20"}
	cfg := getContainerConfig(cmds, vmDir)

	taskID, err := client.RunContainer(imageName, cfg, cName, cID)
	if err != nil {
		t.Errorf("Error Running Container /w TaskID %q: %s", taskID, err)
	} else {
		t.Logf("TaskID: %q", taskID)
	}

	// Run inspect and check the result.
	result, err := client.InspectContainer(cName)
	assert.NilError(t, err)
	assert.Equal(t, taskID, result.ID)
	// Returns image URI - docker.io/library/busybox:latest.
	assert.Check(t, strings.Contains(result.Image, imageName.String()))
	assert.Equal(t, "running", result.Status)

	time.Sleep(time.Second / 4)
	client.KillContainer(cName, "")
	err = client.RemoveContainer(cName)
	if err != nil {
		t.Errorf("Error Removing Container: %s", err)
	}

	// just in case the process is hung -- cleanup
	containerCleanup(client, cName)
}

// containerCleanup ensures that the container is cleaned up.
func containerCleanup(client runtime.Interface, cName string) {
	// just in case the process is hung -- cleanup
	client.KillContainer(cName, "SIGQUIT") //nolint:errcheck // TODO: common constant for SIGQUIT
	client.RemoveContainer(cName)          //nolint:errcheck
}

// getContainerConfig returns a container config.
func getContainerConfig(cmds []string, vmDir string) *runtime.ContainerConfig {
	return &runtime.ContainerConfig{
		Cmd: cmds,
		Binds: []*runtime.Bind{
			runtime.BindBoth(vmDir),
		},
		Devices: []*runtime.Bind{
			runtime.BindBoth("/dev/kvm"),
		},
		Labels: map[string]string{},
	}
}

func TestV2ShimRuntimesHaveBinaryNames(t *testing.T) {
	for _, runtime := range v2ShimRuntimes {
		if v2shim.BinaryName(runtime) == "" {
			t.Errorf("shim binary could not be found -- %q is an invalid runtime/v2/shim", runtime)
		}
	}
}

func TestNewRemoteResolver(t *testing.T) {
	// Use a template for the configuration and get a registry configuration
	// with appropriate protocol.
	templateConfig := `
{
	"auths": {
		"%s://127.5.0.1:5443": {
			"auth": "aHR0cHNfdGVzdHVzZXI6aHR0cHNfdGVzdHBhc3N3b3Jk"
		}
	}
}
`
	getRegistryConfigWithProtocol := func(protocol string) string {
		return fmt.Sprintf(templateConfig, protocol)
	}

	domainRef := "127.5.0.1:5443"

	cases := []struct {
		name               string
		insecureRegistries []string
		registryConfig     string
		wantErr            bool
	}{
		{
			name: "invalid configuration",
			registryConfig: `
{ some invalid json }
`,
			wantErr: true,
		},
		{
			name:           "valid configuration",
			registryConfig: getRegistryConfigWithProtocol("https"),
		},
		{
			name:           "http server address without insecure registries",
			registryConfig: getRegistryConfigWithProtocol("http"),
			wantErr:        true,
		},
		{
			name:               "http server address with insecure registries",
			insecureRegistries: []string{"127.5.0.1:5443"},
			registryConfig:     getRegistryConfigWithProtocol("http"),
		},
	}

	for _, rt := range cases {
		t.Run(rt.name, func(t *testing.T) {
			// Create directory for the registry configuration.
			dir, err := os.MkdirTemp("", "ignite")
			if err != nil {
				t.Fatalf("failed to create storage for ignite: %v", err)
			}
			defer os.RemoveAll(dir)

			// If a registry configuration content is given, write it.
			if len(rt.registryConfig) > 0 {
				configPath := filepath.Join(dir, "config.json")
				writeErr := os.WriteFile(configPath, []byte(rt.registryConfig), 0600)
				assert.NilError(t, writeErr)
				defer os.Remove(configPath)
			}

			// If insecure registries are given, set env vars.
			if len(rt.insecureRegistries) > 0 {
				irValues := strings.Join(rt.insecureRegistries, ",")
				os.Setenv(InsecureRegistriesEnvVar, irValues)
				defer os.Unsetenv(InsecureRegistriesEnvVar)
			}

			_, rrErr := newRemoteResolver(domainRef, dir)
			if (rrErr != nil) != rt.wantErr {
				t.Errorf("expected error %t, actual: %v", rt.wantErr, rrErr)
			}
		})
	}
}
