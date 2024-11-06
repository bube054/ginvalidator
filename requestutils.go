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

// func extractFieldValFromBody(ctx *gin.Context, field string) (string, error) {
// 	if ctx == nil {
// 		return "", ErrExtractionFromNilCtx
// 	}

// 	contentType := ctx.GetHeader("Content-Type")

// 	if contentType == "application/json" {
// 		data, err := ctx.GetRawData()

// 		if err != nil {
// 			return "", err
// 		}

// 		jsonStr := string(data)
// 		validJson := json.Valid([]byte(jsonStr))

// 		if !validJson {
// 			return "", fmt.Errorf("%s is %w", jsonStr, ErrExtractionInvalidJson)
// 		}

// 		result := gjson.Get(jsonStr, field)

// 		return result.String(), nil
// 	}

// 	if contentType == "application/x-www-form-urlencoded" || strings.HasPrefix(contentType, "multipart/form-data") {
// 		return ctx.PostForm(field), nil
// 	}

// 	return "", fmt.Errorf("%s is %w", contentType, ErrExtractionInvalidContentType)
// }

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
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(data))

		return ctx.PostForm(field), nil
	}

	// Invalid content type
	return "", fmt.Errorf("%s is %w", contentType, ErrExtractionInvalidContentType)
}

func extractFieldValFromCookie(ctx *gin.Context, field string) (string, error) {
	if ctx == nil {
		return "", ErrExtractionFromNilCtx
	}

	cookie, err := ctx.Cookie(field)

	if err != nil {
		return "", err
	}

	cookie, err = url.QueryUnescape(cookie)

	if err != nil {
		return "", err
	}

	return cookie, nil
}

func extractFieldValFromHeader(ctx *gin.Context, field string) (string, error) {
	if ctx == nil {
		return "", ErrExtractionFromNilCtx
	}

	header := ctx.GetHeader(field)

	header, err := url.QueryUnescape(header)

	if err != nil {
		return "", err
	}

	return header, nil
}

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

func extractFieldValFromQuery(ctx *gin.Context, field string) (string, error) {
	if ctx == nil {
		return "", ErrExtractionFromNilCtx
	}

	query := ctx.Query(field)

	query, err := url.QueryUnescape(query)

	if err != nil {
		return "", err
	}

	return query, nil
}
