package ginvalidator

import "github.com/gin-gonic/gin"

// valid "github.com/asaskevich/govalidator"

// import (
//   valid "github.com/asaskevich/govalidator"
// )

type modifier struct {
	field        string
	errorMessage string
	location     string
	rules        validationChainRules
	processType  string
}

func (m *modifier) createProcessorFromModifier() validationChain {
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

func (m modifier) Not() validationChain {
	not := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := m.location
		path := field
		newValue := value
		funcName := "IsNot"
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	m.rules = append(m.rules, not)

	return m.createProcessorFromModifier()
}

func (m modifier) Bail() validationChain {
	not := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := m.location
		path := field
		newValue := value
		funcName := "Bail"
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, true)
	}

	m.rules = append(m.rules, not)

	return m.createProcessorFromModifier()
}

type IfFunc func(*gin.Context) bool

func (m modifier) If(ifFunc IfFunc) validationChain {
	iF := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := m.location
		path := field
		newValue := value
		funcName := "If"
		isValid := true
		shouldBail := !ifFunc(ctx)

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, shouldBail)
	}

	m.rules = append(m.rules, iF)

	return m.createProcessorFromModifier()
}
