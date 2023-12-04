package src

import (
	"errors"
	"os"

	contractTypes "github.com/JackalLabs/burrow-contracts/cw1-subkeys/src/types"

	"github.com/CosmWasm/cosmwasm-go/std"
	"github.com/CosmWasm/cosmwasm-go/std/types"
	cw1WhiteList "github.com/JackalLabs/burrow-contracts/cw1-whitelist/src"
)

// var _ std.InstantiateFunc = Instantiate

func Instantiate(deps *std.Deps, env types.Env, info types.MessageInfo, data []byte) (*types.Response, error) {
	res, err := cw1WhiteList.Instantiate(deps, env, info, data)
	if err != nil {
		return nil, err
	}

	// version info for migration info
	// !toreview https://github.com/CosmWasm/cw-plus/blob/main/contracts/cw1-subkeys/src/contract.rs#L33-L44
	name := "cw1-subkeys"
	version := os.Getenv("PKG_VERSION")
	if len(version) == 0 {
		panic("No pkg version found")
	}
	SetContractVersion(deps, name, version)

	return res, nil
}

func SetContractVersion(deps *std.Deps, name string, version string) {
	info := contractTypes.ContractInfo{
		Name:    name,
		Version: version,
	}
	SaveContractInfo(deps.Storage, &info)
}

func Migrate(deps *std.Deps, env types.Env, msg []byte) (*types.Response, error) {
	return nil, errors.New("cannot migrate this contract")
}

func Execute(deps *std.Deps, env types.Env, info types.MessageInfo, data []byte) (*types.Response, error) {
	msg := contractTypes.ExecuteMsg{}
	err := msg.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}

	// we need to find which one is non-empty
	switch {
	case msg.ExecuteRequest != nil:
		return cw1WhiteList.ExecuteExecute(deps, &env, &info, msg.ExecuteRequest)
	case msg.FreezeRequest != nil:
		return cw1WhiteList.ExecuteFreeze(deps, &env, &info, msg.FreezeRequest)
	case msg.UpdateAdminsRequest != nil:
		return cw1WhiteList.ExecuteUpdateAdmins(deps, &env, &info, msg.UpdateAdminsRequest)
	case msg.IncreaseAllowance != nil:
		return executeIncreaseAllowance(deps, &env, &info, msg.IncreaseAllowance)
	case msg.DecreaseAllowance != nil:
		return executeDecreaseAllowance(deps, &env, &info, msg.DecreaseAllowance)
	case msg.SetPermissions != nil:
		return SetPermissions()

	default:
		return nil, types.GenericError("Unknown ExecuteMsg")
	}
}

func Query(deps *std.Deps, env types.Env, data []byte) ([]byte, error) {
	msg := contractTypes.QueryMsg{}
	err := msg.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}

	// we need to find which one is non-empty
	var res std.JSONType
	switch {
	case msg.QueryAdminListRequest != nil:
		res, err = queryAdminList(deps, &env, msg.QueryAdminListRequest)
	case msg.QueryCanExecuteRequest != nil:
		res, err = queryCanExecute(deps, &env, msg.QueryCanExecuteRequest)
	default:
		err = types.GenericError("Unknown QueryMsg " + string(data))
	}
	if err != nil {
		return nil, err
	}

	// if we got a result above, encode it
	bz, err := res.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func executeIncreaseAllowance(deps *std.Deps, env *types.Env, info *types.MessageInfo, msg *contractTypes.IncreaseAllowance) (*types.Response, error) {
	sender := info.Sender
	state, err := cw1WhiteList.LoadState(deps.Storage)
	if err != nil {
		return nil, err
	}

	// check if sender is admin
	if !state.IsAdmin(sender) {
		return nil, errors.New("Unauthorized")
	}

	err = deps.Api.ValidateAddress(msg.Spender)
	if err != nil {
		return nil, err
	}

	// sender can't be spender
	if msg.Spender == sender {
		return nil, errors.New("Cannot Set Your own Account")
	}

	var allow contractTypes.Allowances
	var emptyExpiration contractTypes.Expiration

	prev, err := LoadAllowances(deps.Storage, msg.Spender)

	if msg.Expires != emptyExpiration {
		if msg.Expires.IsExpired(env.Block) {
			return nil, errors.New("setting expired allowance")
		}

		allow.Expires = msg.Expires
		allow.Balance = contractTypes.NativeBalance{
			Coins: []types.Coin{msg.Amount},
		}
	} else if prev.Expires.IsExpired(env.Block) {
		return nil, errors.New("setting expired allowance")
	} else {
		allow.Balance = prev.Balance.AddAssign(msg.Amount)
	}

	SaveAllowances(deps.Storage, msg.Spender, &allow)
	res := &types.Response{
		Attributes: []types.EventAttribute{
			{Key: "action", Value: "increase_allowance"},
			{Key: "owner", Value: sender},
			{Key: "spender", Value: msg.Spender},
			{Key: "denom", Value: msg.Amount.Denom},
			{Key: "amount", Value: msg.Amount.Amount.String()},
		},
	}
	return res, nil
}

func executeDecreaseAllowance(deps *std.Deps, env *types.Env, info *types.MessageInfo, msg *contractTypes.DecreaseAllowance) (*types.Response, error) {
	sender := info.Sender
	state, err := cw1WhiteList.LoadState(deps.Storage)
	if err != nil {
		return nil, err
	}

	// check if sender is admin
	if !state.IsAdmin(sender) {
		return nil, errors.New("Unauthorized")
	}

	err = deps.Api.ValidateAddress(msg.Spender)
	if err != nil {
		return nil, err
	}

	// sender can't be spender
	if msg.Spender == sender {
		return nil, errors.New("Cannot Set Your own Account")
	}

	var allow contractTypes.Allowances
	var emptyExpiration contractTypes.Expiration

	prev, err := LoadAllowances(deps.Storage, msg.Spender)

	if msg.Expires != emptyExpiration {
		if msg.Expires.IsExpired(env.Block) {
			return nil, errors.New("setting expired allowance")
		}

		allow.Expires = msg.Expires
		allow.Balance = contractTypes.NativeBalance{
			Coins: []types.Coin{msg.Amount},
		}
	} else if prev.Expires.IsExpired(env.Block) {
		return nil, errors.New("setting expired allowance")
	} else {
		allow.Balance, _ = prev.Balance.SubSaturating(msg.Amount)
	}

	SaveAllowances(deps.Storage, msg.Spender, &allow)
	res := &types.Response{
		Attributes: []types.EventAttribute{
			{Key: "action", Value: "decrease_allowance"},
			{Key: "owner", Value: sender},
			{Key: "spender", Value: msg.Spender},
			{Key: "denom", Value: msg.Amount.Denom},
			{Key: "amount", Value: msg.Amount.Amount.String()},
		},
	}
	return res, nil
}

func queryAdminList(deps *std.Deps, env *types.Env, msg *contractTypes.QueryAdminListRequest) (*contractTypes.AdminListResponse, error) {
	state, err := LoadState(deps.Storage)
	if err != nil {
		return nil, err
	}

	_ = env
	_ = msg

	return &contractTypes.AdminListResponse{
		Admins:  state.Admins,
		Mutable: state.Mutable,
	}, nil
}

func queryCanExecute(deps *std.Deps, env *types.Env, msg *contractTypes.QueryCanExecuteRequest) (*contractTypes.CanExecuteResponse, error) {
	state, err := LoadState(deps.Storage)
	if err != nil {
		return nil, err
	}

	can := state.IsAdmin(msg.Sender)

	_ = env

	return &contractTypes.CanExecuteResponse{
		CanExecute: can,
	}, nil
}
