package main

import (
	"github.com/orbs-network/orbs-contract-sdk/go/sdk"
)

var CONTRACT = sdk.ContractInfo{
	Name:       "CounterFrom200",
	Permission: sdk.PERMISSION_SCOPE_SERVICE,
	Methods: map[string]sdk.MethodInfo{
		METHOD_INIT.Name:  METHOD_INIT,
		METHOD_ADD.Name:   METHOD_ADD,
		METHOD_GET.Name:   METHOD_GET,
		METHOD_START.Name: METHOD_START,
	},
	InitSingleton: newContract,
}

func newContract(base *sdk.BaseContract) sdk.ContractInstance {
	return &contract{base}
}

type contract struct{ *sdk.BaseContract }

///////////////////////////////////////////////////////////////////////////

var METHOD_INIT = sdk.MethodInfo{
	Name:           "_init",
	External:       false,
	Access:         sdk.ACCESS_SCOPE_READ_WRITE,
	Implementation: (*contract)._init,
}

func (c *contract) _init(ctx sdk.Context) error {
	return c.State.WriteUint64ByKey(ctx, "count", 200)
}

///////////////////////////////////////////////////////////////////////////

var METHOD_ADD = sdk.MethodInfo{
	Name:           "add",
	External:       true,
	Access:         sdk.ACCESS_SCOPE_READ_WRITE,
	Implementation: (*contract).add,
}

func (c *contract) add(ctx sdk.Context, amount uint64) error {
	count, err := c.State.ReadUint64ByKey(ctx, "count")
	if err != nil {
		return err
	}
	count += amount
	return c.State.WriteUint64ByKey(ctx, "count", count)
}

///////////////////////////////////////////////////////////////////////////

var METHOD_GET = sdk.MethodInfo{
	Name:           "get",
	External:       true,
	Access:         sdk.ACCESS_SCOPE_READ_ONLY,
	Implementation: (*contract).get,
}

func (c *contract) get(ctx sdk.Context) (uint64, error) {
	return c.State.ReadUint64ByKey(ctx, "count")
}

///////////////////////////////////////////////////////////////////////////

var METHOD_START = sdk.MethodInfo{
	Name:           "start",
	External:       true,
	Access:         sdk.ACCESS_SCOPE_READ_ONLY,
	Implementation: (*contract).start,
}

func (c *contract) start(ctx sdk.Context) (uint64, error) {
	return 200, nil
}
