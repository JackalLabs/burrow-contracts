package types

import (
	"errors"
	"reflect"
	"sort"
	"time"

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

// Find returns the index and coin with the given denom, or nil if not found.
func (nb NativeBalance) Find(denom string) (int, *types.Coin) {
	for i, c := range nb.Coins {
		if c.Denom == denom {
			return i, &c
		}
	}
	return -1, nil
}

// InsertPos should only be called when denom is not in the wallet.
// It returns the position where denom should be inserted at (via splice).
// It returns -1 if this should be appended.
func (nb NativeBalance) InsertPos(denom string) int {
	for i, c := range nb.Coins {
		if c.Denom >= denom {
			return i
		}
	}
	return -1
}

// AddAssign adds the given coin to NativeBalance.
func (nb NativeBalance) AddAssign(other types.Coin) {
	idx, c := nb.Find(other.Denom)
	if c != nil {
		nb.Coins[idx].Amount = c.Amount.Add(other.Amount)
	} else {
		pos := nb.InsertPos(other.Denom)
		if pos != -1 {
			nb.Coins = append(nb.Coins[:pos], append([]types.Coin{other}, nb.Coins[pos:]...)...)
		} else {
			nb.Coins = append(nb.Coins, other)
		}
	}
}

// EXPIRATION

type Expiration struct {
	/// AtHeight will expire when `env.block.height` >= height
	AtHeight uint64
	/// AtTime will expire when `env.block.time` >= time
	AtTime time.Time
	/// Never will never expire. Used to express the empty variant
	Never bool
}

type Duration struct {
	Height uint64
	Time   uint64
}

// Adds more time to an Expiration
func (e Expiration) Add(d Duration) (Expiration, error) {
	if e.Never {
		return Expiration{Never: true}, nil
	}

	switch {
	case d.Time != 0 && e.AtTime.IsZero():
		return Expiration{}, errors.New("Cannot add height and time")
	case d.Height != 0:
		return Expiration{AtHeight: e.AtHeight + d.Height}, nil
	case d.Time != 0:
		return Expiration{AtTime: e.AtTime.Add(time.Second * time.Duration(d.Time))}, nil
	default:
		return Expiration{}, errors.New("Invalid Duration")
	}
}

// Checks if the expiration is expired
func (e Expiration) IsExpired(block types.BlockInfo) bool {
	switch {
	case e.AtHeight != 0:
		return block.Height >= e.AtHeight
	case !e.AtTime.IsZero():
		blockTime := time.Unix(0, int64(block.Time))
		return blockTime.After(e.AtTime)
	default:
		return false
	}
}

// DURATION //

// Create an expiration after current block
func (d Duration) After(block types.BlockInfo) Expiration {
	switch {
	case d.Height != 0:
		return Expiration{AtHeight: block.Height + d.Height}
	case d.Time != 0:
		duration := time.Second * time.Duration(d.Time)
		blockTime := time.Unix(0, int64(block.Time))
		return Expiration{AtTime: blockTime.Add(duration)}
	default:
		return Expiration{}
	}
}

// Create a Duration slightly larger than the current one, so we can use it to pass expiration point
func (d Duration) PlusOne() Duration {
	switch {
	case d.Height != 0:
		return Duration{Height: d.Height + 1}
	case d.Time != 0:
		return Duration{Time: d.Time + 1}
	default:
		return Duration{}
	}
}

// Add adds two Durations.
func (d Duration) Add(other Duration) (Duration, error) {
	switch {
	case d.Time != 0 && other.Time != 0:
		return Duration{Time: d.Time + other.Time}, nil
	case d.Height != 0 && other.Height != 0:
		return Duration{Height: d.Height + other.Height}, nil
	default:
		return Duration{}, errors.New("Cannot add height and time")
	}
}

// Multiply multiplies a Duration by a scalar.
func (d Duration) Multiply(v uint64) Duration {
	switch {
	case d.Time != 0:
		return Duration{Time: d.Time * v}
	case d.Height != 0:
		return Duration{Height: d.Height * v}
	default:
		return Duration{}
	}
}

// Compares two Expirations
// func (e Expiration) PartialCmp (other Expiration) int {
// 	!todo do we even need this?
// }
