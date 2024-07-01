package ginvalidator

import (
	// valid "github.com/asaskevich/govalidator"
	"strings"

	"github.com/gin-gonic/gin"
)

type sanitizer struct {
	field           string
	errorMessage    string
	location        string
	rules           validationChainRules
	chainMethodType string
}

func (s *sanitizer) createValidationChainFromSanitizer() validationChain {
	return validationChain{
		validator: validator{
			field:        s.field,
			errorMessage: s.errorMessage,
			location:     s.location,
			rules:        s.rules,
		},
		modifier: modifier{
			field:        s.field,
			errorMessage: s.errorMessage,
			location:     s.location,
			rules:        s.rules,
		},
		sanitizer: *s,
	}
}

func (s sanitizer) Default(defaultValue string) validationChain {
	default_ := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := s.location
		path := field

		var newValue = value

		if valueIsNullish(newValue) {
			newValue = defaultValue
		}

		funcName := "Default"
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, default_)

	return s.createValidationChainFromSanitizer()
}

func (s sanitizer) Replace(valuesFrom []string, valueTo string) validationChain {
	replace := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := s.location
		path := field

		var newValue = value

		if valueIsInSlice(newValue, valuesFrom) {
			newValue = valueTo
		}

		funcName := "Replace"
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, replace)

	return s.createValidationChainFromSanitizer()
}

func (s sanitizer) ToLowerCase() validationChain {
	toLowerCase := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := s.location
		path := field

		newValue := strings.ToLower(value)

		funcName := "Replace"
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toLowerCase)

	return s.createValidationChainFromSanitizer()
}

func (s sanitizer) ToUpperCase() validationChain {
	toUpperCase := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := s.location
		path := field

		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createValidationChainFromSanitizer()
}

func (s sanitizer) Blacklist(chars string) validationChain {
	toUpperCase := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := s.location
		path := field

		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createValidationChainFromSanitizer()
}

func (s sanitizer) Escape() validationChain {
	toUpperCase := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := s.location
		path := field

		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createValidationChainFromSanitizer()
}

func (s sanitizer) Unescape() validationChain {
	toUpperCase := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := s.location
		path := field

		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createValidationChainFromSanitizer()
}

func (s sanitizer) LTrim() validationChain {
	toUpperCase := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := s.location
		path := field

		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createValidationChainFromSanitizer()
}

func (s sanitizer) normalizeEmail() validationChain {
	toUpperCase := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := s.location
		path := field

		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createValidationChainFromSanitizer()
}

func (s sanitizer) RTrim() validationChain {
	toUpperCase := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := s.location
		path := field

		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createValidationChainFromSanitizer()
}

func (s sanitizer) toBoolean() validationChain {
	toUpperCase := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := s.location
		path := field

		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createValidationChainFromSanitizer()
}

func (s sanitizer) toDate() validationChain {
	toUpperCase := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := s.location
		path := field

		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createValidationChainFromSanitizer()
}

func (s sanitizer) ToFloat() validationChain {
	toUpperCase := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := s.location
		path := field

		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createValidationChainFromSanitizer()
}

func (s sanitizer) ToInt() validationChain {
	toUpperCase := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := s.location
		path := field

		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createValidationChainFromSanitizer()
}

func (s sanitizer) Trim() validationChain {
	toUpperCase := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := s.location
		path := field

		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createValidationChainFromSanitizer()
}

func (s sanitizer) Whitelist(chars string) validationChain {
	toUpperCase := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := s.location
		path := field

		newValue := strings.ToUpper(value)

		funcName := "Replace"
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createValidationChainFromSanitizer()
}
