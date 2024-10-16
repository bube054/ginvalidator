package ginvalidator

import (
	"fmt"
	"net/http"

	san "github.com/bube054/validatorgo/sanitizer"
	"github.com/gin-gonic/gin"
)

const (
	CustomSanitizerFuncName         string = "Bail"
	BlacklistSanitizerFuncName      string = "Blacklist"
	EscapeSanitizerFuncName         string = "Escape"
	LTrimSanitizerFuncName          string = "LTrim"
	NormalizeEmailSanitizerFuncName string = "NormalizeEmail"
	RTrimSanitizerFuncName          string = "RTrim"
	StripLowSanitizerFuncName       string = "StripLow"
	ToBooleanSanitizerFuncName      string = "ToBoolean"
	ToDateSanitizerFuncName         string = "ToDate"
	ToFloatSanitizerFuncName        string = "ToFloat"
	ToIntSanitizerFuncName          string = "ToInt"
	TrimSanitizerFuncName           string = "Trim"
	UnescapeSanitizerFuncName       string = "Unescape"
	WhitelistSanitizerFuncName      string = "Whitelist"
)

// A sanitizer is simply a piece of the validation chain that can sanitize values from the specified field.
type sanitizer struct {
	field      string             // the field to be specified
	errFmtFunc *ErrFmtFuncHandler // the function to create the error message

	reqLoc            requestLocation  // the HTTP request location (e.g., body, headers, cookies, params, or queries)
	rulesCreatorFuncs ruleCreatorFuncs // the list of functions that creates the validation rules.
}

// recreateValidationChainFromSanitizer takes the previous sanitizer and returns a new validation chain.
func (s *sanitizer) recreateValidationChainFromSanitizer(ruleCreatorFunc ruleCreatorFunc) ValidationChain {
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

// CustomSanitizerFunc defines a function that computes and returns the new sanitized value.
//
// Parameters:
//   - req: The HTTP request context derived from `http.Request`.
//   - initialValue: The original value derived from the specified field.
//   - sanitizedValue: The current sanitized value after applying previous sanitizers.
type CustomSanitizerFunc func(req http.Request, initialValue, sanitizedValue string) string

// CustomSanitizer applies a custom sanitizer function to compute the new sanitized value.
//
// Parameters:
//   - csf: The [CustomSanitizerFunc] used to compute the new sanitized value.
func (s sanitizer) CustomSanitizer(csf CustomSanitizerFunc) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		httpRequest := ctx.Request
		newValue := csf(*httpRequest, initialValue, sanitizedValue)

		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(newValue),
			withValidationChainName(CustomSanitizerFuncName),
			withValidationChainType(sanitizerType),
		)
	}

	return s.recreateValidationChainFromSanitizer(ruleCreator)
}

func (s sanitizer) Blacklist(blacklistedChars string) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		newValue := san.Blacklist(sanitizedValue, blacklistedChars)

		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(newValue),
			withValidationChainName(BlacklistSanitizerFuncName),
			withValidationChainType(sanitizerType),
		)
	}

	return s.recreateValidationChainFromSanitizer(ruleCreator)
}

func (s sanitizer) Escape() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		newValue := san.Escape(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(newValue),
			withValidationChainName(EscapeSanitizerFuncName),
			withValidationChainType(sanitizerType),
		)
	}

	return s.recreateValidationChainFromSanitizer(ruleCreator)
}

func (s sanitizer) LTrim(chars string) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		newValue := san.LTrim(sanitizedValue, chars)

		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(newValue),
			withValidationChainName(LTrimSanitizerFuncName),
			withValidationChainType(sanitizerType),
		)
	}

	return s.recreateValidationChainFromSanitizer(ruleCreator)
}

func (s sanitizer) NormalizeEmail(opts *san.NormalizeEmailOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		newValue := san.NormalizeEmail(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(newValue),
			withValidationChainName(NormalizeEmailSanitizerFuncName),
			withValidationChainType(sanitizerType),
		)
	}

	return s.recreateValidationChainFromSanitizer(ruleCreator)
}

func (s sanitizer) RTrim(chars string) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		newValue := san.RTrim(sanitizedValue, chars)

		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(newValue),
			withValidationChainName(RTrimSanitizerFuncName),
			withValidationChainType(sanitizerType),
		)
	}

	return s.recreateValidationChainFromSanitizer(ruleCreator)
}

func (s sanitizer) StripLow(keepNewLines bool) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		newValue := san.StripLow(sanitizedValue, keepNewLines)

		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(newValue),
			withValidationChainName(StripLowSanitizerFuncName),
			withValidationChainType(sanitizerType),
		)
	}

	return s.recreateValidationChainFromSanitizer(ruleCreator)
}

func (s sanitizer) ToBoolean(strict bool) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		ok := san.ToBoolean(sanitizedValue, strict)
		newValue := "false"

		if ok {
			newValue = "true"
		}

		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(newValue),
			withValidationChainName(ToBooleanSanitizerFuncName),
			withValidationChainType(sanitizerType),
		)
	}

	return s.recreateValidationChainFromSanitizer(ruleCreator)
}

func (s sanitizer) ToDate() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		time := san.ToDate(sanitizedValue)
		newValue := ""

		if time != nil {
			newValue = time.Format("2006-01-02 15:04:05")
		}

		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(newValue),
			withValidationChainName(ToDateSanitizerFuncName),
			withValidationChainType(sanitizerType),
		)
	}

	return s.recreateValidationChainFromSanitizer(ruleCreator)
}

func (s sanitizer) ToFloat() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		float, _ := san.ToFloat(sanitizedValue)
		newValue := fmt.Sprintf("%f", float)

		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(newValue),
			withValidationChainName(ToFloatSanitizerFuncName),
			withValidationChainType(sanitizerType),
		)
	}

	return s.recreateValidationChainFromSanitizer(ruleCreator)
}

func (s sanitizer) ToInt() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		num, _ := san.ToInt(sanitizedValue)
		newValue := fmt.Sprintf("%d", num)

		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(newValue),
			withValidationChainName(ToIntSanitizerFuncName),
			withValidationChainType(sanitizerType),
		)
	}

	return s.recreateValidationChainFromSanitizer(ruleCreator)
}

func (s sanitizer) Trim(chars string) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		newValue := san.Trim(sanitizedValue, chars)

		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(newValue),
			withValidationChainName(TrimSanitizerFuncName),
			withValidationChainType(sanitizerType),
		)
	}

	return s.recreateValidationChainFromSanitizer(ruleCreator)
}

func (s sanitizer) Unescape() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		newValue := san.Unescape(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(newValue),
			withValidationChainName(UnescapeSanitizerFuncName),
			withValidationChainType(sanitizerType),
		)
	}

	return s.recreateValidationChainFromSanitizer(ruleCreator)
}

func (s sanitizer) Whitelist(whitelistedChars string) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		newValue := san.Whitelist(sanitizedValue, whitelistedChars)

		return NewValidationChainRule(
			withIsValid(true),
			withNewValue(newValue),
			withValidationChainName(WhitelistSanitizerFuncName),
			withValidationChainType(sanitizerType),
		)
	}

	return s.recreateValidationChainFromSanitizer(ruleCreator)
}

func newSanitizer(field string, errFmtFunc *ErrFmtFuncHandler, reqLoc requestLocation) sanitizer {
	return sanitizer{
		field:      field,
		errFmtFunc: errFmtFunc,
		reqLoc:     reqLoc,
	}
}
