package types

import (
	"sort"

	"github.com/CosmWasm/cosmwasm-go/std/types"
	cw1WhiteListTypes "github.com/JackalLabs/burrow-contracts/cw1-whitelist/src/types"
)

type ExecuteMsg struct {
	/// Execute requests the contract to re-dispatch all these messages with the
	/// contract's address as sender. Every implementation has it's own logic to
	/// determine in
	ExecuteRequest *ExecuteRequest `json:"execute,omitempty"`

	/// Freeze will make a mutable contract immutable, must be called by an admin
	FreezeRequest *cw1WhiteListTypes.FreezeRequest `json:"freeze,omitempty"`

	/// UpdateAdmins will change the admin set of the contract, must be called by an existing admin,
	/// and only works if the contract is mutable
	UpdateAdminsRequest *cw1WhiteListTypes.UpdateAdminsRequest `json:"update_admins,omitempty"`

	/// Add an allowance to a given subkey (subkey must not be admin)
	IncreaseAllowance *IncreaseAllowance `json:"increase_allowance,omitempty"`

	/// Decreases an allowance for a given subkey (subkey must not be admin)
	DecreaseAllowance *DecreaseAllowance `json:"decrease_allowance,omitempty"`

	// Setups up permissions for a given subkey.
	SetPermissions *SetPermissions `json:"set_permissions,omitempty"`
}

type QueryMsg struct {
	/// Shows all admins and whether or not it is mutable
	QueryAdminListRequest *QueryAdminListRequest `json:"admin_list,omitempty"`

	/// Checks permissions of the caller on this proxy.
	/// If CanExecute returns true then a call to `Execute` with the same message,
	/// before any further state changes, should also succeed.
	QueryCanExecuteRequest *QueryCanExecuteRequest `json:"can_execute,omitempty"`

	/// Get the current allowance for the given subkey (how much it can spend)
	QueryAllowance *QueryAllowance `json:"allowance,omitempty"`

	/// Get the current permissions for the given subkey (how much it can spend)
	QueryPermissions *QueryPermissions `json:"permissions,omitempty"`

	/// Gets all Allowances for this contract
	QueryAllAllowance *QueryAllAllowance `json:"all_allowance,omitempty"`

	/// Gets all Permissions for this contract
	QueryAllPermissions *QueryAllPermissions `json:"all_permissions,omitempty"`
}

// Requests
type ExecuteRequest struct {
	Msgs []types.CosmosMsg `json:"msgs,omitempty"`
}

type UpdateAdminsRequest struct {
	Admins []string `json:"admins,omitempty"`
}

type IncreaseAllowance struct {
	Spender string
	Amount  types.Coin
	Expires Expiration
}
type DecreaseAllowance struct {
	Spender string
	Amount  types.Coin
	Expires Expiration
}
type SetPermissions struct {
	Spender     string
	Permissions Permissions
}

type QueryAdminListRequest struct{}

type QueryCanExecuteRequest struct {
	Sender string          `json:"sender,omitempty"`
	Msg    types.CosmosMsg `json:"msg,omitempty"`
}

type QueryAllowance struct {
	Spender string `json:"spender,omitempty"`
}

type QueryPermissions struct {
	Spender string `json:"spender,omitempty"`
}

type QueryAllAllowance struct {
	StartAfter string `json:"start_after,omitempty"`
	Limit      uint32 `json:"limit,omitempty"`
}

type QueryAllPermissions struct {
	StartAfter string `json:"start_after,omitempty"`
	Limit      uint32 `json:"limit,omitempty"`
}

// Responses
type AdminListResponse struct {
	Admins  []string `json:"admins"`
	Mutable bool     `json:"mutable"`
}

type CanExecuteResponse struct {
	CanExecute bool `json:"can_execute"`
}

// / -Allowance
type AllAllowancesResponse struct {
	Allowances []AllowanceInfo `json:"allowances"`
}

func (r AllAllowancesResponse) Canonical() AllAllowancesResponse {
	for i := range r.Allowances {
		r.Allowances[i] = r.Allowances[i].Canonical()
	}

	sort.Slice(r.Allowances, func(i, j int) bool {
		return r.Allowances[i].CmpBySpender(r.Allowances[j])
	})

	return r
}

type AllowanceInfo struct {
	Spender string        `json:"spender"`
	Balance NativeBalance `json:"balance"`
	Expires Expiration    `json:"expires"`
}

func (i AllowanceInfo) CmpBySpender(other AllowanceInfo) bool {
	return i.Spender < other.Spender
}

func (i AllowanceInfo) Canonical() AllowanceInfo {
	i.Balance.Normalize()
	return i
}

// / -Permission
type AllPermissionsResponse struct {
	Permissions []PermissionInfo `json:"permissions"`
}

type PermissionInfo struct {
	Spender     string      `json:"spender"`
	Permissions Permissions `json:"permissions"`
}

func (i PermissionInfo) CmpBySpender(other PermissionInfo) bool {
	return i.Spender < other.Spender
}
