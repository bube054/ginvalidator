package ginvalidator

import (
	"net/http"
	"regexp"

	vgo "github.com/bube054/validatorgo"
	"github.com/gin-gonic/gin"
)

const (
	CustomValidatorName             string = "CustomValidator"
	ContainsValidatorName           string = "Contains"
	EqualsValidatorName             string = "Equals"
	AbaRoutingValidatorName         string = "AbaRouting"
	AfterValidatorName              string = "After"
	AlphaValidatorName              string = "Alpha"
	AlphanumericValidatorName       string = "Alphanumeric"
	AsciiValidatorName              string = "Ascii"
	BTCAddressValidatorName         string = "BTCAddress"
	Base32ValidatorName             string = "Base32"
	Base58ValidatorName             string = "Base58"
	Base64ValidatorName             string = "Base64"
	BeforeValidatorName             string = "Before"
	BicValidatorName                string = "Bic"
	BooleanValidatorName            string = "Boolean"
	ByteLengthValidatorName         string = "ByteLength"
	CountryCodeValidatorName        string = "CountryCode"
	CreditCardValidatorName         string = "CreditCard"
	CurrencyValidatorName           string = "Currency"
	DataURIValidatorName            string = "DataURI"
	DateValidatorName               string = "Date"
	DecimalValidatorName            string = "Decimal"
	DivisibleByValidatorName        string = "DivisibleBy"
	EANValidatorName                string = "EAN"
	EmailValidatorName              string = "Email"
	EmptyValidatorName              string = "Empty"
	EthereumAddressValidatorName    string = "EthereumAddress"
	FQDNValidatorName               string = "FQDN"
	FloatValidatorName              string = "Float"
	FreightContainerIDValidatorName string = "FreightContainerID"
	FullWidthValidatorName          string = "FullWidth"
	HSLValidatorName                string = "HSL"
	HalfWidthValidatorName          string = "HalfWidth"
	HashValidatorName               string = "Hash"
	HexColorValidatorName           string = "HexColor"
	HexadecimalValidatorName        string = "Hexadecimal"
	IBANValidatorName               string = "IBAN"
	IMEIValidatorName               string = "IMEI"
	IPValidatorName                 string = "IP"
	IPRangeValidatorName            string = "IPRange"
	ISBNValidatorName               string = "ISBN"
	ISINValidatorName               string = "ISIN"
	ISO31661Alpha2ValidatorName     string = "ISO31661Alpha2"
	ISO31661Alpha3ValidatorName     string = "ISO31661Alpha3"
	ISO31661NumericValidatorName    string = "ISO31661Numeric"
	ISO6346ValidatorName            string = "ISO6346"
	ISO6391ValidatorName            string = "ISO6391"
	ISO8601ValidatorName            string = "ISO8601"
	ISRCValidatorName               string = "ISRC"
	ISSNValidatorName               string = "ISSN"
	IdentityCardValidatorName       string = "IdentityCard"
	InValidatorName                 string = "In"
	IntValidatorName                string = "Int"
	ISO4217ValidatorName            string = "ISO4217"
	JSONValidatorName               string = "JSON"
	JWTValidatorName                string = "JWT"
	LatLongValidatorName            string = "LatLong"
	LengthValidatorName             string = "Length"
	LicensePlateValidatorName       string = "LicensePlate"
	LocaleValidatorName             string = "Locale"
	LowerCaseValidatorName          string = "LowerCase"
	LuhnNumberValidatorName         string = "LuhnNumber"
	MD5ValidatorName                string = "MD5"
	MacAddressValidatorName         string = "MacAddress"
	MagnetURIValidatorName          string = "MagnetURI"
	MailtoURIValidatorName          string = "MailtoURI"
	MimeTypeValidatorName           string = "MimeType"
	MobilePhoneValidatorName        string = "MobilePhone"
	MongoIDValidatorName            string = "MongoID"
	MultibyteValidatorName          string = "Multibyte"
	NumericValidatorName            string = "Numeric"
	OctalValidatorName              string = "Octal"
	PassportNumberValidatorName     string = "PassportNumber"
	PortValidatorName               string = "Port"
	PostalCodeValidatorName         string = "PostalCode"
	RFC3339ValidatorName            string = "RFC3339"
	RgbColorValidatorName           string = "RgbColor"
	SemVerValidatorName             string = "SemVer"
	SlugValidatorName               string = "Slug"
	StrongPasswordValidatorName     string = "StrongPassword"
	SurrogatePairValidatorName      string = "SurrogatePair"
	TaxIDValidatorName              string = "TaxID"
	TimeValidatorName               string = "Time"
	ULIDValidatorName               string = "ULID"
	URLValidatorName                string = "URL"
	UUIDValidatorName               string = "UUID"
	UpperCaseValidatorName          string = "UpperCase"
	VATValidatorName                string = "VAT"
	VariableWidthValidatorName      string = "VariableWidth"
	WhitelistedValidatorName        string = "Whitelisted"
	MatchesValidatorName            string = "Matches"
)

// A validator is simply a piece of the validation chain that can validate values from the specified field.
type validator struct {
	field      string             // the field to be specified
	errFmtFunc ErrFmtFuncHandler // the function to create the error message

	reqLoc            requestLocation  // the HTTP request location (e.g., body, headers, cookies, params, or queries)
	rulesCreatorFuncs ruleCreatorFuncs // the list of functions that creates the validation rules.
}

// newValidator creates and returns a new validator.
//
// Parameters:
//   - field: The field to validate from the HTTP request data location (e.g., body, headers, cookies, params, or queries).
//   - errFmtFunc: A function that returns a custom error message. If nil, a generic error message will be used.
//   - reqLoc: The location in the HTTP request from where the field is extracted (e.g., body, headers, cookies, params, or queries).
func newValidator(field string, errFmtFunc ErrFmtFuncHandler, reqLoc requestLocation) validator {
	return validator{
		field:      field,
		errFmtFunc: errFmtFunc,
		reqLoc:     reqLoc,
	}
}

// recreateValidationChainFromValidator takes the previous validator and returns a new validation chain.
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

// CustomValidatorFunc defines a function that evaluates whether the value is valid according to your custom logic.
//
// Parameters:
//   - req: The HTTP request context derived from `http.Request`.
//   - initialValue: The original value derived from the specified field.
//   - sanitizedValue: The current sanitized value after applying previous sanitizers.
type CustomValidatorFunc func(req http.Request, initialValue, sanitizedValue string) bool

// CustomValidator applies a custom validator function.
//
// Parameters:
//   - cvf: The [CustomValidatorFunc] used to evaluate the validity.
func (v validator) CustomValidator(cvf CustomValidatorFunc) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		httpRequest := ctx.Request
		isValid := cvf(*httpRequest, initialValue, sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(CustomValidatorName),
			withValidationChainType(validatorType),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Contains is a validator that checks if the string contains the seed.
//
// This function uses the [Contains] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
//
// [IsContains]: https://pkg.go.dev/github.com/bube054/validatorgo#Contains
func (v validator) Contains(seed string, opts *vgo.ContainsOpt) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid := vgo.Contains(sanitizedValue, seed, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ContainsValidatorName),
			withValidationChainType(validatorType),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Equals is a validator that checks if the string contains the seed.
//
// This function uses the [Equals] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
//
// [IsEquals]: https://pkg.go.dev/github.com/bube054/validatorgo#Equals
func (v validator) Equals(comparison string) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid := vgo.Equals(sanitizedValue, comparison)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(EqualsValidatorName),
			withValidationChainType(validatorType),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// AbaRouting is a validator that checks if the string is an ABA routing number for US bank account / cheque.
//
// This function uses the [IsAbaRouting] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsAbaRouting]: https://pkg.go.dev/github.com/bube054/validatorgo#IsAbaRouting
func (v validator) AbaRouting() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid := vgo.IsAbaRouting(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(AbaRoutingValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsAfter(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(AfterValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsAlpha(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(AlphaValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsAlphanumeric(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(AlphanumericValidatorName),
			withValidationChainType(validatorType),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Ascii validator that checks if the string contains ASCII chars only.
//
// This function uses the [IsAscii] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsAscii]: https://pkg.go.dev/github.com/bube054/validatorgo#IsAscii
func (v validator) Ascii() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid := vgo.IsAscii(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(AbaRoutingValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsBase32(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(Base32ValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsBase58(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(Base58ValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsBase64(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(Base64ValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsBefore(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(BeforeValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsBic(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(BicValidatorName),
			withValidationChainType(validatorType),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Boolean validator that check if the string is a boolean.
//
// This function uses the [IsBoolean] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsBoolean]: https://pkg.go.dev/github.com/bube054/validatorgo#IsBoolean
func (v validator) Boolean(opts *vgo.IsBooleanOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid := vgo.IsBoolean(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(BooleanValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsBTCAddress(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(BTCAddressValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsByteLength(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ByteLengthValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsCreditCard(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(CreditCardValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsCurrency(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(CurrencyValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsDataURI(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(DataURIValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsDate(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(DataURIValidatorName),
			withValidationChainType(validatorType),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Decimal is a validator that check if the string represents a decimal number, such as 0.1, .3, 1.1, 1.00003, 4.0, etc.
//
// This function uses the [IsDecimal] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsDecimal]: https://pkg.go.dev/github.com/bube054/validatorgo#IsDecimal
func (v validator) Decimal(opts *vgo.IsDecimalOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid := vgo.IsDecimal(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(DecimalValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsDivisibleBy(sanitizedValue, num)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(DivisibleByValidatorName),
			withValidationChainType(validatorType),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// EAN is validator that checks if the string is a valid EAN (European Article Number).
//
// This function uses the [IsEAN] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsEAN]: https://pkg.go.dev/github.com/bube054/validatorgo#IsEAN
func (v validator) EAN() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid := vgo.IsEAN(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(EANValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsEmail(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(EmailValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsEmpty(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(EmptyValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsEthereumAddress(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(EthereumAddressValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsFloat(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(FloatValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsFQDN(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(FQDNValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsFreightContainerID(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(FreightContainerIDValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsFullWidth(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(FullWidthValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsHalfWidth(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(HalfWidthValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsHash(sanitizedValue, algorithm)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(HashValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsHexadecimal(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(HexadecimalValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsHexColor(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(HexColorValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsHSL(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(HSLValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsIBAN(sanitizedValue, countryCode)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(IBANValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsIdentityCard(sanitizedValue, locale)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(IdentityCardValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsIMEI(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(IMEIValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsIn(sanitizedValue, values)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(InValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsInt(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(IntValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsIP(sanitizedValue, version)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(IPValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsIPRange(sanitizedValue, version)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(IPRangeValidatorName),
			withValidationChainType(validatorType),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// ISIN is a validator that checks if the string is an ISIN (stock/security identifier).
//
// This function uses the [IsISIN] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsISIN]: https://pkg.go.dev/github.com/bube054/validatorgo#IsISIN
func (v validator) ISIN() ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid := vgo.IsISIN(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(InValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsIso4217(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ISO4217ValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsISO6346(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ISO6346ValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsISO6391(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ISO6391ValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsISO8601(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ISO8601ValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsISO31661Alpha2(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ISO31661Alpha2ValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsISO31661Alpha3(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ISO31661Alpha3ValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsISO31661Numeric(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ISO31661NumericValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsISRC(sanitizedValue, allowHyphens)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ISRCValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsISSN(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ISSNValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsJSON(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(JSONValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsLatLong(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(LatLongValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsLicensePlate(sanitizedValue, locale)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(LicensePlateValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsLocale(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(LocaleValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsLowerCase(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(LowerCaseValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsLuhnNumber(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(LuhnNumberValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsMacAddress(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(MacAddressValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsMagnetURI(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(MagnetURIValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsMailtoURI(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(MailtoURIValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsMD5(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(MD5ValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsMimeType(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(MimeTypeValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsMobilePhone(sanitizedValue, locales, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(MobilePhoneValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsMongoID(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(MongoIDValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsMultibyte(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(MultibyteValidatorName),
			withValidationChainType(validatorType),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}

// Numeric is a validator that check if a string is a number.
//
// This function uses the [IsNumeric] from [validatorgo] package to perform the validation logic.
//
// [validatorgo]: https://pkg.go.dev/github.com/bube054
// [IsNumeric]: https://pkg.go.dev/github.com/bube054/validatorgo#IsNumeric
func (v validator) Numeric(opts *vgo.IsNumericOpts) ValidationChain {
	var ruleCreator ruleCreatorFunc = func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule {
		isValid := vgo.IsNumeric(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(NumericValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsOctal(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(OctalValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsPassportNumber(sanitizedValue, countryCode)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(PassportNumberValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsPort(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(PortValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsPostalCode(sanitizedValue, locale)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(PostalCodeValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsRFC3339(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(RFC3339ValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsRgbColor(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(RgbColorValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsSemVer(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(SemVerValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsSlug(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(SlugValidatorName),
			withValidationChainType(validatorType),
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
		isValid, _ := vgo.IsStrongPassword(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(StrongPasswordValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsTaxID(sanitizedValue, locale)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(TaxIDValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsSurrogatePair(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(SurrogatePairValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsTime(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(TimeValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsULID(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(ULIDValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsUpperCase(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(UpperCaseValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsURL(sanitizedValue, opts)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(URLValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsUUID(sanitizedValue, version)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(UUIDValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsVariableWidth(sanitizedValue)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(VariableWidthValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsVAT(sanitizedValue, countryCode)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(VATValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.IsWhitelisted(sanitizedValue, chars)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(WhitelistedValidatorName),
			withValidationChainType(validatorType),
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
		isValid := vgo.Matches(sanitizedValue, re)

		return NewValidationChainRule(
			withIsValid(isValid),
			withNewValue(sanitizedValue),
			withValidationChainName(MatchesValidatorName),
			withValidationChainType(validatorType),
		)
	}

	return v.recreateValidationChainFromValidator(ruleCreator)
}
