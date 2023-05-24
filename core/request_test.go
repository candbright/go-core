package core

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
	"testing"
)

type Dog struct {
	Name string `json:"name"`
}

func TestGET(t *testing.T) {
	engine := gin.Default()
	engine.POST("/test", func(context *gin.Context) {
		POST[Dog](context, func(receive Dog) (interface{}, error) {
			return receive, nil
		})
	})
	err := engine.Run(":8080")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGETReq(t *testing.T) {
	resp, err := http.Post("http://127.0.0.1:8080/test", "application/json", strings.NewReader(`{
	"name": "lai fu"
}`))
	if err != nil {
		t.Fatal(err)
	}
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(all))
}
