package ginvalidator

// import (
//   valid "github.com/asaskevich/govalidator"
// )

type sanitizer struct {
	field        string
	errorMessage string
	location     string
	rules        validationProcessesRules
}
