package protocol

import (
	"testing"

	nex "github.com/PretendoNetwork/nex-go/v2"
)

// fakePacket carries only an RMC message. Handlers under test are replaced with
// callbacks that panic with a marker before any response is sent, so the test can
// observe which handler a method ID was routed to without a live connection.
type fakePacket struct {
	message *nex.RMCMessage
}

func (p *fakePacket) Sender() nex.ConnectionInterface       { return nil }
func (p *fakePacket) Payload() []byte                       { return nil }
func (p *fakePacket) SetPayload([]byte)                     {}
func (p *fakePacket) RMCMessage() *nex.RMCMessage           { return p.message }
func (p *fakePacket) SetRMCMessage(message *nex.RMCMessage) { p.message = message }

type routedTo string

func request(protocolID uint16, methodID uint32) *fakePacket {
	message := nex.NewRMCRequest(nil)
	message.ProtocolID = protocolID
	message.MethodID = methodID

	return &fakePacket{message: message}
}

func marker(name string) func(error, nex.PacketInterface, uint32) (*nex.RMCMessage, *nex.Error) {
	return func(error, nex.PacketInterface, uint32) (*nex.RMCMessage, *nex.Error) { panic(routedTo(name)) }
}

func newMarkedProtocol() *Protocol {
	protocol := NewProtocol(nil)

	protocol.ClearMyPreviouslyMatchedUserCache = marker("gen7.ClearMyPreviouslyMatchedUserCache")
	protocol.GetAttractionStatus = marker("gen7.GetAttractionStatus")

	// Base-protocol handlers for method IDs that GEN 7 does not override.
	protocol.ClearMyBlockList = marker("base.ClearMyBlockList")
	protocol.GetMyBlockList = marker("base.GetMyBlockList")

	return protocol
}

func dispatch(t *testing.T, protocol *Protocol, packet *fakePacket) (got routedTo) {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			marker, ok := r.(routedTo)
			if !ok {
				t.Fatalf("unexpected panic: %v", r)
			}
			got = marker
		}
	}()

	protocol.HandlePacket(packet)

	return ""
}

func TestGen7MethodsOverrideCollidingBaseIDs(t *testing.T) {
	protocol := newMarkedProtocol()

	cases := map[uint32]routedTo{
		0x30: "gen7.ClearMyPreviouslyMatchedUserCache", // base: FindMatchmakeSessionByGatheringID
		0x31: "gen7.GetAttractionStatus",               // base: FindMatchmakeSessionBySingleGatheringID
	}

	for methodID, want := range cases {
		if got := dispatch(t, protocol, request(ProtocolID, methodID)); got != want {
			t.Errorf("method %#x routed to %q, want %q", methodID, got, want)
		}
	}
}

func TestNonGen7MethodsFallThroughToBase(t *testing.T) {
	protocol := newMarkedProtocol()

	cases := map[uint32]routedTo{
		0x18: "base.GetMyBlockList",
		0x1B: "base.ClearMyBlockList",
	}

	for methodID, want := range cases {
		if got := dispatch(t, protocol, request(ProtocolID, methodID)); got != want {
			t.Errorf("method %#x routed to %q, want %q", methodID, got, want)
		}
	}
}

func TestIgnoresOtherProtocolsAndResponses(t *testing.T) {
	protocol := newMarkedProtocol()

	if got := dispatch(t, protocol, request(ProtocolID+1, MethodGetAttractionStatus)); got != "" {
		t.Errorf("packet for another protocol was handled by %q", got)
	}

	response := request(ProtocolID, MethodGetAttractionStatus)
	response.message.IsRequest = false
	if got := dispatch(t, protocol, response); got != "" {
		t.Errorf("response packet was handled by %q", got)
	}
}

// Method IDs come from the Pretendo wiki (Pokemon Generation 7 matchmake-extension page).
func TestMethodIDsMatchDocumentation(t *testing.T) {
	want := map[string]uint32{
		"ClearMyPreviouslyMatchedUserCache": 48,
		"GetAttractionStatus":               49,
		"GetAttractionStatusWithGroupID":    50,
		"SimpleMatchmake":                   51,
		"EntrySimpleMatchmake":              52,
		"CancelSimpleMatchmakeEntry":        53,
	}
	got := map[string]uint32{
		"ClearMyPreviouslyMatchedUserCache": MethodClearMyPreviouslyMatchedUserCache,
		"GetAttractionStatus":               MethodGetAttractionStatus,
		"GetAttractionStatusWithGroupID":    MethodGetAttractionStatusWithGroupID,
		"SimpleMatchmake":                   MethodSimpleMatchmake,
		"EntrySimpleMatchmake":              MethodEntrySimpleMatchmake,
		"CancelSimpleMatchmakeEntry":        MethodCancelSimpleMatchmakeEntry,
	}

	for name, id := range want {
		if got[name] != id {
			t.Errorf("%s = %d, want %d", name, got[name], id)
		}
	}
	if len(patchedMethods) != len(want) {
		t.Errorf("patchedMethods has %d entries, want %d", len(patchedMethods), len(want))
	}
}
