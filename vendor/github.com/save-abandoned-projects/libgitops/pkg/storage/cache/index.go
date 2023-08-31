package cache

import (
	"github.com/save-abandoned-projects/libgitops/pkg/filter"
	"github.com/save-abandoned-projects/libgitops/pkg/runtime"
	"github.com/save-abandoned-projects/libgitops/pkg/storage"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

type index struct {
	storage storage.Storage
	objects map[schema.GroupVersionKind]map[types.UID]*cacheObject
}

func newIndex(storage storage.Storage) *index {
	return &index{
		storage: storage,
		objects: make(map[schema.GroupVersionKind]map[types.UID]*cacheObject),
	}
}

func (i *index) loadByID(key storage.ObjectKey) (runtime.Object, error) {
	if uids, ok := i.objects[key.GetGVK()]; ok {
		if obj, ok := uids[types.UID(key.GetIdentifier())]; ok {
			log.Tracef("index: cache hit for %s with UID %q", key.GetKind(), key.GetIdentifier())
			return obj.loadFull()
		}
	}

	log.Tracef("index: cache miss for %s with UID %q", key.GetKind(), key.GetIdentifier())
	return nil, nil
}

func (i *index) loadAll() ([]runtime.Object, error) {
	var size uint64

	for gvk := range i.objects {
		size += i.count(gvk)
	}

	all := make([]runtime.Object, 0, size)

	for gvk := range i.objects {
		if objects, err := i.list(storage.NewKindKey(gvk)); err == nil {
			all = append(all, objects...)
		} else {
			return nil, err
		}
	}

	return all, nil
}

func store(i *index, obj runtime.Object, apiType bool) error {
	// If store is called for an invalid Object lacking an UID,
	// panic and print the stack trace. This should never happen.
	if obj.GetUID() == "" {
		panic("Attempt to cache invalid Object: missing UID")
	}

	co, err := newCacheObject(i.storage, obj, apiType)
	if err != nil {
		return err
	}

	gvk := co.object.GetObjectKind().GroupVersionKind()

	if _, ok := i.objects[gvk]; !ok {
		i.objects[gvk] = make(map[types.UID]*cacheObject)
	}

	log.Tracef("index: storing %s object with UID %q, meta: %t", gvk.Kind, obj.GetName(), apiType)
	i.objects[gvk][co.object.GetUID()] = co

	return nil
}

func (i *index) store(obj runtime.Object) error {
	return store(i, obj, false)
}

func (i *index) storeAll(objs []runtime.Object) (err error) {
	for _, obj := range objs {
		if err = i.store(obj); err != nil {
			break
		}
	}

	return
}

func (i *index) storeMeta(obj runtime.PartialObject) error {
	return store(i, obj, true)
}

func (i *index) storeAllMeta(objs []runtime.PartialObject) (err error) {
	for _, obj := range objs {
		if uids, ok := i.objects[obj.GetObjectKind().GroupVersionKind()]; ok {
			if _, ok := uids[obj.GetUID()]; ok {
				continue
			}
		}

		if err = i.storeMeta(obj); err != nil {
			break
		}
	}

	return
}

func (i *index) delete(key storage.ObjectKey) {
	if uids, ok := i.objects[key.GetGVK()]; ok {
		delete(uids, types.UID(key.GetIdentifier()))
	}
}

func (i *index) count(gvk schema.GroupVersionKind) (count uint64) {
	count = uint64(len(i.objects[gvk]))
	log.Tracef("index: counted %d %s object(s)", count, gvk.Kind)
	return
}

func list(i *index, gvk schema.GroupVersionKind) ([]runtime.Object, error) {
	uids := i.objects[gvk]
	list := make([]runtime.Object, 0, len(uids))

	log.Tracef("index: listing %s objects", gvk)
	for _, obj := range uids {
		if result, err := obj.loadFull(); err != nil {
			return nil, err
		} else {
			list = append(list, result)
		}
	}

	return list, nil
}

func listMeta(i *index, gvk schema.GroupVersionKind) ([]runtime.PartialObject, error) {
	uids := i.objects[gvk]
	list := make([]runtime.PartialObject, 0, len(uids))

	log.Tracef("index: listing %s objects meta", gvk)
	for _, obj := range uids {
		if result, err := obj.loadAPI(); err != nil {
			return nil, err
		} else {
			list = append(list, result.(runtime.PartialObject))
		}
	}

	return list, nil
}

func (i *index) list(kind storage.KindKey, opts ...filter.ListOption) ([]runtime.Object, error) {
	return list(i, kind.GetGVK())
}

func (i *index) listMeta(kind storage.KindKey) ([]runtime.PartialObject, error) {
	return listMeta(i, kind.GetGVK())
}
