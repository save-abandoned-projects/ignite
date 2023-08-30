package client

import (
	"errors"
	"fmt"
	api "github.com/save-abandoned-projects/ignite/pkg/apis/ignite"

	"github.com/save-abandoned-projects/libgitops/pkg/filter"
	"github.com/save-abandoned-projects/libgitops/pkg/runtime"
	"github.com/save-abandoned-projects/libgitops/pkg/storage"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

// DynamicClient is an interface for accessing API types generically
type DynamicClient interface {
	// New returns a new Object of its kind
	New() runtime.Object
	// Get returns an Object matching the UID from the storage
	Get(types.UID) (runtime.Object, error)
	// Set saves an Object into the persistent storage
	Set(runtime.Object) error
	// Patch performs a strategic merge patch on the object with
	// the given UID, using the byte-encoded patch given
	Patch(types.UID, []byte) error
	// Find returns an Object based on the given filter, filters can
	// match e.g. the Object's Name, UID or a specific property
	Find(filter filter.ListOption) (runtime.Object, error)
	// FindAll returns multiple Objects based on the given filter, filters can
	// match e.g. the Object's Name, UID or a specific property
	FindAll(filter filter.ListOption) ([]runtime.Object, error)
	// Delete deletes an Object from the storage
	Delete(uid types.UID) error
	// List returns a list of all Objects available
	List() ([]runtime.Object, error)
}

// Dynamic returns the DynamicClient for the Client instance, for the specific kind
func (c *IgniteInternalClient) Dynamic(kind api.Kind) (dc DynamicClient) {
	var ok bool
	gvk := c.gv.WithKind(kind.Title())
	if dc, ok = c.dynamicClients[gvk]; !ok {
		dc = newDynamicClient(c.storage, gvk)
		c.dynamicClients[gvk] = dc
	}

	return
}

// dynamicClient is a struct implementing the DynamicClient interface
// It uses a shared storage instance passed from the Client together with its own Filterer
type dynamicClient struct {
	storage storage.Storage
	gvk     schema.GroupVersionKind
}

// newDynamicClient builds the dynamicClient struct using the storage implementation and a new Filterer
func newDynamicClient(s storage.Storage, gvk schema.GroupVersionKind) DynamicClient {
	return &dynamicClient{
		storage: s,
		gvk:     gvk,
	}
}

// New returns a new Object of its kind
func (c *dynamicClient) New() runtime.Object {
	obj, err := c.storage.Serializer().Defaulter().NewDefaultedObject(c.gvk)
	if err != nil {
		panic(fmt.Sprintf("Client.New must not return an error: %v", err))
	}
	metaObj, ok := obj.(runtime.Object)
	if !ok {
		panic("can't convert to libgitops.runtime.Object")
	}

	return metaObj
}

// Get returns an Object based the given UID
func (c *dynamicClient) Get(uid types.UID) (runtime.Object, error) {
	return c.storage.Get(storage.NewObjectKey(storage.NewKindKey(c.gvk), runtime.NewIdentifier(string(uid))))
}

// Set saves an Object into the persistent storage
func (c *dynamicClient) Set(resource runtime.Object) error {
	return c.storage.Update(resource)
}

// Patch performs a strategic merge patch on the object with
// the given UID, using the byte-encoded patch given
func (c *dynamicClient) Patch(uid types.UID, patch []byte) error {
	return c.storage.Patch(storage.NewObjectKey(storage.NewKindKey(c.gvk), runtime.NewIdentifier(string(uid))), patch)
}

// Find returns an Object based on a given Filter
func (c *dynamicClient) Find(opt filter.ListOption) (runtime.Object, error) {
	objects, err := c.FindAll(opt)
	if err != nil {
		return nil, err
	}

	if len(objects) != 1 {
		return nil, errors.New("ambiguous query: AllFilter used to match single Object")
	}

	return objects[0], nil
}

// FindAll returns multiple Objects based on a given Filter
func (c *dynamicClient) FindAll(opts filter.ListOption) ([]runtime.Object, error) {
	return c.storage.List(storage.NewKindKey(c.gvk), opts)
}

// Delete deletes the Object from the storage
func (c *dynamicClient) Delete(uid types.UID) error {
	return c.storage.Delete(storage.NewObjectKey(storage.NewKindKey(c.gvk), runtime.NewIdentifier(string(uid))))
}

// List returns a list of all Objects available
func (c *dynamicClient) List() ([]runtime.Object, error) {
	return c.storage.List(storage.NewKindKey(c.gvk))
}
