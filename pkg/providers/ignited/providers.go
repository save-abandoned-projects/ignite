package ignited

import (
	"github.com/save-abandoned-projects/ignite/pkg/providers"
	clientprovider "github.com/save-abandoned-projects/ignite/pkg/providers/client"
	manifeststorageprovider "github.com/save-abandoned-projects/ignite/pkg/providers/manifeststorage"
	"github.com/save-abandoned-projects/ignite/pkg/providers/network"
	"github.com/save-abandoned-projects/ignite/pkg/providers/runtime"
)

// Preload providers need to be loaded before flag parsing has finished
var Preload = []providers.ProviderInitFunc{
	manifeststorageprovider.SetManifestStorage, // Use the ManifestStorage implementation, backed by a cache
	clientprovider.SetClient,                   // Set the globally available client
}

// NOTE: Provider initialization is order-dependent!
// E.g. the network plugin depends on the runtime.
var Providers = []providers.ProviderInitFunc{
	runtime.SetRuntime,       // Set the selected runtime
	network.SetNetworkPlugin, // Set the selected network plugin
}
