package plugins

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"guardian/internal/models"
	"guardian/prompt_api"
	"net/http"
)

var (
	ErrForwardRequest       = errors.New("failed to forward request")
	ErrPluginResponseFailed = errors.New("failed to receive a response")
)

type PluginClient interface {
	Forward(ctx context.Context, reqBody *models.PluginRequest) (*models.PluginResponse, error)
}

type HTTPClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
	Forward(ctx context.Context, reqBody *models.PluginRequest) (*models.PluginResponse, error)
}

type HTTPClient struct {
	*http.Client
}

type GRPCClient struct {
	Client prompt_api.PromptServiceClient
}

func NewHTTPClient(client *http.Client) *HTTPClient {
	return &HTTPClient{client}
}

func NewPluginGRPCClient(client *grpc.ClientConn) *GRPCClient {
	return &GRPCClient{
		Client: prompt_api.NewPromptServiceClient(client),
	}
}

func (g *GRPCClient) Forward(ctx context.Context, reqBody *models.PluginRequest) (*models.PluginResponse, error) {
	req := prompt_api.SendPromptRequest{
		Prompt:   reqBody.Prompt,
		Chat:     reqBody.Chat,
		UserID:   reqBody.UserID.String(),
		TargetID: reqBody.TargetID.String(),
	}

	resp, err := g.Client.SendPrompt(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrForwardRequest, err)
	}
	var score uint32
	if respScore, ok := resp.GetOptionalScore().(*prompt_api.SendPromptResponse_Score); ok {
		score = respScore.Score
	}

	return &models.PluginResponse{
		Status: resp.GetStatus(),
		Score:  score,
	}, nil
}

func (h *HTTPClient) Forward(ctx context.Context, reqBody *models.PluginRequest) (*models.PluginResponse, error) {
	marshalledBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal the request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqBody.Address, bytes.NewBuffer(marshalledBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrForwardRequest, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w from: %s", ErrPluginResponseFailed, reqBody.Address)
	}

	var sendResponse models.PluginResponse
	if err := json.NewDecoder(resp.Body).Decode(&sendResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &models.PluginResponse{
		Status: sendResponse.Status,
		Score:  sendResponse.Score,
	}, nil
}
