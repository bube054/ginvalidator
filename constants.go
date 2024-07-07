package ginvalidator

// request locations
const (
	bodyLocation    = "body"
	cookiesLocation = "cookies"
	headersLocation = "headers"
	paramsLocation  = "params"
	queryLocation   = "query"
)

// ctx keys
const (
	bodyLocationStore    = "gin___validator___body___location"
	cookiesLocationStore = "gin___validator___cookies___location"
	headersLocationStore = "gin___validator___headers___location"
	paramsLocationStore  = "gin___validator___params___location"
	queryLocationStore   = "gin___validator___query___location"
)

// method types
const (
	validatorType = "validator"
	modifierType  = "modifier"
	sanitizerType = "sanitizer"
)

// validating, converting to and working with json
const (
	jsonKey = "key"
)

// validation methods and their default error messages
const (
	customFuncName = "custom"
	customErrMsg   = ""

	existsFuncName = "exists"
	existsErrMsg   = "does not exists."

	isArrayFuncName = "isArray"
	isArrayErrMsg   = "is not an array."

	isObjectFuncName = "isObject"
	isObjectErrMsg   = "is not a object."

	isStringFuncName = "isString"
	isStringErrMsg   = "is not a string."

	isNotEmptyFuncName = "isNotEmpty"
	isNotEmptyErrMsg   = "is not empty."

	containsFuncName = "contains"
	containsErrMsg   = "does not contain."

	equalsFuncName = "equals"
	equalsErrMsg   = "does not equal."

	isAfterFuncName = "isAfter"
	isAfterErrMsg   = "is not after."

	isAlphaFuncName = "isAlpha"
	isAlphaErrMsg   = "is not alpha."

	isAlphanumericFuncName = "isAlphanumeric"
	isAlphanumericErrMsg   = "is not alphanumeric."

	isASCIIFuncName = "isASCII"
	isASCIIErrMsg   = "is not ascii."

	isBase32FuncName = "isBase32"
	isBase32ErrMsg   = "is not base32."

	isBase58FuncName = "isBase58"
	isBase58ErrMsg   = "is not base58."

	isBase64FuncName = "isBase64"
	isBase64ErrMsg   = "is not base64."

	isBeforeFuncName = "isBefore"
	isBeforeErrMsg   = "is not before."

	isBICFuncName = "isBIC"
	isBICErrMsg   = "is not bic."

	isBtcAddressFuncName = "isBtcAddress"
	isBtcAddressErrMsg   = "is not btc address."

	isByteLengthFuncName = "isByteLength"
	isByteLengthErrMsg   = "is not byte length."

	isCreditCardFuncName = "isCreditCard"
	isCreditCardErrMsg   = "is not credit card."

	isCurrencyFuncName = "isCurrency"
	isCurrencyErrMsg   = "is not currency."

	isDataURIFuncName = "isDataURI"
	isDataURIErrMsg   = "is not data uri."

	isDateFuncName = "isDate"
	isDateErrMsg   = "is not date."

	isDecimalFuncName = "isDecimal"
	isDecimalErrMsg   = "is not decimal."

	isDivisibleByFuncName = "isDivisibleBy"
	isDivisibleByErrMsg   = "is not divisible by."

	isEmailFuncName = "isEmail"
	isEmailErrMsg   = "is not email."

	isEmptyFuncName = "isEmpty"
	isEmptyErrMsg   = "is not empty."

	isEthereumAddressFuncName = "isEthereumAddress"
	isEthereumAddressErrMsg   = "is not ethereum address."

	isFQDNFuncName = "isFQDN"
	isFQDNErrMsg   = "is not fqdn."

	isFloatFuncName = "isFloat"
	isFloatErrMsg   = "is not float."

	isFullWidthFuncName = "isFullWidth"
	isFullWidthErrMsg   = "is not full width."

	isHalfWidthFuncName = "isHalfWidth"
	isHalfWidthErrMsg   = "is not half width."

	isHashFuncName = "isHash"
	isHashErrMsg   = "is not hash."

	isHexColorFuncName = "isHexColor"
	isHexColorErrMsg   = "is not hex color."

	isHexadecimalFuncName = "isHexadecimal"
	isHexadecimalErrMsg   = "is not hexadecimal."

	isHSLFuncName = "isHSL"
	isHSLErrMsg   = "is not hsl."

	isHSLAFuncName = "isHSLA"
	isHSLAErrMsg   = "is not hsla."

	isIBANFuncName = "isIBAN"
	isIBANErrMsg   = "is not iban."

	isIdentityCardFuncName = "isIdentityCard"
	isIdentityCardErrMsg   = "is not identity card."

	isIMEIFuncName = "isIMEI"
	isIMEIErrMsg   = "is not imei."

	isIPFuncName = "isIP"
	isIPErrMsg   = "is not ip."

	isIPRangeFuncName = "isIPRange"
	isIPRangeErrMsg   = "is not ip range."

	isIPv4FuncName = "isIPv4"
	isIPv4ErrMsg   = "is not ipv4."

	isISBNFuncName = "isISBN"
	isISBNErrMsg   = "is not isbn."

	isISSNFuncName = "isISSN"
	isISSNErrMsg   = "is not issn."

	isISINFuncName = "isISIN"
	isISINErrMsg   = "is not isin."

	isISO6391FuncName = "isISO6391"
	isISO6391ErrMsg   = "is not iso6391."

	isISO8601FuncName = "isISO8601"
	isISO8601ErrMsg   = "is not iso8601."

	isISO31661Alpha2FuncName = "isISO31661Alpha2"
	isISO31661Alpha2ErrMsg   = "is not iso31661 alpha2."

	isISO31661Alpha3FuncName = "isISO31661Alpha3"
	isISO31661Alpha3ErrMsg   = "is not iso31661 alpha3."

	isISO4217FuncName = "isISO4217"
	isISO4217ErrMsg   = "is not iso4217."

	isISRCFuncName = "isISRC"
	isISRCErrMsg   = "is not isrc."

	isInFuncName = "isIn"
	isInErrMsg   = "is not in."

	isIntFuncName = "isInt"
	isIntErrMsg   = "is not int."

	isJSONFuncName = "isJSON"
	isJSONErrMsg   = "is not json."

	isJWTFuncName = "isJWT"
	isJWTErrMsg   = "is not jwt."

	isLatLongFuncName = "isLatLong"
	isLatLongErrMsg   = "is not lat long."

	isLengthFuncName = "isLength"
	isLengthErrMsg   = "is not length."

	isLowercaseFuncName = "isLowercase"
	isLowercaseErrMsg   = "is not lowercase."

	isMACFuncName = "isMAC"
	isMACErrMsg   = "is not mac."

	isLuhnNumberFuncName = "isLuhn"
	isLuhnNumberErrMsg   = "is not luhn number."

	isMagnetURIFuncName = "isMagnetURI"
	isMagnetURIErrMsg   = "is not magnet uri."

	isMACAddressFuncName = "isMACAddress"
	isMACAddressErrMsg   = "is not mac address."

	isMD5FuncName = "isMD5"
	isMD5ErrMsg   = "is not md5."

	isMimeTyeFuncName = "isMime"
	isMimeTyeErrMsg   = "is not mime type."

	isMobilePhoneFuncName = "isMobilePhone"
	isMobilePhoneErrMsg   = "is not mobile phone."

	isMongoIDFuncName = "isMongoID"
	isMongoIDErrMsg   = "is not mongo id."

	isMultibyteFuncName = "isMultibyte"
	isMultibyteErrMsg   = "is not multibyte."

	isNumericFuncName = "isNumeric"
	isNumericErrMsg   = "is not numeric."

	isOctalFuncName = "isOctal"
	isOctalErrMsg   = "is not octal."

	isPassportNumberFuncName = "isPassportNumber"
	isPassportNumberErrMsg   = "is not passport number."

	isPortFuncName = "isPort"
	isPortErrMsg   = "is not port."

	isPostalCodeFuncName = "isPostalCode"
	isPostalCodeErrMsg   = "is not postal code."

	isRgbColorFuncName = "isRgbColor"
	isRgbColorErrMsg   = "is not rgb color."

	isRFC3339FuncName = "isRFC3339"
	isRFC3339ErrMsg   = "is not rfc3339."

	isSemVerFuncName = "isSemVer"
	isSemVerErrMsg   = "is not semver."

	isSlugFuncName = "isSlug"
	isSlugErrMsg   = "is not slug."

	isStrongPasswordFuncName = "isStrongPassword"
	isStrongPasswordErrMsg   = "is not strong password."

	isSurrogatePairFuncName = "isSurrogatePair"
	isSurrogatePairErrMsg   = "is not surrogate pair."

	isTaxIDFuncName = "isTaxID"
	isTaxIDErrMsg   = "is not tax id."

	isTimeFuncName = "isTime"
	isTimeErrMsg   = "is not time."

	isURLFuncName = "isURL"
	isURLErrMsg   = "is not url."

	isUUIDFuncName = "isUUID"
	isUUIDErrMsg   = "is not uuid."

	isUppercaseFuncName = "isUppercase"
	isUppercaseErrMsg   = "is not uppercase."

	isVariableWidthFuncName = "isVariableWidth"
	isVariableWidthErrMsg   = "is not variable width."

	isVATFuncName = "isVAT"
	isVATErrMsg   = "is not vat."

	isWhitelistedFuncName = "isWhitelisted"
	isWhitelistedErrMsg   = "is not whitelisted."

	matchesFuncName = "matches"
	matchesErrMsg   = "does not match."
)

// for modifiers
const (
	notFunc  = "not"
	bailFunc = "bail"
	iFFunc   = "if"
)

// for sanitizers
const (
	customSanitizerFunc = "customSanitizer"
	defaultFunc         = "default"
	replaceFunc         = "replace"
	toArrayFunc         = "toArray"
	toLowerCaseFunc     = "toLowerCase"
	toUpperCaseFunc     = "toUpperCase"
	blacklistFunc       = "blacklist"
	escapeFunc          = "escape"
	unescapeFunc        = "unescape"
	ltrimFunc           = "ltrim"
	normalizeEmailFunc  = "normalizeEmail"
	rtrimFunc           = "rtrim"
	stripLowFunc        = "stripLow"
	toBooleanFunc       = "toBoolean"
	toDateFunc          = "toDate"
	toFloatFunc         = "toFloat"
	toIntFunc           = "toInt"
	trimFunc            = "trim"
	whitelistFunc       = "whitelist"
)
