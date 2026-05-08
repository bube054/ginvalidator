package ginvalidator

import (
	"errors"
	"slices"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
				NewValidationChainError(
					vceWithLocation("body"),
					vceWithMsg("invalid value"),
					vceWithField("invalidField"),
					vceWithValue("value"),
					vceWithOrder(1),
				),
			},
			expectedValidationChainError: []ValidationChainError{
				NewValidationChainError(
					vceWithLocation("body"),
					vceWithMsg("invalid value"),
					vceWithField("invalidField"),
					vceWithValue("value"),
					vceWithOrder(1),
				),
			},
			expectedErr: ErrNilCtxValidationResult,
		},
		{
			name: "Valid validation errors for body location, with 1 item.",
			ctx:  createTestGinCtx(ginCtxReqOpts{}),
			insertedValidationErrors: []ValidationChainError{
				NewValidationChainError(
					vceWithLocation("body"),
					vceWithMsg("invalid value"),
					vceWithField("invalidField"),
					vceWithValue("value"),
					vceWithOrder(1),
				),
			},
			expectedValidationChainError: []ValidationChainError{
				NewValidationChainError(
					vceWithLocation("body"),
					vceWithMsg("invalid value"),
					vceWithField("invalidField"),
					vceWithValue("value"),
					vceWithOrder(1),
				),
			},
			expectedErr: nil,
		},
		{
			name: "Valid validation errors for headers location, with 3 items.",
			ctx:  createTestGinCtx(ginCtxReqOpts{}),
			insertedValidationErrors: []ValidationChainError{
				NewValidationChainError(vceWithLocation("headers"), vceWithMsg("invalid value1"), vceWithField("invalidField1"), vceWithValue("value1"), vceWithOrder(1)),
				NewValidationChainError(vceWithLocation("headers"), vceWithMsg("invalid value2"), vceWithField("invalidField2"), vceWithValue("value2"), vceWithOrder(2)),
				NewValidationChainError(vceWithLocation("headers"), vceWithMsg("invalid value3"), vceWithField("invalidField3"), vceWithValue("value3"), vceWithOrder(3)),
			},
			expectedValidationChainError: []ValidationChainError{
				NewValidationChainError(vceWithLocation("headers"), vceWithMsg("invalid value1"), vceWithField("invalidField1"), vceWithValue("value1"), vceWithOrder(1)),
				NewValidationChainError(vceWithLocation("headers"), vceWithMsg("invalid value2"), vceWithField("invalidField2"), vceWithValue("value2"), vceWithOrder(2)),
				NewValidationChainError(vceWithLocation("headers"), vceWithMsg("invalid value3"), vceWithField("invalidField3"), vceWithValue("value3"), vceWithOrder(3)),
			},
			expectedErr: nil,
		},
		{
			name: "Valid validation errors for cookies location, with 1 item.",
			ctx:  createTestGinCtx(ginCtxReqOpts{}),
			insertedValidationErrors: []ValidationChainError{
				NewValidationChainError(vceWithLocation("cookies"), vceWithMsg("invalid value"), vceWithField("invalidField"), vceWithValue("value"), vceWithOrder(1)),
			},
			expectedValidationChainError: []ValidationChainError{
				NewValidationChainError(vceWithLocation("cookies"), vceWithMsg("invalid value"), vceWithField("invalidField"), vceWithValue("value"), vceWithOrder(1)),
			},
			expectedErr: nil,
		},
		{
			name: "Valid validation errors for params location, with 1 item.",
			ctx:  createTestGinCtx(ginCtxReqOpts{}),
			insertedValidationErrors: []ValidationChainError{
				NewValidationChainError(vceWithLocation("params"), vceWithMsg("invalid value"), vceWithField("invalidField"), vceWithValue("value"), vceWithOrder(1)),
			},
			expectedValidationChainError: []ValidationChainError{
				NewValidationChainError(vceWithLocation("params"), vceWithMsg("invalid value"), vceWithField("invalidField"), vceWithValue("value"), vceWithOrder(1)),
			},
			expectedErr: nil,
		},
		{
			name: "Valid validation errors for query location, with 4 item.",
			ctx:  createTestGinCtx(ginCtxReqOpts{}),
			insertedValidationErrors: []ValidationChainError{
				NewValidationChainError(vceWithLocation("query"), vceWithMsg("invalid value"), vceWithField("invalidField"), vceWithValue("value"), vceWithOrder(1)),
				NewValidationChainError(vceWithLocation("query"), vceWithMsg("invalid value1"), vceWithField("invalidField1"), vceWithValue("value1"), vceWithOrder(2)),
				NewValidationChainError(vceWithLocation("query"), vceWithMsg("invalid value2"), vceWithField("invalidField2"), vceWithValue("value2"), vceWithOrder(3)),
				NewValidationChainError(vceWithLocation("query"), vceWithMsg("invalid value3"), vceWithField("invalidField3"), vceWithValue("value3"), vceWithOrder(4)),
			},
			expectedValidationChainError: []ValidationChainError{
				NewValidationChainError(vceWithLocation("query"), vceWithMsg("invalid value"), vceWithField("invalidField"), vceWithValue("value"), vceWithOrder(1)),
				NewValidationChainError(vceWithLocation("query"), vceWithMsg("invalid value1"), vceWithField("invalidField1"), vceWithValue("value1"), vceWithOrder(2)),
				NewValidationChainError(vceWithLocation("query"), vceWithMsg("invalid value2"), vceWithField("invalidField2"), vceWithValue("value2"), vceWithOrder(3)),
				NewValidationChainError(vceWithLocation("query"), vceWithMsg("invalid value3"), vceWithField("invalidField3"), vceWithValue("value3"), vceWithOrder(4)),
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

func TestSortErrorsByOrder(t *testing.T) {

	tests := []struct {
		name           string
		initial        []ValidationChainError
		expectedResult []ValidationChainError
	}{
		{
			name: "slice of 1 errors",
			initial: []ValidationChainError{
				NewValidationChainError(
					vceWithLocation("body"),
					vceWithMsg(DefaultValChainErrMsg),
					vceWithField("field"),
					vceWithValue("value"),
					vceWithOrder(1)),
			},
			expectedResult: []ValidationChainError{
				NewValidationChainError(
					vceWithLocation("body"),
					vceWithMsg(DefaultValChainErrMsg),
					vceWithField("field"),
					vceWithValue("value"),
					vceWithOrder(1)),
			},
		},
		{
			name: "slice of 2 errors already ordered",
			initial: []ValidationChainError{
				NewValidationChainError(
					vceWithLocation("body"),
					vceWithMsg(DefaultValChainErrMsg),
					vceWithField("field"),
					vceWithValue("value"),
					vceWithOrder(1)),
				NewValidationChainError(
					vceWithLocation("body"),
					vceWithMsg(DefaultValChainErrMsg),
					vceWithField("field"),
					vceWithValue("value"),
					vceWithOrder(2)),
			},
			expectedResult: []ValidationChainError{
				NewValidationChainError(
					vceWithLocation("body"),
					vceWithMsg(DefaultValChainErrMsg),
					vceWithField("field"),
					vceWithValue("value"),
					vceWithOrder(1)),
				NewValidationChainError(
					vceWithLocation("body"),
					vceWithMsg(DefaultValChainErrMsg),
					vceWithField("field"),
					vceWithValue("value"),
					vceWithOrder(2)),
			},
		},
		{
			name: "slice of 2 errors not already ordered",
			initial: []ValidationChainError{
				NewValidationChainError(
					vceWithLocation("body"),
					vceWithMsg(DefaultValChainErrMsg),
					vceWithField("field"),
					vceWithValue("value"),
					vceWithOrder(2)),
				NewValidationChainError(
					vceWithLocation("body"),
					vceWithMsg(DefaultValChainErrMsg),
					vceWithField("field"),
					vceWithValue("value"),
					vceWithOrder(1)),
			},
			expectedResult: []ValidationChainError{
				NewValidationChainError(
					vceWithLocation("body"),
					vceWithMsg(DefaultValChainErrMsg),
					vceWithField("field"),
					vceWithValue("value"),
					vceWithOrder(1)),
				NewValidationChainError(
					vceWithLocation("body"),
					vceWithMsg(DefaultValChainErrMsg),
					vceWithField("field"),
					vceWithValue("value"),
					vceWithOrder(2)),
			},
		},
		{
			name: "slice of 3 errors not already ordered",
			initial: []ValidationChainError{
				NewValidationChainError(
					vceWithLocation("body"),
					vceWithMsg(DefaultValChainErrMsg),
					vceWithField("field"),
					vceWithValue("value"),
					vceWithOrder(3)),
				NewValidationChainError(
					vceWithLocation("body"),
					vceWithMsg(DefaultValChainErrMsg),
					vceWithField("field"),
					vceWithValue("value"),
					vceWithOrder(2)),
				NewValidationChainError(
					vceWithLocation("body"),
					vceWithMsg(DefaultValChainErrMsg),
					vceWithField("field"),
					vceWithValue("value"),
					vceWithOrder(1)),
			},
			expectedResult: []ValidationChainError{
				NewValidationChainError(
					vceWithLocation("body"),
					vceWithMsg(DefaultValChainErrMsg),
					vceWithField("field"),
					vceWithValue("value"),
					vceWithOrder(1)),
				NewValidationChainError(
					vceWithLocation("body"),
					vceWithMsg(DefaultValChainErrMsg),
					vceWithField("field"),
					vceWithValue("value"),
					vceWithOrder(2)),
				NewValidationChainError(
					vceWithLocation("body"),
					vceWithMsg(DefaultValChainErrMsg),
					vceWithField("field"),
					vceWithValue("value"),
					vceWithOrder(3)),
			},
		},
	}

	for _, test := range tests {
		SortValidationErrors(test.initial)
		if !cmp.Equal(test.initial, test.expectedResult, cmpopts.IgnoreUnexported(ValidationChainError{}), cmpopts.EquateEmpty()) {
			t.Errorf("got %+v, wanted %+v", test.initial, test.expectedResult)
		}
	}
}
