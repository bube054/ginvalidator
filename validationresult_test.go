package ginvalidator

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestMessageFallbackAndCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	runChain := func(t *testing.T, chain gin.HandlerFunc) []ValidationChainError {
		t.Helper()
		w := httptest.NewRecorder()
		router := gin.New()

		var result []ValidationChainError
		router.POST("/test", chain, func(ctx *gin.Context) {
			var err error
			result, err = ValidationResult(ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})

		req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(`{"email":"notvalid"}`))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		return result
	}

	t.Run("tier 3: validatorgo message and code when no formatter set", func(t *testing.T) {
		errs := runChain(t, NewBody("email", nil).Chain().Email(nil).Validate())
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %d", len(errs))
		}
		if errs[0].Msg == DefaultValChainErrMsg {
			t.Errorf("expected validatorgo message, got default %q", errs[0].Msg)
		}
		if errs[0].Code == "" {
			t.Error("expected non-empty code from validatorgo")
		}
	})

	t.Run("tier 1: per-chain errFmtFunc overrides message", func(t *testing.T) {
		chain := NewBody("email", func(initial, sanitized, validatorName string) string {
			return "custom: " + validatorName
		}).Chain().Email(nil).Validate()
		errs := runChain(t, chain)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %d", len(errs))
		}
		if errs[0].Msg != "custom: Email" {
			t.Errorf("expected 'custom: Email', got %q", errs[0].Msg)
		}
		if errs[0].Code == "" {
			t.Error("expected non-empty code even with custom formatter")
		}
	})

	t.Run("tier 2: DefaultErrFmtFunc used when no per-chain formatter", func(t *testing.T) {
		origDefault := DefaultErrFmtFunc
		DefaultErrFmtFunc = func(initial, sanitized, validatorName string) string {
			return "default: " + validatorName
		}
		defer func() { DefaultErrFmtFunc = origDefault }()

		errs := runChain(t, NewBody("email", nil).Chain().Email(nil).Validate())
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %d", len(errs))
		}
		if errs[0].Msg != "default: Email" {
			t.Errorf("expected 'default: Email', got %q", errs[0].Msg)
		}
	})

	t.Run("tier 4: DefaultValChainErrMsg for CustomValidator", func(t *testing.T) {
		chain := NewBody("email", nil).Chain().CustomValidator(func(r *http.Request, initialValue, sanitizedValue string) bool {
			return false
		}).Validate()
		errs := runChain(t, chain)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %d", len(errs))
		}
		if errs[0].Msg != DefaultValChainErrMsg {
			t.Errorf("expected %q, got %q", DefaultValChainErrMsg, errs[0].Msg)
		}
		if errs[0].Code != "" {
			t.Errorf("expected empty code for CustomValidator, got %q", errs[0].Code)
		}
	})
}

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

func TestHasErrors(t *testing.T) {
	t.Run("returns false with no errors", func(t *testing.T) {
		ctx := createTestGinCtx(ginCtxReqOpts{})
		saveValidationErrorsToCtx(ctx, []ValidationChainError{})
		if HasErrors(ctx) {
			t.Error("expected false, got true")
		}
	})

	t.Run("returns true with errors", func(t *testing.T) {
		ctx := createTestGinCtx(ginCtxReqOpts{})
		saveValidationErrorsToCtx(ctx, []ValidationChainError{
			NewValidationChainError(vceWithField("email"), vceWithMsg("bad"), vceWithLocation("body"), vceWithValue("x")),
		})
		if !HasErrors(ctx) {
			t.Error("expected true, got false")
		}
	})

	t.Run("returns false on nil ctx", func(t *testing.T) {
		if HasErrors(nil) {
			t.Error("expected false for nil ctx")
		}
	})
}

func TestFirstError(t *testing.T) {
	t.Run("returns nil with no errors", func(t *testing.T) {
		ctx := createTestGinCtx(ginCtxReqOpts{})
		saveValidationErrorsToCtx(ctx, []ValidationChainError{})
		if got := FirstError(ctx); got != nil {
			t.Errorf("expected nil, got %+v", got)
		}
	})

	t.Run("returns first error", func(t *testing.T) {
		ctx := createTestGinCtx(ginCtxReqOpts{})
		saveValidationErrorsToCtx(ctx, []ValidationChainError{
			NewValidationChainError(vceWithField("email"), vceWithMsg("first"), vceWithLocation("body"), vceWithValue("x"), vceWithOrder(1)),
			NewValidationChainError(vceWithField("name"), vceWithMsg("second"), vceWithLocation("body"), vceWithValue("y"), vceWithOrder(2)),
		})
		got := FirstError(ctx)
		if got == nil {
			t.Fatal("expected non-nil")
		}
		if got.Field != "email" || got.Msg != "first" {
			t.Errorf("expected email/first, got %s/%s", got.Field, got.Msg)
		}
	})
}

func TestErrorsByField(t *testing.T) {
	ctx := createTestGinCtx(ginCtxReqOpts{})
	saveValidationErrorsToCtx(ctx, []ValidationChainError{
		NewValidationChainError(vceWithField("email"), vceWithMsg("err1"), vceWithLocation("body"), vceWithValue("x"), vceWithOrder(1)),
		NewValidationChainError(vceWithField("email"), vceWithMsg("err2"), vceWithLocation("body"), vceWithValue("x"), vceWithOrder(2)),
		NewValidationChainError(vceWithField("name"), vceWithMsg("err3"), vceWithLocation("body"), vceWithValue("y"), vceWithOrder(3)),
	})

	grouped, err := ErrorsByField(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(grouped["email"]) != 2 {
		t.Errorf("expected 2 email errors, got %d", len(grouped["email"]))
	}
	if len(grouped["name"]) != 1 {
		t.Errorf("expected 1 name error, got %d", len(grouped["name"]))
	}
}

func TestFirstErrorByField(t *testing.T) {
	ctx := createTestGinCtx(ginCtxReqOpts{})
	saveValidationErrorsToCtx(ctx, []ValidationChainError{
		NewValidationChainError(vceWithField("email"), vceWithMsg("err1"), vceWithLocation("body"), vceWithValue("x"), vceWithOrder(1)),
		NewValidationChainError(vceWithField("email"), vceWithMsg("err2"), vceWithLocation("body"), vceWithValue("x"), vceWithOrder(2)),
		NewValidationChainError(vceWithField("name"), vceWithMsg("err3"), vceWithLocation("body"), vceWithValue("y"), vceWithOrder(3)),
	})

	firsts, err := FirstErrorByField(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(firsts) != 2 {
		t.Errorf("expected 2 fields, got %d", len(firsts))
	}
	if firsts["email"].Msg != "err1" {
		t.Errorf("expected first email error msg 'err1', got %q", firsts["email"].Msg)
	}
	if firsts["name"].Msg != "err3" {
		t.Errorf("expected first name error msg 'err3', got %q", firsts["name"].Msg)
	}
}
