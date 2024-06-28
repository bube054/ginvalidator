package ginvalidator

type param struct {
	field        string
	errorMessage string
}

const defaultParamLocation = "param"

func (p *param) Chain() processor {
	return processor{
		validator: validator{
			field:        p.field,
			errorMessage: p.errorMessage,
			location:     defaultParamLocation,
			rules:        make(validationProcessesRules, 0),
		},
		modifier: modifier{
			field:        p.field,
			errorMessage: p.errorMessage,
			location:     defaultParamLocation,
			rules:        make(validationProcessesRules, 0),
		},
		sanitizer: sanitizer{
			field:        p.field,
			errorMessage: p.errorMessage,
			location:     defaultParamLocation,
			rules:        make(validationProcessesRules, 0),
		},
	}
}

func NewParam(field, errorMessage string) param {
	return param{
		field:        field,
		errorMessage: errorMessage,
	}
}
