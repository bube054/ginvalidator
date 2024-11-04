package ginvalidator

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type ginCtxReqOpts struct {
	method      string
	url         string
	body        string
	contentType string
	headers     map[string]string
	cookies     []*http.Cookie
	params      gin.Params
}

func createTestGinCtx(opts ginCtxReqOpts) *gin.Context {
	gin.SetMode(gin.TestMode)

	if opts.method == "" {
		opts.method = http.MethodPost
	}

	if opts.url == "" {
		opts.url = "/test"
	}

	if opts.contentType == "" {
		opts.contentType = "application/json"
	}

	if opts.headers == nil {
		opts.headers = make(map[string]string)
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest(opts.method, opts.url, bytes.NewBufferString(opts.body))

	req.Header.Set("Content-Type", opts.contentType)

	for name, value := range opts.headers {
		req.Header.Set(name, value)
	}

	for _, cookie := range opts.cookies {
		req.AddCookie(cookie)
	}

	for _, param := range opts.params {
		ctx.Params = append(opts.params, param)
	}

	ctx.Request = req

	return ctx
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return r
}

func TestExtractFieldValFromBody(t *testing.T) {
	tests := []struct {
		name  string
		opts  ginCtxReqOpts
		field string

		value string
		err   error
	}{
		// json extraction
		{name: "Valid json extraction", field: "name", opts: ginCtxReqOpts{body: `{"name":"John"}`, contentType: "application/json"}, value: `John`, err: nil},
		{name: "Nested field extraction", field: "user.name", opts: ginCtxReqOpts{body: `{"user": {"name":"John"}}`, contentType: "application/json"}, value: `John`, err: nil},
		{name: "Missing field", field: "age", opts: ginCtxReqOpts{body: `{"name":"John"}`, contentType: "application/json"}, value: ``, err: nil},
		{name: "Incorrect JSON format", field: "name", opts: ginCtxReqOpts{body: `{"name":"Johnathan"`, contentType: "application/json"}, value: ``, err: ErrExtractionInvalidJson},
		{name: "Array extraction", field: "names.0", opts: ginCtxReqOpts{body: `{"names": ["John", "Doe"]}`, contentType: "application/json"}, value: `John`, err: nil},
		{name: "Deeply nested field extraction", field: "a.b.c.d.e", opts: ginCtxReqOpts{body: `{"a": {"b": {"c": {"d": {"e": "value"}}}}}`, contentType: "application/json"}, value: `value`, err: nil},
		{name: "Field extraction with numeric values", field: "age", opts: ginCtxReqOpts{body: `{"name":"John", "age": 30}`, contentType: "application/json"}, value: `30`, err: nil},
		{name: "Empty JSON object", field: "name", opts: ginCtxReqOpts{body: `{}`, contentType: "application/json"}, value: ``, err: nil},
		{name: "Null field value", field: "name", opts: ginCtxReqOpts{body: `{"name": null}`, contentType: "application/json"}, value: ``, err: nil},

		// // x-www-form-urlencoded extraction
		{name: "Valid x-www-form-urlencoded extraction (name)", field: "name", opts: ginCtxReqOpts{body: `name=John`, contentType: "application/x-www-form-urlencoded"}, value: `John`, err: nil},
		{name: "Valid x-www-form-urlencoded extraction (age)", field: "age", opts: ginCtxReqOpts{body: `name=John&age=30`, contentType: "application/x-www-form-urlencoded"}, value: `30`, err: nil},
		{name: "Valid x-www-form-urlencoded extraction (email)", field: "email", opts: ginCtxReqOpts{body: `name=John&age=30&email=john@example.com`, contentType: "application/x-www-form-urlencoded"}, value: `john@example.com`, err: nil},
		{name: "Invalid x-www-form-urlencoded extraction (missing field)", field: "address", opts: ginCtxReqOpts{body: `name=John&age=30`, contentType: "application/x-www-form-urlencoded"}, value: ``, err: nil},
		{name: "Valid x-www-form-urlencoded extraction (special characters)", field: "name", opts: ginCtxReqOpts{body: `name=John%20Doe&age=30`, contentType: "application/x-www-form-urlencoded"}, value: `John Doe`, err: nil},
		{name: "Valid x-www-form-urlencoded extraction (numeric value as string)", field: "age", opts: ginCtxReqOpts{body: `name=John&age=42`, contentType: "application/x-www-form-urlencoded"}, value: `42`, err: nil},

		// multipart/form-data extraction
		{name: "Valid multipart/form-data extraction (name)", field: "name", opts: ginCtxReqOpts{body: `--boundary\r\nContent-Disposition: form-data; name="name"\r\n\r\nJohn\r\n--boundary--\r\n`, contentType: "multipart/form-data; boundary=--------------------------590299136414163472038474"}, value: `John`, err: nil},
		{name: "Valid multipart/form-data extraction (name)", field: "name", opts: ginCtxReqOpts{body: `name=John`, contentType: "multipart/form-data; boundary=--------------------------590299136414163472038474"}, value: `John`, err: nil},
		{name: "Valid multipart/form-data extraction (age)", field: "age", opts: ginCtxReqOpts{body: `age=31`, contentType: "multipart/form-data; boundary=--------------------------590299136414163472038474"}, value: `30`, err: nil},
		{name: "Valid multipart/form-data extraction (email)", field: "email", opts: ginCtxReqOpts{body: `email=ohn@example.com`, contentType: "multipart/form-data; boundary=--------------------------590299136414163472038474"}, value: `john@example.com`, err: nil},
		{name: "Invalid multipart/form-data extraction (missing field)", field: "address", opts: ginCtxReqOpts{body: `address=`, contentType: "multipart/form-data; boundary=--------------------------590299136414163472038474"}, value: ``, err: nil},
		{name: "Valid multipart/form-data extraction (special characters)", field: "name", opts: ginCtxReqOpts{body: `name=John%20Doe`, contentType: "multipart/form-data; boundary=--------------------------590299136414163472038474"}, value: `John%20Doe`, err: nil},
		{name: "Valid multipart/form-data extraction (numeric value as string)", field: "age", opts: ginCtxReqOpts{body: `age=42`, contentType: "multipart/form-data; boundary=--------------------------590299136414163472038474"}, value: `42`, err: nil},

		{name: "Valid multipart/form-data extraction (name)", field: "name", opts: ginCtxReqOpts{body: `name=John`, contentType: "multipart/form-data; boundary=--------------------------590299136414163472038474"}, value: `John`, err: nil},
		// {name: "Valid multipart/form-data extraction (name)", field: "name", opts: ginCtxReqOpts{body: `--boundary\r\nContent-Disposition: form-data; name="name"\r\n\r\nJohn\r\n--boundary--\r\n`, contentType: "multipart/form-data; boundary=boundary"}, value: `John`, err: nil},
		// {name: "Valid multipart/form-data extraction (age)", field: "age", opts: ginCtxReqOpts{body: `--boundary\r\nContent-Disposition: form-data; name="age"\r\n\r\n31\r\n--boundary--\r\n`, contentType: "multipart/form-data; boundary=boundary"}, value: `31`, err: nil},
		// {name: "Valid multipart/form-data extraction (email)", field: "email", opts: ginCtxReqOpts{body: `--boundary\r\nContent-Disposition: form-data; name="email"\r\n\r\njohn@example.com\r\n--boundary--\r\n`, contentType: "multipart/form-data; boundary=boundary"}, value: `john@example.com`, err: nil},
		// {name: "Invalid multipart/form-data extraction (missing field)", field: "address", opts: ginCtxReqOpts{body: `--boundary\r\nContent-Disposition: form-data; name="address"\r\n\r\n\r\n--boundary--\r\n`, contentType: "multipart/form-data; boundary=boundary"}, value: ``, err: nil},
		// {name: "Valid multipart/form-data extraction (special characters)", field: "name", opts: ginCtxReqOpts{body: `--boundary\r\nContent-Disposition: form-data; name="name"\r\n\r\nJohn%20Doe\r\n--boundary--\r\n`, contentType: "multipart/form-data; boundary=boundary"}, value: `John Doe`, err: nil},
		// {name: "Valid multipart/form-data extraction (numeric value as string)", field: "age", opts: ginCtxReqOpts{body: `--boundary\r\nContent-Disposition: form-data; name="age"\r\n\r\n42\r\n--boundary--\r\n`, contentType: "multipart/form-data; boundary=boundary"}, value: `42`, err: nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := createTestGinCtx(test.opts)
			ans, err := extractFieldValFromBody(test.field, ctx)

			if err != nil {
				if !errors.Is(test.err, test.err) {
					t.Errorf("got error %+v, want %+v", err, test.err)
				}

				return
			}

			if ans != test.value {
				t.Errorf("got %q, want %q", ans, test.value)
			}
		})
	}
}

func TestExtractFieldValFromCookie(t *testing.T) {
	tests := []struct {
		name  string
		opts  ginCtxReqOpts
		field string

		value string
		err   error
	}{
		{name: "Valid cookie extraction", field: "name", opts: ginCtxReqOpts{cookies: []*http.Cookie{{Name: "name", Value: "John"}}}, value: `John`, err: nil},
		{name: "Valid cookie extraction (multiple cookies)", field: "session", opts: ginCtxReqOpts{cookies: []*http.Cookie{{Name: "session", Value: "abc123"}}}, value: `abc123`, err: nil},
		{name: "Valid cookie extraction (empty cookie)", field: "empty", opts: ginCtxReqOpts{cookies: []*http.Cookie{{Name: "empty", Value: ""}}}, value: ``, err: nil},
		{name: "Invalid cookie extraction (missing cookie)", field: "missing", opts: ginCtxReqOpts{cookies: nil}, value: ``, err: nil},
		{name: "Valid cookie extraction (special characters)", field: "name", opts: ginCtxReqOpts{cookies: []*http.Cookie{{Name: "name", Value: "John Doe"}}}, value: `John Doe`, err: nil},
		{name: "Valid cookie extraction (numeric value)", field: "age", opts: ginCtxReqOpts{cookies: []*http.Cookie{{Name: "age", Value: "42"}}}, value: `42`, err: nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := createTestGinCtx(test.opts)
			ans, err := extractFieldValFromCookie(test.field, ctx)

			if err != nil {
				if !errors.Is(test.err, test.err) {
					t.Errorf("got error %+v, want %+v", err, test.err)
				}
			}

			if ans != test.value {
				t.Errorf("got %q, want %q", ans, test.value)
			}
		})
	}
}

func TestExtractFieldValFromHeader(t *testing.T) {
	tests := []struct {
		name  string
		opts  ginCtxReqOpts
		field string

		value string
		err   error
	}{
		{name: "Valid header extraction", field: "name", opts: ginCtxReqOpts{headers: map[string]string{"name": "John"}}, value: `John`, err: nil},
		{name: "Valid header extraction (multiple headers)", field: "session", opts: ginCtxReqOpts{headers: map[string]string{"session": "abc123"}}, value: `abc123`, err: nil},
		{name: "Valid header extraction (empty header)", field: "empty", opts: ginCtxReqOpts{headers: map[string]string{"empty": ""}}, value: ``, err: nil},
		{name: "Invalid header extraction (missing header)", field: "missing", opts: ginCtxReqOpts{headers: map[string]string{}}, value: ``, err: nil},
		{name: "Valid header extraction (special characters)", field: "name", opts: ginCtxReqOpts{headers: map[string]string{"name": "John%20Doe"}}, value: `John Doe`, err: nil},
		{name: "Valid header extraction (numeric value)", field: "age", opts: ginCtxReqOpts{headers: map[string]string{"age": "42"}}, value: `42`, err: nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := createTestGinCtx(test.opts)
			ans, err := extractFieldValFromHeader(test.field, ctx)

			if err != nil {
				if !errors.Is(test.err, test.err) {
					t.Errorf("got error %+v, want %+v", err, test.err)
				}
			}

			if ans != test.value {
				t.Errorf("got %q, want %q", ans, test.value)
			}
		})
	}
}

func TestExtractFieldValFromParam(t *testing.T) {
	tests := []struct {
		name  string
		opts  ginCtxReqOpts
		field string

		value string
		err   error
	}{
		{name: "Valid header extraction", field: "name", opts: ginCtxReqOpts{headers: map[string]string{"name": "John"}, params: gin.Params{gin.Param{Key: "name", Value: "John"}}}, value: `John`, err: nil},
		{name: "Valid header extraction (multiple headers)", field: "session", opts: ginCtxReqOpts{headers: map[string]string{"session": "abc123"}, params: gin.Params{gin.Param{Key: "session", Value: "abc123"}}}, value: `abc123`, err: nil},
		{name: "Valid header extraction (empty header)", field: "empty", opts: ginCtxReqOpts{headers: map[string]string{"empty": ""}, params: gin.Params{gin.Param{Key: "empty", Value: ""}}}, value: ``, err: nil},
		{name: "Invalid header extraction (missing header)", field: "missing", opts: ginCtxReqOpts{headers: map[string]string{}, params: gin.Params{}}, value: ``, err: nil},
		{name: "Valid header extraction (special characters)", field: "name", opts: ginCtxReqOpts{headers: map[string]string{"name": "John%20Doe"}, params: gin.Params{gin.Param{Key: "name", Value: "John%20Doe"}}}, value: `John Doe`, err: nil},
		{name: "Valid header extraction (numeric value)", field: "age", opts: ginCtxReqOpts{headers: map[string]string{"age": "42"}, params: gin.Params{gin.Param{Key: "age", Value: "42"}}}, value: `42`, err: nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := createTestGinCtx(test.opts)
			ans, err := extractFieldValFromParam(test.field, ctx)

			if err != nil {
				if !errors.Is(test.err, test.err) {
					t.Errorf("got error %+v, want %+v", err, test.err)
				}
			}

			if ans != test.value {
				t.Errorf("got %q, want %q", ans, test.value)
			}
		})
	}
}

func TestExtractFieldValFromQuery(t *testing.T) {
	tests := []struct {
		name  string
		opts  ginCtxReqOpts
		field string

		value string
		err   error
	}{
		{name: "Valid query extraction", field: "name", opts: ginCtxReqOpts{url: "/test?name=John"}, value: `John`, err: nil},
		{name: "Valid query extraction (multiple queries)", field: "session", opts: ginCtxReqOpts{url: "/test?session=abc123"}, value: `abc123`, err: nil},
		{name: "Valid query extraction (empty query)", field: "empty", opts: ginCtxReqOpts{url: "/test?empty="}, value: ``, err: nil},
		{name: "Invalid query extraction (missing query)", field: "missing", opts: ginCtxReqOpts{url: "/test"}, value: ``, err: nil},
		{name: "Valid query extraction (special characters)", field: "name", opts: ginCtxReqOpts{url: "/test?name=John%20Doe"}, value: `John Doe`, err: nil},
		{name: "Valid query extraction (numeric value)", field: "age", opts: ginCtxReqOpts{url: "/test?age=42"}, value: `42`, err: nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := createTestGinCtx(test.opts)
			ans, err := extractFieldValFromQuery(test.field, ctx)

			if err != nil {
				if !errors.Is(test.err, test.err) {
					t.Errorf("got error %+v, want %+v", err, test.err)
				}
			}

			if ans != test.value {
				t.Errorf("got %q, want %q", ans, test.value)
			}
		})
	}
}
