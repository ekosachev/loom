package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/ekosachev/loom/internal/domain/models"
)

type OpenRouterAdapter struct {
	client *http.Client
}

func NewOpenRouterAdapter() *OpenRouterAdapter {
	return &OpenRouterAdapter{client: &http.Client{}}
}

type openRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openRouterRequest struct {
	Model    string              `json:"model"`
	Messages []openRouterMessage `json:"messages"`
	Stream   bool                `json:"stream"`
}

func (a *OpenRouterAdapter) StreamCompletion(
	ctx context.Context,
	cfg *models.Config,
	modelID string,
	history []models.Message,
) (<-chan string, <-chan error, error) {
	outCh := make(chan string)
	errCh := make(chan error)

	reqMessages := make([]openRouterMessage, len(history))

	for i, m := range history {
		reqMessages[i] = openRouterMessage{Role: m.Role, Content: m.Content}
	}

	payload := openRouterRequest{
		Model:    modelID,
		Messages: reqMessages,
		Stream:   true,
	}

	jsonData, err := json.Marshal(payload)

	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.Openrouter.Key))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/ekosachev/loom")
	req.Header.Set("X-Title", "Loom")

	resp, err := a.client.Do(req)

	if err != nil {
		return nil, nil, err
	}

	go func() {
		defer resp.Body.Close()

		defer close(outCh)
		defer close(errCh)

		if resp.StatusCode != http.StatusOK {
			errCh <- fmt.Errorf("bad status code: %d", resp.StatusCode)
			return
		}

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			dataStr := strings.TrimPrefix(line, "data: ")
			if dataStr == "[DONE]" {
				break
			}

			var chunk struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
				} `json:"choices"`
			}

			if err := json.Unmarshal([]byte(dataStr), &chunk); err != nil {
				errCh <- err
				return
			}

			if len(chunk.Choices) > 0 {
				content := chunk.Choices[0].Delta.Content
				if content != "" {
					select {
					case outCh <- content:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return outCh, errCh, nil
}

func (a *OpenRouterAdapter) GetMetaForModel(ctx context.Context, cfg models.Config, provider string, slug string) (*models.Model, error) {
	url := fmt.Sprintf("https://openrouter.ai/api/v1/models/%s/%s/endpoints", provider, slug)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.Openrouter.Key))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/ekosachev/loom")
	req.Header.Set("X-Title", "Loom")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("OpenRouterAPI returned an error status code %s", resp.Status)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	type responseSchema struct {
		Data struct {
			Name      string `json:"name"`
			Endpoints []struct {
				SupportedParameters []string `json:"supported_parameters"`
			} `json:"endpoints"`
		} `json:"data"`
	}

	var responseData responseSchema
	if err := json.Unmarshal(respBody, &responseData); err != nil {
		return nil, err
	}

	supportedParams := responseData.Data.Endpoints[0].SupportedParameters
	supportsTools := slices.Contains(supportedParams[:], "tools")

	return &models.Model{
		Name:          responseData.Data.Name,
		Provider:      provider,
		Slug:          slug,
		SupportsTools: supportsTools,
	}, nil
}
