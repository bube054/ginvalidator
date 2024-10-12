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

	// 	data := []byte(`{
	//   "person": {
	//     "name": {
	//       "first": "Leonid",
	//       "last": "Bugaev",
	//       "fullName": "Leonid Bugaev"
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

	body := NewBody("simple_key", nil)
	router.GET("/body", body.CreateChain().Not().Contains("errors", nil).Not().Alphanumeric(nil).Validate(), func(c *gin.Context) {
		errs, err := ValidationResult(c)

		if err != nil {
			fmt.Println("Could not retrieve validation result err:", err)
		} else {
			fmt.Printf("All Errors 🙌🙌🙌🙌 %+v\n", errs)
		}

		c.String(200, "pong")
	})

	return router
}

func TestBody(t *testing.T) {
	// data := []byte(`{"name":"John"}`)
	router := bodySetupRouter()

	body := `{
  "boolean_key": "--- true\n",
  "empty_string_translation": "",
  "key_with_description": "Check it out! This key has a description! (At least in some formats)",
  "key_with_line-break": "This translations contains\na line-break.",
  "nested": {
    "deeply": {
      "key": "Wow, this key is nested even deeper."
    },
    "key": "This key is nested inside a namespace."
  },
  "null_translation": null,
  "pluralized_key": {
    "one": "Only one pluralization found.",
    "other": "Wow, you have pluralizations!",
    "zero": "You have no pluralization."
  },
  "sample_collection": [
    "first item",
    "second item",
    "third item"
  ],
  "simple_key": "Just a simple key 69 with a simple message.",
  "unverified_key": "This translation is not yet 69 verified and waits for it. (In some formats we also export this status)"
	}`

	req, _ := http.NewRequest("GET", "/body", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("request error")
	}
}
