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

// Blacklist is a sanitizer that remove characters that appear in the blacklist.
//
// This function uses the [validatorgo] package to perform the sanitization logic.
//
// Its parameters are according to [Blacklist].
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [Blacklist]: https://pkg.go.dev/github.com/bube054/validatorgo/sanitizer#Blacklist
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

// Escape is a sanitizer that replaces <, >, &, ' and ". with HTML entities.
//
// This function uses the [validatorgo] package to perform the sanitization logic.
//
// Its parameters are according to [Escape].
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [Escape]: https://pkg.go.dev/github.com/bube054/validatorgo/sanitizer#Escape
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

// LTrim is a sanitizer that trims characters (whitespace by default) from the left-side of the input.
//
// This function uses the [validatorgo] package to perform the sanitization logic.
//
// Its parameters are according to [LTrim].
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [LTrim]: https://pkg.go.dev/github.com/bube054/validatorgo/sanitizer#LTrim
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

// NormalizeEmail is a sanitizer that canonicalizes an email address. (This doesn't validate that the input is an email, if you want to validate the email use IsEmail beforehand).
//
// This function uses the [validatorgo] package to perform the sanitization logic.
//
// Its parameters are according to [NormalizeEmail].
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [NormalizeEmail]: https://pkg.go.dev/github.com/bube054/validatorgo/sanitizer#NormalizeEmail
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

// RTrim is a sanitizer that trims characters (whitespace by default) from the right-side of the input.
//
// This function uses the [validatorgo] package to perform the sanitization logic.
//
// Its parameters are according to [RTrim].
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [RTrim]: https://pkg.go.dev/github.com/bube054/validatorgo/sanitizer#RTrim
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

// StripLow is a sanitizer that removes characters with a numerical value < 32 and 127, mostly control characters.
//
// This function uses the [validatorgo] package to perform the sanitization logic.
//
// Its parameters are according to [StripLow].
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [StripLow]: https://pkg.go.dev/github.com/bube054/validatorgo/sanitizer#StripLow
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

// ToBoolean is a A sanitizer that converts the input string to a boolean as s string "true" or "false"
//
// This function uses the [validatorgo] package to perform the sanitization logic.
//
// Its parameters are according to [ToBoolean].
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [ToBoolean]: https://pkg.go.dev/github.com/bube054/validatorgo/sanitizer#ToBoolean
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

// ToDate is a sanitizer that converts the value too a textual representation.
//
// This function uses the [validatorgo] package to perform the sanitization logic.
//
// Its parameters are according to [ToDate].
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [ToDate]: https://pkg.go.dev/github.com/bube054/validatorgo/sanitizer#ToDate
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

// ToFloat is a sanitizer that converts the input string to a float64.
//
// This function uses the [validatorgo] package to perform the sanitization logic.
//
// Its parameters are according to [ToFloat].
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [ToFloat]: https://pkg.go.dev/github.com/bube054/validatorgo/sanitizer#ToFloat
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

// ToInt is a sanitizer that converts the input string to an int and also returns an error if the input is not a int. (Beware of octals)
//
// This function uses the [validatorgo] package to perform the sanitization logic.
//
// Its parameters are according to [ToInt].
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [ToInt]: https://pkg.go.dev/github.com/bube054/validatorgo/sanitizer#ToInt
func (s sanitizer) ToInt() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		num, err := san.ToInt(sanitizedValue)

		fmt.Println(num, err)

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

// Trim is a sanitizer that trim characters (whitespace by default) from both sides of the input.
//
// This function uses the [validatorgo] package to perform the sanitization logic.
//
// Its parameters are according to [Trim].
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [Trim]: https://pkg.go.dev/github.com/bube054/validatorgo/sanitizer#Trim
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

// Unescape is a A sanitizer that replaces HTML encoded entities with <, >, &, ', ", `, \ and /.
//
// This function uses the [validatorgo] package to perform the sanitization logic.
//
// Its parameters are according to [Unescape].
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [Unescape]: https://pkg.go.dev/github.com/bube054/validatorgo/sanitizer#Unescape
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

// Whitelist is a sanitizer that removes characters that do not appear in the whitelist.
//
// This function uses the [validatorgo] package to perform the sanitization logic.
//
// Its parameters are according to [Whitelist].
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [Whitelist]: https://pkg.go.dev/github.com/bube054/validatorgo/sanitizer#Whitelist
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

// newSanitizer creates and returns a new sanitizer.
//
// Parameters:
//   - field: The field to validate from the HTTP request data location (e.g., body, headers, cookies, params, or queries).
//   - errFmtFunc: A function that returns a custom error message. If nil, a generic error message will be used.
//   - reqLoc: The location in the HTTP request from where the field is extracted (e.g., body, headers, cookies, params, or queries).
func newSanitizer(field string, errFmtFunc *ErrFmtFuncHandler, reqLoc requestLocation) sanitizer {
	return sanitizer{
		field:      field,
		errFmtFunc: errFmtFunc,
		reqLoc:     reqLoc,
	}
}
