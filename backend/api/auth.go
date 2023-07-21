package api

import (
	"bar/autogen"
	"bar/internal/config"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
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

var qrCache = cache.New(5*time.Minute, 10*time.Minute)
var stateCache = cache.New(5*time.Minute, 10*time.Minute)

// (GET /account/qr)
func (s *Server) GetAccountQR(c echo.Context, params autogen.GetAccountQRParams) error {
	// Get account from cookie
	sess := s.getUserSess(c)
	accountID, ok := sess.Values["account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// Get account from database
	account, err := s.DBackend.GetAccount(accountID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Delete cookie
			sess.Options.MaxAge = -1
			sess.Save(c.Request(), c.Response())
			return ErrorAccNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	cardPin := fmt.Sprintf("%x", sha256.Sum256([]byte(params.CardPin)))
	if cardPin != account.CardPin {
		return ErrorAccNotFound(c)
	}

	// Generate QR code nonce
	nonce := uuid.NewString()

	// Cache nonce
	qrCache.Set(nonce, accountID, cache.DefaultExpiration)

	conf := config.GetConfig()
	url := fmt.Sprintf("%s/auth/google/begin/%s", conf.ApiConfig.BasePath, nonce)

	// Generate QR code
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return Error500(c)
	}

	r := bytes.NewReader(png)

	logrus.Debugf("QR code generated for account %s: %s", accountID, url)

	// Add headers for caching
	c.Response().Header().Set("Cache-Control", "max-age=300, public")
	c.Response().Header().Set("Expires", time.Now().Add(300*time.Second).Format(time.RFC1123))

	autogen.GetAccountQR200ImagepngResponse{
		Body:          r,
		ContentLength: int64(len(png)),
	}.VisitGetAccountQRResponse(c.Response())
	return nil
}

var scopes = []string{
	"https://www.googleapis.com/auth/userinfo.profile",
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/admin.directory.user.readonly",
}

// (GET /auth/google/begin/{qr_nonce})
func (s *Server) ConnectAccount(c echo.Context, qrNonce string) error {
	// Get account from nonce and delete nonce
	accountID, found := qrCache.Get(qrNonce)
	if !found {
		return ErrorNotAuthenticated(c)
	}
	qrCache.Delete(qrNonce)

	conf := config.GetConfig()

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

	// Cache state
	stateCache.Set(state, accountID, cache.DefaultExpiration)

	hostDomainOption := oauth2.SetAuthURLParam("hd", "telecomnancy.net")
	// Redirect to Google
	url := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline, hostDomainOption)

	return c.Redirect(301, url)
}

type education struct {
	Promo  uint64 `json:"Promotion"`
	Spé    string `json:"Approfondissement"`
	Statut uint64 `json:"Statut"`
}

type googleUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	Link      string `json:"link"`
	Picture   string `json:"picture"`
}

// (GET /auth/google/callback)
func (s *Server) Callback(c echo.Context, params autogen.CallbackParams) error {
	// Get account from state and delete state
	accountID, found := stateCache.Get(params.State)
	if !found {
		return s.CallbackInpromptu(c, params)
	}
	stateCache.Delete(params.State)

	account, err := s.DBackend.GetAccount(accountID.(string))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorAccNotFound(c)
		}
		return Error500(c)
	}

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
		return Error500(c)
	}

	// Get user from Google
	client := oauth2Config.Client(c.Request().Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return Error500(c)
	}
	defer resp.Body.Close()

	usr := &googleUser{}
	err = json.NewDecoder(resp.Body).Decode(usr)
	if err != nil {
		return Error500(c)
	}

	adminService, err := admin.NewService(c.Request().Context(), option.WithTokenSource(oauth2Config.TokenSource(c.Request().Context(), token)))
	if err != nil {
		return Error500(c)
	}

	t, err := adminService.Users.Get(usr.ID).Projection("custom").CustomFieldMask("Education").ViewType("domain_public").Do()
	if err != nil {
		return Error500(c)
	}
	edc := &education{}
	err = json.Unmarshal(t.CustomSchemas["Education"], edc)
	if err != nil {
		return Error500(c)
	}

	account.FirstName = usr.FirstName
	account.LastName = usr.LastName
	account.EmailAddress = usr.Email
	account.GoogleId = usr.ID
	account.GooglePicture = usr.Picture

	err = s.DBackend.UpdateAccount(account)
	if err != nil {
		return Error500(c)
	}

	autogen.Callback200JSONResponse{
		Account: &account.Account,
	}.VisitCallbackResponse(c.Response())
	return nil
}

// (GET /auth/google/callback)
func (s *Server) CallbackInpromptu(c echo.Context, params autogen.CallbackParams) error {

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
		return Error500(c)
	}

	// Get user from Google
	client := oauth2Config.Client(c.Request().Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return Error500(c)
	}
	defer resp.Body.Close()

	usr := &googleUser{}
	err = json.NewDecoder(resp.Body).Decode(usr)
	if err != nil {
		return Error500(c)
	}

	account, err := s.DBackend.GetAccountByGoogle(usr.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorAccNotFound(c)
		}
		return Error500(c)
	}

	adminService, err := admin.NewService(c.Request().Context(), option.WithTokenSource(oauth2Config.TokenSource(c.Request().Context(), token)))
	if err != nil {
		return Error500(c)
	}

	t, err := adminService.Users.Get(usr.ID).Projection("custom").CustomFieldMask("Education").ViewType("domain_public").Do()
	if err != nil {
		return Error500(c)
	}
	edc := &education{}
	err = json.Unmarshal(t.CustomSchemas["Education"], edc)
	if err != nil {
		return Error500(c)
	}

	account.FirstName = usr.FirstName
	account.LastName = usr.LastName
	account.EmailAddress = usr.Email
	account.GoogleId = usr.ID
	account.GooglePicture = usr.Picture

	err = s.DBackend.UpdateAccount(account)
	if err != nil {
		return Error500(c)
	}

	s.SetCookie(c, account)

	autogen.Callback200JSONResponse{
		Account: &account.Account,
	}.VisitCallbackResponse(c.Response())
	return nil
}

// (POST /auth/card)
func (s *Server) ConnectCard(c echo.Context) error {
	// Check that header "X-Local-Token" is set to the local token
	if c.Request().Header.Get("X-Local-Token") != config.GetConfig().ApiConfig.LocalToken {
		return ErrorNotAuthenticated(c)
	}

	var param autogen.ConnectCardJSONBody
	err := c.Bind(&param)
	if err != nil {
		return Error400(c)
	}

	account, err := s.DBackend.GetAccountByCard(param.CardId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorAccNotFound(c)
		}
		return Error500(c)
	}

	// SHA256 hash of the card ID
	hash := sha256.Sum256([]byte(param.CardPin))
	digest := fmt.Sprintf("%x", hash)

	if account.CardPin != digest {
		return ErrorAccNotFound(c)
	}

	s.SetCookie(c, account)

	autogen.ConnectCard200JSONResponse{
		Account: &account.Account,
	}.VisitConnectCardResponse(c.Response())
	return nil
}

// (GET /auth/google)
func (s *Server) ConnectGoogle(c echo.Context) error {
	conf := config.GetConfig()

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

	hostDomainOption := oauth2.SetAuthURLParam("hd", "telecomnancy.net")
	// Redirect to Google
	url := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline, hostDomainOption)

	return c.Redirect(301, url)
}

// (GET /logout)
func (s *Server) Logout(c echo.Context) error {
	s.RemoveCookie(c)

	autogen.Logout204Response{}.VisitLogoutResponse(c.Response())
	return nil
}
