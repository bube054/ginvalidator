package ginvalidator

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	// ErrNilCtxMatchedData is returned when a nil context is passed, preventing extraction of matched data.
	ErrNilCtxMatchedData = errors.New("nil context provided: unable to extract matched data")

	// ErrNoMatchedData is returned when no matched data is found in the context.
	ErrNoMatchedData = errors.New("no matched data available in context")
)

const GinValidatorCtxMatchedDataStoreName string = "__ginvalidator__matched__data__"

// MatchedDataFieldValues is a map of fields and their values for a request location.
type MatchedDataFieldValues map[string]string

// MatchedData is a map of request locations and fields.
// The keys in MatchedData represent the request locations where fields can be found.
// Possible locations include:
//   - "body": Data from the request body.
//   - "cookies": Data from request cookies.
//   - "headers": Data from request headers.
//   - "params": Data from URL parameters.
//   - "queries": Data from URL query parameters.
type MatchedData map[string]MatchedDataFieldValues

// GetMatchedData extracts and returns matched data from various locations in the request context.
// It retrieves fields and values from predefined request locations such as query parameters, body,
// URL parameters, and headers.
//
// Parameters:
//   - ctx: The Gin context, which provides access to the HTTP request and response.
//
// Returns:
//   - MatchedData: A map containing fields and their values organized by request location.
//   - error: An error if there was an issue extracting data from the context; otherwise, nil.
func GetMatchedData(ctx *gin.Context) (MatchedData, error) {
	if ctx == nil {
		return nil, ErrNilCtxMatchedData
	}

	data, ok := ctx.Get(GinValidatorCtxMatchedDataStoreName)

	if !ok {
		return nil, ErrNoMatchedData
	}

	var store MatchedData
	store, ok = data.(MatchedData)

	if !ok {
		return nil, ErrNoMatchedData
	}

	return store, nil
}

// createMatchedDataStore initializes an empty MatchedData store and adds it to the context
// under the key specified by GinValidatorCtxMatchedDataStoreName.
func createMatchedDataStore(ctx *gin.Context) {
	var newStore MatchedData

	ctx.Set(GinValidatorCtxMatchedDataStoreName, newStore)
}

// saveMatchedDataToCtx saves validated/sanitized data into the Gin context under the specified location and field.
func saveMatchedDataToCtx(ctx *gin.Context, location, field, value string) {
	if ctx == nil {
		return
	}

	data, ok := ctx.Get(GinValidatorCtxMatchedDataStoreName)

	if !ok {
		createMatchedDataStore(ctx)
		saveMatchedDataToCtx(ctx, location, field, value)
		return
	}

	var store MatchedData
	store, ok = data.(MatchedData)

	if !ok {
		createMatchedDataStore(ctx)
		saveMatchedDataToCtx(ctx, location, field, value)
		return
	}

	if store == nil {
		store = make(MatchedData)
	}

	specificLocationStore, ok := store[location]

	if !ok {
		specificLocationStore = make(MatchedDataFieldValues)
		store[location] = specificLocationStore
	}

	specificLocationStore[field] = value
	store[location] = specificLocationStore

	ctx.Set(GinValidatorCtxMatchedDataStoreName, store)
}
