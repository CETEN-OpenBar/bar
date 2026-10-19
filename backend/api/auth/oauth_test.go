package auth

import (
	"bar/internal/config"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

func TestProviderFor(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		issuer   string
		wantOIDC bool
		wantErr  string
	}{
		{"unset", "", "", false, ""},
		{"google", "google", "", false, ""},
		{"oidc any case", "OIDC", "http://localhost:8082/default", true, ""},
		{"oidc without issuer", "oidc", "", false, "BAR_OAUTH_OIDC_ISSUER"},
		{"unknown", "github", "", false, "use google or oidc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c config.Config
			c.ApiConfig.BasePath = "http://localhost:8080"
			c.OauthConfig.Provider = tt.provider
			c.OauthConfig.OIDCIssuer = tt.issuer

			p, err := providerFor(c)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("providerFor() error = %v, want it to mention %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if tt.wantOIDC {
				o, ok := p.(*OIDCProvider)
				if !ok {
					t.Fatalf("providerFor() = %T, want *OIDCProvider", p)
				}
				if o.Issuer != tt.issuer || o.RedirectURL != "http://localhost:8080/auth/google/callback" {
					t.Errorf("got issuer %q, redirect %q", o.Issuer, o.RedirectURL)
				}
			} else if _, ok := p.(*GoogleOAuthProvider); !ok {
				t.Fatalf("providerFor() = %T, want *GoogleOAuthProvider", p)
			}
		})
	}
}

func TestGoogleOAuthLinkUnchanged(t *testing.T) {
	u, err := url.Parse((&GoogleOAuthProvider{}).GetOAuthLink("st"))
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	if u.Host != "accounts.google.com" || q.Get("state") != "st" || q.Get("hd") != "telecomnancy.net" || q.Get("access_type") != "offline" {
		t.Errorf("unexpected Google link %s", u)
	}
	for _, s := range []string{"userinfo.email", "userinfo.profile", "admin.directory.user.readonly"} {
		if !strings.Contains(q.Get("scope"), s) {
			t.Errorf("scope %q is missing %s", q.Get("scope"), s)
		}
	}
}
