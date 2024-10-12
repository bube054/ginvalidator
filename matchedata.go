package ginvalidator

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

const GinValidatorSanitizedDataStore string = "__ginvalidator__sanitized__data__"

type sanitizedFields map[string]string
type SanitizedData map[string]sanitizedFields

func MatchedData(ctx *gin.Context) (SanitizedData, error) {
	data, ok := ctx.Get(GinValidatorSanitizedDataStore)

	if !ok {
		return nil, errors.New("no sanitized data yet")
	}

	var store SanitizedData
	store, ok = data.(SanitizedData)

	if !ok {
		return nil, errors.New("no sanitized data yet")
	}

	return store, nil
}

func createNewStore(ctx *gin.Context) {
	var newStore SanitizedData

	ctx.Set(GinValidatorSanitizedDataStore, newStore)
}

func saveSanitizedDataToCtx(ctx *gin.Context, location, field, value string) {
	if ctx == nil {
		return
	}

	data, ok := ctx.Get(GinValidatorSanitizedDataStore)

	if !ok {
		fmt.Println("sanitization store dne, starting to save errs")
		createNewStore(ctx)
		saveSanitizedDataToCtx(ctx, location, field, value)
		return
	}

	var store SanitizedData
	store, ok = data.(SanitizedData)

	if !ok {
		fmt.Println("sanitization store exists but is wrong type")
		createNewStore(ctx)
		saveSanitizedDataToCtx(ctx, location, field, value)
		return
	}

	if store == nil {
		store = make(SanitizedData)
	}

	fmt.Println("Save to sanitization store starting")

	specificLocationStore, ok := store[location]

	if !ok {
		fmt.Println("could not get sanitization location, had to set default")
		specificLocationStore = make(sanitizedFields)
		store[location] = specificLocationStore
	}

	specificLocationStore[field] = value

	store[location] = specificLocationStore

	ctx.Set(GinValidatorSanitizedDataStore, store)

	fmt.Println("Save to sanitization store ending")
}
