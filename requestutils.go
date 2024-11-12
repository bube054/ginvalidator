package ginvalidator

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"net/http"
	"net/url"

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

func extractFieldValFromBody(ctx *gin.Context, field string) (string, error) {
	if ctx == nil {
		return "", ErrFieldExtractionFromNilCtx
	}

	data, err := ctx.GetRawData()
	if err != nil {
		return "", err
	}

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(data))

	contentType := ctx.GetHeader("Content-Type")

	if contentType == "application/json" {
		jsonStr := string(data)
		result := gjson.Get(jsonStr, field)
		return result.String(), nil
	}

	if contentType == "application/x-www-form-urlencoded" || strings.HasPrefix(contentType, "multipart/form-data") {
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

	param := ctx.Param(field)

	param, err := url.QueryUnescape(param)

	if err != nil {
		return "", err
	}

	return param, nil
}

func extractFieldValFromQuery(ctx *gin.Context, field string) (string, error) {
	if ctx == nil {
		return "", ErrFieldExtractionFromNilCtx
	}

	query := ctx.Query(field)

	return query, nil
}

func getOriginalHeaderValue(headers http.Header, key string) string {
	for k, v := range headers {
		if strings.EqualFold(k, key) {
			canonicalKey := http.CanonicalHeaderKey(key)
			log.Printf("Warning: Non-canonical header key '%s' used. Expected '%s'.", key, canonicalKey)
			return v[0]
		}
	}
	return ""
}
