package ginvalidator

// request validator struct for headers e.g e.g ?person=david, validate david.
type query struct {
	field        string
	errorMessage string
}

// query creator of the validation chain.
func (q *query) Chain() validationChain {
	return validationChain{
		validator: validator{
			field:        q.field,
			errorMessage: q.errorMessage,
			location:     queryLocation,
			rules:        make(validationChainRules, 0),
		},
		modifier: modifier{
			field:        q.field,
			errorMessage: q.errorMessage,
			location:     queryLocation,
			rules:        make(validationChainRules, 0),
		},
		sanitizer: sanitizer{
			field:        q.field,
			errorMessage: q.errorMessage,
			location:     queryLocation,
			rules:        make(validationChainRules, 0),
		},
	}
}

// the query struct creator function. which takes in the field to be validated and an errorMessage on failure.
func NewQuery(field, errorMessage string) query {
	return query{
		field:        field,
		errorMessage: errorMessage,
	}
}
