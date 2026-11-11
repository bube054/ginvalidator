# Changelog

All notable changes to `ginvalidator`. Versions follow [semantic versioning](https://semver.org).

## v1.2.0 (2026-11-11)

### Changed
- Gin v1.12.0, up from v1.10.1.
- Requires Go 1.25 or newer, up from 1.23, because Gin v1.12.0 does. Go 1.23 and 1.24 no longer get security fixes from the Go team.

## v1.1.0 (2026-11-11)

### Fixed
- Bodies sent with a `Content-Type` that has parameters or other casing, such as `application/json; charset=utf-8`, are read correctly. Before, every body field read as empty. `+json` types and forms with a charset also work.
- Path parameters are no longer decoded a second time: `100%` no longer errors, and `a+b` no longer becomes `a b`.
- A missing cookie is an empty value, like a missing header or query parameter, instead of an error.
- The body is read once per request and shared by every chain, instead of re-read by each one.
- The library no longer writes to the log.
- validatorgo v1.1.0: fixes a panic in `Email()` on input such as `junk [1.2.3.4]`, and validators run about 50-250x faster.

### Behavior changes
- `GetMatchedData` only contains fields that were sent and passed every rule. Values that failed and fields that were never sent are left out, so `Has` returns false for them.
- A JSON body that isn't valid JSON gives each body chain one error with code `invalid_json` (`InvalidJSONErrCode`), even for `Optional()` chains. Before, broken JSON could pass validation.
- A body that can't be read gives code `unreadable_request` (`UnreadableRequestErrCode`).
- A body that is neither JSON nor a form, such as `text/plain`, has no fields: body fields count as not sent.
- validatorgo v1.1.0 is stricter in `Email`, `IP`, `Port`, `ULID` and `RgbColor`; see [its changelog](https://github.com/bube054/validatorgo/blob/main/CHANGELOG.md).

### Other
- A logo, a contributor guide (`UNDERSTANDING_THE_CODEBASE.md`), and integration tests that check results again.

## v1.0.0 (2026-05-09)

### Breaking changes
- Renamed `DefaultValChainErrMsg` to `DefaultErrMsg` and `ErrFmtFuncHandler` to `ErrFmtFunc`.
- Removed `NewValidationChain`, `NewValidationChainError`, `NewValidationChainRule` and `SortValidationErrors` from the public API.
- Requires Go 1.23 and Gin v1.10.1.

### New
- Errors include validatorgo's message and a machine-readable `Code`, such as `invalid_format`.
- `NewBodyChain`, `NewQueryChain`, `NewParamChain`, `NewHeaderChain` and `NewCookieChain` shorthands.
- `OneOf`: passes if any one group of chains passes.
- `CheckSchema`: declare a request's rules as a `Schema` instead of chains.
- `HasErrors`, `FirstError`, `ErrorsByField`, `FirstErrorByField` and `MatchedData.Has`.
- `DefaultErrFmtFunc`, a package-wide default error formatter.

### Other
- validatorgo v1.0.0, CI, and a README rewritten as a step-by-step tutorial.

## v0.5.0 (2025-01-11)

- validatorgo v0.3.1.

## v0.4.0 (2025-01-11)

- validatorgo v0.3.0.

## v0.3.1 (2025-01-11)

### New
- `Length`, `Array` and `Object` validators.
- validatorgo v0.2.0, up from v0.1.0-alpha.

## v0.3.0 (2024-12-09)

- Same code as v0.2.2.

## v0.2.2 (2024-12-09)

- Documentation on configuring chains.

## v0.2.1 (2024-11-19)

- Documentation on reading validated and sanitized data.

## v0.2.0 (2024-11-13)

### New
- `MatchedData.Get`.

## v0.1.1 (2024-11-13)

- README fixes.

## v0.1.0 (2024-11-12)

First release: validation chains for the body, query, path params, headers and cookies, with validators and sanitizers from validatorgo, the modifiers `Bail`, `Not`, `If`, `Skip` and `Optional`, `ValidationResult` and `GetMatchedData`.
