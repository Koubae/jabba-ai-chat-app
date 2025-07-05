package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/utils"
	"io"
	"log"
	"net/http"
	"time"
)

type AIBotConnector struct {
	AiBotHost           string
	AiBotPort           int
	AiBotTimeoutSeconds int
	AiBotHelloMessage   string
	AiBotURL            string
	httpClient          *http.Client
}

type AIBotRequest struct {
	Message string `json:"message"`
}

type AIBotResponse struct {
	Reply string `json:"reply"`
}

func NewAIBotConnector() *AIBotConnector {
	aiBotHost := utils.GetEnvString("AI_BOT_HOST", "http://localhost")
	aiBotPort := utils.GetEnvInt("AI_BOT_PORT", 20003)
	aiBotTimeoutSeconds := utils.GetEnvInt("AI_BOT_TIMEOUT_SECONDS", 120)
	aiBotHelloMessage := utils.GetEnvString("AI_BOT_HELLO_MESSAGE", "Hello!")

	aiBotURL := fmt.Sprintf("%s:%d", aiBotHost, aiBotPort)

	return &AIBotConnector{
		AiBotHost:           aiBotHost,
		AiBotPort:           aiBotPort,
		AiBotTimeoutSeconds: aiBotTimeoutSeconds,
		AiBotHelloMessage:   aiBotHelloMessage,
		AiBotURL:            aiBotURL,

		httpClient: &http.Client{
			Timeout: time.Duration(aiBotTimeoutSeconds) * time.Second,
		},
	}
}

func (c *AIBotConnector) Hello(ctx context.Context, accessToken string, sessionID string) (*AIBotResponse, error) {
	return c.SendMessage(ctx, accessToken, sessionID, c.AiBotHelloMessage)
}

// SendMessage sends a message to the AI bot and waits for a response
func (c *AIBotConnector) SendMessage(ctx context.Context, accessToken string, sessionID string, message string) (*AIBotResponse, error) {
	requestBody := AIBotRequest{
		Message: message,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/ai/http/bot/send-message/%s", c.AiBotURL, sessionID)

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

	var response AIBotResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}
