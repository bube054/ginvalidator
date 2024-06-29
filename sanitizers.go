package ginvalidator

import (
	// valid "github.com/asaskevich/govalidator"
	"strings"

	"github.com/gin-gonic/gin"
)

type sanitizer struct {
	field        string
	errorMessage string
	location     string
	rules        validationProcessesRules
	processType  string
}

func (s *sanitizer) createProcessorFromSanitizer() processor {
	return processor{
		validator: validator{
			field:        s.field,
			errorMessage: s.errorMessage,
			location:     defaultParamLocation,
			rules:        s.rules,
		},
		modifier: modifier{
			field:        s.field,
			errorMessage: s.errorMessage,
			location:     defaultParamLocation,
			rules:        s.rules,
		},
		sanitizer: *s,
	}
}

func (s sanitizer) Default(defaultValue string) processor {
	default_ := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := s.location
		path := field
		typ := "____"
		var newValue = value

		if valueIsNullish(newValue) {
			newValue = defaultValue
		}

		funcName := "Default"
		isValid := true

		return newValidationProcessResponse(location, "", path, typ, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, default_)

	return s.createProcessorFromSanitizer()
}

func (s sanitizer) Replace(valuesFrom []string, valueTo string) processor {
	replace := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := s.location
		path := field
		typ := "____"
		var newValue = value

		if valueIsInSlice(newValue, valuesFrom) {
			newValue = valueTo
		}

		funcName := "Replace"
		isValid := true

		return newValidationProcessResponse(location, "", path, typ, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, replace)

	return s.createProcessorFromSanitizer()
}

func (s sanitizer) ToLowerCase() processor {
	toLowerCase := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := s.location
		path := field
		typ := "____"
		newValue := strings.ToLower(value)

		funcName := "Replace"
		isValid := true

		return newValidationProcessResponse(location, "", path, typ, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toLowerCase)

	return s.createProcessorFromSanitizer()
}

func (s sanitizer) ToUpperCase() processor {
	toUpperCase := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := s.location
		path := field
		typ := "____"
		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationProcessResponse(location, "", path, typ, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createProcessorFromSanitizer()
}

func (s sanitizer) Blacklist(chars string) processor {
	toUpperCase := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := s.location
		path := field
		typ := "____"
		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationProcessResponse(location, "", path, typ, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createProcessorFromSanitizer()
}

func (s sanitizer) Escape() processor {
	toUpperCase := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := s.location
		path := field
		typ := "____"
		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationProcessResponse(location, "", path, typ, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createProcessorFromSanitizer()
}

func (s sanitizer) Unescape() processor {
	toUpperCase := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := s.location
		path := field
		typ := "____"
		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationProcessResponse(location, "", path, typ, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createProcessorFromSanitizer()
}

func (s sanitizer) LTrim() processor {
	toUpperCase := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := s.location
		path := field
		typ := "____"
		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationProcessResponse(location, "", path, typ, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createProcessorFromSanitizer()
}

func (s sanitizer) normalizeEmail() processor {
	toUpperCase := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := s.location
		path := field
		typ := "____"
		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationProcessResponse(location, "", path, typ, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createProcessorFromSanitizer()
}

func (s sanitizer) RTrim() processor {
	toUpperCase := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := s.location
		path := field
		typ := "____"
		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationProcessResponse(location, "", path, typ, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createProcessorFromSanitizer()
}

func (s sanitizer) toBoolean() processor {
	toUpperCase := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := s.location
		path := field
		typ := "____"
		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationProcessResponse(location, "", path, typ, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createProcessorFromSanitizer()
}

func (s sanitizer) toDate() processor {
	toUpperCase := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := s.location
		path := field
		typ := "____"
		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationProcessResponse(location, "", path, typ, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createProcessorFromSanitizer()
}

func (s sanitizer) ToFloat() processor {
	toUpperCase := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := s.location
		path := field
		typ := "____"
		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationProcessResponse(location, "", path, typ, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createProcessorFromSanitizer()
}

func (s sanitizer) ToInt() processor {
	toUpperCase := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := s.location
		path := field
		typ := "____"
		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationProcessResponse(location, "", path, typ, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createProcessorFromSanitizer()
}

func (s sanitizer) Trim() processor {
	toUpperCase := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := s.location
		path := field
		typ := "____"
		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationProcessResponse(location, "", path, typ, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createProcessorFromSanitizer()
}

func (s sanitizer) Whitelist(chars string) processor {
	toUpperCase := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := s.location
		path := field
		typ := "____"
		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationProcessResponse(location, "", path, typ, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createProcessorFromSanitizer()
}
