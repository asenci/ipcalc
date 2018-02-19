package cmd

import (
	"fmt"
	"math/big"
	"net"
	"strings"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <prefix> [<prefix>...]",
	Short: "Show prefix details",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, s := range args {
			_, prefix, err := ParseCIDR(s)
			if err != nil {
				return err
			}

			maskAddress := net.IP(prefix.Mask)

			wildcardAddress := make(net.IP, len(prefix.Mask))
			for i, n := range prefix.Mask {
				wildcardAddress[i] = ^n
			}

			ones, bits := prefix.Mask.Size()
			maxAddresses := (&big.Int{}).Exp(big.NewInt(2), big.NewInt(int64(bits-ones)), nil)

			firstPlusMax := (&big.Int{}).SetBytes(prefix.IP)
			firstPlusMax.Sub(firstPlusMax, big.NewInt(1)).Add(firstPlusMax, maxAddresses)

			firstPlusMaxBytes := firstPlusMax.Bytes()
			lastAddress := make(net.IP, len(prefix.IP)-len(firstPlusMaxBytes))
			lastAddress = append(lastAddress, firstPlusMaxBytes...)

			fmt.Println("Input:", s)
			fmt.Println("  Prefix:", prefix)
			fmt.Println("  Netmask:", maskAddress)
			fmt.Println("  Wildcard:", wildcardAddress)
			fmt.Println()
			fmt.Println("  Number of addresses:", maxAddresses.String())
			fmt.Println("  First:", prefix.IP)
			fmt.Println("  Last:", lastAddress)

			//padding := 15 - len(ip.String())
			//for i := 0; i < padding; i++ {
			//	fmt.Printf(" ")
			//}
			//i, _ := prefix.Mask.Size()
			//fmt.Println(Bin(Compress(ip), i))

		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func Compress(ip net.IP) []byte {

	if p4 := ip.To4(); p4 != nil {
		return p4
	}

	return ip
}

func Bin(ip []byte, split int) (string, string) {
	var s string
	var o []string

	for _, b := range ip {
		s += fmt.Sprintf("%08b", b)
	}

	switch len(ip) {
	case net.IPv4len:
		o = make([]string, 0, 4)
		for i, j := 0, 8; i < len(s); {
			o = append(o, s[i:j])
			i += 8
			j += 8
		}
		s = strings.Join(o, ".")
		split += split / 8
	case net.IPv6len:
		o = make([]string, 0, 8)
		for i, j := 0, 16; i < len(s); {
			o = append(o, s[i:j])
			i += 16
			j += 16
		}
		s = strings.Join(o, ":")
		split += split / 16
	}

	return s[:split], s[split:]
}

func ParseCIDR(s string) (net.IP, *net.IPNet, error) {
	ip, prefix, err := net.ParseCIDR(s)
	if err != nil {
		if _, ok := err.(*net.ParseError); ok {
			ip, prefix = ParseNonIntMask(s)
			if ip == nil || prefix == nil {
				return nil, nil, err
			}
		} else {
			return nil, nil, err
		}
	}

	return Compress(ip), prefix, nil
}

func ParseNonIntMask(s string) (net.IP, *net.IPNet) {
	hasMask := strings.Index(s, "/")
	if hasMask == -1 {
		return nil, nil
	}

	ss := strings.SplitN(s, "/", 2)

	ip := net.ParseIP(ss[0])
	if ip == nil {
		return nil, nil
	}

	maskAddress := net.ParseIP(ss[1])
	if maskAddress == nil {
		return nil, nil
	}

	mask := net.IPMask(Compress(maskAddress))

	prefix := &net.IPNet{
		IP: Compress(ip.Mask(mask)),
		Mask: mask,
	}

	return ip, prefix

	return nil, nil
}
