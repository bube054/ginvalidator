package ginvalidator

import (
	valid "github.com/asaskevich/govalidator"
)

func wasPreviousRuleNegation(rules validationProcessesRules) bool {
	if len(rules) == 0 {
		return false
	}

	lastRuleIndex := len(rules) - 1
	lastRule := rules[lastRuleIndex]
	response := lastRule("", "", nil)

	if response.funcName == "IsNot" {
		return true
	} else {
		return false
	}
}

func valueIsNullish(value string) bool {
	if value == "" || valid.IsNull(value) || value == "undefined" || value == "NaN" {
		return true
	} else {
		return false
	}
}

func valueIsInSlice(value string, valuesFrom []string) bool {
	for _, val := range valuesFrom {
		if value == val {
			return true
		}
	}

	return false
}
