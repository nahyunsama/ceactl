package llmanalysis

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	Endpoint string
	Model    string
	HTTP     *http.Client
}

func NewClient(endpoint, model string) *Client {
	return &Client{
		Endpoint: endpoint,
		Model:    model,
		HTTP: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// Chat sends systemPrompt/userPrompt to the Ollama /api/chat endpoint
// (POST {model, messages:[system,user], stream:false}) and returns the
// assistant's reply content.
func (c *Client) Chat(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	// TODO: build request body, POST to c.Endpoint+"/api/chat", parse
	// {message: {content: string}} from the response.
	return "", fmt.Errorf("Chat: not implemented")
}
