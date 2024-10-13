package ginvalidator

type Body struct {
	field      string
	errFmtFunc *ErrFmtFuncHandler
}

func (b Body) CreateChain() ValidationChain {
	return NewValidationChain(b.field, b.errFmtFunc, bodyLocation)
}

func NewBody(field string, errFmtFunc *ErrFmtFuncHandler) Body {
	return Body{
		field:      field,
		errFmtFunc: errFmtFunc,
	}
}
