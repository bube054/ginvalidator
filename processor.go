package ginvalidator

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type processor struct {
	validator
	modifier
	sanitizer
}

func (p processor) Validate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		contentTypeHeader := ctx.GetHeader("Content-Type")

		if contentTypeHeader != "application/json" {
			ctx.Next()
		}

		field := p.validator.field
		currentValue := ctx.Param(field)

		responses := make([]validationProcessResponse, 0)
		for _, validator := range p.validator.rules {
			response := validator(currentValue, field, ctx)

			if response.shouldBail || response.funcName == "Bail" {
				break
			}

			currentValue = response.newValue
			responses = append(responses, response)
		}

		fmt.Println("responses:", responses)
		ctx.Set("__GIN__VALIDATOR__PARAM__VALIDATION__PROCESS__RESPONSES__", responses)

		ctx.Next()
	}
}

type validationProcessResponse struct {
	location string
	msg      string
	path     string
	typ      string

	newValue   string
	funcName   string
	isValid    bool
	shouldBail bool
}

func newValidationProcessResponse(location, msg, path, typ, newValue, funcName string, isValid bool, shouldBail bool) validationProcessResponse {
	return validationProcessResponse{
		location:   location,
		msg:        msg,
		path:       path,
		typ:        typ,
		newValue:   newValue,
		funcName:   funcName,
		isValid:    isValid,
		shouldBail: shouldBail,
	}
}

type validationProcessesRule func(value, field string, ctx *gin.Context) validationProcessResponse
type validationProcessesRules []validationProcessesRule
