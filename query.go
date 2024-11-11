package ginvalidator

// Query is used to validate data from the `http.Request` queries.
type Query struct {
	field      string            // the field to be specified
	errFmtFunc ErrFmtFuncHandler // the function to create the error message
}

// Chain initializes a validation chain for the given body field.
// It creates a new ValidationChain object that will validate the specified field
// and format error messages using the provided ErrFmtFuncHandler.
func (q Query) Chain() ValidationChain {
	return NewValidationChain(q.field, q.errFmtFunc, QueryLocation)
}

// NewQuery constructs a Query validator for the given field.
// Returns a [Query] object that can be used to create validation chains.
//
// Parameters:
//   - field: the name of the field to validate.
//   - errFmtFunc: a handler for formatting error messages.
func NewQuery(field string, errFmtFunc ErrFmtFuncHandler) Query {
	return Query{
		field:      field,
		errFmtFunc: errFmtFunc,
	}
}
