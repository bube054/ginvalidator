package ginvalidator

import "github.com/gin-gonic/gin"

type Param struct {
	field      string
	errFmtFunc *ErrFmtFuncHandler
}

func (p Param) CreateChain() ValidationChain {
	return NewValidationChain(p.field, p.errFmtFunc, paramLocation)
}

func NewParam(field string, errFmtFunc *ErrFmtFuncHandler) Param {
	return Param{
		field:      field,
		errFmtFunc: errFmtFunc,
	}
}

func extractParamValue(field string, ctx *gin.Context) (string, error) {
	param := ctx.Param(field)

	return param, nil
}
