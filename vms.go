package ginvalidator

import (
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
)

type VMS struct {
	validator
	modifier
	sanitizer
}

func (v VMS) Validate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bodyBytes, err := ctx.GetRawData()

		if err != nil {
			panic(err)
		}

		field := v.validator.field
		valueByte, _, _, _ := jsonparser.Get(bodyBytes, field)
		initialValue := string(valueByte)
		// sanitizedValue := initialValue
		ruleCreators := v.validator.rulesCreatorFuncs
		rules := make(vmsRules, len(ruleCreators))

		for ind, ruleCreator := range ruleCreators {
			rule := ruleCreator(ctx, initialValue)
			rules[ind] = rule
		}

		fmt.Println(rules)

		ctx.Next()
	}
}

func NewVMS(field string, errFmtFunc *ErrFmtFuncHandler, reqLoc requestLocation) VMS {
	return VMS{
		validator: newValidator(field, errFmtFunc, reqLoc),
		modifier:  newModifier(field, errFmtFunc, reqLoc),
		sanitizer: newSanitizer(field, errFmtFunc, reqLoc),
	}
}
