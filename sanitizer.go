package ginvalidator

import (
	// valid "github.com/asaskevich/govalidator"

	"net/http"
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

type CustomSanitizerFunc func(value string, req http.Request, location string) string

// A sanitizer that makes the newly returned value by the function to become the new value of the field.
func (s sanitizer) CustomSanitizer(customSan CustomSanitizerFunc) validationChain {
	custom := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := s.location
		path := field
		newValue := customSan(value, *ctx.Request, location)
		funcName := customSanitizerFunc
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, custom)

	return s.createValidationChainFromSanitizer()
}

// A sanitizer that replaces the value of the field if it's either an empty string, null, undefined, or NaN.
func (s sanitizer) Default(defaultValue string) validationChain {
	default_ := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := s.location
		path := field
		var newValue = value
		if valueIsNullish(newValue) {
			newValue = defaultValue
		}
		funcName := defaultFunc
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, default_)

	return s.createValidationChainFromSanitizer()
}

// A sanitizer that replaces the value of the field with valueTo whenever the current value is in valuesFrom.
func (s sanitizer) Replace(valuesFrom []string, valueTo string) validationChain {
	replace := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := s.location
		path := field

		var newValue = value

		if valueIsInSlice(newValue, valuesFrom) {
			newValue = valueTo
		}

		funcName := replaceFunc
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

		funcName := toLowerCaseFunc
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

		funcName := toUpperCaseFunc
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

		funcName := blacklistFunc
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

		funcName := escapeFunc
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

		funcName := unescapeFunc
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

		funcName := ltrimFunc
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createValidationChainFromSanitizer()
}

func (s sanitizer) NormalizeEmail() validationChain {
	toUpperCase := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := s.location
		path := field

		newValue := strings.ToUpper(value)

		funcName := normalizeEmailFunc
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

		funcName := rtrimFunc
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createValidationChainFromSanitizer()
}

func (s sanitizer) ToBoolean() validationChain {
	toUpperCase := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := s.location
		path := field

		newValue := strings.ToUpper(value)

		funcName := toBooleanFunc
		isValid := true

		return newValidationChainResponse(location, "", path, newValue, funcName, isValid, false)
	}

	s.rules = append(s.rules, toUpperCase)

	return s.createValidationChainFromSanitizer()
}

func (s sanitizer) ToDate() validationChain {
	toUpperCase := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := s.location
		path := field

		newValue := strings.ToUpper(value)

		funcName := toDateFunc
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

		funcName := toFloatFunc
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
