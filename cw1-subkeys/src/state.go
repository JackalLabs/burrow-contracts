package src

import (
	"errors"

	"github.com/CosmWasm/cosmwasm-go/std"
	contractTypes "github.com/JackalLabs/burrow-contracts/cw1-subkeys/src/types"
)

var (
	PERMISSIONS   = []byte("permissions")
	ALLOWANCES    = []byte("allowances")
	CONTRACT_INFO = []byte("contract_info")
)

func LoadContractInfo(storage std.Storage) (*contractTypes.ContractInfo, error) {
	data := storage.Get(CONTRACT_INFO)
	if data == nil {
		return nil, errors.New("state not found")
	}

	var state contractTypes.ContractInfo
	err := state.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}

	return &state, nil
}

func SaveContractInfo(storage std.Storage, state *contractTypes.ContractInfo) error {
	bz, err := state.MarshalJSON()
	if err != nil {
		return err
	}

	storage.Set(CONTRACT_INFO, bz)

	return nil
}

func LoadPermissions(storage std.Storage) (*contractTypes.Permissions, error) {
	data := storage.Get(PERMISSIONS)
	if data == nil {
		return nil, errors.New("state not found") // TODO(fdymylja): replace when errors API is ready
	}

	var state contractTypes.Permissions
	err := state.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}

	return &state, nil
}

func SavePermissions(storage std.Storage, state *contractTypes.Permissions) error {
	bz, err := state.MarshalJSON()
	if err != nil {
		return err
	}

	storage.Set(PERMISSIONS, bz)

	return nil
}

func LoadAllowances(storage std.Storage) (*contractTypes.Allowances, error) {
	data := storage.Get(ALLOWANCES)
	if data == nil {
		return nil, errors.New("state not found")
	}

	var state contractTypes.Allowances
	err := state.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}

	return &state, nil
}

func SaveAllowances(storage std.Storage, state *contractTypes.Allowances) error {
	bz, err := state.MarshalJSON()
	if err != nil {
		return err
	}

	storage.Set(ALLOWANCES, bz)

	return nil
}
