package ginvalidator

type Sanitizer struct {
	field        string
	errorMessage string
	location     string

	value        string
	currentValue string
	resErrs      []ResponseError
}

func (s Sanitizer) BlackList() Processor {
	return Processor{
		Validator: Validator{},
		Modifier:  Modifier{},
		Sanitizer: Sanitizer{},
	}
}

func (s Sanitizer) LeftTrim() Processor {
	return Processor{
		Validator: Validator{},
		Modifier:  Modifier{},
		Sanitizer: Sanitizer{},
	}
}

func (s Sanitizer) Trim() Processor {
	return Processor{
		Validator: Validator{},
		Modifier:  Modifier{},
		Sanitizer: Sanitizer{},
	}
}
