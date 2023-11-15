package types

import (
	"reflect"
	"sort"

	"github.com/CosmWasm/cosmwasm-go/std/types"
)

// https://docs.rs/cw-utils/1.0.2/src/cw_utils/balance.rs.html

// NATIVE BALANCE
type NativeBalance struct {
	Coins []types.Coin
}

// Normalize sorts the wallet by denom, removes 0 elements, and consolidates duplicate denoms.
func (n *NativeBalance) Normalize() {
	// Drop 0's
	n.Coins = removeZeroAmounts(n.Coins)

	// Sort
	sortCoinsByDenom(n.Coins)

	// Consolidate duplicate denoms
	consolidateDuplicates(n.Coins)
}

// removeZeroAmounts removes coins with amount 0 from the slice.
func removeZeroAmounts(coins []types.Coin) []types.Coin {
	var result []types.Coin
	for _, c := range coins {
		if !reflect.ValueOf(c.Amount).IsZero() {
			result = append(result, c)
		}
	}
	return result
}

// sortCoinsByDenom sorts coins by denom in the slice.
func sortCoinsByDenom(coins []types.Coin) {
	sort.Slice(coins, func(i, j int) bool {
		return coins[i].Denom < coins[j].Denom
	})
}

// consolidateDuplicates consolidates coins with the same denom in the slice.
func consolidateDuplicates(coins []types.Coin) {
	var consolidated []types.Coin
	for i, c := range coins {
		if i > 0 && c.Denom == coins[i-1].Denom {
			this := consolidated[len(consolidated)-1]
			this.Amount.Add(c.Amount)
		} else {
			consolidated = append(consolidated, c)
		}
	}
	copy(coins, consolidated)
	coins = coins[:len(consolidated)]
}

// EXPIRATION

type Expiration struct {
	// !todo
}
