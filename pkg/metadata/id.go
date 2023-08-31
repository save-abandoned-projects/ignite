package metadata

import (
	"fmt"

	api "github.com/save-abandoned-projects/ignite/pkg/apis/ignite"

	"github.com/save-abandoned-projects/ignite/pkg/logs"
	"github.com/save-abandoned-projects/ignite/pkg/providers"
	log "github.com/sirupsen/logrus"
)

// TODO: Get rid of this
var success = make(map[Metadata]bool)

// silent specifies if the ID should be printed, when chaining commands
// silence all but the last command to print the ID only once
func Cleanup(md Metadata, silent bool) error {
	// If success has not been confirmed, remove the generated directory
	kind := md.GetObjectKind().GroupVersionKind().Kind
	if !success[md] {
		if !logs.Quiet {
			log.Infof("Removed %s with name %q and ID %q", kind, md.GetName(), md.GetUID())
		} else if !silent {
			fmt.Println(md.GetUID())
		}

		return providers.Client.Dynamic(api.Kind(kind)).Delete(md.GetUID())
	}

	if !logs.Quiet {
		log.Infof("Created %s with ID %q and name %q", kind, md.GetUID(), md.GetName())
	} else if !silent {
		fmt.Println(md.GetUID())
	}

	return nil
}

// Should be returned as the last command when creating objects
func Success(md Metadata) error {
	success[md] = true
	return nil
}
