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
	ErrFieldExtractionFromNilCtx    = errors.New("gin context is nil")
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
		return "", ErrFieldExtractionFromNilCtx
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
		return "", ErrFieldExtractionFromNilCtx
	}

	header := ctx.GetHeader(field)

	return header, nil
}

// i escaped
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

// auto escaped
func extractFieldValFromQuery(ctx *gin.Context, field string) (string, error) {
	if ctx == nil {
		return "", ErrFieldExtractionFromNilCtx
	}

	query := ctx.Query(field)

	return query, nil
}
