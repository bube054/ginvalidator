package ginvalidator

// request validator struct for dynamic parameters e.g /playable-characters/jax. jax to be validated.
type param struct {
	field        string
	errorMessage string
}

// param creator of the validation chain.
func (p *param) Chain() validationChain {
	return validationChain{
		validator: validator{
			field:        p.field,
			errorMessage: p.errorMessage,
			location:     paramsLocation,
			rules:        make(validationChainRules, 0),
		},
		modifier: modifier{
			field:        p.field,
			errorMessage: p.errorMessage,
			location:     paramsLocation,
			rules:        make(validationChainRules, 0),
		},
		sanitizer: sanitizer{
			field:        p.field,
			errorMessage: p.errorMessage,
			location:     paramsLocation,
			rules:        make(validationChainRules, 0),
		},
	}
}

// the param struct creator function. which takes in the field to be validated and an errorMessage on failure.
func NewParam(field, errorMessage string) param {
	return param{
		field:        field,
		errorMessage: errorMessage,
	}
}
