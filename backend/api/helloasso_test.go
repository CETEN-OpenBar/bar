package api

import (
	"bar/internal/config"
	"testing"
)

func TestHelloAssoConfigured(t *testing.T) {
	tests := []struct {
		name         string
		clientID     string
		clientSecret string
		want         bool
	}{
		{"both empty", "", "", false},
		{"id only", "id", "", false},
		{"secret only", "", "secret", false},
		{"both set", "id", "secret", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c config.Config
			c.HelloAssoConfig.ClientID = tt.clientID
			c.HelloAssoConfig.ClientSecret = tt.clientSecret

			if got := helloAssoConfigured(c); got != tt.want {
				t.Errorf("helloAssoConfigured() = %v, want %v", got, tt.want)
			}
		})
	}
}
