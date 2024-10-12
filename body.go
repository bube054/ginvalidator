package ginvalidator

import (
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

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

func extractBodyValue(field string, ctx *gin.Context) (string, error) {
	data, err := ctx.GetRawData()

	if err != nil {
		return "", err
	}

	json := string(data)

	result := gjson.Get(json, field)

	return result.String(), nil
}
