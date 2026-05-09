package ginvalidator

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestOneOf(t *testing.T) {
	gin.SetMode(gin.TestMode)

	runOneOf := func(t *testing.T, body string, chainGroups ...[]ValidationChain) ([]ValidationChainError, MatchedData) {
		t.Helper()
		w := httptest.NewRecorder()
		router := gin.New()

		var errs []ValidationChainError
		var md MatchedData
		router.POST("/test", OneOf(chainGroups...), func(ctx *gin.Context) {
			errs, _ = ValidationResult(ctx)
			md, _ = GetMatchedData(ctx)
		})

		req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		return errs, md
	}

	t.Run("passes when first group succeeds", func(t *testing.T) {
		errs, md := runOneOf(t, `{"email":"a@b.com"}`,
			[]ValidationChain{NewBody("email", nil).Chain().Email(nil)},
			[]ValidationChain{NewBody("email", nil).Chain().Numeric(nil)},
		)
		if len(errs) != 0 {
			t.Errorf("expected 0 errors, got %d: %+v", len(errs), errs)
		}
		if !md.Has(BodyLocation, "email") {
			t.Error("expected matched data for email")
		}
	})

	t.Run("passes when second group succeeds", func(t *testing.T) {
		errs, _ := runOneOf(t, `{"count":"42"}`,
			[]ValidationChain{NewBody("count", nil).Chain().Alpha(nil)},
			[]ValidationChain{NewBody("count", nil).Chain().Numeric(nil)},
		)
		if len(errs) != 0 {
			t.Errorf("expected 0 errors, got %d: %+v", len(errs), errs)
		}
	})

	t.Run("fails when all groups fail", func(t *testing.T) {
		errs, _ := runOneOf(t, `{"value":"!!!"}`,
			[]ValidationChain{NewBody("value", nil).Chain().Alpha(nil)},
			[]ValidationChain{NewBody("value", nil).Chain().Numeric(nil)},
		)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %d: %+v", len(errs), errs)
		}
		if errs[0].Field != "_oneOf" {
			t.Errorf("expected field '_oneOf', got %q", errs[0].Field)
		}
	})

	t.Run("multi-chain group must all pass", func(t *testing.T) {
		errs, _ := runOneOf(t, `{"a":"hello","b":"world"}`,
			[]ValidationChain{
				NewBody("a", nil).Chain().Alpha(nil),
				NewBody("b", nil).Chain().Numeric(nil),
			},
			[]ValidationChain{
				NewBody("a", nil).Chain().Numeric(nil),
				NewBody("b", nil).Chain().Alpha(nil),
			},
		)
		// Group 1: a=alpha(pass), b=numeric(fail) → fail
		// Group 2: a=numeric(fail), b=alpha(pass) → fail
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %d: %+v", len(errs), errs)
		}
	})

	t.Run("multi-chain group all pass together", func(t *testing.T) {
		errs, _ := runOneOf(t, `{"a":"hello","b":"42"}`,
			[]ValidationChain{
				NewBody("a", nil).Chain().Alpha(nil),
				NewBody("b", nil).Chain().Numeric(nil),
			},
		)
		if len(errs) != 0 {
			t.Errorf("expected 0 errors, got %d: %+v", len(errs), errs)
		}
	})
}
