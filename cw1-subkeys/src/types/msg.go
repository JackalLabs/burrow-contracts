package types

import (
	"github.com/CosmWasm/cosmwasm-go/std/types"
	contractTypes "github.com/JackalLabs/burrow-contracts/cw1-subkeys/src/types"
)

type ExecuteMsg struct {
	/// Execute requests the contract to re-dispatch all these messages with the
	/// contract's address as sender. Every implementation has it's own logic to
	/// determine in
	ExecuteRequest *ExecuteRequest `json:"execute,omitempty"`

	/// Freeze will make a mutable contract immutable, must be called by an admin
	FreezeRequest *FreezeRequest `json:"freeze,omitempty"`

	/// UpdateAdmins will change the admin set of the contract, must be called by an existing admin,
	/// and only works if the contract is mutable
	UpdateAdminsRequest *UpdateAdminsRequest `json:"update_admins,omitempty"`

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
}

// Requests
type ExecuteRequest struct {
	Msgs []types.CosmosMsg `json:"msgs,omitempty"`
}

type FreezeRequest struct{}

type UpdateAdminsRequest struct {
	Admins []string `json:"admins,omitempty"`
}

type QueryAdminListRequest struct{}

type QueryCanExecuteRequest struct {
	Sender string `json:"sender,omitempty"`
	// Msg    types.CosmosMsg `json:"msg,omitempty"`
}

type IncreaseAllowance struct {
	Spender string
	Amount  types.Coin
	// Expires !todo
}
type DecreaseAllowance struct {
	Spender string
	Amount  types.Coin
	// Expires !todo
}
type SetPermissions struct {
	Spender     string
	Permissions contractTypes.Permissions
}

// Responses
type AdminListResponse struct {
	Admins  []string `json:"admins"`
	Mutable bool     `json:"mutable"`
}

type CanExecuteResponse struct {
	CanExecute bool `json:"can_execute"`
}
