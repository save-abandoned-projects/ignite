package cmd

import (
	"io"

	"github.com/save-abandoned-projects/ignite/cmd/ignite/cmd/vmcmd"
	"github.com/spf13/cobra"
)

// NewCmdLogs is an alias for vmcmd.NewCmdLogs
func NewCmdLogs(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdLogs(out)
}
