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
		fmt.Println("field:", field)
		currentValue := ctx.Param(field)
		fmt.Println("value/param:", currentValue)

		fmt.Printf("the rules are %+v\n", p.validator.rules)

		results := make([]validationProcessResponse, 0)
		for i, validator := range p.validator.rules {
			result := validator(currentValue, field)

			currentValue = result.newValue
			results = append(results, result)
			fmt.Printf("result at %d is %v\n", i+1, result)
		}

		ctx.Set("__GIN__VALIDATOR__", results)

		ctx.Next()
	}
}

type validationProcessResponse struct {
	location string
	msg      string
	path     string
	typ      string

	newValue string
	funcName string
	isValid  bool
}

func newValidationProcessResponse(location, msg, path, typ, newValue, funcName string, isValid bool) validationProcessResponse {
	return validationProcessResponse{
		location: location,
		msg:      msg,
		path:     path,
		typ:      typ,
		newValue: newValue,
		funcName: funcName,
		isValid:  isValid,
	}
}

type validationProcessesRule func(value, field string) validationProcessResponse
type validationProcessesRules []validationProcessesRule
