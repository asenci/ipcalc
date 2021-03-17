package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var overlapCmd = &cobra.Command{
	Aliases:       []string{"olp"},
	Use:           "overlap <prefix> <overlapping>",
	Short:         "Check if the specified prefix overlaps with the main prefix",
	SilenceErrors: true,
	SilenceUsage:  true,
	Args:          cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, p := ParseCIDR(args[0])
		if p == nil {
			return fmt.Errorf("invalid prefix: %s", args[0])
		}

		_, q := ParseCIDR(args[1])
		if q == nil {
			return fmt.Errorf("invalid prefix: %s", args[1])
		}

		if p.Overlaps(q) {
			return fmt.Errorf("prefix %s overlaps with prefix %s", q, p)
		}

		fmt.Printf("prefix %s does not overlap with prefix %s\n", q, p)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(overlapCmd)
}
