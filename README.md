# ginvalidator (Pishon)

Well i'm gonna say it, request validation sucks with the [gin](https://gin-gonic.com/) framework. [Binding structs with tags and bulky code 🤮,](https://gin-gonic.com/docs/examples/custom-validators/) thats why i built ginvalidator.
Its almost an exact replica of the popular nodejs and express framework [express-validator](https://github.com/express-validator/express-validator).

ginvalidator is a set of gin middlewares that wraps the extensive collection of validators and sanitizers offered by [govalidator](https://github.com/asaskevich/govalidator) which was derived from the js package [validator.js](https://github.com/validatorjs/validator.js).
It allows you to combine them in many ways so that you can validate and sanitize your gin requests, and offers tools to determine if the request is valid or not, which data was matched according to your validators, and so on.

## It validates request data from these sources/locations.

- The request body: the body of the HTTP request. Can be any value, however objects, arrays and other JavaScript primitives work better.req.body: the body of the HTTP request. Can be any value, however objects, arrays and other JavaScript primitives work better.
- The request cookies: the Cookie header parsed as an object from cookie name to its value.
- The request headers: the headers sent along with the HTTP request.
- The request params: an object from name to value. In express.js, this is parsed from the request path and matched with route definition path, but it can really be anything meaningful coming from the HTTP request.
- The request query's: the portion after the ? in the HTTP request's path, parsed as an object from query parameter name to value.

## Installation

Make sure that [Go is installed on your computer](https://go.dev/doc/install).
Install [gin](https://gin-gonic.com/docs/quickstart/), then type the following command in your terminal:

```go
go get github.com/bube054/ginvalidator
```

Add following line in your \*.go file(s):

```go
import "github.com/bube054/ginvalidator"
```

If you are unhappy using the long ginvalidator, you can do something like this:

```go
import (
  gv "github.com/bube054/ginvalidator"
)
```

## General use code snippet.

```go
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
  gv "github.com/bube054/ginvalidator"
)

func main() {
  router := gin.Default()

  p := gv.NewParam("person", "person can not be empty.")
  router.GET(
    "/hello/:person",
    p.Chain().Not().Empty("person can not be empty.").Validate(),
    func(ctx *gin.Context) {
			person := ctx.Query("person")
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"message": person})
		}
  )

  router.Run() // listen and serve on 0.0.0.0:8080
}
```
