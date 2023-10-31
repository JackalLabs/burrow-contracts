package src

import (
	"encoding/json"
	"testing"

	types2 "github.com/JackalLabs/burrow-contracts/example/src/types"

	"github.com/stretchr/testify/require"

	"github.com/CosmWasm/cosmwasm-go/std"
	"github.com/CosmWasm/cosmwasm-go/std/mock"
	"github.com/CosmWasm/cosmwasm-go/std/types"
)

func mustEncode(t *testing.T, msg interface{}) []byte {
	bz, err := json.Marshal(msg)
	require.NoError(t, err)
	return bz
}

const (
	FUNDER = "creator"
)

// this can be used for a quick setup if you don't have nay other requirements
func defaultInit(t *testing.T, funds []types.Coin) *std.Deps {
	deps := mock.Deps(funds)
	env := mock.Env()
	info := mock.Info(FUNDER, funds)
	initMsg := types2.InitMsg{
		ExampleDetails: "example",
	}
	res, err := Instantiate(deps, env, info, mustEncode(t, initMsg))
	require.NoError(t, err)
	require.NotNil(t, res)
	return deps
}
