package storage

import (
	api "github.com/save-abandoned-projects/ignite/pkg/apis/ignite"
	"github.com/save-abandoned-projects/ignite/pkg/apis/ignite/scheme"
	"github.com/save-abandoned-projects/ignite/pkg/constants"
	"github.com/save-abandoned-projects/ignite/pkg/providers"
	"github.com/save-abandoned-projects/libgitops/pkg/runtime"
	"github.com/save-abandoned-projects/libgitops/pkg/serializer"
	"github.com/save-abandoned-projects/libgitops/pkg/storage"
	"github.com/save-abandoned-projects/libgitops/pkg/storage/cache"
	log "github.com/sirupsen/logrus"
)

func SetGenericStorage() error {
	log.Trace("Initializing the GenericStorage provider...")
	providers.Storage = cache.NewCache(
		storage.NewGenericStorage(
			storage.NewGenericRawStorage(constants.DATA_DIR, api.SchemeGroupVersion, serializer.ContentTypeYAML),
			scheme.Serializer,
			[]runtime.IdentifierFactory{runtime.Metav1NameIdentifier, runtime.ObjectUIDIdentifier}))
	return nil
}
