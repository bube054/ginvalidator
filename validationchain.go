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
			initialValue, extractionErr = extractBodyValue(field, ctx)
			sanitizedValue = initialValue
		}
		if v.validator.reqLoc == 1 {
			initialValue, extractionErr = extractCookieValue(field, ctx)
			sanitizedValue = initialValue
		}
		if v.validator.reqLoc == 2 {
			initialValue, extractionErr = extractHeaderValue(field, ctx)
			sanitizedValue = initialValue
		}
		if v.validator.reqLoc == 3 {
			initialValue, extractionErr = extractParamValue(field, ctx)
			sanitizedValue = initialValue
		}
		if v.validator.reqLoc == 4 {
			initialValue, extractionErr = extractQueryValue(field, ctx)
			sanitizedValue = initialValue
		}

		if extractionErr != nil {
			panic(fmt.Errorf("for request location: %s, could not extract field: %s", reqLoc.string(), field))
		}

		ruleCreators := v.validator.rulesCreatorFuncs
		posErrs := make([]ValidationChainError, 0, len(ruleCreators))

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
				}

				vce := NewValidationChainError(
					VCEWithIsValid(valid),
					VCEWithLocation(location),
					VCEWithMsg(errMsg),
					VCEWithField(field),
					VCEWithValue(initialValue),
					VCEWithSanitizedValue(sanitizedValue),
				)

				posErrs = append(posErrs, vce)
			}

			// rule is for sanitizers
			if rule.validationChainType == 1 {
				sanitizedValue = newValue

				vce := NewValidationChainError(
					VCEWithIsValid(true),
					VCEWithLocation(location),
					VCEWithMsg(""),
					VCEWithField(field),
					VCEWithValue(initialValue),
					VCEWithSanitizedValue(newValue),
				)

				posErrs = append(posErrs, vce)
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

				vce := NewValidationChainError(
					VCEWithIsValid(true),
					VCEWithLocation(location),
					VCEWithMsg(""),
					VCEWithField(field),
					VCEWithValue(initialValue),
					VCEWithSanitizedValue(newValue),
				)

				posErrs = append(posErrs, vce)
			}
		}

		errs := make([]ValidationChainError, 0, cap(posErrs))

		for _, err := range posErrs {
			if !err.isValid {
				errs = append(errs, err)
			}
		}

		fmt.Println("All errors: ", errs)

		saveErrorsToCtx(ctx, errs)

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
