package ginvalidator

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrNilCtxMatchedData = errors.New("cannot get matched data: nil context provided")
	ErrNoMatchedData     = errors.New("no matched data present")
)

const GinValidatorSanitizedDataStore string = "__ginvalidator__sanitized__data__"

type sanitizedFields map[string]string
type SanitizedData map[string]sanitizedFields

func MatchedData(ctx *gin.Context) (SanitizedData, error) {
	if ctx == nil {
		return nil, ErrNilCtxMatchedData
	}

	data, ok := ctx.Get(GinValidatorSanitizedDataStore)

	if !ok {
		return nil, ErrNoMatchedData
	}

	var store SanitizedData
	store, ok = data.(SanitizedData)

	if !ok {
		return nil, ErrNoMatchedData
	}

	// fmt.Println("STORE:", store)
	return store, nil
}

func createSanitizedDataStore(ctx *gin.Context) {
	var newStore SanitizedData

	ctx.Set(GinValidatorSanitizedDataStore, newStore)
}

func saveSanitizedDataToCtx(ctx *gin.Context, location, field, value string) {
	if ctx == nil {
		return
	}

	data, ok := ctx.Get(GinValidatorSanitizedDataStore)

	if !ok {
		// fmt.Println("sanitization store dne, starting to save errs")
		createSanitizedDataStore(ctx)
		saveSanitizedDataToCtx(ctx, location, field, value)
		return
	}

	var store SanitizedData
	store, ok = data.(SanitizedData)

	if !ok {
		// fmt.Println("sanitization store exists but is wrong type")
		createSanitizedDataStore(ctx)
		saveSanitizedDataToCtx(ctx, location, field, value)
		return
	}

	if store == nil {
		// fmt.Println("sanitization store is nil")
		store = make(SanitizedData)
	}

	// fmt.Println("Save to sanitization store starting")

	specificLocationStore, ok := store[location]

	if !ok {
		// fmt.Println("could not get sanitization location, had to set default")
		specificLocationStore = make(sanitizedFields)
		store[location] = specificLocationStore
	}

	specificLocationStore[field] = value

	store[location] = specificLocationStore

	ctx.Set(GinValidatorSanitizedDataStore, store)
	// fmt.Println("Save to sanitization store ending")
}
