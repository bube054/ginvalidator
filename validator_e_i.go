package ginvalidator

import (
	vgo "github.com/bube054/validatorgo"
	"github.com/gin-gonic/gin"
)

// EAN is validator that checks if the string is a valid EAN (European Article Number).
//
// This function uses the [IsEAN] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsEAN]: https://pkg.go.dev/github.com/bube054/validatorgo#IsEAN
func (v validator) EAN() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsEAN(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(EANValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Email is a validator that checks if the string is an email.
//
// This function uses the [IsEmail] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsEmail]: https://pkg.go.dev/github.com/bube054/validatorgo#IsEmail
func (v validator) Email(opts *vgo.IsEmailOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsEmail(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(EmailValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Empty is a validator that checks if the string is an email.
//
// This function uses the [IsEmpty] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsEmpty]: https://pkg.go.dev/github.com/bube054/validatorgo#IsEmpty
func (v validator) Empty(opts *vgo.IsEmptyOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsEmpty(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(EmptyValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// EthereumAddress is a validator checks if the string is an Ethereum address.
//
// This function uses the [IsEthereumAddress] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsEthereumAddress]: https://pkg.go.dev/github.com/bube054/validatorgo#IsEthereumAddress
func (v validator) EthereumAddress() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsEthereumAddress(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(EthereumAddressValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Float is a validator that checks if the string is a float.
//
// This function uses the [IsFloat] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsFloat]: https://pkg.go.dev/github.com/bube054/validatorgo#IsFloat
func (v validator) Float(opts *vgo.IsFloatOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsFloat(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(FloatValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// FQDN is a validator that checks if the string is a fully qualified domain name (e.g. domain.com).
//
// This function uses the [IsFQDN] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsFQDN]: https://pkg.go.dev/github.com/bube054/validatorgo#IsFQDN
func (v validator) FQDN(opts *vgo.IsFQDNOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsFQDN(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(FQDNValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// FreightContainerID is a validator that checks alias for IsISO6346, check if the string is a valid ISO 6346 shipping container identification.
//
// This function uses the [IsFreightContainerID] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsFreightContainerID]: https://pkg.go.dev/github.com/bube054/validatorgo#IsFreightContainerID
func (v validator) FreightContainerID() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsFreightContainerID(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(FreightContainerIDValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// FullWidth validator that checks if the string contains any full-width chars.
//
// This function uses the [IsFullWidth] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsFullWidth]: https://pkg.go.dev/github.com/bube054/validatorgo#IsFullWidth
func (v validator) FullWidth() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsFullWidth(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(FullWidthValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// HalfWidth is a validator that checks if the string contains any half-width chars.
//
// This function uses the [IsHalfWidth] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsHalfWidth]: https://pkg.go.dev/github.com/bube054/validatorgo#IsHalfWidth
func (v validator) HalfWidth() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsHalfWidth(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(HalfWidthValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Hash is a validator that checks if the string is a hash of type algorithm.
//
// This function uses the [IsHash] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsHash]: https://pkg.go.dev/github.com/bube054/validatorgo#IsHash
func (v validator) Hash(algorithm string) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsHash(sanitizedValue, algorithm)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(HashValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Hexadecimal is a validator that checks if the string is a hexadecimal number.
//
// This function uses the [IsHexadecimal] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsHexadecimal]: https://pkg.go.dev/github.com/bube054/validatorgo#IsHexadecimal
func (v validator) Hexadecimal() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsHexadecimal(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(HexadecimalValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// HexColor is a validator that checks if the string is a hexadecimal color.
//
// This function uses the [IsHexColor] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsHexColor]: https://pkg.go.dev/github.com/bube054/validatorgo#IsHexColor
func (v validator) HexColor() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsHexColor(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(HexColorValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// HSL is a validator that checks if the string is an HSL (hue, saturation, lightness, optional alpha) color based on CSS Colors Level 4 specification.
//
// This function uses the [IsHSL] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsHSL]: https://pkg.go.dev/github.com/bube054/validatorgo#IsHSL
func (v validator) HSL() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsHSL(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(HSLValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// IBAN is a validator that checks if the string is an IBAN (International Bank Account Number).
//
// This function uses the [IsIBAN] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsIBAN]: https://pkg.go.dev/github.com/bube054/validatorgo#IsIBAN
func (v validator) IBAN(countryCode string) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsIBAN(sanitizedValue, countryCode)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(IBANValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// IdentityCard is a validator that checks if the string is a valid identity card code.
//
// This function uses the [IsIdentityCard] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsIdentityCard]: https://pkg.go.dev/github.com/bube054/validatorgo#IsIdentityCard
func (v validator) IdentityCard(locale string) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsIdentityCard(sanitizedValue, locale)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(IdentityCardValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// IMEI is a validator that checks if the string is a valid IMEI number.
//
// This function uses the [IsIMEI] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsIMEI]: https://pkg.go.dev/github.com/bube054/validatorgo#IsIMEI
func (v validator) IMEI(opts *vgo.IsIMEIOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsIMEI(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(IMEIValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// In is a validator that checks if the string is in a slice of allowed values.
//
// This function uses the [IsIn] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsIn]: https://pkg.go.dev/github.com/bube054/validatorgo#IsIn
func (v validator) In(values []string) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsIn(sanitizedValue, values)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(InValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Int is a validator that checks if the string is an integer.
//
// This function uses the [IsInt] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsInt]: https://pkg.go.dev/github.com/bube054/validatorgo#IsInt
func (v validator) Int(opts *vgo.IsIntOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsInt(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(IntValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// IP is a validator that checks if the string is an IP (version 4 or 6).
//
// This function uses the [IsIP] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsIP]: https://pkg.go.dev/github.com/bube054/validatorgo#IsIP
func (v validator) IP(version string) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsIP(sanitizedValue, version)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(IPValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// IPRange is a validator that checks if the string is an IPRange (version 4 or 6).
//
// This function uses the [IsIPRange] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsIPRange]: https://pkg.go.dev/github.com/bube054/validatorgo#IsIPRange
func (v validator) IPRange(version string) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsIPRange(sanitizedValue, version)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(IPRangeValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

