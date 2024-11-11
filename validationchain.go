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
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

const DefaultValChainErrMsg string = "Invalid value"

type ValidationChain struct {
	validator
	modifier
	sanitizer
}

func (v ValidationChain) Validate() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var (
			initialValue   string
			sanitizedValue string
			extractionErr  error
		)

		field := v.validator.field
		reqLoc := v.validator.reqLoc
		location := reqLoc.string()
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
			fmt.Println(fmt.Errorf("for request location: %q, could not extract field: %q, err: %w", location, field, extractionErr))
		}

		ruleCreators := v.validator.rulesCreatorFuncs
		valErrs := make([]ValidationChainError, 0, len(ruleCreators))

		numOfPreviousValidatorsFailed := 0 // counter for dealing with previous failed validations, used by bail.
		shouldNegateNextValidator := false // state for dealing with the immediate previous validation validity and negating it, used by not.
		shouldSkipNextValidator := false   // state for dealing with whether to skip next link in the validation chain. used by skip.

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

			var errMsg string

			if errFmtFunc == nil {
				errMsg = DefaultValChainErrMsg
			} else {
				errMsg = errFmtFunc(initialValue, sanitizedValue, vcn)
			}

			// rule is for validators
			if rule.validationChainType == 0 {
				if shouldNegateNextValidator {
					valid = !valid
					shouldNegateNextValidator = false
				}

				if !valid {
					numOfPreviousValidatorsFailed++

					vce := NewValidationChainError(
						vceWithLocation(location),
						vceWithMsg(errMsg),
						vceWithField(field),
						vceWithValue(initialValue),
						vceWithCreatedAt(time.Now()),
					)

					valErrs = append(valErrs, vce)
				}
			}

			// rule is for sanitizers
			if rule.validationChainType == 1 {
				sanitizedValue = newValue
			}

			// rule is for modifiers
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
					if sanitizedValue == "" {
						valErrs = make([]ValidationChainError, 0, len(ruleCreators))
						break
					}
				}
			}
		}

		saveValidationErrorsToCtx(ctx, valErrs)
		saveMatchedDataToCtx(ctx, location, field, sanitizedValue)

		ctx.Next()
	}
}

func NewValidationChain(field string, errFmtFunc ErrFmtFuncHandler, reqLoc requestLocation) ValidationChain {
	return ValidationChain{
		validator: newValidator(field, errFmtFunc, reqLoc),
		modifier:  newModifier(field, errFmtFunc, reqLoc),
		sanitizer: newSanitizer(field, errFmtFunc, reqLoc),
	}
}
