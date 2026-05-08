// Package ginvalidator is a set of [Gin] middlewares that wraps the
// extensive collection of validators and sanitizers offered by [validatorgo].
//
// It allows you to combine them in many ways so that you can validate and sanitize your express requests,
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
	"log"
	"sync/atomic"

	vgo "github.com/bube054/validatorgo"
	"github.com/gin-gonic/gin"
)

var globalErrorOrder uint64

const DefaultValChainErrMsg string = "Invalid value"

// DefaultErrFmtFunc is a package-level fallback error message formatter.
// When set, it is used for any validation chain that does not have its own errFmtFunc.
// If nil (the default), the system falls back to the validatorgo error message,
// then to DefaultValChainErrMsg.
var DefaultErrFmtFunc ErrFmtFuncHandler

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

	if extractionErr != nil {
		log.Printf("Error extracting field %q from request location %q: %v", field, location, extractionErr)
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
			errMsg = DefaultValChainErrMsg
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

				vce := NewValidationChainError(
					vceWithLocation(location),
					vceWithMsg(errMsg),
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
	}
}

func (v ValidationChain) Validate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result := v.validate(ctx)
		saveValidationErrorsToCtx(ctx, result.errors)
		saveMatchedDataToCtx(ctx, result.location, result.field, result.sanitizedValue)
		ctx.Next()
	}
}

func NewValidationChain(field string, errFmtFunc ErrFmtFuncHandler, reqLoc RequestLocation) ValidationChain {
	return ValidationChain{
		validator: newValidator(field, errFmtFunc, reqLoc),
		modifier:  newModifier(field, errFmtFunc, reqLoc),
		sanitizer: newSanitizer(field, errFmtFunc, reqLoc),
	}
}
