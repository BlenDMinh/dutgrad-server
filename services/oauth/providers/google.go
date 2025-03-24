package providers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/BlenDMinh/dutgrad-server/configs"
	"github.com/BlenDMinh/dutgrad-server/services/oauth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOAuthProvider struct {
	config *oauth2.Config
}

func NewGoogleOAuthProvider() *GoogleOAuthProvider {
	config := configs.GetEnv()
	return &GoogleOAuthProvider{
		config: &oauth2.Config{
			ClientID:     config.OAuth.Google.ClientID,
			ClientSecret: config.OAuth.Google.ClientSecret,
			RedirectURL:  config.OAuth.Google.RedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (p *GoogleOAuthProvider) GetConfig() *oauth2.Config {
	return p.config
}

func (p *GoogleOAuthProvider) GetProviderName() string {
	return "google"
}

func (p *GoogleOAuthProvider) GetUserInfo(token *oauth2.Token) (*oauth.OAuthUserInfo, error) {
	client := p.config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	defer resp.Body.Close()

	var googleUser struct {
		Email    string `json:"email"`
		Username string `json:"given_name"`
		Sub      string `json:"sub"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %v", err)
	}

	return &oauth.OAuthUserInfo{
		Email:    googleUser.Email,
		Username: googleUser.Username,
		ID:       googleUser.Sub,
		Provider: "google",
	}, nil
}
