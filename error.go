package ginvalidator

type ErrFmtFuncHandler func(initialValue, sanitizedValue, validatorName string) string
