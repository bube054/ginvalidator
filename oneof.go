package ginvalidator

import (
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

// OneOf runs each group of validation chains and passes if at least one group
// produces no validation errors. If all groups fail, a single error is added
// to the context.
//
// Each argument is a group of chains that must all pass together. The first
// group that passes wins — its matched data is saved and no errors are recorded.
//
// Example:
//
//	router.POST("/login",
//	  ginvalidator.OneOf(
//	    []ginvalidator.ValidationChain{NewBody("email", nil).Chain().Email(nil)},
//	    []ginvalidator.ValidationChain{NewBody("phone", nil).Chain().MobilePhone(nil, "")},
//	  ),
//	  handler,
//	)
func OneOf(chainGroups ...[]ValidationChain) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, group := range chainGroups {
			var groupErrors []ValidationChainError
			var groupResults []chainResult

			for _, chain := range group {
				result := chain.validate(ctx)
				groupErrors = append(groupErrors, result.errors...)
				groupResults = append(groupResults, result)
			}

			if len(groupErrors) == 0 {
				for _, result := range groupResults {
					saveMatchedDataToCtx(ctx, result.location, result.field, result.sanitizedValue)
				}
				ctx.Next()
				return
			}
		}

		order := atomic.AddUint64(&globalErrorOrder, 1)
		oneOfErr := newValidationChainError(
			vceWithLocation(""),
			vceWithMessage("No group in OneOf passed validation"),
			vceWithField("_oneOf"),
			vceWithValue(""),
			vceWithOrder(order),
		)
		saveValidationErrorsToCtx(ctx, []ValidationChainError{oneOfErr})
		ctx.Next()
	}
}
