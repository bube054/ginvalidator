package ginvalidator

import (
	"fmt"
	"net/http"
	"time"

	valid "github.com/asaskevich/govalidator"

	"github.com/gin-gonic/gin"
)

type validator struct {
	field           string
	errorMessage    string
	location        string
	rules           validationChainRules
	chainMethodType string
}

func (v *validator) createValidationChainFromValidator() validationChain {
	return validationChain{
		validator: *v,
		modifier: modifier{
			field:        v.field,
			errorMessage: v.errorMessage,
			location:     v.location,
			rules:        v.rules,
		},
		sanitizer: sanitizer{
			field:        v.field,
			errorMessage: v.errorMessage,
			location:     v.location,
			rules:        v.rules,
		},
	}
}

type customValidatorFunc func(value string, req http.Request, location string, path string) error

func (v validator) Custom(customValidator customValidatorFunc) validationChain {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	custom := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := v.location

		var finalErrMessage string
		customValidatorErr := customValidator(value, *ctx.Request, v.location, ctx.Request.URL.Path)

		if customValidatorErr != nil && customValidatorErr.Error() != "" {
			finalErrMessage = customValidatorErr.Error()
		} else if v.errorMessage != "" {
			finalErrMessage = v.errorMessage
		} else {
			finalErrMessage = fmt.Sprintf("%s is not alpha.", value)
		}

		path := field

		newValue := value
		funcName := "Custom"
		isValid := customValidatorErr == nil

		if previousRuleWasNegation {
			isValid = !isValid

			// if !isValid {
			// 	finalErrMessage = fmt.Sprintf("%s is not alpha.", value)
			// }
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, custom)

	return v.createValidationChainFromValidator()
}

// Min array length and max array length
type ArrayLengthCheckerOpts struct {
	Min uint
	Max uint
}

// Adds a validator to check that a value is an array.
// You will also be able to check in a future version that the array's length is greater than or equal to ArrayLengthCheckerOpts.min and/or that it's less than or equal to ArrayLengthCheckerOpts.max.
func (v validator) IsArray(errorMessage string, lengthChecker *ArrayLengthCheckerOpts) validationChain {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	isArray := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := v.location
		finalErrMessage := getFinalErrorMessage(errorMessage, v.errorMessage, fmt.Sprintf("%s %s", value, isArrayErrMsg))
		path := field
		newValue := value
		funcName := isArrayFuncName
		valueAsJSON := convertValueToJSON(value)
		valueJSONType := getJSONDataType(valueAsJSON)
		isValid := valueJSONType == "array"

		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isArray)

	return v.createValidationChainFromValidator()
}

// Adds a validator to check that a value is an object. For example, {}, { foo: 'bar' } would all pass this validator.
func (v validator) IsObject(errorMessage string) validationChain {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	isObject := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := v.location
		finalErrMessage := getFinalErrorMessage(errorMessage, v.errorMessage, fmt.Sprintf("%s %s", value, isObjectErrMsg))
		path := field
		newValue := value
		funcName := isObjectFuncName
		valueAsJSON := convertValueToJSON(value)
		valueJSONType := getJSONDataType(valueAsJSON)
		isValid := valueJSONType == "object"

		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isObject)

	return v.createValidationChainFromValidator()
}

// Adds a validator to check that a value is a string.
func (v validator) IsString(errorMessage string) validationChain {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	isString := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := v.location
		finalErrMessage := getFinalErrorMessage(errorMessage, v.errorMessage, fmt.Sprintf("%s %s", value, isStringErrMsg))
		path := field
		newValue := value
		funcName := isStringFuncName
		valueAsJSON := convertValueToJSON(value)
		valueJSONType := getJSONDataType(valueAsJSON)
		isValid := valueJSONType == "string"

		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isString)

	return v.createValidationChainFromValidator()
}

// Adds a validator to check that a value is a string that's not empty. This is analogous to .not().isEmpty().
func (v validator) IsNotEmpty(errorMessage string) validationChain {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	isNotEmpty := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := v.location
		finalErrMessage := getFinalErrorMessage(errorMessage, v.errorMessage, fmt.Sprintf("%s %s", value, isStringErrMsg))
		path := field
		newValue := value
		funcName := isStringFuncName
		isValid := value != ""

		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isNotEmpty)

	return v.createValidationChainFromValidator()
}

func (v validator) Contains(errorMessage string, substring string) validationChain {
	contains := func(value, field string, ctx *gin.Context) validationChainResponse {
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

		newValue := value
		funcName := "Contains"
		isValid := valid.Contains(value, substring)

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, contains)

	return v.createValidationChainFromValidator()
}

func (v validator) Equals(errorMessage string, comparison string) validationChain {
	equals := func(value, field string, ctx *gin.Context) validationChainResponse {
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

		newValue := value
		funcName := "Equals"
		isValid := value == comparison

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, equals)

	return v.createValidationChainFromValidator()
}

func (v validator) IsAfter(errorMessage string, comparisonTime time.Time) validationChain {
	isAfter := func(value, field string, ctx *gin.Context) validationChainResponse {
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

		newValue := value
		funcName := "IsAfter"

		isValid := false
		var valueTime interface{}
		valueTime = value

		valueAsTime, isTime := valueTime.(time.Time)

		if isTime {
			isValid = valueAsTime.After(comparisonTime)
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isAfter)

	return v.createValidationChainFromValidator()
}

func (v validator) IsAlpha(errorMessage string) validationChain {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	isASCII := func(value, field string, ctx *gin.Context) validationChainResponse {
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

		newValue := value
		funcName := "IsAlpha"
		isValid := valid.IsAlpha(value)

		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isASCII)

	return v.createValidationChainFromValidator()
}

func (v validator) IsAlphanumeric(errorMessage string) validationChain {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)
	fmt.Println("previousRuleWasNegationIsAlphaNumeric:", previousRuleWasNegation)

	isASCII := func(value, field string, ctx *gin.Context) validationChainResponse {
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

		newValue := value
		funcName := "IsAlphanumeric"
		isValid := valid.IsAlphanumeric(value)

		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isASCII)

	return v.createValidationChainFromValidator()
}

func (v validator) IsASCII(errorMessage string) validationChain {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	isASCII := func(value, field string, ctx *gin.Context) validationChainResponse {
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

		newValue := value
		funcName := "IsASCII"
		isValid := valid.IsASCII(value)

		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isASCII)

	return v.createValidationChainFromValidator()
}

func (v validator) IsBase32(errorMessage string, crockford bool) validationChain {
	// isBase32 := func(value, field string, ctx *gin.Context) validationChainResponse {
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
	//
	// 	newValue := value
	// 	funcName := "IsBase32"
	// 	// isValid := valid.IsBase32(value)

	// 	return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	// }

	// v.rules = append(v.rules, isBase32)

	return v.createValidationChainFromValidator()
}

func (v validator) IsBase64(errorMessage string) validationChain {
	isBase64 := func(value, field string, ctx *gin.Context) validationChainResponse {
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

		newValue := value
		funcName := "IsBase64"
		isValid := valid.IsBase64(value)

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isBase64)

	return v.createValidationChainFromValidator()
}
