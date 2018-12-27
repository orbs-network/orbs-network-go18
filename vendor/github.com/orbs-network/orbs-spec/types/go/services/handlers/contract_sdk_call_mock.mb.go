// AUTO GENERATED FILE (by membufc proto compiler v0.0.21)
package handlers

import (
	"context"
	"github.com/orbs-network/go-mock"
)

/////////////////////////////////////////////////////////////////////////////
// service ContractSdkCallHandler

type MockContractSdkCallHandler struct {
	mock.Mock
}

func (s *MockContractSdkCallHandler) HandleSdkCall(ctx context.Context, input *HandleSdkCallInput) (*HandleSdkCallOutput, error) {
	ret := s.Called(ctx, input)
	if out := ret.Get(0); out != nil {
		return out.(*HandleSdkCallOutput), ret.Error(1)
	} else {
		return nil, ret.Error(1)
	}
}
