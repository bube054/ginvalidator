package ginvalidator

// import (
//   valid "github.com/asaskevich/govalidator"
// )

type modifier struct {
	field        string
	errorMessage string
	location     string
	rules        validationProcessesRules
}
