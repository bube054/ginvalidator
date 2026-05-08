package ginvalidator

// ValidationChainError represents an error that occurred during the validation chain for a request.
// It includes information about the location of the error, the message, the specific field involved,
// the invalid value.
//
// Fields:
//   - Location: The location in the request where the error occurred (e.g., "body", "cookies", "headers", "params", "queries").
//   - Msg: A message describing the validation error.
//   - Field: The name of the field that failed validation.
//   - Value: The invalid value that triggered the validation error.
//   - order: A monotonic counter used internally to preserve insertion order across chains.
type ValidationChainError struct {
	Location string `json:"location"`
	Msg      string `json:"message"`
	Field    string `json:"field"`
	Value    string `json:"value"`
	Code     string `json:"code,omitempty"`
	order    uint64
}

func vceWithLocation(location string) func(*ValidationChainError) {
	return func(vce *ValidationChainError) {
		vce.Location = location
	}
}

func vceWithMsg(msg string) func(*ValidationChainError) {
	return func(vce *ValidationChainError) {
		vce.Msg = msg
	}
}

func vceWithField(field string) func(*ValidationChainError) {
	return func(vce *ValidationChainError) {
		vce.Field = field
	}
}

func vceWithValue(value string) func(*ValidationChainError) {
	return func(vce *ValidationChainError) {
		vce.Value = value
	}
}

func vceWithCode(code string) func(*ValidationChainError) {
	return func(vce *ValidationChainError) {
		vce.Code = code
	}
}

func vceWithOrder(order uint64) func(*ValidationChainError) {
	return func(vce *ValidationChainError) {
		vce.order = order
	}
}

func NewValidationChainError(opts ...func(*ValidationChainError)) ValidationChainError {
	vce := &ValidationChainError{}

	for _, opt := range opts {
		opt(vce)
	}

	return *vce
}

// ErrFmtFuncHandler is a function type used to format validation error messages.
// It takes in the initial and sanitized values of a field, along with the name of the validator
// that triggered the error, and returns a formatted error message as a string.
//
// Parameters:
//   - initialValue: The original value of the field before sanitization.
//   - sanitizedValue: The value of the field after applying sanitization or validation.
//   - validatorName: The name of the validator that was applied and caused the error.
//
// Returns:
//   - A string representing the formatted error message based on the provided values and validator.
type ErrFmtFuncHandler func(initialValue, sanitizedValue, validatorName string) string
