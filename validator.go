package ginvalidator

import (
	"net/http"

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
	Iso4217ValidatorName            string = "Iso4217"
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
	errFmtFunc *ErrFmtFuncHandler // the function to create the error message

	reqLoc            requestLocation  // the HTTP request location (e.g., body, headers, cookies, params, or queries)
	rulesCreatorFuncs ruleCreatorFuncs // the list of functions that creates the validation rules.
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
// [Contains]: https://pkg.go.dev/github.com/bube054/validatorgo#Contains
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
// [Equals]: https://pkg.go.dev/github.com/bube054/validatorgo#Equals
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

// newValidator creates and returns a new validator.
//
// Parameters:
//   - field: The field to validate from the HTTP request data location (e.g., body, headers, cookies, params, or queries).
//   - errFmtFunc: A function that returns a custom error message. If nil, a generic error message will be used.
//   - reqLoc: The location in the HTTP request from where the field is extracted (e.g., body, headers, cookies, params, or queries).
func newValidator(field string, errFmtFunc *ErrFmtFuncHandler, reqLoc requestLocation) validator {
	return validator{
		field:      field,
		errFmtFunc: errFmtFunc,
		reqLoc:     reqLoc,
	}
}
