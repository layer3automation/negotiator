package args

import (
	"github.com/c-robinson/iplib"
)

// CIDRArgs implements the flag.Value interface and allows to parse strings expressing resource quantities.
type CIDRArgs struct {
	CIDRs []iplib.Net4
}

// String returns the stringified map entries.
func (c *CIDRArgs) String() string {
	return ""
}

// Set parses the provided string as a net and put it in the array.
func (c *CIDRArgs) Set(str string) error {
	if c.CIDRs == nil {
		c.CIDRs = make([]iplib.Net4, 0)
	}

	ip, net4, err := iplib.ParseCIDR(str)
	if err != nil {
		return err
	}

	ones, _ := net4.Mask().Size()
	n := iplib.NewNet4(ip, ones)
	c.CIDRs = append(c.CIDRs, n)

	return nil
}

// Type return the type name.
func (c *CIDRArgs) Type() string {
	return "ipNets"
}
