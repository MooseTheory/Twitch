package TwitchApi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type eventSubRequest struct {
	Type      string            `json:"type"`
	Version   string            `json:"version"`
	Condition eventSubCondition `json:"condition"`
	Transport eventSubTransport `json:"transport"`
}

type eventSubCondition struct {
	BroadcasterUserID     string `json:"broadcaster_user_id"`
	FromBroadcasterUserID string `json:"from_broadcaster_user_id"`
	ToBroadcasterUserID   string `json:"to_broadcaster_user_id"`
	RewardID              string `json:"reward_id"`
	OrganizationID        string `json:"organization_id"`
	CategoryID            string `json:"category_id"`
	CampaignID            string `json:"campaign_id"`
	ExtensionClientID     string `json:"extension_client_id"`
	ClientID              string `json:"client_id"`
	UserID                string `json:"user_id"`
}

type eventSubTransport struct {
	Method   string `json:"method"` // supported values: webhook
	Callback string `json:"callback"`
	Secret   string `json:"secret"`
}

type eventSubResponse struct {
	Data         eventSubResponseData `json:"data"`
	Total        int                  `json:"total"`
	TotalCost    int                  `json:"total_cost"`
	MaxTotalCost int                  `json:"max_total_cost"`
}

type eventSubResponseData struct {
	ID        string            `json:"id"`
	Status    string            `json:"status"`
	Type      string            `json:"type"`
	Version   string            `json:"version"`
	Condition eventSubCondition `json:"condition"`
	CreatedAt string            `json:"created_at"`
	Transport eventSubTransport `json:"transport"`
	Cost      int               `json:"cost"`
}

type eventSubscriptionsResponse struct {
	Data         []eventSubResponseData `json:"data"`
	Total        int                    `json:"total"`
	TotalCost    int                    `json:"total_cost"`
	MaxTotalCost int                    `json:"max_total_cost"`
	Pagination   map[string]interface{} `json:"pagination"`
	Cursor       map[string]interface{} `json:"cursor"`
}

func (ta *TwitchApi) getCurrentSubscriptions() error {
	endpointURL := url.URL{
		Scheme: defaultScheme,
		Host:   apiHost,
		Path:   "/helix/eventsub/subscriptions",
	}
	request := http.Request{
		Method: http.MethodGet,
		URL:    &endpointURL,
	}
	request.Header = http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", ta.accessToken)},
		"Client-Id":     []string{ta.ClientID},
	}
	client := http.Client{}
	resp, err := client.Do(&request)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(bodyBytes))
		return fmt.Errorf("bad status code \"%v\"", resp.StatusCode)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result eventSubscriptionsResponse
	err = json.Unmarshal(bodyBytes, &result)
	fmt.Printf("%+v\n", result)
	return err
}

func (ta *TwitchApi) addNewSubscription(channelId string) error {
	// Check if we're already subbed to this.
	endpointURL := url.URL{
		Scheme: defaultScheme,
		Host:   apiHost,
		Path:   "/helix/eventsub/subscriptions",
	}
	request := http.Request{
		Method: http.MethodPost,
		URL:    &endpointURL,
	}
	request.Header = http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", ta.accessToken)},
		"Client-Id":     []string{ta.ClientID},
		"Content-Type":  []string{"application/json"},
	}
	requestJson := eventSubRequest{
		Type:    "stream.online",
		Version: "1",
		Condition: eventSubCondition{
			BroadcasterUserID: channelId,
		},
		Transport: eventSubTransport{
			Method:   "webhook",
			Callback: "https://twitch.trailer.merr.is/event/stream-online",
			Secret:   "some super secure text goes here",
		},
	}
	requestBody, err := json.Marshal(requestJson)
	if err != nil {
		return err
	}
	request.Body = io.NopCloser(bytes.NewReader(requestBody))
	client := http.Client{}
	resp, err := client.Do(&request)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != 201 && resp.StatusCode != 202 {
		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(bodyBytes))
		return fmt.Errorf("bad status code \"%v\"", resp.StatusCode)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result eventSubscriptionsResponse
	err = json.Unmarshal(bodyBytes, &result)
	fmt.Printf("%+v\n", result)
	return err
}
