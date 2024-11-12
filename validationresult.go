package ginvalidator

import (
	"errors"
	"sort"

	"github.com/gin-gonic/gin"
)

var (
	// ErrNilCtxValidationResult is returned when a nil context is provided, making it impossible to extract validation results.
	ErrNilCtxValidationResult = errors.New("nil context provided: unable to extract validation result")

	// ErrNoValidationResult is returned when no validation result is found in the context.
	ErrNoValidationResult = errors.New("validation result not found in context")
)

// GinValidatorCtxErrorsStoreName is the key, where the validation errors are stored.
const GinValidatorCtxErrorsStoreName string = "__ginvalidator__ctx__errors__"

// ctxFieldErrs represents a map where the key is the name of a field and the value is a slice of
// ValidationChainError structs. Each slice holds validation errors associated with that specific field
// in the request.
//
// The map structure allows storing multiple validation errors for each field, helping to track errors
// encountered during validation in different parts of the request (e.g., "body", "cookies", "headers", "params", "queries").
type ctxFieldErrs map[string][]ValidationChainError

// ctxStoreErrs represents a map where the key is a location in the request (e.g., "body", "cookies", "headers", "params", "queries")
// and the value is a ctxFieldErrs, which is a map of fields to their associated validation errors.
//
// This structure allows organizing validation errors by request location, making it easier to track and
// handle errors from different parts of the request context.
type ctxStoreErrs map[string]ctxFieldErrs

// ValidationResult extracts the validation errors from the Gin context.
// It retrieves any validation errors that have occurred during the request processing,
// and returns them as a slice of ValidationChainError structs along with any potential error.
//
// Parameters:
//   - ctx: The Gin context, which provides access to the HTTP request and response, including validation error data.
//
// Returns:
//   - A slice of ValidationChainError: Contains the details of each validation error encountered, including location, field, and message.
//   - error: Returns an error if there is an issue extracting or processing the validation errors; otherwise, nil.
func ValidationResult(ctx *gin.Context) ([]ValidationChainError, error) {
	if ctx == nil {
		return nil, ErrNilCtxValidationResult
	}

	data, ok := ctx.Get(GinValidatorCtxErrorsStoreName)

	if !ok {
		return nil, ErrNoValidationResult
	}

	var store ctxStoreErrs
	store, ok = data.(ctxStoreErrs)

	if !ok {
		return nil, ErrNoValidationResult
	}

	var allErrs []ValidationChainError

	for _, locations := range store {
		for _, errs := range locations {
			allErrs = append(allErrs, errs...)
		}
	}

	// allErrs = RandomizeErrors(allErrs)

	SortValidationErrors(allErrs)

	return allErrs, nil
}

// createErrNewStore initializes an empty ctxStoreErrs store and adds it to the context
// under the key specified by GinValidatorCtxErrorsStoreName.
func createErrNewStore(ctx *gin.Context) {
	var newStore ctxStoreErrs

	ctx.Set(GinValidatorCtxErrorsStoreName, newStore)
}

// saveValidationErrorsToCtx saves validation errors into the Gin context.
func saveValidationErrorsToCtx(ctx *gin.Context, errs []ValidationChainError) {
	if ctx == nil {
		return
	}

	data, ok := ctx.Get(GinValidatorCtxErrorsStoreName)

	if !ok {
		createErrNewStore(ctx)
		saveValidationErrorsToCtx(ctx, errs)
		return
	}

	var store ctxStoreErrs
	store, ok = data.(ctxStoreErrs)

	if !ok {
		createErrNewStore(ctx)
		saveValidationErrorsToCtx(ctx, errs)
		return
	}

	if store == nil {
		store = make(ctxStoreErrs)
	}

	for _, err := range errs {
		field := err.Field
		location := err.Location

		specificLocationStore, ok := store[location]

		if !ok {
			specificLocationStore = make(ctxFieldErrs)
			store[location] = specificLocationStore
		}

		currentErrs, ok := specificLocationStore[field]

		if !ok {
			currentErrs = make([]ValidationChainError, 0)
			specificLocationStore[field] = currentErrs
		}

		currentErrs = append(currentErrs, err)

		specificLocationStore[field] = currentErrs

		store[location] = specificLocationStore

		ctx.Set(GinValidatorCtxErrorsStoreName, store)
	}
}

func SortValidationErrors(errors []ValidationChainError) {
	sort.Slice(errors, func(i, j int) bool {
		if errors[i].createdAt.Before(errors[j].createdAt) {
			return true
		}
		if errors[i].createdAt.Equal(errors[j].createdAt) {
			return errors[i].incId > errors[j].incId
		}
		return false
	})
}

// func randomizeErrors(errors []ValidationChainError) []ValidationChainError {
// 	rand.Seed(time.Now().UnixNano()) // Seed random number generator with current timep
// 	for i := len(errors) - 1; i > 0; i-- {
// 		j := rand.Intn(i + 1)
// 		errors[i], errors[j] = errors[j], errors[i] // Swap elements
// 	}
// 	return errors
// }

