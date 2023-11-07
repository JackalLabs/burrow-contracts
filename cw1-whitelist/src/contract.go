package src

import (
	"errors"

	contractTypes "github.com/JackalLabs/burrow-contracts/cw1-whitelist/src/types"

	"github.com/CosmWasm/cosmwasm-go/std"
	"github.com/CosmWasm/cosmwasm-go/std/types"
)

// var _ std.InstantiateFunc = Instantiate

func Instantiate(deps *std.Deps, env types.Env, info types.MessageInfo, data []byte) (*types.Response, error) {
	deps.Api.Debug("Launching Example! 🚀")

	initMsg := contractTypes.InitMsg{}
	err := initMsg.UnmarshalJSON(data)
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
			{Key: "success", Value: "true"},
		},
	}
	return res, nil
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

func executeExecute(deps *std.Deps, env *types.Env, info *types.MessageInfo, msg *contractTypes.ExecuteRequest) (*types.Response, error) {
	sender := info.Sender

	state, err := LoadState(deps.Storage)
	if err != nil {
		return nil, err
	}

	if !state.IsAdmin(sender) {
		return nil, err
	}

	_ = env

	var messages []types.SubMsg

	for _, msg := range msg.Msgs {
		newSub := types.NewSubMsg(msg)
		messages = append(messages, newSub)
	}

	res := &types.Response{
		Attributes: []types.EventAttribute{
			{Key: "action", Value: "execute"},
		},
		Messages: messages,
	}
	return res, nil
}

func executeFreeze(deps *std.Deps, env *types.Env, info *types.MessageInfo, msg *contractTypes.FreezeRequest) (*types.Response, error) {
	sender := info.Sender

	state, err := LoadState(deps.Storage)
	if err != nil {
		return nil, err
	}

	if !state.IsAdmin(sender) {
		return nil, err
	}

	state.Mutable = false

	err = SaveState(deps.Storage, state)
	if err != nil {
		return nil, err
	}

	res := &types.Response{
		Attributes: []types.EventAttribute{
			{Key: "action", Value: "freeze"},
		},
	}
	return res, nil
}

func executeUpdateAdmins(deps *std.Deps, env *types.Env, info *types.MessageInfo, msg *contractTypes.UpdateAdminsRequest) (*types.Response, error) {
	sender := info.Sender

	state, err := LoadState(deps.Storage)
	if err != nil {
		return nil, err
	}

	if !state.CanModify(sender) {
		return nil, err
	}

	state.Admins = msg.Admins

	err = SaveState(deps.Storage, state)
	if err != nil {
		return nil, err
	}

	res := &types.Response{
		Attributes: []types.EventAttribute{
			{Key: "action", Value: "freeze"},
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
