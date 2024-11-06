package ginvalidator

// Param is used to validate data from the `http.Request` params.
type Param struct {
	field      string            // the field to be specified
	errFmtFunc ErrFmtFuncHandler // the function to create the error message
}

// Chain initializes a validation chain for the given body field.
// It creates a new ValidationChain object that will validate the specified field
// and format error messages using the provided ErrFmtFuncHandler.
func (p Param) Chain() ValidationChain {
	return NewValidationChain(p.field, p.errFmtFunc, paramLocation)
}

// NewParam constructs a Param validator for the given field.
// Returns a [Param] object that can be used to create validation chains.
//
// Parameters:
//   - field: the name of the field to validate.
//   - errFmtFunc: a handler for formatting error messages.
func NewParam(field string, errFmtFunc ErrFmtFuncHandler) Param {
	return Param{
		field:      field,
		errFmtFunc: errFmtFunc,
	}
}
