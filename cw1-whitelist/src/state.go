package src

import (
	"encoding/json"
	"errors"

	"github.com/CosmWasm/cosmwasm-go/std"
	contractTypes "github.com/JackalLabs/burrow-contracts/cw1-whitelist/src/types"
)

var ADMIN_LIST = []byte("admin_list")

func LoadState(storage std.Storage) (*contractTypes.AdminList, error) {
	data := storage.Get(ADMIN_LIST)
	if data == nil {
		return nil, errors.New("state not found") // TODO(fdymylja): replace when errors API is ready
	}

	var state contractTypes.AdminList
	err := json.Unmarshal(data, state)
	if err != nil {
		return nil, err
	}

	return &state, nil
}

func SaveState(storage std.Storage, state *contractTypes.AdminList) error {
	bz, err := json.Marshal(state)
	if err != nil {
		return err
	}

	storage.Set(ADMIN_LIST, bz)

	return nil
}
