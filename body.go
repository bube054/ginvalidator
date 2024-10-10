package ginvalidator

// const bodyRequestLocation = "body"

type Body struct {
	field      string
	errFmtFunc *ErrFmtFuncHandler
}

func (b Body) CreateChain() VMS {
	return NewVMS(b.field, b.errFmtFunc, bodyLocation)
}

func NewBody(field string, errFmtFunc *ErrFmtFuncHandler) Body {
	return Body{
		field:      field,
		errFmtFunc: errFmtFunc,
	}
}
