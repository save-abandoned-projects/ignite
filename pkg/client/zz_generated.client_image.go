
/*
	Note: This file is autogenerated! Do not edit it manually!
	Edit client_image_template.go instead, and run
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

// ImageClient is an interface for accessing Image-specific API objects
type ImageClient interface {
	// New returns a new Image
	New() *api.Image
	// Get returns the Image matching given UID from the storage
	Get(types.UID) (*api.Image, error)
	// Set saves the given Image into persistent storage
	Set(*api.Image) error
	// Patch performs a strategic merge patch on the object with
	// the given UID, using the byte-encoded patch given
	Patch(types.UID, []byte) error
	// Find returns the Image matching the given filter, filters can
	// match e.g. the Object's Name, UID or a specific property
	Find(opt filter.ObjectFilter) (*api.Image, error)
	// FindAll returns multiple Images matching the given filter, filters can
	// match e.g. the Object's Name, UID or a specific property
	FindAll(opts []filter.ListOption) ([]*api.Image, error)
	// Delete deletes the Image with the given UID from the storage
	Delete(uid types.UID) error
	// List returns a list of all Images available
	List() ([]*api.Image, error)
}

// Images returns the ImageClient for the IgniteInternalClient instance
func (c *IgniteInternalClient) Images() ImageClient {
	if c.imageClient == nil {
		c.imageClient = newImageClient(c.storage, c.gv)
	}

	return c.imageClient
}

// imageClient is a struct implementing the ImageClient interface
// It uses a shared storage instance passed from the Client together with its own Filterer
type imageClient struct {
	storage storage.Storage
	gvk     schema.GroupVersionKind
}

// newImageClient builds the imageClient struct using the storage implementation and a new Filterer
func newImageClient(s storage.Storage, gv schema.GroupVersion) ImageClient {
	return &imageClient{
		storage: s,
		gvk:     gv.WithKind(api.KindImage.Title()),
	}
}

// New returns a new Object of its kind
func (c *imageClient) New() *api.Image {
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
	return obj.(*api.Image)
}

// Find returns a single Image based on the given Filter
func (c *imageClient) Find(opt filter.ObjectFilter) (*api.Image, error) {
	log.Tracef("Client.Find; GVK: %v", c.gvk)

	var opts []filter.ListOption
	switch o := opt.(type) {
	case filter.NameFilter, filter.UIDFilter, filter.GvkFilter:
		opts = append(opts, o.(filter.ListOption))
	default:
		return nil, errors.New("bad filter")
	}
	objects, err := c.FindAll(opts)
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

// FindAll returns multiple Images based on the given Filter
func (c *imageClient) FindAll(opts []filter.ListOption) ([]*api.Image, error) {
	log.Tracef("Client.FindAll; GVK: %v", c.gvk)
	matches, err := c.storage.List(storage.NewKindKey(c.gvk), opts...)
	if err != nil {
		return nil, err
	}

	results := make([]*api.Image, 0, len(matches))
	for _, item := range matches {
		results = append(results, item.(*api.Image))
	}

	return results, nil
}

// Get returns the Image matching given UID from the storage
func (c *imageClient) Get(uid types.UID) (*api.Image, error) {
	log.Tracef("Client.Get; UID: %q, GVK: %v", uid, c.gvk)
	object, err := c.storage.Get(storage.NewObjectKey(storage.NewKindKey(c.gvk), runtime.NewIdentifier(string(uid))))
	if err != nil {
		return nil, err
	}

	return object.(*api.Image), nil
}

// Set saves the given Image into the persistent storage
func (c *imageClient) Set(image *api.Image) error {
	log.Tracef("Client.Set; UID: %q, GVK: %v", image.GetUID(), c.gvk)

	return c.storage.Create(image)
}

// Patch performs a strategic merge patch on the object with
// the given UID, using the byte-encoded patch given
func (c *imageClient) Patch(uid types.UID, patch []byte) error {
	return c.storage.Patch(storage.NewObjectKey(storage.NewKindKey(c.gvk), runtime.NewIdentifier(string(uid))), patch)
}

// Delete deletes the Image from the storage
func (c *imageClient) Delete(uid types.UID) error {
	log.Tracef("Client.Delete; UID: %q, GVK: %v", uid, c.gvk)
	return c.storage.Delete(storage.NewObjectKey(storage.NewKindKey(c.gvk), runtime.NewIdentifier(string(uid))))
}

// List returns a list of all Images available
func (c *imageClient) List() ([]*api.Image, error) {
	log.Tracef("Client.List; GVK: %v", c.gvk)
	list, err := c.storage.List(storage.NewKindKey(c.gvk))
	if err != nil {
		return nil, err
	}

	results := make([]*api.Image, 0, len(list))
	for _, item := range list {
		results = append(results, item.(*api.Image))
	}

	return results, nil
}
