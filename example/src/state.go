package src

import (
	"errors"

	"github.com/JackalLabs/burrow-contracts/example/src/types"

	"github.com/CosmWasm/cosmwasm-go/std"
)

var StateKey = []byte("config")

func LoadState(storage std.Storage) (*types.State, error) {
	data := storage.Get(StateKey)
	if data == nil {
		return nil, errors.New("state not found") // TODO(fdymylja): replace when errors API is ready
	}

	var state types.State
	err := state.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func SaveState(storage std.Storage, state *types.State) error {
	bz, err := state.MarshalJSON()
	if err != nil {
		return err
	}

	storage.Set(StateKey, bz)

	return nil
}
