package auth

import (
	"bar/internal/config"
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
)

type googleUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	Link      string `json:"link"`
	Picture   string `json:"picture"`
}

type GoogleOAuthProvider struct {
}

var googleOAuth2Config *oauth2.Config

func getGoogleOAuthConfig() *oauth2.Config {
	// Don't sync it, we always create the same object
	if googleOAuth2Config == nil {
		conf := config.GetConfig()
		clientId := conf.OauthConfig.GoogleClientID
		clientSecret := conf.OauthConfig.GoogleClientSecret
		redirectUrl := conf.ApiConfig.BasePath + "/auth/google/callback"
		googleOAuth2Config = &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			RedirectURL:  redirectUrl,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/admin.directory.user.readonly",
			},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://accounts.google.com/o/oauth2/auth",
				TokenURL: "https://oauth2.googleapis.com/token",
			},
		}
	}
	return googleOAuth2Config
}

func (p *GoogleOAuthProvider) GetOAuthLink(state string) string {
	oAuth2Config := getGoogleOAuthConfig()
	domainParameter := oauth2.SetAuthURLParam("hd", "telecomnancy.net")
	return oAuth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline, domainParameter)
}

func (p *GoogleOAuthProvider) FetchAccountData(requestContext context.Context, oAuthCode string) (*OAuthAccountData, error) {
	oAuth2Config := getGoogleOAuthConfig()
	token, err := oAuth2Config.Exchange(requestContext, oAuthCode)
	if err != nil {
		return nil, err
	}

	// Get user from Google
	client := oAuth2Config.Client(requestContext, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	usr := &googleUser{}
	err = json.NewDecoder(resp.Body).Decode(usr)
	if err != nil {
		return nil, err
	}
	return &OAuthAccountData{
		usr.FirstName,
		usr.LastName,
		usr.Email,
		usr.ID,
		usr.Picture,
	}, nil
}
