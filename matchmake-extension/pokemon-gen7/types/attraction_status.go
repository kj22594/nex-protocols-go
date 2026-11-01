// Package types implements all the types used by the Matchmake Extension (Pokemon GEN 7) protocol
package types

import (
	"fmt"
	"strings"

	"github.com/PretendoNetwork/nex-go/v2/types"
)

// AttractionStatus is a type within the Matchmake Extension (Pokemon GEN 7) protocol
type AttractionStatus struct {
	types.Structure
	MessageInterval         types.UInt16
	OperationFlag           types.UInt8
	ActivePlayerInviteParam types.UInt16
	ActivePlayerJoinParam   types.UInt16
	ExtraParams             types.List[types.UInt32]
}

// WriteTo writes the AttractionStatus to the given writable
func (as AttractionStatus) WriteTo(writable types.Writable) {
	contentWritable := writable.CopyNew()

	as.MessageInterval.WriteTo(contentWritable)
	as.OperationFlag.WriteTo(contentWritable)
	as.ActivePlayerInviteParam.WriteTo(contentWritable)
	as.ActivePlayerJoinParam.WriteTo(contentWritable)
	as.ExtraParams.WriteTo(contentWritable)

	content := contentWritable.Bytes()

	as.WriteHeaderTo(writable, uint32(len(content)))

	writable.Write(content)
}

// ExtractFrom extracts the AttractionStatus from the given readable
func (as *AttractionStatus) ExtractFrom(readable types.Readable) error {
	if err := as.ExtractHeaderFrom(readable); err != nil {
		return fmt.Errorf("failed to extract AttractionStatus header. %s", err.Error())
	}

	if err := as.MessageInterval.ExtractFrom(readable); err != nil {
		return fmt.Errorf("failed to extract AttractionStatus.MessageInterval. %s", err.Error())
	}

	if err := as.OperationFlag.ExtractFrom(readable); err != nil {
		return fmt.Errorf("failed to extract AttractionStatus.OperationFlag. %s", err.Error())
	}

	if err := as.ActivePlayerInviteParam.ExtractFrom(readable); err != nil {
		return fmt.Errorf("failed to extract AttractionStatus.ActivePlayerInviteParam. %s", err.Error())
	}

	if err := as.ActivePlayerJoinParam.ExtractFrom(readable); err != nil {
		return fmt.Errorf("failed to extract AttractionStatus.ActivePlayerJoinParam. %s", err.Error())
	}

	if err := as.ExtraParams.ExtractFrom(readable); err != nil {
		return fmt.Errorf("failed to extract AttractionStatus.ExtraParams. %s", err.Error())
	}

	return nil
}

// Copy returns a new copied instance of AttractionStatus
func (as AttractionStatus) Copy() types.RVType {
	copied := NewAttractionStatus()

	copied.StructureVersion = as.StructureVersion
	copied.MessageInterval = as.MessageInterval.Copy().(types.UInt16)
	copied.OperationFlag = as.OperationFlag.Copy().(types.UInt8)
	copied.ActivePlayerInviteParam = as.ActivePlayerInviteParam.Copy().(types.UInt16)
	copied.ActivePlayerJoinParam = as.ActivePlayerJoinParam.Copy().(types.UInt16)
	copied.ExtraParams = as.ExtraParams.Copy().(types.List[types.UInt32])

	return copied
}

// Equals checks if the given AttractionStatus contains the same data as the current AttractionStatus
func (as AttractionStatus) Equals(o types.RVType) bool {
	if _, ok := o.(AttractionStatus); !ok {
		return false
	}

	other := o.(AttractionStatus)

	if as.StructureVersion != other.StructureVersion {
		return false
	}

	if !as.MessageInterval.Equals(other.MessageInterval) {
		return false
	}

	if !as.OperationFlag.Equals(other.OperationFlag) {
		return false
	}

	if !as.ActivePlayerInviteParam.Equals(other.ActivePlayerInviteParam) {
		return false
	}

	if !as.ActivePlayerJoinParam.Equals(other.ActivePlayerJoinParam) {
		return false
	}

	return as.ExtraParams.Equals(other.ExtraParams)
}

// CopyRef copies the current value of the AttractionStatus
// and returns a pointer to the new copy
func (as AttractionStatus) CopyRef() types.RVTypePtr {
	copied := as.Copy().(AttractionStatus)
	return &copied
}

// Deref takes a pointer to the AttractionStatus
// and dereferences it to the raw value.
// Only useful when working with an instance of RVTypePtr
func (as *AttractionStatus) Deref() types.RVType {
	return *as
}

// String returns the string representation of the AttractionStatus
func (as AttractionStatus) String() string {
	return as.FormatToString(0)
}

// FormatToString pretty-prints the AttractionStatus using the provided indentation level
func (as AttractionStatus) FormatToString(indentationLevel int) string {
	indentationValues := strings.Repeat("\t", indentationLevel+1)
	indentationEnd := strings.Repeat("\t", indentationLevel)

	var b strings.Builder

	b.WriteString("AttractionStatus{\n")
	fmt.Fprintf(&b, "%sMessageInterval: %s,\n", indentationValues, as.MessageInterval)
	fmt.Fprintf(&b, "%sOperationFlag: %s,\n", indentationValues, as.OperationFlag)
	fmt.Fprintf(&b, "%sActivePlayerInviteParam: %s,\n", indentationValues, as.ActivePlayerInviteParam)
	fmt.Fprintf(&b, "%sActivePlayerJoinParam: %s,\n", indentationValues, as.ActivePlayerJoinParam)
	fmt.Fprintf(&b, "%sExtraParams: %s,\n", indentationValues, as.ExtraParams)
	fmt.Fprintf(&b, "%s}", indentationEnd)

	return b.String()
}

// NewAttractionStatus returns a new AttractionStatus
func NewAttractionStatus() AttractionStatus {
	return AttractionStatus{
		MessageInterval:         types.NewUInt16(0),
		OperationFlag:           types.NewUInt8(0),
		ActivePlayerInviteParam: types.NewUInt16(0),
		ActivePlayerJoinParam:   types.NewUInt16(0),
		ExtraParams:             types.NewList[types.UInt32](),
	}
}
