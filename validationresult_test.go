package ginvalidator

import (
	"errors"
	"slices"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestValidationResult(t *testing.T) {
	tests := []struct {
		name                         string
		ctx                          *gin.Context
		insertedValidationErrors     []ValidationChainError
		expectedValidationChainError []ValidationChainError
		expectedErr                  error
	}{
		{
			name: "Nil ctx provided",
			ctx:  nil,
			insertedValidationErrors: []ValidationChainError{
				NewValidationChainError(vceWithLocation("body"), vceWithMsg("invalid value"), vceWithField("invalidField"), vceWithValue("value")),
			},
			expectedValidationChainError: []ValidationChainError{
				NewValidationChainError(vceWithLocation("body"), vceWithMsg("invalid value"), vceWithField("invalidField"), vceWithValue("value")),
			},
			expectedErr: ErrNilCtxValidationResult,
		},
		{
			name: "Valid validation errors for body location, with 1 item.",
			ctx:  createTestGinCtx(ginCtxReqOpts{}),
			insertedValidationErrors: []ValidationChainError{
				NewValidationChainError(vceWithLocation("body"), vceWithMsg("invalid value"), vceWithField("invalidField"), vceWithValue("value")),
			},
			expectedValidationChainError: []ValidationChainError{
				NewValidationChainError(vceWithLocation("body"), vceWithMsg("invalid value"), vceWithField("invalidField"), vceWithValue("value")),
			},
			expectedErr: nil,
		},
		{
			name: "Valid validation errors for headers location, with 3 items.",
			ctx:  createTestGinCtx(ginCtxReqOpts{}),
			insertedValidationErrors: []ValidationChainError{
				NewValidationChainError(vceWithLocation("headers"), vceWithMsg("invalid value1"), vceWithField("invalidField1"), vceWithValue("value1")),
				NewValidationChainError(vceWithLocation("headers"), vceWithMsg("invalid value2"), vceWithField("invalidField2"), vceWithValue("value2")),
				NewValidationChainError(vceWithLocation("headers"), vceWithMsg("invalid value3"), vceWithField("invalidField3"), vceWithValue("value3")),
			},
			expectedValidationChainError: []ValidationChainError{
				NewValidationChainError(vceWithLocation("headers"), vceWithMsg("invalid value1"), vceWithField("invalidField1"), vceWithValue("value1")),
				NewValidationChainError(vceWithLocation("headers"), vceWithMsg("invalid value2"), vceWithField("invalidField2"), vceWithValue("value2")),
				NewValidationChainError(vceWithLocation("headers"), vceWithMsg("invalid value3"), vceWithField("invalidField3"), vceWithValue("value3")),
			},
			expectedErr: nil,
		},
		{
			name: "Valid validation errors for cookies location, with 1 item.",
			ctx:  createTestGinCtx(ginCtxReqOpts{}),
			insertedValidationErrors: []ValidationChainError{
				NewValidationChainError(vceWithLocation("cookies"), vceWithMsg("invalid value"), vceWithField("invalidField"), vceWithValue("value")),
			},
			expectedValidationChainError: []ValidationChainError{
				NewValidationChainError(vceWithLocation("cookies"), vceWithMsg("invalid value"), vceWithField("invalidField"), vceWithValue("value")),
			},
			expectedErr: nil,
		},
		{
			name: "Valid validation errors for params location, with 1 item.",
			ctx:  createTestGinCtx(ginCtxReqOpts{}),
			insertedValidationErrors: []ValidationChainError{
				NewValidationChainError(vceWithLocation("params"), vceWithMsg("invalid value"), vceWithField("invalidField"), vceWithValue("value")),
			},
			expectedValidationChainError: []ValidationChainError{
				NewValidationChainError(vceWithLocation("params"), vceWithMsg("invalid value"), vceWithField("invalidField"), vceWithValue("value")),
			},
			expectedErr: nil,
		},
		{
			name: "Valid validation errors for query location, with 4 item.",
			ctx:  createTestGinCtx(ginCtxReqOpts{}),
			insertedValidationErrors: []ValidationChainError{
				NewValidationChainError(vceWithLocation("query"), vceWithMsg("invalid value"), vceWithField("invalidField"), vceWithValue("value")),
				NewValidationChainError(vceWithLocation("query"), vceWithMsg("invalid value1"), vceWithField("invalidField1"), vceWithValue("value1")),
				NewValidationChainError(vceWithLocation("query"), vceWithMsg("invalid value2"), vceWithField("invalidField2"), vceWithValue("value2")),
				NewValidationChainError(vceWithLocation("query"), vceWithMsg("invalid value3"), vceWithField("invalidField3"), vceWithValue("value3")),
			},
			expectedValidationChainError: []ValidationChainError{
				NewValidationChainError(vceWithLocation("query"), vceWithMsg("invalid value"), vceWithField("invalidField"), vceWithValue("value")),
				NewValidationChainError(vceWithLocation("query"), vceWithMsg("invalid value1"), vceWithField("invalidField1"), vceWithValue("value1")),
				NewValidationChainError(vceWithLocation("query"), vceWithMsg("invalid value2"), vceWithField("invalidField2"), vceWithValue("value2")),
				NewValidationChainError(vceWithLocation("query"), vceWithMsg("invalid value3"), vceWithField("invalidField3"), vceWithValue("value3")),
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			saveValidationErrorsToCtx(test.ctx, test.insertedValidationErrors)

			actualValidationResult, actualErr := ValidationResult(test.ctx)

			if actualErr != nil {
				if !errors.Is(actualErr, test.expectedErr) {
					t.Errorf("got %+v, want %+v", actualErr, test.expectedErr)
				}
			} else {
				if !slices.Equal(actualValidationResult, test.expectedValidationChainError) {
					t.Errorf("got %+v, want %+v", actualValidationResult, test.expectedValidationChainError)
				}
			}
		})
	}
}
