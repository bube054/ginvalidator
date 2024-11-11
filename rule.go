package ginvalidator

import "github.com/gin-gonic/gin"

// validationChainRule represents a rule used in the validation chain, controlling the flow of validation.
type validationChainRule struct {
	isValid             bool                // Indicates whether the value passed the validation rule.
	newValue            string              // The sanitized value to pass to the next rule in the chain.
	validationChainName string              // The name of the validator being applied.
	validationChainType validationChainType // The type of chain (e.g., validator, sanitizer).
	shouldBail          bool                // Determines if validation should stop immediately on failure.
	shouldSkip          bool                // Determines if this chain rule should be skipped.
}

// NewValidationChainRule creates a new validationChainRule with the specified options.
func NewValidationChainRule(opts ...func(*validationChainRule)) validationChainRule {
	validationChainRule := &validationChainRule{}

	for _, opt := range opts {
		opt(validationChainRule)
	}

	return *validationChainRule
}

// withIsValid sets the isValid field, indicating if the value passed the validation rule.
func withIsValid(isValid bool) func(*validationChainRule) {
	return func(vcr *validationChainRule) {
		vcr.isValid = isValid
	}
}

// withNewValue sets the newValue field with the sanitized value to be forwarded.
func withNewValue(newValue string) func(*validationChainRule) {
	return func(vcr *validationChainRule) {
		vcr.newValue = newValue
	}
}

// withValidationChainName sets the validationChainName field to identify the validator.
func withValidationChainName(vcn string) func(*validationChainRule) {
	return func(vcr *validationChainRule) {
		vcr.validationChainName = vcn
	}
}

// withValidationChainType sets the validationChainType field, determining the type of chain.
func withValidationChainType(vct validationChainType) func(*validationChainRule) {
	return func(vcr *validationChainRule) {
		vcr.validationChainType = vct
	}
}

// withShouldBail sets the shouldBail field, specifying if validation should stop immediately on failure.
func withShouldBail(shouldBail bool) func(*validationChainRule) {
	return func(vcr *validationChainRule) {
		vcr.shouldBail = shouldBail
	}
}

// withShouldSkip sets the shouldSkip field, allowing this rule to be skipped.
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

// ruleCreatorFunc defines a function type that generates a validationChainRule.
// It takes a context, the initial value, and the current sanitized value, then returns a validationChainRule.
type ruleCreatorFunc func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule

// ruleCreatorFuncs is a slice of ruleCreatorFunc, allowing multiple rule functions
// to be applied sequentially in a validation chain.
type ruleCreatorFuncs []ruleCreatorFunc
