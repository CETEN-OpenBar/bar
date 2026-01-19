package helloasso

import (
	"bar/internal/config"
	"context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var tokenSource oauth2.TokenSource


func (c *Client) GetToken() (*oauth2.Token, error) {
	if tokenSource != nil {
		return tokenSource.Token()
	}

	conf := config.GetConfig()

	oauthConfig := &clientcredentials.Config{
		ClientID:     conf.HelloAssoConfig.ClientID,
		ClientSecret: conf.HelloAssoConfig.ClientSecret,
		TokenURL:     conf.HelloAssoConfig.URL + "/oauth2/token",
	}

	tokenSource = oauthConfig.TokenSource(context.Background())
	c.Client = oauthConfig.Client(context.Background())

	return tokenSource.Token()
}
