package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/LucaRocco/l3_negotiator/pkg/negotiator"
	argsutils "github.com/LucaRocco/l3_negotiator/pkg/utils/args"
	"github.com/spf13/cobra"
)

func main() {
	var addr string
	props := negotiator.NegotiatorProps{}
	var cidrs argsutils.P2PNetArgs
	var internalNetwork argsutils.CIDRArg
	var overlappingRemappingCIDRs argsutils.CIDRArgs

	var rootCmd = &cobra.Command{
		Use:          os.Args[0],
		Short:        "Layer 3 negotiator",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			props.Mutex.Lock()
			props.P2PNets = cidrs.P2PNets
			if internalNetwork.CIDR != nil {
				props.InternalNetwork = internalNetwork.CIDR
			}
			props.OverlappingRemappingCIDRs = overlappingRemappingCIDRs.CIDRs

			internalMaskOnes, internalMaskBits := props.InternalNetwork.Mask().Size()
			for _, overlappingCIDR := range props.OverlappingRemappingCIDRs {
				overlappingMaskOnes, overlappingMaskBits := overlappingCIDR.Mask().Size()
				if internalMaskOnes != overlappingMaskOnes || internalMaskBits != overlappingMaskBits {
					return fmt.Errorf("%s doesn't have the mask of the internal network", overlappingCIDR.String())
				}
			}
			props.Mutex.Unlock()

			// Run the web server
			return negotiator.SetupRouterAndServeHTTP(addr, context.Background(), &props)
		},
	}

	rootCmd.PersistentFlags().StringVar(&addr, "addr", "0.0.0.0:8000", "set addr in the form host:port where the server will listen on.")
	rootCmd.Flags().Var(&cidrs, "cidr", "set a free cidrs to be negotiated with the remote agent")
	rootCmd.Flags().Var(&internalNetwork, "internal_network", "set the network that you want to make reachable from the other side")
	rootCmd.Flags().StringVar(&props.Interface, "interface", "", "set the interface where the ips and routes will be configured")
	rootCmd.Flags().Var(&overlappingRemappingCIDRs, "remapping_cidr", "set cidrs to be used in case of internal network overlapping with the remote "+
		"internal network. The netmask of each overlapping CIDR must be equals to the internal's one")

	err := rootCmd.Root().MarkFlagRequired("cidr")
	if err != nil {
		log.Fatalf("error: error during marking cidr flag as required: %s", err)
	}

	err = rootCmd.Root().MarkFlagRequired("remapping_cidr")
	if err != nil {
		log.Fatalf("error: error during marking remapping_cidrs flag as required: %s", err)
	}

	err = rootCmd.Root().MarkFlagRequired("interface")
	if err != nil {
		log.Fatalf("error: error during marking interface flag as required: %s", err)
	}

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
