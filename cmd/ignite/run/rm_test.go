package run

import (
	"github.com/save-abandoned-projects/libgitops/pkg/serializer"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"path/filepath"
	"testing"

	"github.com/save-abandoned-projects/libgitops/pkg/runtime"
	"github.com/save-abandoned-projects/libgitops/pkg/storage"
	"github.com/save-abandoned-projects/libgitops/pkg/storage/cache"

	api "github.com/save-abandoned-projects/ignite/pkg/apis/ignite"
	"github.com/save-abandoned-projects/ignite/pkg/apis/ignite/scheme"
	meta "github.com/save-abandoned-projects/ignite/pkg/apis/meta/v1alpha1"
	"github.com/save-abandoned-projects/ignite/pkg/client"
	"github.com/save-abandoned-projects/ignite/pkg/providers"
	"github.com/save-abandoned-projects/ignite/pkg/util"
)

func TestNewRmOptions(t *testing.T) {
	testdataDir := "testdata"

	cases := []struct {
		name        string
		existingVMs []string
		rmFlags     *RmFlags
		vmMatches   []string // argument of NewRmOptions()
		wantMatches []string
		err         bool
	}{
		{
			name:        "rm with vm arg",
			existingVMs: []string{"myvm1", "myvm2", "myvm3"},
			rmFlags:     &RmFlags{},
			vmMatches:   []string{"myvm2"},
			wantMatches: []string{"myvm2"},
		},
		{
			name:        "rm with multiple vm args",
			existingVMs: []string{"myvm1", "myvm2", "myvm3"},
			rmFlags:     &RmFlags{},
			vmMatches:   []string{"myvm2", "myvm3"},
			wantMatches: []string{"myvm2", "myvm3"},
		},
		{
			name:        "error rm non-existing vm",
			existingVMs: []string{"myvm1", "myvm2", "myvm3"},
			rmFlags:     &RmFlags{},
			vmMatches:   []string{"myvm4"},
			err:         false,
		},
		{
			name:        "error rm without any args or config flag",
			existingVMs: []string{"myvm1", "myvm2", "myvm3"},
			rmFlags:     &RmFlags{},
			err:         true,
		},
		{
			name:        "error rm with vm arg and config flag",
			existingVMs: []string{"myvm1"},
			rmFlags:     &RmFlags{ConfigFile: "foo.yaml"},
			vmMatches:   []string{"myvm1"},
			err:         true,
		},
		{
			name:        "rm with config file",
			existingVMs: []string{"myvm1", "myvm2", "myvm3"},
			rmFlags:     &RmFlags{ConfigFile: filepath.Join(testdataDir, "input/rm-vm1.yaml")},
			wantMatches: []string{"myvm2"},
		},
		{
			name:        "error rm config name and uid missing",
			existingVMs: []string{"myvm1"},
			rmFlags:     &RmFlags{ConfigFile: filepath.Join(testdataDir, "input/rm-no-name-uid.yaml")},
			err:         true,
		},
	}

	for _, rt := range cases {
		t.Run(rt.name, func(t *testing.T) {
			// Create storage.
			dir, err := os.MkdirTemp("", "ignite")
			if err != nil {
				t.Fatalf("failed to create storage for ignite: %v", err)
			}
			defer os.RemoveAll(dir)

			storage := cache.NewCache(storage.NewGenericStorage(
				storage.NewGenericRawStorage(dir, api.SchemeGroupVersion, serializer.ContentTypeYAML),
				scheme.Serializer,
				[]runtime.IdentifierFactory{runtime.Metav1NameIdentifier, runtime.ObjectUIDIdentifier}))

			// Create ignite client with the created storage.
			ic := client.NewClient(storage)

			// Create the existing VMs.
			for _, objectName := range rt.existingVMs {
				vm := ic.VMs().New()
				vm.SetName(objectName)

				// Set UID.
				uid, err := util.NewUID()
				if err != nil {
					t.Errorf("failed to generate new UID: %v", err)
				}
				vm.SetUID(types.UID(uid))

				// Set VM image.
				ociRef, err := meta.NewOCIImageRef("foo/bar:latest")
				if err != nil {
					t.Errorf("failed to create new image reference: %v", err)
				}
				img := &api.Image{
					Spec: api.ImageSpec{
						OCI: ociRef,
					},
				}
				vm.SetImage(img)

				// Set Kernel image.
				ociRefKernel, err := meta.NewOCIImageRef("foo/bar:latest")
				if err != nil {
					t.Errorf("failed to create new image reference: %v", err)
				}
				kernel := &api.Kernel{
					Spec: api.KernelSpec{
						OCI: ociRefKernel,
					},
				}
				vm.SetKernel(kernel)

				// Set sandbox image without a helper.
				ociRefSandbox, err := meta.NewOCIImageRef("foo/bar:latest")
				if err != nil {
					t.Errorf("failed to create new image reference: %v", err)
				}
				vm.Spec.Sandbox.OCI = ociRefSandbox

				// Save object.
				if err := ic.VMs().Set(vm); err != nil {
					t.Errorf("failed to store VM object: %v", err)
				}
			}

			// Set provider client used in remove to find VM matches.
			providers.Client = ic

			// Create new rm options using the rmFlags and vmMatches.
			ro, err := rt.rmFlags.NewRmOptions(rt.vmMatches)
			if (err != nil) != rt.err {
				t.Fatalf("expected error %t, actual: %v", rt.err, err)
			}

			// Check if the wanted VMs are in the matched VMs list.
			for _, wantVM := range rt.wantMatches {
				found := false
				for _, vm := range ro.vms {
					if vm.Name == wantVM {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected vm %q to be in remove vm list", wantVM)
				}
			}
		})
	}
}
