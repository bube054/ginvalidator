package ginvalidator

// request validator struct for request bodies e.g {"key": "value"}, key to be validated
type body struct {
	field        string
	errorMessage string
}

// body creator of the validation chain.
func (b *body) Chain() validationChain {
	return validationChain{
		validator: validator{
			field:        b.field,
			errorMessage: b.errorMessage,
			location:     bodyLocation,
			rules:        make(validationChainRules, 0),
		},
		modifier: modifier{
			field:        b.field,
			errorMessage: b.errorMessage,
			location:     bodyLocation,
			rules:        make(validationChainRules, 0),
		},
		sanitizer: sanitizer{
			field:        b.field,
			errorMessage: b.errorMessage,
			location:     bodyLocation,
			rules:        make(validationChainRules, 0),
		},
	}
}

// the body struct creator function. which takes in the field to be validated and an errorMessage on failure.
func NewBody(field, errorMessage string) body {
	return body{
		field:        field,
		errorMessage: errorMessage,
	}
}
