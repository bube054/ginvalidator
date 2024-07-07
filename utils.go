package ginvalidator

import (
	"fmt"
	"strings"

	valid "github.com/asaskevich/govalidator"
	"github.com/buger/jsonparser"
)

func wasPreviousRuleNegation(rules validationChainRules) bool {
	if len(rules) == 0 {
		return false
	}

	lastRuleIndex := len(rules) - 1
	lastRule := rules[lastRuleIndex]
	response := lastRule("", "", nil)

	if response.funcName == notFunc {
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

func splitJSONFieldSelector(value string) ([]string, error) {
	return strings.Split(value, "."), nil
}

func convertValueToJSON(value string) []byte {
	return []byte(fmt.Sprintf(`{"%s":"%s"}`, jsonKey, value))
}

func getJSONDataType(value []byte) string {
	_, typ, _, _ := jsonparser.Get(value, jsonKey)
	return typ.String()
}

func getFinalErrorMessage(errMsgFromChainMethod, errMsgFromConfig, defaultErrMsg string) (finalErrMessage string) {
	if errMsgFromChainMethod != "" {
		finalErrMessage = errMsgFromChainMethod
	} else if errMsgFromConfig != "" {
		finalErrMessage = errMsgFromConfig
	} else {
		finalErrMessage = defaultErrMsg
	}
	return
}
