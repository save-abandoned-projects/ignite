package cache

import (
	"github.com/save-abandoned-projects/libgitops/pkg/runtime"
	"github.com/save-abandoned-projects/libgitops/pkg/storage"
	log "github.com/sirupsen/logrus"
)

type cacheObject struct {
	storage  storage.Storage
	object   runtime.Object
	checksum string
	apiType  bool
}

func newCacheObject(s storage.Storage, object runtime.Object, apiType bool) (c *cacheObject, err error) {
	c = &cacheObject{
		storage: s,
		object:  object,
		apiType: apiType,
	}
	if c.checksum, err = s.Checksum(c.getObjectKey()); err != nil {
		c = nil
	}

	return
}

// loadFull returns the full Object, loading it only if it hasn't been cached before or the checksum has changed
func (c *cacheObject) loadFull() (runtime.Object, error) {
	var checksum string
	reload := c.apiType

	if !reload {
		if chk, err := c.storage.Checksum(c.getObjectKey()); err != nil {
			return nil, err
		} else if chk != c.checksum {
			log.Tracef("cacheObject: %q invalidated, checksum mismatch: %q -> %q", c.object.GetName(), c.checksum, chk)
			checksum = chk
			reload = true
		} else {
			log.Tracef("cacheObject: %q checksum: %q", c.object.GetName(), c.checksum)
		}
	}

	if reload {
		log.Tracef("cacheObject: full load triggered for %q", c.object.GetName())
		obj, err := c.storage.Get(c.getObjectKey())
		if err != nil {
			return nil, err
		}

		// Only apply the change after a successful Get
		c.object = obj
		c.apiType = false

		if len(checksum) > 0 {
			c.checksum = checksum
		}
	}

	return c.object, nil
}

// loadAPI returns the APIType of the Object, loading it only if the checksum has changed
func (c *cacheObject) loadAPI() (runtime.Object, error) {
	if chk, err := c.storage.Checksum(c.getObjectKey()); err != nil {
		return nil, err
	} else if chk != c.checksum {
		log.Tracef("cacheObject: %q invalidated, checksum mismatch: %q -> %q", c.object.GetName(), c.checksum, chk)
		log.Tracef("cacheObject: API load triggered for %q", c.object.GetName())
		obj, err := c.storage.GetMeta(c.getObjectKey())
		if err != nil {
			return nil, err
		}

		// Only apply the change after a successful GetMeta
		c.object = obj
		c.checksum = chk
		c.apiType = true
	} else {
		log.Tracef("cacheObject: %q checksum: %q", c.object.GetName(), c.checksum)
	}

	if c.apiType {
		return c.object, nil
	}

	return runtime.PartialObjectFrom(c.object)
}

func (c *cacheObject) getObjectKey() storage.ObjectKey {
	return storage.NewObjectKey(storage.NewKindKey(c.object.GetObjectKind().GroupVersionKind()), runtime.NewIdentifier(string(c.object.GetUID())))
}
