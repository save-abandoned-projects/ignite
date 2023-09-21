package cmd

import (
	"io"

	"github.com/save-abandoned-projects/ignite/cmd/ignite/cmd/vmcmd"
	"github.com/spf13/cobra"
)

// NewCmdCreate is an alias for vmcmd.NewCmdCreate
func NewCmdCreate(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdCreate(out)
}
