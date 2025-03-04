package metadata

import (
	"github.com/save-abandoned-projects/libgitops/pkg/serializer"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"testing"

	"github.com/save-abandoned-projects/libgitops/pkg/runtime"
	"github.com/save-abandoned-projects/libgitops/pkg/storage"
	"github.com/save-abandoned-projects/libgitops/pkg/storage/cache"

	api "github.com/save-abandoned-projects/ignite/pkg/apis/ignite"
	"github.com/save-abandoned-projects/ignite/pkg/apis/ignite/scheme"
	meta "github.com/save-abandoned-projects/ignite/pkg/apis/meta/v1alpha1"
	"github.com/save-abandoned-projects/ignite/pkg/client"
	"github.com/save-abandoned-projects/ignite/pkg/util"
)

func TestSetLabels(t *testing.T) {
	cases := []struct {
		name       string
		obj        runtime.Object
		labels     []string
		wantLabels map[string]string
		err        bool
	}{
		{
			name: "nil runtime object",
			obj:  nil,
			err:  true,
		},
		{
			name: "valid labels",
			obj:  &api.VM{},
			labels: []string{
				"label1=value1",
				"label2=value2",
				"label3=",
				"label4=value4,label5=value5",
			},
			wantLabels: map[string]string{
				"label1": "value1",
				"label2": "value2",
				"label3": "",
				"label4": "value4,label5=value5",
			},
		},
		{
			name:   "invalid label - key without value",
			obj:    &api.VM{},
			labels: []string{"key1"},
			err:    true,
		},
		{
			name:   "invalid label - empty name",
			obj:    &api.VM{},
			labels: []string{"="},
			err:    true,
		},
	}

	for _, rt := range cases {
		t.Run(rt.name, func(t *testing.T) {
			err := SetLabels(rt.obj, rt.labels)
			if (err != nil) != rt.err {
				t.Errorf("expected error %t, actual: %v", rt.err, err)
			}

			if rt.obj != nil {
				havaLabels := rt.obj.GetLabels()
				// Check the values of all the labels.
				for k, v := range rt.wantLabels {
					if haveV, ok := havaLabels[k]; !ok || haveV != v {
						t.Errorf("expected label key %q to have value %q, actual: %q", k, v, haveV)
					}
				}
			}
		})
	}
}

func TestVerifyUIDOrName(t *testing.T) {
	cases := []struct {
		name            string
		existingObjects []string
		newObject       string
		err             bool
	}{
		{
			name:            "create object with similar names",
			existingObjects: []string{"myvm1", "myvm11", "myvm111"},
			newObject:       "myvm",
		},
		{
			name:            "create object with existing names",
			existingObjects: []string{"myvm1", "myvm2"},
			newObject:       "myvm1",
			err:             true,
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

			storage := cache.NewCache(
				storage.NewGenericStorage(
					storage.NewGenericRawStorage(dir, api.SchemeGroupVersion, serializer.ContentTypeYAML),
					scheme.Serializer,
					[]runtime.IdentifierFactory{runtime.Metav1NameIdentifier, runtime.ObjectUIDIdentifier}))

			// Create ignite client with the created storage.
			ic := client.NewClient(storage)

			// Create existing VM object.
			objectKind := "VM"
			for _, objectName := range rt.existingObjects {
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

			// Check if new object name exists.
			err = verifyUIDOrName(ic, rt.newObject, api.Kind(objectKind))
			if (err != nil) != rt.err {
				t.Errorf("expected error %t, actual: %v", rt.err, err)
			}
		})
	}
}
