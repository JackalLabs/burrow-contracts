package types

import "github.com/CosmWasm/cosmwasm-go/std/types"

type InitMsg struct {
	Admins  []string `json:"admins"`
	Mutable bool     `json:"mutable"`
}

type MigrateMsg struct{}

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
}

type QueryMsg struct {
	QueryAdminListRequest  *QueryAdminListRequest  `json:"admin_list,omitempty"`
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
	Sender string          `json:"sender,omitempty"`
	Msg    types.CosmosMsg `json:"msg,omitempty"`
}

// Responses
type AdminListResponse struct {
	Admins  []string `json:"admins,omitempty"`
	Mutable bool     `json:"mutable,omitempty"`
}

type CanExecuteResponse struct {
	CanExecute bool `json:"can_execute,omitempty"`
}
