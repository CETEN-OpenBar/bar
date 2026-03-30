package api

import (
	"bar/autogen"
	"bar/api/auth"
	"bar/internal/config"
	"bar/internal/models"
	"bar/internal/db"
	"encoding/base64"
	"errors"
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
)

var qrCache = cache.New(5*time.Minute, 10*time.Minute)

type connectionOAuthCallback struct {
	database db.DBackend
	redirectUrl string
}

type linkingOAuthCallback struct {
	database db.DBackend
	accountId string
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
	return auth.InitOAuth(c, connectionOAuthCallback{s.DBackend, rel})
}


// POST /account/qr
// - Ask for user or boarded user
// - Check account pin (passed by POST)
// - If account does not have a cached qrcode
// -   Generate unique qr code storing HOST/auth/google/begin/QR_CODE_ID
// -   Encode qr code in base64 and save in cache (with type 'linking' and id. Currently with 2 entries)
// - Send response with qr code as base64

// (POST /account/qr)
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

// (GET /auth/google/begin/{qr_nonce})
// - Retrieve stored qr code in cache (error if not in cache)
// - Delete from cache
// - Send through websocket "scanned"
// - Get OAuth link
// - Redirect to OAuth link
func (s *Server) ConnectAccount(c echo.Context, qrNonce string) error {
	// Get account from nonce and delete nonce
	data, found := qrCache.Get(qrNonce)
	if !found {
		return ErrorNotAuthenticated(c)
	}

	qrCache.Delete(qrNonce)

	d := data.(*QrCache)

		accountID := d.Data
		qrCache.Delete(accountID.(string))
		BroadcastToRoom(accountID.(string), []byte("scanned"))

	// accountId

	return auth.InitOAuth(c, linkOAuthCallback{accountID})
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
func (s *Server) Callback(ctx echo.Context, params autogen.CallbackParams) error {
	err := auth.ExecuteOAuthCallback(ctx, params.State, params.Code)
	if errors.Is(err, auth.InvalidOAuthStateError) {
	}
	if errors.Is(err, auth.BrokenOAuthCallbackError) {
	}
	// TODO Do something on error
	return err
}

// - Get Token, Generate a Client
// - Retrieive account information (provider specific)
// - Pull account from database (might be provider specific)
// - Update cached account properties
// - Update account in database
// - Save account in session coockie (Currently only if a redirect was specified, might be a bug)
func (callback connectionOAuthCallback) Callback(accountData *auth.OAuthAccountData, ctx echo.Context) error {
	logrus.WithField("account", accountData.EmailAdress).Info("Account logged in using OAuth.")

	requestCtx := ctx.Request().Context()
	database := callback.database
	account, err := database.GetAccountByGoogle(requestCtx, accountData.Id)
	if err != nil {
		account, err = database.GetAccountByEmail(requestCtx, accountData.EmailAdress)
	}
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Redirect to the auth page with an error message
			conf := config.GetConfig()
			return ctx.Redirect(http.StatusPermanentRedirect, conf.ApiConfig.FrontendBasePath+"/auth?noaccount")
		}
		logrus.Error(err)
		return ErrorRedirect(c, "#017")
	}

	account.FirstName = accountData.FirstName
	account.LastName = accountData.LastName
	account.EmailAddress = accountData.EmailAdress
	account.GoogleId = &accountData.Id
	account.GooglePicture = &accountData.PictureLink

	err = database.UpdateAccount(ctx.Request().Context(), account)
	if err != nil {
		logrus.Error(err)
		return ErrorRedirect(c, "#021")
	}

	s.SetCookie(c, account)
	return ctx.Redirect(http.StatusFound, callback.redirectUrl)
}

// OAuth callback for qr_code linking
// - Pull account from database based on id stored in cache (identified by state from oauth)
// - If no account is found, check if the account was just created on board and pull it from there (cf. /auth/card)
// - Get token from OAuth
// - Get client and pull accounts info from OAuth
// - If account state is not on boarded (understand not linked to a card yet)
// -   Set account state to ok
// -   Pull account from database from email
// -   If no account found in db
// -     Add account to db
// -   Else
// -     If account in db has no card id, set it to to previously pulled account
// -   Remove on board session coockie
// - Else
// -   Update account in db
// - Broadcast with websocket "connected"
// - Eventually (?) set account coockie
// - Redirect to url
func (callback linkingOAuthCallback) Callback(accountData *auth.OAuthAccountData, ctx echo.Context) error {
	accountId := callback.accountId

	account, err := s.DBackend.GetAccount(c.Request().Context(), accountId)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			logrus.Error(err)
			return ErrorRedirect(c, "#001")
		}
		// Check if account is onBoard
		acc, found := onBoardCache.Get(accountId)
		if !found {
			logrus.Error(err)
			return ErrorRedirect(c, "#002")
		}
		account = acc.(*models.Account)
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

// (GET /logout)
func (s *Server) Logout(c echo.Context) error {
	s.RemoveCookies(c)

	autogen.Logout204Response{}.VisitLogoutResponse(c.Response())
	return nil
}
