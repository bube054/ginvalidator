package ginvalidator

import "github.com/gin-gonic/gin"

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

	return m.createValidationChainFromModifier()
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

	return m.createValidationChainFromModifier()
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

	return m.createValidationChainFromModifier()
}
