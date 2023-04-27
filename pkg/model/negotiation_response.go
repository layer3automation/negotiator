package model

import (
	"encoding/json"
	"net"

	iputils "github.com/LucaRocco/l3_negotiator/pkg/utils"
	"github.com/c-robinson/iplib"
)

type NegotiationResponse struct {
	Net                iplib.Net4
	FreeIp             net.IP
	AssignedIp         net.IP
	DestinationNetwork *iplib.Net4
}

func (nr *NegotiationResponse) MarshalJSON() ([]byte, error) {
	str := &struct {
		Net                string  `json:"net"`
		FreeIp             string  `json:"freeIp"`
		AssignedIp         string  `json:"assignedIp"`
		DestinationNetwork *string `json:"destinationNetwork,omitempty"`
	}{
		Net:        nr.Net.String(),
		FreeIp:     nr.FreeIp.String(),
		AssignedIp: nr.AssignedIp.String(),
	}

	if nr.DestinationNetwork != nil {
		net4String := iputils.NetToString(*nr.DestinationNetwork)
		str.DestinationNetwork = &net4String
	}
	return json.Marshal(str)
}

func (nr *NegotiationResponse) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Net                string `json:"net"`
		FreeIp             string `json:"freeIp"`
		AssignedIp         string `json:"assignedIp"`
		DestinationNetwork string `json:"destinationNetwork,omitempty"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	nr.Net = iplib.Net4FromStr(aux.Net)
	nr.FreeIp = net.ParseIP(aux.FreeIp)
	nr.AssignedIp = net.ParseIP(aux.AssignedIp)
	net4 := iplib.Net4FromStr(aux.DestinationNetwork)
	nr.DestinationNetwork = &net4

	return nil
}
