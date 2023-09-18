package base

type PayloadKafka struct {
	Body       string `json:"body"`
	Properties []any  `json:"properties"`
	Headers    []any  `json:"headers"`
}
