package auth

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
	"net/http"
	"time"
	"errors"
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

func InitOAuth(ctx echo.Context, callbackFunction OAuthCallbackFunction) error {
	state := uuid.NewString()
	// Note that it could lead to excessive memory usage if someone request it too many times
	connectionCache.SetDefault(state, callbackFunction)
	// TODO Use appropriate OAuth provider
	oauthProvider := GoogleOAuthProvider{}
	url := oauthProvider.GetOAuthLink(state)
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
	// TODO Use appropriate OAuth provider
	oauthProvider := GoogleOAuthProvider{}
	account, err := oauthProvider.FetchAccountData(ctx.Request().Context(), code)
	if err != nil {
		return err
	}
	return callbackFunction.Callback(account, ctx)
}

