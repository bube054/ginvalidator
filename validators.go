package ginvalidator

import (
	"fmt"

	valid "github.com/asaskevich/govalidator"
)

type validator struct {
	field        string
	errorMessage string
	location     string
	rules        validationProcessesRules
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

func (v validator) IsAlpha(errorMessage string) processor {
	isASCII := func(value, field string) validationProcessResponse {
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

		return newValidationProcessResponse(location, finalErrMessage, path, typ, newValue, funcName, isValid)
	}

	v.rules = append(v.rules, isASCII)

	return v.createProcessorFromValidator()
}

func (v validator) IsAlphanumeric(errorMessage string) processor {
	isASCII := func(value, field string) validationProcessResponse {
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

		return newValidationProcessResponse(location, finalErrMessage, path, typ, newValue, funcName, isValid)
	}

	v.rules = append(v.rules, isASCII)

	return v.createProcessorFromValidator()
}

func (v validator) IsASCII(errorMessage string) processor {
	isASCII := func(value, field string) validationProcessResponse {
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

		return newValidationProcessResponse(location, finalErrMessage, path, typ, newValue, funcName, isValid)
	}

	v.rules = append(v.rules, isASCII)

	return v.createProcessorFromValidator()
}
