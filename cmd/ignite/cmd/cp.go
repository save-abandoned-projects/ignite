package cmd

import (
	"io"

	"github.com/save-abandoned-projects/ignite/cmd/ignite/cmd/vmcmd"
	"github.com/spf13/cobra"
)

// NewCmdCP is an alias for vmcmd.NewCmdCP
func NewCmdCP(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdCP(out)
}
