package ginvalidator

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

var (
	// ErrFieldExtractionFromNilCtx occurs when an operation attempts to extract a field from a nil Gin context.
	ErrFieldExtractionFromNilCtx = errors.New("failed to extract field: gin context is nil")

	// ErrExtractionInvalidContentType occurs when the request contains an unsupported or missing Content-Type header.
	ErrExtractionInvalidContentType = errors.New("failed to extract field: unsupported or missing Content-Type header")

	// ErrExtractionInvalidJSON occurs when JSON parsing fails due to malformed JSON in the request body.
	ErrExtractionInvalidJSON = errors.New("failed to extract field: invalid JSON in request body")
)

// RequestLocation defines different locations where data can be extracted from the request.
type RequestLocation int

// Constants representing different locations in a request.
const (
	// BodyLocation represents the request body.
	BodyLocation RequestLocation = iota

	// CookieLocation represents cookies in the request.
	CookieLocation

	// HeaderLocation represents the headers in the request.
	HeaderLocation

	// ParamLocation represents path parameters in the request.
	ParamLocation

	// QueryLocation represents query parameters in the URL of the request.
	QueryLocation
)

// String returns a string representation of the RequestLocation.
func (l RequestLocation) String() string {
	return [...]string{"body", "cookies", "headers", "params", "queries"}[l]
}

type validationChainType int

const (
	validatorType validationChainType = iota
	sanitizerType
	modifierType
)

// ginValidatorCtxBodyStoreName is the gin.Context key for the request body, read
// once per request and shared by every chain.
const ginValidatorCtxBodyStoreName = "__ginvalidator__body__"

// requestBody is the request body as read by the first chain that needed it.
type requestBody struct {
	data []byte

	jsonChecked bool // whether jsonValid has been worked out yet
	jsonValid   bool
}

// validJSON reports whether the body is valid JSON, checking it only once. An
// empty or whitespace-only body counts as valid: it simply has no fields.
func (b *requestBody) validJSON() bool {
	if !b.jsonChecked {
		b.jsonChecked = true
		b.jsonValid = len(bytes.TrimSpace(b.data)) == 0 || gjson.ValidBytes(b.data)
	}
	return b.jsonValid
}

// readRequestBody returns the request body, reading it only the first time it is
// called for a request. The body is put back afterwards so later handlers, and
// ctx.PostForm, can still read it.
func readRequestBody(ctx *gin.Context) (*requestBody, error) {
	if v, ok := ctx.Get(ginValidatorCtxBodyStoreName); ok {
		if body, ok := v.(*requestBody); ok {
			return body, nil
		}
	}

	data, err := ctx.GetRawData()
	if err != nil {
		return nil, err
	}
	ctx.Request.Body = io.NopCloser(bytes.NewReader(data))

	body := &requestBody{data: data}
	ctx.Set(ginValidatorCtxBodyStoreName, body)
	return body, nil
}

func extractFieldValFromBody(ctx *gin.Context, field string) (string, error) {
	if ctx == nil {
		return "", ErrFieldExtractionFromNilCtx
	}

	body, err := readRequestBody(ctx)
	if err != nil {
		return "", err
	}

	contentType := ctx.GetHeader("Content-Type")

	// Compare the media type only, so parameters such as "; charset=utf-8" and
	// differences in case don't matter. An unparsable header leaves it empty.
	mediaType, _, _ := mime.ParseMediaType(contentType)

	// JSON, including structured types such as application/vnd.api+json. gjson does
	// not validate, so the body is checked first: without it, a truncated body like
	// {"email":"a@example.com" still produced values.
	if mediaType == "application/json" || strings.HasSuffix(mediaType, "+json") {
		if !body.validJSON() {
			return "", ErrExtractionInvalidJSON
		}
		return gjson.GetBytes(body.data, field).String(), nil
	}

	if mediaType == "application/x-www-form-urlencoded" || mediaType == "multipart/form-data" {
		return ctx.PostForm(field), nil
	}

	// Invalid content type
	return "", fmt.Errorf("%s is %w", contentType, ErrExtractionInvalidContentType)
}

func extractFieldValFromCookie(ctx *gin.Context, field string) (string, error) {
	if ctx == nil {
		return "", ErrFieldExtractionFromNilCtx
	}

	cookie, err := ctx.Cookie(field)

	// A cookie that wasn't sent is an empty value, like a missing header or query
	// parameter, not an error.
	if errors.Is(err, http.ErrNoCookie) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return cookie, nil
}

func extractFieldValFromHeader(ctx *gin.Context, field string) (string, error) {
	if ctx == nil {
		return "", ErrFieldExtractionFromNilCtx
	}

	header := ctx.GetHeader(field)

	if header == "" {
		return getOriginalHeaderValue(ctx.Request.Header, field), nil
	}

	return header, nil
}

func extractFieldValFromParam(ctx *gin.Context, field string) (string, error) {
	if ctx == nil {
		return "", ErrFieldExtractionFromNilCtx
	}

	// Gin has already decoded the parameter (or, with UseRawPath and
	// UnescapePathValues off, deliberately left it encoded). Decoding it again
	// broke values such as "100%" and turned "+" into a space.
	return ctx.Param(field), nil
}

func extractFieldValFromQuery(ctx *gin.Context, field string) (string, error) {
	if ctx == nil {
		return "", ErrFieldExtractionFromNilCtx
	}

	query := ctx.Query(field)

	return query, nil
}

// fieldPresent reports whether the request actually contains field at loc, as
// opposed to the field being absent, which the extractors also return as "".
func fieldPresent(ctx *gin.Context, loc RequestLocation, field string) bool {
	if ctx == nil || ctx.Request == nil {
		return false
	}

	switch loc {
	case BodyLocation:
		body, err := readRequestBody(ctx)
		if err != nil {
			return false
		}
		mediaType, _, _ := mime.ParseMediaType(ctx.GetHeader("Content-Type"))
		if mediaType == "application/json" || strings.HasSuffix(mediaType, "+json") {
			return body.validJSON() && gjson.GetBytes(body.data, field).Exists()
		}
		_, ok := ctx.GetPostForm(field)
		return ok
	case CookieLocation:
		_, err := ctx.Request.Cookie(field)
		return err == nil
	case HeaderLocation:
		if len(ctx.Request.Header.Values(field)) > 0 {
			return true
		}
		// Like extractFieldValFromHeader, also match keys that were stored without
		// canonical casing, such as ones set directly on the header map.
		for k := range ctx.Request.Header {
			if strings.EqualFold(k, field) {
				return true
			}
		}
		return false
	case ParamLocation:
		_, ok := ctx.Params.Get(field)
		return ok
	case QueryLocation:
		_, ok := ctx.GetQuery(field)
		return ok
	}
	return false
}

// getOriginalHeaderValue finds a header whose key differs from field only in case,
// such as one set directly on the header map without canonical casing.
func getOriginalHeaderValue(headers http.Header, key string) string {
	for k, v := range headers {
		if strings.EqualFold(k, key) && len(v) > 0 {
			return v[0]
		}
	}
	return ""
}
