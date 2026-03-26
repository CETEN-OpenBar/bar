package auth

import (
	"bar/autogen"
	"bar/internal/config"
	"bar/internal/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/color"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

var scopes = []string{
	"https://www.googleapis.com/auth/userinfo.profile",
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/admin.directory.user.readonly",
}

// Init OAuth
// - Get redirected url (and sanitize it)
// - Generate random state
// - Save redirect in cache with state as key
// - Get OAuth link (google specific)
// - Redirect to link

func InitOAuth(c echo.Context, p autogen.ConnectGoogleParams) error {
	conf := config.GetConfig()

	// Get ?r=
	rel := p.R

	// Check if it's a safe redirect (TODO: check if this is correct)
	switch rel {
	case "admin":
		rel = conf.ApiConfig.FrontendBasePath + "/admin"
	case "client/commande":
		rel = conf.ApiConfig.FrontendBasePath + "/client/commande"
	}
	// Init OAuth2 flow with Google
	oauth2Config := oauth2.Config{
		ClientID:     conf.OauthConfig.GoogleClientID,
		ClientSecret: conf.OauthConfig.GoogleClientSecret,
		RedirectURL:  fmt.Sprintf("%s/auth/google/callback", conf.ApiConfig.BasePath),
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}

	// state is not nonce
	state := uuid.NewString()

	redirectCache.Set(state, rel, cache.DefaultExpiration)

	hostDomainOption := oauth2.SetAuthURLParam("hd", "telecomnancy.net")
	// Redirect to Google
	url := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline, hostDomainOption)

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// OAuth callback
// - Get Token, Generate a Client
// - Retrieive account information (provider specific)
// - Pull account from database (might be provider specific)
// - Update cached account properties
// - Update account in database
// - Save account in session coockie (Currently only if a redirect was specified, might be a bug)

func OAuthCallback(c echo.Context, params autogen.CallbackParams) error {
	conf := config.GetConfig()

	// Get token from Google
	oauth2Config := oauth2.Config{
		ClientID:     conf.OauthConfig.GoogleClientID,
		ClientSecret: conf.OauthConfig.GoogleClientSecret,
		RedirectURL:  fmt.Sprintf("%s/auth/google/callback", conf.ApiConfig.BasePath),
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}

	token, err := oauth2Config.Exchange(c.Request().Context(), params.Code)
	if err != nil {
		logrus.Error(err)
		return ErrorRedirect(c, "#014")
	}

	// Get user from Google
	client := oauth2Config.Client(c.Request().Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		logrus.Error(err)
		return ErrorRedirect(c, "#015")
	}
	defer resp.Body.Close()

	usr := &googleUser{}
	err = json.NewDecoder(resp.Body).Decode(usr)
	if err != nil {
		logrus.Error(err)
		return ErrorRedirect(c, "#016")
	}

	account, err := s.DBackend.GetAccountByGoogle(c.Request().Context(), usr.ID)
	if err != nil {
		account, err = s.DBackend.GetAccountByEmail(c.Request().Context(), usr.Email)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				// Redirect to the auth page with an error message
				conf := config.GetConfig()
				return c.Redirect(http.StatusPermanentRedirect, conf.ApiConfig.FrontendBasePath+"/auth?noaccount")
			}
			logrus.Error(err)
			return ErrorRedirect(c, "#017")
		}
	}

	logrus.WithField("account", account.Name()).Info("Account logged in using OAuth.")
	adminService, err := admin.NewService(c.Request().Context(), option.WithTokenSource(oauth2Config.TokenSource(c.Request().Context(), token)))
	if err != nil {
		logrus.Error(err)
		return ErrorRedirect(c, "#018")
	}

	t, err := adminService.Users.Get(usr.ID).Projection("custom").CustomFieldMask("Education").ViewType("domain_public").Do()
	if err != nil {
		logrus.Error(err)
		return ErrorRedirect(c, "#019")
	}
	edc := &education{}
	err = json.Unmarshal(t.CustomSchemas["Education"], edc)
	if err != nil {
		logrus.Error(err)
		return ErrorRedirect(c, "#020")
	}

	account.FirstName = usr.FirstName
	account.LastName = usr.LastName
	account.EmailAddress = usr.Email
	account.GoogleId = &usr.ID
	account.GooglePicture = &usr.Picture

	err = s.DBackend.UpdateAccount(c.Request().Context(), account)
	if err != nil {
		logrus.Error(err)
		return ErrorRedirect(c, "#021")
	}

	r, found := redirectCache.Get(params.State)
	if !found {
		return SuccessRedirect(c)
	}
	redirectCache.Delete(params.State)

	s.SetCookie(c, account)
	return c.Redirect(http.StatusPermanentRedirect, r.(string))
}
