package ginvalidator

import (
	"net/http"
	"testing"

	vgo "github.com/bube054/validatorgo"
)

func TestCustomValidator(t *testing.T) {
	tests := []struct {
		name string

		field      string
		errFmtFunc ErrFmtFuncHandler

		cvf     CustomValidatorFunc
		reqOpts ginCtxReqOpts

		want validationChainRule
	}{
		{
			name:       "Creates a CustomValidator chain rule. Returns true.",
			field:      "name",
			errFmtFunc: nil,
			cvf: func(req http.Request, initialValue, sanitizedValue string) bool {
				return true
			},
			reqOpts: ginCtxReqOpts{body: `{"name": "John"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("John"),
				withValidationChainName(CustomValidatorName),
				withValidationChainType(validatorType),
			),
		},
		{
			name:       "Creates a CustomValidator chain rule. Returns false.",
			field:      "name",
			errFmtFunc: nil,
			cvf: func(req http.Request, initialValue, sanitizedValue string) bool {
				return false
			},
			reqOpts: ginCtxReqOpts{body: `{"name": "John"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(false),
				withNewValue("John"),
				withValidationChainName(CustomValidatorName),
				withValidationChainType(validatorType),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := NewBody(test.field, test.errFmtFunc)
			chain := body.CreateChain()

			vc := chain.CustomValidator(test.cvf)
			vcrs := vc.validator.rulesCreatorFuncs

			if len(vcrs) != 1 {
				t.Errorf("rule creators length invalid.")
				return
			}

			ctx := createTestGinCtx(test.reqOpts)
			vcr := vcrs[0]
			value, _ := extractFieldValFromBody(test.field, ctx)
			r := vcr(ctx, value, value)

			if r != test.want {
				t.Errorf("got %+v, want %+v", r, test.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name string

		field      string
		errFmtFunc ErrFmtFuncHandler

		seed string
		opts *vgo.ContainsOpt

		reqOpts ginCtxReqOpts

		want validationChainRule
	}{
		{
			name:       "Creates a Contains validator chain rule. Returns true.",
			field:      "text",
			errFmtFunc: nil,
			seed: "world",
			opts:  &vgo.ContainsOpt{},
			reqOpts: ginCtxReqOpts{body: `{"text": "Hello world"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("Hello world"),
				withValidationChainName(ContainsValidatorName),
				withValidationChainType(validatorType),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := NewBody(test.field, test.errFmtFunc)
			chain := body.CreateChain()

			vc := chain.Contains(test.seed, test.opts)
			vcrs := vc.validator.rulesCreatorFuncs

			if len(vcrs) != 1 {
				t.Errorf("rule creators length invalid.")
				return
			}

			ctx := createTestGinCtx(test.reqOpts)
			vcr := vcrs[0]
			value, _ := extractFieldValFromBody(test.field, ctx)
			r := vcr(ctx, value, value)

			if r != test.want {
				t.Errorf("got %+v, want %+v", r, test.want)
			}
		})
	}
}