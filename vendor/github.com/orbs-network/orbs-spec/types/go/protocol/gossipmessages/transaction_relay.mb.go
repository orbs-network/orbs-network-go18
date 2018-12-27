// AUTO GENERATED FILE (by membufc proto compiler v0.0.21)
package gossipmessages

import (
	"bytes"
	"fmt"
	"github.com/orbs-network/membuffers/go"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
)

/////////////////////////////////////////////////////////////////////////////
// message TempKillMeTransactionRelay

// reader

type TempKillMeTransactionRelay struct {

	// internal
	// implements membuffers.Message
	_message membuffers.InternalMessage
}

func (x *TempKillMeTransactionRelay) String() string {
	if x == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{}")
}

var _TempKillMeTransactionRelay_Scheme = []membuffers.FieldType{}
var _TempKillMeTransactionRelay_Unions = [][]membuffers.FieldType{}

func TempKillMeTransactionRelayReader(buf []byte) *TempKillMeTransactionRelay {
	x := &TempKillMeTransactionRelay{}
	x._message.Init(buf, membuffers.Offset(len(buf)), _TempKillMeTransactionRelay_Scheme, _TempKillMeTransactionRelay_Unions)
	return x
}

func (x *TempKillMeTransactionRelay) IsValid() bool {
	return x._message.IsValid()
}

func (x *TempKillMeTransactionRelay) Raw() []byte {
	return x._message.RawBuffer()
}

func (x *TempKillMeTransactionRelay) Equal(y *TempKillMeTransactionRelay) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	return bytes.Equal(x.Raw(), y.Raw())
}

// builder

type TempKillMeTransactionRelayBuilder struct {

	// internal
	// implements membuffers.Builder
	_builder               membuffers.InternalBuilder
	_overrideWithRawBuffer []byte
}

func (w *TempKillMeTransactionRelayBuilder) Write(buf []byte) (err error) {
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

func (w *TempKillMeTransactionRelayBuilder) HexDump(prefix string, offsetFromStart membuffers.Offset) (err error) {
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

func (w *TempKillMeTransactionRelayBuilder) GetSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	return w._builder.GetSize()
}

func (w *TempKillMeTransactionRelayBuilder) CalcRequiredSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	w.Write(nil)
	return w._builder.GetSize()
}

func (w *TempKillMeTransactionRelayBuilder) Build() *TempKillMeTransactionRelay {
	buf := make([]byte, w.CalcRequiredSize())
	if w.Write(buf) != nil {
		return nil
	}
	return TempKillMeTransactionRelayReader(buf)
}

func TempKillMeTransactionRelayBuilderFromRaw(raw []byte) *TempKillMeTransactionRelayBuilder {
	return &TempKillMeTransactionRelayBuilder{_overrideWithRawBuffer: raw}
}

/////////////////////////////////////////////////////////////////////////////
// message ForwardedTransactionsMessage (non serializable)

type ForwardedTransactionsMessage struct {
	Sender             *SenderSignature
	SignedTransactions []*protocol.SignedTransaction
}

func (x *ForwardedTransactionsMessage) String() string {
	if x == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{Sender:%s,SignedTransactions:%s,}", x.StringSender(), x.StringSignedTransactions())
}

func (x *ForwardedTransactionsMessage) StringSender() (res string) {
	res = x.Sender.String()
	return
}

func (x *ForwardedTransactionsMessage) StringSignedTransactions() (res string) {
	res = "["
	for _, v := range x.SignedTransactions {
		res += v.String() + ","
	}
	res += "]"
	return
}

/////////////////////////////////////////////////////////////////////////////
// enums

type TransactionsRelayMessageType uint16

const (
	TRANSACTION_RELAY_RESERVED               TransactionsRelayMessageType = 0
	TRANSACTION_RELAY_FORWARDED_TRANSACTIONS TransactionsRelayMessageType = 1
)

func (n TransactionsRelayMessageType) String() string {
	switch n {
	case TRANSACTION_RELAY_RESERVED:
		return "TRANSACTION_RELAY_RESERVED"
	case TRANSACTION_RELAY_FORWARDED_TRANSACTIONS:
		return "TRANSACTION_RELAY_FORWARDED_TRANSACTIONS"
	}
	return "UNKNOWN"
}
