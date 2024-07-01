package ginvalidator

// request validator struct for headers e.g e.g "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
type header struct {
	field        string
	errorMessage string
}

// header creator of the validation chain.
func (p *header) Chain() validationChain {
	return validationChain{
		validator: validator{
			field:           p.field,
			errorMessage:    p.errorMessage,
			location:        headersLocation,
			rules:           make(validationChainRules, 0),
			chainMethodType: validatorType,
		},
		modifier: modifier{
			field:           p.field,
			errorMessage:    p.errorMessage,
			location:        headersLocation,
			rules:           make(validationChainRules, 0),
			chainMethodType: modifierType,
		},
		sanitizer: sanitizer{
			field:           p.field,
			errorMessage:    p.errorMessage,
			location:        headersLocation,
			rules:           make(validationChainRules, 0),
			chainMethodType: sanitizerType,
		},
	}
}

// the header struct creator function. which takes in the field to be validated and an errorMessage on failure.
func NewHeader(field, errorMessage string) header {
	return header{
		field:        field,
		errorMessage: errorMessage,
	}
}
