package providers

import (
	"log"
	"net/http"
	"net/url"

	"github.com/bitly/oauth2_proxy/api"
)

type SlackProvider struct {
	*ProviderData
}

func NewSlackProvider(p *ProviderData) *SlackProvider {
	p.ProviderName = "Slack"
	if p.LoginURL == nil || p.LoginURL.String() == "" {
		p.LoginURL = &url.URL{
			Scheme: "https",
			Host:   "slack.com",
			Path:   "/oauth/authorize",
		}
	}
	if p.RedeemURL == nil || p.RedeemURL.String() == "" {
		p.RedeemURL = &url.URL{
			Scheme: "https",
			Host:   "slack.com",
			Path:   "/api/oauth.access",
		}
	}
	if p.ValidateURL == nil || p.ValidateURL.String() == "" {
		p.ValidateURL = &url.URL{
			Scheme: "https",
			Host:   "slack.com",
			Path:   "/api/users.identity",
		}
	}
	if p.Scope == "" {
		p.Scope = "identity.basic,identity.email"
	}
	return &SlackProvider{ProviderData: p}
}

func (p *SlackProvider) GetEmailAddress(s *SessionState) (string, error) {
	req, err := http.NewRequest("GET",
		p.ValidateURL.String()+"?token="+s.AccessToken, nil)
	if err != nil {
		log.Printf("failed building request %s", err)
		return "", err
	}
	json, err := api.Request(req)
	if err != nil {
		log.Printf("failed making request %s", err)
		return "", err
	}
	return json.Get("user").Get("email").String()
}
