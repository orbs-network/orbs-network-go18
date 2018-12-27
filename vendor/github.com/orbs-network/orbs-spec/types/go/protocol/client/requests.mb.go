// AUTO GENERATED FILE (by membufc proto compiler v0.0.21)
package client

import (
	"bytes"
	"fmt"
	"github.com/orbs-network/membuffers/go"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
)

/////////////////////////////////////////////////////////////////////////////
// message SendTransactionRequest

// reader

type SendTransactionRequest struct {
	// SignedTransaction protocol.SignedTransaction

	// internal
	// implements membuffers.Message
	_message membuffers.InternalMessage
}

func (x *SendTransactionRequest) String() string {
	if x == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{SignedTransaction:%s,}", x.StringSignedTransaction())
}

var _SendTransactionRequest_Scheme = []membuffers.FieldType{membuffers.TypeMessage}
var _SendTransactionRequest_Unions = [][]membuffers.FieldType{}

func SendTransactionRequestReader(buf []byte) *SendTransactionRequest {
	x := &SendTransactionRequest{}
	x._message.Init(buf, membuffers.Offset(len(buf)), _SendTransactionRequest_Scheme, _SendTransactionRequest_Unions)
	return x
}

func (x *SendTransactionRequest) IsValid() bool {
	return x._message.IsValid()
}

func (x *SendTransactionRequest) Raw() []byte {
	return x._message.RawBuffer()
}

func (x *SendTransactionRequest) Equal(y *SendTransactionRequest) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	return bytes.Equal(x.Raw(), y.Raw())
}

func (x *SendTransactionRequest) SignedTransaction() *protocol.SignedTransaction {
	b, s := x._message.GetMessage(0)
	return protocol.SignedTransactionReader(b[:s])
}

func (x *SendTransactionRequest) RawSignedTransaction() []byte {
	return x._message.RawBufferForField(0, 0)
}

func (x *SendTransactionRequest) RawSignedTransactionWithHeader() []byte {
	return x._message.RawBufferWithHeaderForField(0, 0)
}

func (x *SendTransactionRequest) StringSignedTransaction() string {
	return x.SignedTransaction().String()
}

// builder

type SendTransactionRequestBuilder struct {
	SignedTransaction *protocol.SignedTransactionBuilder

	// internal
	// implements membuffers.Builder
	_builder               membuffers.InternalBuilder
	_overrideWithRawBuffer []byte
}

func (w *SendTransactionRequestBuilder) Write(buf []byte) (err error) {
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
	err = w._builder.WriteMessage(buf, w.SignedTransaction)
	if err != nil {
		return
	}
	return nil
}

func (w *SendTransactionRequestBuilder) HexDump(prefix string, offsetFromStart membuffers.Offset) (err error) {
	if w == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = &membuffers.ErrBufferOverrun{}
		}
	}()
	w._builder.Reset()
	err = w._builder.HexDumpMessage(prefix, offsetFromStart, "SendTransactionRequest.SignedTransaction", w.SignedTransaction)
	if err != nil {
		return
	}
	return nil
}

func (w *SendTransactionRequestBuilder) GetSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	return w._builder.GetSize()
}

func (w *SendTransactionRequestBuilder) CalcRequiredSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	w.Write(nil)
	return w._builder.GetSize()
}

func (w *SendTransactionRequestBuilder) Build() *SendTransactionRequest {
	buf := make([]byte, w.CalcRequiredSize())
	if w.Write(buf) != nil {
		return nil
	}
	return SendTransactionRequestReader(buf)
}

func SendTransactionRequestBuilderFromRaw(raw []byte) *SendTransactionRequestBuilder {
	return &SendTransactionRequestBuilder{_overrideWithRawBuffer: raw}
}

/////////////////////////////////////////////////////////////////////////////
// message SendTransactionResponse

// reader

type SendTransactionResponse struct {
	// RequestStatus protocol.RequestStatus
	// TransactionReceipt protocol.TransactionReceipt
	// TransactionStatus protocol.TransactionStatus
	// BlockHeight primitives.BlockHeight
	// BlockTimestamp primitives.TimestampNano

	// internal
	// implements membuffers.Message
	_message membuffers.InternalMessage
}

func (x *SendTransactionResponse) String() string {
	if x == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{RequestStatus:%s,TransactionReceipt:%s,TransactionStatus:%s,BlockHeight:%s,BlockTimestamp:%s,}", x.StringRequestStatus(), x.StringTransactionReceipt(), x.StringTransactionStatus(), x.StringBlockHeight(), x.StringBlockTimestamp())
}

var _SendTransactionResponse_Scheme = []membuffers.FieldType{membuffers.TypeUint16, membuffers.TypeMessage, membuffers.TypeUint16, membuffers.TypeUint64, membuffers.TypeUint64}
var _SendTransactionResponse_Unions = [][]membuffers.FieldType{}

func SendTransactionResponseReader(buf []byte) *SendTransactionResponse {
	x := &SendTransactionResponse{}
	x._message.Init(buf, membuffers.Offset(len(buf)), _SendTransactionResponse_Scheme, _SendTransactionResponse_Unions)
	return x
}

func (x *SendTransactionResponse) IsValid() bool {
	return x._message.IsValid()
}

func (x *SendTransactionResponse) Raw() []byte {
	return x._message.RawBuffer()
}

func (x *SendTransactionResponse) Equal(y *SendTransactionResponse) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	return bytes.Equal(x.Raw(), y.Raw())
}

func (x *SendTransactionResponse) RequestStatus() protocol.RequestStatus {
	return protocol.RequestStatus(x._message.GetUint16(0))
}

func (x *SendTransactionResponse) RawRequestStatus() []byte {
	return x._message.RawBufferForField(0, 0)
}

func (x *SendTransactionResponse) MutateRequestStatus(v protocol.RequestStatus) error {
	return x._message.SetUint16(0, uint16(v))
}

func (x *SendTransactionResponse) StringRequestStatus() string {
	return x.RequestStatus().String()
}

func (x *SendTransactionResponse) TransactionReceipt() *protocol.TransactionReceipt {
	b, s := x._message.GetMessage(1)
	return protocol.TransactionReceiptReader(b[:s])
}

func (x *SendTransactionResponse) RawTransactionReceipt() []byte {
	return x._message.RawBufferForField(1, 0)
}

func (x *SendTransactionResponse) RawTransactionReceiptWithHeader() []byte {
	return x._message.RawBufferWithHeaderForField(1, 0)
}

func (x *SendTransactionResponse) StringTransactionReceipt() string {
	return x.TransactionReceipt().String()
}

func (x *SendTransactionResponse) TransactionStatus() protocol.TransactionStatus {
	return protocol.TransactionStatus(x._message.GetUint16(2))
}

func (x *SendTransactionResponse) RawTransactionStatus() []byte {
	return x._message.RawBufferForField(2, 0)
}

func (x *SendTransactionResponse) MutateTransactionStatus(v protocol.TransactionStatus) error {
	return x._message.SetUint16(2, uint16(v))
}

func (x *SendTransactionResponse) StringTransactionStatus() string {
	return x.TransactionStatus().String()
}

func (x *SendTransactionResponse) BlockHeight() primitives.BlockHeight {
	return primitives.BlockHeight(x._message.GetUint64(3))
}

func (x *SendTransactionResponse) RawBlockHeight() []byte {
	return x._message.RawBufferForField(3, 0)
}

func (x *SendTransactionResponse) MutateBlockHeight(v primitives.BlockHeight) error {
	return x._message.SetUint64(3, uint64(v))
}

func (x *SendTransactionResponse) StringBlockHeight() string {
	return fmt.Sprintf("%s", x.BlockHeight())
}

func (x *SendTransactionResponse) BlockTimestamp() primitives.TimestampNano {
	return primitives.TimestampNano(x._message.GetUint64(4))
}

func (x *SendTransactionResponse) RawBlockTimestamp() []byte {
	return x._message.RawBufferForField(4, 0)
}

func (x *SendTransactionResponse) MutateBlockTimestamp(v primitives.TimestampNano) error {
	return x._message.SetUint64(4, uint64(v))
}

func (x *SendTransactionResponse) StringBlockTimestamp() string {
	return fmt.Sprintf("%s", x.BlockTimestamp())
}

// builder

type SendTransactionResponseBuilder struct {
	RequestStatus      protocol.RequestStatus
	TransactionReceipt *protocol.TransactionReceiptBuilder
	TransactionStatus  protocol.TransactionStatus
	BlockHeight        primitives.BlockHeight
	BlockTimestamp     primitives.TimestampNano

	// internal
	// implements membuffers.Builder
	_builder               membuffers.InternalBuilder
	_overrideWithRawBuffer []byte
}

func (w *SendTransactionResponseBuilder) Write(buf []byte) (err error) {
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
	w._builder.WriteUint16(buf, uint16(w.RequestStatus))
	err = w._builder.WriteMessage(buf, w.TransactionReceipt)
	if err != nil {
		return
	}
	w._builder.WriteUint16(buf, uint16(w.TransactionStatus))
	w._builder.WriteUint64(buf, uint64(w.BlockHeight))
	w._builder.WriteUint64(buf, uint64(w.BlockTimestamp))
	return nil
}

func (w *SendTransactionResponseBuilder) HexDump(prefix string, offsetFromStart membuffers.Offset) (err error) {
	if w == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = &membuffers.ErrBufferOverrun{}
		}
	}()
	w._builder.Reset()
	w._builder.HexDumpUint16(prefix, offsetFromStart, "SendTransactionResponse.RequestStatus", uint16(w.RequestStatus))
	err = w._builder.HexDumpMessage(prefix, offsetFromStart, "SendTransactionResponse.TransactionReceipt", w.TransactionReceipt)
	if err != nil {
		return
	}
	w._builder.HexDumpUint16(prefix, offsetFromStart, "SendTransactionResponse.TransactionStatus", uint16(w.TransactionStatus))
	w._builder.HexDumpUint64(prefix, offsetFromStart, "SendTransactionResponse.BlockHeight", uint64(w.BlockHeight))
	w._builder.HexDumpUint64(prefix, offsetFromStart, "SendTransactionResponse.BlockTimestamp", uint64(w.BlockTimestamp))
	return nil
}

func (w *SendTransactionResponseBuilder) GetSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	return w._builder.GetSize()
}

func (w *SendTransactionResponseBuilder) CalcRequiredSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	w.Write(nil)
	return w._builder.GetSize()
}

func (w *SendTransactionResponseBuilder) Build() *SendTransactionResponse {
	buf := make([]byte, w.CalcRequiredSize())
	if w.Write(buf) != nil {
		return nil
	}
	return SendTransactionResponseReader(buf)
}

func SendTransactionResponseBuilderFromRaw(raw []byte) *SendTransactionResponseBuilder {
	return &SendTransactionResponseBuilder{_overrideWithRawBuffer: raw}
}

/////////////////////////////////////////////////////////////////////////////
// message CallMethodRequest

// reader

type CallMethodRequest struct {
	// Transaction protocol.Transaction

	// internal
	// implements membuffers.Message
	_message membuffers.InternalMessage
}

func (x *CallMethodRequest) String() string {
	if x == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{Transaction:%s,}", x.StringTransaction())
}

var _CallMethodRequest_Scheme = []membuffers.FieldType{membuffers.TypeMessage}
var _CallMethodRequest_Unions = [][]membuffers.FieldType{}

func CallMethodRequestReader(buf []byte) *CallMethodRequest {
	x := &CallMethodRequest{}
	x._message.Init(buf, membuffers.Offset(len(buf)), _CallMethodRequest_Scheme, _CallMethodRequest_Unions)
	return x
}

func (x *CallMethodRequest) IsValid() bool {
	return x._message.IsValid()
}

func (x *CallMethodRequest) Raw() []byte {
	return x._message.RawBuffer()
}

func (x *CallMethodRequest) Equal(y *CallMethodRequest) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	return bytes.Equal(x.Raw(), y.Raw())
}

func (x *CallMethodRequest) Transaction() *protocol.Transaction {
	b, s := x._message.GetMessage(0)
	return protocol.TransactionReader(b[:s])
}

func (x *CallMethodRequest) RawTransaction() []byte {
	return x._message.RawBufferForField(0, 0)
}

func (x *CallMethodRequest) RawTransactionWithHeader() []byte {
	return x._message.RawBufferWithHeaderForField(0, 0)
}

func (x *CallMethodRequest) StringTransaction() string {
	return x.Transaction().String()
}

// builder

type CallMethodRequestBuilder struct {
	Transaction *protocol.TransactionBuilder

	// internal
	// implements membuffers.Builder
	_builder               membuffers.InternalBuilder
	_overrideWithRawBuffer []byte
}

func (w *CallMethodRequestBuilder) Write(buf []byte) (err error) {
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
	err = w._builder.WriteMessage(buf, w.Transaction)
	if err != nil {
		return
	}
	return nil
}

func (w *CallMethodRequestBuilder) HexDump(prefix string, offsetFromStart membuffers.Offset) (err error) {
	if w == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = &membuffers.ErrBufferOverrun{}
		}
	}()
	w._builder.Reset()
	err = w._builder.HexDumpMessage(prefix, offsetFromStart, "CallMethodRequest.Transaction", w.Transaction)
	if err != nil {
		return
	}
	return nil
}

func (w *CallMethodRequestBuilder) GetSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	return w._builder.GetSize()
}

func (w *CallMethodRequestBuilder) CalcRequiredSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	w.Write(nil)
	return w._builder.GetSize()
}

func (w *CallMethodRequestBuilder) Build() *CallMethodRequest {
	buf := make([]byte, w.CalcRequiredSize())
	if w.Write(buf) != nil {
		return nil
	}
	return CallMethodRequestReader(buf)
}

func CallMethodRequestBuilderFromRaw(raw []byte) *CallMethodRequestBuilder {
	return &CallMethodRequestBuilder{_overrideWithRawBuffer: raw}
}

/////////////////////////////////////////////////////////////////////////////
// message CallMethodResponse

// reader

type CallMethodResponse struct {
	// RequestStatus protocol.RequestStatus
	// OutputArgumentArray primitives.PackedArgumentArray
	// OutputEventsArray primitives.PackedEventsArray
	// CallMethodResult protocol.ExecutionResult
	// BlockHeight primitives.BlockHeight
	// BlockTimestamp primitives.TimestampNano

	// internal
	// implements membuffers.Message
	_message membuffers.InternalMessage
}

func (x *CallMethodResponse) String() string {
	if x == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{RequestStatus:%s,OutputArgumentArray:%s,OutputEventsArray:%s,CallMethodResult:%s,BlockHeight:%s,BlockTimestamp:%s,}", x.StringRequestStatus(), x.StringOutputArgumentArray(), x.StringOutputEventsArray(), x.StringCallMethodResult(), x.StringBlockHeight(), x.StringBlockTimestamp())
}

var _CallMethodResponse_Scheme = []membuffers.FieldType{membuffers.TypeUint16, membuffers.TypeBytes, membuffers.TypeBytes, membuffers.TypeUint16, membuffers.TypeUint64, membuffers.TypeUint64}
var _CallMethodResponse_Unions = [][]membuffers.FieldType{}

func CallMethodResponseReader(buf []byte) *CallMethodResponse {
	x := &CallMethodResponse{}
	x._message.Init(buf, membuffers.Offset(len(buf)), _CallMethodResponse_Scheme, _CallMethodResponse_Unions)
	return x
}

func (x *CallMethodResponse) IsValid() bool {
	return x._message.IsValid()
}

func (x *CallMethodResponse) Raw() []byte {
	return x._message.RawBuffer()
}

func (x *CallMethodResponse) Equal(y *CallMethodResponse) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	return bytes.Equal(x.Raw(), y.Raw())
}

func (x *CallMethodResponse) RequestStatus() protocol.RequestStatus {
	return protocol.RequestStatus(x._message.GetUint16(0))
}

func (x *CallMethodResponse) RawRequestStatus() []byte {
	return x._message.RawBufferForField(0, 0)
}

func (x *CallMethodResponse) MutateRequestStatus(v protocol.RequestStatus) error {
	return x._message.SetUint16(0, uint16(v))
}

func (x *CallMethodResponse) StringRequestStatus() string {
	return x.RequestStatus().String()
}

func (x *CallMethodResponse) OutputArgumentArray() primitives.PackedArgumentArray {
	return primitives.PackedArgumentArray(x._message.GetBytes(1))
}

func (x *CallMethodResponse) RawOutputArgumentArray() []byte {
	return x._message.RawBufferForField(1, 0)
}

func (x *CallMethodResponse) RawOutputArgumentArrayWithHeader() []byte {
	return x._message.RawBufferWithHeaderForField(1, 0)
}

func (x *CallMethodResponse) MutateOutputArgumentArray(v primitives.PackedArgumentArray) error {
	return x._message.SetBytes(1, []byte(v))
}

func (x *CallMethodResponse) StringOutputArgumentArray() string {
	return fmt.Sprintf("%s", x.OutputArgumentArray())
}

func (x *CallMethodResponse) OutputEventsArray() primitives.PackedEventsArray {
	return primitives.PackedEventsArray(x._message.GetBytes(2))
}

func (x *CallMethodResponse) RawOutputEventsArray() []byte {
	return x._message.RawBufferForField(2, 0)
}

func (x *CallMethodResponse) RawOutputEventsArrayWithHeader() []byte {
	return x._message.RawBufferWithHeaderForField(2, 0)
}

func (x *CallMethodResponse) MutateOutputEventsArray(v primitives.PackedEventsArray) error {
	return x._message.SetBytes(2, []byte(v))
}

func (x *CallMethodResponse) StringOutputEventsArray() string {
	return fmt.Sprintf("%s", x.OutputEventsArray())
}

func (x *CallMethodResponse) CallMethodResult() protocol.ExecutionResult {
	return protocol.ExecutionResult(x._message.GetUint16(3))
}

func (x *CallMethodResponse) RawCallMethodResult() []byte {
	return x._message.RawBufferForField(3, 0)
}

func (x *CallMethodResponse) MutateCallMethodResult(v protocol.ExecutionResult) error {
	return x._message.SetUint16(3, uint16(v))
}

func (x *CallMethodResponse) StringCallMethodResult() string {
	return x.CallMethodResult().String()
}

func (x *CallMethodResponse) BlockHeight() primitives.BlockHeight {
	return primitives.BlockHeight(x._message.GetUint64(4))
}

func (x *CallMethodResponse) RawBlockHeight() []byte {
	return x._message.RawBufferForField(4, 0)
}

func (x *CallMethodResponse) MutateBlockHeight(v primitives.BlockHeight) error {
	return x._message.SetUint64(4, uint64(v))
}

func (x *CallMethodResponse) StringBlockHeight() string {
	return fmt.Sprintf("%s", x.BlockHeight())
}

func (x *CallMethodResponse) BlockTimestamp() primitives.TimestampNano {
	return primitives.TimestampNano(x._message.GetUint64(5))
}

func (x *CallMethodResponse) RawBlockTimestamp() []byte {
	return x._message.RawBufferForField(5, 0)
}

func (x *CallMethodResponse) MutateBlockTimestamp(v primitives.TimestampNano) error {
	return x._message.SetUint64(5, uint64(v))
}

func (x *CallMethodResponse) StringBlockTimestamp() string {
	return fmt.Sprintf("%s", x.BlockTimestamp())
}

// builder

type CallMethodResponseBuilder struct {
	RequestStatus       protocol.RequestStatus
	OutputArgumentArray primitives.PackedArgumentArray
	OutputEventsArray   primitives.PackedEventsArray
	CallMethodResult    protocol.ExecutionResult
	BlockHeight         primitives.BlockHeight
	BlockTimestamp      primitives.TimestampNano

	// internal
	// implements membuffers.Builder
	_builder               membuffers.InternalBuilder
	_overrideWithRawBuffer []byte
}

func (w *CallMethodResponseBuilder) Write(buf []byte) (err error) {
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
	w._builder.WriteUint16(buf, uint16(w.RequestStatus))
	w._builder.WriteBytes(buf, []byte(w.OutputArgumentArray))
	w._builder.WriteBytes(buf, []byte(w.OutputEventsArray))
	w._builder.WriteUint16(buf, uint16(w.CallMethodResult))
	w._builder.WriteUint64(buf, uint64(w.BlockHeight))
	w._builder.WriteUint64(buf, uint64(w.BlockTimestamp))
	return nil
}

func (w *CallMethodResponseBuilder) HexDump(prefix string, offsetFromStart membuffers.Offset) (err error) {
	if w == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = &membuffers.ErrBufferOverrun{}
		}
	}()
	w._builder.Reset()
	w._builder.HexDumpUint16(prefix, offsetFromStart, "CallMethodResponse.RequestStatus", uint16(w.RequestStatus))
	w._builder.HexDumpBytes(prefix, offsetFromStart, "CallMethodResponse.OutputArgumentArray", []byte(w.OutputArgumentArray))
	w._builder.HexDumpBytes(prefix, offsetFromStart, "CallMethodResponse.OutputEventsArray", []byte(w.OutputEventsArray))
	w._builder.HexDumpUint16(prefix, offsetFromStart, "CallMethodResponse.CallMethodResult", uint16(w.CallMethodResult))
	w._builder.HexDumpUint64(prefix, offsetFromStart, "CallMethodResponse.BlockHeight", uint64(w.BlockHeight))
	w._builder.HexDumpUint64(prefix, offsetFromStart, "CallMethodResponse.BlockTimestamp", uint64(w.BlockTimestamp))
	return nil
}

func (w *CallMethodResponseBuilder) GetSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	return w._builder.GetSize()
}

func (w *CallMethodResponseBuilder) CalcRequiredSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	w.Write(nil)
	return w._builder.GetSize()
}

func (w *CallMethodResponseBuilder) Build() *CallMethodResponse {
	buf := make([]byte, w.CalcRequiredSize())
	if w.Write(buf) != nil {
		return nil
	}
	return CallMethodResponseReader(buf)
}

func CallMethodResponseBuilderFromRaw(raw []byte) *CallMethodResponseBuilder {
	return &CallMethodResponseBuilder{_overrideWithRawBuffer: raw}
}

/////////////////////////////////////////////////////////////////////////////
// message GetTransactionStatusRequest

// reader

type GetTransactionStatusRequest struct {
	// ProtocolVersion primitives.ProtocolVersion
	// VirtualChainId primitives.VirtualChainId
	// TransactionTimestamp primitives.TimestampNano
	// Txhash primitives.Sha256

	// internal
	// implements membuffers.Message
	_message membuffers.InternalMessage
}

func (x *GetTransactionStatusRequest) String() string {
	if x == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{ProtocolVersion:%s,VirtualChainId:%s,TransactionTimestamp:%s,Txhash:%s,}", x.StringProtocolVersion(), x.StringVirtualChainId(), x.StringTransactionTimestamp(), x.StringTxhash())
}

var _GetTransactionStatusRequest_Scheme = []membuffers.FieldType{membuffers.TypeUint32, membuffers.TypeUint32, membuffers.TypeUint64, membuffers.TypeBytes}
var _GetTransactionStatusRequest_Unions = [][]membuffers.FieldType{}

func GetTransactionStatusRequestReader(buf []byte) *GetTransactionStatusRequest {
	x := &GetTransactionStatusRequest{}
	x._message.Init(buf, membuffers.Offset(len(buf)), _GetTransactionStatusRequest_Scheme, _GetTransactionStatusRequest_Unions)
	return x
}

func (x *GetTransactionStatusRequest) IsValid() bool {
	return x._message.IsValid()
}

func (x *GetTransactionStatusRequest) Raw() []byte {
	return x._message.RawBuffer()
}

func (x *GetTransactionStatusRequest) Equal(y *GetTransactionStatusRequest) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	return bytes.Equal(x.Raw(), y.Raw())
}

func (x *GetTransactionStatusRequest) ProtocolVersion() primitives.ProtocolVersion {
	return primitives.ProtocolVersion(x._message.GetUint32(0))
}

func (x *GetTransactionStatusRequest) RawProtocolVersion() []byte {
	return x._message.RawBufferForField(0, 0)
}

func (x *GetTransactionStatusRequest) MutateProtocolVersion(v primitives.ProtocolVersion) error {
	return x._message.SetUint32(0, uint32(v))
}

func (x *GetTransactionStatusRequest) StringProtocolVersion() string {
	return fmt.Sprintf("%s", x.ProtocolVersion())
}

func (x *GetTransactionStatusRequest) VirtualChainId() primitives.VirtualChainId {
	return primitives.VirtualChainId(x._message.GetUint32(1))
}

func (x *GetTransactionStatusRequest) RawVirtualChainId() []byte {
	return x._message.RawBufferForField(1, 0)
}

func (x *GetTransactionStatusRequest) MutateVirtualChainId(v primitives.VirtualChainId) error {
	return x._message.SetUint32(1, uint32(v))
}

func (x *GetTransactionStatusRequest) StringVirtualChainId() string {
	return fmt.Sprintf("%s", x.VirtualChainId())
}

func (x *GetTransactionStatusRequest) TransactionTimestamp() primitives.TimestampNano {
	return primitives.TimestampNano(x._message.GetUint64(2))
}

func (x *GetTransactionStatusRequest) RawTransactionTimestamp() []byte {
	return x._message.RawBufferForField(2, 0)
}

func (x *GetTransactionStatusRequest) MutateTransactionTimestamp(v primitives.TimestampNano) error {
	return x._message.SetUint64(2, uint64(v))
}

func (x *GetTransactionStatusRequest) StringTransactionTimestamp() string {
	return fmt.Sprintf("%s", x.TransactionTimestamp())
}

func (x *GetTransactionStatusRequest) Txhash() primitives.Sha256 {
	return primitives.Sha256(x._message.GetBytes(3))
}

func (x *GetTransactionStatusRequest) RawTxhash() []byte {
	return x._message.RawBufferForField(3, 0)
}

func (x *GetTransactionStatusRequest) RawTxhashWithHeader() []byte {
	return x._message.RawBufferWithHeaderForField(3, 0)
}

func (x *GetTransactionStatusRequest) MutateTxhash(v primitives.Sha256) error {
	return x._message.SetBytes(3, []byte(v))
}

func (x *GetTransactionStatusRequest) StringTxhash() string {
	return fmt.Sprintf("%s", x.Txhash())
}

// builder

type GetTransactionStatusRequestBuilder struct {
	ProtocolVersion      primitives.ProtocolVersion
	VirtualChainId       primitives.VirtualChainId
	TransactionTimestamp primitives.TimestampNano
	Txhash               primitives.Sha256

	// internal
	// implements membuffers.Builder
	_builder               membuffers.InternalBuilder
	_overrideWithRawBuffer []byte
}

func (w *GetTransactionStatusRequestBuilder) Write(buf []byte) (err error) {
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
	w._builder.WriteUint32(buf, uint32(w.ProtocolVersion))
	w._builder.WriteUint32(buf, uint32(w.VirtualChainId))
	w._builder.WriteUint64(buf, uint64(w.TransactionTimestamp))
	w._builder.WriteBytes(buf, []byte(w.Txhash))
	return nil
}

func (w *GetTransactionStatusRequestBuilder) HexDump(prefix string, offsetFromStart membuffers.Offset) (err error) {
	if w == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = &membuffers.ErrBufferOverrun{}
		}
	}()
	w._builder.Reset()
	w._builder.HexDumpUint32(prefix, offsetFromStart, "GetTransactionStatusRequest.ProtocolVersion", uint32(w.ProtocolVersion))
	w._builder.HexDumpUint32(prefix, offsetFromStart, "GetTransactionStatusRequest.VirtualChainId", uint32(w.VirtualChainId))
	w._builder.HexDumpUint64(prefix, offsetFromStart, "GetTransactionStatusRequest.TransactionTimestamp", uint64(w.TransactionTimestamp))
	w._builder.HexDumpBytes(prefix, offsetFromStart, "GetTransactionStatusRequest.Txhash", []byte(w.Txhash))
	return nil
}

func (w *GetTransactionStatusRequestBuilder) GetSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	return w._builder.GetSize()
}

func (w *GetTransactionStatusRequestBuilder) CalcRequiredSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	w.Write(nil)
	return w._builder.GetSize()
}

func (w *GetTransactionStatusRequestBuilder) Build() *GetTransactionStatusRequest {
	buf := make([]byte, w.CalcRequiredSize())
	if w.Write(buf) != nil {
		return nil
	}
	return GetTransactionStatusRequestReader(buf)
}

func GetTransactionStatusRequestBuilderFromRaw(raw []byte) *GetTransactionStatusRequestBuilder {
	return &GetTransactionStatusRequestBuilder{_overrideWithRawBuffer: raw}
}

/////////////////////////////////////////////////////////////////////////////
// message GetTransactionStatusResponse

// reader

type GetTransactionStatusResponse struct {
	// RequestStatus protocol.RequestStatus
	// TransactionReceipt protocol.TransactionReceipt
	// TransactionStatus protocol.TransactionStatus
	// BlockHeight primitives.BlockHeight
	// BlockTimestamp primitives.TimestampNano

	// internal
	// implements membuffers.Message
	_message membuffers.InternalMessage
}

func (x *GetTransactionStatusResponse) String() string {
	if x == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{RequestStatus:%s,TransactionReceipt:%s,TransactionStatus:%s,BlockHeight:%s,BlockTimestamp:%s,}", x.StringRequestStatus(), x.StringTransactionReceipt(), x.StringTransactionStatus(), x.StringBlockHeight(), x.StringBlockTimestamp())
}

var _GetTransactionStatusResponse_Scheme = []membuffers.FieldType{membuffers.TypeUint16, membuffers.TypeMessage, membuffers.TypeUint16, membuffers.TypeUint64, membuffers.TypeUint64}
var _GetTransactionStatusResponse_Unions = [][]membuffers.FieldType{}

func GetTransactionStatusResponseReader(buf []byte) *GetTransactionStatusResponse {
	x := &GetTransactionStatusResponse{}
	x._message.Init(buf, membuffers.Offset(len(buf)), _GetTransactionStatusResponse_Scheme, _GetTransactionStatusResponse_Unions)
	return x
}

func (x *GetTransactionStatusResponse) IsValid() bool {
	return x._message.IsValid()
}

func (x *GetTransactionStatusResponse) Raw() []byte {
	return x._message.RawBuffer()
}

func (x *GetTransactionStatusResponse) Equal(y *GetTransactionStatusResponse) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	return bytes.Equal(x.Raw(), y.Raw())
}

func (x *GetTransactionStatusResponse) RequestStatus() protocol.RequestStatus {
	return protocol.RequestStatus(x._message.GetUint16(0))
}

func (x *GetTransactionStatusResponse) RawRequestStatus() []byte {
	return x._message.RawBufferForField(0, 0)
}

func (x *GetTransactionStatusResponse) MutateRequestStatus(v protocol.RequestStatus) error {
	return x._message.SetUint16(0, uint16(v))
}

func (x *GetTransactionStatusResponse) StringRequestStatus() string {
	return x.RequestStatus().String()
}

func (x *GetTransactionStatusResponse) TransactionReceipt() *protocol.TransactionReceipt {
	b, s := x._message.GetMessage(1)
	return protocol.TransactionReceiptReader(b[:s])
}

func (x *GetTransactionStatusResponse) RawTransactionReceipt() []byte {
	return x._message.RawBufferForField(1, 0)
}

func (x *GetTransactionStatusResponse) RawTransactionReceiptWithHeader() []byte {
	return x._message.RawBufferWithHeaderForField(1, 0)
}

func (x *GetTransactionStatusResponse) StringTransactionReceipt() string {
	return x.TransactionReceipt().String()
}

func (x *GetTransactionStatusResponse) TransactionStatus() protocol.TransactionStatus {
	return protocol.TransactionStatus(x._message.GetUint16(2))
}

func (x *GetTransactionStatusResponse) RawTransactionStatus() []byte {
	return x._message.RawBufferForField(2, 0)
}

func (x *GetTransactionStatusResponse) MutateTransactionStatus(v protocol.TransactionStatus) error {
	return x._message.SetUint16(2, uint16(v))
}

func (x *GetTransactionStatusResponse) StringTransactionStatus() string {
	return x.TransactionStatus().String()
}

func (x *GetTransactionStatusResponse) BlockHeight() primitives.BlockHeight {
	return primitives.BlockHeight(x._message.GetUint64(3))
}

func (x *GetTransactionStatusResponse) RawBlockHeight() []byte {
	return x._message.RawBufferForField(3, 0)
}

func (x *GetTransactionStatusResponse) MutateBlockHeight(v primitives.BlockHeight) error {
	return x._message.SetUint64(3, uint64(v))
}

func (x *GetTransactionStatusResponse) StringBlockHeight() string {
	return fmt.Sprintf("%s", x.BlockHeight())
}

func (x *GetTransactionStatusResponse) BlockTimestamp() primitives.TimestampNano {
	return primitives.TimestampNano(x._message.GetUint64(4))
}

func (x *GetTransactionStatusResponse) RawBlockTimestamp() []byte {
	return x._message.RawBufferForField(4, 0)
}

func (x *GetTransactionStatusResponse) MutateBlockTimestamp(v primitives.TimestampNano) error {
	return x._message.SetUint64(4, uint64(v))
}

func (x *GetTransactionStatusResponse) StringBlockTimestamp() string {
	return fmt.Sprintf("%s", x.BlockTimestamp())
}

// builder

type GetTransactionStatusResponseBuilder struct {
	RequestStatus      protocol.RequestStatus
	TransactionReceipt *protocol.TransactionReceiptBuilder
	TransactionStatus  protocol.TransactionStatus
	BlockHeight        primitives.BlockHeight
	BlockTimestamp     primitives.TimestampNano

	// internal
	// implements membuffers.Builder
	_builder               membuffers.InternalBuilder
	_overrideWithRawBuffer []byte
}

func (w *GetTransactionStatusResponseBuilder) Write(buf []byte) (err error) {
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
	w._builder.WriteUint16(buf, uint16(w.RequestStatus))
	err = w._builder.WriteMessage(buf, w.TransactionReceipt)
	if err != nil {
		return
	}
	w._builder.WriteUint16(buf, uint16(w.TransactionStatus))
	w._builder.WriteUint64(buf, uint64(w.BlockHeight))
	w._builder.WriteUint64(buf, uint64(w.BlockTimestamp))
	return nil
}

func (w *GetTransactionStatusResponseBuilder) HexDump(prefix string, offsetFromStart membuffers.Offset) (err error) {
	if w == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = &membuffers.ErrBufferOverrun{}
		}
	}()
	w._builder.Reset()
	w._builder.HexDumpUint16(prefix, offsetFromStart, "GetTransactionStatusResponse.RequestStatus", uint16(w.RequestStatus))
	err = w._builder.HexDumpMessage(prefix, offsetFromStart, "GetTransactionStatusResponse.TransactionReceipt", w.TransactionReceipt)
	if err != nil {
		return
	}
	w._builder.HexDumpUint16(prefix, offsetFromStart, "GetTransactionStatusResponse.TransactionStatus", uint16(w.TransactionStatus))
	w._builder.HexDumpUint64(prefix, offsetFromStart, "GetTransactionStatusResponse.BlockHeight", uint64(w.BlockHeight))
	w._builder.HexDumpUint64(prefix, offsetFromStart, "GetTransactionStatusResponse.BlockTimestamp", uint64(w.BlockTimestamp))
	return nil
}

func (w *GetTransactionStatusResponseBuilder) GetSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	return w._builder.GetSize()
}

func (w *GetTransactionStatusResponseBuilder) CalcRequiredSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	w.Write(nil)
	return w._builder.GetSize()
}

func (w *GetTransactionStatusResponseBuilder) Build() *GetTransactionStatusResponse {
	buf := make([]byte, w.CalcRequiredSize())
	if w.Write(buf) != nil {
		return nil
	}
	return GetTransactionStatusResponseReader(buf)
}

func GetTransactionStatusResponseBuilderFromRaw(raw []byte) *GetTransactionStatusResponseBuilder {
	return &GetTransactionStatusResponseBuilder{_overrideWithRawBuffer: raw}
}

/////////////////////////////////////////////////////////////////////////////
// message GetTransactionReceiptProofRequest

// reader

type GetTransactionReceiptProofRequest struct {
	// ProtocolVersion primitives.ProtocolVersion
	// VirtualChainId primitives.VirtualChainId
	// TransactionTimestamp primitives.TimestampNano
	// Txhash primitives.Sha256

	// internal
	// implements membuffers.Message
	_message membuffers.InternalMessage
}

func (x *GetTransactionReceiptProofRequest) String() string {
	if x == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{ProtocolVersion:%s,VirtualChainId:%s,TransactionTimestamp:%s,Txhash:%s,}", x.StringProtocolVersion(), x.StringVirtualChainId(), x.StringTransactionTimestamp(), x.StringTxhash())
}

var _GetTransactionReceiptProofRequest_Scheme = []membuffers.FieldType{membuffers.TypeUint32, membuffers.TypeUint32, membuffers.TypeUint64, membuffers.TypeBytes}
var _GetTransactionReceiptProofRequest_Unions = [][]membuffers.FieldType{}

func GetTransactionReceiptProofRequestReader(buf []byte) *GetTransactionReceiptProofRequest {
	x := &GetTransactionReceiptProofRequest{}
	x._message.Init(buf, membuffers.Offset(len(buf)), _GetTransactionReceiptProofRequest_Scheme, _GetTransactionReceiptProofRequest_Unions)
	return x
}

func (x *GetTransactionReceiptProofRequest) IsValid() bool {
	return x._message.IsValid()
}

func (x *GetTransactionReceiptProofRequest) Raw() []byte {
	return x._message.RawBuffer()
}

func (x *GetTransactionReceiptProofRequest) Equal(y *GetTransactionReceiptProofRequest) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	return bytes.Equal(x.Raw(), y.Raw())
}

func (x *GetTransactionReceiptProofRequest) ProtocolVersion() primitives.ProtocolVersion {
	return primitives.ProtocolVersion(x._message.GetUint32(0))
}

func (x *GetTransactionReceiptProofRequest) RawProtocolVersion() []byte {
	return x._message.RawBufferForField(0, 0)
}

func (x *GetTransactionReceiptProofRequest) MutateProtocolVersion(v primitives.ProtocolVersion) error {
	return x._message.SetUint32(0, uint32(v))
}

func (x *GetTransactionReceiptProofRequest) StringProtocolVersion() string {
	return fmt.Sprintf("%s", x.ProtocolVersion())
}

func (x *GetTransactionReceiptProofRequest) VirtualChainId() primitives.VirtualChainId {
	return primitives.VirtualChainId(x._message.GetUint32(1))
}

func (x *GetTransactionReceiptProofRequest) RawVirtualChainId() []byte {
	return x._message.RawBufferForField(1, 0)
}

func (x *GetTransactionReceiptProofRequest) MutateVirtualChainId(v primitives.VirtualChainId) error {
	return x._message.SetUint32(1, uint32(v))
}

func (x *GetTransactionReceiptProofRequest) StringVirtualChainId() string {
	return fmt.Sprintf("%s", x.VirtualChainId())
}

func (x *GetTransactionReceiptProofRequest) TransactionTimestamp() primitives.TimestampNano {
	return primitives.TimestampNano(x._message.GetUint64(2))
}

func (x *GetTransactionReceiptProofRequest) RawTransactionTimestamp() []byte {
	return x._message.RawBufferForField(2, 0)
}

func (x *GetTransactionReceiptProofRequest) MutateTransactionTimestamp(v primitives.TimestampNano) error {
	return x._message.SetUint64(2, uint64(v))
}

func (x *GetTransactionReceiptProofRequest) StringTransactionTimestamp() string {
	return fmt.Sprintf("%s", x.TransactionTimestamp())
}

func (x *GetTransactionReceiptProofRequest) Txhash() primitives.Sha256 {
	return primitives.Sha256(x._message.GetBytes(3))
}

func (x *GetTransactionReceiptProofRequest) RawTxhash() []byte {
	return x._message.RawBufferForField(3, 0)
}

func (x *GetTransactionReceiptProofRequest) RawTxhashWithHeader() []byte {
	return x._message.RawBufferWithHeaderForField(3, 0)
}

func (x *GetTransactionReceiptProofRequest) MutateTxhash(v primitives.Sha256) error {
	return x._message.SetBytes(3, []byte(v))
}

func (x *GetTransactionReceiptProofRequest) StringTxhash() string {
	return fmt.Sprintf("%s", x.Txhash())
}

// builder

type GetTransactionReceiptProofRequestBuilder struct {
	ProtocolVersion      primitives.ProtocolVersion
	VirtualChainId       primitives.VirtualChainId
	TransactionTimestamp primitives.TimestampNano
	Txhash               primitives.Sha256

	// internal
	// implements membuffers.Builder
	_builder               membuffers.InternalBuilder
	_overrideWithRawBuffer []byte
}

func (w *GetTransactionReceiptProofRequestBuilder) Write(buf []byte) (err error) {
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
	w._builder.WriteUint32(buf, uint32(w.ProtocolVersion))
	w._builder.WriteUint32(buf, uint32(w.VirtualChainId))
	w._builder.WriteUint64(buf, uint64(w.TransactionTimestamp))
	w._builder.WriteBytes(buf, []byte(w.Txhash))
	return nil
}

func (w *GetTransactionReceiptProofRequestBuilder) HexDump(prefix string, offsetFromStart membuffers.Offset) (err error) {
	if w == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = &membuffers.ErrBufferOverrun{}
		}
	}()
	w._builder.Reset()
	w._builder.HexDumpUint32(prefix, offsetFromStart, "GetTransactionReceiptProofRequest.ProtocolVersion", uint32(w.ProtocolVersion))
	w._builder.HexDumpUint32(prefix, offsetFromStart, "GetTransactionReceiptProofRequest.VirtualChainId", uint32(w.VirtualChainId))
	w._builder.HexDumpUint64(prefix, offsetFromStart, "GetTransactionReceiptProofRequest.TransactionTimestamp", uint64(w.TransactionTimestamp))
	w._builder.HexDumpBytes(prefix, offsetFromStart, "GetTransactionReceiptProofRequest.Txhash", []byte(w.Txhash))
	return nil
}

func (w *GetTransactionReceiptProofRequestBuilder) GetSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	return w._builder.GetSize()
}

func (w *GetTransactionReceiptProofRequestBuilder) CalcRequiredSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	w.Write(nil)
	return w._builder.GetSize()
}

func (w *GetTransactionReceiptProofRequestBuilder) Build() *GetTransactionReceiptProofRequest {
	buf := make([]byte, w.CalcRequiredSize())
	if w.Write(buf) != nil {
		return nil
	}
	return GetTransactionReceiptProofRequestReader(buf)
}

func GetTransactionReceiptProofRequestBuilderFromRaw(raw []byte) *GetTransactionReceiptProofRequestBuilder {
	return &GetTransactionReceiptProofRequestBuilder{_overrideWithRawBuffer: raw}
}

/////////////////////////////////////////////////////////////////////////////
// message GetTransactionReceiptProofResponse

// reader

type GetTransactionReceiptProofResponse struct {
	// RequestStatus protocol.RequestStatus
	// PackedProof primitives.PackedReceiptProof
	// TransactionStatus protocol.TransactionStatus
	// BlockHeight primitives.BlockHeight
	// BlockTimestamp primitives.TimestampNano
	// PackedReceipt primitives.PackedReceipt

	// internal
	// implements membuffers.Message
	_message membuffers.InternalMessage
}

func (x *GetTransactionReceiptProofResponse) String() string {
	if x == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{RequestStatus:%s,PackedProof:%s,TransactionStatus:%s,BlockHeight:%s,BlockTimestamp:%s,PackedReceipt:%s,}", x.StringRequestStatus(), x.StringPackedProof(), x.StringTransactionStatus(), x.StringBlockHeight(), x.StringBlockTimestamp(), x.StringPackedReceipt())
}

var _GetTransactionReceiptProofResponse_Scheme = []membuffers.FieldType{membuffers.TypeUint16, membuffers.TypeBytes, membuffers.TypeUint16, membuffers.TypeUint64, membuffers.TypeUint64, membuffers.TypeBytes}
var _GetTransactionReceiptProofResponse_Unions = [][]membuffers.FieldType{}

func GetTransactionReceiptProofResponseReader(buf []byte) *GetTransactionReceiptProofResponse {
	x := &GetTransactionReceiptProofResponse{}
	x._message.Init(buf, membuffers.Offset(len(buf)), _GetTransactionReceiptProofResponse_Scheme, _GetTransactionReceiptProofResponse_Unions)
	return x
}

func (x *GetTransactionReceiptProofResponse) IsValid() bool {
	return x._message.IsValid()
}

func (x *GetTransactionReceiptProofResponse) Raw() []byte {
	return x._message.RawBuffer()
}

func (x *GetTransactionReceiptProofResponse) Equal(y *GetTransactionReceiptProofResponse) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	return bytes.Equal(x.Raw(), y.Raw())
}

func (x *GetTransactionReceiptProofResponse) RequestStatus() protocol.RequestStatus {
	return protocol.RequestStatus(x._message.GetUint16(0))
}

func (x *GetTransactionReceiptProofResponse) RawRequestStatus() []byte {
	return x._message.RawBufferForField(0, 0)
}

func (x *GetTransactionReceiptProofResponse) MutateRequestStatus(v protocol.RequestStatus) error {
	return x._message.SetUint16(0, uint16(v))
}

func (x *GetTransactionReceiptProofResponse) StringRequestStatus() string {
	return x.RequestStatus().String()
}

func (x *GetTransactionReceiptProofResponse) PackedProof() primitives.PackedReceiptProof {
	return primitives.PackedReceiptProof(x._message.GetBytes(1))
}

func (x *GetTransactionReceiptProofResponse) RawPackedProof() []byte {
	return x._message.RawBufferForField(1, 0)
}

func (x *GetTransactionReceiptProofResponse) RawPackedProofWithHeader() []byte {
	return x._message.RawBufferWithHeaderForField(1, 0)
}

func (x *GetTransactionReceiptProofResponse) MutatePackedProof(v primitives.PackedReceiptProof) error {
	return x._message.SetBytes(1, []byte(v))
}

func (x *GetTransactionReceiptProofResponse) StringPackedProof() string {
	return fmt.Sprintf("%s", x.PackedProof())
}

func (x *GetTransactionReceiptProofResponse) TransactionStatus() protocol.TransactionStatus {
	return protocol.TransactionStatus(x._message.GetUint16(2))
}

func (x *GetTransactionReceiptProofResponse) RawTransactionStatus() []byte {
	return x._message.RawBufferForField(2, 0)
}

func (x *GetTransactionReceiptProofResponse) MutateTransactionStatus(v protocol.TransactionStatus) error {
	return x._message.SetUint16(2, uint16(v))
}

func (x *GetTransactionReceiptProofResponse) StringTransactionStatus() string {
	return x.TransactionStatus().String()
}

func (x *GetTransactionReceiptProofResponse) BlockHeight() primitives.BlockHeight {
	return primitives.BlockHeight(x._message.GetUint64(3))
}

func (x *GetTransactionReceiptProofResponse) RawBlockHeight() []byte {
	return x._message.RawBufferForField(3, 0)
}

func (x *GetTransactionReceiptProofResponse) MutateBlockHeight(v primitives.BlockHeight) error {
	return x._message.SetUint64(3, uint64(v))
}

func (x *GetTransactionReceiptProofResponse) StringBlockHeight() string {
	return fmt.Sprintf("%s", x.BlockHeight())
}

func (x *GetTransactionReceiptProofResponse) BlockTimestamp() primitives.TimestampNano {
	return primitives.TimestampNano(x._message.GetUint64(4))
}

func (x *GetTransactionReceiptProofResponse) RawBlockTimestamp() []byte {
	return x._message.RawBufferForField(4, 0)
}

func (x *GetTransactionReceiptProofResponse) MutateBlockTimestamp(v primitives.TimestampNano) error {
	return x._message.SetUint64(4, uint64(v))
}

func (x *GetTransactionReceiptProofResponse) StringBlockTimestamp() string {
	return fmt.Sprintf("%s", x.BlockTimestamp())
}

func (x *GetTransactionReceiptProofResponse) PackedReceipt() primitives.PackedReceipt {
	return primitives.PackedReceipt(x._message.GetBytes(5))
}

func (x *GetTransactionReceiptProofResponse) RawPackedReceipt() []byte {
	return x._message.RawBufferForField(5, 0)
}

func (x *GetTransactionReceiptProofResponse) RawPackedReceiptWithHeader() []byte {
	return x._message.RawBufferWithHeaderForField(5, 0)
}

func (x *GetTransactionReceiptProofResponse) MutatePackedReceipt(v primitives.PackedReceipt) error {
	return x._message.SetBytes(5, []byte(v))
}

func (x *GetTransactionReceiptProofResponse) StringPackedReceipt() string {
	return fmt.Sprintf("%s", x.PackedReceipt())
}

// builder

type GetTransactionReceiptProofResponseBuilder struct {
	RequestStatus     protocol.RequestStatus
	PackedProof       primitives.PackedReceiptProof
	TransactionStatus protocol.TransactionStatus
	BlockHeight       primitives.BlockHeight
	BlockTimestamp    primitives.TimestampNano
	PackedReceipt     primitives.PackedReceipt

	// internal
	// implements membuffers.Builder
	_builder               membuffers.InternalBuilder
	_overrideWithRawBuffer []byte
}

func (w *GetTransactionReceiptProofResponseBuilder) Write(buf []byte) (err error) {
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
	w._builder.WriteUint16(buf, uint16(w.RequestStatus))
	w._builder.WriteBytes(buf, []byte(w.PackedProof))
	w._builder.WriteUint16(buf, uint16(w.TransactionStatus))
	w._builder.WriteUint64(buf, uint64(w.BlockHeight))
	w._builder.WriteUint64(buf, uint64(w.BlockTimestamp))
	w._builder.WriteBytes(buf, []byte(w.PackedReceipt))
	return nil
}

func (w *GetTransactionReceiptProofResponseBuilder) HexDump(prefix string, offsetFromStart membuffers.Offset) (err error) {
	if w == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = &membuffers.ErrBufferOverrun{}
		}
	}()
	w._builder.Reset()
	w._builder.HexDumpUint16(prefix, offsetFromStart, "GetTransactionReceiptProofResponse.RequestStatus", uint16(w.RequestStatus))
	w._builder.HexDumpBytes(prefix, offsetFromStart, "GetTransactionReceiptProofResponse.PackedProof", []byte(w.PackedProof))
	w._builder.HexDumpUint16(prefix, offsetFromStart, "GetTransactionReceiptProofResponse.TransactionStatus", uint16(w.TransactionStatus))
	w._builder.HexDumpUint64(prefix, offsetFromStart, "GetTransactionReceiptProofResponse.BlockHeight", uint64(w.BlockHeight))
	w._builder.HexDumpUint64(prefix, offsetFromStart, "GetTransactionReceiptProofResponse.BlockTimestamp", uint64(w.BlockTimestamp))
	w._builder.HexDumpBytes(prefix, offsetFromStart, "GetTransactionReceiptProofResponse.PackedReceipt", []byte(w.PackedReceipt))
	return nil
}

func (w *GetTransactionReceiptProofResponseBuilder) GetSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	return w._builder.GetSize()
}

func (w *GetTransactionReceiptProofResponseBuilder) CalcRequiredSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	w.Write(nil)
	return w._builder.GetSize()
}

func (w *GetTransactionReceiptProofResponseBuilder) Build() *GetTransactionReceiptProofResponse {
	buf := make([]byte, w.CalcRequiredSize())
	if w.Write(buf) != nil {
		return nil
	}
	return GetTransactionReceiptProofResponseReader(buf)
}

func GetTransactionReceiptProofResponseBuilderFromRaw(raw []byte) *GetTransactionReceiptProofResponseBuilder {
	return &GetTransactionReceiptProofResponseBuilder{_overrideWithRawBuffer: raw}
}

/////////////////////////////////////////////////////////////////////////////
// enums
