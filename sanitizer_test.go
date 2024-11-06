package ginvalidator

import (
	"net/http"
	"testing"

	san "github.com/bube054/validatorgo/sanitizer"
)

func TestCustomSanitizer(t *testing.T) {
	tests := []struct {
		name string

		field      string
		errFmtFunc ErrFmtFuncHandler

		csf     CustomSanitizerFunc
		reqOpts ginCtxReqOpts

		want validationChainRule
	}{
		{
			name:       "Creates a CustomSanitizer chain rule. Returns the validated value.",
			field:      "name",
			errFmtFunc: nil,
			csf: func(req http.Request, initialValue, sanitizedValue string) string {
				return initialValue
			},
			reqOpts: ginCtxReqOpts{body: `{"name": "John"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("John"),
				withValidationChainName(CustomSanitizerFuncName),
				withValidationChainType(sanitizerType),
			),
		},
		{
			name:       "Creates a CustomSanitizer chain rule. Returns the an empty string.",
			field:      "name",
			errFmtFunc: nil,
			csf: func(req http.Request, initialValue, sanitizedValue string) string {
				return ""
			},
			reqOpts: ginCtxReqOpts{body: `{"name": "John"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue(""),
				withValidationChainName(CustomSanitizerFuncName),
				withValidationChainType(sanitizerType),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := NewBody(test.field, test.errFmtFunc)
			chain := body.CreateChain()

			vc := chain.CustomSanitizer(test.csf)
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

func TestBlacklist(t *testing.T) {
	tests := []struct {
		name string

		field      string
		errFmtFunc ErrFmtFuncHandler

		blacklistedChars string
		reqOpts          ginCtxReqOpts

		want validationChainRule
	}{
		{
			name:             "Creates a Blacklist sanitizer chain rule.",
			field:            "name",
			errFmtFunc:       nil,
			blacklistedChars: "0-9",
			reqOpts:          ginCtxReqOpts{body: `{"name": "John109"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("John"),
				withValidationChainName(BlacklistSanitizerFuncName),
				withValidationChainType(sanitizerType),
			),
		},
		{
			name:             "Creates a Blacklist sanitizer chain rule.",
			field:            "name",
			errFmtFunc:       nil,
			blacklistedChars: "[a-zA-Z]",
			reqOpts:          ginCtxReqOpts{body: `{"name": "John109"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("109"),
				withValidationChainName(BlacklistSanitizerFuncName),
				withValidationChainType(sanitizerType),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := NewBody(test.field, test.errFmtFunc)
			chain := body.CreateChain()

			vc := chain.Blacklist(test.blacklistedChars)
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

func TestEscape(t *testing.T) {
	tests := []struct {
		name string

		field      string
		errFmtFunc ErrFmtFuncHandler

		reqOpts ginCtxReqOpts

		want validationChainRule
	}{
		{
			name:       "Creates an Escape sanitizer chain rule.",
			field:      "name",
			errFmtFunc: nil,
			reqOpts:    ginCtxReqOpts{body: `{"name": "<John"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("&lt;John"),
				withValidationChainName(EscapeSanitizerFuncName),
				withValidationChainType(sanitizerType),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := NewBody(test.field, test.errFmtFunc)
			chain := body.CreateChain()

			vc := chain.Escape()
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

func TestLTrim(t *testing.T) {
	tests := []struct {
		name string

		field      string
		errFmtFunc ErrFmtFuncHandler

		chars   string
		reqOpts ginCtxReqOpts

		want validationChainRule
	}{
		{
			name:       "Creates an LTrim sanitizer chain rule.",
			field:      "name",
			errFmtFunc: nil,
			chars:      "",
			reqOpts:    ginCtxReqOpts{body: `{"name": "John"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("John"),
				withValidationChainName(LTrimSanitizerFuncName),
				withValidationChainType(sanitizerType),
			),
		},
		{
			name:       "Sanitized value returned is left trimmed",
			field:      "name",
			errFmtFunc: nil,
			chars:      " ",
			reqOpts:    ginCtxReqOpts{body: `{"name": " John "}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("John "),
				withValidationChainName(LTrimSanitizerFuncName),
				withValidationChainType(sanitizerType),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := NewBody(test.field, test.errFmtFunc)
			chain := body.CreateChain()

			vc := chain.LTrim(test.chars)
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

func TestNormalizeEmail(t *testing.T) {
	tests := []struct {
		name string

		field      string
		errFmtFunc ErrFmtFuncHandler

		opts    *san.NormalizeEmailOpts
		reqOpts ginCtxReqOpts

		want validationChainRule
	}{
		{
			name:       "Creates a NormalizeEmail sanitizer chain rule.",
			field:      "email",
			errFmtFunc: nil,
			opts:       nil,
			reqOpts:    ginCtxReqOpts{body: `{"email": "Example@Example.com"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("example@example.com"),
				withValidationChainName(NormalizeEmailSanitizerFuncName),
				withValidationChainType(sanitizerType),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := NewBody(test.field, test.errFmtFunc)
			chain := body.CreateChain()

			vc := chain.NormalizeEmail(test.opts)
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

func TestRTrim(t *testing.T) {
	tests := []struct {
		name string

		field      string
		errFmtFunc ErrFmtFuncHandler

		chars   string
		reqOpts ginCtxReqOpts

		want validationChainRule
	}{
		{
			name:       "Creates an RTrim sanitizer chain rule.",
			field:      "name",
			errFmtFunc: nil,
			chars:      "",
			reqOpts:    ginCtxReqOpts{body: `{"name": "John"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("John"),
				withValidationChainName(RTrimSanitizerFuncName),
				withValidationChainType(sanitizerType),
			),
		},
		{
			name:       "Creates an RTrim sanitizer chain rule.",
			field:      "name",
			errFmtFunc: nil,
			chars:      " ",
			reqOpts:    ginCtxReqOpts{body: `{"name": " John "}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue(" John"),
				withValidationChainName(RTrimSanitizerFuncName),
				withValidationChainType(sanitizerType),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := NewBody(test.field, test.errFmtFunc)
			chain := body.CreateChain()

			vc := chain.RTrim(test.chars)
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

func TestStripLow(t *testing.T) {
	tests := []struct {
		name string

		field      string
		errFmtFunc ErrFmtFuncHandler

		keepNewLines bool
		reqOpts      ginCtxReqOpts

		want validationChainRule
	}{
		{
			name:         "Creates a StripLow sanitizer chain rule.",
			field:        "name",
			errFmtFunc:   nil,
			keepNewLines: false,
			reqOpts:      ginCtxReqOpts{body: `{"name": "Hello\nWorld"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("HelloWorld"),
				withValidationChainName(StripLowSanitizerFuncName),
				withValidationChainType(sanitizerType),
			),
		},
		{
			name:         "Creates a StripLow sanitizer chain rule.",
			field:        "name",
			errFmtFunc:   nil,
			keepNewLines: true,
			reqOpts:      ginCtxReqOpts{body: `{"name": " John\n"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue(" John\n"),
				withValidationChainName(StripLowSanitizerFuncName),
				withValidationChainType(sanitizerType),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := NewBody(test.field, test.errFmtFunc)
			chain := body.CreateChain()

			vc := chain.StripLow(test.keepNewLines)
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

func TestToBoolean(t *testing.T) {
	tests := []struct {
		name string

		field      string
		errFmtFunc ErrFmtFuncHandler

		strict bool
		reqOpts      ginCtxReqOpts

		want validationChainRule
	}{
		{
			name:         "Creates a Toboolean sanitizer chain rule.",
			field:        "name",
			errFmtFunc:   nil,
			strict: false,
			reqOpts:      ginCtxReqOpts{body: `{"name": "true"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("true"),
				withValidationChainName(ToBooleanSanitizerFuncName),
				withValidationChainType(sanitizerType),
			),
		},
		{
			name:         "Creates a Toboolean sanitizer chain rule.",
			field:        "name",
			errFmtFunc:   nil,
			strict: true,
			reqOpts:      ginCtxReqOpts{body: `{"name": "false"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("false"),
				withValidationChainName(ToBooleanSanitizerFuncName),
				withValidationChainType(sanitizerType),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := NewBody(test.field, test.errFmtFunc)
			chain := body.CreateChain()

			vc := chain.ToBoolean(test.strict)
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

func TestToDate(t *testing.T) {
	tests := []struct {
		name string

		field      string
		errFmtFunc ErrFmtFuncHandler

		reqOpts ginCtxReqOpts

		want validationChainRule
	}{
		{
			name:       "Creates a ToDate sanitizer chain rule.",
			field:      "date",
			errFmtFunc: nil,
			reqOpts:    ginCtxReqOpts{body: `{"date": "Mon Jan  2 15:04:05 2006"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("2006-01-02 15:04:05"),
				withValidationChainName(ToDateSanitizerFuncName),
				withValidationChainType(sanitizerType),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := NewBody(test.field, test.errFmtFunc)
			chain := body.CreateChain()

			vc := chain.ToDate()
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

func TestToFloat(t *testing.T) {
	tests := []struct {
		name string

		field      string
		errFmtFunc ErrFmtFuncHandler

		reqOpts ginCtxReqOpts

		want validationChainRule
	}{
		{
			name:       "Creates a ToFloat sanitizer chain rule.",
			field:      "flt",
			errFmtFunc: nil,
			reqOpts:    ginCtxReqOpts{body: `{"flt": "123"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("123.000000"),
				withValidationChainName(ToFloatSanitizerFuncName),
				withValidationChainType(sanitizerType),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := NewBody(test.field, test.errFmtFunc)
			chain := body.CreateChain()

			vc := chain.ToFloat()
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

func TestToInt(t *testing.T) {
	tests := []struct {
		name string

		field      string
		errFmtFunc ErrFmtFuncHandler

		reqOpts ginCtxReqOpts

		want validationChainRule
	}{
		{
			name:       "Creates a ToInt sanitizer chain rule.",
			field:      "int",
			errFmtFunc: nil,
			reqOpts:    ginCtxReqOpts{body: `{"int": "123"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("123"),
				withValidationChainName(ToIntSanitizerFuncName),
				withValidationChainType(sanitizerType),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := NewBody(test.field, test.errFmtFunc)
			chain := body.CreateChain()

			vc := chain.ToInt()
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

func TestTrim(t *testing.T) {
	tests := []struct {
		name string

		field      string
		errFmtFunc ErrFmtFuncHandler

		chars   string
		reqOpts ginCtxReqOpts

		want validationChainRule
	}{
		{
			name:       "Creates a Trim sanitizer chain rule.",
			field:      "name",
			errFmtFunc: nil,
			chars:      "",
			reqOpts:    ginCtxReqOpts{body: `{"name": "John"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("John"),
				withValidationChainName(TrimSanitizerFuncName),
				withValidationChainType(sanitizerType),
			),
		},
		{
			name:       "Creates a Trim sanitizer chain rule.",
			field:      "name",
			errFmtFunc: nil,
			chars:      " ",
			reqOpts:    ginCtxReqOpts{body: `{"name": " John "}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("John"),
				withValidationChainName(TrimSanitizerFuncName),
				withValidationChainType(sanitizerType),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := NewBody(test.field, test.errFmtFunc)
			chain := body.CreateChain()

			vc := chain.Trim(test.chars)
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

func TestUnescape(t *testing.T) {
	tests := []struct {
		name string

		field      string
		errFmtFunc ErrFmtFuncHandler

		reqOpts ginCtxReqOpts

		want validationChainRule
	}{
		{
			name:       "Creates an Unescape sanitizer chain rule.",
			field:      "name",
			errFmtFunc: nil,
			reqOpts:    ginCtxReqOpts{body: `{"name": "&lt;John"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("<John"),
				withValidationChainName(UnescapeSanitizerFuncName),
				withValidationChainType(sanitizerType),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := NewBody(test.field, test.errFmtFunc)
			chain := body.CreateChain()

			vc := chain.Unescape()
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

func TestWhitelist(t *testing.T) {
	tests := []struct {
		name string

		field      string
		errFmtFunc ErrFmtFuncHandler

		whitelistedChars string
		reqOpts          ginCtxReqOpts

		want validationChainRule
	}{
		{
			name:             "Creates a Whitelist sanitizer chain rule.",
			field:            "name",
			errFmtFunc:       nil,
			whitelistedChars: "0-9",
			reqOpts:          ginCtxReqOpts{body: `{"name": "John109"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("109"),
				withValidationChainName(WhitelistSanitizerFuncName),
				withValidationChainType(sanitizerType),
			),
		},
		{
			name:             "Creates a Whitelist sanitizer chain rule.",
			field:            "name",
			errFmtFunc:       nil,
			whitelistedChars: "[a-zA-Z]",
			reqOpts:          ginCtxReqOpts{body: `{"name": "John109"}`, contentType: "application/json"},
			want: NewValidationChainRule(
				withIsValid(true),
				withNewValue("John"),
				withValidationChainName(WhitelistSanitizerFuncName),
				withValidationChainType(sanitizerType),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := NewBody(test.field, test.errFmtFunc)
			chain := body.CreateChain()

			vc := chain.Whitelist(test.whitelistedChars)
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