package ginvalidator

// Header is used to validate data from the `http.Request` headers.
type Header struct {
	field      string            // the field to be specified
	errFmtFunc ErrFmtFunc // the function to create the error message
}

// Chain initializes a validation chain for the given body field.
// It creates a new ValidationChain object that will validate the specified field
// and format error messages using the provided ErrFmtFunc.
func (h Header) Chain() ValidationChain {
	return newValidationChain(h.field, h.errFmtFunc, HeaderLocation)
}

// NewHeader constructs a Header validator for the given field.
// Returns a [Header] object that can be used to create validation chains.
//
// Parameters:
//   - field: the name of the field to validate.
//   - errFmtFunc: a handler for formatting error messages.
func NewHeader(field string, errFmtFunc ErrFmtFunc) Header {
	return Header{
		field:      field,
		errFmtFunc: errFmtFunc,
	}
}

// NewHeaderChain is a shorthand for NewHeader(field, errFmtFunc).Chain().
func NewHeaderChain(field string, errFmtFunc ErrFmtFunc) ValidationChain {
	return NewHeader(field, errFmtFunc).Chain()
}
