// Package protocol implements the Pokemon GEN 7 Subscription protocol
package protocol

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-go/v2/globals"
	subscription "github.com/PretendoNetwork/nex-protocols-go/v2/subscription"
	subscription_types "github.com/PretendoNetwork/nex-protocols-go/v2/subscription/types"
	"golang.org/x/exp/slices"
)

const (
	// ProtocolID is the protocol ID for the Subscription (Pokemon GEN 7) protocol. ID is the same as the Subscription protocol
	ProtocolID = 0x75

	// MethodUnknownMethod0xE is the method ID for the method UnknownMethod0xE
	MethodUnknownMethod0xE = 0xE

	// MethodUnknownMethod0xF is the method ID for the method UnknownMethod0xF
	MethodUnknownMethod0xF = 0xF
)

var patchedMethods = []uint32{
	MethodUnknownMethod0xE,
	MethodUnknownMethod0xF,
}

type subscriptionProtocol = subscription.Protocol

// Protocol stores all the RMC method handlers for the Subscription (Pokemon GEN 7) protocol and listens for requests
// Embeds the Subscription protocol
type Protocol struct {
	endpoint nex.EndpointInterface
	subscriptionProtocol
	UnknownMethod0xE func(err error, packet nex.PacketInterface, callID uint32, unknown types.UInt32, param subscription_types.SubscriptionData) (*nex.RMCMessage, *nex.Error)
	UnknownMethod0xF func(err error, packet nex.PacketInterface, callID uint32, unknown types.UInt32, param subscription_types.SubscriptionData) (*nex.RMCMessage, *nex.Error)
}

// HandlePacket sends the packet to the correct RMC method handler
func (protocol *Protocol) HandlePacket(packet nex.PacketInterface) {
	message := packet.RMCMessage()

	if !message.IsRequest || message.ProtocolID != ProtocolID {
		return
	}

	if !slices.Contains(patchedMethods, message.MethodID) {
		protocol.subscriptionProtocol.HandlePacket(packet)
		return
	}

	switch message.MethodID {
	case MethodUnknownMethod0xE:
		protocol.handleUnknownMethod0xE(packet)
	case MethodUnknownMethod0xF:
		protocol.handleUnknownMethod0xF(packet)
	default:
		errMessage := fmt.Sprintf("Unsupported Subscription (Pokemon GEN 7) method ID: %#v\n", message.MethodID)
		err := nex.NewError(nex.ResultCodes.Core.NotImplemented, errMessage)

		globals.RespondError(packet, ProtocolID, err)
		globals.Logger.Warning(err.Message)
	}
}

// NewProtocol returns a new Subscription (Pokemon GEN 7) protocol
func NewProtocol(endpoint nex.EndpointInterface) *Protocol {
	protocol := &Protocol{endpoint: endpoint}
	protocol.SetEndpoint(endpoint)

	return protocol
}
