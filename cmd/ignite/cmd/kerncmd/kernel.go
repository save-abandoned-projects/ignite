package kerncmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/save-abandoned-projects/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/save-abandoned-projects/ignite/cmd/ignite/run"
	"github.com/spf13/cobra"
)

// NewCmdKernel handles kernel-related functionality via its subcommands
// This command by itself lists available kernels
func NewCmdKernel(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kernel",
		Short: "Manage VM kernels",
		Long: dedent.Dedent(`
			Groups together functionality for managing VM kernels.
			Calling this command alone lists all available kernels.
		`),
		Aliases: []string{"kernels"},
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(func() error {
				ko, err := run.NewKernelsOptions()
				if err != nil {
					return err
				}

				return run.Kernels(ko)
			}())
		},
	}

	cmd.AddCommand(NewCmdImport(out))
	cmd.AddCommand(NewCmdLs(out))
	cmd.AddCommand(NewCmdRm(out))
	return cmd
}
