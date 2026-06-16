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
	key    string
}

func NewOpenRouterAdapter(key string) *OpenRouterAdapter {
	return &OpenRouterAdapter{client: &http.Client{}, key: key}
}

func (a *OpenRouterAdapter) StreamCompletion(
	ctx context.Context,
	request models.CompletionRequest,
) (<-chan models.StreamEvent, error) {
	eventCh := make(chan models.StreamEvent)

	reqMessages := make([]openRouterMessage, len(request.ThreadHistory))

	for i, m := range request.ThreadHistory {
		reqMessages[i] = openRouterMessage{Role: m.Role, Content: m.Content}
		if m.ToolCallID != nil {
			reqMessages[i].ToolCallID = *m.ToolCallID
		}
	}

	toolsPreprocessed := make([]requestToolDTO, len(request.Tools))
	for i, t := range request.Tools {
		toolsPreprocessed[i] = preprocessTool(t)
	}

	modelID := fmt.Sprintf("%s/%s", request.Model.Provider, request.Model.Slug)
	payload := openRouterRequestDTO{
		Model:      modelID,
		Messages:   reqMessages,
		Stream:     true,
		ToolChoice: "auto",
		Tools:      toolsPreprocessed,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.key))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/ekosachev/loom")
	req.Header.Set("X-Title", "Loom")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	go func() {
		defer resp.Body.Close()
		defer close(eventCh)

		if resp.StatusCode != http.StatusOK {
			eventCh <- models.StreamEvent{
				Type: models.EventError,
				Err:  err,
			}
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		toolCalls := map[int]*toolCallDTO{}

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
						Content   string        `json:"content"`
						ToolCalls []toolCallDTO `json:"tool_calls"`
					} `json:"delta"`
				} `json:"choices"`
				Usage *struct {
					CompletionTokens int     `json:"completion_tokens"`
					PromptTokens     int     `json:"prompt_tokens"`
					Cost             float64 `json:"cost"`
				} `json:"usage"`
			}

			if err := json.Unmarshal([]byte(dataStr), &chunk); err != nil {
				eventCh <- models.StreamEvent{
					Type: models.EventError,
					Err:  err,
				}
				return
			}

			if usage := chunk.Usage; usage != nil {
				eventCh <- models.StreamEvent{
					Type: models.EventUsageInfo,
					Usage: &models.UsageInfo{
						PromptTokens:     usage.PromptTokens,
						CompletionTokens: usage.CompletionTokens,
						Cost:             usage.Cost,
					},
				}
			}

			if len(chunk.Choices) > 0 {
				content := chunk.Choices[0].Delta.Content
				chunkToolCalls := chunk.Choices[0].Delta.ToolCalls

				for _, call := range chunkToolCalls {
					if _, ok := toolCalls[call.Index]; !ok {
						toolCalls[call.Index] = &call
					} else {
						toolCalls[call.Index].combineChunks(call)
					}
				}

				if content != "" {
					select {
					case eventCh <- models.StreamEvent{
						Type: models.EventText,
						Text: content,
					}:
					case <-ctx.Done():
						eventCh <- models.StreamEvent{Type: models.EventDone}
						break
					}
				}
			}
		}

		for _, call := range toolCalls {
			eventCh <- models.StreamEvent{
				Type: models.EventToolCall,
				ToolCall: &models.ToolCall{
					Name:      call.Function.Name,
					ID:        call.ID,
					Arguments: call.Function.Arguments,
				},
			}
		}
	}()

	return eventCh, nil
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
