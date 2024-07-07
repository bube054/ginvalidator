package ginvalidator

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	valid "github.com/asaskevich/govalidator"
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
)

func ExampleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(ExampleMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	paramErrMsg := "person is not valid."
	p := NewParam("person", paramErrMsg)

	cookieErrMsg := "cookie is not valid."
	c := NewCookie("myCookie", cookieErrMsg)

	bodyErrMsg := "company is not valid."
	b := NewBody("company", bodyErrMsg)

	headerErrMsg := "Auth not valid"
	h := NewHeader("auth", headerErrMsg)

	queryErrMsg := "id is not valid"
	q := NewQuery("id", queryErrMsg)

	router.GET("/hello/:person",
		b.Chain().
			// IsAlpha("person is not alpha").
			// Not().
			// Custom(func(value string, req http.Request, location string) error {
			// 	return nil
			// }).
			IsObject("company is not valid.").
			// IsArray(bodyErrMsg, &ArrayLengthCheckerOpts{}).
			// IsNotEmpty(paramErrMsg).
			// Contains(err, "lighter").
			// IsASCII("person is not ascii").
			// Bail().
			// Not().
			// IsAlphanumeric("").
			Validate(),
		p.Chain().
			// IsAlpha("person is not alpha").
			// Not().
			// Custom(func(value string, req http.Request, location string) error {
			// 	return nil
			// }).
			// IsArray("person is an array", &ArrayLengthCheckerOpts{}).
			// IsNotEmpty(paramErrMsg).
			// Contains(err, "lighter").
			Not().
			IsASCII("person is not ascii").
			// Bail().
			// IsAlphanumeric("").
			Validate(),
		c.Chain().
			IsBOOLEAN(cookieErrMsg, true).
			// IsAlpha(cookieErrMsg).
			// Not().
			// Custom(func(value string, req http.Request, location string) error {
			// 	return nil
			// }).
			// IsArray("person is an array", &ArrayLengthCheckerOpts{}).
			// IsNotEmpty(paramErrMsg).
			// Contains(err, "lighter").
			// IsASCII("person is not ascii").
			// Bail().
			// Not().
			// IsAlphanumeric("").
			Validate(),
		q.Chain().
			// IsBOOLEAN(cookieErrMsg, true).
			// IsAlpha(cookieErrMsg).
			// Not().
			// Custom(func(value string, req http.Request, location string) error {
			// 	return nil
			// }).
			// IsArray("person is an array", &ArrayLengthCheckerOpts{}).
			// IsNotEmpty(paramErrMsg).
			Contains(queryErrMsg, "619").
			// IsASCII("person is not ascii").
			// Bail().
			// Not().
			// IsAlphanumeric("").
			Validate(),
		h.Chain().
			// IsBOOLEAN(headerErrMsg, true).
			// IsAlpha(cookieErrMsg).
			// Not().
			Custom(func(value string, req http.Request, location string) error {
				if value != "1234567891" {
					return errors.New("value is not 1234567891.")
				}

				return nil
			}).
			// IsArray("person is an array", &ArrayLengthCheckerOpts{}).
			// IsNotEmpty(paramErrMsg).
			// Contains(err, "lighter").
			// IsASCII("person is not ascii").
			// Bail().
			// Not().
			// IsAlphanumeric("").
			Validate(),

		func(ctx *gin.Context) {
			person := ctx.Query("person")

			errs := NewValidationResult(ctx)

			log.Printf("%+v\n", errs)

			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"message": person})
		})
	return router
}

func TestExampleMiddleware(t *testing.T) {
	router := setupRouter()

	// Test the /test route
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	fmt.Println("isValid", valid.IsASCII("你好，世界！"))

	// fmt.Printf("json field(s) %v", splitFieldOnPeriod("name"))

	// 	assert.Equal(t, http.StatusOK, w.Code)
	// 	assert.Equal(t, `{"message":"success"}`, w.Body.String())

	data := []byte(`{
	  "person": {
	    "name": {
	      "first": "Leonid",
	      "last": "Bugaev",
	      "fullName": "Leonid Bugaev",
	    },
	    "github": {
	      "handle": "buger",
	      "followers": 109
	    },
	    "avatars": [
	      { "url": "https://avatars1.githubusercontent.com/u/14009?v=3&s=460", "type": "thumbnail" }
	    ]
	  },
	  "company": {
	    "name": "Acme"
	  }
	}`)

	_ = data

	ty := "person"
	key, typ, _, _ := jsonparser.Get([]byte(fmt.Sprintf(`{"key":"%s"}`, ty)), "key")
	// key, typ, _, _ := jsonparser.Get(data, "person", "avatars", "[0]", "url")

	fmt.Printf("key is %s while datatype is %v", key, typ)
}

func TestParamMiddleware(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	body := `{
  "person": {
    "name": {
      "first": "Leonid",
      "last": "Bugaev",
      "fullName": "Leonid Bugaev"
    },
    "github": {
      "handle": "buger",
      "followers": 109
    },
    "avatars": [
      { "url": "https://avatars1.githubusercontent.com/u/14009?v=3&s=460", "type": "thumbnail" }
    ]
  },
  "company": {
    "name": "Acme"
  }
}`
	req, _ := http.NewRequest("GET", "/hello/light-speed?id=619", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Auth", "123456789")

	// Setting a cookie
	cookie := &http.Cookie{
		Name:  "myCookie",
		Value: "cookieValue",
		Path:  "/",
	}
	req.AddCookie(cookie)

	router.ServeHTTP(w, req)
}
