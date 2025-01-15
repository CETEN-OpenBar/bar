package webhook

import (
	"bar/internal/config"
	"bar/internal/models"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type EmbedFooter struct {
	Text      string `json:"text"`
	IconURL   string `json:"icon_url,omitempty"`
	Timestamp string `json:"timestamp"`
}

type DiscordEmbed struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Color       int         `json:"color"`
	Footer      EmbedFooter `json:"footer"`
	Timestamp   string      `json:"timestamp"`
}

type DiscordWebhook struct {
	Content   string         `json:"content,omitempty"`
	Username  string         `json:"username,omitempty"`
	AvatarURL string         `json:"avatar_url,omitempty"`
	Embeds    []DiscordEmbed `json:"embeds,omitempty"`
}

func SendDiscordWebhook(restock models.Restock) error {

	c := config.GetConfig()
	webhookURL := c.DiscordWebhookURL
	avatarURL := "https://bar.telecomnancy.net/logo.png" // The URL of the bot's avatar
	botUsername := "BarBot"                              // The name that appears as the webhook sender

	var formattedPastries string
	for i, item := range restock.Items {
		if i > 0 {
			formattedPastries += "\n"
		}
		formattedPastries += "**" + item.ItemName + " : " + strconv.FormatUint(item.AmountOfBundle*item.AmountPerBundle, 10) + "**"
	}

	// Create footer text with creator name
	footerText := "Updated by " + restock.CreatedByName

	currentTime := time.Now().Format(time.RFC3339)

	embed := DiscordEmbed{
		Title:       "Reappro Viennoiseries",
		Description: formattedPastries,
		Color:       16776960, // Yellow color for croissant
		Footer: EmbedFooter{
			Text: footerText,
		},
		Timestamp: currentTime, // This will be properly formatted by Discord
	}

	payload := DiscordWebhook{
		Username:  botUsername, // This is the bot's display name
		AvatarURL: avatarURL,
		Embeds:    []DiscordEmbed{embed},
	}

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		logrus.Errorf("unexpected status code: %d", resp.StatusCode)
		return errors.New("unexpected status code")
	}

	return nil
}
