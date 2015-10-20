package wctx

import (
	"golang.org/x/net/context"

	"net/http"
	"net/http/httptest"
	"testing"
)

// Ensures that Handler and HandleFund are what we think they are.
var a HandleFunc = func(ctx context.Context, w http.ResponseWriter, r *http.Request) {}
var _ Handler = a
var _ http.Handler = New()

func TestGet(t *testing.T) {
	r := New()
	r.GET("/page/:pageid/:mod", func(c context.Context, w http.ResponseWriter, req *http.Request) {
		p, ok := FromContext(c)
		if !ok {
			t.Error("cannot get params")
		}
		if pid := p.ByName("pageid"); pid != "12" {
			t.Errorf("wrong pageid, want 12, got %s", pid)
		}
		if mod := p.ByName("mod"); mod != "sub" {
			t.Errorf("wrong mod, want `sub`, got %q", mod)
		}
	})
	rw := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/page/12/sub", nil)
	r.ServeHTTP(rw, req)
	if rw.Code != 200 {
		t.Errorf("http response code want: 200, got: %d", rw.Code)
	}
}
