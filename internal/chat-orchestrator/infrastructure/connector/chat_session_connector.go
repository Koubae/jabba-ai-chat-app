package connector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/model"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/utils"
	"golang.org/x/net/context"
	"io"
	"log"
	"net/http"
	"time"
)

type ChatSessionConnector struct {
	Host           string
	Port           int
	TimeoutSeconds int
	URL            string
	httpClient     *http.Client
}

type Request struct {
	SessionID string `json:"session_id"`
	Name      string `json:"name"`
	MemberID  string `json:"member_id"`
	Channel   string `json:"channel"`
}

// type Response struct {
// 	ChatURL       string     `json:"chat_url"`
// 	ID            string     `json:"id"`
// 	ApplicationId string     `json:"application_id"`
// 	Name          string     `json:"name"`
// 	Created       *time.Time `json:"created"`
// 	Updated       *time.Time `json:"updated"`
// }

func NewChatSessionConnector() *ChatSessionConnector {
	host := utils.GetEnvString("CHAT_SESSION_HOST", "http://localhost")
	port := utils.GetEnvInt("CHAT_SESSION_PORT", 20002)
	timeoutSeconds := utils.GetEnvInt("CHAT_SESSION_TIMEOUT_SECONDS", 120)

	return &ChatSessionConnector{
		Host:           host,
		Port:           port,
		TimeoutSeconds: timeoutSeconds,
		URL:            fmt.Sprintf("%s:%d", host, port),

		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
	}
}

func (c *ChatSessionConnector) CreateSession(
	ctx context.Context,
	accessToken string,
	sessionID string,
	name string,
	memberID string,
	channel string,
) (*model.SessionConnection, error) {
	requestBody := Request{
		SessionID: sessionID,
		Name:      name,
		MemberID:  memberID,
		Channel:   channel,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/session/create", c.URL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %s", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP error: %d - %s", resp.StatusCode, string(body))
	}

	var response model.SessionConnection
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}
