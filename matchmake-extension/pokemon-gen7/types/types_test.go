package types

import (
	"bytes"
	"testing"

	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
)

// Expected byte layouts are written out by hand from the field tables at
// https://nintendo-wiki.pretendo.network/docs/nex/protocols/matchmake-extension/pokemon-generation-7
// so these tests check the wire format, not just encode/decode symmetry.

func streamSettings(structureHeader bool) *nex.ByteStreamSettings {
	settings := nex.NewByteStreamSettings()
	settings.UseStructureHeader = structureHeader

	return settings
}

func encode(t *testing.T, value types.RVType, structureHeader bool) []byte {
	t.Helper()

	stream := nex.NewByteStreamOut(nex.NewLibraryVersions(), streamSettings(structureHeader))
	value.WriteTo(stream)

	return stream.Bytes()
}

func newAttractionStatus() AttractionStatus {
	status := NewAttractionStatus()
	status.MessageInterval = types.NewUInt16(0x0102)
	status.OperationFlag = types.NewUInt8(0x03)
	status.ActivePlayerInviteParam = types.NewUInt16(0x0405)
	status.ActivePlayerJoinParam = types.NewUInt16(0x0607)
	status.ExtraParams = types.NewList[types.UInt32]()
	status.ExtraParams = append(status.ExtraParams, types.NewUInt32(0x0A0B0C0D), types.NewUInt32(1))

	return status
}

var attractionStatusBody = []byte{
	0x02, 0x01, // MessageInterval (u16 LE)
	0x03,       // OperationFlag (u8)
	0x05, 0x04, // ActivePlayerInviteParam (u16 LE)
	0x07, 0x06, // ActivePlayerJoinParam (u16 LE)
	0x02, 0x00, 0x00, 0x00, // ExtraParams count (u32 LE)
	0x0D, 0x0C, 0x0B, 0x0A, // ExtraParams[0]
	0x01, 0x00, 0x00, 0x00, // ExtraParams[1]
}

func TestAttractionStatusWireFormat(t *testing.T) {
	got := encode(t, newAttractionStatus(), false)

	if !bytes.Equal(got, attractionStatusBody) {
		t.Fatalf("encoded AttractionStatus mismatch\n got: % X\nwant: % X", got, attractionStatusBody)
	}
}

func TestAttractionStatusStructureHeader(t *testing.T) {
	got := encode(t, newAttractionStatus(), true)

	want := append([]byte{0x00, byte(len(attractionStatusBody)), 0x00, 0x00, 0x00}, attractionStatusBody...)
	if !bytes.Equal(got, want) {
		t.Fatalf("encoded AttractionStatus with header mismatch\n got: % X\nwant: % X", got, want)
	}

	decoded := NewAttractionStatus()
	if err := decoded.ExtractFrom(nex.NewByteStreamIn(got, nex.NewLibraryVersions(), streamSettings(true))); err != nil {
		t.Fatalf("ExtractFrom failed: %v", err)
	}

	if !decoded.Equals(newAttractionStatus()) {
		t.Fatalf("decoded AttractionStatus differs:\n%s", decoded.FormatToString(0))
	}
}

func TestAttractionStatusEmptyExtraParams(t *testing.T) {
	status := NewAttractionStatus()

	got := encode(t, status, false)
	want := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	if !bytes.Equal(got, want) {
		t.Fatalf("zero-value AttractionStatus mismatch\n got: % X\nwant: % X", got, want)
	}
}

func newHostInfo() SimpleMatchmakeHostInfo {
	info := NewSimpleMatchmakeHostInfo()
	info.PID = types.NewPID(0x11223344)
	info.SessionKey = types.NewBuffer([]byte{0xAA, 0xBB, 0xCC})
	info.StationURLs = types.NewList[types.StationURL]()
	info.StationURLs = append(info.StationURLs, types.NewStationURL("prudp:/address=1.2.3.4"))

	return info
}

func TestSimpleMatchmakeHostInfoWireFormat(t *testing.T) {
	got := encode(t, newHostInfo(), false)

	url := "prudp:/address=1.2.3.4"
	want := []byte{
		0x44, 0x33, 0x22, 0x11, // PID (u32 LE on 3DS)
		0x03, 0x00, 0x00, 0x00, 0xAA, 0xBB, 0xCC, // SessionKey: Buffer (u32 length + bytes)
		0x01, 0x00, 0x00, 0x00, // StationURLs count
		byte(len(url) + 1), 0x00, // StationURL: String (u16 length incl. NUL)
	}
	want = append(want, url...)
	want = append(want, 0x00)

	if !bytes.Equal(got, want) {
		t.Fatalf("encoded SimpleMatchmakeHostInfo mismatch\n got: % X\nwant: % X", got, want)
	}

	decoded := NewSimpleMatchmakeHostInfo()
	if err := decoded.ExtractFrom(nex.NewByteStreamIn(got, nex.NewLibraryVersions(), streamSettings(false))); err != nil {
		t.Fatalf("ExtractFrom failed: %v", err)
	}

	if !decoded.Equals(newHostInfo()) {
		t.Fatalf("decoded SimpleMatchmakeHostInfo differs:\n%s", decoded.FormatToString(0))
	}
}

// Malformed client packets must produce an error, never a panic or a partially-trusted value.
func TestTruncatedInputReturnsError(t *testing.T) {
	cases := map[string]struct {
		data   []byte
		target interface{ ExtractFrom(types.Readable) error }
	}{
		"AttractionStatus missing list":    {attractionStatusBody[:7], &AttractionStatus{}},
		"AttractionStatus short list":      {attractionStatusBody[:len(attractionStatusBody)-2], &AttractionStatus{}},
		"HostInfo buffer longer than data": {[]byte{1, 0, 0, 0, 0xFF, 0, 0, 0, 0xAA}, &SimpleMatchmakeHostInfo{}},
		"HostInfo empty":                   {[]byte{}, &SimpleMatchmakeHostInfo{}},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			stream := nex.NewByteStreamIn(tc.data, nex.NewLibraryVersions(), streamSettings(false))
			if err := tc.target.ExtractFrom(stream); err == nil {
				t.Fatal("expected an error for truncated input, got nil")
			}
		})
	}
}

func TestStructureHeaderLengthLongerThanData(t *testing.T) {
	data := []byte{0x00, 0xFF, 0x00, 0x00, 0x00, 0x01}

	err := (&AttractionStatus{}).ExtractFrom(nex.NewByteStreamIn(data, nex.NewLibraryVersions(), streamSettings(true)))
	if err == nil {
		t.Fatal("expected an error when structure header length exceeds data, got nil")
	}
}

func TestCopyIsDeep(t *testing.T) {
	original := newHostInfo()
	copied := original.Copy().(SimpleMatchmakeHostInfo)

	copied.StationURLs[0] = types.NewStationURL("prudp:/address=9.9.9.9")

	if original.Equals(copied) {
		t.Fatal("mutating the copy changed equality; Copy is not independent")
	}
	if original.StationURLs[0].URL() == copied.StationURLs[0].URL() {
		t.Fatal("Copy shares StationURLs backing array with original")
	}
}
