package src

import (
	"errors"

	contractTypes "github.com/JackalLabs/burrow-contracts/example/src/types"

	"github.com/CosmWasm/cosmwasm-go/std"
	"github.com/CosmWasm/cosmwasm-go/std/types"
)

var _ std.InstantiateFunc = Instantiate

func Instantiate(deps *std.Deps, env types.Env, info types.MessageInfo, msg []byte) (*types.Response, error) {
	deps.Api.Debug("Launching Example! 🚀")

	initMsg := contractTypes.InitMsg{}
	err := initMsg.UnmarshalJSON(msg)
	if err != nil {
		return nil, err
	}

	state := contractTypes.State{
		ExampleStateField: initMsg.ExampleDetails,
	}

	err = SaveState(deps.Storage, &state)
	if err != nil {
		return nil, err
	}
	res := &types.Response{
		Attributes: []types.EventAttribute{
			{"example_details", initMsg.ExampleDetails},
		},
	}
	return res, nil
}

func Migrate(deps *std.Deps, env types.Env, msg []byte) (*types.Response, error) {
	return nil, errors.New("cannot migrate this contract")
}

func Execute(deps *std.Deps, env types.Env, info types.MessageInfo, data []byte) (*types.Response, error) {
	msg := contractTypes.HandleMsg{}
	err := msg.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}

	// we need to find which one is non-empty
	switch {
	case msg.ExampleMsg != nil:
		return executeExample(deps, &env, &info, msg.ExampleMsg)
	default:
		return nil, types.GenericError("Unknown HandleMsg")
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
