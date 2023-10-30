package types

type InitMsg struct {
	ExampleDetails string `json:"example_details"`
}

type MigrateMsg struct{}

type HandleMsg struct {
	ExampleMsg *ExampleMsgReqeust `json:"example_msg,omitempty"`
}

type QueryMsg struct {
	ExampleQuery *ExampleQueryRequest `json:"example_query,omitempty"`
}

type ExampleMsgReqeust struct {
	ExampleField string `json:"example_field,omitempty"`
}

type ExampleMsgResponse struct {
	ExampleField string `json:"example_field,omitempty"`
}

type ExampleQueryRequest struct {
	ExampleField string `json:"example_field,omitempty"`
}

type ExampleQueryResponse struct {
	ExampleField string `json:"example_field,omitempty"`
}
