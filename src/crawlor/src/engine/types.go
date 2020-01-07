package engine

type Request struct {
	url        string
	ParserFunc func([]byte) ParserResult
}

type ParserResult struct {
	Requests []Request
	items []interface{}
}
