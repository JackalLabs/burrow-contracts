package types

type Permissions struct {
	Delegate   bool `json:"delegate"`
	Redelegate bool `json:"redelegate"`
	Undelegate bool `json:"undelegate"`
	Withdraw   bool `json:"withdraw"`
}

type Allowances struct {
	Balance NativeBalance `json:"native_balance"`
	Expires Expiration    `json:"expiration"`
}

type ContractInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type BigMap struct {
	Keys map[string][]byte
}
