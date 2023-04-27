package negotiator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/LucaRocco/l3_negotiator/pkg/configurator"
	model "github.com/LucaRocco/l3_negotiator/pkg/model"
	"github.com/LucaRocco/l3_negotiator/pkg/utils"
	"github.com/c-robinson/iplib"
)

type NegotiatorProps struct {
	Mutex                     sync.RWMutex
	P2PNets                   []iplib.Net4
	InternalNetwork           *iplib.Net4
	Interface                 string
	OverlappingRemappingCIDRs []iplib.Net4
}

type Negotiation struct {
	RemoteAgent   string `json:"remoteAgent"`
	RemoteNetwork string `json:"remoteNetwork"`
}

type Negotiator struct {
	props        *NegotiatorProps
	ctx          context.Context
	configurator *configurator.Configurator
}

func NewNegotiator(ctx context.Context, props *NegotiatorProps) *Negotiator {
	return &Negotiator{
		props:        props,
		configurator: configurator.NewConfigurator(ctx),
		ctx:          ctx,
	}
}

func (n *Negotiator) startNegotiation(negotiation Negotiation) error {
	log.Printf("Negotiation request received: %+v", negotiation)

	n.props.Mutex.RLock()
	js, err := json.Marshal(model.NewNegotiationRequest(n.props.P2PNets, n.props.InternalNetwork))
	n.props.Mutex.RUnlock()

	if err != nil {
		return err
	}
	log.Printf("Asking negotitation to [%s] with payload %s", negotiation.RemoteAgent, string(js[:]))
	res, err := http.Post(negotiation.RemoteAgent, "application/json", bytes.NewBuffer(js))
	if err != nil {
		return err
	}

	log.Printf("Response received from [%s] with status %d", negotiation.RemoteAgent, res.StatusCode)
	if res.StatusCode == http.StatusOK {
		nr := model.NegotiationResponse{}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		log.Printf("payload: %s", string(body[:]))
		json.Unmarshal(body, &nr)

		//TODO: delete the chosen /30

		n.configurator.AddIPToInterface(&nr.FreeIp, 30, n.props.Interface)
		destinationNetwork := nr.DestinationNetwork.IP()
		n.configurator.AddRoute(&destinationNetwork, &nr.AssignedIp)
	}
	return nil
}

func (n *Negotiator) handleNegotiation(nr model.NegotiationRequest) (*model.NegotiationResponse, error) {
	var commonNetwork *iplib.Net4

	for _, remoteP2P := range nr.Cidrs {
		if utils.Contains(n.props.P2PNets, remoteP2P) {
			commonNetwork = &remoteP2P
			break
		}
	}

	if commonNetwork == nil {
		return nil, fmt.Errorf("No common network found for the P2P network")
	}

	ipToAssignLocally := commonNetwork.FirstAddress()
	ipToSendRemotely := commonNetwork.LastAddress()
	n.configurator.AddIPToInterface(&ipToAssignLocally, 30, n.props.Interface)
	response := &model.NegotiationResponse{
		Net:        *commonNetwork,
		FreeIp:     ipToSendRemotely,
		AssignedIp: ipToAssignLocally,
	}

	if nr.DestinationNetwork != nil && nr.DestinationNetwork.String() == n.props.InternalNetwork.String() {
		log.Printf("Handling overlapping. Both sides have the same internal network: %s", nr.DestinationNetwork.String())
		var remappingNetwork iplib.Net4
		if len(n.props.OverlappingRemappingCIDRs) > 0 {
			remappingNetwork = n.props.OverlappingRemappingCIDRs[0]
			log.Printf("remapping internal network with %s", remappingNetwork.String())
		} else {
			return nil, fmt.Errorf("No remapping network left")
		}

		n.configurator.ConfigureNat(&configurator.NatConfiguration{
			LocalNetwork:   remappingNetwork.String(),
			RemoteNetwork:  nr.DestinationNetwork.String(),
			LocalInterface: n.props.Interface,
			PacketMark:     "0x100",
			TableNumber:    "100",
			NextHop:        ipToSendRemotely.String(),
		})

		response.DestinationNetwork = &remappingNetwork
	} else {
		destinationNetwork := nr.DestinationNetwork.IP()
		n.configurator.AddRoute(&destinationNetwork, &ipToSendRemotely)
		response.DestinationNetwork = n.props.InternalNetwork
	}

	return response, nil
}
