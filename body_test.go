package ginvalidator

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"strings"
	"testing"

	// vgo "github.com/bube054/validatorgo"
	"github.com/gin-gonic/gin"
)

func TestBodyValidationChain(t *testing.T) {
	body := `{
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
		// {
		// 	name:        "Test Validator(pass).",
		// 	method:      "POST",
		// 	url:         "/test",
		// 	body:        body,
		// 	contentType: "application/json",
		// 	customValidatorsChain: []gin.HandlerFunc{
		// 		NewBody("name.first", nil).CreateChain().Ascii().Validate(),
		// 	},
		// 	validationResult:    []ValidationChainError{},
		// 	validationResultErr: nil,
		// 	matchedData:         MatchedData{"body": matchedDataFieldValues{"name.first": "Tom"}},
		// 	matchedDataErr:      nil,
		// },
		// {
		// 	name:        "Test Validator(fail).",
		// 	method:      "POST",
		// 	url:         "/test",
		// 	body:        body,
		// 	contentType: "application/json",
		// 	customValidatorsChain: []gin.HandlerFunc{
		// 		NewBody("name.first", nil).CreateChain().Numeric(nil).Validate(),
		// 	},
		// 	validationResult: []ValidationChainError{
		// 		{Location: "body", Msg: defaultValidationChainErrorMessage, Field: "name.first", Value: "Tom"},
		// 	},
		// 	validationResultErr: nil,
		// 	matchedData:         MatchedData{"body": matchedDataFieldValues{"name.first": "Tom"}},
		// 	matchedDataErr:      nil,
		// },
		{
			name:        "Test Custom Validator(pass).",
			method:      "POST",
			url:         "/test",
			body:        body,
			contentType: "application/json",
			customValidatorsChain: []gin.HandlerFunc{
				NewBody("name.first", nil).CreateChain().CustomValidator(
					func(req http.Request, initialValue, sanitizedValue string) bool {
						data, err := io.ReadAll(req.Body)
						if err != nil {
							panic(fmt.Errorf("Custom validator could not read req body err: %w", err))
						}

						fmt.Println("data:", string(data))
						// fmt.Println("data1", req.Host)
						// fmt.Printf("%+v \n", req)
						// fmt.Println(req.G)

						if string(data) != body {
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
		// {
		// 	name:        "Test Validator(fail).",
		// 	method:      "POST",
		// 	url:         "/test",
		// 	body:        body,
		// 	contentType: "application/json",
		// 	customValidatorsChain: []gin.HandlerFunc{
		// 		NewBody("name.first", nil).CreateChain().Numeric(nil).Validate(),
		// 	},
		// 	validationResult: []ValidationChainError{
		// 		{Location: "body", Msg: defaultValidationChainErrorMessage, Field: "name.first", Value: "Tom"},
		// 	},
		// 	validationResultErr: nil,
		// 	matchedData:         MatchedData{"body": matchedDataFieldValues{"name.first": "Tom"}},
		// 	matchedDataErr:      nil,
		// },
		// {
		// 	name:        "Test Sanitizer.",
		// 	method:      "POST",
		// 	url:         "/test",
		// 	body:        body,
		// 	contentType: "application/json",
		// 	customValidatorsChain: []gin.HandlerFunc{
		// 		NewBody("name.last", nil).CreateChain().Whitelist("a-z").Validate(),
		// 	},
		// 	validationResult:    []ValidationChainError{},
		// 	validationResultErr: nil,
		// 	matchedData:         MatchedData{"body": matchedDataFieldValues{"name.last": "nderson"}},
		// 	matchedDataErr:      nil,
		// },
		// For "application/x-www-form-urlencoded
		// For "multipart/form-data"
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			router := setupRouter()

			assertHandler := func(ctx *gin.Context) {
				validationResult, err := ValidationResult(ctx)

				if test.validationResultErr != nil {
					if !errors.Is(test.validationResultErr, err) {
						t.Errorf("got %v, wanted %v", err, test.validationResultErr)
					}
				} else {
					if !slices.Equal(test.validationResult, validationResult) {
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

// func bodySetupRouter() *gin.Engine {
// 	router := setupRouter()

// 	body := NewBody("message", nil)
// 	router.GET("/body",
// 		body.
// 			CreateChain().
// 			Contains("errors", nil).
// 			Blacklist("0-9").
// 			Alphanumeric(nil).
// 			Blacklist("0-9").
// 			Validate(),
// 		func(c *gin.Context) {
// 			errs, err := ValidationResult(c)

// 			if err != nil {
// 				fmt.Println("Could not retrieve validation result err:", err)
// 			} else {
// 				fmt.Printf("All Errors 🙌🙌🙌🙌 %+v\n", errs)
// 			}

// 			data, err := MatchedData(c)

// 			if err != nil {
// 				fmt.Println(err)
// 			}

// 			fmt.Println("data:", data)

// 			c.String(200, "pong")
// 		},
// 	)

// 	return router
// }

// func TestBody(t *testing.T) {
// 	// data := []byte(`{"name":"John"}`)
// 	router := bodySetupRouter()

// 	body := `{
// 		"name": {"first": "Tom", "last": "Anderson"},
// 		"age":37,
// 		"message": "A good saying is 7 comes after ate.",
// 		"children": ["Sara","Alex","Jack"],
// 		"fav.movie": "Deer Hunter",
// 		"friends": [
// 			{"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
// 			{"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
// 			{"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
// 		]
// 	}`

// 	req, _ := http.NewRequest("GET", "/body", strings.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	if w.Code != http.StatusOK {
// 		t.Errorf("request error")
// 	}
// }

// 1. Permutations of the Sequence
// Since order matters, you have
// 3
// !
// =
// 6
// 3!=6 permutations:

// Modifier -> Sanitizer -> Validator
// Modifier -> Validator -> Sanitizer
// Sanitizer -> Modifier -> Validator
// Sanitizer -> Validator -> Modifier
// Validator -> Modifier -> Sanitizer
// Validator -> Sanitizer -> Modifier
// 2. Combinations of Two Components
// You may also need to test cases where only two of the components are used together:

// Modifier -> Sanitizer
// Modifier -> Validator
// Sanitizer -> Validator
// Sanitizer -> Modifier
// Validator -> Modifier
// Validator -> Sanitizer
// 3. Single Component Tests
// Each component alone should also be tested to ensure individual behavior:

// Modifier Only
// Sanitizer Only
// Validator Only
// 4. Edge Cases and Special Considerations
// Empty Input: Test how each combination handles empty or nil input.
// Invalid Input: Ensure proper handling and error messages.
// Boundary Values: If your validator has rules like length checks or range limits, make sure you cover these.
