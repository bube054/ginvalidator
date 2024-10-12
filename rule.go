package ginvalidator

import "github.com/gin-gonic/gin"

type validationChainRule struct {
	isValid             bool
	newValue            string
	validationChainName string
	validationChainType validationChainType
	shouldBail          bool
	shouldSkip          bool
}

func NewValidationChainRule(opts ...func(*validationChainRule)) validationChainRule {
	validationChainRule := &validationChainRule{}

	for _, opt := range opts {
		opt(validationChainRule)
	}

	return *validationChainRule
}

func withIsValid(isValid bool) func(*validationChainRule) {
	return func(vcr *validationChainRule) {
		vcr.isValid = isValid
	}
}

func withNewValue(newValue string) func(*validationChainRule) {
	return func(vcr *validationChainRule) {
		vcr.newValue = newValue
	}
}

func withValidationChainName(vcn string) func(*validationChainRule) {
	return func(vcr *validationChainRule) {
		vcr.validationChainName = vcn
	}
}

func withValidationChainType(vct validationChainType) func(*validationChainRule) {
	return func(vcr *validationChainRule) {
		vcr.validationChainType = vct
	}
}

func withShouldBail(shouldBail bool) func(*validationChainRule) {
	return func(vcr *validationChainRule) {
		vcr.shouldBail = shouldBail
	}
}

func withShouldSkip(shouldSkip bool) func(*validationChainRule) {
	return func(vcr *validationChainRule) {
		vcr.shouldSkip = shouldSkip
	}
}

// func NewValidationChainRule(isValid bool, newValue string, validationChainName string, validationChainType string, shouldBail bool, shouldNegate bool) validationChainRule {
// 	return validationChainRule{
// 		isValid:      isValid,
// 		newValue:     newValue,
// 		validationChainName:      validationChainName,
// 		validationChainType:          validationChainType,
// 		shouldBail:   shouldBail,
// 		shouldNegate: shouldNegate,
// 	}
// }

type ruleCreatorFunc func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule

type ruleCreatorFuncs []ruleCreatorFunc
