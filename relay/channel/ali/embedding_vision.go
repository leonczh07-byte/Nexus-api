package ali

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/types"

	"github.com/gin-gonic/gin"
)

const tongyiEmbeddingVisionPrefix = "tongyi-embedding-vision"

func isTongyiEmbeddingVisionModel(model string) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(model)), tongyiEmbeddingVisionPrefix)
}

type TongyiEmbeddingContent struct {
	Text  string `json:"text,omitempty"`
	Image string `json:"image,omitempty"`
}

type TongyiEmbeddingNativeRequest struct {
	Model string `json:"model"`
	Input struct {
		Contents []TongyiEmbeddingContent `json:"contents"`
	} `json:"input"`
}

type TongyiEmbeddingNativeResponse struct {
	Output struct {
		Embeddings []struct {
			Embedding []float64 `json:"embedding"`
			Index     int       `json:"index"`
		} `json:"embeddings"`
	} `json:"output"`
	Usage struct {
		InputTokens int `json:"input_tokens"`
	} `json:"usage"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

func convertEmbeddingRequest2Tongyi(info *relaycommon.RelayInfo, request dto.EmbeddingRequest) (any, error) {
	texts := make([]string, 0)
	switch v := request.Input.(type) {
	case string:
		if strings.TrimSpace(v) != "" {
			texts = append(texts, v)
		}
	case []any:
		for _, item := range v {
			if s, ok := item.(string); ok && strings.TrimSpace(s) != "" {
				texts = append(texts, s)
			}
		}
	}
	if len(texts) == 0 {
		return nil, errors.New("embedding input text is required")
	}

	native := TongyiEmbeddingNativeRequest{Model: info.UpstreamModelName}
	native.Input.Contents = make([]TongyiEmbeddingContent, 0, len(texts))
	for _, t := range texts {
		native.Input.Contents = append(native.Input.Contents, TongyiEmbeddingContent{Text: t})
	}
	return native, nil
}

func handleTongyiEmbeddingResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError) {
	body, readErr := io.ReadAll(resp.Body)
	resp.Body.Close()
	if readErr != nil {
		return nil, types.NewErrorWithStatusCode(
			errors.New("failed to read tongyi embedding response"),
			types.ErrorCodeReadResponseBodyFailed,
			http.StatusInternalServerError,
		)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, types.NewErrorWithStatusCode(
			fmt.Errorf("tongyi embedding upstream error (status %d): %s", resp.StatusCode, string(body)),
			types.ErrorCodeBadResponse,
			resp.StatusCode,
		)
	}

	var native TongyiEmbeddingNativeResponse
	if unmarshalErr := json.Unmarshal(body, &native); unmarshalErr != nil {
		return nil, types.NewErrorWithStatusCode(
			errors.New("failed to parse tongyi embedding response"),
			types.ErrorCodeBadResponseBody,
			http.StatusInternalServerError,
		)
	}
	if len(native.Output.Embeddings) == 0 {
		msg := native.Message
		if msg == "" {
			msg = "empty embeddings in tongyi response"
		}
		return nil, types.NewErrorWithStatusCode(
			errors.New(msg),
			types.ErrorCodeBadResponse,
			http.StatusBadGateway,
		)
	}

	items := make([]dto.EmbeddingResponseItem, 0, len(native.Output.Embeddings))
	for i, e := range native.Output.Embeddings {
		items = append(items, dto.EmbeddingResponseItem{
			Object:    "embedding",
			Index:     i,
			Embedding: e.Embedding,
		})
	}
	promptTokens := native.Usage.InputTokens
	oaiResp := dto.EmbeddingResponse{
		Object: "list",
		Data:   items,
		Model:  info.OriginModelName,
	}
	oaiResp.Usage = dto.Usage{
		PromptTokens:     promptTokens,
		CompletionTokens: 0,
		TotalTokens:      promptTokens,
	}
	c.JSON(http.StatusOK, oaiResp)

	usage = &dto.Usage{
		PromptTokens:     promptTokens,
		CompletionTokens: 0,
		TotalTokens:      promptTokens,
	}
	return usage, nil
}
