package futuapi

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestFutuAPI_initConnect(t *testing.T) {
	c := NewFutuAPIT(1, "jinxtest")
	ctx := context.Background()
	r, err := c.initConnect(ctx, "127.0.0.1.11111")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
}

type indexHandler struct {
	content string
}

func (ih *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, ih.content)
}

func TestHttp(t *testing.T) {
	http.Handle("/", &indexHandler{content: "jinx"})
	http.ListenAndServe(":8081", nil)
}
