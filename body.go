package ginvalidator

// Body is used to validate data from the `http.Request` body.
type Body struct {
	field      string             // the field to be specified
	errFmtFunc *ErrFmtFuncHandler // the function to create the error message
}

// CreateChain initializes a validation chain for the given body field.
// It creates a new ValidationChain object that will validate the specified field
// and format error messages using the provided ErrFmtFuncHandler.
func (b Body) CreateChain() ValidationChain {
	return NewValidationChain(b.field, b.errFmtFunc, bodyLocation)
}

// NewBody constructs a Body validator for the given field.
// Returns a [Body] object that can be used to create validation chains.
//
// Parameters:
//   - field: the name of the field to validate.
//   - errFmtFunc: a handler for formatting error messages.
func NewBody(field string, errFmtFunc *ErrFmtFuncHandler) Body {
	return Body{
		field:      field,
		errFmtFunc: errFmtFunc,
	}
}