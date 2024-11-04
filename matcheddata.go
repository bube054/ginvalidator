package ginvalidator

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrNilCtxMatchedData = errors.New("cannot get matched data: nil context provided")
	ErrNoMatchedData     = errors.New("no matched data present")
)

const GinValidatorCtxMatchedDataStoreName string = "__ginvalidator__matched__data__"

type matchedDataFieldValues map[string]string
type MatchedData map[string]matchedDataFieldValues

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

func createSanitizedDataStore(ctx *gin.Context) {
	var newStore MatchedData

	ctx.Set(GinValidatorCtxMatchedDataStoreName, newStore)
}

func saveSanitizedDataToCtx(ctx *gin.Context, location, field, value string) {
	if ctx == nil {
		return
	}

	data, ok := ctx.Get(GinValidatorCtxMatchedDataStoreName)

	if !ok {
		createSanitizedDataStore(ctx)
		saveSanitizedDataToCtx(ctx, location, field, value)
		return
	}

	var store MatchedData
	store, ok = data.(MatchedData)

	if !ok {
		createSanitizedDataStore(ctx)
		saveSanitizedDataToCtx(ctx, location, field, value)
		return
	}

	if store == nil {
		store = make(MatchedData)
	}

	specificLocationStore, ok := store[location]

	if !ok {
		specificLocationStore = make(matchedDataFieldValues)
		store[location] = specificLocationStore
	}

	specificLocationStore[field] = value
	store[location] = specificLocationStore

	ctx.Set(GinValidatorCtxMatchedDataStoreName, store)
}
