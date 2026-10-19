package auth

import (
	"bar/internal/config"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type OAuthAccountData struct {
	FirstName   string
	LastName    string
	EmailAdress string
	Id          string
	PictureLink string
}

var InvalidOAuthStateError = errors.New("Invalid OAuth state")
var BrokenOAuthCallbackError = errors.New("Broken OAuth callback")

type OAuthProvider interface {
	GetOAuthLink(state string) string
	FetchAccountData(requestContext context.Context, oAuthCode string) (*OAuthAccountData, error)
}

type OAuthCallbackFunction interface {
	Callback(*OAuthAccountData, echo.Context) error
}

var connectionCache = cache.New(5*time.Minute, 10*time.Minute)

var ProviderUnavailableError = errors.New("OAuth provider unavailable")

// providerFor picks the OAuth provider named by BAR_OAUTH_PROVIDER.
func providerFor(c config.Config) (OAuthProvider, error) {
	switch strings.ToLower(c.OauthConfig.Provider) {
	case "", "google":
		return &GoogleOAuthProvider{}, nil
	case "oidc":
		if c.OauthConfig.OIDCIssuer == "" {
			return nil, errors.New("BAR_OAUTH_OIDC_ISSUER is required when BAR_OAUTH_PROVIDER is oidc")
		}
		return &OIDCProvider{
			Issuer:       c.OauthConfig.OIDCIssuer,
			ClientID:     c.OauthConfig.OIDCClientID,
			ClientSecret: c.OauthConfig.OIDCClientSecret,
			RedirectURL:  c.ApiConfig.BasePath + "/auth/google/callback",
		}, nil
	default:
		return nil, fmt.Errorf("unknown BAR_OAUTH_PROVIDER %q, use google or oidc", c.OauthConfig.Provider)
	}
}

// activeProvider is built once from the config, and replaced in tests.
var activeProvider = sync.OnceValues(func() (OAuthProvider, error) {
	return providerFor(config.GetConfig())
})

// CheckProvider logs which OAuth provider logins go through, or why there
// isn't one. Logins fail until the config is fixed.
func CheckProvider() {
	p, err := activeProvider()
	if err != nil {
		logrus.Error(err)
		return
	}
	if o, ok := p.(*OIDCProvider); ok {
		logrus.Infof("OAuth provider: oidc (%s)", o.Issuer)
		return
	}
	logrus.Info("OAuth provider: google")
}

func InitOAuth(ctx echo.Context, callbackFunction OAuthCallbackFunction) error {
	state := uuid.NewString()
	// Note that it could lead to excessive memory usage if someone request it too many times
	oauthProvider, err := activeProvider()
	if err != nil {
		return err
	}
	url := oauthProvider.GetOAuthLink(state)
	if url == "" {
		return ProviderUnavailableError
	}
	connectionCache.SetDefault(state, callbackFunction)
	return ctx.Redirect(http.StatusFound, url)
}

func ExecuteOAuthCallback(ctx echo.Context, state string, code string) error {
	rawCallbackFunction, found := connectionCache.Get(state)
	if !found {
		// Ignore all request with invalid state for security reasons
		return InvalidOAuthStateError
	}
	connectionCache.Delete(state)
	callbackFunction, ok := rawCallbackFunction.(OAuthCallbackFunction)
	if !ok {
		// Should not happen. But who knows
		return BrokenOAuthCallbackError
	}
	oauthProvider, err := activeProvider()
	if err != nil {
		return err
	}
	account, err := oauthProvider.FetchAccountData(ctx.Request().Context(), code)
	if err != nil {
		return err
	}
	return callbackFunction.Callback(account, ctx)
}
