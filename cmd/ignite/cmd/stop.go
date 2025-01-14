package cmd

import (
	"io"

	"github.com/save-abandoned-projects/ignite/cmd/ignite/cmd/vmcmd"
	"github.com/spf13/cobra"
)

// NewCmdStop is an alias for vmcmd.NewCmdStop
func NewCmdStop(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdStop(out)
}
