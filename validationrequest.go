package ginvalidator

import "github.com/gin-gonic/gin"

type ValidationResult struct {
}

func NewValidationResult(ctx *gin.Context) (results []ReturnableValidationChainResponse) {
	// for body
	value, exists := ctx.Get(bodyLocationStore)
	store, ok := value.(CtxStore)

	if exists && ok {
		for field, responses := range store {
			for _, response := range responses {
				if !response.isValid {
					results = append(results, ReturnableValidationChainResponse{location: bodyLocation, msg: response.msg, path: field, value: response.newValue})
					break
				}
			}
		}
	}

	// for cookies
	value, exists = ctx.Get(cookiesLocationStore)
	store, ok = value.(CtxStore)

	if exists && ok {
		for field, responses := range store {
			for _, response := range responses {
				if !response.isValid {
					results = append(results, ReturnableValidationChainResponse{location: cookiesLocation, msg: response.msg, path: field, value: response.newValue})
					break
				}
			}
		}
	}

	// for headers
	value, exists = ctx.Get(headersLocationStore)
	store, ok = value.(CtxStore)

	if exists && ok {
		for field, responses := range store {
			for _, response := range responses {
				if !response.isValid {
					results = append(results, ReturnableValidationChainResponse{location: headersLocation, msg: response.msg, path: field, value: response.newValue})
					break
				}
			}
		}
	}

	// for params
	value, exists = ctx.Get(paramsLocationStore)
	store, ok = value.(CtxStore)

	if exists && ok {
		for field, responses := range store {
			for _, response := range responses {
				if !response.isValid {
					results = append(results, ReturnableValidationChainResponse{location: paramsLocation, msg: response.msg, path: field, value: response.newValue})
					break
				}
			}
		}
	}

	// for query
	value, exists = ctx.Get(queryLocationStore)
	store, ok = value.(CtxStore)

	if exists && ok {
		for field, responses := range store {
			for _, response := range responses {
				if !response.isValid {
					results = append(results, ReturnableValidationChainResponse{location: bodyLocation, msg: response.msg, path: field, value: response.newValue})
					break
				}
			}
		}
	}

	return
}
