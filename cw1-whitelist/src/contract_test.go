package src

import (
	"encoding/json"
	"testing"

	contractTypes "github.com/JackalLabs/burrow-contracts/cw1-whitelist/src/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CosmWasm/cosmwasm-go/std"
	"github.com/CosmWasm/cosmwasm-go/std/math"
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

var FUND = []types.Coin{types.NewCoin(math.NewUint128FromUint64(98765), "ujkl")}

// this can be used for a quick setup if you don't have any other requirements
func defaultInit(t *testing.T, funds []types.Coin) (*std.Deps, types.Env) {
	deps := mock.Deps(funds)
	env := mock.Env()
	info := mock.Info(FUNDER, funds)
	initMsg := contractTypes.InitMsg{
		Admins:  []string{"alice", "bob", "charlie"},
		Mutable: true,
	}
	res, err := Instantiate(deps, env, info, mustEncode(t, initMsg))
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, res.Attributes[0].Key, "success")
	return deps, env
}

func TestInitAndQuery(t *testing.T) {
	deps, env := defaultInit(t, FUND)

	// Query Admin List
	qmsg := []byte(`{"admin_list":{}}`)
	data, err := Query(deps, env, qmsg)

	require.NoError(t, err)

	var qres contractTypes.AdminListResponse

	err = json.Unmarshal(data, &qres)
	require.NoError(t, err)
	assert.Equal(t, []string{"alice", "bob", "charlie"}, qres.Admins)
}
