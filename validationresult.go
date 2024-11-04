package ginvalidator

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrNilCtxValidationResult     = errors.New("nil ctx result")
	ErrNoValidationResult         = errors.New("can not get validation result")
)

// GinValidatorCtxErrorsStoreName is the key, where the validation errors are stored.
const GinValidatorCtxErrorsStoreName string = "__ginvalidator__ctx__errors__"

type ctxFieldErrs map[string][]ValidationChainError
type ctxStoreErrs map[string]ctxFieldErrs

// ValidationResult extracts the validation errors from gin's ctx.
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

	// createErrNewStore(ctx)

	return allErrs, nil
}

func createErrNewStore(ctx *gin.Context) {
	var newStore ctxStoreErrs

	ctx.Set(GinValidatorCtxErrorsStoreName, newStore)
}

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
		// fmt.Println("store exists but is wrong type")
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
			// fmt.Println("could not get location, had to set default")
			specificLocationStore = make(ctxFieldErrs)
			store[location] = specificLocationStore
		}

		currentErrs, ok := specificLocationStore[field]

		if !ok {
			// fmt.Println("could not get errors, had to set default")
			currentErrs = make([]ValidationChainError, 0)
			specificLocationStore[field] = currentErrs
		}

		currentErrs = append(currentErrs, err)

		specificLocationStore[field] = currentErrs

		store[location] = specificLocationStore

		// fmt.Println("Save to store starting")
		ctx.Set(GinValidatorCtxErrorsStoreName, store)
		// fmt.Println("Save to store ending")
	}
}
