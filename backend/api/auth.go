package api

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

var qrCache = cache.New(5*time.Minute, 10*time.Minute)
var stateCache = cache.New(5*time.Minute, 10*time.Minute)
var redirectCache = cache.New(5*time.Minute, 10*time.Minute)

type StateCache struct {
	Type string
	Data interface{}
}

type QrCache struct {
	Type string
	Data interface{}
}

// (GET /account/qr)
func (s *Server) GetAccountQR(c echo.Context) error {
	// Get account from cookie
	account, err := MustGetUserOrOnBoard(c)
	if err != nil {
		return nil
	}

	var params autogen.GetAccountQRJSONBody
	err = c.Bind(&params)
	if err != nil {
		return Error400(c)
	}

	if !account.VerifyPin(params.CardPin) {
		return ErrorAccNotFound(c)
	}

	b64, found := qrCache.Get(account.Id.String())
	if !found {
		// Generate QR code nonce
		nonce := uuid.NewString()

		// Cache nonce
		qrCache.Set(nonce, &QrCache{
			Type: "linking",
			Data: account.Id.String(),
		}, cache.DefaultExpiration)

		conf := config.GetConfig()
		url := fmt.Sprintf("%s/auth/google/begin/%s", conf.ApiConfig.BasePath, nonce)
		qr, err := qrcode.New(url, qrcode.Medium)
		if err != nil {
			return Error500(c)
		}
		qr.BackgroundColor = color.RGBA{R: 255, G: 255, B: 255, A: 0}
		// Generate QR code
		png, err := qr.PNG(200)
		if err != nil {
			return Error500(c)
		}
		b64 = base64.StdEncoding.EncodeToString(png)
		qrCache.Set(account.Id.String(), b64, cache.DefaultExpiration)

		logrus.Debugf("QR code generated for account %s: %s", account.Id.String(), url)
	}

	// Convert to base64
	r := strings.NewReader(b64.(string))

	autogen.GetAccountQR200ImagepngResponse{
		ContentLength: int64(r.Len()),
		Body:          r,
	}.VisitGetAccountQRResponse(c.Response())
	return nil
}

// (GET /account/qr)
func (s *Server) GetAccountQRWebsocket(c echo.Context) error {
	_, err := MustGetUserOrOnBoard(c)
	if err != nil {
		return nil
	}

	return LinkUpgrade(c)
}

var scopes = []string{
	"https://www.googleapis.com/auth/userinfo.profile",
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/admin.directory.user.readonly",
}

// (GET /auth/google/begin/{qr_nonce})
func (s *Server) ConnectAccount(c echo.Context, qrNonce string) error {
	// Get account from nonce and delete nonce
	data, found := qrCache.Get(qrNonce)
	if !found {
		return ErrorNotAuthenticated(c)
	}

	qrCache.Delete(qrNonce)

	d := data.(*QrCache)

	if d.Type == "linking" {
		accountID := d.Data
		qrCache.Delete(accountID.(string))
		BroadcastToRoom(accountID.(string), []byte("scanned"))
	} else if d.Type == "qr_auth" {
		uid := d.Data
		BroadcastToRoom(uid.(string), []byte("scanned"))
	}

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
	stateCache.Set(state, &StateCache{
		Type: d.Type,
		Data: d.Data,
	}, cache.DefaultExpiration)

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

func ErrorRedirect(c echo.Context, err string) error {
	conf := config.GetConfig()
	return c.Redirect(http.StatusPermanentRedirect, conf.ApiConfig.FrontendBasePath+"/borne/mobile?rt=authError&rm="+err)
}

func SuccessRedirect(c echo.Context) error {
	conf := config.GetConfig()
	return c.Redirect(http.StatusPermanentRedirect, conf.ApiConfig.FrontendBasePath+"/borne/mobile?rt=authSuccess")
}

// (GET /auth/google/callback)
func (s *Server) Callback(c echo.Context, params autogen.CallbackParams) error {
	// Get account from state and delete state
	data, found := stateCache.Get(params.State)
	if !found {
		// This callback is used when connecting to the admin panel for example.
		// The users clicks a button to log in with Google.
		return s.CallbackInpromptu(c, params)
	}
	stateCache.Delete(params.State)

	state := data.(*StateCache)
	switch state.Type {
	case "qr_auth":
		// Used when connecting to a borne with the QR Code displayed
		return s.CallbackQRAuth(c, params, state)
	case "linking":
		// Used when linking an account to a Google account
		return s.CallbackLinking(c, params, state)
	default:
		// Default fallback that should not happen
		return s.CallbackLinking(c, params, state)
	}
}

func (s *Server) CallbackLinking(c echo.Context, params autogen.CallbackParams, state *StateCache) error {
	accountID := state.Data

	conf := config.GetConfig()

	account, err := s.DBackend.GetAccount(c.Request().Context(), accountID.(string))
	if err != nil {
		if err != mongo.ErrNoDocuments {
			logrus.Error(err)
			return ErrorRedirect(c, "#001")
		}
		// Check if account is onBoard
		acc, found := onBoardCache.Get(accountID.(string))
		if !found {
			logrus.Error(err)
			return ErrorRedirect(c, "#002")
		}
		account = acc.(*models.Account)
	}

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
		return ErrorRedirect(c, "#003")
	}

	// Get user from Google
	client := oauth2Config.Client(c.Request().Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		logrus.Error(err)
		return ErrorRedirect(c, "#004")
	}
	defer resp.Body.Close()

	usr := &googleUser{}
	err = json.NewDecoder(resp.Body).Decode(usr)
	if err != nil {
		logrus.Error(err)
		return ErrorRedirect(c, "#005")
	}

	adminService, err := admin.NewService(c.Request().Context(), option.WithTokenSource(oauth2Config.TokenSource(c.Request().Context(), token)))
	if err != nil {
		logrus.Error(err)
		return ErrorRedirect(c, "#006")
	}

	t, err := adminService.Users.Get(usr.ID).Projection("custom").CustomFieldMask("Education").ViewType("domain_public").Do()
	if err != nil {
		logrus.Error(err)
		return ErrorRedirect(c, "#007")
	}
	edc := &education{}
	err = json.Unmarshal(t.CustomSchemas["Education"], edc)
	if err != nil {
		logrus.Error(err)
		return ErrorRedirect(c, "#008")
	}

	account.FirstName = usr.FirstName
	account.LastName = usr.LastName
	account.EmailAddress = usr.Email
	account.GoogleId = &usr.ID
	account.GooglePicture = &usr.Picture

	if account.State == autogen.AccountNotOnBoarded {
		account.State = autogen.AccountOK

		// Check if an account with this Google ID and no Card ID exists
		acc, err := s.DBackend.GetAccountByEmail(c.Request().Context(), usr.Email)
		if err != nil {
			if err != mongo.ErrNoDocuments {
				logrus.Error(err)
				return ErrorRedirect(c, "#009")
			}

			err = s.DBackend.CreateAccount(c.Request().Context(), account)
			if err != nil {
				logrus.Error(err)
				return ErrorRedirect(c, "#010")
			}
		} else {
			if acc.CardId == nil {
				acc.CardId = account.CardId
			}

			err = s.DBackend.UpdateAccount(c.Request().Context(), acc)
			if err != nil {
				logrus.Error(err)
				return ErrorRedirect(c, "#011")
			}

			account = acc

			account.FirstName = usr.FirstName
			account.LastName = usr.LastName
			account.EmailAddress = usr.Email
			account.GoogleId = &usr.ID
			account.GooglePicture = &usr.Picture
		}

		// Delete ONBOARD cookie
		s.RemoveOnBoardCookie(c)
	} else {
		err = s.DBackend.UpdateAccount(c.Request().Context(), account)
		if err != nil {
			logrus.Error(err)
			return ErrorRedirect(c, "#012")
		}
	}

	BroadcastToRoom(accountID.(string), []byte("connected"))

	r, found := redirectCache.Get(params.State)
	if !found {
		return SuccessRedirect(c)
	}
	redirectCache.Delete(params.State)

	s.SetCookie(c, account)
	return c.Redirect(http.StatusPermanentRedirect, r.(string))
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

// (GET /auth/google)
func (s *Server) ConnectGoogle(c echo.Context, p autogen.ConnectGoogleParams) error {
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

// (GET /logout)
func (s *Server) Logout(c echo.Context) error {
	s.RemoveCookies(c)

	autogen.Logout204Response{}.VisitLogoutResponse(c.Response())
	return nil
}
