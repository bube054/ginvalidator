# Understanding the Codebase

ginvalidator is a Gin middleware library for validating and sanitizing HTTP request fields. It wraps [validatorgo](https://github.com/bube054/validatorgo) in a fluent chain API. This document explains how the code is organized and the best order to read it.

## The One Concept You Need

Everything in this repo revolves around one type:

```go
type ruleCreatorFunc func(ctx *gin.Context, initialValue, sanitizedValue string) validationChainRule
```

A `ruleCreatorFunc` is a closure that, given a request context and the current field value, returns a result: did validation pass? what's the new sanitized value? Every validator, sanitizer, and modifier in the library just creates one of these closures and appends it to a list. When a request arrives, the list runs in order.

Once you understand this, you understand the whole codebase.

## Recommended Reading Order

### 1. `validationerror.go` — the output shape

Start here. It's small and shows you `ValidationChainError` — the struct the entire library exists to produce. Fields: `Location`, `Field`, `Value`, `Message`, `Code`. Now you know what validation produces.

### 2. `rule.go` — the building block

This defines `validationChainRule` (what a single rule produces) and `ruleCreatorFunc` (the closure type). It uses a functional options pattern (`withIsValid`, `withNewValue`, etc.) to construct rules. Every other file creates these.

### 3. `validator.go` — the pattern (read one method)

Open this file and read a single method like `Contains()`. Don't read all 87 validators. The pattern is always:
1. Call a `validatorgo` function
2. Wrap the result in a `ruleCreatorFunc`
3. Append it to the chain via `recreateValidationChainFromValidator()`

Sanitizers in `sanitizer.go` and modifiers in `modifier.go` follow the same shape.

### 4. `validationchain.go` — the engine

This is the most important file. The `validate()` method:
1. Extracts the field value from the request
2. Loops through every `ruleCreatorFunc` in order
3. For validators: checks `isValid`, collects errors
4. For sanitizers: updates the running `sanitizedValue`
5. For modifiers: adjusts control flow (bail, negate, skip)

The `Validate()` method wraps this into a `gin.HandlerFunc`.

### 5. `requestutils.go` — how data gets in

Defines `RequestLocation` (body, query, param, header, cookie) and extraction functions for each. JSON bodies use `gjson` for path-based access (e.g., `"address.city"`). Form bodies use `ctx.PostForm()`. The body is re-wrapped after reading so downstream handlers can still access it.

### 6. `validationresult.go` and `matcheddata.go` — how data gets out

- `validationresult.go`: stores errors in the Gin context, retrieves them via `ValidationResult()`, `HasErrors()`, `FirstError()`, `ErrorsByField()`
- `matcheddata.go`: stores sanitized field values, retrieves them via `GetMatchedData()`

Both use string keys on `gin.Context` to store nested maps.

### 7. `body.go`, `query.go`, `param.go`, `header.go`, `cookie.go` — entry points

These are thin. Each one just calls `newValidationChain()` with the right `RequestLocation`. Read one, skip the rest.

### 8. `oneof.go` and `checkschema.go` — advanced features

- `OneOf()`: middleware that passes if at least one group of chains has zero errors
- `CheckSchema()`: declarative schema-based alternative to fluent chains

## Mental Model

```
Entry points              Chain building              Execution             Results
─────────────────    ──────────────────────    ──────────────────    ─────────────────
body.go                validator.go              validationchain.go    validationresult.go
query.go               sanitizer.go              requestutils.go       matcheddata.go
param.go               modifier.go
header.go
cookie.go

NewBody("email")  →  .Email().Trim().Bail()  →  .Validate() runs   →  ValidationResult(ctx)
                     (appends closures)          closures on request    GetMatchedData(ctx)
```

## Data Flow on Each Request

1. Gin calls the `HandlerFunc` returned by `.Validate()`
2. `validate()` extracts the field value from the request location
3. Each `ruleCreatorFunc` runs in order against `(initialValue, sanitizedValue)`
4. Errors are stored in `ctx` under `"__ginvalidator__ctx__errors__"` (nested map: location → field → errors)
5. Sanitized values are stored under `"__ginvalidator__ctx__matched__data__"` (map: location → field → value)
6. `ctx.Next()` is called — validation never blocks the request
7. Your handler calls `ValidationResult(ctx)` or `GetMatchedData(ctx)` to read results

## File Groups

| Files | Purpose |
|-------|---------|
| `body.go`, `query.go`, `param.go`, `header.go`, `cookie.go` | Create chains for each request location |
| `validator.go`, `validator_a_d.go`, `validator_e_i.go`, `validator_is_m.go`, `validator_n_z.go` | 87+ validators (alphabetically split) |
| `sanitizer.go` | 13 sanitizers |
| `modifier.go` | 5 modifiers: Bail, Not, Optional, If, Skip |
| `validationchain.go` | Core execution loop and middleware conversion |
| `rule.go` | Rule struct and closure type |
| `requestutils.go` | Field extraction from requests |
| `validationresult.go` | Error storage and retrieval |
| `matcheddata.go` | Sanitized data storage and retrieval |
| `validationerror.go` | Error struct and formatting |
| `oneof.go` | OneOf middleware |
| `checkschema.go` | Schema-based validation |

## Running Tests

```bash
go vet ./...                    # lint
go test -count=1 -race ./...   # run all tests (same as CI)
```

Tests follow two patterns:
- **Unit tests** (`validator_test.go`, `sanitizer_test.go`): create a chain, extract the single `ruleCreatorFunc`, run it, compare the `validationChainRule` output
- **Integration tests** (`body_test.go`, `query_test.go`, etc.): spin up a Gin router with `httptest.NewRecorder`, send a request, check `ValidationResult()` and `GetMatchedData()` in the handler
