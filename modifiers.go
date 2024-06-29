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
	rules        validationProcessesRules
	processType  string
}

func (m *modifier) createProcessorFromModifier() processor {
	return processor{
		validator: validator{
			field:        m.field,
			errorMessage: m.errorMessage,
			location:     defaultParamLocation,
			rules:        m.rules,
		},
		modifier: *m,
		sanitizer: sanitizer{
			field:        m.field,
			errorMessage: m.errorMessage,
			location:     defaultParamLocation,
			rules:        m.rules,
		},
	}
}

func (m modifier) Not() processor {
	not := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := m.location
		path := field
		typ := "____"
		newValue := value
		funcName := "IsNot"
		isValid := true

		return newValidationProcessResponse(location, "", path, typ, newValue, funcName, isValid, false)
	}

	m.rules = append(m.rules, not)

	return m.createProcessorFromModifier()
}

func (m modifier) Bail() processor {
	not := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := m.location
		path := field
		typ := "____"
		newValue := value
		funcName := "Bail"
		isValid := true

		return newValidationProcessResponse(location, "", path, typ, newValue, funcName, isValid, true)
	}

	m.rules = append(m.rules, not)

	return m.createProcessorFromModifier()
}

type IfFunc func(*gin.Context) bool

func (m modifier) If(ifFunc IfFunc) processor {
	iF := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := m.location
		path := field
		typ := "____"
		newValue := value
		funcName := "If"
		isValid := true
		shouldBail := !ifFunc(ctx)

		return newValidationProcessResponse(location, "", path, typ, newValue, funcName, isValid, shouldBail)
	}

	m.rules = append(m.rules, iF)

	return m.createProcessorFromModifier()
}
