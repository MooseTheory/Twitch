package TwitchApi

import (
	"fmt"
	"time"
)

var (
	defaultScheme = "https"
	authHost      = "id.twitch.tv"
	apiHost       = "api.twitch.tv"
)

type TwitchApi struct {
	ClientID        string
	ClientSecret    string
	accessToken     string
	tokenValidUntil time.Time
}

func (ta *TwitchApi) Connect() error {
	// Check and make sure we have API keys
	if ta.ClientID == "" {
		return fmt.Errorf("client ID is not set")
	}
	if ta.ClientSecret == "" {
		return fmt.Errorf("client secret is not set")
	}
	// If we already have an access token see if it is still valid.
	if ta.accessToken != "" {
		isValid, err := ta.isTokenValid()
		if err != nil {
			return err
		}
		if isValid {
			return nil
		}
	}
	// OK, we have keys, but our token either doesn't exist, or is expired
	// Try to get a fresh token.
	err := ta.fetchAccessToken()
	if err != nil {
		return err
	}
	return nil
}

func (ta *TwitchApi) GetExistingSubs() {
	err := ta.addNewSubscription("57576022")
	if err != nil {
		panic(err)
	}
	err = ta.getCurrentSubscriptions()
	if err != nil {
		panic(err)
	}
}

func (ta *TwitchApi) PrintInfo() {
	fmt.Printf("%+v\n", ta)
}
