package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Action struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Url   string `json:"url"`
	Style string `json:"style"`
}

type Attachment struct {
	Fallback     string   `json:"fallback"`
	Color        string   `json:"color"`
	PreText      string   `json:"pretext"`
	AuthorName   string   `json:"author_name"`
	AuthorLink   string   `json:"author_link"`
	AuthorIcon   string   `json:"author_icon"`
	Title        string   `json:"title"`
	TitleLink    string   `json:"title_link"`
	Text         string   `json:"text"`
	ImageUrl     string   `json:"image_url"`
	Fields       []Field  `json:"fields"`
	Footer       string   `json:"footer"`
	FooterIcon   string   `json:"footer_icon"`
	Timestamp    int64    `json:"ts"`
	MarkdownIn   []string `json:"mrkdwn_in"`
	Actions      []Action `json:"actions"`
	CallbackID   string   `json:"callback_id"`
	ThumbnailUrl string   `json:"thumb_url"`
}

type Payload struct {
	Parse       string       `json:"parse,omitempty"`
	Username    string       `json:"username,omitempty"`
	IconUrl     string       `json:"icon_url,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Text        string       `json:"text,omitempty"`
	LinkNames   string       `json:"link_names,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	UnfurlLinks bool         `json:"unfurl_links,omitempty"`
	UnfurlMedia bool         `json:"unfurl_media,omitempty"`
	Markdown    bool         `json:"mrkdwn,omitempty"`
}

type Slack struct {
	URL     string
	Channel string
}

// func SendSlackNotification(slackURL string, slackChannel string, slackPayload Payload) error {
func (s *Slack) SendNotification(p Payload) error {
	// チャンネルの上書き
	if s.Channel != "" {
		p.Channel = s.Channel
	}

	payloadBytes, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack Payload: %w", err)
	}

	req, err := http.NewRequest("POST", s.URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create new HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %d", res.StatusCode)
	}

	return nil
}
