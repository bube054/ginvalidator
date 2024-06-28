// params.go
package ginvalidator

import (
	"github.com/gin-gonic/gin"
)

var defaultParamLocation = "query"

func NewParam(ctx *gin.Context, field, errorMessage string) Processor {
	value := ctx.Query(field)

	return Processor{
		Validator: Validator{
			field:        field,
			errorMessage: errorMessage,
			value:        value,
			currentValue: value,
			resErrs:      make([]ResponseError, 0),
			location:     defaultParamLocation,
		},
		Modifier: Modifier{
			field:        field,
			errorMessage: errorMessage,
			value:        value,
			currentValue: value,
			resErrs:      make([]ResponseError, 0),
			location:     defaultParamLocation,
		},
		Sanitizer: Sanitizer{
			field:        field,
			errorMessage: errorMessage,
			value:        value,
			currentValue: value,
			resErrs:      make([]ResponseError, 0),
			location:     defaultParamLocation,
		},
	}
}
