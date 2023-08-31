package cache

import (
	"github.com/save-abandoned-projects/libgitops/pkg/filter"
	"github.com/save-abandoned-projects/libgitops/pkg/runtime"
	"github.com/save-abandoned-projects/libgitops/pkg/serializer"
	"github.com/save-abandoned-projects/libgitops/pkg/storage"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Cache is an intermediate caching layer, which conforms to Storage
// Typically you back the cache with an actual storage
type Cache interface {
	storage.Storage
	// Flush is used to write the state of the entire cache to storage
	// Warning: this is a very expensive operation
	Flush() error
}

type cache struct {
	// storage is the backing Storage for the cache
	// used to look up non-cached Objects
	storage storage.Storage

	// index caches the Objects by GroupVersionKind and UID
	// This guarantees uniqueness when looking up a specific Object
	index *index
}

func (c *cache) Get(key storage.ObjectKey) (obj runtime.Object, err error) {
	log.Tracef("cache: Get %s with UID %q", key.GetKind(), key.GetIdentifier())

	// If the requested Object resides in the cache, return it
	if obj, err = c.index.loadByID(key); err != nil || obj != nil {
		return
	}

	// Request the Object from the storage
	obj, err = c.storage.Get(key)
	// If no errors occurred, cache it
	if err == nil {
		err = c.index.store(obj)
	}

	return
}

func (c *cache) GetMeta(key storage.ObjectKey) (obj runtime.PartialObject, err error) {
	log.Tracef("cache: GetMeta %s with UID %q", key.GetKind(), key.GetIdentifier())

	obj, err = c.storage.GetMeta(key)

	// If no errors occurred while loading, store the Object in the cache
	if err == nil {
		err = c.index.storeMeta(obj)
	}

	return
}

func (c *cache) List(kind storage.KindKey, opts ...filter.ListOption) ([]runtime.Object, error) {
	return c.list(kind, c.storage.List, c.index.list, c.index.storeAll)
}

func (c *cache) ListMeta(kind storage.KindKey) ([]runtime.PartialObject, error) {
	return c.listMeta(kind, c.storage.ListMeta, c.index.listMeta, c.index.storeAllMeta)
}

func (s *cache) Find(kind storage.KindKey, opts ...filter.ListOption) (runtime.Object, error) {
	//TODO implement me
	panic("implement me")
}

func (c *cache) Checksum(key storage.ObjectKey) (string, error) {
	// The cache is transparent about the checksums
	return c.storage.Checksum(key)
}

func (c *cache) Count(kind storage.KindKey) (uint64, error) {
	// The cache is transparent about how many items it has cached
	return c.storage.Count(kind)
}

func (c *cache) ObjectKeyFor(obj runtime.Object) (storage.ObjectKey, error) {
	return storage.NewObjectKey(storage.NewKindKey(obj.GetObjectKind().GroupVersionKind()), runtime.NewIdentifier(string(obj.GetUID()))), nil
}

func (c *cache) Create(obj runtime.Object) error {
	c.index.store(obj)
	return c.storage.Create(obj)
}

func (c *cache) Update(obj runtime.Object) error {
	c.index.store(obj)
	return c.storage.Update(obj)
}

func (c *cache) Patch(key storage.ObjectKey, patch []byte) error {
	// TODO: For now patches are always flushed, the cache will load the updated Object on-demand on access
	return c.storage.Patch(storage.NewObjectKey(storage.NewKindKey(key.GetGVK()), runtime.NewIdentifier(key.GetIdentifier())), patch)
}

func (c *cache) Delete(key storage.ObjectKey) error {
	log.Tracef("cache: Delete %s with UID %q", key.GetKind(), key.GetIdentifier())

	// Delete the given Object from the cache and storage
	c.index.delete(key)
	return c.storage.Delete(key)
}

var _ Cache = &cache{}

func NewCache(backingStorage storage.Storage) Cache {
	c := &cache{
		storage: backingStorage,
		index:   newIndex(backingStorage),
	}

	return c
}

func (s *cache) Serializer() serializer.Serializer {
	return s.storage.Serializer()
}

//func (c *cache) New(gvk schema.GroupVersionKind) (runtime.Object, error) {
//	// Request the storage to create the Object. The
//	// newly generated Object has not got an UID which
//	// is required for indexing, so just return it
//	// without storing it into the cache
//	return c.storage.New(gvk)
//}

func (c *cache) Set(gvk schema.GroupVersionKind, obj runtime.Object) error {
	log.Tracef("cache: Set %s with UID %q", gvk.Kind, obj.GetUID())

	// Store the changed Object in the cache
	if err := c.index.store(obj); err != nil {
		return err
	}

	// TODO: For now the cache always flushes, we might add automatic flushing later
	return c.storage.Update(obj)
}

type listFunc func(kind storage.KindKey, opts ...filter.ListOption) ([]runtime.Object, error)
type cacheStoreFunc func([]runtime.Object) error
type listMetaFunc func(kind storage.KindKey) ([]runtime.PartialObject, error)
type cacheMetaStoreFunc func([]runtime.PartialObject) error

// list is a common handler for List and ListMeta
func (c *cache) list(kind storage.KindKey, slf, clf listFunc, csf cacheStoreFunc) (objs []runtime.Object, err error) {
	var storageCount uint64
	if storageCount, err = c.storage.Count(storage.NewObjectKey(kind, runtime.NewIdentifier(string(rune(0))))); err != nil {
		return
	}

	if c.index.count(kind.GetGVK()) != storageCount {
		log.Tracef("cache: miss when listing: %s", kind.GetGVK())
		// If the cache doesn't track all the Objects, request them from the storage
		if objs, err = slf(storage.NewObjectKey(kind, runtime.NewIdentifier(string(rune(0))))); err != nil {
			// If no errors occurred, store the Objects in the cache
			err = csf(objs)
		}
	} else {
		log.Tracef("cache: hit when listing: %s", kind.GetGVK())
		// If the cache tracks everything, return the cache's contents
		objs, err = clf(storage.NewObjectKey(kind, runtime.NewIdentifier(string(rune(0)))))
	}

	return
}

// list is a common handler for List and ListMeta
func (c *cache) listMeta(kind storage.KindKey, slf, clf listMetaFunc, cmsf cacheMetaStoreFunc) (objs []runtime.PartialObject, err error) {
	var storageCount uint64
	if storageCount, err = c.storage.Count(storage.NewObjectKey(kind, runtime.NewIdentifier(string(rune(0))))); err != nil {
		return
	}

	if c.index.count(kind.GetGVK()) != storageCount {
		log.Tracef("cache: miss when listing: %s", kind.GetGVK())
		// If the cache doesn't track all the Objects, request them from the storage
		if objs, err = slf(storage.NewKindKey(kind.GetGVK())); err != nil {
			// If no errors occurred, store the Objects in the cache
			err = cmsf(objs)
		}
	} else {
		log.Tracef("cache: hit when listing: %s", kind.GetGVK())
		// If the cache tracks everything, return the cache's contents
		objs, err = clf(kind)
	}

	return
}

func (c *cache) RawStorage() storage.RawStorage {
	return c.storage.RawStorage()
}

func (c *cache) Close() error {
	return c.storage.Close()
}

func (c *cache) Flush() error {
	// Load the entire cache
	allObjects, err := c.index.loadAll()
	if err != nil {
		return err
	}

	for _, obj := range allObjects {
		// Request the storage to save each Object
		if err := c.storage.Update(obj); err != nil {
			return err
		}
	}

	return nil
}
