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

	"github.com/gin-gonic/gin"
)

type ValidationChain struct {
	validator
	modifier
	sanitizer
}

func (v ValidationChain) Validate() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		field := v.validator.field
		reqLoc := v.validator.reqLoc
		location := reqLoc.string()
		errFmtFunc := v.validator.errFmtFunc

		var initialValue, sanitizedValue string
		var extractionErr error

		if v.validator.reqLoc == 0 {
			initialValue, extractionErr = extractFieldValFromBody(field, ctx)
			sanitizedValue = initialValue
		}
		if v.validator.reqLoc == 1 {
			initialValue, extractionErr = extractFieldValFromCookie(field, ctx)
			sanitizedValue = initialValue
		}
		if v.validator.reqLoc == 2 {
			initialValue, extractionErr = extractFieldValFromHeader(field, ctx)
			sanitizedValue = initialValue
		}
		if v.validator.reqLoc == 3 {
			initialValue, extractionErr = extractFieldValFromParam(field, ctx)
			sanitizedValue = initialValue
		}
		if v.validator.reqLoc == 4 {
			initialValue, extractionErr = extractFieldValFromQuery(field, ctx)
			sanitizedValue = initialValue
		}

		if extractionErr != nil {
			panic(fmt.Errorf("for request location: %s, could not extract field: %s", reqLoc.string(), field))
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

			if shouldNegateNextValidator {
				valid = !valid
				shouldNegateNextValidator = false
			}

			newValue := rule.newValue
			shouldBail := rule.shouldBail
			shouldSkip := rule.shouldSkip

			var errMsg string

			if errFmtFunc == nil {
				errMsg = "Invalid value"
			} else {
				eff := *errFmtFunc
				errMsg = eff(initialValue, sanitizedValue, vcn)
			}

			// rule is for validators
			if rule.validationChainType == 0 {
				if !valid {
					numOfPreviousValidatorsFailed++

					vce := NewValidationChainError(
						vceWithLocation(location),
						vceWithMsg(errMsg),
						vceWithField(field),
						vceWithValue(initialValue),
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
			}
		}

		saveErrorsToCtx(ctx, valErrs)
		saveSanitizedDataToCtx(ctx, location, field, sanitizedValue)

		ctx.Next()
	}
}

func NewValidationChain(field string, errFmtFunc *ErrFmtFuncHandler, reqLoc requestLocation) ValidationChain {
	return ValidationChain{
		validator: newValidator(field, errFmtFunc, reqLoc),
		modifier:  newModifier(field, errFmtFunc, reqLoc),
		sanitizer: newSanitizer(field, errFmtFunc, reqLoc),
	}
}
