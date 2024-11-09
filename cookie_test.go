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

func TestCookieValidationChain(t *testing.T) {
	cookies := []*http.Cookie{
		{Name: "PHPSESSID", Value: "f25g9kvjlou432vmc0ht"},
		{Name: "JSESSIONID", Value: "D4E4B8CD58F4B5205E013B0B4467D5DF"},
		{Name: "ASP.NET_SessionId", Value: "aspxJQRF2Z1j40N5oFYVtGye"},
		{Name: "auth_token", Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."},
		{Name: "access_token", Value: "ya29.A0AfH6SMB..."},
		{Name: "_ga", Value: "GA1.2.123456789.1621593405"},
		{Name: "_gid", Value: "GA1.2.987654321.1621593405"},
		{Name: "_gat", Value: "1"},
		{Name: "_fat", Value: "2"},
		{Name: "_fbp", Value: "fb.1.1621593405.123456789"},
		{Name: "_hjid", Value: "1f341234-56fe-48a2-a112-3e7cf1134567"},
		{Name: "user_pref", Value: "theme=dark&lang=en"},
		{Name: "currency", Value: "USD"},
		{Name: "locale", Value: "en_US"},
		{Name: "csrf_token", Value: "A1B2C3D4E5F6G7H8I9"},
		{Name: "XSRF-TOKEN", Value: "abc123xyz456"},
		{Name: "IDE", Value: "A9E9fb9tY2H48S"},
		{Name: "fr", Value: "0aX7v9nZ7EfLXN"},
		{Name: "cart_id", Value: "b10a8db164e0754105b7a99be72e3fe5"},
		{Name: "cart_total", Value: "25.99"},
	}

	tests := []struct {
		name        string
		method      string
		url         string
		cookies     []*http.Cookie
		contentType string

		customValidatorsChain []gin.HandlerFunc
		validationResult      []ValidationChainError
		validationResultErr   error
		matchedData           MatchedData
		matchedDataErr        error
	}{
		// For "application/json"
		{
			name:        "Test Validator(pass).",
			method:      "POST",
			url:         "/test",
			cookies:     cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("PHPSESSID", nil).Chain().Ascii().Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Validator(fail).",
			method:      "POST",
			url:         "/test",
			cookies:     cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("PHPSESSID", nil).Chain().Numeric(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: "PHPSESSID", Value: "f25g9kvjlou432vmc0ht"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Validator(pass) multiple fields.",
			method:      "POST",
			url:         "/test",
			cookies:     cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("PHPSESSID", nil).Chain().Ascii().Validate(),
				NewCookie("JSESSIONID", nil).Chain().Ascii().Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht", "JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Validator(fail) multiple fields.",
			method:      "POST",
			url:         "/test",
			cookies:     cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("PHPSESSID", nil).Chain().Numeric(nil).Validate(),
				NewCookie("JSESSIONID", nil).Chain().Currency(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: "PHPSESSID", Value: "f25g9kvjlou432vmc0ht"},
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: "JSESSIONID", Value: "D4E4B8CD58F4B5205E013B0B4467D5DF"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht", "JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Validator(fail) multiple fields, with error message.",
			method:      "POST",
			url:         "/test",
			cookies:     cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("PHPSESSID", func(initialValue, sanitizedValue, validatorName string) string {
					return validatorName
				}).Chain().Numeric(nil).Validate(),
				NewCookie("JSESSIONID", func(initialValue, sanitizedValue, validatorName string) string {
					return validatorName
				}).Chain().Currency(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "cookies", Msg: "Numeric", Field: "PHPSESSID", Value: "f25g9kvjlou432vmc0ht"},
				{Location: "cookies", Msg: "Currency", Field: "JSESSIONID", Value: "D4E4B8CD58F4B5205E013B0B4467D5DF"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht", "JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Custom Validator(pass).",
			method:      "POST",
			url:         "/test",
			cookies:     cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("PHPSESSID", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool { return true },
				).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Custom Validator(fail).",
			method:      "POST",
			url:         "/test",
			cookies:     cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("PHPSESSID", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: "PHPSESSID", Value: "f25g9kvjlou432vmc0ht"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Custom Validator(pass) multiple fields.",
			method:      "POST",
			url:         "/test",
			cookies:     cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("PHPSESSID", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).Validate(),
				NewCookie("JSESSIONID", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht", "JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Custom Validator(fail) multiple fields.",
			method:      "POST",
			url:         "/test",
			cookies:     cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("PHPSESSID", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).Validate(),
				NewCookie("JSESSIONID", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: "PHPSESSID", Value: "f25g9kvjlou432vmc0ht"},
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: "JSESSIONID", Value: "D4E4B8CD58F4B5205E013B0B4467D5DF"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht", "JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Sanitizer.",
			method:      "POST",
			url:         "/test",
			cookies:        cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("locale", nil).Chain().Whitelist("a-z").Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"locale": "en"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Sanitizer multiple fields.",
			method:      "POST",
			url:         "/test",
			cookies:        cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("locale", nil).Chain().Whitelist("a-z").Validate(),
				NewCookie("currency", nil).Chain().Whitelist("a-z").Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"locale": "en", "currency": ""}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Custom Sanitizer.",
			method:      "POST",
			url:         "/test",
			cookies:        cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("JSESSIONID", nil).Chain().CustomSanitizer(
					func(req *http.Request, initialValue, sanitizedValue string) string {
						return "custom-sanitizer"
					},
				).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"JSESSIONID": "custom-sanitizer"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Custom Sanitizer multiple fields.",
			method:      "POST",
			url:         "/test",
			cookies:        cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("PHPSESSID", nil).Chain().CustomSanitizer(
					func(req *http.Request, initialValue, sanitizedValue string) string {
						return "f25g9kvjlou432vmc0ht"
					},
				).Validate(),
				NewCookie("JSESSIONID", nil).Chain().CustomSanitizer(
					func(req *http.Request, initialValue, sanitizedValue string) string {
						return "D4E4B8CD58F4B5205E013B0B4467D5DF"
					},
				).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht", "JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Modifier(bail).",
			method:      "POST",
			url:         "/test",
			cookies:        cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("JSESSIONID", nil).Chain().Alpha(nil).Bail().LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: "JSESSIONID", Value: "D4E4B8CD58F4B5205E013B0B4467D5DF"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Modifier(bail) multiple fields.",
			method:      "POST",
			url:         "/test",
			cookies:        cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("JSESSIONID", nil).Chain().Alpha(nil).Bail().LowerCase().Validate(),
				NewCookie("_fat", nil).Chain().Alpha(nil).Bail().LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: "JSESSIONID", Value: "D4E4B8CD58F4B5205E013B0B4467D5DF"},
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: "_fat", Value: "2"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"JSESSIONID": "D4E4B8CD58F4B5205E013B0B4467D5DF", "_fat": "2"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Modifier(if/bail).",
			method:      "POST",
			url:         "/test",
			cookies:        cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("PHPSESSID", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: "PHPSESSID", Value: "f25g9kvjlou432vmc0ht"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"PHPSESSID": "f25g9kvjlou432vmc0ht"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Modifier(if/bail) multiple fields.",
			method:      "POST",
			url:         "/test",
			cookies:        cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("auth_token", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).LowerCase().Validate(),
				NewCookie("_gat", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: "auth_token", Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."},
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: "_gat", Value: "1"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"auth_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...", "_gat": "1"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Modifier(if/proceed).",
			method:      "POST",
			url:         "/test",
			cookies:        cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("fr", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: "fr", Value: "0aX7v9nZ7EfLXN"},
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: "fr", Value: "0aX7v9nZ7EfLXN"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"fr": "0aX7v9nZ7EfLXN"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Modifier(if/proceed) multiple fields.",
			method:      "POST",
			url:         "/test",
			cookies:     cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("csrf_token", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).LowerCase().Validate(),
				NewCookie("_gat", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: "csrf_token", Value: "A1B2C3D4E5F6G7H8I9"},
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: "csrf_token", Value: "A1B2C3D4E5F6G7H8I9"},
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: "_gat", Value: "1"},
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: "_gat", Value: "1"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"csrf_token": "A1B2C3D4E5F6G7H8I9", "_gat": "1"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Not(false -> true).",
			method:      "POST",
			url:         "/test",
			cookies:     cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("_gat", nil).Chain().Not().Alpha(nil).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"_gat": "1"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Not(false -> true) multiple fields.",
			method:      "POST",
			url:         "/test",
			cookies:     cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("_gat", nil).Chain().Not().Alpha(nil).Validate(),
				NewCookie("_fat", nil).Chain().Not().Alpha(nil).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"_gat": "1", "_fat": "2"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Not(true -> false).",
			method:      "POST",
			url:         "/test",
			cookies:     cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("_gat", nil).Chain().Not().Numeric(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: "_gat", Value: "1"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"_gat": "1"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Not(true -> false) multiple fields.",
			method:      "POST",
			url:         "/test",
			cookies:     cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("_gat", nil).Chain().Not().Numeric(nil).Validate(),
				NewCookie("_fat", nil).Chain().Not().Numeric(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: "_gat", Value: "1"},
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: "_fat", Value: "2"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"_gat": "1", "_fat": "2"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Skip(true).",
			method:      "POST",
			url:         "/test",
			cookies:     cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie(`auth_token`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).Numeric(nil).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{`auth_token`: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Skip(true) multiple fields.",
			method:      "POST",
			url:         "/test",
			cookies:     cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie(`currency`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).Numeric(nil).Validate(),
				NewCookie(`_gat`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).Numeric(nil).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{`currency`: "USD", "_gat": "1"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Skip(false).",
			method:      "POST",
			url:         "/test",
			cookies:     cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie(`currency`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).Numeric(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: `currency`, Value: "USD"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{`currency`: "USD"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Skip(true) just one.",
			method:      "POST",
			url:         "/test",
			cookies:     cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie(`auth_token`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return true
					},
				).Numeric(nil).Currency(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: `auth_token`, Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{`auth_token`: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Skip(false) just one.",
			method:      "POST",
			url:         "/test",
			cookies:     cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie(`access_token`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						return false
					},
				).Numeric(nil).Currency(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: `access_token`, Value: "ya29.A0AfH6SMB..."},
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: `access_token`, Value: "ya29.A0AfH6SMB..."},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{`access_token`: "ya29.A0AfH6SMB..."}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Optional(present)",
			method:      "POST",
			url:         "/test",
			cookies:     cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("_gat", nil).Chain().Alpha(nil).Optional().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "cookies", Msg: defaultValChainErrMsg, Field: "_gat", Value: "1"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"_gat": "1"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Optional(not-present)",
			method:      "POST",
			url:         "/test",
			cookies:     cookies,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewCookie("tax_id", nil).Chain().Alpha(nil).Optional().Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"cookies": matchedDataFieldValues{"tax_id": ""}},
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
			req.Header.Set("Content-Type", test.contentType)

			for _, cookie := range test.cookies {
				req.AddCookie(cookie)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		})
	}
}
