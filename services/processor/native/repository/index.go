package repository

import (
	sdkContext "github.com/orbs-network/orbs-contract-sdk/go/context"
	"github.com/orbs-network/orbs-network-go/services/processor/native/repository/ASBEthereum"
	"github.com/orbs-network/orbs-network-go/services/processor/native/repository/BenchmarkContract"
	"github.com/orbs-network/orbs-network-go/services/processor/native/repository/BenchmarkToken"
	"github.com/orbs-network/orbs-network-go/services/processor/native/repository/ERC20Proxy"
	"github.com/orbs-network/orbs-network-go/services/processor/native/repository/_Deployments"
	"github.com/orbs-network/orbs-network-go/services/processor/native/repository/_GlobalPreOrder"
	"github.com/orbs-network/orbs-network-go/services/processor/native/repository/_Info"
)

var PreBuiltContracts = map[string]*sdkContext.ContractInfo{
	globalpreorder_systemcontract.CONTRACT_NAME: {
		PublicMethods: globalpreorder_systemcontract.PUBLIC,
		Permission:    sdkContext.PERMISSION_SCOPE_SYSTEM,
	},
	deployments_systemcontract.CONTRACT_NAME: {
		PublicMethods: deployments_systemcontract.PUBLIC,
		Permission:    sdkContext.PERMISSION_SCOPE_SYSTEM,
	},
	info_systemcontract.CONTRACT_NAME: {
		PublicMethods: info_systemcontract.PUBLIC,
		Permission:    sdkContext.PERMISSION_SCOPE_SYSTEM,
	},
	benchmarkcontract.CONTRACT_NAME: {
		PublicMethods: benchmarkcontract.PUBLIC,
		SystemMethods: benchmarkcontract.SYSTEM,
		Permission:    sdkContext.PERMISSION_SCOPE_SERVICE,
	},
	benchmarktoken.CONTRACT_NAME: {
		PublicMethods: benchmarktoken.PUBLIC,
		SystemMethods: benchmarktoken.SYSTEM,
		Permission:    sdkContext.PERMISSION_SCOPE_SERVICE,
	},
	asb_ether.CONTRACT_NAME: {
		PublicMethods: asb_ether.PUBLIC,
		SystemMethods: asb_ether.SYSTEM,
		EventsMethods: asb_ether.EVENTS,
		Permission:    sdkContext.PERMISSION_SCOPE_SERVICE,
	},
	erc20proxy.CONTRACT_NAME: {
		PublicMethods: erc20proxy.PUBLIC,
		SystemMethods: erc20proxy.SYSTEM,
		Permission:    sdkContext.PERMISSION_SCOPE_SERVICE,
	},
	// add new pre-built native system contracts here
}
