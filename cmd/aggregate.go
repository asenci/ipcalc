package cmd

import (
	"fmt"
	"strings"

	"github.com/mikioh/ipaddr"
	"github.com/spf13/cobra"
)

var aggregateCmd = &cobra.Command{
	Aliases: []string{"agg"},
	Use:   "aggregate <prefix> [<prefix>...]",
	Short: "Aggregate the specified prefixes",
	SilenceErrors: true,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var ps []ipaddr.Prefix
		var agg []string

		for _, s := range args {
			_, p := ParseCIDR(s)
			if p == nil {
				return fmt.Errorf("invalid prefix: %s", s)
			}

			ps = append(ps, *p)
		}

		for _, p := range ipaddr.Aggregate(ps) {
			agg = append(agg, p.String())
		}

		fmt.Printf("Prefixes: %s\n", strings.Join(agg, ", "))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(aggregateCmd)
}
