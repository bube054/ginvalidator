# ginvalidator

<img alt="Tag" src="https://img.shields.io/badge/tag-v0.1.0-blue?labelColor=gray"> <img alt="Go Version" src="https://img.shields.io/badge/Go->=1.21-00ADD8?labelColor=gray"> <img alt="Reference" src="https://img.shields.io/badge/-reference-00ADD8?logo=go&labelColor=gray"> <img alt="Tests" src="https://img.shields.io/badge/tests-passing-brightgreen?logo=github&labelColor=gray"> <img alt="Go Report" src="https://img.shields.io/badge/go_report-A%2B-00ADD8"> <img alt="Coverage" src="https://img.shields.io/badge/coverage-87.30%25-brightgreen?logo=codecov"> <img alt="Contributors" src="https://img.shields.io/badge/contributors-1-blueviolet"> <img alt="License" src="https://img.shields.io/badge/license-MIT-yellow">

## Overview

ginvalidator is a set of [Gin](https://github.com/gin-gonic/gin) middlewares that wraps the extensive collection of validators and sanitizers offered by my other open source package [validatorgo](https://github.com/bube054/validatorgo).

It allows you to combine them in many ways so that you can validate and sanitize your express requests, and offers tools to determine if the request is valid or not, which data was matched according to your validators, and so on.

It is based on the popular js/express library [express-validator](https://github.com/express-validator/express-validator)

## Support

This version of ginvalidator requires that your application is running on [Go](https://go.dev/dl/) 1.16+.
It's also verified to work with [gin](https://github.com/gin-gonic/gin) 1.x.x.

## Rationale

Why not use?

- _Handwritten Validators_:
  You could write your own validation logic manually, but that gets repetitive and messy fast. Every time you need a new validation, you’re writing the same kind of code over and over. It’s easy to make mistakes, and it’s a pain to maintain.
- _Gin's Built-in Model Binding and Validation_:
  Gin has validation built in, but it’s not ideal for everyone. Struct tags are limiting and make your code harder to read, especially when you need complex rules. Plus, the validation gets tied too tightly to your models, which isn't great for flexibility.
- _Other Libraries (like [Galidator](github.com/golodash/galidator))_:
  There are other libraries out there, but they often feel too complex for what they do. They require more setup and work than you’d expect, especially when you just want a simple, straightforward solution for validation.

## Installation

Create an empty folder then run.

```
$ go mod init example.com/learning
go: creating new go.mod: module example.com/learning
```

Using go get install gin.

```
$ go get -u github.com/gin-gonic/gin
```

Also using go get install ginvalidator.

```
$ go get -u github.com/bube054/ginvalidator
```

## Getting Started

One of the best ways to learn something is by example! So let's roll the sleeves up and get some coding happening.

### Setup

The first thing that one needs is a gin server running. Let's implement one that pings someone; for this, create a main.go then add the following code:

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.GET("/ping", func(ctx *gin.Context) {
        person := ctx.Query("person")
        ctx.String(http.StatusOK, "Hello, %s!", person)
    })
    r.Run() // listen and serve on 0.0.0.0:8080
}
```

Now run this file by executing go run main.go on your terminal.

```
$ go run main.go
```

The HTTP server should be running, and you can open http://localhost:8080/ping?person=John to salute John!

> 💡 **Tip:**
> You can use [Air](https://blog.logrocket.com/using-air-go-implement-live-reload/) with Go and gin to implement live reload. These automatically restart the server whenever a file is changed, so you don't have to do this yourself!

### Adding a validator

So the server is working, but there are problems with it. Most notably, you don't want to ping someone when the person's name is not set.
For example, going to http://localhost:8080/ping will print "Hello, ".

That's where ginvalidator comes in handy. It provides validators, sanitizers and modifiers that are used to validate your request.
Let's add a validator and a modifier that checks that the person query string cannot be empty, with the validator named Empty and modifier named Not:

```go
 import (
   "net/http"

   "github.com/gin-gonic/gin"
   gv "github.com/bube054/ginvalidator"
 )

 func main() {
    r := gin.Default()
    r.GET("/ping", gv.NewQuery("person", nil).Chain().Not().Empty(nil).Validate(), func(ctx *gin.Context) {
        person := ctx.Query("person")
        ctx.String(http.StatusOK, "Hello, %s!", person)
    })
    r.Run()
  }
```

Now, restart your server, and go to http://localhost:8080/ping again. Hmm, it still prints "Hello, !"... why?

### Handling validation errors

ginvalidator validation chain dos not report validation errors to users automatically.
The reason for this is simple: as you add more validators, or for more fields, how do you want to collect the errors? Do you want a list of all errors, only one per field, only one overall...?

So the next obvious step is to change the above code again, this time verifying the validation result with the [ValidationResult]() function:

```go
 import (
   "net/http"

   "github.com/gin-gonic/gin"
   gv "github.com/bube054/ginvalidator"
 )

 func main() {
    r := gin.Default()
    r.GET("/ping",
      gv.NewQuery("person", nil).Chain().Not().Empty(nil).Validate(),
      func(ctx *gin.Context) {
        result, err := gv.ValidationResult(ctx)
        if err != nil {
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "message": "The server encountered an unexpected error.",
        })
      }

      if len(result) != 0 {
        ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
            "errors": result,
        })
      }

      person := ctx.Query("person")
      ctx.String(http.StatusOK, "Hello, %s!", person)
    })
    r.Run()
  }
```

Now if you access http://localhost:8080/ping again, what you'll see is the following JSON content:

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
- it's located in the query string (location: "query");
- the error message that was given was Invalid value.

This is a better scenario, but it can still be improved. Let's continue.

### Creating better error messages

All request location validators accept an optional second argument, a function used to format error messages. If `nil` is provided, a generic error message will be used instead, as seen above.

```go
 import (
   "net/http"

   "github.com/gin-gonic/gin"
   gv "github.com/bube054/ginvalidator"
 )

 func main() {
    r := gin.Default()
    r.GET("/ping",
      gv.NewQuery("person",
      func(initialValue, sanitizedValue, validatorName string) string {
        return "Please enter your name."
      },
      ).Chain().Not().Empty(nil).Validate(),
      func(ctx *gin.Context) {
        result, err := gv.ValidationResult(ctx)
        if err != nil {
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "message": "The server encountered an unexpected error.",
        })
      }

      if len(result) != 0 {
        ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
            "errors": result,
        })
      }

      person := ctx.Query("person")
      ctx.String(http.StatusOK, "Hello, %s!", person)
    })
    r.Run()
  }
```

Now if you access http://localhost:8080/ping again, what you'll see is the following JSON content:

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

### Sanitizing inputs

While the user can no longer send empty person names, it can still inject HTML into your page! This is known as the Cross-Site Scripting vulnerability (XSS).
Let's see how it works. Go to http://localhost:8080/ping?person=<b>John</b>, and you should see "Hello, <b>John</b>!".
While this example is fine, an attacker could change the person query string to a \<script> tag which loads its own JavaScript that could be harmful.
In this scenario, one way to mitigate the issue with express-validator is to use a sanitizer, most specifically Escape, which transforms special HTML characters with others that can be represented as text.

```go
 func main() {
    r := gin.Default()
    r.GET("/ping",
      gv.NewQuery("person",
      func(initialValue, sanitizedValue, validatorName string) string {
        return "Please enter your name."
      },
      ).Chain().Not().Empty(nil).Validate(),
     func(ctx *gin.Context) {
      result, err := gv.ValidationResult(ctx)
      if err != nil {
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "message": "The server encountered an unexpected error.",
        })
      }

      if len(result) != 0 {
        ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
            "errors": result,
        })
      }

      person := ctx.Query("person")
      ctx.String(http.StatusOK, fmt.Sprintf("Hello, %s!", person))
    })
    r.Run()
  }
```

Now, if you restart the server and refresh the page, what you'll see is "Hello, &lt;b&gt;John&lt;/b&gt;!". Our example page is no longer vulnerable to XSS!

### Accessing validated data

You can use [GetMatchedData](), which automatically collects all data that ginvalidator has validated and/or sanitized:

```go
 func main() {
    r := gin.Default()
    r.GET("/ping",
      gv.NewQuery("person",
      func(initialValue, sanitizedValue, validatorName string) string {
        return "Please enter your name."
      },
      ).Chain().Not().Empty(nil).Validate(),
     func(ctx *gin.Context) {
      result, err := gv.ValidationResult(ctx)
      if err != nil {
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "message": "The server encountered an unexpected error.",
        })
      }

      if len(result) != 0 {
        ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
            "errors": result,
        })
      }

      data, err := gv.GetMatchedData(ctx)
      if err != nil {
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "message": "The server encountered an unexpected error.",
        })
      }

      person, ok := data["queries"]["person"]
      if !ok {
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "message": "The server encountered an unexpected error.",
        })
      }

      ctx.String(http.StatusOK, "Hello, %s!", person)
    })
    r.Run()
  }
```

## The Validation Chain

The validation chain is one of the main concepts in ginvalidator, therefore it's useful to learn about it, so that you can use it effectively.

But don't worry: if you've read through the [Getting Started guide](#getting-started), you have already used validation chains without even noticing!

### What are validation chains?

Validation chains are created using the following functions, each targeting a specific location in the HTTP request:

- `NewBody`: Validates data from the `http.Request` body. Its location is `"body"`.
- `NewCookie`: Validates data from the `http.Request` cookies. Its location is `"cookies"`.
- `NewHeader`: Validates data from the `http.Request` headers. Its location is `"headers"`.
- `NewParam`: Validates data from the Gin route parameters. Its location is `"params"`.
- `NewQuery`: Validates data from the `http.Request` query parameters. Its location is `"queries"`.

They have this name because they wrap the value of a field with validations (or sanitizations), and each of its methods returns itself.
This pattern is usually called [method chaining](https://en.wikipedia.org/wiki/Method_chaining), hence why the name validation chain.

Validation chains not only have a number of useful methods for defining validations and sanitizations, but they also have methods `Validate()` which return gin middleware functions.

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

If the value is invalid, an error is recorded for that field using some error message. This validation error can then be retrieved at a later point in the route handler and returned to the user.

`Sanitizers` transform the field value. They are useful to remove noise from the value and perhaps even to provide some basic line of defense against threats.

Sanitizers persist the updated fields value back into the gin Contexts, so that it's usable by other ginvalidator functions, your own route handler code, and even other middlewares.

`Modifiers` define how validation chains behave when they are run.

### Standard validators/sanitizers

All of the functionality exposed by the validation chain actually comes from [validatorgo](github.com/bube054/validatorgo), one of my other open source go packages which specializes in string validation/sanitation. Please checkout, star and share 🙏🙏🙏, Thank You.

This includes all of validatorgo validators and sanitizers, from commonly used `IsEmail`, `IsLength`, and `Trim` to the more niche `IsISBN`, `IsMultibyte` and `StripLow`!

These are called standard validators and standard sanitizers in ginvalidator. But with the `Is` prefix.

Because validatorgo only works with strings, ginvalidator will always convert fields with a standard validator/sanitizer to string first.

### Chaining order

The order in which you call methods on a validation chain usually matters.
They are almost always run in the order that they are specified, therefore you can tell what a validation chain will do just by reading its definition, from first chained method to last.

Take the following snippet as an example:

```go
// Validate if search_query is not empty, then trim it
NewQuery("search_query", nil).Not().Empty().Trim("");
```

In this case, if the user passes a search_query value that is composed of whitespaces only, it won't be empty, therefore the validation passes. But since the .trim() sanitizer is there, the whitespaces will be removed, and the field will become empty, so you actually end up with a false positive.

Now, compare it with the below snippet:

```go
// Trim search_query, then validate if it's not empty
NewQuery("search_query", nil).Not().Empty().Trim("");
```

This chain will more sensibly remove whitespaces, and then validate if the value is not empty.

### Reusing validation chains

<!-- Validation chains are mutable.
This means that calling methods on one will cause the original chain object to be updated, just like any references to it. -->

If you wish to reuse the same chain, it's a good idea to return them from functions:

<!-- One exception to this rule is .optional(): It can be placed at any point in the chain and it will mark the chain as optional the same way. -->

```go
func createEmailChain() gv.ValidationChain {
  return gv.NewBody("email", nil).Chain().Email(nil)
}

r.POST("/login", createEmailChain().Validate(), handleLoginRoute)
r.POST("/signup", createEmailChain().Validate(), handleSignupRoute)
```

## Field Selection

In ginvalidator, a field is any value that is either validated or sanitized. It is string.

Pretty much every function or value returned by express-validator reference fields in some way. For this reason, it's important to understand the field path syntax both for when selecting fields for validation, and when accessing the validation errors or validated data.

### Syntax

- `Body` fields are only valid for the following Content-Types:
  - application/json: Uses [gjson path syntax](https://github.com/tidwall/gjson?tab=readme-ov-file#path-syntax) for value extraction.
  - application/x-www-form-urlencoded
  - multipart/form-data

- `Query` fields correspond to URL search parameters, and their values are automatically escaped by Gin.
- `Param` fields represent URL path parameters, and their values are automatically escaped by ginvalidator.
- `Header` fields are HTTP request headers, and their values are not escaped. A log warning will appear if you provide a non-canonical header key.
- `Cookies` fields are HTTP cookies, and their values are automatically escaped by Gin.


## Customizing express-validator

If the server you're building is anything but a very simple one, you'll need validators, sanitizers and error messages beyond the ones built into express-validator sooner or later.

### Custom Validators and Sanitizers

A classic need that ginvalidator can't fulfill for you, and that you might run into, is validating whether an e-mail address is in use or not when a user signing up.

It's possible to do this in ginvalidator by implementing a custom validator.

Custom validators are also methods available on the validation chain, that receives a special function [CustomValidatorFunc](), and have to returns a boolean that will determine if the field is valid or not.

Custom sanitizers are also methods available on the validation chain, that receives a special function [CustomSanitizerFunc](), and have to returns the new sanitized value.

## Implementing a custom validator
Custom validators must return a booleans. true to indicate that the field is valid, or false to indicate it's invalid.

Custom validators can be asynchronous by using goroutines and a `sync.WaitGroup` to handle concurrent operations. Within the validator, you can spin up goroutines for each asynchronous task, adding each task to the WaitGroup. Once all tasks complete, the validator should return a boolean.

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
        func(req *http.Request, initialValue, sanitizedValue string) bool {
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
    CustomValidator(func(req *http.Request, initialValue, sanitizedValue string) bool {
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
Custom sanitizers don't have many rules. Whatever the value that they return, is the new value that the field will acquire.
Custom sanitizers can also be asynchronous by using goroutines and a `sync.WaitGroup` to handle concurrent operations.

```go
r.POST("/user/:id", gv.NewParam("id", nil).Chain().CustomSanitizer(
  func(req *http.Request, initialValue, sanitizedValue string) string {
    return strings.Repeat(sanitizedValue, 3) // some string manipulation
  },
).Validate(),
  func(ctx *gin.Context) {
    // Handle request
  },
)
```

### Error Messages
Whenever a field value is invalid, an error message is recorded for it.
The default error message is `"Invalid value"`, which is not descriptive at all of what the error is, so you might need to customize it. You can customize by

```go
gv.NewBody("email", func(initialValue, sanitizedValue, validatorName string) string {

  switch validatorName {
  case gv.EmailValidatorName:
    return "Email is not valid."
  case gv.EmptyValidatorName:
    return "Email is empty."
  default:
    return gv.DefaultValChainErrMsg
  }

}).Chain().Not().Empty(nil).Email(nil).Validate()
```

`initialValue` is the original value extracted from the request (before any sanitization).
`sanitizedValue` is the value after it has been sanitized (if applicable).
`validatorName` is the name of the validator that failed, which helps identify the validation rule that did not pass.

### Maintainers

- [bube054](https://github.com/bube054) - Attah Gbubemi David (author)

# Other related projects

- [ginvalidator](https://github.com/bube054/ginvalidator)
- [echovalidator](https://github.com/bube054/echovalidator)
- [fibervalidator](https://github.com/bube054/fibervalidator)
- [chivalidator](https://github.com/bube054/chivalidator)

# License

This project is licensed under the [MIT](https://opensource.org/license/mit). See the [LICENSE](https://github.com/bube054/validatorgo/blob/master/LICENSE) file for details.
