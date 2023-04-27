package model

import (
	"encoding/json"

	utils "github.com/LucaRocco/l3_negotiator/pkg/utils"
	"github.com/c-robinson/iplib"
)

type NegotiationRequest struct {
	Cidrs              []iplib.Net4 `json:"cidrs"`
	DestinationNetwork *iplib.Net4  `json:"destinationNetwork"`
}

func NewNegotiationRequest(cidrs []iplib.Net4, destinationNetwork *iplib.Net4) *NegotiationRequest {
	nr := &NegotiationRequest{
		Cidrs: cidrs,
	}
	if destinationNetwork != nil {
		nr.DestinationNetwork = destinationNetwork
	}

	return nr
}

func (nr *NegotiationRequest) MarshalJSON() ([]byte, error) {
	str := &struct {
		Cidrs              []string `json:"cidrs"`
		DestinationNetwork *string  `json:"destinationNetwork,omitempty"`
	}{
		Cidrs: utils.NetsToStrings(nr.Cidrs),
	}

	if nr.DestinationNetwork != nil {
		net4String := utils.NetToString(*nr.DestinationNetwork)
		str.DestinationNetwork = &net4String
	}
	return json.Marshal(str)
}

func (nr *NegotiationRequest) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Cidrs              []string `json:"cidrs"`
		DestinationNetwork *string  `json:"destinationNetwork,omitempty"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	nr.Cidrs = utils.StringsToNets(aux.Cidrs)
	if aux.DestinationNetwork != nil {
		net4 := utils.StringToNet(*aux.DestinationNetwork)
		nr.DestinationNetwork = &net4
	}

	return nil
}
