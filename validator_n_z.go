package ginvalidator

import (
	"regexp"

	vgo "github.com/bube054/validatorgo"
	"github.com/gin-gonic/gin"
)

// Numeric is a validator that checks if a string is a number.
//
// This function uses the [IsNumeric] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsNumeric]: https://pkg.go.dev/github.com/bube054/validatorgo#IsNumeric
func (v validator) Numeric(opts *vgo.IsNumericOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsNumeric(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(NumericValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Octal is a  validator to check that a value is a json object.
//
// This function uses the [IsObject] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsObject]: https://pkg.go.dev/github.com/bube054/validatorgo#IsObject
func (v validator) Object(opts *vgo.IsObjectOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsObject(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(OctalValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Octal is a validator that checks if the string is a valid octal number.
//
// This function uses the [IsOctal] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsOctal]: https://pkg.go.dev/github.com/bube054/validatorgo#IsOctal
func (v validator) Octal() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsOctal(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(OctalValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// PassportNumber is a validator that checks if the string is a valid passport number.
//
// This function uses the [IsPassportNumber] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsPassportNumber]: https://pkg.go.dev/github.com/bube054/validatorgo#IsPassportNumber
func (v validator) PassportNumber(countryCode string) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsPassportNumber(sanitizedValue, countryCode)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(PassportNumberValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Port is a validator that checks if the string is a valid port number.
//
// This function uses the [IsPort] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsPort]: https://pkg.go.dev/github.com/bube054/validatorgo#IsPort
func (v validator) Port() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsPort(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(PortValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// PostalCode is a validator that checks if the string is a postal code.
//
// This function uses the [IsPostalCode] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsPostalCode]: https://pkg.go.dev/github.com/bube054/validatorgo#IsPostalCode
func (v validator) PostalCode(locale string) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsPostalCode(sanitizedValue, locale)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(PostalCodeValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// RFC3339 is a validator that checks if the string is a valid RFC 3339 date.
//
// This function uses the [IsRFC3339] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsRFC3339]: https://pkg.go.dev/github.com/bube054/validatorgo#IsRFC3339
func (v validator) RFC3339() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsRFC3339(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(RFC3339ValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// RgbColor is a validator that checks if the string is a rgb or rgba color.
//
// This function uses the [IsRgbColor] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsRgbColor]: https://pkg.go.dev/github.com/bube054/validatorgo#IsRgbColor
func (v validator) RgbColor(opts *vgo.IsRgbOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsRgbColor(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(RgbColorValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// SemVer is a validator that checks if the string is a Semantic Versioning Specification (SemVer).
//
// This function uses the [IsSemVer] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsSemVer]: https://pkg.go.dev/github.com/bube054/validatorgo#IsSemVer
func (v validator) SemVer() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsSemVer(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(SemVerValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Slug is a validator that checks if the string is of type slug.
//
// This function uses the [IsSlug] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsSlug]: https://pkg.go.dev/github.com/bube054/validatorgo#IsSlug
func (v validator) Slug() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsSlug(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(SlugValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// StrongPassword is a validator that checks if the string is of type strongPassword.
//
// This function uses the [IsStrongPassword] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsStrongPassword]: https://pkg.go.dev/github.com/bube054/validatorgo#IsStrongPassword
func (v validator) StrongPassword(opts *vgo.IsStrongPasswordOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, _, vErr := vgo.IsStrongPassword(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(StrongPasswordValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// TaxID is a validator that checks if the string is a valid Tax Identification Number.
//
// This function uses the [IsTaxID] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsTaxID]: https://pkg.go.dev/github.com/bube054/validatorgo#IsTaxID
func (v validator) TaxID(locale string) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsTaxID(sanitizedValue, locale)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(TaxIDValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// SurrogatePair is a validator that checks if the string contains any surrogate pairs chars.
//
// This function uses the [IsSurrogatePair] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsSurrogatePair]: https://pkg.go.dev/github.com/bube054/validatorgo#IsSurrogatePair
func (v validator) SurrogatePair() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsSurrogatePair(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(SurrogatePairValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Time is a validator that checks if the string is a valid time e.g. 23:01:59
//
// This function uses the [IsTime] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsTime]: https://pkg.go.dev/github.com/bube054/validatorgo#IsTime
func (v validator) Time(opts *vgo.IsTimeOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsTime(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(TimeValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// ULID is a validator that checks if the string is a ULID.
//
// This function uses the [IsULID] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsULID]: https://pkg.go.dev/github.com/bube054/validatorgo#IsULID
func (v validator) ULID() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsULID(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ULIDValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// UpperCase is a validator that checks if the string is uppercase.
//
// This function uses the [IsUpperCase] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsUpperCase]: https://pkg.go.dev/github.com/bube054/validatorgo#IsUpperCase
func (v validator) UpperCase() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsUpperCase(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(UpperCaseValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// URL is a validator that checks if the string is URL.
//
// This function uses the [IsURL] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsURL]: https://pkg.go.dev/github.com/bube054/validatorgo#IsURL
func (v validator) URL(opts *vgo.IsURLOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsURL(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(URLValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// UUID is a validator that checks if the string is an RFC9562 UUID.
//
// This function uses the [IsUUID] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsUUID]: https://pkg.go.dev/github.com/bube054/validatorgo#IsUUID
func (v validator) UUID(version string) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsUUID(sanitizedValue, version)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(UUIDValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// VariableWidth is a validator that checks if the string contains a mixture of full and half-width chars.
//
// This function uses the [IsVariableWidth] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsVariableWidth]: https://pkg.go.dev/github.com/bube054/validatorgo#IsVariableWidth
func (v validator) VariableWidth() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsVariableWidth(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(VariableWidthValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// VAT is a validator that checks if the string is a valid VAT.
//
// This function uses the [IsVAT] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsVAT]: https://pkg.go.dev/github.com/bube054/validatorgo#IsVAT
func (v validator) VAT(countryCode string) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsVAT(sanitizedValue, countryCode)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(VATValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Whitelisted is a validator that checks if the string consists only of characters that appear in the whitelist chars.
//
// This function uses the [IsWhitelisted] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsWhitelisted]: https://pkg.go.dev/github.com/bube054/validatorgo#IsWhitelisted
func (v validator) Whitelisted(chars string) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsWhitelisted(sanitizedValue, chars)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(WhitelistedValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Matches is a validator that checks if the string matches the regex.
//
// This function uses the [IsMatches] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsMatches]: https://pkg.go.dev/github.com/bube054/validatorgo#IsMatches
func (v validator) Matches(re *regexp.Regexp) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.Matches(sanitizedValue, re)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(MatchesValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}
