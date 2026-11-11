// Package ginvalidator is a set of [Gin] middlewares that wraps the
// extensive collection of validators and sanitizers offered by [validatorgo].
//
// It allows you to combine them in many ways so that you can validate and sanitize your Gin requests,
// and offers tools to determine if the request is valid or not, which data was matched according to your validators, and so on.
//
// It is based on the popular js/express library [express-validator]
//
// [Gin]: https://github.com/gin-gonic/gin
// [validatorgo]: https://github.com/bube054/validatorgo
// [express-validator]: https://github.com/express-validator/express-validator
package ginvalidator

import (
	"errors"
	"sync/atomic"

	vgo "github.com/bube054/validatorgo"
	"github.com/gin-gonic/gin"
)

var globalErrorOrder uint64

const DefaultErrMsg string = "Invalid value"

// InvalidJSONErrMsg and InvalidJSONErrCode are the message and code of the error a
// body chain records when the request body is not valid JSON.
const (
	InvalidJSONErrMsg  string = "Invalid JSON in request body"
	InvalidJSONErrCode string = "invalid_json"
)

// UnreadableRequestErrMsg and UnreadableRequestErrCode are the message and code of
// the error a chain records when its value could not be read from the request at
// all, for example because the client disconnected while sending the body.
const (
	UnreadableRequestErrMsg  string = "Could not read the request"
	UnreadableRequestErrCode string = "unreadable_request"
)

// requestErrorResult is the result of a chain whose value could not be read: one
// error, no rules run, and nothing for the matched data.
func requestErrorResult(location, field, msg, code string) chainResult {
	vce := newValidationChainError(
		vceWithLocation(location),
		vceWithMessage(msg),
		vceWithField(field),
		vceWithCode(code),
		vceWithOrder(atomic.AddUint64(&globalErrorOrder, 1)),
	)
	return chainResult{
		errors:   []ValidationChainError{vce},
		location: location,
		field:    field,
	}
}

// DefaultErrFmtFunc is a package-level fallback error message formatter.
// When set, it is used for any validation chain that does not have its own errFmtFunc.
// If nil (the default), the system falls back to the validatorgo error message,
// then to DefaultErrMsg.
var DefaultErrFmtFunc ErrFmtFunc

type ValidationChain struct {
	validator
	modifier
	sanitizer
}

type chainResult struct {
	errors         []ValidationChainError
	location       string
	field          string
	sanitizedValue string
	matched        bool // the field was sent and passed every rule
}

func (v ValidationChain) validate(ctx *gin.Context) chainResult {
	var (
		initialValue   string
		sanitizedValue string
		extractionErr  error
	)

	field := v.validator.field
	reqLoc := v.validator.reqLoc
	location := reqLoc.String()
	errFmtFunc := v.validator.errFmtFunc

	switch v.validator.reqLoc {
	case 0:
		initialValue, extractionErr = extractFieldValFromBody(ctx, field)
	case 1:
		initialValue, extractionErr = extractFieldValFromCookie(ctx, field)
	case 2:
		initialValue, extractionErr = extractFieldValFromHeader(ctx, field)
	case 3:
		initialValue, extractionErr = extractFieldValFromParam(ctx, field)
	case 4:
		initialValue, extractionErr = extractFieldValFromQuery(ctx, field)
	}

	sanitizedValue = initialValue

	// Problems reading the request are reported as errors on the chain, never
	// logged, so the application decides what to do with them.
	switch {
	case extractionErr == nil:
	case errors.Is(extractionErr, ErrExtractionInvalidContentType):
		// The body isn't JSON or a form (or there is no body), so it has no fields
		// this chain can read. The field is treated as not sent, like a missing
		// query parameter: Optional passes and other rules fail as usual.
	case errors.Is(extractionErr, ErrExtractionInvalidJSON):
		// A body that isn't valid JSON has no trustworthy fields, so the chain
		// reports that instead of running its rules. This also applies to Optional
		// chains, which would otherwise treat every field as missing and pass.
		return requestErrorResult(location, field, InvalidJSONErrMsg, InvalidJSONErrCode)
	default:
		// The request could not be read at all, for example because the client
		// disconnected or the body is over a size limit.
		return requestErrorResult(location, field, UnreadableRequestErrMsg, UnreadableRequestErrCode)
	}

	ruleCreators := v.validator.rulesCreatorFuncs
	valErrs := make([]ValidationChainError, 0, len(ruleCreators))

	numOfPreviousValidatorsFailed := 0
	shouldNegateNextValidator := false
	shouldSkipNextValidator := false

	for _, ruleCreator := range ruleCreators {
		if shouldSkipNextValidator {
			shouldSkipNextValidator = false
			continue
		}

		rule := ruleCreator(ctx, initialValue, sanitizedValue)
		vcn := rule.validationChainName
		valid := rule.isValid
		newValue := rule.newValue
		shouldBail := rule.shouldBail
		shouldSkip := rule.shouldSkip
		validationErr := rule.validationErr

		var errMsg string

		switch {
		case errFmtFunc != nil:
			errMsg = errFmtFunc(initialValue, sanitizedValue, vcn)
		case DefaultErrFmtFunc != nil:
			errMsg = DefaultErrFmtFunc(initialValue, sanitizedValue, vcn)
		case validationErr != nil:
			var ve *vgo.ValidationError
			if errors.As(validationErr, &ve) {
				errMsg = ve.Message
			} else {
				errMsg = validationErr.Error()
			}
		default:
			errMsg = DefaultErrMsg
		}

		if rule.validationChainType == 0 {
			if shouldNegateNextValidator {
				valid = !valid
				shouldNegateNextValidator = false
			}

			if !valid {
				numOfPreviousValidatorsFailed++

				order := atomic.AddUint64(&globalErrorOrder, 1)

				var code string
				if validationErr != nil {
					var ve *vgo.ValidationError
					if errors.As(validationErr, &ve) {
						code = ve.Code
					}
				}

				vce := newValidationChainError(
					vceWithLocation(location),
					vceWithMessage(errMsg),
					vceWithField(field),
					vceWithValue(initialValue),
					vceWithCode(code),
					vceWithOrder(order),
				)

				valErrs = append(valErrs, vce)
			}
		}

		if rule.validationChainType == 1 {
			sanitizedValue = newValue
		}

		if rule.validationChainType == 2 {
			vcn := rule.validationChainName

			if vcn == "Bail" {
				if numOfPreviousValidatorsFailed > 0 {
					break
				}
			}

			if vcn == "If" {
				if shouldBail {
					break
				}
			}

			if vcn == "Not" {
				shouldNegateNextValidator = true
			}

			if vcn == "Skip" {
				shouldSkipNextValidator = shouldSkip
			}

			if vcn == "Optional" {
				if initialValue == "" {
					valErrs = make([]ValidationChainError, 0)
					break
				}
			}
		}
	}

	return chainResult{
		errors:         valErrs,
		location:       location,
		field:          field,
		sanitizedValue: sanitizedValue,
		matched:        len(valErrs) == 0 && fieldPresent(ctx, reqLoc, field),
	}
}

func (v ValidationChain) Validate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		saveChainResultToCtx(ctx, v.validate(ctx))
		ctx.Next()
	}
}

func newValidationChain(field string, errFmtFunc ErrFmtFunc, reqLoc RequestLocation) ValidationChain {
	return ValidationChain{
		validator: newValidator(field, errFmtFunc, reqLoc),
		modifier:  newModifier(field, errFmtFunc, reqLoc),
		sanitizer: newSanitizer(field, errFmtFunc, reqLoc),
	}
}
