package ginvalidator

// Header is used to validate data from the `http.Request` headers.
type Header struct {
	field      string             // the field to be specified
	errFmtFunc ErrFmtFuncHandler // the function to create the error message
}

// CreateChain initializes a validation chain for the given body field.
// It creates a new ValidationChain object that will validate the specified field
// and format error messages using the provided ErrFmtFuncHandler.
func (h Header) CreateChain() ValidationChain {
	return NewValidationChain(h.field, h.errFmtFunc, headerLocation)
}

// NewHeader constructs a Header validator for the given field.
// Returns a [Header] object that can be used to create validation chains.
//
// Parameters:
//   - field: the name of the field to validate.
//   - errFmtFunc: a handler for formatting error messages.
func NewHeader(field string, errFmtFunc ErrFmtFuncHandler) Header {
	return Header{
		field:      field,
		errFmtFunc: errFmtFunc,
	}
}
