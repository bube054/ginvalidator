package ginvalidator

import "github.com/gin-gonic/gin"

type vmsRule struct {
	isValid      bool
	newValue     string
	vmsName      string
	typ          string
	shouldBail   bool
	bailLevel    string // "chain" || "request"
	shouldNegate bool
	shouldHide   bool
	optional     bool
}

func NewVMSRule(opts ...func(*vmsRule)) vmsRule {
	vmsRule := &vmsRule{}

	for _, opt := range opts {
		opt(vmsRule)
	}

	return *vmsRule
}

func withIsValid(isValid bool) func(*vmsRule) {
	return func(vr *vmsRule) {
		vr.isValid = isValid
	}
}

func withNewValue(newValue string) func(*vmsRule) {
	return func(vr *vmsRule) {
		vr.newValue = newValue
	}
}

func withVMSName(vmsName string) func(*vmsRule) {
	return func(vr *vmsRule) {
		vr.vmsName = vmsName
	}
}

func withTyp(typ string) func(*vmsRule) {
	return func(vr *vmsRule) {
		vr.typ = typ
	}
}

func withShouldBail(shouldBail bool) func(*vmsRule) {
	return func(vr *vmsRule) {
		vr.shouldBail = shouldBail
	}
}

func withBailLevel(bailLevel string) func(*vmsRule) {
	return func(vr *vmsRule) {
		switch bailLevel {
		case "chain", "request":
			vr.bailLevel = bailLevel
		default:
			vr.bailLevel = "chain"
		}
	}
}

func withShouldNegate(shouldNegate bool) func(*vmsRule) {
	return func(vr *vmsRule) {
		vr.shouldNegate = shouldNegate
	}
}

func withShouldHide(shouldHide bool) func(*vmsRule) {
	return func(vr *vmsRule) {
		vr.shouldHide = shouldHide
	}
}

func withOptional(optional bool) func(*vmsRule) {
	return func(vr *vmsRule) {
		vr.optional = optional
	}
}

// func NewVMSRule(isValid bool, newValue string, vmsName string, typ string, shouldBail bool, shouldNegate bool) vmsRule {
// 	return vmsRule{
// 		isValid:      isValid,
// 		newValue:     newValue,
// 		vmsName:      vmsName,
// 		typ:          typ,
// 		shouldBail:   shouldBail,
// 		shouldNegate: shouldNegate,
// 	}
// }

type vmsRules []vmsRule

type ruleCreatorFunc func(ctx *gin.Context, value string) vmsRule

type ruleCreatorFuncs []ruleCreatorFunc
