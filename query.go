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
			field:           q.field,
			errorMessage:    q.errorMessage,
			location:        queryLocation,
			rules:           make(validationChainRules, 0),
			chainMethodType: validatorType,
		},
		modifier: modifier{
			field:           q.field,
			errorMessage:    q.errorMessage,
			location:        queryLocation,
			rules:           make(validationChainRules, 0),
			chainMethodType: modifierType,
		},
		sanitizer: sanitizer{
			field:           q.field,
			errorMessage:    q.errorMessage,
			location:        queryLocation,
			rules:           make(validationChainRules, 0),
			chainMethodType: sanitizerType,
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
