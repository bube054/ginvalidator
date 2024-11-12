package ginvalidator

import (
	"net/http"
	"testing"
)

func TestBail(t *testing.T) {
	tests := []struct {
		name string

		field      string
		errFmtFunc ErrFmtFuncHandler

		reqOpts ginCtxReqOpts

		want validationChainRule
	}{
		{
			name:       "Creates a Bail modifier validation chain rule.",
			field:      "name",
			errFmtFunc: nil,
			reqOpts:    ginCtxReqOpts{body: `{"name": "John"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("John"),
				withValidationChainName(BailModifierFuncName),
				withValidationChainType(modifierType),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := NewBody(test.field, test.errFmtFunc)
			chain := body.Chain()

			vc := chain.Bail()
			vcrs := vc.validator.rulesCreatorFuncs

			if len(vcrs) != 1 {
				t.Errorf("rule creators length invalid.")
				return
			}

			ctx := createTestGinCtx(test.reqOpts)
			vcr := vcrs[0]
			value, _ := extractFieldValFromBody(ctx, test.field)
			r := vcr(ctx, value, value)

			if r != test.want {
				t.Errorf("got %+v, want %+v", r, test.want)
			}
		})
	}
}

func TestIf(t *testing.T) {
	tests := []struct {
		name string

		imf        IfModifierFunc
		field      string
		errFmtFunc ErrFmtFuncHandler

		reqOpts ginCtxReqOpts

		want validationChainRule
	}{
		{
			name: "Creates an If modifier validation chain rule. It returns true, breaking the chain.",
			imf: func(req *http.Request, initialValue, sanitizedValue string) bool {
				return true
			},
			field:      "name",
			errFmtFunc: nil,
			reqOpts:    ginCtxReqOpts{body: `{"name": "John"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("John"),
				withValidationChainName(IfModifierFuncName),
				withValidationChainType(modifierType),
				withShouldBail(true),
			),
		},
		{
			name: "Creates an If modifier validation chain rule. It returns false, continuing the chain.",
			imf: func(req *http.Request, initialValue, sanitizedValue string) bool {
				return false
			},
			field:      "name",
			errFmtFunc: nil,
			reqOpts:    ginCtxReqOpts{body: `{"name": "John"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("John"),
				withValidationChainName(IfModifierFuncName),
				withValidationChainType(modifierType),
				withShouldBail(false),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := NewBody(test.field, test.errFmtFunc)
			chain := body.Chain()

			vc := chain.If(test.imf)
			vcrs := vc.validator.rulesCreatorFuncs

			if len(vcrs) != 1 {
				t.Errorf("rule creators length invalid.")
				return
			}

			ctx := createTestGinCtx(test.reqOpts)
			vcr := vcrs[0]
			value, _ := extractFieldValFromBody(ctx, test.field)
			r := vcr(ctx, value, value)

			if r != test.want {
				t.Errorf("got %+v, want %+v", r, test.want)
			}
		})
	}
}

func TestNot(t *testing.T) {
	tests := []struct {
		name string

		field      string
		errFmtFunc ErrFmtFuncHandler

		reqOpts ginCtxReqOpts

		want validationChainRule
	}{
		{
			name:       "Creates a Not modifier validation chain rule.",
			field:      "name",
			errFmtFunc: nil,
			reqOpts:    ginCtxReqOpts{body: `{"name": "John"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("John"),
				withValidationChainName(NotModifierFuncName),
				withValidationChainType(modifierType),
				withShouldBail(false),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := NewBody(test.field, test.errFmtFunc)
			chain := body.Chain()

			vc := chain.Not()
			vcrs := vc.validator.rulesCreatorFuncs

			if len(vcrs) != 1 {
				t.Errorf("rule creators length invalid.")
				return
			}

			ctx := createTestGinCtx(test.reqOpts)
			vcr := vcrs[0]
			value, _ := extractFieldValFromBody(ctx, test.field)
			r := vcr(ctx, value, value)

			if r != test.want {
				t.Errorf("got %+v, want %+v", r, test.want)
			}
		})
	}
}

func TestSkip(t *testing.T) {
	tests := []struct {
		name string

		smf        SkipModifierFunc
		field      string
		errFmtFunc ErrFmtFuncHandler

		reqOpts ginCtxReqOpts

		want validationChainRule
	}{
		{
			name: "Creates a Skip modifier validation chain rule. It returns true, skipping the next chain rule.",
			smf: func(req *http.Request, initialValue, sanitizedValue string) bool {
				return true
			},
			field:      "name",
			errFmtFunc: nil,
			reqOpts:    ginCtxReqOpts{body: `{"name": "John"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("John"),
				withValidationChainName(SkipModifierFuncName),
				withValidationChainType(modifierType),
				withShouldBail(false),
				withShouldSkip(true),
			),
		},
		{
			name: "Creates a Skip modifier validation chain rule. It returns false, continuing to the next chain rule.",
			smf: func(req *http.Request, initialValue, sanitizedValue string) bool {
				return false
			},
			field:      "name",
			errFmtFunc: nil,
			reqOpts:    ginCtxReqOpts{body: `{"name": "John"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("John"),
				withValidationChainName(SkipModifierFuncName),
				withValidationChainType(modifierType),
				withShouldBail(false),
				withShouldSkip(false),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := NewBody(test.field, test.errFmtFunc)
			chain := body.Chain()

			vc := chain.Skip(test.smf)
			vcrs := vc.validator.rulesCreatorFuncs

			if len(vcrs) != 1 {
				t.Errorf("rule creators length invalid.")
				return
			}

			ctx := createTestGinCtx(test.reqOpts)
			vcr := vcrs[0]
			value, _ := extractFieldValFromBody(ctx, test.field)
			r := vcr(ctx, value, value)

			if r != test.want {
				t.Errorf("got %+v, want %+v", r, test.want)
			}
		})
	}
}

func TestOptional(t *testing.T) {
	tests := []struct {
		name string

		smf        SkipModifierFunc
		field      string
		errFmtFunc ErrFmtFuncHandler

		reqOpts ginCtxReqOpts

		want validationChainRule
	}{
		{
			name: "Creates an Optional modifier validation chain rule.",

			field:      "name",
			errFmtFunc: nil,
			reqOpts:    ginCtxReqOpts{body: `{"name": "John"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("John"),
				withValidationChainName(OptionalModifierFuncName),
				withValidationChainType(modifierType),
				withShouldBail(false),
				withShouldSkip(false),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := NewBody(test.field, test.errFmtFunc)
			chain := body.Chain()

			vc := chain.Optional()
			vcrs := vc.validator.rulesCreatorFuncs

			if len(vcrs) != 1 {
				t.Errorf("rule creators length invalid.")
				return
			}

			ctx := createTestGinCtx(test.reqOpts)
			vcr := vcrs[0]
			value, _ := extractFieldValFromBody(ctx, test.field)
			r := vcr(ctx, value, value)

			if r != test.want {
				t.Errorf("got %+v, want %+v", r, test.want)
			}
		})
	}
}
