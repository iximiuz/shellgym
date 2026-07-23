package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/iximiuz/labs-content/tools/shellgym/skills"
)

// newSkillsCmd dumps embedded authoring skills (markdown, ready to be
// dropped into .claude/skills/<name>/SKILL.md or similar).
func newSkillsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "skills [name]",
		Short: "Print an embedded authoring skill (markdown); no args = list available skills",
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := skills.FS.ReadDir(".")
			if err != nil {
				return err
			}
			names := make([]string, 0, len(entries))
			for _, e := range entries {
				names = append(names, strings.TrimSuffix(e.Name(), ".md"))
			}
			sort.Strings(names)
			if len(args) == 0 {
				for _, n := range names {
					fmt.Fprintln(cmd.OutOrStdout(), n)
				}
				return nil
			}
			raw, err := skills.FS.ReadFile(args[0] + ".md")
			if err != nil {
				return fmt.Errorf("unknown skill %q (available: %s)", args[0], strings.Join(names, ", "))
			}
			fmt.Fprint(cmd.OutOrStdout(), string(raw))
			return nil
		},
	}
}
