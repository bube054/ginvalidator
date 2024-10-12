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

func (s *sanitizer) recreateVMSFromSanitizer(ruleCreatorFunc ruleCreatorFunc) ValidationChain {
	newRulesCreatorFunc := append(s.rulesCreatorFuncs, ruleCreatorFunc)

	return ValidationChain{
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
type CustomSanitizerFunc func(req http.Request, initialValue, sanitizedValue string) string

func (s sanitizer) CustomSanitizer(csf CustomSanitizerFunc) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		httpRequest := ctx.Request
		newValue := csf(*httpRequest, initialValue, sanitizedValue)

		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(newValue),
			withValidationChainName("CustomSanitizer"),
			withValidationChainType(sanitizerType),
		)
	}

	return s.recreateVMSFromSanitizer(ruleCreator)
}

// built in sanitizers end here

func (s sanitizer) Blacklist(blacklistedChars string) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		newValue := san.Blacklist(sanitizedValue, blacklistedChars)

		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(newValue),
			withValidationChainName("Blacklist"),
			withValidationChainType(sanitizerType),
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
