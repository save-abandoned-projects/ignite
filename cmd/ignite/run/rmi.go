package run

import (
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/types"

	api "github.com/save-abandoned-projects/ignite/pkg/apis/ignite"
	"github.com/save-abandoned-projects/ignite/pkg/operations/lookup"
	"github.com/save-abandoned-projects/ignite/pkg/providers"
	"github.com/save-abandoned-projects/libgitops/pkg/filter"
	log "github.com/sirupsen/logrus"
)

type RmiFlags struct {
	Force bool
}

type RmiOptions struct {
	*RmiFlags
	images []*api.Image
	allVMs []*api.VM
}

func (rf *RmiFlags) NewRmiOptions(imageMatches []string) (*RmiOptions, error) {
	ro := &RmiOptions{RmiFlags: rf}
	for _, match := range imageMatches {
		images, err := providers.Client.Images().FindAll([]filter.ListOption{
			filter.NameFilter{Name: match, MatchPrefix: true},
			filter.UIDFilter{UID: types.UID(match), MatchPrefix: true}},
		)
		if err != nil {
			return nil, err
		}
		ro.images = append(ro.images, images...)
	}

	var err error
	ro.allVMs, err = getAllVMs()
	if err != nil {
		return nil, err
	}

	return ro, nil
}

func Rmi(ro *RmiOptions) error {
	for _, image := range ro.images {
		for _, vm := range ro.allVMs {
			imageUID, err := lookup.ImageUIDForVM(vm, providers.Client)
			if err != nil {
				log.Warnf("Could not lookup image UID for VM %q: %v", vm.GetUID(), err)
				continue
			}

			// Check if there's any VM using this image
			if imageUID == image.GetUID() {
				if ro.Force {
					// Force-kill and remove the VM used by this image
					if err := Rm(&RmOptions{
						&RmFlags{Force: true},
						[]*api.VM{vm},
					}); err != nil {
						return err
					}
				} else {
					return fmt.Errorf("unable to remove, image %q is in use by VM %q", image.GetUID(), vm.GetUID())
				}
			}
		}

		if err := os.RemoveAll(image.ObjectPath()); err != nil {
			return fmt.Errorf("unable to remove directory for %s %q: %v", image.GetObjectKind(), image.GetUID(), err)
		}

		fmt.Println(image.GetUID())
	}

	return nil
}
