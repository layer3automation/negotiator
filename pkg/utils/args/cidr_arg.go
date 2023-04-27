package args

import (
	"github.com/c-robinson/iplib"
)

// CIDRArg implements the flag.Value interface and allows to parse strings expressing resource quantities.
type CIDRArg struct {
	CIDR *iplib.Net4
}

// String returns the stringified map entries.
func (c *CIDRArg) String() string {
	if c.CIDR != nil {
		return c.CIDR.String()
	}
	return ""
}

// Set parses the provided string as a net and put it in the array.
func (c *CIDRArg) Set(str string) error {
	ip, net4, err := iplib.ParseCIDR(str)
	if err != nil {
		return err
	}

	ones, _ := net4.Mask().Size()
	n := iplib.NewNet4(ip, ones)
	c.CIDR = &n

	return nil
}

// Type return the type name.
func (c *CIDRArg) Type() string {
	return "ipNets"
}
