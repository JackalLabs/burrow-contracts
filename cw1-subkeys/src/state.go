package src

import (
	"errors"
	"strconv"

	"github.com/CosmWasm/cosmwasm-go/std"
	contractTypes "github.com/JackalLabs/burrow-contracts/cw1-subkeys/src/types"
)

var (
	PERMISSIONS_MAP = []byte("permissions_map")
	ALLOWANCES_MAP  = []byte("allowances_map")
	CONTRACT_INFO   = []byte("contract_info")
)

func LoadPermissions(storage std.Storage, key string) (*contractTypes.Permissions, error) {
	// Load storage for map
	// Map should return a key
	// Load storage for the value of that key
	data := storage.Get(PERMISSIONS_MAP)
	if data == nil {
		return nil, errors.New("PERMISSIONS_MAP not found")
	}

	var bigMap contractTypes.BigMap
	err := bigMap.UnmarshalJSON(data)
	if err != nil {
		return nil, errors.New("bigMap not found")
	}

	byteKey := bigMap.Keys[key]
	if byteKey == nil {
		return nil, errors.New("byteKey doesn't exist")
	}

	data = storage.Get(byteKey)
	if data == nil {
		return nil, errors.New("byteKey doesn't have any value")
	}

	var permissions contractTypes.Permissions
	err = permissions.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}
	return &permissions, nil
}

func SavePermissions(storage std.Storage, spender string, permissions *contractTypes.Permissions) error {
	bz, err := permissions.MarshalJSON()
	if err != nil {
		return err
	}

	// save permission with spender as byteKey
	byteKey := []byte("permissions_" + spender)
	storage.Set(byteKey, bz)

	// save this byteKey to bigMap - update PERMISSIONS_MAP
	data := storage.Get(PERMISSIONS_MAP)
	if data == nil {
		return errors.New("PERMISSIONS_MAP not found")
	}

	var bigMap contractTypes.BigMap
	err = bigMap.UnmarshalJSON(data)
	if err != nil {
		return errors.New("bigMap not found")
	}

	bigMap.Keys[spender] = byteKey
	bz, err = bigMap.MarshalJSON()
	if err != nil {
		return err
	}
	storage.Set(PERMISSIONS_MAP, bz)

	return nil
}

func LoadAllowances(storage std.Storage, key string) (*contractTypes.Allowances, error) {
	// Load storage for map
	// Map should return a key
	// Load storage for the value of that key
	data := storage.Get(ALLOWANCES_MAP)
	if data == nil {
		return nil, errors.New("ALLOWANCES_MAP not found")
	}

	var bigMap contractTypes.BigMap
	err := bigMap.UnmarshalJSON(data)
	if err != nil {
		return nil, errors.New("bigMap not found")
	}

	byteKey := bigMap.Keys[key]
	if byteKey == nil {
		return nil, errors.New("byteKey doesn't exist")
	}

	data = storage.Get(byteKey)
	if data == nil {
		return nil, errors.New("byteKey doesn't have any value")
	}

	var allowances contractTypes.Allowances
	err = allowances.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}
	return &allowances, nil
}

func SaveAllowances(storage std.Storage, spender string, allowances *contractTypes.Allowances) error {
	bz, err := allowances.MarshalJSON()
	if err != nil {
		return err
	}

	// save allowances with spender as byteKey
	byteKey := []byte("allowances_" + spender)
	storage.Set(byteKey, bz)

	// save this byteKey to bigMap - load then update ALLOWANCES_MAP
	data := storage.Get(ALLOWANCES_MAP)
	if data == nil {
		return errors.New("ALLOWANCES_MAP not found")
	}

	var bigMap contractTypes.BigMap
	err = bigMap.UnmarshalJSON(data)
	if err != nil {
		return errors.New("bigMap not found")
	}

	bigMap.Keys[spender] = byteKey
	bz, err = bigMap.MarshalJSON()
	if err != nil {
		return err
	}
	storage.Set(ALLOWANCES_MAP, bz)

	return nil
}

func LoadAllAllowances(storage std.Storage, limit int) (*contractTypes.AllAllowancesResponse, error) {
	// load big map that contains all the keys to the allowances
	data := storage.Get(ALLOWANCES_MAP)
	if data == nil {
		return nil, errors.New("ALLOWANCES_MAP not found")
	}

	var bigMap contractTypes.BigMap
	err := bigMap.UnmarshalJSON(data)
	if err != nil {
		return nil, errors.New("bigMap not found")
	}

	var allAllow contractTypes.AllAllowancesResponse

	for i, byteKey := range bigMap.Keys {
		if i == strconv.Itoa(limit) {
			break
		}
		// load a single allowance now
		data = storage.Get(byteKey)
		if data == nil {
			return nil, errors.New("byteKey doesn't have any value")
		}

		var allow contractTypes.Allowances
		err = allow.UnmarshalJSON(data)
		if err != nil {
			return nil, err
		}

		allowInfo := contractTypes.AllowanceInfo{
			Spender: i,
			Balance: allow.Balance,
			Expires: allow.Expires,
		}

		// add allowance to all allowances
		allAllow.Allowances = append(allAllow.Allowances, allowInfo)
	}

	return &allAllow, nil
}

func LoadAllPermissions(storage std.Storage, limit int) (*contractTypes.AllPermissionsResponse, error) {
	// load big map that contains all the keys to the permissions
	data := storage.Get(PERMISSIONS_MAP)
	if data == nil {
		return nil, errors.New("PERMISSIONS_MAP not found")
	}

	var bigMap contractTypes.BigMap
	err := bigMap.UnmarshalJSON(data)
	if err != nil {
		return nil, errors.New("bigMap not found")
	}

	var allPerm contractTypes.AllPermissionsResponse

	for i, byteKey := range bigMap.Keys {
		if i == strconv.Itoa(limit) {
			break
		}
		// load a single permissions now
		data = storage.Get(byteKey)
		if data == nil {
			return nil, errors.New("byteKey doesn't have any value")
		}

		var perm contractTypes.Permissions
		err = perm.UnmarshalJSON(data)
		if err != nil {
			return nil, err
		}

		permInfo := contractTypes.PermissionInfo{
			Spender:     i,
			Permissions: perm,
		}

		// add allowance to all allowances
		allPerm.Permissions = append(allPerm.Permissions, permInfo)
	}

	return &allPerm, nil
}

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
