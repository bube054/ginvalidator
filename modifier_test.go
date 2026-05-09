package ginvalidator

import (
	"net/http"
	"testing"
)

func TestBail(t *testing.T) {
	tests := []struct {
		name string

		field      string
		errFmtFunc ErrFmtFunc

		reqOpts ginCtxReqOpts

		want validationChainRule
	}{
		{
			name:       "Creates a Bail modifier validation chain rule.",
			field:      "name",
			errFmtFunc: nil,
			reqOpts:    ginCtxReqOpts{body: `{"name": "John"}`, contentType: "application/json"},
			want: newValidationChainRule(
				withIsValid(true),
				withNewValue("John"),
				withValidationChainName(BailModifierName),
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
		errFmtFunc ErrFmtFunc

		reqOpts ginCtxReqOpts

		want validationChainRule
	}{
		{
			name: "Creates an If modifier validation chain rule. It returns true, breaking the chain.",
			imf: func(r *http.Request, initialValue, sanitizedValue string) bool {
				return true
			},
			field:      "name",
			errFmtFunc: nil,
			reqOpts:    ginCtxReqOpts{body: `{"name": "John"}`, contentType: "application/json"},
			want: newValidationChainRule(
				withIsValid(true),
				withNewValue("John"),
				withValidationChainName(IfModifierName),
				withValidationChainType(modifierType),
				withShouldBail(true),
			),
		},
		{
			name: "Creates an If modifier validation chain rule. It returns false, continuing the chain.",
			imf: func(r *http.Request, initialValue, sanitizedValue string) bool {
				return false
			},
			field:      "name",
			errFmtFunc: nil,
			reqOpts:    ginCtxReqOpts{body: `{"name": "John"}`, contentType: "application/json"},
			want: newValidationChainRule(
				withIsValid(true),
				withNewValue("John"),
				withValidationChainName(IfModifierName),
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
		errFmtFunc ErrFmtFunc

		reqOpts ginCtxReqOpts

		want validationChainRule
	}{
		{
			name:       "Creates a Not modifier validation chain rule.",
			field:      "name",
			errFmtFunc: nil,
			reqOpts:    ginCtxReqOpts{body: `{"name": "John"}`, contentType: "application/json"},
			want: newValidationChainRule(
				withIsValid(true),
				withNewValue("John"),
				withValidationChainName(NotModifierName),
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
		errFmtFunc ErrFmtFunc

		reqOpts ginCtxReqOpts

		want validationChainRule
	}{
		{
			name: "Creates a Skip modifier validation chain rule. It returns true, skipping the next chain rule.",
			smf: func(r *http.Request, initialValue, sanitizedValue string) bool {
				return true
			},
			field:      "name",
			errFmtFunc: nil,
			reqOpts:    ginCtxReqOpts{body: `{"name": "John"}`, contentType: "application/json"},
			want: newValidationChainRule(
				withIsValid(true),
				withNewValue("John"),
				withValidationChainName(SkipModifierName),
				withValidationChainType(modifierType),
				withShouldBail(false),
				withShouldSkip(true),
			),
		},
		{
			name: "Creates a Skip modifier validation chain rule. It returns false, continuing to the next chain rule.",
			smf: func(r *http.Request, initialValue, sanitizedValue string) bool {
				return false
			},
			field:      "name",
			errFmtFunc: nil,
			reqOpts:    ginCtxReqOpts{body: `{"name": "John"}`, contentType: "application/json"},
			want: newValidationChainRule(
				withIsValid(true),
				withNewValue("John"),
				withValidationChainName(SkipModifierName),
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
		errFmtFunc ErrFmtFunc

		reqOpts ginCtxReqOpts

		want validationChainRule
	}{
		{
			name: "Creates an Optional modifier validation chain rule.",

			field:      "name",
			errFmtFunc: nil,
			reqOpts:    ginCtxReqOpts{body: `{"name": "John"}`, contentType: "application/json"},
			want: newValidationChainRule(
				withIsValid(true),
				withNewValue("John"),
				withValidationChainName(OptionalModifierName),
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
