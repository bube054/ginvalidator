// params.go
package ginvalidator

import (
	// valid "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"fmt"
)

type Param struct {
	field        string
	errorMessage string

	value string
}

func NewParam(field, errorMessage string) *Param {
	return &Param{
		field:        field,
		errorMessage: errorMessage,
	}
}

func (p *Param) IsALPHA() *Param {
	return p
}

func (p *Param) IsASCII() *Param {
	return p
}

func (p *Param) Validate() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		contentType := ctx.GetHeader("Content-Type")
		fmt.Println("contentType:", contentType)

		if contentType != "application/json" {
			ctx.Next()
		}

		ctx.Next()
	}
}
