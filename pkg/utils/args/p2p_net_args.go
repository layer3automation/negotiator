package args

import (
	"github.com/c-robinson/iplib"
)

// P2PNetArgs implements the flag.Value interface and allows to parse strings expressing resource quantities.
type P2PNetArgs struct {
	P2PNets []iplib.Net4
}

// String returns the stringified map entries.
func (c *P2PNetArgs) String() string {
	return ""
}

// Set parses the provided string as a net and put it in the array.
func (c *P2PNetArgs) Set(str string) error {
	if c.P2PNets == nil {
		c.P2PNets = make([]iplib.Net4, 0)
	}

	ip, net4, err := iplib.ParseCIDR(str)
	if err != nil {
		return err
	}

	ones, _ := net4.Mask().Size()
	n := iplib.NewNet4(ip, ones)

	nets, err := n.Subnet(30)
	if err != nil {
		return err
	}
	c.P2PNets = append(c.P2PNets, nets...)

	return nil
}

// Type return the type name.
func (c *P2PNetArgs) Type() string {
	return "ipNets"
}
