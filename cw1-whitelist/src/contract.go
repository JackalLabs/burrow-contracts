package src

import (
	"encoding/json"
	"errors"

	contractTypes "github.com/JackalLabs/burrow-contracts/cw1-whitelist/src/types"

	"github.com/CosmWasm/cosmwasm-go/std"
	"github.com/CosmWasm/cosmwasm-go/std/types"
)

var _ std.InstantiateFunc = Instantiate

func Instantiate(deps *std.Deps, env types.Env, info types.MessageInfo, data []byte) (*types.Response, error) {
	deps.Api.Debug("Launching Example! 🚀")

	initMsg := contractTypes.InitMsg{}
	// err := initMsg.UnmarshalJSON(msg)
	err := json.Unmarshal(data, initMsg)
	if err != nil {
		return nil, err
	}

	state := contractTypes.AdminList{
		Admins:  initMsg.Admins,
		Mutable: initMsg.Mutable,
	}

	err = SaveState(deps.Storage, &state)
	if err != nil {
		return nil, err
	}
	res := &types.Response{
		Attributes: []types.EventAttribute{
			{"success", true},
		},
	}
	return res, nil
}

func Migrate(deps *std.Deps, env types.Env, msg []byte) (*types.Response, error) {
	return nil, errors.New("cannot migrate this contract")
}

func Execute(deps *std.Deps, env types.Env, info types.MessageInfo, data []byte) (*types.Response, error) {
	msg := contractTypes.ExecuteMsg{}
	// err := msg.UnmarshalJSON(data)
	err := json.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	// we need to find which one is non-empty
	switch {
	case msg.ExecuteRequest != nil:
		return executeExecute(deps, &env, &info, msg.ExecuteRequest)
	case msg.FreezeRequest != nil:
		return executeFreeze(deps, &env, &info, msg.FreezeRequest)
	case msg.UpdateAdminsRequest != nil:
		return executeUpdateAdmins(deps, &env, &info, msg.UpdateAdminsRequest)
	default:
		return nil, types.GenericError("Unknown ExecuteMsg")
	}
}

func Query(deps *std.Deps, env types.Env, data []byte) ([]byte, error) {
	msg := contractTypes.QueryMsg{}
	// err := msg.UnmarshalJSON(data)
	err := json.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	// we need to find which one is non-empty
	var res std.JSONType
	switch {
	case msg.ExampleQuery != nil:
		res, err = queryExample(deps, &env, msg.ExampleQuery)
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

func executeExecute(deps *std.Deps, env *types.Env, info *types.MessageInfo, msg *contractTypes.ExecuteRequest) (*types.Response, error) {
	sender := info.Sender

	state, err := LoadState(deps.Storage)
	if err != nil {
		return nil, err
	}

	state.isAdmin(sender)

	err := example(deps, info.Sender)
	if err != nil {
		return nil, err
	}

	_ = env

	deps.Api.Debug(msg.Msgs)

	return &types.Response{
		Attributes: []types.EventAttribute{
			{"action", "example"},
		},
	}, nil
}

func executeFreeze(deps *std.Deps, env *types.Env, info *types.MessageInfo, msg *contractTypes.ExecuteRequest) (*types.Response, error) {
	// !todo
}

func executeUpdateAdmins(deps *std.Deps, env *types.Env, info *types.MessageInfo, msg *contractTypes.ExecuteRequest) (*types.Response, error) {
	// !todo
}
