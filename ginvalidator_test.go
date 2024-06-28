package ginvalidator

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	valid "github.com/asaskevich/govalidator"
	// "github.com/buger/jsonparser"
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

	router.GET("/hello", func(ctx *gin.Context) {
		p := NewParam(ctx, "person", "person not provided")
		p.IsAlphanumeric("person is not alphanumeric.").IsASCII("person is not an ascii character.").IsALPHA("person is not an alphanumeric.")

		fmt.Printf("Validation errors: %+v\n", p.GetErrors())

		ctx.Next()
	}, func(c *gin.Context) {
		person := c.Query("person")
		c.JSON(http.StatusOK, gin.H{"message": person})
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

	// 	assert.Equal(t, http.StatusOK, w.Code)
	// 	assert.Equal(t, `{"message":"success"}`, w.Body.String())

	// 	data := []byte(`{
	//   "person": {
	//     "name": {
	//       "first": "Leonid",
	//       "last": "Bugaev",
	//       "fullName": "Leonid Bugaev",
	//     },
	//     "github": {
	//       "handle": "buger",
	//       "followers": 109
	//     },
	//     "avatars": [
	//       { "url": "https://avatars1.githubusercontent.com/u/14009?v=3&s=460", "type": "thumbnail" }
	//     ]
	//   },
	//   "company": {
	//     "name": "Acme"
	//   }
	// }`)

	// 	key, _, _, _ := jsonparser.Get(data, "person", "avatars", "[0]", "url")

	// 	fmt.Println("Key:", string(key))
}

func TestParamMiddleware(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/hello?person=jason(*)", nil)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	// fmt.Println("Response:", w.Body.String())

	// assert.Equal(t, http.StatusOK, w.Code)
	// assert.Equal(t, `{"message":"david"}`, w.Body.String())
}
