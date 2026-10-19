package api

import (
	"bar/autogen"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestCallbackRedirectsOnError(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=unknown&code=x", nil)
	c := echo.New().NewContext(req, rec)

	if err := (&Server{}).Callback(c, autogen.CallbackParams{State: "unknown", Code: "x"}); err != nil {
		t.Fatalf("Callback() = %v, want a redirect", err)
	}
	if loc := rec.Header().Get("Location"); !strings.Contains(loc, "rt=authError") {
		t.Errorf("redirected to %q, want the auth error page", loc)
	}
}
