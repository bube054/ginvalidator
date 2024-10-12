package ginvalidator

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type modifier struct {
	field      string
	errFmtFunc *ErrFmtFuncHandler

	reqLoc            requestLocation
	rulesCreatorFuncs ruleCreatorFuncs
}

func (m *modifier) recreateVMSFromModifier(ruleCreatorFunc ruleCreatorFunc) ValidationChain {
	newRulesCreatorFunc := append(m.rulesCreatorFuncs, ruleCreatorFunc)

	return ValidationChain{
		validator: validator{
			field:             m.field,
			reqLoc:            m.reqLoc,
			errFmtFunc:        m.errFmtFunc,
			rulesCreatorFuncs: newRulesCreatorFunc,
		},
		modifier: modifier{
			field:             m.field,
			reqLoc:            m.reqLoc,
			errFmtFunc:        m.errFmtFunc,
			rulesCreatorFuncs: newRulesCreatorFunc,
		},
		sanitizer: sanitizer{
			field:             m.field,
			reqLoc:            m.reqLoc,
			errFmtFunc:        m.errFmtFunc,
			rulesCreatorFuncs: newRulesCreatorFunc,
		},
	}
}

func (m modifier) Bail() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {

		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(sanitizedValue),
			withValidationChainName("Bail"),
			withValidationChainType(modifierType),
			withShouldBail(false),
		)
	}

	return m.recreateVMSFromModifier(ruleCreator)
}

type IfModifierFunc func(req http.Request, initialValue, sanitizedValue string) bool

func (m modifier) If(imf IfModifierFunc) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		httpRequest := ctx.Request
		shouldBail := imf(*httpRequest, initialValue, sanitizedValue)

		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(sanitizedValue),
			withValidationChainName("If"),
			withValidationChainType(modifierType),
			withShouldBail(shouldBail),
		)
	}

	return m.recreateVMSFromModifier(ruleCreator)
}

func (m modifier) Not() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(sanitizedValue),
			withValidationChainName("Not"),
			withValidationChainType(modifierType),
			withShouldBail(false),
		)
	}

	return m.recreateVMSFromModifier(ruleCreator)
}

type SkipModifierFunc func(req http.Request, initialValue, sanitizedValue string) bool

func (m modifier) Skip(smf SkipModifierFunc) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		httpRequest := ctx.Request
		shouldSkip := smf(*httpRequest, initialValue, sanitizedValue)

		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(sanitizedValue),
			withValidationChainName("If"),
			withValidationChainType(modifierType),
			withShouldBail(false),
			withShouldSkip(shouldSkip),
		)
	}

	return m.recreateVMSFromModifier(ruleCreator)
}

func newModifier(field string, errFmtFunc *ErrFmtFuncHandler, reqLoc requestLocation) modifier {
	return modifier{
		field:      field,
		errFmtFunc: errFmtFunc,
		reqLoc:     reqLoc,
	}
}
