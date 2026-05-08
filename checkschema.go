package ginvalidator

import (
	"sort"

	"github.com/gin-gonic/gin"
)

// SchemaField describes how a single field should be validated within a [Schema].
type SchemaField struct {
	// In specifies which request location the field comes from (body, cookies,
	// headers, params, or queries).
	In RequestLocation

	// ErrFmtFunc is an optional per-field error message formatter.
	// When nil, the package-level fallback rules apply.
	ErrFmtFunc ErrFmtFuncHandler

	// Optional, when true, skips validation if the field is empty.
	Optional bool

	// Build receives a fresh [ValidationChain] and returns the configured
	// chain with validators, sanitizers, and modifiers attached.
	// Use .Bail() within the Build function to stop on the first failure.
	// If nil the chain runs with no validators (always passes).
	Build func(ValidationChain) ValidationChain
}

// Schema maps field names to their validation configuration.
type Schema map[string]SchemaField

// CheckSchema creates a single [gin.HandlerFunc] from a declarative schema.
// Fields are processed in sorted order for deterministic error ordering.
//
// Example:
//
//	router.POST("/register",
//	  ginvalidator.CheckSchema(ginvalidator.Schema{
//	    "email": {In: ginvalidator.BodyLocation, Build: func(vc ginvalidator.ValidationChain) ginvalidator.ValidationChain {
//	      return vc.Email(nil)
//	    }},
//	    "name": {In: ginvalidator.BodyLocation, Optional: true, Build: func(vc ginvalidator.ValidationChain) ginvalidator.ValidationChain {
//	      return vc.Alpha(nil)
//	    }},
//	  }),
//	  handler,
//	)
func CheckSchema(schema Schema) gin.HandlerFunc {
	fields := make([]string, 0, len(schema))
	for f := range schema {
		fields = append(fields, f)
	}
	sort.Strings(fields)

	return func(ctx *gin.Context) {
		for _, field := range fields {
			sf := schema[field]
			vc := NewValidationChain(field, sf.ErrFmtFunc, sf.In)

			if sf.Optional {
				vc = vc.Optional()
			}

			if sf.Build != nil {
				vc = sf.Build(vc)
			}

			result := vc.validate(ctx)
			saveValidationErrorsToCtx(ctx, result.errors)
			saveMatchedDataToCtx(ctx, result.location, result.field, result.sanitizedValue)
		}
		ctx.Next()
	}
}
