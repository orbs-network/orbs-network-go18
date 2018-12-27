// AUTO GENERATED FILE (by membufc proto compiler v0.0.21)
package gossipmessages

import (
	"bytes"
	"fmt"
	"github.com/orbs-network/membuffers/go"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
)

/////////////////////////////////////////////////////////////////////////////
// message TempKillMeBenchmarkConsensus

// reader

type TempKillMeBenchmarkConsensus struct {

	// internal
	// implements membuffers.Message
	_message membuffers.InternalMessage
}

func (x *TempKillMeBenchmarkConsensus) String() string {
	if x == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{}")
}

var _TempKillMeBenchmarkConsensus_Scheme = []membuffers.FieldType{}
var _TempKillMeBenchmarkConsensus_Unions = [][]membuffers.FieldType{}

func TempKillMeBenchmarkConsensusReader(buf []byte) *TempKillMeBenchmarkConsensus {
	x := &TempKillMeBenchmarkConsensus{}
	x._message.Init(buf, membuffers.Offset(len(buf)), _TempKillMeBenchmarkConsensus_Scheme, _TempKillMeBenchmarkConsensus_Unions)
	return x
}

func (x *TempKillMeBenchmarkConsensus) IsValid() bool {
	return x._message.IsValid()
}

func (x *TempKillMeBenchmarkConsensus) Raw() []byte {
	return x._message.RawBuffer()
}

func (x *TempKillMeBenchmarkConsensus) Equal(y *TempKillMeBenchmarkConsensus) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	return bytes.Equal(x.Raw(), y.Raw())
}

// builder

type TempKillMeBenchmarkConsensusBuilder struct {

	// internal
	// implements membuffers.Builder
	_builder               membuffers.InternalBuilder
	_overrideWithRawBuffer []byte
}

func (w *TempKillMeBenchmarkConsensusBuilder) Write(buf []byte) (err error) {
	if w == nil {
		return
	}
	w._builder.NotifyBuildStart()
	defer w._builder.NotifyBuildEnd()
	defer func() {
		if r := recover(); r != nil {
			err = &membuffers.ErrBufferOverrun{}
		}
	}()
	if w._overrideWithRawBuffer != nil {
		return w._builder.WriteOverrideWithRawBuffer(buf, w._overrideWithRawBuffer)
	}
	w._builder.Reset()
	return nil
}

func (w *TempKillMeBenchmarkConsensusBuilder) HexDump(prefix string, offsetFromStart membuffers.Offset) (err error) {
	if w == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = &membuffers.ErrBufferOverrun{}
		}
	}()
	w._builder.Reset()
	return nil
}

func (w *TempKillMeBenchmarkConsensusBuilder) GetSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	return w._builder.GetSize()
}

func (w *TempKillMeBenchmarkConsensusBuilder) CalcRequiredSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	w.Write(nil)
	return w._builder.GetSize()
}

func (w *TempKillMeBenchmarkConsensusBuilder) Build() *TempKillMeBenchmarkConsensus {
	buf := make([]byte, w.CalcRequiredSize())
	if w.Write(buf) != nil {
		return nil
	}
	return TempKillMeBenchmarkConsensusReader(buf)
}

func TempKillMeBenchmarkConsensusBuilderFromRaw(raw []byte) *TempKillMeBenchmarkConsensusBuilder {
	return &TempKillMeBenchmarkConsensusBuilder{_overrideWithRawBuffer: raw}
}

/////////////////////////////////////////////////////////////////////////////
// message BenchmarkConsensusCommitMessage (non serializable)

type BenchmarkConsensusCommitMessage struct {
	BlockPair *protocol.BlockPairContainer
}

func (x *BenchmarkConsensusCommitMessage) String() string {
	if x == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{BlockPair:%s,}", x.StringBlockPair())
}

func (x *BenchmarkConsensusCommitMessage) StringBlockPair() (res string) {
	res = x.BlockPair.String()
	return
}

/////////////////////////////////////////////////////////////////////////////
// message BenchmarkConsensusCommittedMessage (non serializable)

type BenchmarkConsensusCommittedMessage struct {
	Status *BenchmarkConsensusStatus
	Sender *SenderSignature
}

func (x *BenchmarkConsensusCommittedMessage) String() string {
	if x == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{Status:%s,Sender:%s,}", x.StringStatus(), x.StringSender())
}

func (x *BenchmarkConsensusCommittedMessage) StringStatus() (res string) {
	res = x.Status.String()
	return
}

func (x *BenchmarkConsensusCommittedMessage) StringSender() (res string) {
	res = x.Sender.String()
	return
}

/////////////////////////////////////////////////////////////////////////////
// message BenchmarkConsensusStatus

// reader

type BenchmarkConsensusStatus struct {
	// LastCommittedBlockHeight primitives.BlockHeight

	// internal
	// implements membuffers.Message
	_message membuffers.InternalMessage
}

func (x *BenchmarkConsensusStatus) String() string {
	if x == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{LastCommittedBlockHeight:%s,}", x.StringLastCommittedBlockHeight())
}

var _BenchmarkConsensusStatus_Scheme = []membuffers.FieldType{membuffers.TypeUint64}
var _BenchmarkConsensusStatus_Unions = [][]membuffers.FieldType{}

func BenchmarkConsensusStatusReader(buf []byte) *BenchmarkConsensusStatus {
	x := &BenchmarkConsensusStatus{}
	x._message.Init(buf, membuffers.Offset(len(buf)), _BenchmarkConsensusStatus_Scheme, _BenchmarkConsensusStatus_Unions)
	return x
}

func (x *BenchmarkConsensusStatus) IsValid() bool {
	return x._message.IsValid()
}

func (x *BenchmarkConsensusStatus) Raw() []byte {
	return x._message.RawBuffer()
}

func (x *BenchmarkConsensusStatus) Equal(y *BenchmarkConsensusStatus) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	return bytes.Equal(x.Raw(), y.Raw())
}

func (x *BenchmarkConsensusStatus) LastCommittedBlockHeight() primitives.BlockHeight {
	return primitives.BlockHeight(x._message.GetUint64(0))
}

func (x *BenchmarkConsensusStatus) RawLastCommittedBlockHeight() []byte {
	return x._message.RawBufferForField(0, 0)
}

func (x *BenchmarkConsensusStatus) MutateLastCommittedBlockHeight(v primitives.BlockHeight) error {
	return x._message.SetUint64(0, uint64(v))
}

func (x *BenchmarkConsensusStatus) StringLastCommittedBlockHeight() string {
	return fmt.Sprintf("%s", x.LastCommittedBlockHeight())
}

// builder

type BenchmarkConsensusStatusBuilder struct {
	LastCommittedBlockHeight primitives.BlockHeight

	// internal
	// implements membuffers.Builder
	_builder               membuffers.InternalBuilder
	_overrideWithRawBuffer []byte
}

func (w *BenchmarkConsensusStatusBuilder) Write(buf []byte) (err error) {
	if w == nil {
		return
	}
	w._builder.NotifyBuildStart()
	defer w._builder.NotifyBuildEnd()
	defer func() {
		if r := recover(); r != nil {
			err = &membuffers.ErrBufferOverrun{}
		}
	}()
	if w._overrideWithRawBuffer != nil {
		return w._builder.WriteOverrideWithRawBuffer(buf, w._overrideWithRawBuffer)
	}
	w._builder.Reset()
	w._builder.WriteUint64(buf, uint64(w.LastCommittedBlockHeight))
	return nil
}

func (w *BenchmarkConsensusStatusBuilder) HexDump(prefix string, offsetFromStart membuffers.Offset) (err error) {
	if w == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = &membuffers.ErrBufferOverrun{}
		}
	}()
	w._builder.Reset()
	w._builder.HexDumpUint64(prefix, offsetFromStart, "BenchmarkConsensusStatus.LastCommittedBlockHeight", uint64(w.LastCommittedBlockHeight))
	return nil
}

func (w *BenchmarkConsensusStatusBuilder) GetSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	return w._builder.GetSize()
}

func (w *BenchmarkConsensusStatusBuilder) CalcRequiredSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	w.Write(nil)
	return w._builder.GetSize()
}

func (w *BenchmarkConsensusStatusBuilder) Build() *BenchmarkConsensusStatus {
	buf := make([]byte, w.CalcRequiredSize())
	if w.Write(buf) != nil {
		return nil
	}
	return BenchmarkConsensusStatusReader(buf)
}

func BenchmarkConsensusStatusBuilderFromRaw(raw []byte) *BenchmarkConsensusStatusBuilder {
	return &BenchmarkConsensusStatusBuilder{_overrideWithRawBuffer: raw}
}

/////////////////////////////////////////////////////////////////////////////
// enums
