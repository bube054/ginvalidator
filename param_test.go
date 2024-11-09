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

func TestParamValidationChain(t *testing.T) {

	params := gin.Params{
		gin.Param{Key: "PHPSESSID", Value: "f25g9kvjlou432vmc0ht"},
		gin.Param{Key: "JSESSIONID", Value: "D4E4B8CD58F4B5205E013B0B4467D5DF"},
		gin.Param{Key: "ASP.NET_SessionId", Value: "aspxJQRF2Z1j40N5oFYVtGye"},
		gin.Param{Key: "auth_token", Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."},
		gin.Param{Key: "access_token", Value: "ya29.A0AfH6SMB..."},
		gin.Param{Key: "_ga", Value: "GA1.2.123456789.1621593405"},
		gin.Param{Key: "_gid", Value: "GA1.2.987654321.1621593405"},
		gin.Param{Key: "_gat", Value: "1"},
		gin.Param{Key: "_fat", Value: "2"},
		gin.Param{Key: "_fbp", Value: "fb.1.1621593405.123456789"},
		gin.Param{Key: "_hjid", Value: "1f341234-56fe-48a2-a112-3e7cf1134567"},
		gin.Param{Key: "user_pref", Value: "theme=dark&lang=en"},
		gin.Param{Key: "currency", Value: "USD"},
		gin.Param{Key: "locale", Value: "en_US"},
		gin.Param{Key: "csrf_token", Value: "A1B2C3D4E5F6G7H8I9"},
		gin.Param{Key: "XSRF-TOKEN", Value: "abc123xyz456"},
		gin.Param{Key: "IDE", Value: "A9E9fb9tY2H48S"},
		gin.Param{Key: "fr", Value: "0aX7v9nZ7EfLXN"},
		gin.Param{Key: "cart_id", Value: "b10a8db164e0754105b7a99be72e3fe5"},
		gin.Param{Key: "cart_total", Value: "25.99"},
	}

	tests := []struct {
		name   string
		method string
		url    string
		params gin.Params

		customValidatorsChain []gin.HandlerFunc
		validationResult      []ValidationChainError
		validationResultErr   error
		matchedData           MatchedData
		matchedDataErr        error
	}{
		{
			name:   "(Test Validator)(pass).",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("PHPSESSID", nil).Chain().Ascii().Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Validator(fail).",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("PHPSESSID", nil).Chain().Numeric(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "params", Msg: defaultValChainErrMsg, Field: "PHPSESSID", Value: "f25g9kvjlou432vmc0ht"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Validator(pass) multiple fields.",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("PHPSESSID", nil).Chain().Ascii().Validate(),
				NewParam("JSESSIONID", nil).Chain().Ascii().Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht", "JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Validator(fail) multiple fields.",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("PHPSESSID", nil).Chain().Numeric(nil).Validate(),
				NewParam("JSESSIONID", nil).Chain().Currency(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "params", Msg: defaultValChainErrMsg, Field: "PHPSESSID", Value: "f25g9kvjlou432vmc0ht"},
				{Location: "params", Msg: defaultValChainErrMsg, Field: "JSESSIONID", Value: "D4E4B8CD58F4B5205E013B0B4467D5DF"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht", "JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Validator(fail) multiple fields, with error message.",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("PHPSESSID", func(initialValue, sanitizedValue, validatorName string) string {
					return validatorName
				}).Chain().Numeric(nil).Validate(),
				NewParam("JSESSIONID", func(initialValue, sanitizedValue, validatorName string) string {
					return validatorName
				}).Chain().Currency(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "params", Msg: "Numeric", Field: "PHPSESSID", Value: "f25g9kvjlou432vmc0ht"},
				{Location: "params", Msg: "Currency", Field: "JSESSIONID", Value: "D4E4B8CD58F4B5205E013B0B4467D5DF"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht", "JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Custom Validator(pass).",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("PHPSESSID", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool { return true },
				).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Custom Validator(fail).",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("PHPSESSID", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "params", Msg: defaultValChainErrMsg, Field: "PHPSESSID", Value: "f25g9kvjlou432vmc0ht"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Custom Validator(pass) multiple fields.",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("PHPSESSID", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).Validate(),
				NewParam("JSESSIONID", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht", "JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Custom Validator(fail) multiple fields.",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("PHPSESSID", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).Validate(),
				NewParam("JSESSIONID", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "params", Msg: defaultValChainErrMsg, Field: "PHPSESSID", Value: "f25g9kvjlou432vmc0ht"},
				{Location: "params", Msg: defaultValChainErrMsg, Field: "JSESSIONID", Value: "D4E4B8CD58F4B5205E013B0B4467D5DF"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht", "JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Sanitizer.",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("locale", nil).Chain().Whitelist("a-z").Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"locale": "en"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Sanitizer multiple fields.",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("locale", nil).Chain().Whitelist("a-z").Validate(),
				NewParam("currency", nil).Chain().Whitelist("a-z").Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"locale": "en", "currency": ""}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Custom Sanitizer.",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("JSESSIONID", nil).Chain().CustomSanitizer(
					func(req *http.Request, initialValue, sanitizedValue string) string {
						return "custom-sanitizer"
					},
				).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"JSESSIONID": "custom-sanitizer"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Custom Sanitizer multiple fields.",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("PHPSESSID", nil).Chain().CustomSanitizer(
					func(req *http.Request, initialValue, sanitizedValue string) string {
						return "f25g9kvjlou432vmc0ht"
					},
				).Validate(),
				NewParam("JSESSIONID", nil).Chain().CustomSanitizer(
					func(req *http.Request, initialValue, sanitizedValue string) string {
						return "D4E4B8CD58F4B5205E013B0B4467D5DF"
					},
				).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht", "JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Modifier(bail).",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("JSESSIONID", nil).Chain().Alpha(nil).Bail().LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "params", Msg: defaultValChainErrMsg, Field: "JSESSIONID", Value: "D4E4B8CD58F4B5205E013B0B4467D5DF"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Modifier(bail) multiple fields.",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("JSESSIONID", nil).Chain().Alpha(nil).Bail().LowerCase().Validate(),
				NewParam("_fat", nil).Chain().Alpha(nil).Bail().LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "params", Msg: defaultValChainErrMsg, Field: "JSESSIONID", Value: "D4E4B8CD58F4B5205E013B0B4467D5DF"},
				{Location: "params", Msg: defaultValChainErrMsg, Field: "_fat", Value: "2"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF", "_fat": "2"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Modifier(if/bail).",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("PHPSESSID", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "params", Msg: defaultValChainErrMsg, Field: "PHPSESSID", Value: "f25g9kvjlou432vmc0ht"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Modifier(if/bail) multiple fields.",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("auth_token", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).LowerCase().Validate(),
				NewParam("_gat", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "params", Msg: defaultValChainErrMsg, Field: "auth_token", Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."},
				{Location: "params", Msg: defaultValChainErrMsg, Field: "_gat", Value: "1"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"auth_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...", "_gat": "1"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Modifier(if/proceed).",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("fr", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "params", Msg: defaultValChainErrMsg, Field: "fr", Value: "0aX7v9nZ7EfLXN"},
				{Location: "params", Msg: defaultValChainErrMsg, Field: "fr", Value: "0aX7v9nZ7EfLXN"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"fr": "0aX7v9nZ7EfLXN"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Modifier(if/proceed) multiple fields.",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("csrf_token", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).LowerCase().Validate(),
				NewParam("_gat", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "params", Msg: defaultValChainErrMsg, Field: "csrf_token", Value: "A1B2C3D4E5F6G7H8I9"},
				{Location: "params", Msg: defaultValChainErrMsg, Field: "csrf_token", Value: "A1B2C3D4E5F6G7H8I9"},
				{Location: "params", Msg: defaultValChainErrMsg, Field: "_gat", Value: "1"},
				{Location: "params", Msg: defaultValChainErrMsg, Field: "_gat", Value: "1"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"csrf_token": "A1B2C3D4E5F6G7H8I9", "_gat": "1"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Not(false -> true).",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("_gat", nil).Chain().Not().Alpha(nil).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"_gat": "1"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Not(false -> true) multiple fields.",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("_gat", nil).Chain().Not().Alpha(nil).Validate(),
				NewParam("_fat", nil).Chain().Not().Alpha(nil).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"_gat": "1", "_fat": "2"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Not(true -> false).",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("_gat", nil).Chain().Not().Numeric(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "params", Msg: defaultValChainErrMsg, Field: "_gat", Value: "1"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"_gat": "1"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Not(true -> false) multiple fields.",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("_gat", nil).Chain().Not().Numeric(nil).Validate(),
				NewParam("_fat", nil).Chain().Not().Numeric(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "params", Msg: defaultValChainErrMsg, Field: "_gat", Value: "1"},
				{Location: "params", Msg: defaultValChainErrMsg, Field: "_fat", Value: "2"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"_gat": "1", "_fat": "2"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Skip(true).",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam(`auth_token`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).Numeric(nil).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{`auth_token`: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Skip(true) multiple fields.",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam(`currency`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).Numeric(nil).Validate(),
				NewParam(`_gat`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).Numeric(nil).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{`currency`: "USD", "_gat": "1"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Skip(false).",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam(`currency`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).Numeric(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "params", Msg: defaultValChainErrMsg, Field: `currency`, Value: "USD"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{`currency`: "USD"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Skip(true) just one.",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam(`auth_token`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).Numeric(nil).Currency(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "params", Msg: defaultValChainErrMsg, Field: `auth_token`, Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{`auth_token`: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Skip(false) just one.",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam(`access_token`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).Numeric(nil).Currency(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "params", Msg: defaultValChainErrMsg, Field: `access_token`, Value: "ya29.A0AfH6SMB..."},
				{Location: "params", Msg: defaultValChainErrMsg, Field: `access_token`, Value: "ya29.A0AfH6SMB..."},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{`access_token`: "ya29.A0AfH6SMB..."}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Optional(present)",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("_gat", nil).Chain().Alpha(nil).Optional().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "params", Msg: defaultValChainErrMsg, Field: "_gat", Value: "1"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"_gat": "1"}},
			matchedDataErr:      nil,
		},
		{
			name:   "Test Optional(not-present)",
			method: "POST",
			url:    "/test",
			params: params,
			customValidatorsChain: []gin.HandlerFunc{
				NewParam("tax_id", nil).Chain().Alpha(nil).Optional().Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"params": matchedDataFieldValues{"tax_id": ""}},
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

			dynamicPath, finalURL := generateDynamicURL(test.params)

			if test.method == "GET" {
				router.GET(dynamicPath, append(test.customValidatorsChain, assertHandler)...)
			} else if test.method == "POST" {
				router.POST(dynamicPath, append(test.customValidatorsChain, assertHandler)...)
			} else if test.method == "PUT" {
				router.PUT(dynamicPath, append(test.customValidatorsChain, assertHandler)...)
			} else if test.method == "PATCH" {
				router.PATCH(dynamicPath, append(test.customValidatorsChain, assertHandler)...)
			} else {
				t.Errorf("invalid http request method: %s", test.method)
			}

			req, _ := http.NewRequest(test.method, finalURL, strings.NewReader(""))

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		})
	}
}

func generateDynamicURL(params gin.Params) (string, string) {
	var pathTemplate []string
	var pathValues []string

	for _, param := range params {
		pathTemplate = append(pathTemplate, ":"+param.Key)
		pathValues = append(pathValues, param.Value)
	}

	dynamicPath := "/" + strings.Join(pathTemplate, "/")
	finalURL := "/" + strings.Join(pathValues, "/")

	return dynamicPath, finalURL
}
