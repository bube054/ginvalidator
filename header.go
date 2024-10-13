package ginvalidator

type Header struct {
	field      string
	errFmtFunc *ErrFmtFuncHandler
}

func (h Header) CreateChain() ValidationChain {
	return NewValidationChain(h.field, h.errFmtFunc, headerLocation)
}

func NewHeader(field string, errFmtFunc *ErrFmtFuncHandler) Header {
	return Header{
		field:      field,
		errFmtFunc: errFmtFunc,
	}
}
