package ginvalidator

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// serveBody sends body as JSON to a route made of handlers followed by a final
// handler, and returns what that final handler saw.
func serveBody(t *testing.T, body string, handlers ...gin.HandlerFunc) (errs []ValidationChainError, handlerBody string) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	final := func(ctx *gin.Context) {
		errs, _ = ValidationResult(ctx)
		data, _ := io.ReadAll(ctx.Request.Body)
		handlerBody = string(data)
	}
	r.POST("/test", append(handlers, final)...)

	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), req)
	return errs, handlerBody
}

func TestInvalidJSONBody(t *testing.T) {
	broken := `{"email":"a@example.com"`

	errs, _ := serveBody(t, broken,
		NewBodyChain("email", nil).Email(nil).Validate(),
		NewBodyChain("bio", nil).Optional().Alpha(nil).Validate(),
	)

	if len(errs) != 2 {
		t.Fatalf("got %d errors, want one per chain: %+v", len(errs), errs)
	}
	for _, e := range errs {
		if e.Code != InvalidJSONErrCode || e.Message != InvalidJSONErrMsg || e.Location != "body" {
			t.Errorf("got %+v, want an invalid JSON error", e)
		}
	}
	if errs[0].Field != "email" || errs[1].Field != "bio" {
		t.Errorf("got fields %q and %q, want email and bio", errs[0].Field, errs[1].Field)
	}
}

func TestEmptyJSONBodyHasNoFields(t *testing.T) {
	errs, _ := serveBody(t, "",
		NewBodyChain("bio", nil).Optional().Alpha(nil).Validate(),
	)
	if len(errs) != 0 {
		t.Errorf("an empty body should count as no fields, got %+v", errs)
	}
}

func TestBodyIsReadOncePerRequest(t *testing.T) {
	// After the first chain, empty the request body. The second chain must still see
	// the field, because it uses the copy the first chain read.
	drain := func(ctx *gin.Context) {
		io.ReadAll(ctx.Request.Body)
		ctx.Next()
	}

	errs, _ := serveBody(t, `{"a":"x","b":"y"}`,
		NewBodyChain("a", nil).Equals("x").Validate(),
		drain,
		NewBodyChain("b", nil).Equals("y").Validate(),
	)
	if len(errs) != 0 {
		t.Errorf("got %+v, want no errors", errs)
	}
}

func TestBodyStillReadableAfterChains(t *testing.T) {
	body := `{"a":"x","b":"y"}`
	_, handlerBody := serveBody(t, body,
		NewBodyChain("a", nil).Equals("x").Validate(),
		NewBodyChain("b", nil).Equals("y").Validate(),
	)
	if handlerBody != body {
		t.Errorf("handler read %q, want %q", handlerBody, body)
	}
}

// errReader fails every read, like a client that disconnects mid-body.
type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }

func serveRequest(t *testing.T, req *http.Request, handlers ...gin.HandlerFunc) (errs []ValidationChainError, md MatchedData) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/test", append(handlers, func(ctx *gin.Context) {
		errs, _ = ValidationResult(ctx)
		md, _ = GetMatchedData(ctx)
	})...)
	r.ServeHTTP(httptest.NewRecorder(), req)
	return errs, md
}

func TestUnreadableBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/test", errReader{})
	req.Header.Set("Content-Type", "application/json")

	errs, md := serveRequest(t, req, NewBodyChain("email", nil).Optional().Email(nil).Validate())

	if len(errs) != 1 || errs[0].Code != UnreadableRequestErrCode || errs[0].Field != "email" {
		t.Fatalf("got %+v, want one %s error for email", errs, UnreadableRequestErrCode)
	}
	if md.Has(BodyLocation, "email") {
		t.Errorf("an unreadable field must not be in the matched data")
	}
}

func TestUnsupportedContentTypeHasNoBodyFields(t *testing.T) {
	newReq := func() *http.Request {
		req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("email=a@example.com"))
		req.Header.Set("Content-Type", "text/plain")
		return req
	}

	// An optional field is simply not sent
	errs, md := serveRequest(t, newReq(), NewBodyChain("email", nil).Optional().Email(nil).Validate())
	if len(errs) != 0 || md.Has(BodyLocation, "email") {
		t.Errorf("optional: got errors %+v and matched %v, want neither", errs, md)
	}

	// A required field fails its rule with the normal message, not a request error
	errs, _ = serveRequest(t, newReq(), NewBodyChain("email", nil).Email(nil).Validate())
	if len(errs) != 1 || errs[0].Code == UnreadableRequestErrCode || errs[0].Code == InvalidJSONErrCode {
		t.Errorf("required: got %+v, want one ordinary validation error", errs)
	}
}

func TestNothingIsLogged(t *testing.T) {
	var buf strings.Builder
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/test/:id",
		NewBodyChain("email", nil).Email(nil).Validate(),
		NewCookieChain("session", nil).Optional().Ascii().Validate(),
		NewHeaderChain("x-request-id", nil).Optional().Ascii().Validate(),
		NewParamChain("id", nil).Numeric(nil).Validate(),
		NewQueryChain("page", nil).Optional().Int(nil).Validate(),
		func(ctx *gin.Context) {},
	)

	for _, ctype := range []string{"text/plain", "", "application/json"} {
		req := httptest.NewRequest(http.MethodPost, "/test/100%25", strings.NewReader(`{"email":`))
		if ctype != "" {
			req.Header.Set("Content-Type", ctype)
		}
		req.Header["x-request-id"] = []string{"abc"} // not in canonical case
		r.ServeHTTP(httptest.NewRecorder(), req)
	}

	if buf.Len() != 0 {
		t.Errorf("the library wrote to the log:\n%s", buf.String())
	}
}
