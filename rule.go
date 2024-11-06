package ginvalidator

import "github.com/gin-gonic/gin"

// A validationChainRule is a rule that will control the the final flow in the validate function.
type validationChainRule struct {
	isValid             bool                // the values validity
	newValue            string              // the new sanitized value to be forwarded to the new rule creator func
	validationChainName string              // the name of the validator
	validationChainType validationChainType // the type of chain
	shouldBail          bool                // whether to  end the validation cahin
	shouldSkip          bool                // whether to skip a chain
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
