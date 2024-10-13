package ginvalidator

type Cookie struct {
	field      string
	errFmtFunc *ErrFmtFuncHandler
}

func (c Cookie) CreateChain() ValidationChain {
	return NewValidationChain(c.field, c.errFmtFunc, cookieLocation)
}

func NewCookie(field string, errFmtFunc *ErrFmtFuncHandler) Cookie {
	return Cookie{
		field:      field,
		errFmtFunc: errFmtFunc,
	}
}
