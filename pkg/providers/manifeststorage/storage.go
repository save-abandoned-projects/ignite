package manifeststorage

import (
	api "github.com/save-abandoned-projects/ignite/pkg/apis/ignite"
	"github.com/save-abandoned-projects/ignite/pkg/apis/ignite/scheme"
	"github.com/save-abandoned-projects/ignite/pkg/constants"
	"github.com/save-abandoned-projects/ignite/pkg/providers"
	"github.com/save-abandoned-projects/libgitops/pkg/runtime"
	"github.com/save-abandoned-projects/libgitops/pkg/serializer"
	"github.com/save-abandoned-projects/libgitops/pkg/storage"
	"github.com/save-abandoned-projects/libgitops/pkg/storage/cache"
	"github.com/save-abandoned-projects/libgitops/pkg/storage/sync"
	"github.com/save-abandoned-projects/libgitops/pkg/storage/watch"
	log "github.com/sirupsen/logrus"
)

var ManifestStorage sync.SyncStorage

func SetManifestStorage() (err error) {
	log.Trace("Initializing the ManifestStorage provider...")
	ss, err := NewTwoWayManifestStorage(constants.MANIFEST_DIR, constants.DATA_DIR, scheme.Serializer)
	if err != nil {
		return err
	}

	providers.Storage = cache.NewCache(ss)

	return
}

func NewTwoWayManifestStorage(manifestDir, dataDir string, ser serializer.Serializer) (storage.Storage, error) {
	ws, err := watch.NewGenericWatchStorage(storage.NewGenericStorage(
		storage.NewGenericMappedRawStorage(manifestDir),
		scheme.Serializer,
		[]runtime.IdentifierFactory{runtime.Metav1NameIdentifier}))
	if err != nil {
		return nil, err
	}

	ss := sync.NewSyncStorage(
		storage.NewGenericStorage(
			storage.NewGenericRawStorage(dataDir, api.SchemeGroupVersion, serializer.ContentTypeYAML),
			ser,
			[]runtime.IdentifierFactory{runtime.Metav1NameIdentifier}),
		ws)

	return ss, nil
}
