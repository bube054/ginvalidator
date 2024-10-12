package ginvalidator

import (
	"github.com/gin-gonic/gin"
)

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

func extractCookieValue(field string, ctx *gin.Context) (string, error) {
	cookie, err := ctx.Cookie(field)

	if err != nil {
		return "", err
	}

	return cookie, nil
}
