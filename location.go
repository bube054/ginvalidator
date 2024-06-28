package ginvalidator

type location int

const (
	body = iota
	cookies
	headers
	params
	query
)
