
/*
	Note: This file is autogenerated! Do not edit it manually!
	Edit client_vm_template.go instead, and run
	hack/generate-client.sh afterwards.
*/

package client

import (
	"errors"
	"fmt"

	api "github.com/save-abandoned-projects/ignite/pkg/apis/ignite"
	"github.com/save-abandoned-projects/libgitops/pkg/filter"
	"github.com/save-abandoned-projects/libgitops/pkg/runtime"
	"github.com/save-abandoned-projects/libgitops/pkg/storage"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

// VMClient is an interface for accessing VM-specific API objects
type VMClient interface {
	// New returns a new VM
	New() *api.VM
	// Get returns the VM matching given UID from the storage
	Get(types.UID) (*api.VM, error)
	// Set saves the given VM into persistent storage
	Set(*api.VM) error
	// Patch performs a strategic merge patch on the object with
	// the given UID, using the byte-encoded patch given
	Patch(types.UID, []byte) error
	// Find returns the VM matching the given filter, filters can
	// match e.g. the Object's Name, UID or a specific property
	Find(opt filter.ObjectFilter) (*api.VM, error)
	// FindAll returns multiple VMs matching the given filter, filters can
	// match e.g. the Object's Name, UID or a specific property
	FindAll(opts filter.ListOptions) ([]*api.VM, error)
	// Delete deletes the VM with the given UID from the storage
	Delete(uid types.UID) error
	// List returns a list of all VMs available
	List() ([]*api.VM, error)
}

// VMs returns the VMClient for the IgniteInternalClient instance
func (c *IgniteInternalClient) VMs() VMClient {
	if c.vmClient == nil {
		c.vmClient = newVMClient(c.storage, c.gv)
	}

	return c.vmClient
}

// vmClient is a struct implementing the VMClient interface
// It uses a shared storage instance passed from the Client together with its own Filterer
type vmClient struct {
	storage storage.Storage
	gvk     schema.GroupVersionKind
}

// newVMClient builds the vmClient struct using the storage implementation and a new Filterer
func newVMClient(s storage.Storage, gv schema.GroupVersion) VMClient {
	return &vmClient{
		storage: s,
		gvk:     gv.WithKind(api.KindVM.Title()),
	}
}

// New returns a new Object of its kind
func (c *vmClient) New() *api.VM {
	log.Tracef("Client.New; GVK: %v", c.gvk)

	obj, err := c.storage.Serializer().Defaulter().NewDefaultedObject(c.gvk)
	if err != nil {
		panic(fmt.Sprintf("Client.New must not return an error: %v", err))
	}
	// Cast to runtime.Object, and make sure it works
	metaObj, ok := obj.(runtime.Object)
	if !ok {
		panic("can't convert to libgitops.runtime.Object")
	}
	// Set the desired gvk from the caller of this Object
	// In practice, this means, although we created an internal type,
	// from defaulting external TypeMeta information was set. Set the
	// desired gvk here so it's correctly handled in all code that gets
	// the gvk from the Object
	metaObj.GetObjectKind().SetGroupVersionKind(c.gvk)
	return obj.(*api.VM)
}

// Find returns a single VM based on the given Filter
func (c *vmClient) Find(opt filter.ObjectFilter) (*api.VM, error) {
	log.Tracef("Client.Find; GVK: %v", c.gvk)

	objects, err := c.FindAll(filter.ListOptions{Filters: []filter.ListFilter{filter.ObjectToListFilter(opt)}})
	if err != nil {
		return nil, err
	}

	if len(objects) == 0 {
		return nil, nil
	}

	if len(objects) != 1 {
		return nil, errors.New("ambiguous query: AllFilter used to match single Object")
	}

	return objects[0], nil
}

// FindAll returns multiple VMs based on the given Filter
func (c *vmClient) FindAll(opts filter.ListOptions) ([]*api.VM, error) {
	log.Tracef("Client.FindAll; GVK: %v", c.gvk)
	matches, err := c.storage.List(storage.NewKindKey(c.gvk), nil)
	if err != nil {
		return nil, err
	}

	for _, filter := range opts.Filters {
		matches, err = filter.Filter(matches...)
		if err != nil {
			return nil, err
		}
	}
	results := make([]*api.VM, 0, len(matches))
	for _, item := range matches {
		results = append(results, item.(*api.VM))
	}

	return results, nil
}

// Get returns the VM matching given UID from the storage
func (c *vmClient) Get(uid types.UID) (*api.VM, error) {
	log.Tracef("Client.Get; UID: %q, GVK: %v", uid, c.gvk)
	object, err := c.storage.Get(storage.NewObjectKey(storage.NewKindKey(c.gvk), runtime.NewIdentifier(string(uid))))
	if err != nil {
		return nil, err
	}

	return object.(*api.VM), nil
}

// Set saves the given VM into the persistent storage
func (c *vmClient) Set(vm *api.VM) error {
	log.Tracef("Client.Set; UID: %q, GVK: %v", vm.GetUID(), c.gvk)

	return c.storage.Create(vm)
}

// Patch performs a strategic merge patch on the object with
// the given UID, using the byte-encoded patch given
func (c *vmClient) Patch(uid types.UID, patch []byte) error {
	return c.storage.Patch(storage.NewObjectKey(storage.NewKindKey(c.gvk), runtime.NewIdentifier(string(uid))), patch)
}

// Delete deletes the VM from the storage
func (c *vmClient) Delete(uid types.UID) error {
	log.Tracef("Client.Delete; UID: %q, GVK: %v", uid, c.gvk)
	return c.storage.Delete(storage.NewObjectKey(storage.NewKindKey(c.gvk), runtime.NewIdentifier(string(uid))))
}

// List returns a list of all VMs available
func (c *vmClient) List() ([]*api.VM, error) {
	log.Tracef("Client.List; GVK: %v", c.gvk)
	list, err := c.storage.List(storage.NewKindKey(c.gvk))
	if err != nil {
		return nil, err
	}

	results := make([]*api.VM, 0, len(list))
	for _, item := range list {
		results = append(results, item.(*api.VM))
	}

	return results, nil
}
