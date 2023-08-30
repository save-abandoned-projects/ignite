
/*
	Note: This file is autogenerated! Do not edit it manually!
	Edit client_kernel_template.go instead, and run
	hack/generate-client.sh afterwards.
*/

package client

import (
	"errors"
	"fmt"
	"github.com/save-abandoned-projects/libgitops/pkg/filter"
	"github.com/save-abandoned-projects/libgitops/pkg/runtime"
	"github.com/save-abandoned-projects/libgitops/pkg/storage"
	log "github.com/sirupsen/logrus"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

// KernelClient is an interface for accessing Kernel-specific API objects
type KernelClient interface {
	// New returns a new Kernel
	New() *api.Kernel
	// Get returns the Kernel matching given UID from the storage
	Get(types.UID) (*api.Kernel, error)
	// Set saves the given Kernel into persistent storage
	Set(*api.Kernel) error
	// Patch performs a strategic merge patch on the object with
	// the given UID, using the byte-encoded patch given
	Patch(types.UID, []byte) error
	// Find returns the Kernel matching the given filter, filters can
	// match e.g. the Object's Name, UID or a specific property
	Find(opt filter.ObjectFilter) (*api.Kernel, error)
	// FindAll returns multiple Kernels matching the given filter, filters can
	// match e.g. the Object's Name, UID or a specific property
	FindAll(opts filter.ListOptions) ([]*api.Kernel, error)
	// Delete deletes the Kernel with the given UID from the storage
	Delete(uid types.UID) error
	// List returns a list of all Kernels available
	List() ([]*api.Kernel, error)
}

// Kernels returns the KernelClient for the IgniteInternalClient instance
func (c *IgniteInternalClient) Kernels() KernelClient {
	if c.kernelClient == nil {
		c.kernelClient = newKernelClient(c.storage, c.gv)
	}

	return c.kernelClient
}

// kernelClient is a struct implementing the KernelClient interface
// It uses a shared storage instance passed from the Client together with its own Filterer
type kernelClient struct {
	storage storage.Storage
	gvk     schema.GroupVersionKind
}

// newKernelClient builds the kernelClient struct using the storage implementation and a new Filterer
func newKernelClient(s storage.Storage, gv schema.GroupVersion) KernelClient {
	return &kernelClient{
		storage: s,
		gvk:     gv.WithKind(api.KindKernel.Title()),
	}
}

// New returns a new Object of its kind
func (c *kernelClient) New() *api.Kernel {
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
	return obj.(*api.Kernel)
}

// Find returns a single Kernel based on the given Filter
func (c *kernelClient) Find(opt filter.ObjectFilter) (*api.Kernel, error) {
	log.Tracef("Client.Find; GVK: %v", c.gvk)

	objects, err := c.FindAll(filter.ListOptions{Filters: []filter.ListFilter{filter.ObjectToListFilter(opt)}})
	if err != nil {
		return nil, err
	}

	if len(objects) != 1 {
		return nil, errors.New("ambiguous query: AllFilter used to match single Object")
	}

	return objects[0], nil
}

// FindAll returns multiple Kernels based on the given Filter
func (c *kernelClient) FindAll(opts filter.ListOptions) ([]*api.Kernel, error) {
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
	results := make([]*api.Kernel, 0, len(matches))
	for _, item := range matches {
		results = append(results, item.(*api.Kernel))
	}

	return results, nil
}

// Get returns the Kernel matching given UID from the storage
func (c *kernelClient) Get(uid types.UID) (*api.Kernel, error) {
	log.Tracef("Client.Get; UID: %q, GVK: %v", uid, c.gvk)
	object, err := c.storage.Get(storage.NewObjectKey(storage.NewKindKey(c.gvk), runtime.NewIdentifier(string(uid))))
	if err != nil {
		return nil, err
	}

	return object.(*api.Kernel), nil
}

// Set saves the given Kernel into the persistent storage
func (c *kernelClient) Set(kernel *api.Kernel) error {
	log.Tracef("Client.Set; UID: %q, GVK: %v", kernel.GetUID(), c.gvk)

	return c.storage.Update(kernel)
}

// Patch performs a strategic merge patch on the object with
// the given UID, using the byte-encoded patch given
func (c *kernelClient) Patch(uid types.UID, patch []byte) error {
	return c.storage.Patch(storage.NewObjectKey(storage.NewKindKey(c.gvk), runtime.NewIdentifier(string(uid))), patch)
}

// Delete deletes the Kernel from the storage
func (c *kernelClient) Delete(uid types.UID) error {
	log.Tracef("Client.Delete; UID: %q, GVK: %v", uid, c.gvk)
	return c.storage.Delete(storage.NewObjectKey(storage.NewKindKey(c.gvk), runtime.NewIdentifier(string(uid))))
}

// List returns a list of all Kernels available
func (c *kernelClient) List() ([]*api.Kernel, error) {
	log.Tracef("Client.List; GVK: %v", c.gvk)
	list, err := c.storage.List(storage.NewKindKey(c.gvk))
	if err != nil {
		return nil, err
	}

	results := make([]*api.Kernel, 0, len(list))
	for _, item := range list {
		results = append(results, item.(*api.Kernel))
	}

	return results, nil
}
