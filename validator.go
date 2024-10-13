package ginvalidator

import (
	"net/http"

	vgo "github.com/bube054/validatorgo"
	"github.com/gin-gonic/gin"
)

type validator struct {
	field      string
	errFmtFunc *ErrFmtFuncHandler

	reqLoc            requestLocation
	rulesCreatorFuncs ruleCreatorFuncs
}

func (v *validator) recreateValidationChainFromValidator(ruleCreatorFunc ruleCreatorFunc) ValidationChain {
	newRulesCreatorFunc := append(v.rulesCreatorFuncs, ruleCreatorFunc)

	return ValidationChain{
		validator: validator{
			field:             v.field,
			reqLoc:            v.reqLoc,
			errFmtFunc:        v.errFmtFunc,
			rulesCreatorFuncs: newRulesCreatorFunc,
		},
		modifier: modifier{
			field:             v.field,
			reqLoc:            v.reqLoc,
			errFmtFunc:        v.errFmtFunc,
			rulesCreatorFuncs: newRulesCreatorFunc,
		},
		sanitizer: sanitizer{
			field:             v.field,
			reqLoc:            v.reqLoc,
			errFmtFunc:        v.errFmtFunc,
			rulesCreatorFuncs: newRulesCreatorFunc,
		},
	}
}

type CustomValidatorFunc func(req http.Request, initialValue, sanitizedValue string) bool

func (v validator) CustomValidator(cvf CustomValidatorFunc) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		httpRequest := ctx.Request
		isValid := cvf(*httpRequest, initialValue, sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName("CustomValidator"),
			withValidationChainType(validatorType),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

func (v validator) Contains(seed string, opts *vgo.ContainsOpt) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid := vgo.Contains(sanitizedValue, seed, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName("Contains"),
			withValidationChainType(validatorType),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

func (v validator) Equals(comparison string) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid := vgo.Equals(sanitizedValue, comparison)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName("Equals"),
			withValidationChainType(validatorType),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

func (v validator) AbaRouting() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid := vgo.IsAbaRouting(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName("AbaRouting"),
			withValidationChainType(validatorType),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

func (v validator) After(opts *vgo.IsAfterOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid := vgo.IsAfter(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName("After"),
			withValidationChainType(validatorType),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

func (v validator) Alphanumeric(opts *vgo.IsAlphanumericOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid := vgo.IsAlphanumeric(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName("Alphanumeric"),
			withValidationChainType(validatorType),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

func (v validator) Ascii() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid := vgo.IsAscii(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName("Ascii"),
			withValidationChainType(validatorType),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

func (v validator) Base32(opts *vgo.IsBase32Opts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid := vgo.IsBase32(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName("Base32"),
			withValidationChainType(validatorType),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

func (v validator) Base58() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid := vgo.IsBase58(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName("Base58"),
			withValidationChainType(validatorType),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

func newValidator(field string, errFmtFunc *ErrFmtFuncHandler, reqLoc requestLocation) validator {
	return validator{
		field:      field,
		errFmtFunc: errFmtFunc,
		reqLoc:     reqLoc,
	}
}
