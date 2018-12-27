package virtualmachine

import "github.com/orbs-network/orbs-spec/types/go/primitives"

type keyValuePair struct {
	key     []byte
	value   []byte
	isDirty bool
}

type contractTransientState struct {
	pairs map[string]*keyValuePair
}

type transientState struct {
	contracts map[primitives.ContractName]*contractTransientState
}

func newTransientState() *transientState {
	return &transientState{
		contracts: make(map[primitives.ContractName]*contractTransientState),
	}
}

func (t *transientState) getValue(contract primitives.ContractName, key []byte) ([]byte, bool) {
	c, found := t.contracts[contract]
	if !found {
		return nil, false
	}
	k := keyForMap(key)
	pair, found := c.pairs[k]
	if found {
		return pair.value, found
	} else {
		return nil, found
	}
}

func (t *transientState) setValue(contract primitives.ContractName, key []byte, value []byte, isDirty bool) {
	c, found := t.contracts[contract]
	if !found {
		c = &contractTransientState{
			pairs: make(map[string]*keyValuePair),
		}
		t.contracts[contract] = c
	}
	k := keyForMap(key)
	pair, found := c.pairs[k]
	if found {
		pair.value = value
		pair.isDirty = isDirty
	} else {
		c.pairs[k] = &keyValuePair{key, value, isDirty}
	}
}

func (t *transientState) forDirty(contract primitives.ContractName, f func(key []byte, value []byte)) {
	c, found := t.contracts[contract]
	if found {
		for _, pair := range c.pairs {
			if pair.isDirty {
				f(pair.key, pair.value)
			}
		}
	}
}

func (t *transientState) mergeIntoTransientState(masterTransientState *transientState) {
	for contractName, _ := range t.contracts {
		t.forDirty(contractName, func(key []byte, value []byte) {
			masterTransientState.setValue(contractName, key, value, true)
		})
	}
}

func keyForMap(key []byte) string {
	return string(key) // TODO(v1): improve to create a version without copy (unsafe cast)
}
