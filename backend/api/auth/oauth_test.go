package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

type stubProvider struct{ link string }

func (p stubProvider) GetOAuthLink(state string) string { return p.link }

func (p stubProvider) FetchAccountData(context.Context, string) (*OAuthAccountData, error) {
	return &OAuthAccountData{}, nil
}

func useProvider(t *testing.T, p OAuthProvider, err error) {
	t.Helper()
	old := activeProvider
	activeProvider = func() (OAuthProvider, error) { return p, err }
	t.Cleanup(func() { activeProvider = old })
}

func newEchoContext() (echo.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/auth/google", nil)
	return echo.New().NewContext(req, rec), rec
}

func TestInitOAuthRefusesEmptyLink(t *testing.T) {
	useProvider(t, stubProvider{link: ""}, nil)
	c, rec := newEchoContext()

	if err := InitOAuth(c, nil); err == nil {
		t.Errorf("InitOAuth() = nil, redirected to %q", rec.Header().Get("Location"))
	}
}

func TestInitOAuthRedirectsToProvider(t *testing.T) {
	useProvider(t, stubProvider{link: "http://idp/authorize"}, nil)
	c, rec := newEchoContext()

	if err := InitOAuth(c, nil); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusFound || rec.Header().Get("Location") != "http://idp/authorize" {
		t.Errorf("got %d to %q", rec.Code, rec.Header().Get("Location"))
	}
}
