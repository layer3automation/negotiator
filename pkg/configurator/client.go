package configurator

import (
	"context"
	"log"
	"net"
	"strconv"

	configuratorClient "github.com/layer3automation/linux_configuration_agent/configurator"
	"google.golang.org/grpc"
)

type Configurator struct {
	client configuratorClient.ConfigurationAgentServiceClient
	ctx    context.Context
}

type NatConfiguration struct {
	LocalNetwork   string
	RemoteNetwork  string
	LocalInterface string
	TableNumber    string
	PacketMark     string
	NextHop        string
}

func NewConfigurator(ctx context.Context) *Configurator {
	conn, err := grpc.Dial("0.0.0.0:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to gRPC server: %s", err)
	}

	client := configuratorClient.NewConfigurationAgentServiceClient(conn)

	return &Configurator{
		client: client,
		ctx:    ctx,
	}
}

func (c *Configurator) AddIPToInterface(ip *net.IP, mask int, interfaceName string) (*configuratorClient.Result, error) {
	log.Printf("AssignIPToInterface called: %s, %s", ip, interfaceName)

	ipAssignment := configuratorClient.IPAssignment{
		Ip:        ip.String() + "/" + strconv.Itoa(mask),
		Interface: interfaceName,
	}
	return c.client.AddIPToInterface(c.ctx, &ipAssignment)
}

func (c *Configurator) AddRoute(destinationNetwork *net.IP, nextHop *net.IP) (*configuratorClient.Result, error) {
	log.Printf("AddRoute called: destination %s, nextHop %s", destinationNetwork, nextHop)

	route := configuratorClient.Route{
		DestinationNetwork: destinationNetwork.String(),
		NextHop:            nextHop.String(),
	}

	return c.client.AddRoute(c.ctx, &route)
}

func (c *Configurator) ConfigureNat(in *NatConfiguration) (*configuratorClient.Result, error) {
	log.Printf("ConfigureNat called: %+v", in)

	natConf := configuratorClient.NatConfiguration{
		LocalNetwork:   in.LocalNetwork,
		RemoteNetwork:  in.RemoteNetwork,
		LocalInterface: in.LocalInterface,
		TableNumber:    in.TableNumber,
		PacketMark:     in.PacketMark,
		NextHop:        in.NextHop,
	}
	return c.client.ConfigureNat(c.ctx, &natConf)
}
