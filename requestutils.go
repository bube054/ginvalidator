package ginvalidator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

var (
	ErrExtractionFromNilCtx         = errors.New("gin context is nil")
	ErrExtractionInvalidContentType = errors.New("invalid content-type header")
	ErrExtractionInvalidJson        = errors.New("invalid json")
)

type requestLocation int

const (
	bodyLocation requestLocation = iota
	cookieLocation
	headerLocation
	paramLocation
	queryLocation
)

func (l requestLocation) string() string {
	return [...]string{"body", "cookies", "headers", "params", "query"}[l]
}

type validationChainType int

const (
	validatorType validationChainType = iota
	sanitizerType
	modifierType
)

// not escaped
func extractFieldValFromBody(ctx *gin.Context, field string) (string, error) {
	if ctx == nil {
		return "", ErrExtractionFromNilCtx
	}

	data, err := ctx.GetRawData()
	if err != nil {
		return "", err
	}

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(data))

	contentType := ctx.GetHeader("Content-Type")

	if contentType == "application/json" {
		jsonStr := string(data)

		if !json.Valid(data) {
			return "", fmt.Errorf("%s is %w", jsonStr, ErrExtractionInvalidJson)
		}

		result := gjson.Get(jsonStr, field)
		return result.String(), nil
	}

	if contentType == "application/x-www-form-urlencoded" || strings.HasPrefix(contentType, "multipart/form-data") {
		return ctx.PostForm(field), nil
	}

	// Invalid content type
	return "", fmt.Errorf("%s is %w", contentType, ErrExtractionInvalidContentType)
}

// auto escape
func extractFieldValFromCookie(ctx *gin.Context, field string) (string, error) {
	if ctx == nil {
		return "", ErrExtractionFromNilCtx
	}

	cookie, err := ctx.Cookie(field)

	if err != nil {
		return "", err
	}

	return cookie, nil
}

// not escaped
func extractFieldValFromHeader(ctx *gin.Context, field string) (string, error) {
	if ctx == nil {
		return "", ErrExtractionFromNilCtx
	}

	header := ctx.GetHeader(field)

	return header, nil
}

// i escaped
func extractFieldValFromParam(ctx *gin.Context, field string) (string, error) {
	if ctx == nil {
		return "", ErrExtractionFromNilCtx
	}

	param := ctx.Param(field)

	param, err := url.QueryUnescape(param)

	if err != nil {
		return "", err
	}

	return param, nil
}

// auto escaped
func extractFieldValFromQuery(ctx *gin.Context, field string) (string, error) {
	if ctx == nil {
		return "", ErrExtractionFromNilCtx
	}

	query := ctx.Query(field)

	return query, nil
}

// Body (JSON, URL-encoded, multipart):

// JSON: Values are already parsed and extracted from JSON, so no URL decoding is needed.
// URL-encoded: Express’s express.urlencoded() middleware automatically URL-decodes data if it’s encoded in application/x-www-form-urlencoded format.
// Multipart: Multipart form data (e.g., file uploads) is handled separately, typically by using a package like multer. For URL-decoding, you may need to apply custom decoding logic if there are URL-encoded values within the multipart data fields.
// Headers:

// Headers are usually raw strings, and express-validator doesn’t apply URL decoding. For fields where URL encoding is expected (like custom headers with encoded values), you’ll need to manually decode them.
// Cookies:

// Cookies are extracted as raw strings, with no URL decoding applied. If cookies are URL-encoded, you’ll need to decode them manually.
// Query:

// Query string parameters are typically URL-decoded by Express’s express.query() middleware, so express-validator should receive decoded values.