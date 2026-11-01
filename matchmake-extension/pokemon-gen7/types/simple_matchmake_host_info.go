// Package types implements all the types used by the Matchmake Extension (Pokemon GEN 7) protocol
package types

import (
	"fmt"
	"strings"

	"github.com/PretendoNetwork/nex-go/v2/types"
)

// SimpleMatchmakeHostInfo is a type within the Matchmake Extension (Pokemon GEN 7) protocol
type SimpleMatchmakeHostInfo struct {
	types.Structure
	PID         types.PID
	SessionKey  types.Buffer
	StationURLs types.List[types.StationURL]
}

// WriteTo writes the SimpleMatchmakeHostInfo to the given writable
func (smhi SimpleMatchmakeHostInfo) WriteTo(writable types.Writable) {
	contentWritable := writable.CopyNew()

	smhi.PID.WriteTo(contentWritable)
	smhi.SessionKey.WriteTo(contentWritable)
	smhi.StationURLs.WriteTo(contentWritable)

	content := contentWritable.Bytes()

	smhi.WriteHeaderTo(writable, uint32(len(content)))

	writable.Write(content)
}

// ExtractFrom extracts the SimpleMatchmakeHostInfo from the given readable
func (smhi *SimpleMatchmakeHostInfo) ExtractFrom(readable types.Readable) error {
	if err := smhi.ExtractHeaderFrom(readable); err != nil {
		return fmt.Errorf("failed to extract SimpleMatchmakeHostInfo header. %s", err.Error())
	}

	if err := smhi.PID.ExtractFrom(readable); err != nil {
		return fmt.Errorf("failed to extract SimpleMatchmakeHostInfo.PID. %s", err.Error())
	}

	if err := smhi.SessionKey.ExtractFrom(readable); err != nil {
		return fmt.Errorf("failed to extract SimpleMatchmakeHostInfo.SessionKey. %s", err.Error())
	}

	if err := smhi.StationURLs.ExtractFrom(readable); err != nil {
		return fmt.Errorf("failed to extract SimpleMatchmakeHostInfo.StationURLs. %s", err.Error())
	}

	return nil
}

// Copy returns a new copied instance of SimpleMatchmakeHostInfo
func (smhi SimpleMatchmakeHostInfo) Copy() types.RVType {
	copied := NewSimpleMatchmakeHostInfo()

	copied.StructureVersion = smhi.StructureVersion
	copied.PID = smhi.PID.Copy().(types.PID)
	copied.SessionKey = smhi.SessionKey.Copy().(types.Buffer)
	copied.StationURLs = smhi.StationURLs.Copy().(types.List[types.StationURL])

	return copied
}

// Equals checks if the given SimpleMatchmakeHostInfo contains the same data as the current SimpleMatchmakeHostInfo
func (smhi SimpleMatchmakeHostInfo) Equals(o types.RVType) bool {
	if _, ok := o.(SimpleMatchmakeHostInfo); !ok {
		return false
	}

	other := o.(SimpleMatchmakeHostInfo)

	if smhi.StructureVersion != other.StructureVersion {
		return false
	}

	if !smhi.PID.Equals(other.PID) {
		return false
	}

	if !smhi.SessionKey.Equals(other.SessionKey) {
		return false
	}

	return smhi.StationURLs.Equals(other.StationURLs)
}

// CopyRef copies the current value of the SimpleMatchmakeHostInfo
// and returns a pointer to the new copy
func (smhi SimpleMatchmakeHostInfo) CopyRef() types.RVTypePtr {
	copied := smhi.Copy().(SimpleMatchmakeHostInfo)
	return &copied
}

// Deref takes a pointer to the SimpleMatchmakeHostInfo
// and dereferences it to the raw value.
// Only useful when working with an instance of RVTypePtr
func (smhi *SimpleMatchmakeHostInfo) Deref() types.RVType {
	return *smhi
}

// String returns the string representation of the SimpleMatchmakeHostInfo
func (smhi SimpleMatchmakeHostInfo) String() string {
	return smhi.FormatToString(0)
}

// FormatToString pretty-prints the SimpleMatchmakeHostInfo using the provided indentation level
func (smhi SimpleMatchmakeHostInfo) FormatToString(indentationLevel int) string {
	indentationValues := strings.Repeat("\t", indentationLevel+1)
	indentationEnd := strings.Repeat("\t", indentationLevel)

	var b strings.Builder

	b.WriteString("SimpleMatchmakeHostInfo{\n")
	fmt.Fprintf(&b, "%sPID: %s,\n", indentationValues, smhi.PID)
	fmt.Fprintf(&b, "%sSessionKey: %s,\n", indentationValues, smhi.SessionKey)
	fmt.Fprintf(&b, "%sStationURLs: %s,\n", indentationValues, smhi.StationURLs)
	fmt.Fprintf(&b, "%s}", indentationEnd)

	return b.String()
}

// NewSimpleMatchmakeHostInfo returns a new SimpleMatchmakeHostInfo
func NewSimpleMatchmakeHostInfo() SimpleMatchmakeHostInfo {
	return SimpleMatchmakeHostInfo{
		PID:         types.NewPID(0),
		SessionKey:  types.NewBuffer(nil),
		StationURLs: types.NewList[types.StationURL](),
	}
}
