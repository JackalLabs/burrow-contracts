package types

import "slices"

type State struct {
	ExampleStateField string `json:"example_field"`
}

type AdminList struct {
	Admins  []string `json:"admins"`
	Mutable bool     `json:"mutable"`
}

func (a AdminList) IsAdmin(addr string) bool {
	contain := slices.Contains(a.Admins, addr)
	return contain
}

func (a AdminList) CanModify(addr string) bool {
	return a.IsAdmin(addr) && a.Mutable
}
