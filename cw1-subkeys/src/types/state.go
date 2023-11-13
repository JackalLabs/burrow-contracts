package types

import "github.com/CosmWasm/cosmwasm-go/std/types"

// https://docs.rs/cw-utils/1.0.2/src/cw_utils/balance.rs.html
type NativeBalance struct {
	Coins []types.Coin
}
type Permissions struct {
	Delegate   bool `json:"delegate"`
	Redelegate bool `json:"redelegate"`
	Undelegate bool `json:"undelegate"`
	Withdraw   bool `json:"withdraw"`
}

type Allowances struct {
	Balance NativeBalance `json:"native_balance"`
	// !todo
	// Expires  Expiration `json:"expiration"`
}
