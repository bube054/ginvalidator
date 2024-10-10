package ginvalidator

import (
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

	body := NewBody("name", nil)
	router.GET("/body", body.CreateChain().Contains("y", nil).Validate(), func(c *gin.Context) {
		c.String(200, "pong")
	})

	return router
}

func TestBody(t *testing.T) {
	// data := []byte(`{"name":"John"}`)
	router := bodySetupRouter()

	body := `{"name":"John"}`
	req, _ := http.NewRequest("GET", "/body", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("request error")
	}
}
