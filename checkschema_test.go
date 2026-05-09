package ginvalidator

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCheckSchema(t *testing.T) {
	gin.SetMode(gin.TestMode)

	runSchema := func(t *testing.T, body string, schema Schema) ([]ValidationChainError, MatchedData) {
		t.Helper()
		w := httptest.NewRecorder()
		router := gin.New()

		var errs []ValidationChainError
		var md MatchedData
		router.POST("/test", CheckSchema(schema), func(ctx *gin.Context) {
			errs, _ = ValidationResult(ctx)
			md, _ = GetMatchedData(ctx)
		})

		req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		return errs, md
	}

	t.Run("all fields pass", func(t *testing.T) {
		errs, md := runSchema(t, `{"email":"a@b.com","name":"Alice"}`, Schema{
			"email": {In: BodyLocation, Build: func(vc ValidationChain) ValidationChain { return vc.Email(nil) }},
			"name":  {In: BodyLocation, Build: func(vc ValidationChain) ValidationChain { return vc.Alpha(nil) }},
		})
		if len(errs) != 0 {
			t.Errorf("expected 0 errors, got %d: %+v", len(errs), errs)
		}
		if !md.Has(BodyLocation, "email") || !md.Has(BodyLocation, "name") {
			t.Error("expected matched data for both fields")
		}
	})

	t.Run("field fails validation", func(t *testing.T) {
		errs, _ := runSchema(t, `{"email":"notanemail","name":"Alice"}`, Schema{
			"email": {In: BodyLocation, Build: func(vc ValidationChain) ValidationChain { return vc.Email(nil) }},
			"name":  {In: BodyLocation, Build: func(vc ValidationChain) ValidationChain { return vc.Alpha(nil) }},
		})
		if len(errs) != 1 {
			t.Errorf("expected 1 error, got %d: %+v", len(errs), errs)
		}
		if len(errs) > 0 && errs[0].Field != "email" {
			t.Errorf("expected error on 'email', got %q", errs[0].Field)
		}
	})

	t.Run("optional field skipped when empty", func(t *testing.T) {
		errs, _ := runSchema(t, `{"email":"a@b.com"}`, Schema{
			"email": {In: BodyLocation, Build: func(vc ValidationChain) ValidationChain { return vc.Email(nil) }},
			"nick":  {In: BodyLocation, Optional: true, Build: func(vc ValidationChain) ValidationChain { return vc.Alpha(nil) }},
		})
		if len(errs) != 0 {
			t.Errorf("expected 0 errors, got %d: %+v", len(errs), errs)
		}
	})

	t.Run("bail stops on first error per field", func(t *testing.T) {
		errs, _ := runSchema(t, `{"val":"!!!"}`, Schema{
			"val": {In: BodyLocation, Build: func(vc ValidationChain) ValidationChain {
				return vc.Alpha(nil).Bail().Numeric(nil)
			}},
		})
		if len(errs) != 1 {
			t.Errorf("expected 1 error (bail), got %d: %+v", len(errs), errs)
		}
	})

	t.Run("custom error formatter", func(t *testing.T) {
		errs, _ := runSchema(t, `{"age":"abc"}`, Schema{
			"age": {
				In: BodyLocation,
				ErrFmtFunc: func(initial, sanitized, validatorName string) string {
					return "bad age"
				},
				Build: func(vc ValidationChain) ValidationChain { return vc.Numeric(nil) },
			},
		})
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %d", len(errs))
		}
		if errs[0].Message != "bad age" {
			t.Errorf("expected 'bad age', got %q", errs[0].Message)
		}
	})

	t.Run("nil build runs with no validators", func(t *testing.T) {
		errs, md := runSchema(t, `{"x":"anything"}`, Schema{
			"x": {In: BodyLocation},
		})
		if len(errs) != 0 {
			t.Errorf("expected 0 errors, got %d: %+v", len(errs), errs)
		}
		if !md.Has(BodyLocation, "x") {
			t.Error("expected matched data for x")
		}
	})
}
