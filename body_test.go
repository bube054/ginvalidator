package ginvalidator

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	// vgo "github.com/bube054/validatorgo"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestBodyValidationChain(t *testing.T) {
	jsonBody := `{
		"name": {"first": "Tom", "last": "Anderson"},
		"age":37,
		"message": "A good saying is 7 comes after ate.",
		"children": ["Sara","Alex","Jack"],
		"fav.movie": "Deer Hunter",
		"friends": [
			{"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
			{"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
			{"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
		]
	}`

	urlEncodedBody := "name=John&age=30&email=john@example.com"

	formDataBody := `--6c55825619090769257acc8079eeba85e5e9874c7116b71da35065945dd9
Content-Disposition: form-data; name="name"

John
--6c55825619090769257acc8079eeba85e5e9874c7116b71da35065945dd9
Content-Disposition: form-data; name="email"

john@example.com
--6c55825619090769257acc8079eeba85e5e9874c7116b71da35065945dd9--
`

	tests := []struct {
		name        string
		method      string
		url         string
		body        string
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
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("name.first", nil).Chain().Ascii().Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"name.first": "Tom"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Validator(fail).",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("name.first", nil).Chain().Numeric(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "body", Msg: defaultValChainErrMsg, Field: "name.first", Value: "Tom"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"name.first": "Tom"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Validator(pass) multiple fields.",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("name.first", nil).Chain().Ascii().Validate(),
				NewBody("name.last", nil).Chain().Alpha(nil).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"name.first": "Tom", "name.last": "Anderson"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Validator(fail) multiple fields.",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("name.first", nil).Chain().Numeric(nil).Validate(),
				NewBody("name.last", nil).Chain().Currency(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "body", Msg: defaultValChainErrMsg, Field: "name.first", Value: "Tom"},
				{Location: "body", Msg: defaultValChainErrMsg, Field: "name.last", Value: "Anderson"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"name.first": "Tom", "name.last": "Anderson"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Custom Validator(pass).",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("name.first", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						data, err := io.ReadAll(req.Body)
						if err != nil {
							panic(fmt.Errorf("Custom validator could not read req body err: %w", err))
						}

						if string(data) != jsonBody {
							panic(fmt.Errorf("Custom validator req bodies do not match body: %s", data))
						}

						if initialValue != "Tom" {
							panic(fmt.Errorf("Custom validator initial value is invalid value: %s", initialValue))
						}

						if sanitizedValue != "Tom" {
							panic(fmt.Errorf("Custom validator sanitized value is invalid value: %s", sanitizedValue))
						}

						return true
					},
				).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"name.first": "Tom"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Custom Validator(fail).",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("name.first", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						data, err := io.ReadAll(req.Body)
						if err != nil {
							panic(fmt.Errorf("Custom validator could not read req body err: %w", err))
						}

						if string(data) != jsonBody {
							panic(fmt.Errorf("Custom validator req bodies do not match body: %s", data))
						}

						if initialValue != "Tom" {
							panic(fmt.Errorf("Custom validator initial value is invalid value: %s", initialValue))
						}

						if sanitizedValue != "Tom" {
							panic(fmt.Errorf("Custom validator sanitized value is invalid value: %s", sanitizedValue))
						}

						return false
					},
				).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "body", Msg: defaultValChainErrMsg, Field: "name.first", Value: "Tom"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"name.first": "Tom"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Custom Validator(pass) multiple fields.",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("name.first", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						data, err := io.ReadAll(req.Body)
						if err != nil {
							panic(fmt.Errorf("Custom validator could not read req body err: %w", err))
						}

						req.Body = io.NopCloser(bytes.NewBuffer(data))

						if string(data) != jsonBody {
							panic(fmt.Errorf("Custom validator req bodies do not match body: %s", data))
						}

						if initialValue != "Tom" {
							panic(fmt.Errorf("Custom validator initial value is invalid value: %s", initialValue))
						}

						if sanitizedValue != "Tom" {
							panic(fmt.Errorf("Custom validator sanitized value is invalid value: %s", sanitizedValue))
						}

						return true
					},
				).Validate(),
				NewBody("name.last", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						data, err := io.ReadAll(req.Body)
						if err != nil {
							panic(fmt.Errorf("Custom validator could not read req body err: %w", err))
						}

						if string(data) != jsonBody {
							panic(fmt.Errorf("Custom validator req bodies do not match body: %s", data))
						}

						if initialValue != "Anderson" {
							panic(fmt.Errorf("Custom validator initial value is invalid value: %s", initialValue))
						}

						if sanitizedValue != "Anderson" {
							panic(fmt.Errorf("Custom validator sanitized value is invalid value: %s", sanitizedValue))
						}

						return true
					},
				).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"name.first": "Tom", "name.last": "Anderson"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Custom Validator(fail) multiple fields.",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("name.first", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						data, err := io.ReadAll(req.Body)
						if err != nil {
							panic(fmt.Errorf("Custom validator could not read req body err: %w", err))
						}

						req.Body = io.NopCloser(bytes.NewBuffer(data))

						if string(data) != jsonBody {
							panic(fmt.Errorf("Custom validator req bodies do not match body: %s", data))
						}

						if initialValue != "Tom" {
							panic(fmt.Errorf("Custom validator initial value is invalid value: %s", initialValue))
						}

						if sanitizedValue != "Tom" {
							panic(fmt.Errorf("Custom validator sanitized value is invalid value: %s", sanitizedValue))
						}

						return false
					},
				).Validate(),
				NewBody("name.last", nil).Chain().CustomValidator(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						data, err := io.ReadAll(req.Body)
						if err != nil {
							panic(fmt.Errorf("Custom validator could not read req body err: %w", err))
						}

						if string(data) != jsonBody {
							panic(fmt.Errorf("Custom validator req bodies do not match body: %s", data))
						}

						if initialValue != "Anderson" {
							panic(fmt.Errorf("Custom validator initial value is invalid value: %s", initialValue))
						}

						if sanitizedValue != "Anderson" {
							panic(fmt.Errorf("Custom validator sanitized value is invalid value: %s", sanitizedValue))
						}

						return false
					},
				).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "body", Msg: defaultValChainErrMsg, Field: "name.first", Value: "Tom"},
				{Location: "body", Msg: defaultValChainErrMsg, Field: "name.last", Value: "Anderson"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"name.first": "Tom", "name.last": "Anderson"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Sanitizer.",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("name.last", nil).Chain().Whitelist("a-z").Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"name.last": "nderson"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Sanitizer multiple fields.",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("name.last", nil).Chain().Whitelist("a-z").Validate(),
				NewBody("name.first", nil).Chain().Whitelist("a-z").Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"name.last": "nderson", "name.first": "om"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Custom Sanitizer.",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("name.last", nil).Chain().CustomSanitizer(
					func(req *http.Request, initialValue, sanitizedValue string) string {
						data, err := io.ReadAll(req.Body)
						if err != nil {
							panic(fmt.Errorf("Custom Sanitizer could not read req body err: %w", err))
						}

						if string(data) != jsonBody {
							panic(fmt.Errorf("Custom Sanitizer req bodies do not match body: %s", data))
						}

						if initialValue != "Anderson" {
							panic(fmt.Errorf("Custom Sanitizer initial value is invalid value: %s", initialValue))
						}

						if sanitizedValue != "Anderson" {
							panic(fmt.Errorf("Custom Sanitizer sanitized value is invalid value: %s", sanitizedValue))
						}

						return "custom-sanitizer"
					},
				).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"name.last": "custom-sanitizer"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Custom Sanitizer multiple fields.",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("name.last", nil).Chain().CustomSanitizer(
					func(req *http.Request, initialValue, sanitizedValue string) string {
						data, err := io.ReadAll(req.Body)
						if err != nil {
							panic(fmt.Errorf("Custom Sanitizer could not read req body err: %w", err))
						}

						req.Body = io.NopCloser(bytes.NewBuffer(data))

						if string(data) != jsonBody {
							panic(fmt.Errorf("Custom Sanitizer req bodies do not match body: %s", data))
						}

						if initialValue != "Anderson" {
							panic(fmt.Errorf("Custom Sanitizer initial value is invalid value: %s", initialValue))
						}

						if sanitizedValue != "Anderson" {
							panic(fmt.Errorf("Custom Sanitizer sanitized value is invalid value: %s", sanitizedValue))
						}

						return "Tom"
					},
				).Validate(),
				NewBody("name.first", nil).Chain().CustomSanitizer(
					func(req *http.Request, initialValue, sanitizedValue string) string {
						data, err := io.ReadAll(req.Body)
						if err != nil {
							panic(fmt.Errorf("Custom Sanitizer could not read req body err: %w", err))
						}

						if string(data) != jsonBody {
							panic(fmt.Errorf("Custom Sanitizer req bodies do not match body: %s", data))
						}

						if initialValue != "Tom" {
							panic(fmt.Errorf("Custom Sanitizer initial value is invalid value: %s", initialValue))
						}

						if sanitizedValue != "Tom" {
							panic(fmt.Errorf("Custom Sanitizer sanitized value is invalid value: %s", sanitizedValue))
						}

						return "Anderson"
					},
				).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"name.last": "Tom", "name.first": "Anderson"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Modifier(bail).",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("message", nil).Chain().Alpha(nil).Bail().LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "body", Msg: defaultValChainErrMsg, Field: "message", Value: "A good saying is 7 comes after ate."},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"message": "A good saying is 7 comes after ate."}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Modifier(bail) multiple fields.",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("message", nil).Chain().Alpha(nil).Bail().LowerCase().Validate(),
				NewBody("friends.0.age", nil).Chain().Alpha(nil).Bail().LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "body", Msg: defaultValChainErrMsg, Field: "message", Value: "A good saying is 7 comes after ate."},
				{Location: "body", Msg: defaultValChainErrMsg, Field: "friends.0.age", Value: "44"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"message": "A good saying is 7 comes after ate.", "friends.0.age": "44"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Modifier(if/bail).",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("message", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						data, err := io.ReadAll(req.Body)
						if err != nil {
							panic(fmt.Errorf("If modifier could not read req body err: %w", err))
						}

						if string(data) != jsonBody {
							panic(fmt.Errorf("If modifier req bodies do not match body: %s", data))
						}

						if initialValue != "A good saying is 7 comes after ate." {
							panic(fmt.Errorf("If modifier initial value is invalid value: %s", initialValue))
						}

						if sanitizedValue != "A good saying is 7 comes after ate." {
							panic(fmt.Errorf("If modifier sanitized value is invalid value: %s", sanitizedValue))
						}

						return true
					},
				).LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "body", Msg: defaultValChainErrMsg, Field: "message", Value: "A good saying is 7 comes after ate."},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"message": "A good saying is 7 comes after ate."}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Modifier(if/bail) multiple fields.",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("message", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						data, err := io.ReadAll(req.Body)
						if err != nil {
							panic(fmt.Errorf("If modifier could not read req body err: %w", err))
						}

						req.Body = io.NopCloser(bytes.NewBuffer(data))

						if string(data) != jsonBody {
							panic(fmt.Errorf("If modifier req bodies do not match body: %s", data))
						}

						if initialValue != "A good saying is 7 comes after ate." {
							panic(fmt.Errorf("If modifier initial value is invalid value: %s", initialValue))
						}

						if sanitizedValue != "A good saying is 7 comes after ate." {
							panic(fmt.Errorf("If modifier sanitized value is invalid value: %s", sanitizedValue))
						}

						return true
					},
				).LowerCase().Validate(),
				NewBody("age", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						data, err := io.ReadAll(req.Body)
						if err != nil {
							panic(fmt.Errorf("If modifier could not read req body err: %w", err))
						}

						if string(data) != jsonBody {
							panic(fmt.Errorf("If modifier req bodies do not match body: %s", data))
						}

						if initialValue != "37" {
							panic(fmt.Errorf("If modifier initial value is invalid value: %s", initialValue))
						}

						if sanitizedValue != "37" {
							panic(fmt.Errorf("If modifier sanitized value is invalid value: %s", sanitizedValue))
						}

						return true
					},
				).LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "body", Msg: defaultValChainErrMsg, Field: "message", Value: "A good saying is 7 comes after ate."},
				{Location: "body", Msg: defaultValChainErrMsg, Field: "age", Value: "37"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"message": "A good saying is 7 comes after ate.", "age": "37"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Modifier(if/proceed).",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("message", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						data, err := io.ReadAll(req.Body)
						if err != nil {
							panic(fmt.Errorf("If modifier could not read req body err: %w", err))
						}

						if string(data) != jsonBody {
							panic(fmt.Errorf("If modifier req bodies do not match body: %s", data))
						}

						if initialValue != "A good saying is 7 comes after ate." {
							panic(fmt.Errorf("If modifier initial value is invalid value: %s", initialValue))
						}

						if sanitizedValue != "A good saying is 7 comes after ate." {
							panic(fmt.Errorf("If modifier sanitized value is invalid value: %s", sanitizedValue))
						}

						return false
					},
				).LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "body", Msg: defaultValChainErrMsg, Field: "message", Value: "A good saying is 7 comes after ate."},
				{Location: "body", Msg: defaultValChainErrMsg, Field: "message", Value: "A good saying is 7 comes after ate."},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"message": "A good saying is 7 comes after ate."}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Modifier(if/proceed) multiple fields.",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("message", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						data, err := io.ReadAll(req.Body)
						if err != nil {
							panic(fmt.Errorf("If modifier could not read req body err: %w", err))
						}

						req.Body = io.NopCloser(bytes.NewBuffer(data))

						if string(data) != jsonBody {
							panic(fmt.Errorf("If modifier req bodies do not match body: %s", data))
						}

						if initialValue != "A good saying is 7 comes after ate." {
							panic(fmt.Errorf("If modifier initial value is invalid value: %s", initialValue))
						}

						if sanitizedValue != "A good saying is 7 comes after ate." {
							panic(fmt.Errorf("If modifier sanitized value is invalid value: %s", sanitizedValue))
						}

						return false
					},
				).LowerCase().Validate(),
				NewBody("age", nil).Chain().Alpha(nil).If(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						data, err := io.ReadAll(req.Body)
						if err != nil {
							panic(fmt.Errorf("If modifier could not read req body err: %w", err))
						}

						req.Body = io.NopCloser(bytes.NewBuffer(data))

						if string(data) != jsonBody {
							panic(fmt.Errorf("If modifier req bodies do not match body: %s", data))
						}

						if initialValue != "37" {
							panic(fmt.Errorf("If modifier initial value is invalid value: %s", initialValue))
						}

						if sanitizedValue != "37" {
							panic(fmt.Errorf("If modifier sanitized value is invalid value: %s", sanitizedValue))
						}

						return false
					},
				).LowerCase().Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "body", Msg: defaultValChainErrMsg, Field: "message", Value: "A good saying is 7 comes after ate."},
				{Location: "body", Msg: defaultValChainErrMsg, Field: "message", Value: "A good saying is 7 comes after ate."},
				{Location: "body", Msg: defaultValChainErrMsg, Field: "age", Value: "37"},
				{Location: "body", Msg: defaultValChainErrMsg, Field: "age", Value: "37"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"message": "A good saying is 7 comes after ate.", "age": "37"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Not(false -> true).",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("age", nil).Chain().Not().Alpha(nil).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"age": "37"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Not(false -> true) multiple fields.",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("age", nil).Chain().Not().Alpha(nil).Validate(),
				NewBody("message", nil).Chain().Not().Alpha(nil).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"age": "37", "message": "A good saying is 7 comes after ate."}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Not(true -> false).",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("age", nil).Chain().Not().Numeric(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "body", Msg: defaultValChainErrMsg, Field: "age", Value: "37"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"age": "37"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Not(true -> false) multiple fields.",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("age", nil).Chain().Not().Numeric(nil).Validate(),
				NewBody("friends.1.age", nil).Chain().Not().Numeric(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "body", Msg: defaultValChainErrMsg, Field: "age", Value: "37"},
				{Location: "body", Msg: defaultValChainErrMsg, Field: "friends.1.age", Value: "68"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"age": "37", "friends.1.age": "68"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Skip(true).",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody(`fav\.movie`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						data, err := io.ReadAll(req.Body)
						if err != nil {
							panic(fmt.Errorf("Skip modifier could not read req body err: %w", err))
						}

						if string(data) != jsonBody {
							panic(fmt.Errorf("Skip modifier req bodies do not match body: %s", data))
						}

						if initialValue != "Deer Hunter" {
							panic(fmt.Errorf("Skip modifier initial value is invalid value: %s", initialValue))
						}

						if sanitizedValue != "Deer Hunter" {
							panic(fmt.Errorf("Skip modifier sanitized value is invalid value: %s", sanitizedValue))
						}

						return true
					},
				).Numeric(nil).Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{`fav\.movie`: "Deer Hunter"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Skip(true) multiple fields.",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody(`fav\.movie`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						data, err := io.ReadAll(req.Body)
						if err != nil {
							panic(fmt.Errorf("Skip modifier could not read req body err: %w", err))
						}

						req.Body = io.NopCloser(bytes.NewBuffer(data))

						if string(data) != jsonBody {
							panic(fmt.Errorf("Skip modifier req bodies do not match body: %s", data))
						}

						if initialValue != "Deer Hunter" {
							panic(fmt.Errorf("Skip modifier initial value is invalid value: %s", initialValue))
						}

						if sanitizedValue != "Deer Hunter" {
							panic(fmt.Errorf("Skip modifier sanitized value is invalid value: %s", sanitizedValue))
						}

						return true
					},
				).Numeric(nil).Validate(),
				NewBody(`age`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						data, err := io.ReadAll(req.Body)
						if err != nil {
							panic(fmt.Errorf("Skip modifier could not read req body err: %w", err))
						}

						req.Body = io.NopCloser(bytes.NewBuffer(data))

						if string(data) != jsonBody {
							panic(fmt.Errorf("Skip modifier req bodies do not match body: %s", data))
						}

						if initialValue != "37" {
							panic(fmt.Errorf("Skip modifier initial value is invalid value: %s", initialValue))
						}

						if sanitizedValue != "37" {
							panic(fmt.Errorf("Skip modifier sanitized value is invalid value: %s", sanitizedValue))
						}

						return true
					},
				).Numeric(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				// {Location: "body", Msg: defaultValChainErrMsg, Field: `fav\.movie`, Value: "Deer Hunter"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{`fav\.movie`: "Deer Hunter", "age": "37"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Skip(false).",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody(`fav\.movie`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						data, err := io.ReadAll(req.Body)
						if err != nil {
							panic(fmt.Errorf("Skip modifier could not read req body err: %w", err))
						}

						if string(data) != jsonBody {
							panic(fmt.Errorf("Skip modifier req bodies do not match body: %s", data))
						}

						if initialValue != "Deer Hunter" {
							panic(fmt.Errorf("Skip modifier initial value is invalid value: %s", initialValue))
						}

						if sanitizedValue != "Deer Hunter" {
							panic(fmt.Errorf("Skip modifier sanitized value is invalid value: %s", sanitizedValue))
						}

						return false
					},
				).Numeric(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "body", Msg: defaultValChainErrMsg, Field: `fav\.movie`, Value: "Deer Hunter"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{`fav\.movie`: "Deer Hunter"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Skip(true) just one.",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody(`fav\.movie`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						data, err := io.ReadAll(req.Body)
						if err != nil {
							panic(fmt.Errorf("Skip modifier could not read req body err: %w", err))
						}

						if string(data) != jsonBody {
							panic(fmt.Errorf("Skip modifier req bodies do not match body: %s", data))
						}

						if initialValue != "Deer Hunter" {
							panic(fmt.Errorf("Skip modifier initial value is invalid value: %s", initialValue))
						}

						if sanitizedValue != "Deer Hunter" {
							panic(fmt.Errorf("Skip modifier sanitized value is invalid value: %s", sanitizedValue))
						}

						return true
					},
				).Numeric(nil).Currency(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "body", Msg: defaultValChainErrMsg, Field: `fav\.movie`, Value: "Deer Hunter"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{`fav\.movie`: "Deer Hunter"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Skip(false) just one.",
			method:      "POST",
			url:         "/test",
			body:        jsonBody,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody(`fav\.movie`, nil).Chain().Skip(
					func(req *http.Request, initialValue, sanitizedValue string) bool {
						data, err := io.ReadAll(req.Body)
						if err != nil {
							panic(fmt.Errorf("Skip modifier could not read req body err: %w", err))
						}

						if string(data) != jsonBody {
							panic(fmt.Errorf("Skip modifier req bodies do not match body: %s", data))
						}

						if initialValue != "Deer Hunter" {
							panic(fmt.Errorf("Skip modifier initial value is invalid value: %s", initialValue))
						}

						if sanitizedValue != "Deer Hunter" {
							panic(fmt.Errorf("Skip modifier sanitized value is invalid value: %s", sanitizedValue))
						}

						return false
					},
				).Numeric(nil).Currency(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "body", Msg: defaultValChainErrMsg, Field: `fav\.movie`, Value: "Deer Hunter"},
				{Location: "body", Msg: defaultValChainErrMsg, Field: `fav\.movie`, Value: "Deer Hunter"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{`fav\.movie`: "Deer Hunter"}},
			matchedDataErr:      nil,
		},
		// For "application/x-www-form-urlencoded
		{
			name:        "Test Validator(pass).",
			method:      "POST",
			url:         "/test",
			body:        urlEncodedBody,
			contentType: "application/x-www-form-urlencoded",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("name", nil).Chain().Ascii().Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"name": "John"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Validator(fail).",
			method:      "POST",
			url:         "/test",
			body:        urlEncodedBody,
			contentType: "application/x-www-form-urlencoded",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("email", nil).Chain().Numeric(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "body", Msg: defaultValChainErrMsg, Field: "email", Value: "john@example.com"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"email": "john@example.com"}},
			matchedDataErr:      nil,
		},
		// For "multipart/form-data"
		{
			name:        "Test Validator(pass).",
			method:      "POST",
			url:         "/test",
			body:        formDataBody,
			contentType: "multipart/form-data; boundary=6c55825619090769257acc8079eeba85e5e9874c7116b71da35065945dd9",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("name", nil).Chain().Ascii().Validate(),
			},
			validationResult:    []ValidationChainError{},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"name": "John"}},
			matchedDataErr:      nil,
		},
		{
			name:        "Test Validator(fail).",
			method:      "POST",
			url:         "/test",
			body:        formDataBody,
			contentType: "multipart/form-data; boundary=6c55825619090769257acc8079eeba85e5e9874c7116b71da35065945dd9",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("email", nil).Chain().Numeric(nil).Validate(),
			},
			validationResult: []ValidationChainError{
				{Location: "body", Msg: defaultValChainErrMsg, Field: "email", Value: "john@example.com"},
			},
			validationResultErr: nil,
			matchedData:         MatchedData{"body": matchedDataFieldValues{"email": "john@example.com"}},
			matchedDataErr:      nil,
		},
	}

	for _, test := range tests {
		test := test
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

			req, _ := http.NewRequest(test.method, test.url, strings.NewReader(test.body))
			req.Header.Set("Content-Type", test.contentType)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		})
	}
}
