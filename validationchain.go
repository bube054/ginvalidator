package ginvalidator

import (
	"fmt"
	"io"

	// "io"
	"log"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
)

// the chain of either validators, sanitizers and modifiers
type validationChain struct {
	validator
	modifier
	sanitizer
}

func (vc validationChain) Validate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		field := vc.getFieldToValidate()
		sanitizedValue, err := vc.getValueToValidate(ctx)
		responses := make([]validationChainResponse, 0)

		if err != nil {
			log.Println("a validation error has occurred:", err)
			ctx.Next()
		}

		rules := vc.getValidationRules()
		for i, rules := range rules {
			response := rules(sanitizedValue, field, ctx)

			if i > 0 {
				previousResponse := responses[len(responses)-1]
				if !previousResponse.isValid && response.shouldBail {
					break
				}
			}

			if response.shouldBail {
				break
			}

			sanitizedValue = response.newValue
			responses = append(responses, response)
		}

		vc.saveResponsesToStore(ctx, responses)

		// fmt.Printf("responses: %+v\n", responses)
		ctx.Next()
	}
}

// gets the validation rules
func (vc validationChain) getValidationRules() validationChainRules {
	return vc.validator.rules // could have also access modifier or sanitizer.
}

// gets the validation location
func (vc validationChain) getValidationLocation() string {
	return vc.validator.location // could have also access modifier or sanitizer.
}

// gets the validation field
func (vc validationChain) getFieldToValidate() string {
	return vc.validator.field // could have also access modifier or sanitizer.
}

// get the field to be validated from validation location
func (vc validationChain) getValueToValidate(ctx *gin.Context) (string, error) {
	field := vc.getFieldToValidate()
	switch vc.getValidationLocation() {
	case bodyLocation:
		keys, err := splitJSONFieldSelector(field)
		if err != nil {
			return "", err
		}
		reqBody, err := ctx.Request.GetBody()
		if err != nil {
			return "", err
		}
		reqBodyBytes, err := io.ReadAll(reqBody)
		if err != nil {
			return "", err
		}
		key, _, _, err := jsonparser.Get(reqBodyBytes, keys...)
		if err != nil {
			return "", err
		}
		return string(key), nil
	case cookiesLocation:
		return ctx.Cookie(field)
	case headersLocation:
		return ctx.GetHeader(field), nil
	case paramsLocation:
		return ctx.Param(field), nil
	case queryLocation:
		return ctx.Query(field), nil
	default:
		return "", fmt.Errorf("invalid request location for %s", field)
	}
}

func (vc validationChain) getStoreName() string {
	switch vc.getValidationLocation() {
	case bodyLocation:
		return bodyLocationStore
	case cookiesLocation:
		return cookiesLocationStore
	case headersLocation:
		return headersLocationStore
	case paramsLocation:
		return paramsLocationStore
	case queryLocation:
		return queryLocationStore
	default:
		return ""
	}
}

type CtxStore map[string][]validationChainResponse

func (vc validationChain) saveResponsesToStore(ctx *gin.Context, responses []validationChainResponse) {
	defaultCtxStore := make(CtxStore)
	storeName := vc.getStoreName()
	field := vc.getFieldToValidate()

	value, exists := ctx.Get(storeName)

	if !exists {
		defaultCtxStore[field] = responses
		ctx.Set(storeName, defaultCtxStore)
	} else {
		store, ok := value.(CtxStore)
		if !ok {
			store = make(CtxStore)
		}
		store[field] = responses
		ctx.Set(storeName, store)
	}
}
