package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// fakeIssuer serves the bits of an OIDC provider that OIDCProvider uses.
func fakeIssuer(t *testing.T, userinfo map[string]string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"issuer":                 srv.URL,
			"authorization_endpoint": srv.URL + "/authorize",
			"token_endpoint":         srv.URL + "/token",
			"userinfo_endpoint":      srv.URL + "/userinfo",
		})
	})
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"access_token":"at","token_type":"Bearer"}`))
	})
	mux.HandleFunc("/userinfo", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer at" {
			http.Error(w, "bad token", http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(userinfo)
	})
	return srv
}

var alice = map[string]string{
	"sub":         "alice-sub",
	"email":       "alice@telecomnancy.net",
	"given_name":  "alice",
	"family_name": "Test",
	"picture":     "http://example.com/alice.png",
}

func newTestOIDC(issuer string) *OIDCProvider {
	return &OIDCProvider{
		Issuer:       issuer,
		ClientID:     "bar-dev",
		ClientSecret: "dev",
		RedirectURL:  "http://localhost:8080/auth/google/callback",
	}
}

func TestOIDCLink(t *testing.T) {
	srv := fakeIssuer(t, alice)

	link := newTestOIDC(srv.URL).GetOAuthLink("st")
	u, err := url.Parse(link)
	if err != nil || link == "" {
		t.Fatalf("GetOAuthLink() = %q, %v", link, err)
	}
	if want, _ := url.Parse(srv.URL); u.Host != want.Host || u.Path != "/authorize" {
		t.Errorf("link points at %s%s, want %s/authorize", u.Host, u.Path, want.Host)
	}
	q := u.Query()
	for k, want := range map[string]string{
		"client_id":     "bar-dev",
		"redirect_uri":  "http://localhost:8080/auth/google/callback",
		"state":         "st",
		"scope":         "openid email profile",
		"response_type": "code",
	} {
		if q.Get(k) != want {
			t.Errorf("%s = %q, want %q", k, q.Get(k), want)
		}
	}
}

func TestOIDCFetchAccountData(t *testing.T) {
	srv := fakeIssuer(t, alice)

	got, err := newTestOIDC(srv.URL).FetchAccountData(context.Background(), "code")
	if err != nil {
		t.Fatal(err)
	}
	want := OAuthAccountData{
		FirstName:   "alice",
		LastName:    "Test",
		EmailAdress: "alice@telecomnancy.net",
		Id:          "alice-sub",
		PictureLink: "http://example.com/alice.png",
	}
	if *got != want {
		t.Errorf("FetchAccountData() = %+v, want %+v", *got, want)
	}
}

func TestOIDCFetchAccountDataRejectsEmptySubject(t *testing.T) {
	srv := fakeIssuer(t, map[string]string{"email": "nobody@telecomnancy.net"})

	if _, err := newTestOIDC(srv.URL).FetchAccountData(context.Background(), "code"); err == nil {
		t.Error("FetchAccountData() accepted a userinfo response without sub")
	}
}

func TestOIDCUnreachableIssuer(t *testing.T) {
	p := newTestOIDC("http://127.0.0.1:1")

	if link := p.GetOAuthLink("st"); link != "" {
		t.Errorf("GetOAuthLink() = %q, want empty", link)
	}
	if _, err := p.FetchAccountData(context.Background(), "code"); err == nil {
		t.Error("FetchAccountData() returned no error")
	}
}

func TestOIDCRetriesDiscoveryAfterFailure(t *testing.T) {
	srv := fakeIssuer(t, alice)
	p := newTestOIDC("http://127.0.0.1:1")
	p.GetOAuthLink("st")

	p.Issuer = srv.URL
	if link := p.GetOAuthLink("st"); link == "" {
		t.Error("GetOAuthLink() still empty once the issuer is reachable")
	}
}
