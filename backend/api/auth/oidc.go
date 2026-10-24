package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// OIDCProvider logs in through any OpenID Connect issuer, such as the dev
// server started by docker compose. Endpoints come from the issuer's
// discovery document, fetched on first use.
type OIDCProvider struct {
	Issuer       string
	ClientID     string
	ClientSecret string
	RedirectURL  string

	mu          sync.Mutex
	conf        *oauth2.Config
	userinfoURL string
}

type oidcDiscovery struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	UserinfoEndpoint      string `json:"userinfo_endpoint"`
}

type oidcUser struct {
	Sub        string `json:"sub"`
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
}

var oidcHTTPClient = &http.Client{Timeout: 5 * time.Second}

// discover returns the cached OAuth config, fetching it if needed. A failure
// isn't cached, so the next login tries again.
func (p *OIDCProvider) discover(ctx context.Context) (*oauth2.Config, string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.conf != nil {
		return p.conf, p.userinfoURL, nil
	}

	url := strings.TrimSuffix(p.Issuer, "/") + "/.well-known/openid-configuration"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}
	resp, err := oidcHTTPClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("OIDC discovery at %s: %s", url, resp.Status)
	}

	var d oidcDiscovery
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return nil, "", fmt.Errorf("OIDC discovery at %s: %w", url, err)
	}
	if d.AuthorizationEndpoint == "" || d.TokenEndpoint == "" || d.UserinfoEndpoint == "" {
		return nil, "", fmt.Errorf("OIDC discovery at %s: missing endpoints", url)
	}

	p.conf = &oauth2.Config{
		ClientID:     p.ClientID,
		ClientSecret: p.ClientSecret,
		RedirectURL:  p.RedirectURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  d.AuthorizationEndpoint,
			TokenURL: d.TokenEndpoint,
		},
	}
	p.userinfoURL = d.UserinfoEndpoint
	return p.conf, p.userinfoURL, nil
}

// GetOAuthLink returns "" when the issuer can't be reached.
func (p *OIDCProvider) GetOAuthLink(state string) string {
	conf, _, err := p.discover(context.Background())
	if err != nil {
		logrus.Error(err)
		return ""
	}
	return conf.AuthCodeURL(state)
}

func (p *OIDCProvider) FetchAccountData(requestContext context.Context, oAuthCode string) (*OAuthAccountData, error) {
	conf, userinfoURL, err := p.discover(requestContext)
	if err != nil {
		return nil, err
	}
	ctx := context.WithValue(requestContext, oauth2.HTTPClient, oidcHTTPClient)
	token, err := conf.Exchange(ctx, oAuthCode)
	if err != nil {
		return nil, err
	}

	resp, err := conf.Client(ctx, token).Get(userinfoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OIDC userinfo: %s", resp.Status)
	}

	usr := &oidcUser{}
	if err := json.NewDecoder(resp.Body).Decode(usr); err != nil {
		return nil, err
	}
	// Accounts are matched on this ID, so an empty one must never get through
	if usr.Sub == "" {
		return nil, errors.New("OIDC userinfo: no sub")
	}
	return &OAuthAccountData{
		FirstName:   usr.GivenName,
		LastName:    usr.FamilyName,
		EmailAdress: usr.Email,
		Id:          usr.Sub,
		PictureLink: usr.Picture,
	}, nil
}
