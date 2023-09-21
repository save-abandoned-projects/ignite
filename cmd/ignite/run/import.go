package run

import (
	"github.com/save-abandoned-projects/ignite/cmd/ignite/cmd/cmdutil"
	api "github.com/save-abandoned-projects/ignite/pkg/apis/ignite"
	meta "github.com/save-abandoned-projects/ignite/pkg/apis/meta/v1alpha1"
	"github.com/save-abandoned-projects/ignite/pkg/config"
	"github.com/save-abandoned-projects/ignite/pkg/metadata"
	"github.com/save-abandoned-projects/ignite/pkg/operations"
	"github.com/save-abandoned-projects/ignite/pkg/providers"
	"github.com/save-abandoned-projects/ignite/pkg/util"
)

func ImportImage(source string) (image *api.Image, err error) {
	// Populate the runtime provider.
	if err := config.SetAndPopulateProviders(providers.RuntimeName, providers.NetworkPluginName); err != nil {
		return nil, err
	}

	cmdutil.ResolveRegistryConfigDir()

	ociRef, err := meta.NewOCIImageRef(source)
	if err != nil {
		return
	}

	image, err = operations.FindOrImportImage(providers.Client, ociRef)
	if err != nil {
		return
	}
	defer util.DeferErr(&err, func() error { return metadata.Cleanup(image, false) })

	err = metadata.Success(image)

	return
}

func ImportKernel(source string) (kernel *api.Kernel, err error) {
	// Populate the runtime provider.
	if err := config.SetAndPopulateProviders(providers.RuntimeName, providers.NetworkPluginName); err != nil {
		return nil, err
	}

	cmdutil.ResolveRegistryConfigDir()

	ociRef, err := meta.NewOCIImageRef(source)
	if err != nil {
		return
	}

	kernel, err = operations.FindOrImportKernel(providers.Client, ociRef)
	if err != nil {
		return
	}
	defer util.DeferErr(&err, func() error { return metadata.Cleanup(kernel, false) })

	err = metadata.Success(kernel)

	return
}
