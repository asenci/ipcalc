package cmd

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/mikioh/ipaddr"
	"github.com/spf13/cobra"
)

var (
	compressIPv6 bool
	verbose bool
)

var rootCmd = &cobra.Command{
	Version: "0.1a",
	Use:   "ipcalc [flags] <prefix> [prefix...]",
	Long: `ipcalc - IPv6-enabled CIDR calculator

Default action is to show the prefixes details`,
	DisableFlagsInUseLine: true,
	SilenceErrors: true,
	//SilenceUsage:  true,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sb := &strings.Builder{}

		for i, s := range args {
			ip, p := ParseCIDR(s)
			if ip == nil || p == nil {
				return fmt.Errorf("invalid prefix: %s", s)
			}

			if i != 0 {
				sb.WriteByte('\n')
			}

			fmt.Fprintf(sb, "Prefix: %s\n", p)
			fmt.Fprintf(sb, "  Number of addresses: %d\n", p.NumNodes())
			sb.WriteByte('\n')
			fmt.Fprintf(sb, "  Netmask:  %s\n", Explode(p.Mask))
			fmt.Fprintf(sb, "  Wildcard: %s\n", Explode(p.Hostmask()))
			sb.WriteByte('\n')
			fmt.Fprintf(sb, "  First:    %s\n", Explode(p.IP))
			if ! bytes.Equal(ip, p.IP) {
				fmt.Fprintf(sb, "  Input:    %s\n", Explode(ip))
			}
			fmt.Fprintf(sb, "  Last:     %s\n", Explode(p.Last()))

			if verbose {
				maskLen, _ := p.Mask.Size()

				sb.WriteByte('\n')
				fmt.Fprintf(sb, "  First:    %s\n", Bin(p.IP, maskLen))
				if ! bytes.Equal(ip, p.IP) {
					fmt.Fprintf(sb, "  Input:    %s\n", Bin(ip, maskLen))
				}
				fmt.Fprintf(sb, "  Last:     %s\n", Bin(p.Last(), maskLen))
			}
		}

		fmt.Print(sb)
		return nil
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&compressIPv6, "compress", "c", false, "zero-compress IPv6 addresses using the \"::\" notation")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print additional information")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Bin(ip []byte, split int) string {
	var o []rune

	sb := &strings.Builder{}
	for _, b := range ip {
		fmt.Fprintf(sb, "%08b", b)
	}
	s := sb.String()

	nip := net.IP(ip)
	if nip.To4() != nil {
		for i, c := range s[len(s)-32:] {
			if i != 0 {
				if i%8 == 0 {
					o = append(o, '.')
				}
				if i == split {
					o = append(o, ' ')
				}
			}
			o = append(o, c)
		}
	} else {
		for i, c := range s {
			if i != 0 {
				if i%16 == 0 {
					o = append(o, ':')
				}
				if i == split {
					o = append(o, ' ')
				}
			}
			o = append(o, c)
		}
	}

	return string(o)
}

func ParseCIDR(s string) (net.IP, *ipaddr.Prefix) {
	cidrMaskRegex := regexp.MustCompile(`/\d+$`)

	if ! cidrMaskRegex.MatchString(s) {
		i := strings.Index(s, "/")
		if i < 0 {
			return nil, nil
		}

		ip := s[:i]
		mask := s[i+1:]

		var m []byte
		m = net.ParseIP(mask)
		if m == nil {
			return nil, nil
		}
		if v4 := net.IP(m).To4(); v4 != nil {
			m = v4
		}

		ones, _ := net.IPMask(m).Size()

		s = fmt.Sprintf("%s/%d", ip, ones)
	}

	ip, n, err := net.ParseCIDR(s)
	if err != nil {
		return nil, nil
	}
	return ip, ipaddr.NewPrefix(n)
}

func Explode(ip []byte) string {
	nip := net.IP(ip)

	if compressIPv6 || nip.To4() != nil {
		return nip.String()
	}

	sb := &strings.Builder{}
	for i := 0; i < net.IPv6len; i += 2 {
		if i != 0 {
			sb.WriteByte(':')
		}

		fmt.Fprintf(sb, "%02x%02x", ip[i], ip[i+1])
	}

	return sb.String()
}

