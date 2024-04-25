package oauth

import (
	"net"
	"net/http"
	"time"
)

type OauthService struct {
	httpClient *http.Client
	github     Github
}

type Config struct {
	Github Github
}

func NewOauthService(config Config) *OauthService {
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   time.Second,
			ResponseHeaderTimeout: time.Second,
		},
	}
	return &OauthService{
		httpClient: client,
		github: Github{
			ClientId:     config.Github.ClientId,
			ClientSecret: config.Github.ClientSecret,
			Url:          config.Github.Url,
			RedirectUrl:  config.Github.RedirectUrl,
		},
	}
}
