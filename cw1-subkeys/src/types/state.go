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
