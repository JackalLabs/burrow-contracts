package types

import "github.com/CosmWasm/cosmwasm-go/std/types"

type InitMsg struct {
	Admins  string `json:"admins"`
	Mutable bool   `json:"mutable"`
}

type MigrateMsg struct{}

type ExecuteMsg struct {
	/// Execute requests the contract to re-dispatch all these messages with the
	/// contract's address as sender. Every implementation has it's own logic to
	/// determine in
	Execute *Execute `json:"execute,omitempty"`

	/// Freeze will make a mutable contract immutable, must be called by an admin
	Freeze *Freeze `json:"freeze,omitempty"`

	/// UpdateAdmins will change the admin set of the contract, must be called by an existing admin,
	/// and only works if the contract is mutable
	UpdateAdmins *UpdateAdmins `json:"update_admins,omitempty"`
}

type QueryMsg struct {
	AdminList  *AdminList  `json:"admin_list,omitempty"`
	CanExecute *CanExecute `json:"can_execute,omitempty"`
}

// Requests
type Execute struct {
	Msgs []types.CosmosMsg `json:"msgs,omitempty"`
}

type Freeze struct{}

type UpdateAdmins struct {
	Admins []string `json:"admins,omitempty"`
}

type AdminList struct{}

type CanExecute struct {
	Sender string          `json:"sender,omitempty"`
	Msg    types.CosmosMsg `json:"msg,omitempty"`
}

// Responses
type AdminListResponse struct {
	Admins  string `json:"admins,omitempty"`
	Mutable bool   `json:"mutable,omitempty"`
}

// type ExampleQueryResponse struct {
// 	ExampleField string `json:"example_field,omitempty"`
// }
