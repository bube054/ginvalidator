package ginvalidator

import "time"

type ValidationChainError struct {
	Location  string
	Msg       string
	Field     string
	Value     string
	createdAt time.Time
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

func vceWithCreatedAt(time time.Time) func(*ValidationChainError) {
	return func(vce *ValidationChainError) {
		vce.createdAt = time
	}
}

// func vceWithSanitizedValue(sanitizedValue string) func(*ValidationChainError) {
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
