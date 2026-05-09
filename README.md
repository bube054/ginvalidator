# ginvalidator

[![CI](https://github.com/bube054/ginvalidator/actions/workflows/ci.yml/badge.svg)](https://github.com/bube054/ginvalidator/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/bube054/ginvalidator.svg)](https://pkg.go.dev/github.com/bube054/ginvalidator)
[![Go Report Card](https://goreportcard.com/badge/github.com/bube054/ginvalidator)](https://goreportcard.com/report/github.com/bube054/ginvalidator)
![License](https://img.shields.io/github/license/bube054/ginvalidator)

Middleware-based request validation for [Gin](https://github.com/gin-gonic/gin). Chain validators, sanitizers and modifiers on any request field, then collect the errors however you want.

Inspired by [express-validator](https://github.com/express-validator/express-validator). All built-in validators and sanitizers are powered by [validatorgo](https://github.com/bube054/validatorgo).

## Requirements

- [Go](https://go.dev/dl/) 1.22+
- [Gin](https://github.com/gin-gonic/gin) 1.x

## Install

```bash
go get github.com/bube054/ginvalidator
```

## Getting Started

Let's build a small signup API together, step by step. By the end you'll know how to validate fields, read errors, sanitize input and access the cleaned-up data.

### Step 1 — A bare Gin server

Create a new folder, run `go mod init example.com/signup`, then create `main.go`:

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/signup", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "welcome aboard"})
	})

	r.Run()
}
```

Run it:

```bash
go run main.go
```

The server is now listening on `http://localhost:8080`. Every time you change your code, you'll need to stop the server (Ctrl+C) and run `go run main.go` again. If that gets annoying, check out [Air](https://www.bytesizego.com/blog/golang-air) — it watches your files and restarts the server automatically on every save.

Now open a **separate terminal** and send a request:

```bash
curl -X POST http://localhost:8080/signup
```

> **Windows users:** PowerShell's `curl` is actually an alias for `Invoke-WebRequest` — avoid it. Use **Git Bash**, **WSL**, or **Command Prompt** instead. Some curl commands in this guide use bash syntax (`\` for line breaks, single quotes around JSON) that won't work in Command Prompt. Whenever that's the case, we'll show both the bash version and a Command Prompt version right below it. If you're new to curl, [this guide](https://everything.curl.dev/cmdline) is a good starting point.

You should see `{"message":"welcome aboard"}`. The route works, but it accepts literally anything — no validation at all. Let's fix that.

### Step 2 — Adding validators

Install ginvalidator:

```bash
go get github.com/bube054/ginvalidator
```

Now update `main.go` to validate an `email` and a `username` field from the request body:

```go
package main

import (
	"net/http"

	gv "github.com/bube054/ginvalidator"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/signup",
		gv.NewBodyChain("email", nil).
			Not().Empty(nil).
			Bail().
			Email(nil).
			Validate(),
		gv.NewBodyChain("username", nil).
			Not().Empty(nil).
			Bail().
			Alphanumeric(nil).
			Validate(),
		func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "welcome aboard"})
		},
	)

	r.Run()
}
```

Let's unpack what's going on:

- `gv.NewBodyChain("email", nil)` creates a validation chain for the `"email"` field in the request body. The second argument is an optional error formatter — `nil` means "use defaults".
- `.Not().Empty(nil)` means "this field must not be empty". `Empty` checks if a string is empty, and `Not()` flips the result — so an empty string fails.
- `.Bail()` tells the chain to stop if anything before it failed. Without this, the chain would keep going and run `Email` on an empty string, which would give you a second, redundant error.
- `.Email(nil)` checks that the value is a valid email address.
- `.Validate()` finishes the chain and returns a `gin.HandlerFunc` that you plug into the route.

Each chain is its own middleware. Gin runs them left to right before your handler.

> **Note:** `gv` is used as an alias for `ginvalidator` throughout these examples. You'll also see `vgo` used as an alias for `validatorgo` later on.

Now restart the server (Ctrl+C, then `go run main.go` — or let [Air](https://www.bytesizego.com/blog/golang-air) handle it) and try a bad request:

```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"email": "nope", "username": ""}'
```

On Windows Command Prompt:

```cmd
curl -X POST http://localhost:8080/signup ^
  -H "Content-Type: application/json" ^
  -d "{\"email\": \"nope\", \"username\": \"\"}"
```

Hmm, you still get `{"message":"welcome aboard"}`. Why? Because ginvalidator records errors but doesn't reject the request for you — that's your job. Let's handle the errors.

### Step 3 — Checking for errors

Update the handler to check for validation errors before responding:

```go
r.POST("/signup",
	gv.NewBodyChain("email", nil).
		Not().Empty(nil).
		Bail().
		Email(nil).
		Validate(),
	gv.NewBodyChain("username", nil).
		Not().Empty(nil).
		Bail().
		Alphanumeric(nil).
		Validate(),
	func(ctx *gin.Context) {
		result, err := gv.ValidationResult(ctx)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			ctx.Abort()
			return
		}

		if len(result) > 0 {
			ctx.IndentedJSON(http.StatusUnprocessableEntity, gin.H{
				"errors": result,
			})
			ctx.Abort()
			return
		}

		ctx.IndentedJSON(http.StatusOK, gin.H{"message": "welcome aboard"})
	},
)
```

`ValidationResult(ctx)` pulls every recorded error out of the Gin context and returns them as a sorted slice. If there are errors, we send them back. If not, we proceed.

Restart the server (Ctrl+C, then `go run main.go` — or let [Air](https://www.bytesizego.com/blog/golang-air) handle it) and try the same bad request:

```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"email": "nope", "username": ""}'
```

On Windows Command Prompt:

```cmd
curl -X POST http://localhost:8080/signup ^
  -H "Content-Type: application/json" ^
  -d "{\"email\": \"nope\", \"username\": \"\"}"
```

> **Tip:** If the JSON comes back as a single hard-to-read line, Gin has `ctx.IndentedJSON` which pretty-prints the output. We use it in this example so you can read the response easily.

Now you get:

```json
{
  "errors": [
    {
      "location": "body",
      "message": "invalid email",
      "field": "email",
      "value": "nope",
      "code": "invalid_format"
    },
    {
      "location": "body",
      "message": "Invalid value",
      "field": "username",
      "value": ""
    }
  ]
}
```

Two errors, one for each field. Let's look at what each piece means:

- **email** `"nope"`: `Not().Empty()` passed (it's not empty), so the chain continued past `Bail()`. Then `Email()` failed — `"nope"` isn't a valid email. The `message` (`"invalid email"`) and `code` (`"invalid_format"`) were provided automatically by [validatorgo](https://github.com/bube054/validatorgo).
- **username** `""`: `Not().Empty()` failed — the string IS empty, and `Not()` flipped that into a failure. `Bail()` stopped the chain right there, so `Alphanumeric` never ran. Since `Empty` technically *passed* before `Not()` negated it, there's no validatorgo error message to use, so it falls back to `"Invalid value"`.

Try a valid request now:

```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"email": "john@example.com", "username": "john123"}'
```

On Windows Command Prompt:

```cmd
curl -X POST http://localhost:8080/signup ^
  -H "Content-Type: application/json" ^
  -d "{\"email\": \"john@example.com\", \"username\": \"john123\"}"
```

`{"message":"welcome aboard"}` — it works!

### Step 4 — Better error messages

That `"Invalid value"` for the username isn't very helpful. Let's write a custom error formatter. It's a function you pass as the second argument to `NewBodyChain`:

```go
gv.NewBodyChain("username",
	func(initialValue, sanitizedValue, validatorName string) string {
		switch validatorName {
		case gv.EmptyValidatorName:
			return "Username can't be blank."
		case gv.AlphanumericValidatorName:
			return "Username can only contain letters and numbers."
		default:
			return "Invalid username."
		}
	},
).
	Not().Empty(nil).
	Bail().
	Alphanumeric(nil).
	Validate()
```

The function receives three things:

- `initialValue` — the original value from the request
- `sanitizedValue` — the value after any sanitizers have run (more on this later)
- `validatorName` — which validator failed, e.g. `"Empty"`, `"Email"`, `"Alphanumeric"`

You can use `validatorName` to return different messages for different failures on the same field. The full list of validator names is in the [ginvalidator constants](https://pkg.go.dev/github.com/bube054/ginvalidator#pkg-constants).

### Step 5 — Sanitizing input

Validation tells you if the data is good. Sanitization cleans the data up.

Right now, a user could send something like `<script>alert('hacked')</script>` as a username. If you ever render that in HTML, you've got a [Cross-Site Scripting (XSS)](https://owasp.org/www-community/attacks/xss/) vulnerability. The `Escape` sanitizer fixes this by converting special HTML characters into safe equivalents (e.g. `<` becomes `&lt;`):

```go
gv.NewBodyChain("username", nil).
	Trim("").
	Not().Empty(nil).
	Escape().          // HTML-escape to prevent XSS
	Validate()
```

Try it — send a script tag as the username:

```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"email": "john@example.com", "username": "<script>alert(1)</script>"}'
```

On Windows Command Prompt (`<` and `>` are redirection operators in cmd, so they need to be escaped with `^`):

```cmd
curl -X POST http://localhost:8080/signup ^
  -H "Content-Type: application/json" ^
  -d "{\"email\": \"john@example.com\", \"username\": \"^<script^>alert(1)^</script^>\"}"
```

You'll get `{"message":"welcome aboard"}` — the username isn't empty, so it passes. But behind the scenes, `Escape` transformed the value into `&lt;script&gt;alert(1)&lt;&#x2F;script&gt;` before storing it. We'll see the actual sanitized value in the next step when we read matched data.

We also added `Trim` before the empty check. Order matters here — if someone sends `"   "` (just spaces), we want to strip those first, *then* check if it's empty. If we did it the other way around, `"   "` would pass the empty check and we'd end up with an empty username.

> **Important:** Sanitizers don't modify the original `http.Request`. To access the sanitized value, you need to use `GetMatchedData` — we'll do that next.

### Step 6 — Reading the validated data

After validation and sanitization, ginvalidator stores the final values for you. Pull them out with `GetMatchedData`:

```go
r.POST("/signup",
	gv.NewBodyChain("email", nil).
		Not().Empty(nil).
		Bail().
		Email(nil).
		Validate(),
	gv.NewBodyChain("username", nil).
		Trim("").
		Not().Empty(nil).
		Escape().
		Validate(),
	func(ctx *gin.Context) {
		result, err := gv.ValidationResult(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		if len(result) > 0 {
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"errors": result,
			})
			return
		}

		data, err := gv.GetMatchedData(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		email, _ := data.Get(gv.BodyLocation, "email")
		username, _ := data.Get(gv.BodyLocation, "username")

		// email and username are validated + sanitized, safe to use
		ctx.JSON(http.StatusOK, gin.H{
			"message":  "welcome aboard",
			"email":    email,
			"username": username,
		})
	},
)
```

`data.Get(location, fieldName)` returns the final value and a boolean. The first argument is the location constant — `BodyLocation`, `QueryLocation`, `ParamLocation`, `HeaderLocation`, or `CookieLocation` — and the second is the field name.

That's the core loop: **create chains** → **check errors** → **read matched data**. Everything else builds on this.

---

## Validation chains in depth

We've been using `NewBodyChain` so far. There are five of these, one for each place data can come from in an HTTP request:

| Constructor | Location constant | Reads from |
|---|---|---|
| `NewBodyChain` | `BodyLocation` | Request body (JSON, form, multipart) |
| `NewQueryChain` | `QueryLocation` | URL query parameters |
| `NewParamChain` | `ParamLocation` | Gin route parameters (`:id`) |
| `NewHeaderChain` | `HeaderLocation` | HTTP headers |
| `NewCookieChain` | `CookieLocation` | Cookies |

They all work the same way — the only difference is where they look for the field value.

> **Just a shortcut:** Each `NewXChain` is shorthand for `NewX(...).Chain()`. So `NewBodyChain("email", nil)` is the same as `NewBody("email", nil).Chain()`, `NewQueryChain("q", nil)` is the same as `NewQuery("q", nil).Chain()`, and so on for `NewParam`, `NewHeader` and `NewCookie`. The longer form is still there if you ever need to hold on to the intermediate value.

A few things worth knowing about field names:

- **JSON body fields use [GJSON path syntax](https://github.com/tidwall/gjson#path-syntax)** — so for `{"user":{"profile":{"email":"a@b.c"}}}`, the field name is `"user.profile.email"`. You're not limited to top-level keys.
- **Body extraction switches on `Content-Type`** — JSON uses GJSON paths, while `application/x-www-form-urlencoded` and `multipart/form-data` use plain form field names.
- **Headers must be in canonical form** — use `"Content-Type"`, not `"content-type"`. ginvalidator will log a warning if you pass a non-canonical key.

Here's a route that validates a route parameter and an optional query parameter:

```go
r.GET("/users/:id",
	gv.NewParamChain("id", nil).
		Not().Empty(nil).
		Bail().
		Numeric(nil).
		Validate(),
	gv.NewQueryChain("fields", nil).
		Optional().
		Alpha(nil).
		Validate(),
	func(ctx *gin.Context) {
		result, err := gv.ValidationResult(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		if len(result) > 0 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": result})
			return
		}

		data, _ := gv.GetMatchedData(ctx)
		id, _ := data.Get(gv.ParamLocation, "id")
		ctx.JSON(http.StatusOK, gin.H{"user_id": id})
	},
)
```

Try it out — these GET requests work the same in bash and Command Prompt:

```bash
# valid: numeric id, alpha fields
curl http://localhost:8080/users/123?fields=name

# valid: numeric id, fields omitted (Optional skips the chain)
curl http://localhost:8080/users/123

# invalid: non-numeric id
curl http://localhost:8080/users/abc

# invalid: fields contains digits
curl http://localhost:8080/users/123?fields=name1
```

### Reusing chains

If you use the same validation in multiple places, wrap it in a function:

```go
func validateEmail() gin.HandlerFunc {
	return gv.NewBodyChain("email", nil).
		Not().Empty(nil).
		Bail().
		Email(nil).
		Validate()
}

r.POST("/login", validateEmail(), loginHandler)
r.POST("/signup", validateEmail(), signupHandler)
```

## Validators

Validators check if a field value meets some criteria. When one fails, an error is recorded. It doesn't reject the request — you decide what to do with the errors in your handler.

All built-in validators are powered by [validatorgo](https://github.com/bube054/validatorgo) — my other open source library that I originally built for this very project. ginvalidator wraps each one as a chain method: validatorgo's `IsEmail` becomes `.Email()`, `IsAlphanumeric` becomes `.Alphanumeric()`, `IsEmpty` becomes `.Empty()`, and so on (the `Is` prefix is dropped).

Here are some commonly used ones:

- `Email`, `URL`, `IP`, `UUID` — format checks
- `Alpha`, `Alphanumeric`, `Numeric`, `Int`, `Float` — character/number checks
- `Empty`, `Contains`, `Equals`, `Matches` — string checks
- `CreditCard`, `MobilePhone`, `PostalCode` — domain-specific checks
- `StrongPassword` — password strength
- `JSON`, `Boolean`, `Date`, `ISO8601` — data type checks

The full list is in the [ginvalidator constants](https://pkg.go.dev/github.com/bube054/ginvalidator#pkg-constants). For detailed documentation on what each validator does and its options, see [validatorgo on pkg.go.dev](https://pkg.go.dev/github.com/bube054/validatorgo).

### Validator options

Many validators accept an options struct to tweak their behavior. These structs come from [validatorgo](https://pkg.go.dev/github.com/bube054/validatorgo). Pass `nil` when the defaults are fine (which is most of the time):

```go
gv.NewBodyChain("email", nil).
	Email(nil).
	Validate()
```

When you need to customize, import validatorgo and pass the struct. Note that these structs use pointer types (`*bool`, `*int`, `*string`) so the library can tell the difference between "not set" and "set to the zero value" — see [validatorgo's option fields docs](https://github.com/bube054/validatorgo#option-fields) for the full rationale:

```go
import vgo "github.com/bube054/validatorgo"

gv.NewBodyChain("email", nil).
	Email(&vgo.IsEmailOpts{
		RequireTld:    vgo.Bool(false), // helper to create a *bool pointer
		HostWhitelist: []string{"gmail.com", "yahoo.com"},
	}).
	Validate()
```

validatorgo provides little helpers like `vgo.Bool`, `vgo.String`, `vgo.Int` and `vgo.Float64` that take a value and return a pointer to it — so you don't have to declare an intermediate variable just to take its address. If you use validatorgo structs directly in your code, run `go mod tidy` so Go adds it as a direct dependency.

### CustomValidator

For logic that built-in validators can't cover — like checking if an email is already taken:

```go
// pretend this hits a database
func userExistsByEmail(email string) bool {
	taken := map[string]bool{"john@example.com": true}
	return taken[email]
}

r.POST("/create-user",
	gv.NewBodyChain("email", nil).
		Not().Empty(nil).
		Bail().
		Email(nil).
		CustomValidator(func(r *http.Request, initialValue, sanitizedValue string) bool {
			return !userExistsByEmail(sanitizedValue)
		}).
		Validate(),
	func(ctx *gin.Context) {
		result, _ := gv.ValidationResult(ctx)
		if len(result) > 0 {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"errors": result})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{"message": "user created"})
	},
)
```

Try it:

```bash
# fresh email → passes
curl -X POST http://localhost:8080/create-user \
  -H "Content-Type: application/json" \
  -d '{"email": "jane@example.com"}'
# {"message":"user created"}

# email already taken → CustomValidator returns false
curl -X POST http://localhost:8080/create-user \
  -H "Content-Type: application/json" \
  -d '{"email": "john@example.com"}'
```

The function signature:

```go
type CustomValidatorFunc func(r *http.Request, initialValue, sanitizedValue string) bool
```

- `r` — the raw HTTP request
- `initialValue` — the value before any sanitization
- `sanitizedValue` — the value after sanitizers that ran earlier in the chain
- Return `true` if valid, `false` if not

`CustomValidator` doesn't produce error `code`s (since there's no validatorgo validator behind it).

## Sanitizers

Sanitizers transform the field value. The transformed value is what later validators in the chain see, and what you get back from `GetMatchedData`.

Built-in sanitizers (also from [validatorgo](https://github.com/bube054/validatorgo)):

`Trim`, `LTrim`, `RTrim`, `Escape`, `Unescape`, `Blacklist`, `Whitelist`, `NormalizeEmail`, `StripLow`, `ToBoolean`, `ToDate`, `ToFloat`, `ToInt`

We already saw `Trim` and `Escape` in the Getting Started section. Here's another example — normalizing an email address:

```go
r.POST("/subscribe",
	gv.NewBodyChain("email", nil).
		Trim("").
		NormalizeEmail(nil).
		Not().Empty(nil).
		Bail().
		Email(nil).
		Validate(),
	func(ctx *gin.Context) {
		result, _ := gv.ValidationResult(ctx)
		if len(result) > 0 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": result})
			return
		}

		data, _ := gv.GetMatchedData(ctx)
		email, _ := data.Get(gv.BodyLocation, "email")
		// email is trimmed + normalized (e.g. "John@Gmail.com" → "john@gmail.com")

		ctx.JSON(http.StatusOK, gin.H{"subscribed": email})
	},
)
```

Try it with a messy email:

```bash
curl -X POST http://localhost:8080/subscribe \
  -H "Content-Type: application/json" \
  -d '{"email": "  John@Gmail.com  "}'
```

On Windows Command Prompt:

```cmd
curl -X POST http://localhost:8080/subscribe ^
  -H "Content-Type: application/json" ^
  -d "{\"email\": \"  John@Gmail.com  \"}"
```

You'll get back:

```json
{"subscribed":"john@gmail.com"}
```

The whitespace was trimmed and the email was normalized to lowercase before being stored.

> **Remember:** Sanitizers don't modify the original `http.Request`. Always use `GetMatchedData` to read sanitized values.

### CustomSanitizer

```go
type CustomSanitizerFunc func(r *http.Request, initialValue, sanitizedValue string) string
```

Whatever string you return becomes the new value:

```go
r.POST("/articles",
	gv.NewBodyChain("slug", nil).
		CustomSanitizer(func(r *http.Request, initialValue, sanitizedValue string) string {
			return strings.ToLower(strings.ReplaceAll(sanitizedValue, " ", "-"))
		}).
		Validate(),
	func(ctx *gin.Context) {
		data, _ := gv.GetMatchedData(ctx)
		slug, _ := data.Get(gv.BodyLocation, "slug")
		ctx.JSON(http.StatusOK, gin.H{"slug": slug})
	},
)
```

Try it:

```bash
curl -X POST http://localhost:8080/articles \
  -H "Content-Type: application/json" \
  -d '{"slug": "My Blog Post"}'
# {"slug":"my-blog-post"}
```

## Modifiers

Modifiers don't validate or transform — they control how the chain behaves.

### Not

Negates the next validator. A pass becomes a fail, a fail becomes a pass.

```go
// "this field must NOT be empty"
r.POST("/check-name",
	gv.NewBodyChain("name", nil).
		Not().Empty(nil).
		Validate(),
	func(ctx *gin.Context) {
		if gv.HasErrors(ctx) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
	},
)
```

Try it:

```bash
# empty name → fails
curl -X POST http://localhost:8080/check-name \
  -H "Content-Type: application/json" \
  -d '{"name": ""}'
# {"error":"name is required"}

# non-empty name → passes
curl -X POST http://localhost:8080/check-name \
  -H "Content-Type: application/json" \
  -d '{"name": "alice"}'
# {"message":"ok"}
```

Step by step: `Empty` checks if the string is empty. If the string IS empty, `Empty` returns true ("yes it's empty"). `Not()` flips that to false — validation fails. If the string is NOT empty, `Empty` returns false, `Not()` flips it to true — validation passes.

`Not()` only affects the single validator right after it.

### Bail

Stops the chain if any previous validator failed. Put it after any validator where continuing doesn't make sense:

```go
r.POST("/transfer",
	gv.NewBodyChain("amount", nil).
		Not().Empty(nil).
		Bail().              // stop if empty
		Numeric(nil).
		Bail().              // stop if not numeric
		CustomValidator(func(r *http.Request, initialValue, sanitizedValue string) bool {
			// no point hitting the DB if the value isn't even a number
			return sanitizedValue == "100" // pretend only 100 is affordable
		}).
		Validate(),
	func(ctx *gin.Context) {
		if gv.HasErrors(ctx) {
			result, _ := gv.ValidationResult(ctx)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": result})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "transferred"})
	},
)
```

Try it:

```bash
# empty → first Bail stops the chain. only one error.
curl -X POST http://localhost:8080/transfer \
  -H "Content-Type: application/json" \
  -d '{"amount": ""}'

# not numeric → passes Empty, fails Numeric. second Bail stops the chain.
curl -X POST http://localhost:8080/transfer \
  -H "Content-Type: application/json" \
  -d '{"amount": "abc"}'

# numeric but insufficient → all the way to CustomValidator.
curl -X POST http://localhost:8080/transfer \
  -H "Content-Type: application/json" \
  -d '{"amount": "50"}'

# happy path
curl -X POST http://localhost:8080/transfer \
  -H "Content-Type: application/json" \
  -d '{"amount": "100"}'
# {"message":"transferred"}
```

Without `Bail()`, every validator would run on every request, and you'd get errors stacked up (plus you'd hit the DB even when the value isn't a number). You can use `Bail()` multiple times in the same chain.

### Optional

Skips the entire chain if the field is empty. You can put it anywhere in the chain — position doesn't matter.

```go
// bio is optional — if empty, no validators run.
// if the user does send a value, it must be alpha characters only.
r.POST("/profile",
	gv.NewBodyChain("bio", nil).
		Optional().
		Alpha(nil).
		Validate(),
	func(ctx *gin.Context) {
		if gv.HasErrors(ctx) {
			result, _ := gv.ValidationResult(ctx)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": result})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "saved"})
	},
)
```

Try it:

```bash
# bio omitted → Optional skips the chain → passes
curl -X POST http://localhost:8080/profile \
  -H "Content-Type: application/json" \
  -d '{}'
# {"message":"saved"}

# bio is alpha → passes
curl -X POST http://localhost:8080/profile \
  -H "Content-Type: application/json" \
  -d '{"bio": "hello"}'
# {"message":"saved"}

# bio has digits → Alpha fails
curl -X POST http://localhost:8080/profile \
  -H "Content-Type: application/json" \
  -d '{"bio": "hello123"}'
```

### If

Conditionally stops the chain based on a function you provide. Return `true` to stop (bail out), `false` to continue.

```go
type IfModifierFunc func(r *http.Request, initialValue, sanitizedValue string) bool
```

```go
// only validate the discount code if the user is a premium member
r.POST("/checkout",
	gv.NewBodyChain("discountCode", nil).
		If(func(r *http.Request, initialValue, sanitizedValue string) bool {
			// return true = stop the chain, false = keep going
			return r.Header.Get("X-User-Tier") != "premium"
		}).
		Not().Empty(nil).
		Bail().
		Alphanumeric(nil).
		Validate(),
	func(ctx *gin.Context) {
		if gv.HasErrors(ctx) {
			result, _ := gv.ValidationResult(ctx)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": result})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "checked out"})
	},
)
```

Try it:

```bash
# no premium header → If returns true → chain stops → empty discount is fine
curl -X POST http://localhost:8080/checkout \
  -H "Content-Type: application/json" \
  -d '{"discountCode": ""}'
# {"message":"checked out"}

# premium user with empty code → If returns false → validators run → fails
curl -X POST http://localhost:8080/checkout \
  -H "Content-Type: application/json" \
  -H "X-User-Tier: premium" \
  -d '{"discountCode": ""}'

# premium user with valid code → passes
curl -X POST http://localhost:8080/checkout \
  -H "Content-Type: application/json" \
  -H "X-User-Tier: premium" \
  -d '{"discountCode": "SAVE20"}'
# {"message":"checked out"}
```

If `If` returns `true`, the chain stops right there — no validators after it run, no errors recorded.

### Skip

Skips just the next item in the chain (validator, sanitizer, or modifier). Return `true` to skip it, `false` to run it.

```go
type SkipModifierFunc func(r *http.Request, initialValue, sanitizedValue string) bool
```

```go
// skip the length check when the request comes from an admin
r.POST("/post-message",
	gv.NewBodyChain("message", nil).
		Not().Empty(nil).
		Skip(func(r *http.Request, initialValue, sanitizedValue string) bool {
			return r.Header.Get("X-Role") == "admin"
		}).
		ByteLength(&vgo.IsByteLengthOpts{Max: vgo.Int(10)}). // skipped if Skip returned true
		Validate(),
	func(ctx *gin.Context) {
		if gv.HasErrors(ctx) {
			result, _ := gv.ValidationResult(ctx)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": result})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "posted"})
	},
)
```

Try it:

```bash
# regular user with short message → passes
curl -X POST http://localhost:8080/post-message \
  -H "Content-Type: application/json" \
  -d '{"message": "hello"}'
# {"message":"posted"}

# regular user with long message → fails ByteLength
curl -X POST http://localhost:8080/post-message \
  -H "Content-Type: application/json" \
  -d '{"message": "this message is way too long"}'

# admin with the same long message → ByteLength is skipped → passes
curl -X POST http://localhost:8080/post-message \
  -H "Content-Type: application/json" \
  -H "X-Role: admin" \
  -d '{"message": "this message is way too long"}'
# {"message":"posted"}
```

Unlike `If`, `Skip` only skips the next item — everything after it still runs.

## Error messages

When a validator fails, ginvalidator picks the error message using this priority:

1. **Per-chain formatter** — the function you pass as the second argument to `NewBodyChain`, `NewQueryChain`, etc. (we covered this in [Step 4](#step-4--better-error-messages))
2. **`DefaultErrFmtFunc`** — a package-level formatter you can set once for your whole app
3. **validatorgo message** — the [validatorgo](https://github.com/bube054/validatorgo) validator returns a `ValidationError` with a `Message` field (like `"invalid email"`). If nothing above is set, this is used.
4. **`"Invalid value"`** — the last-resort fallback

### DefaultErrFmtFunc

If you're tired of writing formatters for every single chain, set a global default. It applies to any chain that doesn't have its own formatter:

```go
func main() {
	gv.DefaultErrFmtFunc = func(initialValue, sanitizedValue, validatorName string) string {
		return fmt.Sprintf("%s check failed", validatorName)
	}

	r := gin.Default()
	// ... all your routes will use this formatter as the default
	r.Run()
}
```

Per-chain formatters still win when present.

### Error codes

You may have noticed the `code` field in some of the error responses earlier. When a built-in validator fails, [validatorgo](https://github.com/bube054/validatorgo) returns a `ValidationError` that looks like this:

```go
// from validatorgo
type ValidationError struct {
	Validator string // e.g. "IsEmail"
	Code      string // e.g. "invalid_format"
	Message   string // e.g. "invalid email"
}
```

ginvalidator reads the `Code` and `Message` from this error and puts them into your validation results. The `code` field is `omitempty` in JSON, so it only shows up when there's actually a code. `CustomValidator` doesn't produce codes since there's no validatorgo validator behind it.

Understanding [validatorgo's error types](https://pkg.go.dev/github.com/bube054/validatorgo) will help you make the most of these codes — they're handy for i18n or building client-side error handling.

## Reading errors

We've been using `ValidationResult` to get all errors as a slice. That works, but ginvalidator also has helpers for common patterns.

For the examples below, assume each handler is plugged into a route like this so you have something to curl against:

```go
r.POST("/signup",
	gv.NewBodyChain("email", nil).
		Not().Empty(nil).
		Bail().
		Email(nil).
		Validate(),
	gv.NewBodyChain("username", nil).
		Not().Empty(nil).
		Bail().
		Alphanumeric(nil).
		Validate(),
	signupHandler, // <-- swap in the handler from each example below
)
```

And the bad request we'll use to trigger errors:

```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"email": "nope", "username": ""}'
```

### HasErrors

The simplest check — just a boolean:

```go
func signupHandler(ctx *gin.Context) {
	if gv.HasErrors(ctx) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "validation failed",
		})
		return
	}

	// all good, do your thing
	ctx.JSON(http.StatusCreated, gin.H{"message": "signed up"})
}
```

Response:

```json
{"error":"validation failed"}
```

### FirstError

When you only want to show one error at a time:

```go
func signupHandler(ctx *gin.Context) {
	if err := gv.FirstError(ctx); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"field":   err.Field,
			"message": err.Msg,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "signed up"})
}
```

Response (only the first error — the email one):

```json
{"field":"email","message":"invalid email"}
```

`FirstError` returns a pointer to the first `ValidationChainError`, or `nil` if everything passed.

### ErrorsByField

Groups all errors by field name — useful when your UI shows a list of errors under each form input:

```go
func signupHandler(ctx *gin.Context) {
	grouped, err := gv.ErrorsByField(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	if len(grouped) > 0 {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"errors": grouped})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "signed up"})
}
```

Response:

```json
{
  "errors": {
    "email": [
      {"location":"body","message":"invalid email","field":"email","value":"nope","code":"invalid_format"}
    ],
    "username": [
      {"location":"body","message":"Invalid value","field":"username","value":""}
    ]
  }
}
```

### FirstErrorByField

Same idea, but at most one error per field. Common for "show one error per input" UIs:

```go
func signupHandler(ctx *gin.Context) {
	firsts, err := gv.FirstErrorByField(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	if len(firsts) > 0 {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"errors": firsts})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "signed up"})
}
```

Response (note each field maps to a single error object, not an array):

```json
{
  "errors": {
    "email":    {"location":"body","message":"invalid email","field":"email","value":"nope","code":"invalid_format"},
    "username": {"location":"body","message":"Invalid value","field":"username","value":""}
  }
}
```

## Matched data

We covered `GetMatchedData` in [Step 6](#step-6--reading-the-validated-data), but here's a quick recap of the two methods it gives you:

**`Get(location, field)`** — returns the value and a boolean:

```go
data, _ := gv.GetMatchedData(ctx)
email, ok := data.Get(gv.BodyLocation, "email")
if !ok {
	// field wasn't matched
}
```

**`Has(location, field)`** — just checks if the field was matched, without pulling the value. Useful for optional fields:

```go
if data.Has(gv.BodyLocation, "bio") {
	bio, _ := data.Get(gv.BodyLocation, "bio")
	// user sent a bio, do something with it
}
```

## OneOf

Sometimes a request is valid if *any one of several groups* of validations passes. A login that accepts either an email or a phone number is a classic example:

```go
r.POST("/login",
	gv.OneOf(
		[]gv.ValidationChain{
			gv.NewBodyChain("email", nil).
				Not().Empty(nil).
				Email(nil),
		},
		[]gv.ValidationChain{
			gv.NewBodyChain("phone", nil).
				Not().Empty(nil).
				MobilePhone(nil, ""),
		},
	),
	func(ctx *gin.Context) {
		if gv.HasErrors(ctx) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "provide either a valid email or phone number",
			})
			return
		}

		data, _ := gv.GetMatchedData(ctx)
		if data.Has(gv.BodyLocation, "email") {
			email, _ := data.Get(gv.BodyLocation, "email")
			ctx.JSON(http.StatusOK, gin.H{"login_via": "email", "value": email})
			return
		}
		phone, _ := data.Get(gv.BodyLocation, "phone")
		ctx.JSON(http.StatusOK, gin.H{"login_via": "phone", "value": phone})
	},
)
```

If at least one group produces zero errors, the request passes and that group's matched data is saved. If every group fails, a single error with field `"_oneOf"` is recorded.

You can put multiple chains in one group — they all have to pass for that group to count:

```go
gv.OneOf(
	// group 1: both name and email must be valid
	[]gv.ValidationChain{
		gv.NewBodyChain("name", nil).
			Not().Empty(nil),
		gv.NewBodyChain("email", nil).
			Not().Empty(nil).
			Email(nil),
	},
	// group 2: just a username
	[]gv.ValidationChain{
		gv.NewBodyChain("username", nil).
			Not().Empty(nil).
			Alphanumeric(nil),
	},
)
```

## CheckSchema

When a route has a lot of fields, writing individual chains for each one gets long. `CheckSchema` lets you define everything in a single map and gives you back one middleware:

```go
r.POST("/register",
	gv.CheckSchema(gv.Schema{
		"email": {
			In: gv.BodyLocation,
			Build: func(vc gv.ValidationChain) gv.ValidationChain {
				return vc.Not().Empty(nil).Bail().Email(nil)
			},
		},
		"username": {
			In: gv.BodyLocation,
			Build: func(vc gv.ValidationChain) gv.ValidationChain {
				return vc.Not().Empty(nil).Bail().Alphanumeric(nil)
			},
		},
		"bio": {
			In:       gv.BodyLocation,
			Optional: true,
			Build: func(vc gv.ValidationChain) gv.ValidationChain {
				return vc.Trim("").Escape()
			},
		},
	}),
	func(ctx *gin.Context) {
		if gv.HasErrors(ctx) {
			result, _ := gv.ValidationResult(ctx)
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"errors": result})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "registered"})
	},
)
```

Try it:

```bash
# all good — bio is optional and gets escaped
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email": "jane@example.com", "username": "jane123", "bio": "hi there"}'
# {"message":"registered"}

# bio omitted — Optional skips its chain, no errors
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email": "jane@example.com", "username": "jane123"}'
# {"message":"registered"}

# multiple things wrong — both email and username fail
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email": "nope", "username": "jane!"}'
```

On Windows Command Prompt:

```cmd
curl -X POST http://localhost:8080/register ^
  -H "Content-Type: application/json" ^
  -d "{\"email\": \"jane@example.com\", \"username\": \"jane123\", \"bio\": \"hi there\"}"
```

Each field in the schema gets a `SchemaField`:

- **`In`** — where the field comes from (`BodyLocation`, `QueryLocation`, `ParamLocation`, `HeaderLocation`, `CookieLocation`)
- **`Optional`** — if `true`, skip validation when the field is empty
- **`ErrFmtFunc`** — per-field error formatter (same type as the second argument to `NewBodyChain`)
- **`Build`** — receives a fresh `ValidationChain`, return it with your validators/sanitizers/modifiers attached. Use `Bail()` inside `Build` to stop on first failure. If `Build` is `nil`, the field always passes.

Fields are processed in alphabetical order, so errors come back in a predictable order.

## Maintainers

- [bube054](https://github.com/bube054) — **Attah Gbubemi David** (author)

## License

[MIT](LICENSE)
