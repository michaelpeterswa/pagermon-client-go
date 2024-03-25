// Package pagermomn provides a client and helper functions for sending messages to a PagerMon server
package pagermon

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/michaelpeterswa/go-lib/multimonng"
)

type PagerMonClient struct {
	httpClient *http.Client
	apiKey     string
	baseUrl    string
}

func NewPagerMonClient(httpClient *http.Client, apiKey string, baseUrl string) *PagerMonClient {
	return &PagerMonClient{httpClient: httpClient, apiKey: apiKey, baseUrl: strings.TrimSuffix(baseUrl, "/")}
}

type PagerMonMessage struct {
	Address  string `json:"address"`
	Message  string `json:"message"`
	CurrTime int64  `json:"datetime"`
	Source   string `json:"source"`
}

func NewPagerMonMessage(currTime time.Time, address string, source string, message string) *PagerMonMessage {
	return &PagerMonMessage{CurrTime: currTime.Unix(), Address: address, Source: source, Message: message}
}

func (pc *PagerMonClient) getMessageEndpoint() string {
	return pc.baseUrl + "/api/messages"
}

func messageToForm(message *PagerMonMessage) url.Values {
	form := url.Values{}
	form.Add("address", message.Address)
	form.Add("message", message.Message)
	form.Add("source", message.Source)
	form.Add("datetime", strconv.FormatInt(message.CurrTime, 10))
	return form
}

func (pc *PagerMonClient) SendMessage(ctx context.Context, message *PagerMonMessage) error {
	req, err := http.NewRequestWithContext(ctx, "POST", pc.getMessageEndpoint(), strings.NewReader(messageToForm(message).Encode()))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.ParseForm()

	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("User-Agent", "pagermon-client-go")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("apikey", pc.apiKey)

	resp, err := pc.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %s", resp.Status)
	}

	return nil
}

func MultimonNGMessageToPagerMonMessage(m *multimonng.MultimonNGMessage, identifier string) *PagerMonMessage {
	return NewPagerMonMessage(time.Now(), m.Address, identifier, m.Alpha)
}
