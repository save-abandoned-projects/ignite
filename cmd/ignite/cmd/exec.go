package cmd

import (
	"io"

	"github.com/save-abandoned-projects/ignite/cmd/ignite/cmd/vmcmd"
	"github.com/spf13/cobra"
)

// NewCmdExec is an alias for vmcmd.NewCmdExec
func NewCmdExec(out io.Writer, err io.Writer, in io.Reader) *cobra.Command {
	return vmcmd.NewCmdExec(out, err, in)
}
