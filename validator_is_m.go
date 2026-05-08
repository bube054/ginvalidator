// Validator methods in alphabetical order: ISIN through Multibyte.

package ginvalidator

import (
	vgo "github.com/bube054/validatorgo"
	"github.com/gin-gonic/gin"
)

// ISIN is a validator that checks if the string is an ISIN (stock/security identifier).
//
// This function uses the [IsISIN] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsISIN]: https://pkg.go.dev/github.com/bube054/validatorgo#IsISIN
func (v validator) ISIN() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsISIN(sanitizedValue)

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

// ISO4217 is a validator that checks if the string is a valid ISO 4217 officially assigned.
//
// This function uses the [IsISO4217] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsISO4217]: https://pkg.go.dev/github.com/bube054/validatorgo#IsISO4217
func (v validator) ISO4217() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsIso4217(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ISO4217ValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// ISO6346 is a validator that checks if the string is a valid ISO 6346 shipping container identification.
//
// This function uses the [IsISO6346] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsISO6346]: https://pkg.go.dev/github.com/bube054/validatorgo#IsISO6346
func (v validator) ISO6346() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsISO6346(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ISO6346ValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// ISO6391 is a validator that checks if the string is a valid ISO 639-1 language code.
//
// This function uses the [IsISO6391] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsISO6391]: https://pkg.go.dev/github.com/bube054/validatorgo#IsISO6391
func (v validator) ISO6391() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsISO6391(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ISO6391ValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// ISO8601 is a validator that checks if the string is a valid ISO 8601 date.
//
// This function uses the [IsISO8601] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsISO8601]: https://pkg.go.dev/github.com/bube054/validatorgo#IsISO8601
func (v validator) ISO8601(opts *vgo.IsISO8601Opts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsISO8601(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ISO8601ValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// ISO31661Alpha2 is a validator that checks if the string is a valid ISO 3166-1 alpha-2 officially assigned country code.
//
// This function uses the [IsISO31661Alpha2] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsISO31661Alpha2]: https://pkg.go.dev/github.com/bube054/validatorgo#IsISO31661Alpha2
func (v validator) ISO31661Alpha2() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsISO31661Alpha2(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ISO31661Alpha2ValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// ISO31661Alpha3 is a validator that checks if the string is a valid ISO 3166-1 alpha-2 officially assigned country code.
//
// This function uses the [IsISO31661Alpha3] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsISO31661Alpha3]: https://pkg.go.dev/github.com/bube054/validatorgo#IsISO31661Alpha3
func (v validator) ISO31661Alpha3() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsISO31661Alpha3(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ISO31661Alpha3ValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// ISO31661Numeric is a validator that checks check if the string is a valid ISO 3166-1 numeric officially assigned country code.
//
// This function uses the [IsISO31661Numeric] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsISO31661Numeric]: https://pkg.go.dev/github.com/bube054/validatorgo#IsISO31661Numeric
func (v validator) ISO31661Numeric() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsISO31661Numeric(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ISO31661NumericValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// ISRC is a validator that checks if the string is an ISRC.
//
// This function uses the [IsISRC] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsISRC]: https://pkg.go.dev/github.com/bube054/validatorgo#IsISRC
func (v validator) ISRC(allowHyphens bool) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsISRC(sanitizedValue, allowHyphens)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ISRCValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// ISSN is a validator that checks if the string is an ISSN.
//
// This function uses the [IsISSN] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsISSN]: https://pkg.go.dev/github.com/bube054/validatorgo#IsISSN
func (v validator) ISSN(opts *vgo.IsISSNOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsISSN(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ISSNValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// JSON is a validator that checks if the string is an JSON.
//
// This function uses the [IsJSON] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsJSON]: https://pkg.go.dev/github.com/bube054/validatorgo#IsJSON
func (v validator) JSON() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsJSON(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(JSONValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// LatLong is a validator that checks if the string is a valid latitude-longitude coordinate.
//
// This function uses the [IsLatLong] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsLatLong]: https://pkg.go.dev/github.com/bube054/validatorgo#IsLatLong
func (v validator) LatLong(opts *vgo.IsLatLongOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsLatLong(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(LatLongValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// length is a validator that checks if the string's length falls in a range.
//
// This function uses the [IsLength] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsLength]: https://pkg.go.dev/github.com/bube054/validatorgo#IsLength
func (v validator) Length(opts *vgo.IsLengthOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsLength(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(LengthValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// LicensePlate is a validator that checks if the string matches the format of a country's license plate.
//
// This function uses the [IsLicensePlate] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsLicensePlate]: https://pkg.go.dev/github.com/bube054/validatorgo#IsLicensePlate
func (v validator) LicensePlate(locale string) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsLicensePlate(sanitizedValue, locale)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(LicensePlateValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Locale is a validator that checks if the string is a locale.
//
// This function uses the [IsLocale] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsLocale]: https://pkg.go.dev/github.com/bube054/validatorgo#IsLocale
func (v validator) Locale() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsLocale(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(LocaleValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// LowerCase is a validator that checks if the string is lowercase.
//
// This function uses the [IsLowerCase] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsLowerCase]: https://pkg.go.dev/github.com/bube054/validatorgo#IsLowerCase
func (v validator) LowerCase() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsLowerCase(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(LowerCaseValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// LuhnNumber is a validator that checks if the string passes the Luhn algorithm check.
//
// This function uses the [IsLuhnNumber] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsLuhnNumber]: https://pkg.go.dev/github.com/bube054/validatorgo#IsLuhnNumber
func (v validator) LuhnNumber() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsLuhnNumber(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(LuhnNumberValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// MacAddress is a validator that checks if the string is a MAC address.
//
// This function uses the [IsMacAddress] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsMacAddress]: https://pkg.go.dev/github.com/bube054/validatorgo#IsMacAddress
func (v validator) MacAddress(opts *vgo.IsMacAddressOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsMacAddress(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(MacAddressValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// MagnetURI is a validator that checks if the string is a Magnet URI format.
//
// This function uses the [IsMagnetURI] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsMagnetURI]: https://pkg.go.dev/github.com/bube054/validatorgo#IsMagnetURI
func (v validator) MagnetURI() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsMagnetURI(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(MagnetURIValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// MailtoURI is a validator that checks if the string is a Mailto URI format.
//
// This function uses the [IsMailtoURI] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsMailtoURI]: https://pkg.go.dev/github.com/bube054/validatorgo#IsMailtoURI
func (v validator) MailtoURI(opts *vgo.IsMailToURIOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsMailtoURI(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(MailtoURIValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// MD5 is a validator that checks if the string is a MD5 hash.
//
// This function uses the [IsMD5] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsMD5]: https://pkg.go.dev/github.com/bube054/validatorgo#IsMD5
func (v validator) MD5() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsMD5(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(MD5ValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// MimeType is a validator that checks if the string matches to a valid MIME type format.
//
// This function uses the [IsMimeType] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsMimeType]: https://pkg.go.dev/github.com/bube054/validatorgo#IsMimeType
func (v validator) MimeType() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsMimeType(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(MimeTypeValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// MobilePhone is a validator that checks if the string is a mobile phone number.
//
// This function uses the [IsMobilePhone] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsMobilePhone]: https://pkg.go.dev/github.com/bube054/validatorgo#IsMobilePhone
func (v validator) MobilePhone(locales []string, opts *vgo.IsMobilePhoneOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsMobilePhone(sanitizedValue, locales, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(MobilePhoneValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// MongoID is a validator that checks if the string is a valid hex-encoded representation of a MongoDB ObjectId.
//
// This function uses the [IsMongoID] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsMongoID]: https://pkg.go.dev/github.com/bube054/validatorgo#IsMongoID
func (v validator) MongoID() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsMongoID(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(MongoIDValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Multibyte is a validator that checks if the string contains one or more multibyte chars.
//
// This function uses the [IsMultibyte] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsMultibyte]: https://pkg.go.dev/github.com/bube054/validatorgo#IsMultibyte
func (v validator) Multibyte() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid, vErr := vgo.IsMultibyte(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(MultibyteValidatorName),
			withValidationChainType(validatorType),
			withValidationErr(vErr),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}
