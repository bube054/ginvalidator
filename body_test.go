package ginvalidator

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return r
}

func bodySetupRouter() *gin.Engine {
	router := setupRouter()

	body := NewBody("message", nil)
	router.GET("/body",
		body.
			CreateChain().
			Contains("errors", nil).
			Blacklist("0-9").
			Alphanumeric(nil).
			Blacklist("0-9").
			Validate(),
		func(c *gin.Context) {
			errs, err := ValidationResult(c)

			if err != nil {
				fmt.Println("Could not retrieve validation result err:", err)
			} else {
				fmt.Printf("All Errors 🙌🙌🙌🙌 %+v\n", errs)
			}

			data, err := MatchedData(c)

			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("data:", data)

			c.String(200, "pong")
		},
	)

	return router
}

func TestBody(t *testing.T) {
	// data := []byte(`{"name":"John"}`)
	router := bodySetupRouter()

	body := `{
		"name": {"first": "Tom", "last": "Anderson"},
		"age":37,
		"message": "A good saying is 7 comes after ate.",
		"children": ["Sara","Alex","Jack"],
		"fav.movie": "Deer Hunter",
		"friends": [
			{"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
			{"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
			{"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
		]
	}`

	req, _ := http.NewRequest("GET", "/body", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("request error")
	}
}
