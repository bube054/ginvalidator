package ginvalidator

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type modifier struct {
	field           string
	errorMessage    string
	location        string
	rules           validationChainRules
	chainMethodType string
}

func (m *modifier) createValidationChainFromModifier() validationChain {
	return validationChain{
		validator: validator{
			field:        m.field,
			errorMessage: m.errorMessage,
			location:     m.location,
			rules:        m.rules,
		},
		modifier: *m,
		sanitizer: sanitizer{
			field:        m.field,
			errorMessage: m.errorMessage,
			location:     m.location,
			rules:        m.rules,
		},
	}
}

// A modifier that negates the result of the next validator in the chain.
func (m modifier) Not() validationChain {
	not := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := m.location
		path := field
		newValue := value
		funcName := notFunc
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	m.rules = append(m.rules, not)

	return m.createValidationChainFromModifier()
}

// A modifier that stops running the validation chain if any of the previous validators failed.
// This is useful to prevent a custom validator that touches a database or external API from running when you know it will fail.
func (m modifier) Bail() validationChain {
	bail := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := m.location
		path := field
		newValue := value
		funcName := bailFunc
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, true)
	}

	m.rules = append(m.rules, bail)

	return m.createValidationChainFromModifier()
}

type Continue func(value string, req http.Request, location string) bool

// A modifier that adds a condition on whether the validation chain should continue running on a field or not.
func (m modifier) If(ifFunc Continue) validationChain {
	iF := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := m.location
		path := field
		newValue := value
		funcName := iFFunc
		isValid := true
		shouldBail := !ifFunc(value, *ctx.Request, location)

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, shouldBail)
	}

	m.rules = append(m.rules, iF)

	return m.createValidationChainFromModifier()
}
