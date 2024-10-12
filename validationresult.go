package ginvalidator

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const GinValidatorCtxStoreName string = "__ginvalidator__ctx__errors__"

type fieldErrsMap map[string][]ValidationChainError
type storeErrsMap map[string]fieldErrsMap

func ValidationResult(ctx *gin.Context) ([]ValidationChainError, error) {
	if ctx == nil {
		return nil, nil
	}

	data, ok := ctx.Get(GinValidatorCtxStoreName)

	if !ok {
		fmt.Println("store dne")
		return nil, nil
	}

	var store storeErrsMap
	store, ok = data.(storeErrsMap)

	if !ok {
		fmt.Println("store exists but is wrong type")
		return nil, nil
	}

	var allErrs []ValidationChainError

	for _, locations := range store {
		for _, errs := range locations {
			allErrs = append(allErrs, errs...)
		}
	}

	createNewStore(ctx)

	return allErrs, nil
}

func createNewStore(ctx *gin.Context) {
	var newStore storeErrsMap

	ctx.Set(GinValidatorCtxStoreName, newStore)
}

func saveErrorsToCtx(ctx *gin.Context, errs []ValidationChainError) {
	if ctx == nil {
		return
	}

	data, ok := ctx.Get(GinValidatorCtxStoreName)

	if !ok {
		fmt.Println("store dne, starting to save errs")
		createNewStore(ctx)
		saveErrorsToCtx(ctx, errs)
		return
	}

	var store storeErrsMap
	store, ok = data.(storeErrsMap)

	if !ok {
		fmt.Println("store exists but is wrong type")
		createNewStore(ctx)
		saveErrorsToCtx(ctx, errs)
		return
	}

	if store == nil {
		store = make(storeErrsMap)
	}

	for _, err := range errs {
		field := err.Field
		location := err.Location

		specificLocationStore, ok := store[location]

		if !ok {
			fmt.Println("could not get location, had to set default")
			specificLocationStore = make(fieldErrsMap)
			store[location] = specificLocationStore
		}

		currentErrs, ok := specificLocationStore[field]

		if !ok {
			fmt.Println("could not get errors, had to set default")
			currentErrs = make([]ValidationChainError, 0)
			specificLocationStore[field] = currentErrs
		}

		currentErrs = append(currentErrs, err)

		specificLocationStore[field] = currentErrs

		store[location] = specificLocationStore

		fmt.Println("Save to store starting")
		ctx.Set(GinValidatorCtxStoreName, store)
		fmt.Println("Save to store ending")
	}
}
