// shellgym - a background daemon that turns a Linux box into an
// interactive command-line trainer (Shell Gym).
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/iximiuz/labs-content/tools/shellgym/internal/checkclient"
)

var version = "dev"

func main() {
	root := &cobra.Command{
		Use:           "shellgym",
		Short:         "Shell Gym - an interactive Linux command-line trainer",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(
		newServeCmd(),
		newValidateCmd(),
		newSolveCmd(),
		newCheckCmd(),
		newSkillsCmd(),
		&cobra.Command{
			Use:   "version",
			Short: "Print the version",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println(version)
			},
		},
	)
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "shellgym:", err)
		os.Exit(1)
	}
}

// check is the hidden entry point behind the PATH shims that task scripts
// call (wait_cwd, wait_exec, hint_exit, ...).
func newCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "check <name> [args...]",
		Short:              "Run a built-in check (internal, used via PATH shims)",
		Hidden:             true,
		Args:               cobra.MinimumNArgs(1),
		DisableFlagParsing: true, // checks parse their own --timeout/--now
		RunE: func(cmd *cobra.Command, args []string) error {
			os.Exit(checkclient.Main(args[0], args[1:]))
			return nil
		},
	}
	return cmd
}
