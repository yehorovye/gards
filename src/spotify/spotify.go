package spotify

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const baseURL = "https://api.spotify.com/v1"

type Client struct {
	Token string
	HTTP  *http.Client
}

type Image struct {
	URL    string
	Height int
	Width  int
}

type Track struct {
	ID       string
	Name     string
	Artists  []string
	Album    string
	Duration int
	URL      string
	Images   []Image
}

func New(token string) *Client {
	return &Client{
		Token: token,
		HTTP:  &http.Client{Timeout: 5 * time.Second},
	}
}

func GetClientCredentialsToken(clientID, clientSecret string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.Header.Set("Authorization", "Basic "+auth)

	httpClient := &http.Client{Timeout: 5 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("spotify: failed to fetch token: " + resp.Status)
	}

	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

func (c *Client) GetTrackFromURL(link string) (*Track, error) {
	id, err := extractID(link)
	if err != nil {
		return nil, err
	}
	return c.GetTrackByID(id)
}

func (c *Client) GetTrackByID(id string) (*Track, error) {
	req, err := http.NewRequest("GET", baseURL+"/tracks/"+id, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("spotify: " + resp.Status)
	}

	var d struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Album struct {
			Name   string `json:"name"`
			Images []struct {
				URL    string `json:"url"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"images"`
		} `json:"album"`
		Artists      []struct{ Name string } `json:"artists"`
		DurationMs   int                     `json:"duration_ms"`
		ExternalUrls map[string]string       `json:"external_urls"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return nil, err
	}

	var names []string
	for _, a := range d.Artists {
		names = append(names, a.Name)
	}

	var images []Image
	for _, img := range d.Album.Images {
		images = append(images, Image{
			URL:    img.URL,
			Height: img.Height,
			Width:  img.Width,
		})
	}

	return &Track{
		ID:       d.ID,
		Name:     d.Name,
		Artists:  names,
		Album:    d.Album.Name,
		Duration: d.DurationMs,
		URL:      d.ExternalUrls["spotify"],
		Images:   images,
	}, nil
}

func extractID(link string) (string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	if !strings.Contains(u.Host, "spotify.com") || !strings.HasPrefix(u.Path, "/track/") {
		return "", errors.New("invalid track url")
	}
	parts := strings.Split(u.Path, "/")
	if len(parts) < 3 {
		return "", errors.New("invalid track id")
	}
	return parts[2], nil
}
