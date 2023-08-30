package metadata

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	api "github.com/save-abandoned-projects/ignite/pkg/apis/ignite"
	"github.com/save-abandoned-projects/ignite/pkg/client"
	"github.com/save-abandoned-projects/ignite/pkg/constants"
	"github.com/save-abandoned-projects/ignite/pkg/providers"
	"github.com/save-abandoned-projects/ignite/pkg/util"
	"github.com/save-abandoned-projects/libgitops/pkg/filter"
	"github.com/save-abandoned-projects/libgitops/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

var (
	nameRegex = regexp.MustCompile(`^[a-z-_0-9.:/@]*$`)
	uidRegex  = regexp.MustCompile(`^[a-z0-9]{16}$`)
)

// ErrNilObject is returned when the runtime object being initialized is nil.
var ErrNilObject = errors.New("object cannot be nil when initializing runtime data")

type Metadata interface {
	runtime.Object
}

// SetNameAndUID sets the name and UID for an object if that isn't set automatically
func SetNameAndUID(obj runtime.Object, c *client.Client) error {
	if obj == nil {
		return ErrNilObject
	}

	if c == nil {
		c = providers.Client
	}

	// Generate or validate the given UID, if any
	if err := processUID(obj, c); err != nil {
		return err
	}

	// Generate or validate the given name, if any
	return processName(obj, c)
}

// SetLabels metadata labels for a given object.
func SetLabels(obj runtime.Object, labels []string) error {
	if obj == nil {
		return ErrNilObject
	}

	for _, label := range labels {
		kv := strings.SplitN(label, "=", 2)
		// Check the length of key/val. Must be exactly 2.
		if len(kv) != 2 {
			return fmt.Errorf("invalid label value %q: supported syntax: --label <key>=<value>", label)
		}
		// Label name must not be empty.
		if kv[0] == "" {
			return fmt.Errorf("invalid label %q, name empty", label)
		}
		obj.SetLabels(map[string]string{kv[0]: kv[1]})
	}
	return nil
}

// processUID a new 8-byte ID and handles directory creation/deletion
func processUID(obj runtime.Object, c *client.Client) error {
	uid := string(obj.GetUID())

	kind := api.Kind(obj.GetObjectKind().GroupVersionKind().Kind)
	// Validate the given UID if set
	if len(uid) > 0 {
		// Verify that if specified
		if !uidRegex.MatchString(uid) {
			return fmt.Errorf("invalid UID %q: does not match required format %s", uid, uidRegex.String())
		}

		// Make sure there isn't any duplicate names
		if err := verifyUIDOrName(c, uid, kind); err != nil {
			return err
		}
	} else {
		// No UID set, generate one
		var err error
		for {
			if uid, err = util.NewUID(); err != nil {
				return fmt.Errorf("failed to generate ID: %v", err)
			}

			// If the generated UID is unique break the generator loop
			if err := verifyUIDOrName(c, uid, kind); err == nil {
				// Set the generated UID to the object
				obj.SetUID(types.UID(uid))
				break
			}
		}
	}

	// Create the directory for the specified UID
	// TODO: Move this kind of functionality into pkg/storage
	dir := path.Join(constants.DATA_DIR, string(bytes.ToLower([]byte(kind))), uid)
	if err := os.MkdirAll(dir, constants.DATA_DIR_PERM); err != nil {
		return fmt.Errorf("failed to create directory for ID %q: %v", uid, err)
	}

	return nil
}

func processName(obj runtime.Object, c *client.Client) error {
	name := obj.GetName()
	kind := api.Kind(obj.GetObjectKind().GroupVersionKind().Kind)

	// Enforce a latest tag for images and kernels. Also,
	// images and kernels must have their name set at this stage
	if kind == api.KindImage || kind == api.KindKernel {
		if len(name) == 0 {
			// this should not happen, programmer error
			return fmt.Errorf("%s name must not be unset", kind)
		}
	} else if len(name) == 0 { // If some other kind's name is empty, set a random name
		name = util.RandomName()
	}

	// Validate the name with the regexp
	if !nameRegex.MatchString(name) {
		return fmt.Errorf("invalid name %q: does not match required format %s", name, nameRegex.String())
	}

	// Make sure there isn't any duplicate names
	if err := verifyUIDOrName(c, name, kind); err != nil {
		return err
	}

	// write the desired name to the object
	obj.SetName(name)
	return nil
}

func verifyUIDOrName(c *client.Client, match string, kind api.Kind) error {
	_, err := c.Dynamic(kind).Find(filter.NameFilter{Name: match})
	if err != nil {
		return err
	}
	return nil
}
