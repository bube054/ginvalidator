package ginvalidator

type ResponseError struct {
	location string
	msg      string
	path     string
	typ      string
}

type Processor struct {
	Validator
	Modifier
	Sanitizer
}

func (p *Processor) GetErrors() []ResponseError {
	return p.Validator.resErrs
}
