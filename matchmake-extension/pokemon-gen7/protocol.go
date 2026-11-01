// Package protocol implements the Pokemon GEN 7 Matchmake Extension protocol
package protocol

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-go/v2/globals"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"
	"golang.org/x/exp/slices"
)

// Pokemon Sun/Moon/Ultra Sun/Ultra Moon reuse method IDs 0x30-0x35, which the base
// Matchmake Extension protocol assigns to other methods. Requests for these IDs are
// routed to the GEN 7 handlers below; all other IDs fall through to the base protocol.
const (
	// ProtocolID is the protocol ID for the Matchmake Extension (Pokemon GEN 7) protocol. ID is the same as the Matchmake Extension protocol
	ProtocolID = 0x6D

	// MethodClearMyPreviouslyMatchedUserCache is the method ID for the method ClearMyPreviouslyMatchedUserCache
	MethodClearMyPreviouslyMatchedUserCache = 0x30

	// MethodGetAttractionStatus is the method ID for the method GetAttractionStatus
	MethodGetAttractionStatus = 0x31

	// MethodGetAttractionStatusWithGroupID is the method ID for the method GetAttractionStatusWithGroupID
	MethodGetAttractionStatusWithGroupID = 0x32

	// MethodSimpleMatchmake is the method ID for the method SimpleMatchmake
	MethodSimpleMatchmake = 0x33

	// MethodEntrySimpleMatchmake is the method ID for the method EntrySimpleMatchmake
	MethodEntrySimpleMatchmake = 0x34

	// MethodCancelSimpleMatchmakeEntry is the method ID for the method CancelSimpleMatchmakeEntry
	MethodCancelSimpleMatchmakeEntry = 0x35
)

var patchedMethods = []uint32{
	MethodClearMyPreviouslyMatchedUserCache,
	MethodGetAttractionStatus,
	MethodGetAttractionStatusWithGroupID,
	MethodSimpleMatchmake,
	MethodEntrySimpleMatchmake,
	MethodCancelSimpleMatchmakeEntry,
}

type matchmakeExtensionProtocol = matchmake_extension.Protocol

// Protocol stores all the RMC method handlers for the Matchmake Extension (Pokemon GEN 7) protocol and listens for requests
// Embeds the Matchmake Extension protocol
type Protocol struct {
	endpoint nex.EndpointInterface
	matchmakeExtensionProtocol
	ClearMyPreviouslyMatchedUserCache func(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error)
	GetAttractionStatus               func(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error)
	GetAttractionStatusWithGroupID    func(err error, packet nex.PacketInterface, callID uint32, groupID types.UInt32) (*nex.RMCMessage, *nex.Error)
	SimpleMatchmake                   func(err error, packet nex.PacketInterface, callID uint32, groupID types.UInt32) (*nex.RMCMessage, *nex.Error)
	EntrySimpleMatchmake              func(err error, packet nex.PacketInterface, callID uint32, groupID types.UInt32, sessionKey types.Buffer) (*nex.RMCMessage, *nex.Error)
	CancelSimpleMatchmakeEntry        func(err error, packet nex.PacketInterface, callID uint32, groupID types.UInt32) (*nex.RMCMessage, *nex.Error)
}

// HandlePacket sends the packet to the correct RMC method handler
func (protocol *Protocol) HandlePacket(packet nex.PacketInterface) {
	message := packet.RMCMessage()

	if !message.IsRequest || message.ProtocolID != ProtocolID {
		return
	}

	if !slices.Contains(patchedMethods, message.MethodID) {
		protocol.matchmakeExtensionProtocol.HandlePacket(packet)
		return
	}

	switch message.MethodID {
	case MethodClearMyPreviouslyMatchedUserCache:
		protocol.handleClearMyPreviouslyMatchedUserCache(packet)
	case MethodGetAttractionStatus:
		protocol.handleGetAttractionStatus(packet)
	case MethodGetAttractionStatusWithGroupID:
		protocol.handleGetAttractionStatusWithGroupID(packet)
	case MethodSimpleMatchmake:
		protocol.handleSimpleMatchmake(packet)
	case MethodEntrySimpleMatchmake:
		protocol.handleEntrySimpleMatchmake(packet)
	case MethodCancelSimpleMatchmakeEntry:
		protocol.handleCancelSimpleMatchmakeEntry(packet)
	default:
		errMessage := fmt.Sprintf("Unsupported MatchmakeExtension (Pokemon GEN 7) method ID: %#v\n", message.MethodID)
		err := nex.NewError(nex.ResultCodes.Core.NotImplemented, errMessage)

		globals.RespondError(packet, ProtocolID, err)
		globals.Logger.Warning(err.Message)
	}
}

// NewProtocol returns a new Matchmake Extension (Pokemon GEN 7) protocol
func NewProtocol(endpoint nex.EndpointInterface) *Protocol {
	protocol := &Protocol{endpoint: endpoint}
	protocol.SetEndpoint(endpoint)

	return protocol
}
