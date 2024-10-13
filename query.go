package ginvalidator

type Query struct {
	field      string
	errFmtFunc *ErrFmtFuncHandler
}

func (q Query) CreateChain() ValidationChain {
	return NewValidationChain(q.field, q.errFmtFunc, queryLocation)
}

func NewQuery(field string, errFmtFunc *ErrFmtFuncHandler) Query {
	return Query{
		field:      field,
		errFmtFunc: errFmtFunc,
	}
}
