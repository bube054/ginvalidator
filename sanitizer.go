package ginvalidator

import (
	"net/http"

	san "github.com/bube054/validatorgo/sanitizer"
	"github.com/gin-gonic/gin"
)

type sanitizer struct {
	field      string
	errFmtFunc *ErrFmtFuncHandler

	reqLoc            requestLocation
	rulesCreatorFuncs ruleCreatorFuncs
}

func (s *sanitizer) recreateVMSFromSanitizer(ruleCreatorFunc ruleCreatorFunc) VMS {
	newRulesCreatorFunc := append(s.rulesCreatorFuncs, ruleCreatorFunc)

	return VMS{
		validator: validator{
			field:             s.field,
			reqLoc:            s.reqLoc,
			errFmtFunc:        s.errFmtFunc,
			rulesCreatorFuncs: newRulesCreatorFunc,
		},
		modifier: modifier{
			field:             s.field,
			reqLoc:            s.reqLoc,
			errFmtFunc:        s.errFmtFunc,
			rulesCreatorFuncs: newRulesCreatorFunc,
		},
		sanitizer: sanitizer{
			field:             s.field,
			reqLoc:            s.reqLoc,
			errFmtFunc:        s.errFmtFunc,
			rulesCreatorFuncs: newRulesCreatorFunc,
		},
	}
}

// built in sanitizers start here
type CustomSanitizerFunc func(req http.Request, value string) string

func (s sanitizer) CustomSanitizer(csf CustomSanitizerFunc) VMS {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, value string) vmsRule {
		httpRequest := ctx.Request
		newValue := csf(*httpRequest, value)

		return NewVMSRule(
			withIsValid(true),
			withNewValue(newValue),
			withVMSName("CustomSanitizer"),
			withTyp("sanitizer"),
			withShouldBail(false),
			withShouldNegate(false),
			withShouldHide(false),
			withOptional(false),
		)
	}

	return s.recreateVMSFromSanitizer(ruleCreator)
}

// built in sanitizers end here

func (s sanitizer) Blacklist(blacklistedChars string) VMS {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, value string) vmsRule {
		newValue := san.Blacklist(value, blacklistedChars)

		return NewVMSRule(
			withIsValid(true),
			withNewValue(newValue),
			withVMSName("Blacklist"),
			withTyp("sanitizer"),
			withShouldBail(false),
			withShouldNegate(false),
			withShouldHide(false),
			withOptional(false),
		)
	}

	return s.recreateVMSFromSanitizer(ruleCreator)
}

func newSanitizer(field string, errFmtFunc *ErrFmtFuncHandler, reqLoc requestLocation) sanitizer {
	return sanitizer{
		field:      field,
		errFmtFunc: errFmtFunc,
		reqLoc:     reqLoc,
	}
}
