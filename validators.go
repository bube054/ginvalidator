package ginvalidator

import (
	"fmt"
	"net/http"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type validator struct {
	field        string
	errorMessage string
	location     string
	rules        validationProcessesRules
	processType  string
}

func (v *validator) createProcessorFromValidator() processor {
	return processor{
		validator: *v,
		modifier: modifier{
			field:        v.field,
			errorMessage: v.errorMessage,
			location:     defaultParamLocation,
			rules:        v.rules,
		},
		sanitizer: sanitizer{
			field:        v.field,
			errorMessage: v.errorMessage,
			location:     defaultParamLocation,
			rules:        v.rules,
		},
	}
}

type customValidatorFunc func(value string, req http.Request, location string, path string) error

func (v validator) Custom(customValidator customValidatorFunc) processor {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	custom := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := v.location

		var finalErrMessage string
		customValidatorErr := customValidator(value, *ctx.Request, v.location, ctx.Request.URL.Path)

		if customValidatorErr == nil {
			finalErrMessage = ""
		} else {
			if customValidatorErr.Error() != "" {
				finalErrMessage = customValidatorErr.Error()
			} else if v.errorMessage != "" {
				finalErrMessage = v.errorMessage
			} else {
				finalErrMessage = fmt.Sprintf("%s is not alpha.", value)
			}
		}

		path := field
		typ := "____"
		newValue := value
		funcName := "Custom"
		isValid := customValidatorErr == nil

		if previousRuleWasNegation {
			isValid = !isValid

			// if !isValid {
			// 	finalErrMessage = fmt.Sprintf("%s is not alpha.", value)
			// }
		}

		return newValidationProcessResponse(location, finalErrMessage, path, typ, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, custom)

	return v.createProcessorFromValidator()
}

func (v validator) IsArray(errorMessage string) processor {
	return v.createProcessorFromValidator()
}

func (v validator) IsObject(errorMessage string) processor {
	return v.createProcessorFromValidator()
}

func (v validator) IsString(errorMessage string) processor {
	return v.createProcessorFromValidator()
}

func (v validator) IsNotEmpty(errorMessage string) processor {
	isNotEmpty := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := v.location
		var finalErrMessage string
		if errorMessage != "" {
			finalErrMessage = errorMessage
		} else if v.errorMessage != "" {
			finalErrMessage = v.errorMessage
		} else {
			finalErrMessage = fmt.Sprintf("%s is not empty.", value)
		}
		path := field
		typ := "____"
		newValue := value
		funcName := "IsNotEmpty"
		isValid := value != ""

		return newValidationProcessResponse(location, finalErrMessage, path, typ, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isNotEmpty)

	return v.createProcessorFromValidator()
}

// standard validators

func (v validator) Contains(errorMessage string, substring string) processor {
	contains := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := v.location
		var finalErrMessage string
		if errorMessage != "" {
			finalErrMessage = errorMessage
		} else if v.errorMessage != "" {
			finalErrMessage = v.errorMessage
		} else {
			finalErrMessage = fmt.Sprintf("%s is not empty.", value)
		}
		path := field
		typ := "____"
		newValue := value
		funcName := "Contains"
		isValid := valid.Contains(value, substring)

		return newValidationProcessResponse(location, finalErrMessage, path, typ, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, contains)

	return v.createProcessorFromValidator()
}

func (v validator) Equals(errorMessage string, comparison string) processor {
	equals := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := v.location
		var finalErrMessage string
		if errorMessage != "" {
			finalErrMessage = errorMessage
		} else if v.errorMessage != "" {
			finalErrMessage = v.errorMessage
		} else {
			finalErrMessage = fmt.Sprintf("%s is not empty.", value)
		}
		path := field
		typ := "____"
		newValue := value
		funcName := "Equals"
		isValid := value == comparison

		return newValidationProcessResponse(location, finalErrMessage, path, typ, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, equals)

	return v.createProcessorFromValidator()
}

func (v validator) IsAfter(errorMessage string, comparisonTime time.Time) processor {
	isAfter := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := v.location
		var finalErrMessage string
		if errorMessage != "" {
			finalErrMessage = errorMessage
		} else if v.errorMessage != "" {
			finalErrMessage = v.errorMessage
		} else {
			finalErrMessage = fmt.Sprintf("%s is not after %s.", value, comparisonTime)
		}
		path := field
		typ := "____"
		newValue := value
		funcName := "IsAfter"

		isValid := false
		var valueTime interface{}
		valueTime = value

		valueAsTime, isTime := valueTime.(time.Time)

		if isTime {
			isValid = valueAsTime.After(comparisonTime)
		}

		return newValidationProcessResponse(location, finalErrMessage, path, typ, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isAfter)

	return v.createProcessorFromValidator()
}

func (v validator) IsAlpha(errorMessage string) processor {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	isASCII := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := v.location
		var finalErrMessage string
		if errorMessage != "" {
			finalErrMessage = errorMessage
		} else if v.errorMessage != "" {
			finalErrMessage = v.errorMessage
		} else {
			finalErrMessage = fmt.Sprintf("%s is not alpha.", value)
		}
		path := field
		typ := "____"
		newValue := value
		funcName := "IsAlpha"
		isValid := valid.IsAlpha(value)

		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationProcessResponse(location, finalErrMessage, path, typ, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isASCII)

	return v.createProcessorFromValidator()
}

func (v validator) IsAlphanumeric(errorMessage string) processor {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)
	fmt.Println("previousRuleWasNegationIsAlphaNumeric:", previousRuleWasNegation)

	isASCII := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := v.location
		var finalErrMessage string
		if errorMessage != "" {
			finalErrMessage = errorMessage
		} else if v.errorMessage != "" {
			finalErrMessage = v.errorMessage
		} else {
			finalErrMessage = fmt.Sprintf("%s is not alphanumeric.", value)
		}
		path := field
		typ := "____"
		newValue := value
		funcName := "IsAlphanumeric"
		isValid := valid.IsAlphanumeric(value)

		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationProcessResponse(location, finalErrMessage, path, typ, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isASCII)

	return v.createProcessorFromValidator()
}

func (v validator) IsASCII(errorMessage string) processor {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	isASCII := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := v.location
		var finalErrMessage string
		if errorMessage != "" {
			finalErrMessage = errorMessage
		} else if v.errorMessage != "" {
			finalErrMessage = v.errorMessage
		} else {
			finalErrMessage = fmt.Sprintf("%s is not ascii.", value)
		}
		path := field
		typ := "____"
		newValue := value
		funcName := "IsASCII"
		isValid := valid.IsASCII(value)

		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationProcessResponse(location, finalErrMessage, path, typ, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isASCII)

	return v.createProcessorFromValidator()
}

func (v validator) IsBase32(errorMessage string, crockford bool) processor {
	// isBase32 := func(value, field string, ctx *gin.Context) validationProcessResponse {
	// 	location := v.location
	// 	var finalErrMessage string
	// 	if errorMessage != "" {
	// 		finalErrMessage = errorMessage
	// 	} else if v.errorMessage != "" {
	// 		finalErrMessage = v.errorMessage
	// 	} else {
	// 		finalErrMessage = fmt.Sprintf("%s is not base32.", value)
	// 	}
	// 	path := field
	// 	typ := "____"
	// 	newValue := value
	// 	funcName := "IsBase32"
	// 	// isValid := valid.IsBase32(value)

	// 	return newValidationProcessResponse(location, finalErrMessage, path, typ, newValue, funcName, isValid, false)
	// }

	// v.rules = append(v.rules, isBase32)

	return v.createProcessorFromValidator()
}

func (v validator) IsBase64(errorMessage string) processor {
	isBase64 := func(value, field string, ctx *gin.Context) validationProcessResponse {
		location := v.location
		var finalErrMessage string
		if errorMessage != "" {
			finalErrMessage = errorMessage
		} else if v.errorMessage != "" {
			finalErrMessage = v.errorMessage
		} else {
			finalErrMessage = fmt.Sprintf("%s is not base64.", value)
		}
		path := field
		typ := "____"
		newValue := value
		funcName := "IsBase64"
		isValid := valid.IsBase64(value)

		return newValidationProcessResponse(location, finalErrMessage, path, typ, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isBase64)

	return v.createProcessorFromValidator()
}
