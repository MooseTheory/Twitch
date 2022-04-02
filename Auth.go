package TwitchApi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type oauthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type oauthTokenValidationResponse struct {
	ClientID  string   `json:"client_id"`
	Login     string   `json:"login"`
	Scopes    []string `json:"scopes"`
	UserID    string   `json:"user_id"`
	ExpiresIn int      `json:"expires_in"`
}

func (ta *TwitchApi) fetchAccessToken() error {
	endpointURL := url.URL{
		Scheme: defaultScheme,
		Host:   authHost,
		Path:   "/oauth2/token",
	}
	urlValues := url.Values{}
	urlValues.Set("client_id", ta.ClientID)
	urlValues.Add("client_secret", ta.ClientSecret)
	urlValues.Add("grant_type", "client_credentials")
	endpointURL.RawQuery = urlValues.Encode()

	request := http.Request{
		Method: http.MethodPost,
		URL:    &endpointURL,
	}
	client := http.Client{}
	resp, err := client.Do(&request)
	var result oauthTokenResponse
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code \"%v\"", resp.StatusCode)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return err
	}
	ta.accessToken = result.AccessToken
	ta.tokenValidUntil = time.Now().Add(time.Duration(result.ExpiresIn * int(time.Second)))
	return nil
}

func (ta *TwitchApi) isTokenValid() (bool, error) {
	if ta.tokenValidUntil.Before(time.Now()) {
		return false, nil
	}
	endpointURL := url.URL{
		Scheme: defaultScheme,
		Host:   authHost,
		Path:   "/oauth2/validate",
	}
	request := http.Request{
		Method: http.MethodGet,
		URL:    &endpointURL,
	}
	request.Header = http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", ta.accessToken)},
	}
	client := http.Client{}
	resp, err := client.Do(&request)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("bad status code \"%v\"", resp.StatusCode)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var result oauthTokenValidationResponse
	err = json.Unmarshal(bodyBytes, &result)
	return result.ClientID != "" && err != nil, err
}
