package oauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Github struct {
	ClientId     string
	ClientSecret string
	Url          string
	RedirectUrl  string
	UserUrl      string
}

type GithubUserData struct {
	Email     string `json:"email"`
	AvatarUrl string `json:"avatar_url"`
	Login     string `json:"login"`
	Name      string `json:"name"`
}

func (service OauthService) GetGithubAccessToken(code string) (string, error) {
	data := map[string]string{
		"client_id":     service.github.ClientId,
		"client_secret": service.github.ClientSecret,
		"code":          code,
	}

	requestJSON, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		service.github.Url,
		bytes.NewBuffer(requestJSON),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := service.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	var body struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
	}

	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return "", err
	}

	return body.AccessToken, nil
}

// email, avatar_url, login, name
func (service OauthService) GetGithubData(token string) (*GithubUserData, error) {

	req, err := http.NewRequest(
		http.MethodGet,
		service.github.UserUrl,
		nil,
	)
	if err != nil {
		return nil, err
	}

	authHeader := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", authHeader)

	res, err := service.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var userData GithubUserData

	err = json.NewDecoder(res.Body).Decode(&userData)
	if err != nil {
		return nil, err
	}

	return &userData, nil
}
