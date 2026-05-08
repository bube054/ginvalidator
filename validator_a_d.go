package ginvalidator

import (
	vgo "github.com/bube054/validatorgo"
	"github.com/gin-gonic/gin"
)

// AbaRouting is a validator that checks if the string is an ABA routing number for US bank account / cheque.
//
// This function uses the [IsAbaRouting] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsAbaRouting]: https://pkg.go.dev/github.com/bube054/validatorgo#IsAbaRouting
func (v validator) AbaRouting() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsAbaRouting(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(AbaRoutingValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// After is a validator that checks if the string is a date that is after the specified date.
//
// This function uses the [IsAfter] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsAfter]: https://pkg.go.dev/github.com/bube054/validatorgo#IsAfter
func (v validator) After(opts *vgo.IsAfterOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsAfter(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(AfterValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Alpha is a validator that checks if the string contains only letters (a-zA-Z).
//
// This function uses the [IsAlpha] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsAlpha]: https://pkg.go.dev/github.com/bube054/validatorgo#IsAlpha
func (v validator) Alpha(opts *vgo.IsAlphaOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsAlpha(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(AlphaValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Alphanumeric is a validator that checks if the string contains only letters and numbers (a-zA-Z0-9).
//
// This function uses the [IsAlphanumeric] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsAlphanumeric]: https://pkg.go.dev/github.com/bube054/validatorgo#IsAlphanumeric
func (v validator) Alphanumeric(opts *vgo.IsAlphanumericOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsAlphanumeric(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(AlphanumericValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Base32 is a validator to check that a value is an array.
//
// This function uses the [IsArray] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsArray]: https://pkg.go.dev/github.com/bube054/validatorgo#IsArray
func (v validator) Array(opts *vgo.IsArrayOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsArray(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ArrayValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Ascii is a validator that checks if the string contains ASCII chars only.
//
// This function uses the [IsAscii] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsAscii]: https://pkg.go.dev/github.com/bube054/validatorgo#IsAscii
func (v validator) Ascii() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsAscii(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(AbaRoutingValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Base32 is a validator that checks if the string is base32 encoded.
//
// This function uses the [IsBase32] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsBase32]: https://pkg.go.dev/github.com/bube054/validatorgo#IsBase32
func (v validator) Base32(opts *vgo.IsBase32Opts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsBase32(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(Base32ValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Base58 is a validator that checks if the string is base32 encoded.
//
// This function uses the [IsBase58] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsBase58]: https://pkg.go.dev/github.com/bube054/validatorgo#IsBase58
func (v validator) Base58() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsBase58(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(Base58ValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Base64 is a validator that checks if the string is base64 encoded.
//
// This function uses the [IsBase64] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsBase64]: https://pkg.go.dev/github.com/bube054/validatorgo#IsBase64
func (v validator) Base64(opts *vgo.IsBase64Opts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsBase64(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(Base64ValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Before is a validator that checks if the string is a date that is before the specified date.
//
// This function uses the [IsBefore] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsBefore]: https://pkg.go.dev/github.com/bube054/validatorgo#IsBefore
func (v validator) Before(opts *vgo.IsBeforeOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsBefore(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(BeforeValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Bic is a validator that checks if the string is a BIC (Bank Identification Code) or SWIFT code.
//
// This function uses the [IsBic] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsBic]: https://pkg.go.dev/github.com/bube054/validatorgo#IsBic
func (v validator) Bic() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsBic(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(BicValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Boolean validator that checks if the string is a boolean.
//
// This function uses the [IsBoolean] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsBoolean]: https://pkg.go.dev/github.com/bube054/validatorgo#IsBoolean
func (v validator) Boolean(opts *vgo.IsBooleanOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsBoolean(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(BooleanValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// BTCAddress is a validator that checks if the string is a valid BTC address.
//
// This function uses the [IsBTCAddress] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsBTCAddress]: https://pkg.go.dev/github.com/bube054/validatorgo#IsBTCAddress
func (v validator) BTCAddress() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsBTCAddress(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(BTCAddressValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// ByteLength is a validator that checks if the string's length (in UTF-8 bytes) falls in a range.
//
// This function uses the [IsByteLength] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsByteLength]: https://pkg.go.dev/github.com/bube054/validatorgo#IsByteLength
func (v validator) ByteLength(opts *vgo.IsByteLengthOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsByteLength(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ByteLengthValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// CreditCard is a validator that checks if the string is a credit card number.
//
// This function uses the [IsCreditCard] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsCreditCard]: https://pkg.go.dev/github.com/bube054/validatorgo#IsCreditCard
func (v validator) CreditCard(opts *vgo.IsCreditCardOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsCreditCard(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(CreditCardValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Currency is a validator that checks if the string is a valid currency amount.
//
// This function uses the [IsCurrency] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsCurrency]: https://pkg.go.dev/github.com/bube054/validatorgo#IsCurrency
func (v validator) Currency(opts *vgo.IsCurrencyOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsCurrency(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(CurrencyValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// DataURI is a validator that checks if the string is a data uri format.
//
// This function uses the [IsDataURI] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsDataURI]: https://pkg.go.dev/github.com/bube054/validatorgo#IsDataURI
func (v validator) DataURI() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsDataURI(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(DataURIValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Date is a validator that checks if the string is a valid date. e.g. 2002-07-15.
//
// This function uses the [IsDate] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsDate]: https://pkg.go.dev/github.com/bube054/validatorgo#IsDate
func (v validator) Date(opts *vgo.IsDateOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsDate(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(DataURIValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Decimal is a validator that checks if the string represents a decimal number, such as 0.1, .3, 1.1, 1.00003, 4.0, etc.
//
// This function uses the [IsDecimal] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsDecimal]: https://pkg.go.dev/github.com/bube054/validatorgo#IsDecimal
func (v validator) Decimal(opts *vgo.IsDecimalOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsDecimal(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(DecimalValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// DivisibleBy is a validator thats checks if the string is a number(integer not a floating point) that is divisible by another(integer not a floating point).
//
// This function uses the [IsDivisibleBy] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsDivisibleBy]: https://pkg.go.dev/github.com/bube054/validatorgo#IsDivisibleBy
func (v validator) DivisibleBy(num int) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsDivisibleBy(sanitizedValue, num)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(DivisibleByValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}
