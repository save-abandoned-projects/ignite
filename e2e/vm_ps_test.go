package e2e

import (
	"bytes"
	"fmt"
	"github.com/save-abandoned-projects/libgitops/pkg/serializer"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"gotest.tools/assert"

	"github.com/save-abandoned-projects/ignite/e2e/util"
	api "github.com/save-abandoned-projects/ignite/pkg/apis/ignite"
	"github.com/save-abandoned-projects/ignite/pkg/apis/ignite/scheme"
	"github.com/save-abandoned-projects/ignite/pkg/constants"
)

// TestVMpsWithOutdatedStatus checks if outdated status indicators are printed
// in ps output when the VM manifest status on disk don't match with actual
// status.
func TestVMpsWithOutdatedStatus(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	vmName := "e2e-test-ignite-ps-outdated"

	igniteCmd := util.NewCommand(t, igniteBin)

	defer igniteCmd.New().
		With("rm", "-f", vmName).
		Run()

	igniteCmd.New().
		With("run", "--name="+vmName).
		With(util.DefaultVMImage).
		With("--ssh").
		Run()

	// Filter the VM and obtain the UID.
	nameFilter := fmt.Sprintf("{{.ObjectMeta.Name}}=%s", vmName)
	psUIDCmd := igniteCmd.New().
		With("ps", "-a").
		With("-f", nameFilter).
		With("-t", "{{.ObjectMeta.UID}}")
	psUIDOut, psUIDErr := psUIDCmd.Cmd.CombinedOutput()
	assert.NilError(t, psUIDErr, fmt.Sprintf("ps: \n%q\n%s", psUIDCmd.Cmd, psUIDOut))

	uid := strings.TrimSpace(string(psUIDOut))

	// Update the VM manifest with false info.
	metadata_path := filepath.Join(constants.VM_DIR, uid, "metadata.json")
	vm := &api.VM{}
	err := scheme.Serializer.Decoder().DecodeInto(serializer.NewJSONFrameReader(serializer.FromFile(metadata_path)), vm)
	vm.Status.Running = false

	var vmBytes bytes.Buffer
	err = scheme.Serializer.Encoder().Encode(serializer.NewJSONFrameWriter(&vmBytes), vm)
	assert.NilError(t, err)
	assert.NilError(t, ioutil.WriteFile(metadata_path, vmBytes.Bytes(), 0644))

	// Revert the false data for proper cleanup.
	// NOTE: This is needed because ignite rm believes the VM manifest status
	// instead of checking for actual status itself.
	defer func() {
		vm.Status.Running = true
		var vmBytes bytes.Buffer
		err = scheme.Serializer.Encoder().Encode(serializer.NewJSONFrameWriter(&vmBytes), vm)
		assert.NilError(t, err)
		assert.NilError(t, ioutil.WriteFile(metadata_path, vmBytes.Bytes(), 0644))
	}()

	// Run ps -a and look for the outdated status info.
	psCmd := igniteCmd.New().
		With("ps", "-a")

	psOut, psErr := psCmd.Cmd.CombinedOutput()
	assert.NilError(t, psErr, fmt.Sprintf("ps: \n%q\n%s", psCmd.Cmd, psOut))

	psOutString := string(psOut)
	// Check for the outdated status and the note about it.
	assert.Check(t, strings.Contains(psOutString, "*Up"))
	assert.Check(t, strings.Contains(psOutString, "The symbol *"))
}
