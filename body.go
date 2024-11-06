package ginvalidator

// Body is used to validate data from the `http.Request` body.
type Body struct {
	field      string             // the field to be specified
	errFmtFunc ErrFmtFuncHandler // the function to create the error message
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
//   - field: the name of the field to validate. It uses [gjson] for its json field extraction syntax.
//   - errFmtFunc: a handler for formatting error messages.
//
// [gjson]: https://github.com/tidwall/gjson?tab=readme-ov-file#path-syntax
func NewBody(field string, errFmtFunc ErrFmtFuncHandler) Body {
	return Body{
		field:      field,
		errFmtFunc: errFmtFunc,
	}
}
