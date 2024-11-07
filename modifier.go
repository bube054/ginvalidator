package ginvalidator

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	BailModifierFuncName string = "Bail"
	IfModifierFuncName   string = "If"
	NotModifierFuncName  string = "Not"
	SkipModifierFuncName string = "Skip"
)

// A modifier is simply a piece of the validation chain that can manipulate the whole validation chain.
type modifier struct {
	field      string             // the field to be specified
	errFmtFunc ErrFmtFuncHandler // the function to create the error message

	reqLoc            requestLocation  // the HTTP request location (e.g., body, headers, cookies, params, or queries)
	rulesCreatorFuncs ruleCreatorFuncs // the list of functions that creates the validation rules.
}

// recreateValidationChainFromModifier takes the previous modifier and returns a new validation chain.
func (m *modifier) recreateValidationChainFromModifier(ruleCreatorFunc ruleCreatorFunc) ValidationChain {
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

// Bail is a modifier that stops running the validation chain if any of the previous validators failed.
//
// This is useful to prevent a custom validator that touches a database or external API from running when you know it will fail.
//
// .Bail() can be used multiple times in the same validation chain if desired.
func (m modifier) Bail() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {

		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(sanitizedValue),
			withValidationChainName(BailModifierFuncName),
			withValidationChainType(modifierType),
			withShouldBail(false),
		)
	}

	return m.recreateValidationChainFromModifier(ruleCreator)
}

// IfModifierFunc defines a function that determines whether the validation chain should stop or continue.
// It returns `true` if the chain should stop, or `false` if it should continue.
//
// Parameters:
//   - req: the HTTP request context derived from `http.Request`.
//   - initialValue: the original value derived from the specified field.
//   - sanitizedValue: the current sanitized value after applying previous sanitizers.
type IfModifierFunc func(req http.Request, initialValue, sanitizedValue string) bool

// If adds a conditional check to decide whether the validation chain should continue for a field.
//
// The condition is evaluated by the provided [IfModifierFunc] and the result determines
// if the validation chain should bail out (`true`) or proceed (`false`).
//
// Parameters:
//   - imf: The [IfModifierFunc] used to evaluate the condition.
func (m modifier) If(imf IfModifierFunc) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		httpRequest := ctx.Request
		shouldBail := imf(*httpRequest, initialValue, sanitizedValue)

		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(sanitizedValue),
			withValidationChainName(IfModifierFuncName),
			withValidationChainType(modifierType),
			withShouldBail(shouldBail),
		)
	}

	return m.recreateValidationChainFromModifier(ruleCreator)
}

// Not negates the result of the next validator in the chain.
func (m modifier) Not() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(sanitizedValue),
			withValidationChainName(NotModifierFuncName),
			withValidationChainType(modifierType),
			withShouldBail(false),
		)
	}

	return m.recreateValidationChainFromModifier(ruleCreator)
}

// SkipModifierFunc defines a function that determines wwhether the next validator, modifier or sanitizer in validation chain should be skipped.
// It returns `true` if the next chain should skipped, or `false` if it should continue.
//
// Parameters:
//   - req: the HTTP request context derived from `http.Request`.
//   - initialValue: the original value derived from the specified field.
//   - sanitizedValue: the current sanitized value after applying previous sanitizers.
type SkipModifierFunc func(req http.Request, initialValue, sanitizedValue string) bool

// Skip adds a conditional check to decide whether the next validator, modifier or sanitizer in validation chain should be skipped.
//
// The condition is evaluated by the provided [SkipModifierFunc] and the result determines
// if the next link in validation chain should be skipped out (`true`) or proceed (`false`).
//
// Parameters:
//   - smf: The [SkipModifierFunc] used to evaluate the condition.
func (m modifier) Skip(smf SkipModifierFunc) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		httpRequest := ctx.Request
		shouldSkip := smf(*httpRequest, initialValue, sanitizedValue)

		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(sanitizedValue),
			withValidationChainName(SkipModifierFuncName),
			withValidationChainType(modifierType),
			withShouldBail(false),
			withShouldSkip(shouldSkip),
		)
	}

	return m.recreateValidationChainFromModifier(ruleCreator)
}

// newModifier creates and returns a new modifier.
//
// Parameters:
//   - field: The field to validate from the HTTP request data location (e.g., body, headers, cookies, params, or queries).
//   - errFmtFunc: A function that returns a custom error message. If nil, a generic error message will be used.
//   - reqLoc: The location in the HTTP request from where the field is extracted (e.g., body, headers, cookies, params, or queries).
func newModifier(field string, errFmtFunc ErrFmtFuncHandler, reqLoc requestLocation) modifier {
	return modifier{
		field:      field,
		errFmtFunc: errFmtFunc,
		reqLoc:     reqLoc,
	}
}
