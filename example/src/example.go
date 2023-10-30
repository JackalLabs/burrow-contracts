package src

import (
	contractTypes "github.com/JackalLabs/burrow-contracts/example/src/types"

	"github.com/CosmWasm/cosmwasm-go/std"
	"github.com/CosmWasm/cosmwasm-go/std/types"
)

func executeExample(deps *std.Deps, env *types.Env, info *types.MessageInfo, msg *contractTypes.ExampleMsgReqeust) (*types.Response, error) {
	err := example(deps, info.Sender)
	if err != nil {
		return nil, err
	}

	_ = env

	deps.Api.Debug(msg.ExampleField)

	return &types.Response{
		Attributes: []types.EventAttribute{
			{"action", "example"},
		},
	}, nil
}

func example(deps *std.Deps, sender types.HumanAddress) error {

	deps.Storage.Set([]byte(sender), []byte("example!"))

	return nil
}

func queryExample(deps *std.Deps, env *types.Env, msg *contractTypes.ExampleQueryRequest) (*contractTypes.ExampleQueryResponse, error) {
	state, err := LoadState(deps.Storage)
	if err != nil {
		return nil, err
	}

	_ = state
	_ = env
	_ = msg

	return &contractTypes.ExampleQueryResponse{
		ExampleField: "example response",
	}, nil
}
