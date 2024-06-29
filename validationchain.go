package ginvalidator

import (
	"fmt"
	"io"
	"log"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
)

// the chain of either validators, sanitizers and modifiers
type validationChain struct {
	validator
	modifier
	sanitizer
}

func (vc validationChain) Validate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// contentTypeHeader := ctx.GetHeader("Content-Type")

		// if contentTypeHeader != "application/json" {
		// 	ctx.Next()
		// }

		// field := vc.validator.field
		// currentValue := ctx.Param(field)

		// responses := make([]validationChainResponse, 0)
		// for _, validator := range vc.validator.rules {
		// 	response := validator(currentValue, field, ctx)

		// 	if response.shouldBail || response.funcName == "Bail" {
		// 		break
		// 	}

		// 	currentValue = response.newValue
		// 	responses = append(responses, response)
		// }

		// fmt.Println("responses:", responses)
		// ctx.Set("__GIN__VALIDATOR__PARAM__VALIDATION__PROCESS__RESPONSES__", responses)
		field := vc.getValidationField()
		sanitizedValue, err := vc.getValueToValidate(ctx)
		responses := make([]validationChainResponse, 0)

		if err != nil {
			// handle errors more efficiently later.
			log.Println("a validation error has occurred:", err)
			ctx.Next()
		}

		for _, rules := range vc.getValidationRules() {
			response := rules(sanitizedValue, field, ctx)

			if response.shouldBail || response.funcName == "Bail" {
				break
			}

			sanitizedValue = response.newValue
			responses = append(responses, response)
		}

		ctx.Next()
	}
}

// gets the validation field
func (vc validationChain) getValidationField() string {
	return vc.validator.field // could have also access modifier or sanitizer.
}

// gets the validation rules
func (vc validationChain) getValidationRules() validationChainRules {
	return vc.validator.rules // could have also access modifier or sanitizer.
}

// gets the validation location
func (vc validationChain) getValidationLocation() string {
	return vc.validator.location // could have also access modifier or sanitizer.
}

// get the field to be validated from validation location
func (vc validationChain) getValueToValidate(ctx *gin.Context) (string, error) {
	field := vc.getValidationField()
	switch vc.getValidationLocation() {
	case bodyLocation:
		keys, err := splitJSONFieldSelector(field)
		if err != nil {
			return "", err
		}
		reqBody, err := ctx.Request.GetBody()
		if err != nil {
			return "", err
		}
		reqBodyBytes, err := io.ReadAll(reqBody)
		if err != nil {
			return "", err
		}
		key, _, _, err := jsonparser.Get(reqBodyBytes, keys...)
		if err != nil {
			return "", err
		}
		return string(key), nil
	case cookiesLocation:
		return ctx.Cookie(field)
	case headersLocation:
		return ctx.GetHeader(field), nil
	case paramsLocation:
		return ctx.Param(field), nil
	case queryLocation:
		return ctx.Query(field), nil
	default:
		return "", fmt.Errorf("invalid request location for %s.", field)
	}
}
