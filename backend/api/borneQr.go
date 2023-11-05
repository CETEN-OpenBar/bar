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
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
)

// (GET /auth/qr)
func (s *Server) GetBorneAuthQRWebsocket(c echo.Context) error {
	// Check that header "X-Local-Token" is set to the local token
	if c.Request().Header.Get("X-Local-Token") != config.GetConfig().ApiConfig.LocalToken {
		return ErrorNotAuthenticated(c)
	}

	return AuthUpgrade(c)
}

func AuthUpgrade(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return Error500(c)
	}

	uid := uuid.New().String()

	room := GetWSRoom(uid)
	room.Add(conn)

	// Define a ticker to send a new QR code every 5 minutes starting now
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			err := RefreshQRCode(conn, uid)
			if err != nil {
				break
			}
		}
	}()

	RefreshQRCode(conn, uid)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	room.Remove(conn)

	return nil
}

func RefreshQRCode(conn *websocket.Conn, uid string) error {
	// Generate QR code nonce
	nonce := uuid.NewString()

	// Cache nonce
	qrCache.Set(nonce, &QrCache{
		Type: "qr_auth",
		Data: uid,
	}, cache.DefaultExpiration)

	conf := config.GetConfig()
	url := fmt.Sprintf("%s/auth/google/begin/%s", conf.ApiConfig.BasePath, nonce)
	qr, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		logrus.WithError(err).Error("failed to generate QR code")
		return err
	}
	qr.BackgroundColor = color.RGBA{R: 255, G: 255, B: 255, A: 0}
	// Generate QR code
	png, err := qr.PNG(200)
	if err != nil {
		logrus.WithError(err).Error("failed to generate QR code")
		return err
	}
	b64 := base64.StdEncoding.EncodeToString(png)

	logrus.Debugf("QR code generated for borne %s: %s", uid, url)

	data := map[string]string{
		"type": "qr_code",
		"data": b64,
	}

	// Send the QR code to the client
	err = conn.WriteJSON(data)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) CallbackQRAuth(c echo.Context, params autogen.CallbackParams, state *StateCache) error {
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
		return ErrorRedirect(c, "#022")
	}

	// Get user from Google
	client := oauth2Config.Client(c.Request().Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		logrus.Error(err)
		return ErrorRedirect(c, "#023")
	}
	defer resp.Body.Close()

	usr := &googleUser{}
	err = json.NewDecoder(resp.Body).Decode(usr)
	if err != nil {
		logrus.Error(err)
		return ErrorRedirect(c, "#024")
	}

	account, err := s.DBackend.GetAccountByGoogle(c.Request().Context(), usr.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorAccNotFound(c)
		}
		logrus.Error(err)
		return ErrorRedirect(c, "#025")
	}

	uid := state.Data
	room := GetWSRoom(uid.(string))
	nonce := uuid.New().String()
	connectionCache.Set(nonce, account, cache.DefaultExpiration)
	room.Broadcast([]byte(`{"type": "nonce", "data": "` + nonce + `"}`))

	return c.Redirect(http.StatusPermanentRedirect, conf.ApiConfig.FrontendBasePath+"/borne/connected")
}

var connectionCache = cache.New(5*time.Minute, 10*time.Minute)

// (POST /auth/qr)
func (s *Server) PostBorneAuthQR(c echo.Context) error {
	var bdy autogen.PostBorneAuthQRJSONBody
	err := c.Bind(&bdy)
	if err != nil {
		return Error400(c)
	}

	// Check that header "X-Local-Token" is set to the local token
	if c.Request().Header.Get("X-Local-Token") != config.GetConfig().ApiConfig.LocalToken {
		return ErrorNotAuthenticated(c)
	}

	// Check that the nonce is in the cache
	cached, found := qrCache.Get(bdy.Nonce)
	if !found {
		return Error400(c)
	}

	account := cached.(*models.Account)

	s.SetCookie(c, account)

	autogen.PostBorneAuthQR200JSONResponse{
		Account: &account.Account,
	}.VisitPostBorneAuthQRResponse(c.Response())
	return nil
}
