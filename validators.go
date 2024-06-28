package ginvalidator

import (
	valid "github.com/asaskevich/govalidator"
)

type Validator struct {
	field        string
	errorMessage string
	location     string

	value        string
	currentValue string
	resErrs      []ResponseError
}

func (v Validator) IsALPHA(errMsg string) Processor {
	var finalErrMsg = errMsg
	var resErrs = v.resErrs
	var resErr ResponseError

	if errMsg == "" && v.errorMessage != "" {
		finalErrMsg = v.errorMessage
	}

	if !valid.IsAlpha(v.currentValue) {
		resErr = ResponseError{location: v.location, msg: finalErrMsg, path: v.field, typ: "field"}
		resErrs = append(resErrs, resErr)
	}

	return Processor{
		Validator: Validator{field: v.field, errorMessage: v.errorMessage, location: v.location, value: v.value, currentValue: v.currentValue, resErrs: resErrs},
		Modifier:  Modifier{field: v.field, errorMessage: v.errorMessage, location: v.location, value: v.value, currentValue: v.currentValue, resErrs: resErrs},
		Sanitizer: Sanitizer{field: v.field, errorMessage: v.errorMessage, location: v.location, value: v.value, currentValue: v.currentValue, resErrs: resErrs},
	}
}

func (v Validator) IsASCII(errMsg string) Processor {
	var finalErrMsg = errMsg
	var resErrs = v.resErrs
	var resErr ResponseError

	if errMsg == "" && v.errorMessage != "" {
		finalErrMsg = v.errorMessage
	}

	if !valid.IsASCII(v.currentValue) {
		resErr = ResponseError{location: v.location, msg: finalErrMsg, path: v.field, typ: "field"}
		resErrs = append(resErrs, resErr)
	}

	return Processor{
		Validator: Validator{field: v.field, errorMessage: v.errorMessage, location: v.location, value: v.value, currentValue: v.currentValue, resErrs: resErrs},
		Modifier:  Modifier{field: v.field, errorMessage: v.errorMessage, location: v.location, value: v.value, currentValue: v.currentValue, resErrs: resErrs},
		Sanitizer: Sanitizer{field: v.field, errorMessage: v.errorMessage, location: v.location, value: v.value, currentValue: v.currentValue, resErrs: resErrs},
	}
}

func (v Validator) IsAlphanumeric(errMsg string) Processor {

	var finalErrMsg = errMsg
	var resErrs = v.resErrs
	var resErr ResponseError

	if errMsg == "" && v.errorMessage != "" {
		finalErrMsg = v.errorMessage
	}

	if !valid.IsAlphanumeric(v.currentValue) {
		resErr = ResponseError{location: v.location, msg: finalErrMsg, path: v.field, typ: "field"}
		resErrs = append(resErrs, resErr)
	}

	return Processor{
		Validator: Validator{field: v.field, errorMessage: v.errorMessage, location: v.location, value: v.value, currentValue: v.currentValue, resErrs: resErrs},
		Modifier:  Modifier{field: v.field, errorMessage: v.errorMessage, location: v.location, value: v.value, currentValue: v.currentValue, resErrs: resErrs},
		Sanitizer: Sanitizer{field: v.field, errorMessage: v.errorMessage, location: v.location, value: v.value, currentValue: v.currentValue, resErrs: resErrs},
	}
}
