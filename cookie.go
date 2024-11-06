package ginvalidator

// Cookie is used to validate data from the `http.Request` cookies.
type Cookie struct {
	field      string             // the field to be specified
	errFmtFunc ErrFmtFuncHandler // the function to create the error message
}

// CreateChain initializes a validation chain for the given body field.
// It creates a new ValidationChain object that will validate the specified field
// and format error messages using the provided ErrFmtFuncHandler.
func (c Cookie) CreateChain() ValidationChain {
	return NewValidationChain(c.field, c.errFmtFunc, cookieLocation)
}

// NewCookie constructs a Cookie validator for the given field.
// Returns a [Cookie] object that can be used to create validation chains.
//
// Parameters:
//   - field: the name of the field to validate.
//   - errFmtFunc: a handler for formatting error messages.
func NewCookie(field string, errFmtFunc ErrFmtFuncHandler) Cookie {
	return Cookie{
		field:      field,
		errFmtFunc: errFmtFunc,
	}
}
