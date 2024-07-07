package ginvalidator

import (
	"fmt"
	"net/http"
	"strings"
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

type CustomValidatorFunc func(value string, req http.Request, location string) error

// Adds a custom validator function to the chain.
// The field value will be valid if:
// The custom validator returns nil error;
// If the custom validator returns an error.
// A common use case for .custom() is to verify that an e-mail address doesn't already exists. If it does, return an error:
func (v validator) Custom(customValidator CustomValidatorFunc) validationChain {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	custom := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := v.location

		var finalErrMessage string
		customValidatorErr := customValidator(value, *ctx.Request, v.location)

		if customValidatorErr != nil && customValidatorErr.Error() != "" {
			finalErrMessage = customValidatorErr.Error()
		} else if v.errorMessage != "" {
			finalErrMessage = v.errorMessage
		} else {
			finalErrMessage = fmt.Sprintf("%s is not alpha.", value)
		}

		path := field

		newValue := value
		funcName := customFuncName
		isValid := customValidatorErr == nil

		if previousRuleWasNegation {
			isValid = !isValid
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
		finalErrMessage := getFinalErrorMessage(errorMessage, v.errorMessage, isArrayErrMsg)
		path := field
		newValue := value
		funcName := isArrayFuncName
		valueAsJSON := convertValueToJSON(value)
		valueJSONType := getJSONDataType(valueAsJSON)
		fmt.Println("valueJSONType:", valueJSONType)
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
		finalErrMessage := getFinalErrorMessage(errorMessage, v.errorMessage, isObjectErrMsg)
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
		finalErrMessage := getFinalErrorMessage(errorMessage, v.errorMessage, isStringErrMsg)
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
		finalErrMessage := getFinalErrorMessage(errorMessage, v.errorMessage, isNotEmptyErrMsg)
		path := field
		newValue := value
		funcName := isNotEmptyFuncName
		isValid := value != ""

		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isNotEmpty)

	return v.createValidationChainFromValidator()
}

// Takes the usual error message and checks if the value contains the substring.
func (v validator) Contains(errorMessage string, substring string) validationChain {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	contains := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := v.location
		var finalErrMessage string
		if errorMessage != "" {
			finalErrMessage = errorMessage
		} else if v.errorMessage != "" {
			finalErrMessage = v.errorMessage
		} else {
			finalErrMessage = containsErrMsg
		}
		path := field

		newValue := value
		funcName := containsFuncName
		isValid := valid.Contains(value, substring)

		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, contains)

	return v.createValidationChainFromValidator()
}

// Takes the usual error message and checks if the value is same as the comparison.
func (v validator) Equals(errorMessage string, comparison string) validationChain {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	equals := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := v.location
		var finalErrMessage string
		if errorMessage != "" {
			finalErrMessage = errorMessage
		} else if v.errorMessage != "" {
			finalErrMessage = v.errorMessage
		} else {
			finalErrMessage = equalsErrMsg
		}
		path := field

		newValue := value
		funcName := equalsFuncName
		isValid := value == comparison

		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, equals)

	return v.createValidationChainFromValidator()
}

// Takes the usual error message and checks if the value is after the comparison time.
func (v validator) IsAfter(errorMessage string, comparisonTime time.Time) validationChain {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	isAfter := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := v.location
		var finalErrMessage string
		if errorMessage != "" {
			finalErrMessage = errorMessage
		} else if v.errorMessage != "" {
			finalErrMessage = v.errorMessage
		} else {
			finalErrMessage = isAfterErrMsg
		}
		path := field

		newValue := value
		funcName := isAfterFuncName

		isValid := false
		var valueTime interface{} = value

		valueAsTime, isTime := valueTime.(time.Time)

		if isTime {
			isValid = valueAsTime.After(comparisonTime)
		}

		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isAfter)

	return v.createValidationChainFromValidator()
}

// Takes the usual error message and checks if the string contains only letters (a-zA-Z). Empty string is valid.
func (v validator) IsAlpha(errorMessage string) validationChain {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	isAlpha := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := v.location
		var finalErrMessage string
		if errorMessage != "" {
			finalErrMessage = errorMessage
		} else if v.errorMessage != "" {
			finalErrMessage = v.errorMessage
		} else {
			finalErrMessage = isAlphaErrMsg
		}
		path := field

		newValue := value
		funcName := isAlphaFuncName
		isValid := valid.IsAlpha(value)

		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isAlpha)

	return v.createValidationChainFromValidator()
}

// Takes the usual error message and checks if the string contains only letters and numbers. Empty string is valid.
func (v validator) IsAlphanumeric(errorMessage string) validationChain {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	isAlphaNumeric := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := v.location
		var finalErrMessage string
		if errorMessage != "" {
			finalErrMessage = errorMessage
		} else if v.errorMessage != "" {
			finalErrMessage = v.errorMessage
		} else {
			finalErrMessage = isAlphanumericErrMsg
		}
		path := field

		newValue := value
		funcName := isAlphaFuncName
		isValid := valid.IsAlphanumeric(value)

		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isAlphaNumeric)

	return v.createValidationChainFromValidator()
}

// Takes the usual error message and checks if the string contains ASCII chars only. Empty string is valid.
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
			finalErrMessage = isASCIIErrMsg
		}
		path := field

		newValue := value
		funcName := isASCIIFuncName
		isValid := valid.IsASCII(value)
		fmt.Println("test", previousRuleWasNegation)

		if previousRuleWasNegation {
			fmt.Println("isValid", isValid)
			isValid = !isValid
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isASCII)

	return v.createValidationChainFromValidator()
}

func (v validator) IsBase32(errorMessage string, crockford bool) validationChain {
	// previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

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

// Takes the usual error message and checks if a string is base64 encoded.
func (v validator) IsBase64(errorMessage string) validationChain {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	isBase64 := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := v.location
		var finalErrMessage string
		if errorMessage != "" {
			finalErrMessage = errorMessage
		} else if v.errorMessage != "" {
			finalErrMessage = v.errorMessage
		} else {
			finalErrMessage = isBase64ErrMsg
		}
		path := field

		newValue := value
		funcName := isBase64FuncName
		isValid := valid.IsBase64(value)

		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isBase64)

	return v.createValidationChainFromValidator()
}

// Takes the usual error message and checks if the value is before the comparison time.
func (v validator) IsBefore(errorMessage string, comparisonTime time.Time) validationChain {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	isBefore := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := v.location
		var finalErrMessage string
		if errorMessage != "" {
			finalErrMessage = errorMessage
		} else if v.errorMessage != "" {
			finalErrMessage = v.errorMessage
		} else {
			finalErrMessage = isBeforeErrMsg
		}
		path := field

		newValue := value
		funcName := isBeforeFuncName

		isValid := false
		var valueTime interface{} = value

		valueAsTime, isTime := valueTime.(time.Time)

		if isTime {
			isValid = valueAsTime.Before(comparisonTime)
		}

		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isBefore)

	return v.createValidationChainFromValidator()
}

func (v validator) IsBIC(errorMessage string) validationChain {
	// previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	// isBase64 := func(value, field string, ctx *gin.Context) validationChainResponse {
	// 	location := v.location
	// 	var finalErrMessage string
	// 	if errorMessage != "" {
	// 		finalErrMessage = errorMessage
	// 	} else if v.errorMessage != "" {
	// 		finalErrMessage = v.errorMessage
	// 	} else {
	// 		finalErrMessage = isBase64ErrMsg
	// 	}
	// 	path := field

	// 	newValue := value
	// 	funcName := isBase64FuncName
	// 	isValid := valid.IsBIC(value)
	// 	if previousRuleWasNegation {
	// 	isValid = !isValid
	// }

	// 	return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	// }

	// v.rules = append(v.rules, isBase64)

	return v.createValidationChainFromValidator()
}

// Takes the usual error message and checks if the value is a boolean.
func (v validator) IsBOOLEAN(errorMessage string, strictNess bool) validationChain {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	isBoolean := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := v.location
		finalErrMessage := getFinalErrorMessage(errorMessage, v.errorMessage, isBeforeErrMsg)
		path := field
		newValue := value
		funcName := isBeforeFuncName
		var isValid bool

		if strictNess {
			isValid = value == "true" || value == "false"
		} else {
			isValid = strings.ToLower(value) == "true" || strings.ToLower(value) == "false"
		}

		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isBoolean)

	return v.createValidationChainFromValidator()
}

// Takes the usual error message and checks if the value is a boolean.
func (v validator) IsBtcAddress(errorMessage string) validationChain {
	// previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	// isBtcAddress := func(value, field string, ctx *gin.Context) validationChainResponse {
	// 	location := v.location
	// 	finalErrMessage := getFinalErrorMessage(errorMessage, v.errorMessage, isBtcAddressErrMsg)
	// 	path := field
	// 	newValue := value
	// 	funcName := isBtcAddressFuncName
	// 	isValid := valid.isB(value)

	// 	return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	// }

	// v.rules = append(v.rules, isBtcAddress)

	return v.createValidationChainFromValidator()
}

// Takes the usual error message and checks if the string's length (in bytes) falls in a range.
func (v validator) IsByteLength(errorMessage string, min int, max int) validationChain {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	isByteLength := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := v.location
		finalErrMessage := getFinalErrorMessage(errorMessage, v.errorMessage, isByteLengthErrMsg)
		path := field
		newValue := value
		funcName := isByteLengthFuncName
		isValid := valid.IsByteLength(value, min, max)

		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isByteLength)

	return v.createValidationChainFromValidator()
}

// Takes the usual error message and checks if the string is a credit card.
func (v validator) IsCreditCard(errorMessage string) validationChain {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	isCreditCard := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := v.location
		finalErrMessage := getFinalErrorMessage(errorMessage, v.errorMessage, isCreditCardErrMsg)
		path := field
		newValue := value
		funcName := isCreditCardFuncName
		isValid := valid.IsCreditCard(newValue)

		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isCreditCard)

	return v.createValidationChainFromValidator()
}

// Takes the usual error message and checks if the string is a credit card.
func (v validator) IsCurrency(errorMessage string, card string) validationChain {
	// previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	// isCurrency := func(value, field string, ctx *gin.Context) validationChainResponse {
	// 	location := v.location
	// 	finalErrMessage := getFinalErrorMessage(errorMessage, v.errorMessage, isCurrencyErrMsg)
	// 	path := field
	// 	newValue := value
	// 	funcName := isCurrencyFuncName
	// 	isValid := valid.isCur

	// 	if previousRuleWasNegation {
	// 		isValid = !isValid
	// 	}

	// 	return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	// }

	// v.rules = append(v.rules, isCurrency)

	return v.createValidationChainFromValidator()
}

// Takes the usual error message and checks if a string is base64 encoded data URI such as an image.
func (v validator) IsDataURI(errorMessage string) validationChain {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	isDataURI := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := v.location
		finalErrMessage := getFinalErrorMessage(errorMessage, v.errorMessage, isDataURIErrMsg)
		path := field
		newValue := value
		funcName := isDataURIFuncName
		isValid := valid.IsDataURI(newValue)

		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isDataURI)

	return v.createValidationChainFromValidator()
}

// Takes the usual error message and checks if a string is base64 encoded data URI such as an image.
func (v validator) IsDate(errorMessage string) validationChain {
	previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	isDate := func(value, field string, ctx *gin.Context) validationChainResponse {
		location := v.location
		finalErrMessage := getFinalErrorMessage(errorMessage, v.errorMessage, isDateErrMsg)
		path := field
		newValue := value
		funcName := isDateFuncName
		_, err := time.Parse("2006-01-02", newValue)
		isValid := err == nil
		if previousRuleWasNegation {
			isValid = !isValid
		}

		return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	}

	v.rules = append(v.rules, isDate)

	return v.createValidationChainFromValidator()
}

// Takes the usual error message and checks if a string is base64 encoded data URI such as an image.
func (v validator) IsDecimal(errorMessage string) validationChain {
	// previousRuleWasNegation := wasPreviousRuleNegation(v.rules)

	// isDecimal := func(value, field string, ctx *gin.Context) validationChainResponse {
	// 	location := v.location
	// 	finalErrMessage := getFinalErrorMessage(errorMessage, v.errorMessage, isDecimalErrMsg)
	// 	path := field
	// 	newValue := value
	// 	funcName := isDecimalFuncName
	// 	isValid := valid.IsDec
	// 	if previousRuleWasNegation {
	// 		isValid = !isValid
	// 	}

	// 	return newValidationChainResponse(location, finalErrMessage, path, newValue, funcName, isValid, false)
	// }

	// v.rules = append(v.rules, isDecimal)

	return v.createValidationChainFromValidator()
}
