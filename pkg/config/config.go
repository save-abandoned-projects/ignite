package config

import (
	"fmt"
	"os"

	api "github.com/save-abandoned-projects/ignite/pkg/apis/ignite"
	"github.com/save-abandoned-projects/ignite/pkg/apis/ignite/scheme"
	"github.com/save-abandoned-projects/ignite/pkg/constants"
	"github.com/save-abandoned-projects/ignite/pkg/network"
	"github.com/save-abandoned-projects/ignite/pkg/providers"
	"github.com/save-abandoned-projects/ignite/pkg/providers/ignite"
	"github.com/save-abandoned-projects/ignite/pkg/runtime"
	"github.com/save-abandoned-projects/libgitops/pkg/serializer"
	log "github.com/sirupsen/logrus"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
)

// ApplyConfiguration merges the given configurations with the default ignite
// configurations.
func ApplyConfiguration(configPath string) error {
	var configFilePath string

	if configPath != "" {
		configFilePath = configPath
	} else {
		// Check the default config location.
		if _, err := os.Stat(constants.IGNITE_CONFIG_FILE); !os.IsNotExist(err) {
			log.Debugf("Found default ignite configuration file %s", constants.IGNITE_CONFIG_FILE)
			configFilePath = constants.IGNITE_CONFIG_FILE
		}
	}

	if configFilePath != "" {
		log.Debugf("Using ignite configuration file %s", configFilePath)
		var err error
		providers.ComponentConfig, err = getConfigFromFile(configFilePath)
		if err != nil {
			return err
		}

		// Set providers runtime and network plugin if found in config
		// and not set explicitly via flags.
		if providers.ComponentConfig.Spec.Runtime != "" && providers.RuntimeName == "" {
			providers.RuntimeName = providers.ComponentConfig.Spec.Runtime
		}
		if providers.ComponentConfig.Spec.NetworkPlugin != "" && providers.NetworkPluginName == "" {
			providers.NetworkPluginName = providers.ComponentConfig.Spec.NetworkPlugin
		}
		if providers.ComponentConfig.Spec.IDPrefix != "" && providers.IDPrefix == "" {
			providers.IDPrefix = providers.ComponentConfig.Spec.IDPrefix
		}
	} else {
		log.Debugln("Using ignite default configurations")
	}

	// Set the default runtime and network-plugin if it's not set by
	// now.
	if providers.RuntimeName == "" {
		providers.RuntimeName = runtime.RuntimeContainerd
	}
	if providers.NetworkPluginName == "" {
		providers.NetworkPluginName = network.PluginCNI
	}
	if providers.IDPrefix == "" {
		providers.IDPrefix = constants.IGNITE_PREFIX
	}

	return nil
}

// getConfigFromFile reads a config file and returns ignite configuration.
func getConfigFromFile(configPath string) (*api.Configuration, error) {
	gvs := scheme.Scheme.PrioritizedVersionsForGroup(api.GroupName)
	obj, err := scheme.Scheme.New(gvs[0].WithKind(api.KindConfiguration.Title()))
	if err != err {
		return nil, err
	}

	// TODO: Fix libgitops DecodeFileInto to not allow empty files.
	if err := scheme.Serializer.Decoder().DecodeInto(serializer.NewYAMLFrameReader(serializer.FromFile(configPath)), obj); err != nil {
		return nil, err
	}

	scheme.Scheme.Default(obj)
	// Ensure the read configuration is valid. If a file contains Kind and
	// APIVersion, it's a valid config file. Empty file is invalid.
	// NOTE: This is a workaround for libgitops allowing to decode of empty file.
	if obj.GetObjectKind().GroupVersionKind().Kind == "" || obj.GetObjectKind().GroupVersionKind().Version == "" {
		return nil, fmt.Errorf("invalid config file, Kind and APIVersion must be set")
	}

	if internalObj, err := scheme.Scheme.ConvertToVersion(obj, k8sruntime.InternalGroupVersioner); err != nil {
		return nil, err
	} else {
		return internalObj.(*api.Configuration), nil
	}
}

// SetAndPopulateProviders sets and populates the providers.
func SetAndPopulateProviders(runtimeName runtime.Name, networkPlugin network.PluginName) error {
	providers.RuntimeName = runtimeName
	providers.NetworkPluginName = networkPlugin
	return providers.Populate(ignite.Providers)
}
