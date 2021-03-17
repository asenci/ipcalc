package cmd

import (
	"fmt"
	"net"
	"strings"

	"github.com/mikioh/ipaddr"
	"github.com/spf13/cobra"
)

var summarizeCmd = &cobra.Command{
	Aliases:       []string{"sum"},
	Use:           "summarize <first address> <last address>",
	Short:         "Summarize the specified address range",
	SilenceErrors: true,
	Args:          cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var sum []string

		first := net.ParseIP(args[0])
		if first == nil {
			return fmt.Errorf("invalid address: %s", args[0])
		}

		last := net.ParseIP(args[1])
		if last == nil {
			return fmt.Errorf("invalid address: %s", args[1])
		}

		for _, p := range ipaddr.Summarize(first, last) {
			sum = append(sum, p.String())
		}

		if len(sum) == 0 {
			return fmt.Errorf("unable to summarize the specified address range")
		}

		fmt.Printf("> %s%s%s\n", Purple, strings.Join(sum, ", "), Reset)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(summarizeCmd)
}
