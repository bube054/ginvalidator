package ginvalidator

type Modifier struct {
	field        string
	errorMessage string
	location     string

	value        string
	currentValue string
	resErrs      []ResponseError
}

func (m Modifier) ReplacePattern() Processor {
	return Processor{
		Validator: Validator{},
		Modifier:  Modifier{},
		Sanitizer: Sanitizer{},
	}
}

func (m Modifier) Reverse() Processor {
	return Processor{
		Validator: Validator{},
		Modifier:  Modifier{},
		Sanitizer: Sanitizer{},
	}
}

func (m Modifier) RightTrim() Processor {
	return Processor{
		Validator: Validator{},
		Modifier:  Modifier{},
		Sanitizer: Sanitizer{},
	}
}
