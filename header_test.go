package ginvalidator

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestHeaderValidationChain(t *testing.T) {
	header := http.Header{
		"PHPSESSID":         {"f25g9kvjlou432vmc0ht"},
		"JSESSIONID":        {"D4E4B8CD58F4B5205E013B0B4467D5DF"},
		"ASP.NET_SessionId": {"aspxJQRF2Z1j40N5oFYVtGye"},
		"auth_token":        {"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."},
		"access_token":      {"ya29.A0AfH6SMB..."},
		"_ga":               {"GA1.2.123456789.1621593405"},
		"_gid":              {"GA1.2.987654321.1621593405"},
		"_gat":              {"1"},
		"_fat":              {"2"},
		"_fbp":              {"fb.1.1621593405.123456789"},
		"_hjid":             {"1f341234-56fe-48a2-a112-3e7cf1134567"},
		"user_pref":         {"theme=dark&lang=en"},
		"currency":          {"USD"},
		"locale":            {"en_US"},
		"csrf_token":        {"A1B2C3D4E5F6G7H8I9"},
		"XSRF-TOKEN":        {"abc123xyz456"},
		"IDE":               {"A9E9fb9tY2H48S"},
		"fr":                {"0aX7v9nZ7EfLXN"},
		"cart_id":           {"b10a8db164e0754105b7a99be72e3fe5"},
		"cart_total":        {"25.99"},
	}

	tests := []struct {
		name    string
		method  string
		url     string
		headers http.Header

		customValidatorsChain []gin.HandlerFunc
		validationResult      []ValidationChainError
		validationResultErr   error
		matchedData           MatchedData
		matchedDataErr        error
	}{
		{
			name:    "(Test Validator)(pass).",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("PHPSESSID", nil).Chain().Ascii().Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Validator(fail).",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("PHPSESSID", nil).Chain().Numeric(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: "PHPSESSID", Value: "f25g9kvjlou432vmc0ht"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Validator(pass) multiple fields.",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("PHPSESSID", nil).Chain().Ascii().Validate(),
				NewHeader("JSESSIONID", nil).Chain().Ascii().Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht", "JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Validator(fail) multiple fields.",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("PHPSESSID", nil).Chain().Numeric(nil).Validate(),
				NewHeader("JSESSIONID", nil).Chain().Currency(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: "PHPSESSID", Value: "f25g9kvjlou432vmc0ht"},
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: "JSESSIONID", Value: "D4E4B8CD58F4B5205E013B0B4467D5DF"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht", "JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Validator(fail) multiple fields, with error message.",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("PHPSESSID", func(initialValue, sanitizedValue, validatorName string) string {
					return validatorName
				}).Chain().Numeric(nil).Validate(),
				NewHeader("JSESSIONID", func(initialValue, sanitizedValue, validatorName string) string {
					return validatorName
				}).Chain().Currency(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "headers", Msg: "Numeric", Field: "PHPSESSID", Value: "f25g9kvjlou432vmc0ht"},
				{Location: "headers", Msg: "Currency", Field: "JSESSIONID", Value: "D4E4B8CD58F4B5205E013B0B4467D5DF"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht", "JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Custom Validator(pass).",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("PHPSESSID", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool { return true },
				).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Custom Validator(fail).",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("PHPSESSID", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: "PHPSESSID", Value: "f25g9kvjlou432vmc0ht"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Custom Validator(pass) multiple fields.",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("PHPSESSID", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).Validate(),
				NewHeader("JSESSIONID", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht", "JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Custom Validator(fail) multiple fields.",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("PHPSESSID", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).Validate(),
				NewHeader("JSESSIONID", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: "PHPSESSID", Value: "f25g9kvjlou432vmc0ht"},
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: "JSESSIONID", Value: "D4E4B8CD58F4B5205E013B0B4467D5DF"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht", "JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Sanitizer.",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("locale", nil).Chain().Whitelist("a-z").Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"locale": "en"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Sanitizer multiple fields.",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("locale", nil).Chain().Whitelist("a-z").Validate(),
				NewHeader("currency", nil).Chain().Whitelist("a-z").Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"locale": "en", "currency": ""}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Custom Sanitizer.",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("JSESSIONID", nil).Chain().CustomSanitizer(
					func(req *http.Request, initialValue, sanitizedValue string) string {
						return "custom-sanitizer"
					},
				).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"JSESSIONID": "custom-sanitizer"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Custom Sanitizer multiple fields.",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("PHPSESSID", nil).Chain().CustomSanitizer(
					func(req *http.Request, initialValue, sanitizedValue string) string {
						return "f25g9kvjlou432vmc0ht"
					},
				).Validate(),
				NewHeader("JSESSIONID", nil).Chain().CustomSanitizer(
					func(req *http.Request, initialValue, sanitizedValue string) string {
						return "D4E4B8CD58F4B5205E013B0B4467D5DF"
					},
				).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht", "JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Modifier(bail).",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("JSESSIONID", nil).Chain().Alpha(nil).Bail().LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: "JSESSIONID", Value: "D4E4B8CD58F4B5205E013B0B4467D5DF"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Modifier(bail) multiple fields.",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("JSESSIONID", nil).Chain().Alpha(nil).Bail().LowerCase().Validate(),
				NewHeader("_fat", nil).Chain().Alpha(nil).Bail().LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: "JSESSIONID", Value: "D4E4B8CD58F4B5205E013B0B4467D5DF"},
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: "_fat", Value: "2"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF", "_fat": "2"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Modifier(if/bail).",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("PHPSESSID", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: "PHPSESSID", Value: "f25g9kvjlou432vmc0ht"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Modifier(if/bail) multiple fields.",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("auth_token", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).LowerCase().Validate(),
				NewHeader("_gat", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: "auth_token", Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."},
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: "_gat", Value: "1"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"auth_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...", "_gat": "1"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Modifier(if/proceed).",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("fr", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: "fr", Value: "0aX7v9nZ7EfLXN"},
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: "fr", Value: "0aX7v9nZ7EfLXN"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"fr": "0aX7v9nZ7EfLXN"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Modifier(if/proceed) multiple fields.",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("csrf_token", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).LowerCase().Validate(),
				NewHeader("_gat", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: "csrf_token", Value: "A1B2C3D4E5F6G7H8I9"},
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: "csrf_token", Value: "A1B2C3D4E5F6G7H8I9"},
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: "_gat", Value: "1"},
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: "_gat", Value: "1"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"csrf_token": "A1B2C3D4E5F6G7H8I9", "_gat": "1"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Not(false -> true).",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("_gat", nil).Chain().Not().Alpha(nil).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"_gat": "1"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Not(false -> true) multiple fields.",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("_gat", nil).Chain().Not().Alpha(nil).Validate(),
				NewHeader("_fat", nil).Chain().Not().Alpha(nil).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"_gat": "1", "_fat": "2"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Not(true -> false).",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("_gat", nil).Chain().Not().Numeric(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: "_gat", Value: "1"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"_gat": "1"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Not(true -> false) multiple fields.",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("_gat", nil).Chain().Not().Numeric(nil).Validate(),
				NewHeader("_fat", nil).Chain().Not().Numeric(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: "_gat", Value: "1"},
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: "_fat", Value: "2"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"_gat": "1", "_fat": "2"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Skip(true).",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader(`auth_token`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).Numeric(nil).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{`auth_token`: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Skip(true) multiple fields.",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader(`currency`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).Numeric(nil).Validate(),
				NewHeader(`_gat`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).Numeric(nil).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{`currency`: "USD", "_gat": "1"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Skip(false).",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader(`currency`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).Numeric(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: `currency`, Value: "USD"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{`currency`: "USD"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Skip(true) just one.",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader(`auth_token`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).Numeric(nil).Currency(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: `auth_token`, Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{`auth_token`: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Skip(false) just one.",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader(`access_token`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).Numeric(nil).Currency(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: `access_token`, Value: "ya29.A0AfH6SMB..."},
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: `access_token`, Value: "ya29.A0AfH6SMB..."},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{`access_token`: "ya29.A0AfH6SMB..."}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Optional(present)",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("_gat", nil).Chain().Alpha(nil).Optional().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "headers", Msg: DefaultValChainErrMsg, Field: "_gat", Value: "1"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"_gat": "1"}},
			matchedDataErr:      nil,
		},
		{
			name:    "Test Optional(not-present)",
			method:  "POST",
			url:     "/test",
			headers: header,
			customValidatorsChain: []gin.HandlerFunc{
				NewHeader("tax_id", nil).Chain().Alpha(nil).Optional().Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"headers": MatchedDataFieldValues{"tax_id": ""}},
			matchedDataErr:      nil,
		},
	}

	for _, test := range tests {
		// test := test
		t.Run(test.name, func(t *testing.T) {
			// t.Parallel()
			router := setupRouter()

			assertHandler := func(ctx *gin.Context) {
				validationResult, err := ValidationResult(ctx)

				if test.validationResultErr != nil {
					if !errors.Is(test.validationResultErr, err) {
						t.Errorf("got %+v, wanted %+v", err, test.validationResultErr)
					}
				} else {
					if !cmp.Equal(test.validationResult, validationResult, cmpopts.IgnoreUnexported(ValidationChainError{}), cmpopts.EquateEmpty()) {
						t.Errorf("got %+v, wanted %+v", validationResult, test.validationResult)
					}
				}

				matchedData, err := GetMatchedData(ctx)

				if test.matchedDataErr != nil {
					if !errors.Is(test.matchedDataErr, err) {
						t.Errorf("got error %v, wanted error %v", err, test.matchedDataErr)
					}
				} else {
					if !reflect.DeepEqual(test.matchedData, matchedData) {
						t.Errorf("got map %+v, wanted map %+v", matchedData, test.matchedData)
					}
				}
			}

			if test.method == "GET" {
				router.GET(test.url, append(test.customValidatorsChain, assertHandler)...)
			} else if test.method == "POST" {
				router.POST(test.url, append(test.customValidatorsChain, assertHandler)...)
			} else if test.method == "PUT" {
				router.PUT(test.url, append(test.customValidatorsChain, assertHandler)...)
			} else if test.method == "PATCH" {
				router.PATCH(test.url, append(test.customValidatorsChain, assertHandler)...)
			} else {
				t.Errorf("invalid http request method: %s", test.method)
			}

			req, _ := http.NewRequest(test.method, test.url, strings.NewReader(""))
			req.Header = test.headers

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		})
	}
}
