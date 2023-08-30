package cmd

import (
	"io"

	"github.com/save-abandoned-projects/ignite/cmd/ignite/cmd/vmcmd"
	"github.com/spf13/cobra"
)

// NewCmdSSH is an alias for vmcmd.NewCmdSSH
func NewCmdSSH(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdSSH(out)
}
