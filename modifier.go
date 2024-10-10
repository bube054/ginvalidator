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

func (m *modifier) recreateVMSFromModifier(ruleCreatorFunc ruleCreatorFunc) VMS {
	newRulesCreatorFunc := append(m.rulesCreatorFuncs, ruleCreatorFunc)

	return VMS{
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

type BailModifierOpts struct {
	level string // "chain" || "request"
}

func (m modifier) Bail(opts BailModifierOpts) VMS {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, value string) vmsRule {

		return NewVMSRule(
			withIsValid(true),
			withNewValue(value),
			withVMSName("Bail"),
			withTyp("modifier"),
			withShouldBail(false),
			withBailLevel(opts.level),
			withShouldNegate(false),
			withShouldHide(false),
			withOptional(false),
		)
	}

	return m.recreateVMSFromModifier(ruleCreator)
}

type IfModifierFunc func(req http.Request, value string) bool

func (m modifier) If(imf IfModifierFunc) VMS {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, value string) vmsRule {
		httpRequest := ctx.Request
		shouldBail := imf(*httpRequest, value)

		return NewVMSRule(
			withIsValid(true),
			withNewValue(value),
			withVMSName("If"),
			withTyp("modifier"),
			withShouldBail(shouldBail),
			withShouldNegate(false),
			withShouldHide(false),
			withOptional(false),
		)
	}

	return m.recreateVMSFromModifier(ruleCreator)
}

func (m modifier) Not() VMS {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, value string) vmsRule {
		return NewVMSRule(
			withIsValid(true),
			withNewValue(value),
			withVMSName("Not"),
			withTyp("modifier"),
			withShouldBail(false),
			withShouldNegate(true),
			withShouldHide(false),
			withOptional(false),
		)
	}

	return m.recreateVMSFromModifier(ruleCreator)
}

func (m modifier) Optional() VMS {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, value string) vmsRule {
		return NewVMSRule(
			withIsValid(true),
			withNewValue(value),
			withVMSName("Optional"),
			withTyp("modifier"),
			withShouldBail(false),
			withShouldNegate(false),
			withShouldHide(false),
			withOptional(true),
		)
	}

	return m.recreateVMSFromModifier(ruleCreator)
}

func (m modifier) Hide() VMS {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, value string) vmsRule {
		return NewVMSRule(
			withIsValid(true),
			withNewValue(value),
			withVMSName("Hide"),
			withTyp("modifier"),
			withShouldBail(false),
			withShouldNegate(false),
			withShouldHide(true),
			withOptional(false),
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
