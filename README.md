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

## Setup

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

## Adding a validator

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

## Handling validation errors

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

## Sanitizing inputs

While the user can no longer send empty person names, it can still inject HTML into your page! This is known as the Cross-Site Scripting vulnerability (XSS).
Let's see how it works. Go to http://localhost:8080/ping?person=<b>John</b>, and you should see "Hello, <b>John</b>!".
While this example is fine, an attacker could change the person query string to a \<script> tag which loads its own JavaScript that could be harmful.
In this scenario, one way to mitigate the issue with express-validator is to use a sanitizer, most specifically Escape, which transforms special HTML characters with others that can be represented as text.

```go
 func main() {
    r := gin.Default()
    r.GET("/ping",
     gv.NewQuery("person", nil).Chain().Not().Empty(nil).Escape().Validate(),
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

## Accessing validated data

You can use [GetMatchedData](), which automatically collects all data that ginvalidator has validated and/or sanitized:

```go
 func main() {
    r := gin.Default()
    r.GET("/ping",
     gv.NewQuery("person", nil).Chain().Not().Empty(nil).Escape().Validate(),
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

      person := data["queries"]["person"]
      ctx.String(http.StatusOK, "Hello, %s!", person)
    })
    r.Run()
  }
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
