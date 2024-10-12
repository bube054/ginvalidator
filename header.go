package ginvalidator

import "github.com/gin-gonic/gin"

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

func extractHeaderValue(field string, ctx *gin.Context) (string, error) {
	header := ctx.Request.Header[field][0]

	return header, nil
}
