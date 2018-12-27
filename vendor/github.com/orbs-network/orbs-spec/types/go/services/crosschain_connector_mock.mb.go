// AUTO GENERATED FILE (by membufc proto compiler v0.0.21)
package services

import (
	"context"
	"github.com/orbs-network/go-mock"
)

/////////////////////////////////////////////////////////////////////////////
// service CrosschainConnector

type MockCrosschainConnector struct {
	mock.Mock
}

func (s *MockCrosschainConnector) EthereumCallContract(ctx context.Context, input *EthereumCallContractInput) (*EthereumCallContractOutput, error) {
	ret := s.Called(ctx, input)
	if out := ret.Get(0); out != nil {
		return out.(*EthereumCallContractOutput), ret.Error(1)
	} else {
		return nil, ret.Error(1)
	}
}

func (s *MockCrosschainConnector) EthereumGetTransactionLogs(ctx context.Context, input *EthereumGetTransactionLogsInput) (*EthereumGetTransactionLogsOutput, error) {
	ret := s.Called(ctx, input)
	if out := ret.Get(0); out != nil {
		return out.(*EthereumGetTransactionLogsOutput), ret.Error(1)
	} else {
		return nil, ret.Error(1)
	}
}
