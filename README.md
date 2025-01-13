# ginvalidator

<!-- <img alt="Tag" src="https://img.shields.io/badge/tag-v0.1.0-blue?labelColor=gray"> <img alt="Go Version" src="https://img.shields.io/badge/Go->=1.21-00ADD8?labelColor=gray"> <img alt="Reference" src="https://img.shields.io/badge/-reference-00ADD8?logo=go&labelColor=gray"> <img alt="Tests" src="https://img.shields.io/badge/tests-passing-brightgreen?logo=github&labelColor=gray"> <img alt="Go Report" src="https://img.shields.io/badge/go_report-A%2B-00ADD8"> <img alt="Coverage" src="https://img.shields.io/badge/coverage-87.30%25-brightgreen?logo=codecov"> <img alt="Contributors" src="https://img.shields.io/badge/contributors-1-blueviolet"> <img alt="License" src="https://img.shields.io/badge/license-MIT-yellow"> -->

## Overview

ginvalidator is a set of [Gin](https://github.com/gin-gonic/gin) middlewares that wraps the extensive collection of validators and sanitizers offered by my other open source package [validatorgo](https://github.com/bube054/validatorgo). It also uses the popular open-source package [gjson](https://github.com/tidwall/gjson) for JSON field syntax, providing efficient querying and extraction of data from JSON objects.

It allows you to combine them in many ways so that you can validate and sanitize your Gin requests, and offers tools to determine if the request is valid or not, which data was matched according to your validators.

It is based on the popular js/express library [express-validator](https://github.com/express-validator/express-validator)

## Support

This version of ginvalidator requires that your application is running on [Go](https://go.dev/dl/) 1.16+.
It's also verified to work with [Gin](https://github.com/gin-gonic/gin) 1.x.x.

## Rationale

Why not use?

- _Handwritten Validators_:
  You could write your own validation logic manually, but that gets repetitive and messy fast. Every time you need a new validation, you’re writing the same kind of code over and over. It’s easy to make mistakes, and it’s a pain to maintain.
- _Gin's Built-in Model Binding and Validation_:
  Gin has validation built in, but it’s not ideal for everyone. Struct tags are limiting and make your code harder to read, especially when you need complex rules. Plus, the validation gets tied too tightly to your models, which isn't great for flexibility.
- _Other Libraries (like [Galidator](github.com/golodash/galidator))_:
  There are other libraries out there, but they often feel too complex for what they do. They require more setup and work than you’d expect, especially when you just want a simple, straightforward solution for validation.

## Installation

Make sure you have [Go](https://go.dev/dl/) installed on your machine.

### Step 1: Create a New Go Module

1. Create an empty folder with a name of your choice.
2. Open a terminal, navigate (`cd`) into that folder, and initialize a new Go module:

```bash
go mod init example.com/tutorial
```

### Step 2: Install Required Packages

Use go get to install the necessary packages.

1. Install Gin:

```bash
go get -u github.com/gin-gonic/gin
```

2. Install ginvalidator:

```bash
go get -u github.com/bube054/ginvalidator
```

> 📝 **Note:**  
`ginvalidator` uses `validatorgo` as a dependency. In this tutorial, we will directly import `validatorgo` into the project. Remember to run `go mod tidy` to add it as a direct dependency.

## Getting Started

One of the best ways to learn something is by example! So let's roll the sleeves up and get some coding happening.

### Setup

The first thing that one needs is a Gin server running. Let's implement one that says hi to someone; for this, create a `main.go` then add the following code:

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.GET("/hello", func(ctx *gin.Context) {
        person := ctx.Query("person")
        ctx.String(http.StatusOK, "Hello, %s!", person)
    })

    r.Run() // listen and serve on 0.0.0.0:8080
}
```

Now run this file by executing `go run main.go` on your terminal.

```bash
go run main.go
```

The HTTP server should be running, and you can open http://localhost:8080/hello?person=John to salute John!

> 💡 **Tip:**
> You can use [Air](https://github.com/air-verse/air) with Go and Gin to implement live reload. This automatically restart the server whenever a file is changed, so you don't have to do this yourself!

### Adding a validator

So the server is working, but there are problems with it. Most notably, you don't want to say hello to someone when the person's name is not set.
For example, going to http://localhost:8080/hello will print `"Hello, "`.

That's where `ginvalidator` and also `validatorgo` come in handy. 
`ginvalidator` provides `validators`, `sanitizers` and `modifiers` that are used to validate your request. 
`validatorgo` provides a set of configuration structs to help you customize `validators` and `sanitizers`. These structs, like `vgo.IsEmptyOpts{IgnoreWhitespace: false}`, allow you to fine-tune the behavior of each validation and or sanitization step. By passing these configuration options into chain methods, you can precisely control how the input is processed and validated.
Let's add a validator and a modifier that checks that the person query string cannot be empty, with the validator named Empty and modifier named Not:

```go
package main

import (
    "net/http"

    gv "github.com/bube054/ginvalidator"
    vgo "github.com/bube054/validatorgo"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.GET("/hello", 
        gv.NewQuery("person", nil).
        Chain().
        Not().
        Empty(&vgo.IsEmptyOpts{IgnoreWhitespace: false}).
        Validate(), func(ctx *gin.Context) {
            person := ctx.Query("person")
            ctx.String(http.StatusOK, "Hello, %s!", person)
        })

    r.Run()
}
```

> 📝 **Note:**  
> For brevity, `gv` is used as an alias for `ginvalidator` and `vgo` is used as an alias for `validatorgo` in the code examples.

Now, restart your server, and go to http://localhost:8080/hello again. Hmm, it still prints `"Hello, !"`... why?

### Handling validation errors

`ginvalidator` validation chain does not report validation errors to users automatically.
The reason for this is simple: as you add more `validators`, or for more `fields`, how do you want to collect the `errors`? Do you want a list of all `errors`, only one per `field`, only one overall...?

So the next obvious step is to change the above code again, this time verifying the validation result with the `ValidationResult` function:

```go
package main

import (
    "net/http"

    gv "github.com/bube054/ginvalidator"
    vgo "github.com/bube054/validatorgo"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.GET("/hello",
        gv.NewQuery("person", nil).
            Chain().
            Not().
            Empty(&vgo.IsEmptyOpts{IgnoreWhitespace: false}).
            Validate(),
        func(ctx *gin.Context) {
            result, err := gv.ValidationResult(ctx)
            if err != nil {
                ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                    "message": "The server encountered an unexpected error.",
                })
                return
            }

            if len(result) != 0 {
                ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
                    "errors": result,
                })
                return
            }

            person := ctx.Query("person")
            ctx.String(http.StatusOK, "Hello, %s!", person)
        })

    r.Run()
}
```

Now, if you access http://localhost:8080/hello again, you’ll see the following JSON content, formatted for clarity:

```
{
  "errors": [
    {
      "location": "queries",
      "message": "Invalid value",
      "field": "person",
      "value": ""
    }
  ]
}
```

Now, what this is telling us is that

- there's been exactly one error in this request;
- this field is called person;
- it's located in the query string (location: `"queries"`);
- the error message that was given was `"Invalid value"`.

This is a better scenario, but it can still be improved. Let's continue.

### Creating better error messages

All request location validators accept an optional second argument, which is a function used to format the error message. If `nil` is provided, a default, generic error message will be used, as shown in the example above.

```go
package main

import (
    "net/http"

    gv "github.com/bube054/ginvalidator"
    vgo "github.com/bube054/validatorgo"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.GET("/hello",
        gv.NewQuery("person",
            func(initialValue, sanitizedValue, validatorName string) string {
                return "Please enter your name."
            },
        ).Chain().
            Not().
            Empty(&vgo.IsEmptyOpts{IgnoreWhitespace: false}).
            Validate(),
        func(ctx *gin.Context) {
            result, err := gv.ValidationResult(ctx)
            if err != nil {
                ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                    "message": "The server encountered an unexpected error.",
                })
                return
            }

            if len(result) != 0 {
                ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
                    "errors": result,
                })
                return
            }

            person := ctx.Query("person")
            ctx.String(http.StatusOK, "Hello, %s!", person)
        })

    r.Run()
}
```

Now if you access http://localhost:8080/hello again, what you'll see is the following JSON content, with the new error message:

```
{
  "errors": [
    {
      "location": "queries",
      "message": "Please enter your name.",
      "field": "person",
      "value": ""
    }
  ]
}
```

### Accessing validated/sanitized data

You can use `GetMatchedData`, which automatically collects all data that `ginvalidator` has validated and/or sanitized. This data can then be accessed using the `Get` method of `MatchedData`:

```go
package main

import (
    "fmt"
    "net/http"

    gv "github.com/bube054/ginvalidator"
    vgo "github.com/bube054/validatorgo"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.GET(
        "/hello",
        gv.NewQuery(
            "person",
            func(initialValue, sanitizedValue, validatorName string) string {
                return "Please enter your name."
            },
        ).Chain().
            Not().
            Empty(&vgo.IsEmptyOpts{IgnoreWhitespace: false}).
            Validate(),
        func(ctx *gin.Context) {
            result, err := gv.ValidationResult(ctx)
            if err != nil {
                ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                    "message": "The server encountered an unexpected error.",
                })
                return
            }

            if len(result) != 0 {
                ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
                    "errors": result,
                })
                return
            }

            data, err := gv.GetMatchedData(ctx)
            if err != nil {
                ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                    "message": "The server encountered an unexpected error.",
                })
                return
            }

            person, ok := data.Get(gv.QueryLocation, "person")
            if !ok {
                ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
                    "message": fmt.Sprintf(
                        "The server could not find 'person' in the expected location: %s. Also please ensure you're using the correct location, such as Body, Header, Cookie, Query, or Param.",
                        gv.QueryLocation,
                    ),
                })
                return
            }

            ctx.String(http.StatusOK, "Hello, %s!", person)
        },
    )

    r.Run()
}
```

open http://localhost:8080/hello?person=John to salute John!

### Available Data Locations 🚩

The following are the valid data locations you can use:  
- **`BodyLocation`**: Represents the request body.  
- **`CookieLocation`**: Represents cookies in the request.  
- **`QueryLocation`**: Represents query parameters in the URL.  
- **`ParamLocation`**: Represents path parameters in the request.  
- **`HeaderLocation`**: Represents the headers in the request.  

Each of these locations includes a `String` method that returns the location where validated/sanitized data is stored.


### Sanitizing inputs

While the user can no longer send empty person names, it can still inject HTML into your page! This is known as the [Cross-Site Scripting vulnerability (XSS)](https://www.youtube.com/watch?si=zN2bDf-xT-h5wL4a&v=z4LhLJnmoZ0&feature=youtu.be).
Let's see how it works. Go to `http://localhost:8080/hello?person=<b>John</b>`, and you should see "Hello, <b>John</b>!".
While this example is fine, an attacker could change the person query string to a \<script> tag which loads its own JavaScript that could be harmful.
In this scenario, one way to mitigate the issue with ginvalidator is to use a sanitizer, most specifically `Escape`, which transforms special HTML characters with others that can be represented as text.

```go
package main

import (
    "fmt"
    "net/http"

    gv "github.com/bube054/ginvalidator"
    vgo "github.com/bube054/validatorgo"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.GET("/hello",
        gv.NewQuery("person",
            func(initialValue, sanitizedValue, validatorName string) string {
                return "Please enter your name."
            },
        ).Chain().
            Not().
            Empty(&vgo.IsEmptyOpts{IgnoreWhitespace: false}).
            Escape(). // Added sanitizer
            Validate(),
        func(ctx *gin.Context) {
            result, err := gv.ValidationResult(ctx)
            if err != nil {
                ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                    "message": "The server encountered an unexpected error.",
                })
                return
            }

            if len(result) != 0 {
                ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
                    "errors": result,
                })
                return
            }

            data, err := gv.GetMatchedData(ctx)
            if err != nil {
                ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                    "message": "The server encountered an unexpected error.",
                })
                return
            }

            person, ok := data.Get(gv.QueryLocation, "person")
            if !ok {
                ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
                    "message": fmt.Sprintf(
                        "The server could not find 'person' in the expected location: %s. Also please ensure you're using the correct location, such as Body, Header, Cookie, Query, or Param.",
                        gv.QueryLocation,
                    ),
                })
                return
            }

            ctx.String(http.StatusOK, "Hello, %s!", person)
        })

    r.Run()
}
```

Now, if you restart the server and refresh the page, what you'll see is "Hello, \&lt;b\&gt;John\&lt;b\&gt;!". Our example page is no longer vulnerable to XSS!

> ⚠️ **Caution:**  
> `ginvalidator` does not modify `http.Request` values during sanitization. To access sanitized data, always use the `GetMatchedData` function.

## The Validation Chain

The [validation chain](https://pkg.go.dev/github.com/bube054/ginvalidator#ValidationChain) is one of the main concepts in ginvalidator, therefore it's useful to learn about it, so that you can use it effectively.

But don't worry: if you've read through the [Getting Started guide](#getting-started), you have already used validation chains without even noticing!

### What are validation chains?

Validation chains are created using the following functions, each targeting a specific location in the HTTP request:

- `NewBody`: Validates data from the `http.Request` body. Its location is `BodyLocation`.
- `NewCookie`: Validates data from the `http.Request` cookies. Its location is `CookieLocation`.
- `NewHeader`: Validates data from the `http.Request` headers. Its location is `HeaderLocation`.
- `NewParam`: Validates data from the Gin route parameters. Its location is `ParamLocation`.
- `NewQuery`: Validates data from the `http.Request` query parameters. Its location is `QueryLocation`.

They have this name because they wrap the value of a field with validations (or sanitizations), and each of its methods returns itself.
This pattern is usually called [method chaining](https://en.wikipedia.org/wiki/Method_chaining), hence why the name validation chain.

Validation chains not only have a number of useful methods for defining validations, sanitizations and modifications but they also have methods `Validate` which returns the `Gin` middleware handler function.

This is an example of how validation chains are usually used, and how you can read one:

```go
  r.Get(
    "newsletter",
    // For the `email` field in ctx.GetRawData()...
    gv.NewBody("email", nil)
    // the actual validation chain
    .Chain()
    // ...mark the field as optional
    .Optional()
    // ...and when it's present, trim its value, then validate it as an email address
    .Trim("")
    .Email(nil),
    maybeSubscribeToNewsletter,
  )
```

### Features

A validation chain has three kinds of methods: `validators`, `sanitizers` and `modifiers`.

`Validators` determine if the value of a request field is valid. This means checking if the field is in the format that you expect it to be. For example, if you're building a sign up form, your requirements could be that the username must be an e-mail address, and that passwords must be at least 8 characters long.

If the value is invalid, an error is recorded for that field using some error message. This validation error can then be retrieved at a later point in the Gin route handler and returned to the user.

They are:

* [CustomValidator](https://pkg.go.dev/github.com/bube054/ginvalidator#CustomValidator)
* [Contains](https://pkg.go.dev/github.com/bube054/ginvalidator#Contains)
* [Equals](https://pkg.go.dev/github.com/bube054/ginvalidator#Equals)
* [AbaRouting](https://pkg.go.dev/github.com/bube054/ginvalidator#AbaRouting)
* [After](https://pkg.go.dev/github.com/bube054/ginvalidator#After)
* [Alpha](https://pkg.go.dev/github.com/bube054/ginvalidator#Alpha)
* [Alphanumeric](https://pkg.go.dev/github.com/bube054/ginvalidator#Alphanumeric)
* [Ascii](https://pkg.go.dev/github.com/bube054/ginvalidator#Ascii)
* [Base32](https://pkg.go.dev/github.com/bube054/ginvalidator#Base32)
* [Base58](https://pkg.go.dev/github.com/bube054/ginvalidator#Base58)
* [Base64](https://pkg.go.dev/github.com/bube054/ginvalidator#Base64)
* [Before](https://pkg.go.dev/github.com/bube054/ginvalidator#Before)
* [Bic](https://pkg.go.dev/github.com/bube054/ginvalidator#Bic)
* [Boolean](https://pkg.go.dev/github.com/bube054/ginvalidator#Boolean)
* [BTCAddress](https://pkg.go.dev/github.com/bube054/ginvalidator#BTCAddress)
* [ByteLength](https://pkg.go.dev/github.com/bube054/ginvalidator#ByteLength)
* [CreditCard](https://pkg.go.dev/github.com/bube054/ginvalidator#CreditCard)
* [Currency](https://pkg.go.dev/github.com/bube054/ginvalidator#Currency)
* [DataURI](https://pkg.go.dev/github.com/bube054/ginvalidator#DataURI)
* [Date](https://pkg.go.dev/github.com/bube054/ginvalidator#Date)
* [Decimal](https://pkg.go.dev/github.com/bube054/ginvalidator#Decimal)
* [DivisibleBy](https://pkg.go.dev/github.com/bube054/ginvalidator#DivisibleBy)
* [EAN](https://pkg.go.dev/github.com/bube054/ginvalidator#EAN)
* [Email](https://pkg.go.dev/github.com/bube054/ginvalidator#Email)
* [Empty](https://pkg.go.dev/github.com/bube054/ginvalidator#Empty)
* [EthereumAddress](https://pkg.go.dev/github.com/bube054/ginvalidator#EthereumAddress)
* [Float](https://pkg.go.dev/github.com/bube054/ginvalidator#Float)
* [FQDN](https://pkg.go.dev/github.com/bube054/ginvalidator#FQDN)
* [FreightContainerID](https://pkg.go.dev/github.com/bube054/ginvalidator#FreightContainerID)
* [FullWidth](https://pkg.go.dev/github.com/bube054/ginvalidator#FullWidth)
* [HalfWidth](https://pkg.go.dev/github.com/bube054/ginvalidator#HalfWidth)
* [Hash](https://pkg.go.dev/github.com/bube054/ginvalidator#Hash)
* [Hexadecimal](https://pkg.go.dev/github.com/bube054/ginvalidator#Hexadecimal)
* [HexColor](https://pkg.go.dev/github.com/bube054/ginvalidator#HexColor)
* [HSL](https://pkg.go.dev/github.com/bube054/ginvalidator#HSL)
* [IBAN](https://pkg.go.dev/github.com/bube054/ginvalidator#IBAN)
* [IdentityCard](https://pkg.go.dev/github.com/bube054/ginvalidator#IdentityCard)
* [IMEI](https://pkg.go.dev/github.com/bube054/ginvalidator#IMEI)
* [In](https://pkg.go.dev/github.com/bube054/ginvalidator#In)
* [Int](https://pkg.go.dev/github.com/bube054/ginvalidator#Int)
* [IP](https://pkg.go.dev/github.com/bube054/ginvalidator#IP)
* [IPRange](https://pkg.go.dev/github.com/bube054/ginvalidator#IPRange)
* [ISIN](https://pkg.go.dev/github.com/bube054/ginvalidator#ISIN)
* [ISO4217](https://pkg.go.dev/github.com/bube054/ginvalidator#ISO4217)
* [ISO6346](https://pkg.go.dev/github.com/bube054/ginvalidator#ISO6346)
* [ISO6391](https://pkg.go.dev/github.com/bube054/ginvalidator#ISO6391)
* [ISO8601](https://pkg.go.dev/github.com/bube054/ginvalidator#ISO8601)
* [ISO31661Alpha2](https://pkg.go.dev/github.com/bube054/ginvalidator#ISO31661Alpha2)
* [ISO31661Alpha3](https://pkg.go.dev/github.com/bube054/ginvalidator#ISO31661Alpha3)
* [ISO31661Numeric](https://pkg.go.dev/github.com/bube054/ginvalidator#ISO31661Numeric)
* [ISRC](https://pkg.go.dev/github.com/bube054/ginvalidator#ISRC)
* [ISSN](https://pkg.go.dev/github.com/bube054/ginvalidator#ISSN)
* [JSON](https://pkg.go.dev/github.com/bube054/ginvalidator#JSON)
* [LatLong](https://pkg.go.dev/github.com/bube054/ginvalidator#LatLong)
* [LicensePlate](https://pkg.go.dev/github.com/bube054/ginvalidator#LicensePlate)
* [Locale](https://pkg.go.dev/github.com/bube054/ginvalidator#Locale)
* [LowerCase](https://pkg.go.dev/github.com/bube054/ginvalidator#LowerCase)
* [LuhnNumber](https://pkg.go.dev/github.com/bube054/ginvalidator#LuhnNumber)
* [MacAddress](https://pkg.go.dev/github.com/bube054/ginvalidator#MacAddress)
* [MagnetURI](https://pkg.go.dev/github.com/bube054/ginvalidator#MagnetURI)
* [MailtoURI](https://pkg.go.dev/github.com/bube054/ginvalidator#MailtoURI)
* [MD5](https://pkg.go.dev/github.com/bube054/ginvalidator#MD5)
* [MimeType](https://pkg.go.dev/github.com/bube054/ginvalidator#MimeType)
* [MobilePhone](https://pkg.go.dev/github.com/bube054/ginvalidator#MobilePhone)
* [MongoID](https://pkg.go.dev/github.com/bube054/ginvalidator#MongoID)
* [Multibyte](https://pkg.go.dev/github.com/bube054/ginvalidator#Multibyte)
* [Numeric](https://pkg.go.dev/github.com/bube054/ginvalidator#Numeric)
* [Octal](https://pkg.go.dev/github.com/bube054/ginvalidator#Octal)
* [PassportNumber](https://pkg.go.dev/github.com/bube054/ginvalidator#PassportNumber)
* [Port](https://pkg.go.dev/github.com/bube054/ginvalidator#Port)
* [PostalCode](https://pkg.go.dev/github.com/bube054/ginvalidator#PostalCode)
* [RFC3339](https://pkg.go.dev/github.com/bube054/ginvalidator#RFC3339)
* [RgbColor](https://pkg.go.dev/github.com/bube054/ginvalidator#RgbColor)
* [SemVer](https://pkg.go.dev/github.com/bube054/ginvalidator#SemVer)
* [Slug](https://pkg.go.dev/github.com/bube054/ginvalidator#Slug)
* [StrongPassword](https://pkg.go.dev/github.com/bube054/ginvalidator#StrongPassword)
* [TaxID](https://pkg.go.dev/github.com/bube054/ginvalidator#TaxID)
* [SurrogatePair](https://pkg.go.dev/github.com/bube054/ginvalidator#SurrogatePair)
* [Time](https://pkg.go.dev/github.com/bube054/ginvalidator#Time)
* [ULID](https://pkg.go.dev/github.com/bube054/ginvalidator#ULID)
* [UpperCase](https://pkg.go.dev/github.com/bube054/ginvalidator#UpperCase)
* [URL](https://pkg.go.dev/github.com/bube054/ginvalidator#URL)
* [UUID](https://pkg.go.dev/github.com/bube054/ginvalidator#UUID)
* [VariableWidth](https://pkg.go.dev/github.com/bube054/ginvalidator#VariableWidth)
* [VAT](https://pkg.go.dev/github.com/bube054/ginvalidator#VAT)
* [Whitelisted](https://pkg.go.dev/github.com/bube054/ginvalidator#Whitelisted)
* [Matches](https://pkg.go.dev/github.com/bube054/ginvalidator#Matches)

`Sanitizers` transform the field value. They are useful to remove noise from the value and perhaps even to provide some basic line of defense against threats.

Sanitizers persist the updated fields value back into the Gin Contexts, so that it's usable by other ginvalidator functions, your own route handler code, and even other middlewares.

They are:
* [CustomSanitizer](https://pkg.go.dev/github.com/bube054/ginvalidator#CustomSanitizer)
* [Blacklist](https://pkg.go.dev/github.com/bube054/ginvalidator#Blacklist)
* [Escape](https://pkg.go.dev/github.com/bube054/ginvalidator#Escape)
* [LTrim](https://pkg.go.dev/github.com/bube054/ginvalidator#LTrim)
* [NormalizeEmail](https://pkg.go.dev/github.com/bube054/ginvalidator#NormalizeEmail)
* [RTrim](https://pkg.go.dev/github.com/bube054/ginvalidator#RTrim)
* [StripLow](https://pkg.go.dev/github.com/bube054/ginvalidator#StripLow)
* [ToBoolean](https://pkg.go.dev/github.com/bube054/ginvalidator#ToBoolean)
* [ToDate](https://pkg.go.dev/github.com/bube054/ginvalidator#ToDate)
* [ToFloat](https://pkg.go.dev/github.com/bube054/ginvalidator#ToFloat)
* [ToInt](https://pkg.go.dev/github.com/bube054/ginvalidator#ToInt)
* [Trim](https://pkg.go.dev/github.com/bube054/ginvalidator#Trim)
* [Unescape](https://pkg.go.dev/github.com/bube054/ginvalidator#Unescape)
* [Whitelist](https://pkg.go.dev/github.com/bube054/ginvalidator#Whitelist)

`Modifiers` define how validation chains behave when they are run.

They are:
* [Bail](https://pkg.go.dev/github.com/bube054/ginvalidator#Bail)
* [If](https://pkg.go.dev/github.com/bube054/ginvalidator#If)
* [Not](https://pkg.go.dev/github.com/bube054/ginvalidator#Not)
* [Skip](https://pkg.go.dev/github.com/bube054/ginvalidator#Skip)
* [Optional](https://pkg.go.dev/github.com/bube054/ginvalidator#Optional)

> 📝 **Note:**  
> These methods are thoroughly documented using GoDoc within the pkg.go.dev ginvalidator [documentation](https://pkg.go.dev/github.com/bube054). If any details are unclear, you may also want to refer to related functions within the `validatorgo` package for additional context, which I’ll be explaining below.

### Standard validators/sanitizers
All of the functionality exposed by the validation chain actually comes from [validatorgo](github.com/bube054/validatorgo), one of my other open source go packages which specializes in string validation/sanitation. Please check it out, star and share 🙏🙏🙏, Thank You.

This includes all of `validatorgo` validators and sanitizers, from commonly used `IsEmail`, `IsLength`, and `Trim` to the more niche `IsISBN`, `IsMultibyte` and `StripLow`!

These are called standard `validators` and standard `sanitizers` in ginvalidator. But without the `Is` prefix from `validatorgo`.

### Chaining order

The order in which you call methods on a validation chain usually matters.
They are almost always run in the order that they are specified, therefore you can tell what a validation chain will do just by reading its definition, from first chained method to last.

Take the following snippet as an example:

```go
// Validate if search_query is not empty, then trim it
NewQuery("search_query", nil).Chain().Not().Empty().Trim("").Validate();
```

In this case, if the user passes a `"search_query"` value that is composed of whitespaces only, it won't be empty, therefore the validation passes. But since the `.Trim()` sanitizer is there, the whitespace's will be removed, and the field will become empty, so you actually end up with a false positive.

Now, compare it with the below snippet:

```go
// Trim search_query, then validate if it's not empty
NewQuery("search_query", nil).Chain().Trim("").Not().Empty().Validate();
```

This chain will more sensibly remove whitespace's, and then validate if the value is not empty.

One exception to this rule is `.Optional()`: It can be placed at any point in the chain and it will mark the chain as optional.

### Reusing validation chains

If you wish to reuse the same chain, it's a good idea to return them from functions:

```go
func createEmailValidator() gin.HandlerFunc {
  return gv.NewBody("email", nil).Chain().Email(nil).Validate()
}

func handleLoginRoute(ctx *gin.Context) {
  // Handle login route
}

func handleSignupRoute(ctx *gin.Context) {
  // Handle signup route
}

r.POST("/login", createEmailValidator(), handleLoginRoute)
r.POST("/signup", createEmailValidator(), handleSignupRoute)
```

## Field Selection

In `ginvalidator`, a field is any value that is either validated or sanitized and it is string.

Pretty much every function or value returned by ginvalidator reference fields in some way. For this reason, it's important to understand the field path syntax both for when selecting fields for validation, and when accessing the validation errors or validated data.

### Syntax

- **`Body` fields** are only valid for the following Content-Types:

   `application/json`: This uses [GJSON path syntax](https://github.com/tidwall/gjson#path-syntax) for extracting values. Please refer to the linked documentation for details.


    - **Example**:
      ```json
      {
        "user": {
          "name": "John",
          "email": "john.doe@example.com"
        }
      }
      ```
      With path `user.name`, the extracted value would be `"John"`.

   `application/x-www-form-urlencoded`: Typically used for HTML form submissions. Fields are submitted as key-value pairs in the body.

    - **Example**:
      ```
      Content-Type: application/x-www-form-urlencoded
      ```
      Body:
      ```
      name=John&email=john.doe@example.com
      ```
      Field `"name"` would have the value `"John"`, and `"email"` would have the value `"john.doe@example.com"`.

   `multipart/form-data`: Commonly used for file uploads or when submitting form data with files.

    - **Example**:

      ```
      Content-Type: multipart/form-data
      ```

      Body:

      ```
      --boundary
      Content-Disposition: form-data; name="name"

      John
      --boundary
      Content-Disposition: form-data; name="file"; filename="resume.pdf"
      Content-Type: application/pdf

      [binary data]
      --boundary--
      ```

      Field `"name"` would have the value `"John"`, and `"file"` would be the uploaded file.

- **`Query` fields** correspond to URL search parameters, and their values are automatically url unescaped by Gin.  
  **Examples:**

  - Field: `"name"`, Value: `"John"`
    ```
    /hello?name=John
    ```
  - Field: `"full_name"`, Value: `"John Doe"`
    ```
    /hello?full_name=John%20Doe
    ```

- **`Param` fields** represent URL path parameters, and their values are automatically unescaped by `ginvalidator`.  
  **Example:**

  - Field: `"id"`, Value: `"123"`
    ```
    /users/:id
    ```

- **`Header` fields** are HTTP request headers, and their values are not unescaped. A log warning will appear if you provide a non-canonical header key.  
  **Example:**

  - Field: `"User-Agent"`, Value: `"Mozilla/5.0"`
    ```
    Header: "User-Agent", Value: "Mozilla/5.0"
    ```

- **`Cookies` fields** are HTTP cookies, and their values are automatically url unescaped by Gin.  
  **Example:**
  - Field: `"session_id"`, Value: `"abc 123"`
    ```
    Cookie: "session_id=abc%20123"
    ```

## Customizing express-validator

If the server you're building is anything but a very simple one, you'll need `validators`, `sanitizers` and error messages beyond the ones built into ginvalidator sooner or later.

### Custom Validators and Sanitizers

A classic need that ginvalidator can't fulfill for you, and that you might run into, is validating whether an e-mail address is in use or not when a user signing up.

It's possible to do this in ginvalidator by implementing a custom validator.

A `CustomValidator` is a method available on the validation chain, that receives a special function [CustomValidatorFunc](), and have to returns a boolean that will determine if the field is valid or not.

A `CustomSanitizer` is also a method available on the validation chain, that receives a special function [CustomSanitizerFunc](), and have to returns the new sanitized value.

## Implementing a custom validator

A `CustomValidator` can be asynchronous by using goroutines and a `sync.WaitGroup` to handle concurrent operations. Within the validator, you can spin up goroutines for each asynchronous task, adding each task to the WaitGroup. Once all tasks complete, the validator should return a boolean.

For example, in order to check that an e-mail is not in use:

```go
func isUserPresent(email string) bool {
    return email == "existing@example.com"
}

r.POST("/create-user",
    gv.
      NewBody("email", nil).
      Chain().
      CustomValidator(
        func(r *http.Request, initialValue, sanitizedValue string) bool {
          var exists bool
          var wg sync.WaitGroup
          wg.Add(1)

          go func() {
            defer wg.Done()
            exists = isUserPresent(sanitizedValue)
          }()

          wg.Wait()

          return !exists
        },
      ).
      Validate(),

    func(ctx *gin.Context) {
      // Handle the request
    },
)
```

Or maybe you could also verify that the password matches the repeat:

```go
type createUser struct {
  Password             string `json:"password"`
  PasswordConfirmation string `json:"passwordConfirmation"`
}

r.POST("/create-user",
    gv.NewBody("password", nil).
    Chain().
    Matches(regexp.MustCompile(`^[A-Za-z\d]{8,}$`)).
    Validate(),
  gv.NewBody("passwordConfirmation", nil).
    Chain().
    CustomValidator(func(r *http.Request, initialValue, sanitizedValue string) bool {
      data, err := io.ReadAll(req.Body)
      if err != nil {
      return false
      }

      // Refill the request body to allow further reads, if needed.
      req.Body = io.NopCloser(bytes.NewBuffer(data))

      var user createUser
      json.Unmarshal(data, &user)

      return sanitizedValue == user.PasswordConfirmation
    }).
    Validate(),
  func(ctx *gin.Context) {
    // Handle request
  },
)
```

> ⚠️ **Caution:**
> If the request body will be accessed multiple times—whether in the same validation chain, in another validation chain for the same request context, or in subsequent handlers—ensure you reset the request body after each read. Failing to do so can lead to errors or missing data when the body is read again.

### Implementing a custom sanitizer

`CustomSanitizer` don't have many rules. Whatever the value that they return, is the new value that the field will acquire.
Custom sanitizers can also be asynchronous by using goroutines and a `sync.WaitGroup` to handle concurrent operations.

```go
r.POST("/user/:id",
    gv.NewParam("id", nil).
      Chain().
      CustomSanitizer(
        func(r *http.Request, initialValue, sanitizedValue string) string {
          return strings.Repeat(sanitizedValue, 3) // some string manipulation
        },
      ).
      Validate(),

    func(ctx *gin.Context) {
      // Handle request
    },
)
```

### Error Messages

Whenever a field value is invalid, an error message is recorded for it.
The default error message is `"Invalid value"`, which is not descriptive at all of what the error is, so you might need to customize it. You can customize by

```go
gv.NewBody("email",
    func(initialValue, sanitizedValue, validatorName string) string {
        switch validatorName {
        case gv.EmailValidatorName:
            return "Email is not valid."
        case gv.EmptyValidatorName:
            return "Email is empty."
        default:
            return gv.DefaultValChainErrMsg
        }
    },
).
Chain().
Not().Empty(nil).
Email(nil).
Validate()
```

- `initialValue` is the original value extracted from the request (before any sanitization).
- `sanitizedValue` is the value after it has been sanitized (if applicable).
- `validatorName` is the name of the validator that failed, which helps identify the validation rule that did not pass.

For a complete list of validator names, refer to the [ginvalidator constants](https://pkg.go.dev/github.com/bube054/ginvalidator#pkg-constants).


### Maintainers

- [bube054](https://github.com/bube054) - **Attah Gbubemi David (author)**

<!-- # Other related projects

- [ginvalidator](https://github.com/bube054/ginvalidator)
- [echovalidator](https://github.com/bube054/echovalidator)
- [fibervalidator](https://github.com/bube054/fibervalidator)
- [chivalidator](https://github.com/bube054/chivalidator) -->

## License

This project is licensed under the [MIT](https://opensource.org/license/mit). See the [LICENSE](https://github.com/bube054/validatorgo/blob/master/LICENSE) file for details.
