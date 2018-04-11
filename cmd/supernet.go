package cmd

import (

"fmt"

"github.com/mikioh/ipaddr"
"github.com/spf13/cobra"

)

var supernetCmd = &cobra.Command{
	Aliases: []string{"sup"},
	Use:   "supernet <prefix> [<prefix>...]",
	Short: "Find the shortest common prefix for the specified prefixes",
	SilenceErrors: true,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var ps []ipaddr.Prefix

		for _, s := range args {
			_, p := ParseCIDR(s)
			if p == nil {
				return fmt.Errorf("invalid prefix: %s", s)
			}

			ps = append(ps, *p)
		}

		fmt.Printf("Prefix: %s\n", ipaddr.Supernet(ps))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(supernetCmd)
}
