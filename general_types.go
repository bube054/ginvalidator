package ginvalidator

import (
	"github.com/gin-gonic/gin"
)

// the one that will readable by the packages user
type ReturnableValidationChainResponse struct {
	location string // where the value to be validated lives e.g body, cookies, headers, params, query.
	msg      string // the validation error message. empty string for no errors.
	path     string // the str
	value    string
}

// for internal use only
type validationChainResponse struct {
	ReturnableValidationChainResponse
	newValue   string // is the initial value to be validated unless a modifier has be called.
	funcName   string // the name of the validator, could be a validator, modifier, sanitizer.
	isValid    bool   // whether the validation was valid or not.
	shouldBail bool   // whether the validation should stop at once.
}

// creates a new validator
func newValidationChainResponse(location, msg, path, newValue, funcName string, isValid bool, shouldBail bool) validationChainResponse {
	return validationChainResponse{
		ReturnableValidationChainResponse: ReturnableValidationChainResponse{
			location: location,
			msg:      msg,
			path:     path,
		},
		newValue:   newValue,
		funcName:   funcName,
		isValid:    isValid,
		shouldBail: shouldBail,
	}
}

// the func/slice of functions for returning a validator response, should be called in the final validator function.
type validationChainRule func(value, field string, ctx *gin.Context) validationChainResponse
type validationChainRules []validationChainRule
