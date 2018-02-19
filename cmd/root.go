package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	verbose bool
)

var rootCmd = &cobra.Command{
	Version: "0.1a",
	Use:   "ipcalc",
	Short: "IP address/CIDR calculator",
	Long: `An IPv6-enabled CIDR calculator`,
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose,"verbose", "v", false, "Print additional information")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

