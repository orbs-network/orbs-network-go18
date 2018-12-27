package httpserver

import (
	"bytes"
	"github.com/orbs-network/go-mock"
	"github.com/orbs-network/orbs-network-go/instrumentation/log"
	"github.com/orbs-network/orbs-network-go/instrumentation/metric"
	"github.com/orbs-network/orbs-network-go/test/builders"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/orbs-network/orbs-spec/types/go/protocol/client"
	"github.com/orbs-network/orbs-spec/types/go/services"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func makeServer(papiMock *services.MockPublicApi) HttpServer {
	logger := log.GetLogger().WithOutput(log.NewFormattingOutput(os.Stdout, log.NewHumanReadableFormatter()))

	return NewHttpServer("", logger, papiMock, metric.NewRegistry())
}

func TestHttpServer_Robots(t *testing.T) {
	s := makeServer(nil)

	req, _ := http.NewRequest("Get", "/robots.txt", nil)
	rec := httptest.NewRecorder()
	s.(*server).robots(rec, req)

	expectedResponse := `User-agent: *
Disallow: /`

	require.Equal(t, http.StatusOK, rec.Code, "should succeed")
	require.Equal(t, "text/plain", rec.Header().Get("Content-Type"), "should have our content type")
	require.Equal(t, expectedResponse, rec.Body.String(), "should have text value")
}

func TestHttpServerSendTxHandler_Basic(t *testing.T) {
	papiMock := &services.MockPublicApi{}
	response := &client.SendTransactionResponseBuilder{
		RequestStatus:      protocol.REQUEST_STATUS_COMPLETED,
		TransactionReceipt: nil,
		TransactionStatus:  protocol.TRANSACTION_STATUS_COMMITTED,
		BlockHeight:        1,
		BlockTimestamp:     primitives.TimestampNano(time.Now().Nanosecond()),
	}

	papiMock.When("SendTransaction", mock.Any, mock.Any).Times(1).Return(&services.SendTransactionOutput{ClientResponse: response.Build()})

	s := makeServer(papiMock)

	request := (&client.SendTransactionRequestBuilder{
		SignedTransaction: builders.TransferTransaction().Builder(),
	}).Build()

	req, _ := http.NewRequest("POST", "", bytes.NewReader(request.Raw()))
	rec := httptest.NewRecorder()
	s.(*server).sendTransactionHandler(rec, req)

	require.Equal(t, http.StatusOK, rec.Code, "should succeed")
}

func TestHttpServerSendTxHandler_Error(t *testing.T) {
	papiMock := &services.MockPublicApi{}

	papiMock.When("SendTransaction", mock.Any, mock.Any).Times(1).Return(nil, errors.Errorf("stam"))

	s := makeServer(papiMock)

	request := (&client.SendTransactionRequestBuilder{
		SignedTransaction: builders.TransferTransaction().Builder(),
	}).Build()

	req, _ := http.NewRequest("POST", "", bytes.NewReader(request.Raw()))
	rec := httptest.NewRecorder()
	s.(*server).sendTransactionHandler(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code, "should fail with 500")
}

func TestHttpServerCallMethod_Basic(t *testing.T) {
	papiMock := &services.MockPublicApi{}
	response := &client.CallMethodResponseBuilder{
		RequestStatus:       protocol.REQUEST_STATUS_COMPLETED,
		OutputArgumentArray: nil,
		CallMethodResult:    protocol.EXECUTION_RESULT_SUCCESS,
		BlockHeight:         1,
		BlockTimestamp:      primitives.TimestampNano(time.Now().Nanosecond()),
	}

	papiMock.When("CallMethod", mock.Any, mock.Any).Times(1).Return(&services.CallMethodOutput{ClientResponse: response.Build()})

	s := makeServer(papiMock)

	request := (&client.CallMethodRequestBuilder{
		Transaction: &protocol.TransactionBuilder{},
	}).Build()

	req, _ := http.NewRequest("POST", "", bytes.NewReader(request.Raw()))
	rec := httptest.NewRecorder()
	s.(*server).callMethodHandler(rec, req)

	require.Equal(t, http.StatusOK, rec.Code, "should succeed")
	// actual values are checked in the server_test.go as unit test of internal WriteMembuffResponse
}

func TestHttpServerCallMethod_Error(t *testing.T) {
	papiMock := &services.MockPublicApi{}

	papiMock.When("CallMethod", mock.Any, mock.Any).Times(1).Return(nil, errors.Errorf("stam"))

	s := makeServer(papiMock)

	request := (&client.CallMethodRequestBuilder{
		Transaction: &protocol.TransactionBuilder{},
	}).Build()

	req, _ := http.NewRequest("POST", "", bytes.NewReader(request.Raw()))
	rec := httptest.NewRecorder()
	s.(*server).callMethodHandler(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code, "should fail with 500")
	// actual values are checked in the server_test.go as unit test of internal writeErrorResponseAndLog
}

func TestHttpServerGetTx_Basic(t *testing.T) {
	papiMock := &services.MockPublicApi{}
	response := &client.GetTransactionStatusResponseBuilder{
		RequestStatus:      protocol.REQUEST_STATUS_COMPLETED,
		TransactionReceipt: nil,
		TransactionStatus:  protocol.TRANSACTION_STATUS_COMMITTED,
		BlockHeight:        1,
		BlockTimestamp:     primitives.TimestampNano(time.Now().Nanosecond()),
	}

	papiMock.When("GetTransactionStatus", mock.Any, mock.Any).Times(1).Return(&services.GetTransactionStatusOutput{ClientResponse: response.Build()})

	s := makeServer(papiMock)

	request := (&client.GetTransactionStatusRequestBuilder{}).Build()

	req, _ := http.NewRequest("POST", "", bytes.NewReader(request.Raw()))
	rec := httptest.NewRecorder()
	s.(*server).getTransactionStatusHandler(rec, req)

	require.Equal(t, http.StatusOK, rec.Code, "should succeed")
	// actual values are checked in the server_test.go as unit test of internal WriteMembuffResponse
}

func TestHttpServerGetTx_Error(t *testing.T) {
	papiMock := &services.MockPublicApi{}

	papiMock.When("GetTransactionStatus", mock.Any, mock.Any).Times(1).Return(nil, errors.Errorf("stam"))

	s := makeServer(papiMock)

	request := (&client.GetTransactionStatusRequestBuilder{}).Build()

	req, _ := http.NewRequest("POST", "", bytes.NewReader(request.Raw()))
	rec := httptest.NewRecorder()
	s.(*server).getTransactionStatusHandler(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code, "should fail with 500")
	// actual values are checked in the server_test.go as unit test of internal writeErrorResponseAndLog
}

func TestHttpServerGetReceipt_Basic(t *testing.T) {
	papiMock := &services.MockPublicApi{}
	response := &client.GetTransactionReceiptProofResponseBuilder{
		RequestStatus:     protocol.REQUEST_STATUS_COMPLETED,
		PackedProof:       nil,
		PackedReceipt:     nil,
		TransactionStatus: protocol.TRANSACTION_STATUS_COMMITTED,
		BlockHeight:       1,
		BlockTimestamp:    primitives.TimestampNano(time.Now().Nanosecond()),
	}

	papiMock.When("GetTransactionReceiptProof", mock.Any, mock.Any).Times(1).Return(&services.GetTransactionReceiptProofOutput{ClientResponse: response.Build()})

	s := makeServer(papiMock)

	request := (&client.GetTransactionReceiptProofRequestBuilder{}).Build()

	req, _ := http.NewRequest("POST", "", bytes.NewReader(request.Raw()))
	rec := httptest.NewRecorder()
	s.(*server).getTransactionReceiptProofHandler(rec, req)

	require.Equal(t, http.StatusOK, rec.Code, "should succeed")
	// actual values are checked in the server_test.go as unit test of internal WriteMembuffResponse
}

func TestHttpServerGetReceipt_Error(t *testing.T) {
	papiMock := &services.MockPublicApi{}

	papiMock.When("GetTransactionReceiptProof", mock.Any, mock.Any).Times(1).Return(nil, errors.Errorf("stam"))

	s := makeServer(papiMock)

	request := (&client.GetTransactionReceiptProofRequestBuilder{}).Build()

	req, _ := http.NewRequest("POST", "", bytes.NewReader(request.Raw()))
	rec := httptest.NewRecorder()
	s.(*server).getTransactionReceiptProofHandler(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code, "should fail with 500")
	// actual values are checked in the server_test.go as unit test of internal writeErrorResponseAndLog
}
