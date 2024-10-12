package ginvalidator

type ValidationChainError struct {
	Location       string
	Msg            string
	Field          string
	Value          string
}

func VCEWithLocation(location string) func(*ValidationChainError) {
	return func(vce *ValidationChainError) {
		vce.Location = location
	}
}

func VCEWithMsg(msg string) func(*ValidationChainError) {
	return func(vce *ValidationChainError) {
		vce.Msg = msg
	}
}

func VCEWithField(field string) func(*ValidationChainError) {
	return func(vce *ValidationChainError) {
		vce.Field = field
	}
}

func VCEWithValue(value string) func(*ValidationChainError) {
	return func(vce *ValidationChainError) {
		vce.Value = value
	}
}

// func VCEWithSanitizedValue(sanitizedValue string) func(*ValidationChainError) {
// 	return func(vce *ValidationChainError) {
// 		vce.SanitizedValue = sanitizedValue
// 	}
// }

func NewValidationChainError(opts ...func(*ValidationChainError)) ValidationChainError {
	vce := &ValidationChainError{}

	for _, opt := range opts {
		opt(vce)
	}

	return *vce
}

type ErrFmtFuncHandler func(initialValue, sanitizedValue, validatorName string) string
