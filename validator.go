package ginvalidator

import (
	"net/http"

	vgo "github.com/bube054/validatorgo"
	"github.com/gin-gonic/gin"
)

type validator struct {
	field      string
	errFmtFunc *ErrFmtFuncHandler

	reqLoc            requestLocation
	rulesCreatorFuncs ruleCreatorFuncs
}

func (v *validator) recreateVMSFromValidator(ruleCreatorFunc ruleCreatorFunc) VMS {
	newRulesCreatorFunc := append(v.rulesCreatorFuncs, ruleCreatorFunc)

	return VMS{
		validator: validator{
			field:             v.field,
			reqLoc:            v.reqLoc,
			errFmtFunc:        v.errFmtFunc,
			rulesCreatorFuncs: newRulesCreatorFunc,
		},
		modifier: modifier{
			field:             v.field,
			reqLoc:            v.reqLoc,
			errFmtFunc:        v.errFmtFunc,
			rulesCreatorFuncs: newRulesCreatorFunc,
		},
		sanitizer: sanitizer{
			field:             v.field,
			reqLoc:            v.reqLoc,
			errFmtFunc:        v.errFmtFunc,
			rulesCreatorFuncs: newRulesCreatorFunc,
		},
	}
}

// built in validators start here
type CustomValidatorFunc func(req http.Request, value string) bool

func (v validator) CustomValidator(cvf CustomValidatorFunc) VMS {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, value string) vmsRule {
		httpRequest := ctx.Request
		isValid := cvf(*httpRequest, value)

		return NewVMSRule(
			withIsValid(isValid),
			withNewValue(value),
			withVMSName("CustomValidator"),
			withTyp("validator"),
			withShouldBail(false),
			withShouldNegate(false),
			withShouldHide(false),
			withOptional(false),
		)
	}

	return v.recreateVMSFromValidator(ruleCreator)
}

// func (v validator) NotEmpty() VMS {
// 	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, value string) vmsRule {
// 		isValid := !vgo.IsEmpty(value, nil)

// 		return NewVMSRule(isValid, value, "NotEmpty", "validator", false, false)
// 	}

// 	return v.recreateVMSFromValidator(ruleCreator)
// }

// built in validators end here

// imported validators ends here
func (v validator) Contains(seed string, opts *vgo.ContainsOpt) VMS {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, value string) vmsRule {
		isValid := vgo.Contains(value, seed, opts)

		return NewVMSRule(
			withIsValid(isValid),
			withNewValue(value),
			withVMSName("Contains"),
			withTyp("validator"),
			withShouldBail(false),
			withShouldNegate(false),
			withShouldHide(false),
			withOptional(false),
		)
	}

	return v.recreateVMSFromValidator(ruleCreator)
}

func (v validator) Equals(comparison string) VMS {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, value string) vmsRule {
		isValid := vgo.Equals(value, comparison)

		return NewVMSRule(
			withIsValid(isValid),
			withNewValue(value),
			withVMSName("Equals"),
			withTyp("validator"),
			withShouldBail(false),
			withShouldNegate(false),
			withShouldHide(false),
			withOptional(false),
		)
	}

	return v.recreateVMSFromValidator(ruleCreator)
}

// imported validators ends here

func newValidator(field string, errFmtFunc *ErrFmtFuncHandler, reqLoc requestLocation) validator {
	return validator{
		field:      field,
		errFmtFunc: errFmtFunc,
		reqLoc:     reqLoc,
	}
}
