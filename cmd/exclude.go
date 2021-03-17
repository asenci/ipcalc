package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var excludeCmd = &cobra.Command{
	Aliases:       []string{"exc"},
	Use:           "exclude <prefix> <excluded>",
	Short:         "Exclude the specified prefix from the main prefix",
	SilenceErrors: true,
	Args:          cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var exc []string

		_, p := ParseCIDR(args[0])
		if p == nil {
			return fmt.Errorf("invalid prefix: %s", args[0])
		}

		_, q := ParseCIDR(args[1])
		if q == nil {
			return fmt.Errorf("invalid prefix: %s", args[1])
		}

		for _, sp := range p.Exclude(q) {
			exc = append(exc, sp.String())
		}

		fmt.Printf("> %s%s%s\n", Purple, strings.Join(exc, ", "), Reset)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(excludeCmd)
}
