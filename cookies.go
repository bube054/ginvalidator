package ginvalidator

// request validator struct for request cookies e.g "AFC": AGTF6HpKzga79RYx0BeAvTSzZZ2_f24n5nO7zQSgn2KKZS55E7Jm1VqX4xs
type cookie struct {
	field        string
	errorMessage string
}

// cookie creator of the validation chain.
func (c *cookie) Chain() validationChain {
	return validationChain{
		validator: validator{
			field:           c.field,
			errorMessage:    c.errorMessage,
			location:        cookiesLocation,
			rules:           make(validationChainRules, 0),
			chainMethodType: validatorType,
		},
		modifier: modifier{
			field:           c.field,
			errorMessage:    c.errorMessage,
			location:        cookiesLocation,
			rules:           make(validationChainRules, 0),
			chainMethodType: modifierType,
		},
		sanitizer: sanitizer{
			field:           c.field,
			errorMessage:    c.errorMessage,
			location:        cookiesLocation,
			rules:           make(validationChainRules, 0),
			chainMethodType: sanitizerType,
		},
	}
}

// the cookie struct creator function. which takes in the field to be validated and an errorMessage on failure.
func NewCookie(field, errorMessage string) cookie {
	return cookie{
		field:        field,
		errorMessage: errorMessage,
	}
}
