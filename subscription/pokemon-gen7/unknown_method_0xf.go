// Package protocol implements the Pokemon GEN 7 Subscription protocol
package protocol

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-go/v2/globals"
	subscription_types "github.com/PretendoNetwork/nex-protocols-go/v2/subscription/types"
)

func (protocol *Protocol) handleUnknownMethod0xF(packet nex.PacketInterface) {
	if protocol.UnknownMethod0xF == nil {
		err := nex.NewError(nex.ResultCodes.Core.NotImplemented, "Subscription::UnknownMethod0xF not implemented")

		globals.Logger.Warning(err.Message)
		globals.RespondError(packet, ProtocolID, err)

		return
	}

	request := packet.RMCMessage()
	callID := request.CallID
	parameters := request.Parameters
	endpoint := packet.Sender().Endpoint()
	parametersStream := nex.NewByteStreamIn(parameters, endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	var unknown types.UInt32
	var param subscription_types.SubscriptionData

	err := unknown.ExtractFrom(parametersStream)
	if err != nil {
		_, rmcError := protocol.UnknownMethod0xF(fmt.Errorf("failed to read unknown from parameters. %s", err.Error()), packet, callID, unknown, param)
		if rmcError != nil {
			globals.RespondError(packet, ProtocolID, rmcError)
		}

		return
	}

	err = param.ExtractFrom(parametersStream)
	if err != nil {
		_, rmcError := protocol.UnknownMethod0xF(fmt.Errorf("failed to read param from parameters. %s", err.Error()), packet, callID, unknown, param)
		if rmcError != nil {
			globals.RespondError(packet, ProtocolID, rmcError)
		}

		return
	}

	rmcMessage, rmcError := protocol.UnknownMethod0xF(nil, packet, callID, unknown, param)
	if rmcError != nil {
		globals.RespondError(packet, ProtocolID, rmcError)
		return
	}

	globals.Respond(packet, rmcMessage)
}
