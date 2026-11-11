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

// Get retrieves a specific field's value from a given request location within MatchedData.
// Parameters:
//   - loc: The request location to search in (e.g., "body", "cookies", "headers", "params", "queries").
//   - field: The name of the field to retrieve.
//
// Returns:
//   - The value associated with the specified field at the given location (as a string).
//   - A boolean indicating if the field exists (true if found, false if not).
func (md MatchedData) Get(loc RequestLocation, field string) (string, bool) {
	value, ok := md[loc.String()][field]
	return value, ok
}

// Has reports whether a field exists at the given request location.
func (md MatchedData) Has(loc RequestLocation, field string) bool {
	_, ok := md[loc.String()][field]
	return ok
}

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

// matchedDataStore returns the request's matched data store, creating an empty
// one if no chain has stored anything yet.
func matchedDataStore(ctx *gin.Context) MatchedData {
	if data, ok := ctx.Get(GinValidatorCtxMatchedDataStoreName); ok {
		if store, ok := data.(MatchedData); ok && store != nil {
			return store
		}
	}
	store := make(MatchedData)
	ctx.Set(GinValidatorCtxMatchedDataStoreName, store)
	return store
}

// saveMatchedDataToCtx saves a validated, sanitized value under its location and field.
func saveMatchedDataToCtx(ctx *gin.Context, location, field, value string) {
	if ctx == nil {
		return
	}

	store := matchedDataStore(ctx)
	if store[location] == nil {
		store[location] = make(MatchedDataFieldValues)
	}
	store[location][field] = value
}

// saveChainResultToCtx records a chain's errors and, if the field was sent and
// passed every rule, its final value. Failed and missing fields are left out of
// the matched data, so Get and Has only report values that are safe to use.
//
// The matched data store is created even when nothing is saved, so
// GetMatchedData works for every request a chain has run on.
func saveChainResultToCtx(ctx *gin.Context, result chainResult) {
	if ctx == nil {
		return
	}

	saveValidationErrorsToCtx(ctx, result.errors)
	matchedDataStore(ctx)
	if result.matched {
		saveMatchedDataToCtx(ctx, result.location, result.field, result.sanitizedValue)
	}
}
