# ginvalidator

<img alt="Tag" src="https://img.shields.io/badge/tag-v0.1.0-blue?labelColor=gray"> <img alt="Go Version" src="https://img.shields.io/badge/Go->=1.21-00ADD8?labelColor=gray"> <img alt="Reference" src="https://img.shields.io/badge/-reference-00ADD8?logo=go&labelColor=gray"> <img alt="Tests" src="https://img.shields.io/badge/tests-passing-brightgreen?logo=github&labelColor=gray"> <img alt="Go Report" src="https://img.shields.io/badge/go_report-A%2B-00ADD8"> <img alt="Coverage" src="https://img.shields.io/badge/coverage-87.30%25-brightgreen?logo=codecov"> <img alt="Contributors" src="https://img.shields.io/badge/contributors-1-blueviolet"> <img alt="License" src="https://img.shields.io/badge/license-MIT-yellow">

## Overview

ginvalidator is a set of [Gin](https://github.com/gin-gonic/gin) middlewares that wraps the extensive collection of validators and sanitizers offered by [validatorgo](https://github.com/bube054/validatorgo).

It allows you to combine them in many ways so that you can validate and sanitize your express requests, and offers tools to determine if the request is valid or not, which data was matched according to your validators, and so on.

It is based on the popular js/express library [express-validator](https://github.com/express-validator/express-validator)

## Support
This version of ginvalidator requires that your application is running on [Go](https://go.dev/dl/) 1.21+.
It's also verified to work with [gin](https://github.com/gin-gonic/gin) 1.x.x.

## Rationale
Why not use?

* *Handwritten Validators*:
  You could write your own validation logic manually, but that gets repetitive and messy fast. Every time you need a new validation, you’re writing the same kind of code over and over. It’s easy to make mistakes, and it’s a pain to maintain.
* *Gin's Built-in Model Binding and Validation*:
  Gin has validation built in, but it’s not ideal for everyone. Struct tags are limiting and make your code harder to read, especially when you need complex rules. Plus, the validation gets tied too tightly to your models, which isn't great for flexibility.
* *Other Libraries (like [Galidator](github.com/golodash/galidator))*:
  There are other libraries out there, but they often feel too complex for what they do. They require more setup and work than you’d expect, especially when you just want a simple, straightforward solution for validation.


## Installation

Using go get.

```
 go get github.com/bube054/ginvalidator
```

Then import the package into your own code.

```go
 import (
   "fmt"
   "github.com/bube054/ginvalidator"
 )
```

If you are unhappy using the long validatorgo package name, you can do this.

```go
 import (
   "fmt"
   gv "github.com/bube054/ginvalidator"
 )
```

# Maintainers

- [bube054](https://github.com/bube054) - Attah Gbubemi David (author)

# Other related projects

- [ginvalidator](https://github.com/bube054/ginvalidator)
- [echovalidator](https://github.com/bube054/echovalidator)
- [fibervalidator](https://github.com/bube054/fibervalidator)
- [chivalidator](https://github.com/bube054/chivalidator)

# License

This project is licensed under the [MIT](https://opensource.org/license/mit). See the [LICENSE](https://github.com/bube054/validatorgo/blob/master/LICENSE) file for details.
