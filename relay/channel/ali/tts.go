package ali

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/types"

	"github.com/gin-gonic/gin"
)

const aliTTSContextKeyCharCount = "ali_tts_char_count"

type AliTTSInput struct {
	Text  string `json:"text"`
	Voice string `json:"voice,omitempty"`
}

type AliTTSParameters struct {
	Format string  `json:"format,omitempty"`
	Speed  float64 `json:"speed,omitempty"`
}

type AliTTSRequest struct {
	Model      string            `json:"model"`
	Input      AliTTSInput       `json:"input"`
	Parameters *AliTTSParameters `json:"parameters,omitempty"`
}

type AliTTSNativeResponse struct {
	Output struct {
		Audio struct {
			URL string `json:"url"`
		} `json:"audio"`
	} `json:"output"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

func convertAudioRequest2AliTTS(c *gin.Context, info *relaycommon.RelayInfo, request dto.AudioRequest) (io.Reader, error) {
	if strings.TrimSpace(request.Input) == "" {
		return nil, errors.New("tts input text is required")
	}
	voice := strings.TrimSpace(request.Voice)
	if voice == "" {
		voice = "Cherry"
	}
	aliReq := AliTTSRequest{
		Model: info.UpstreamModelName,
		Input: AliTTSInput{Text: request.Input, Voice: voice},
	}
	params := &AliTTSParameters{}
	if request.ResponseFormat != "" {
		params.Format = request.ResponseFormat
	}
	if request.Speed != nil {
		params.Speed = *request.Speed
	}
	if params.Format != "" || params.Speed != 0 {
		aliReq.Parameters = params
	}
	c.Set(aliTTSContextKeyCharCount, utf8.RuneCountInString(request.Input))
	body, err := json.Marshal(aliReq)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(body), nil
}

func audioContentTypeFromURL(audioURL string) string {
	lower := strings.ToLower(audioURL)
	switch {
	case strings.Contains(lower, ".mp3"):
		return "audio/mpeg"
	case strings.Contains(lower, ".wav"):
		return "audio/wav"
	case strings.Contains(lower, ".aac"):
		return "audio/aac"
	case strings.Contains(lower, ".ogg"), strings.Contains(lower, ".opus"):
		return "audio/ogg"
	case strings.Contains(lower, ".flac"):
		return "audio/flac"
	}
	return "audio/wav"
}

func handleAliTTSResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError) {
	body, readErr := io.ReadAll(resp.Body)
	resp.Body.Close()
	if readErr != nil {
		return nil, types.NewErrorWithStatusCode(
			errors.New("failed to read ali tts response"),
			types.ErrorCodeReadResponseBodyFailed,
			http.StatusInternalServerError,
		)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, types.NewErrorWithStatusCode(
			fmt.Errorf("ali tts upstream error (status %d): %s", resp.StatusCode, string(body)),
			types.ErrorCodeBadResponse,
			resp.StatusCode,
		)
	}

	var native AliTTSNativeResponse
	if unmarshalErr := json.Unmarshal(body, &native); unmarshalErr != nil {
		return nil, types.NewErrorWithStatusCode(
			errors.New("failed to parse ali tts response"),
			types.ErrorCodeBadResponseBody,
			http.StatusInternalServerError,
		)
	}
	if native.Output.Audio.URL == "" {
		msg := native.Message
		if msg == "" {
			msg = "empty audio url in ali tts response"
		}
		return nil, types.NewErrorWithStatusCode(
			errors.New(msg),
			types.ErrorCodeBadResponse,
			http.StatusBadGateway,
		)
	}

	client := &http.Client{Timeout: 120 * time.Second}
	audioResp, fetchErr := client.Get(native.Output.Audio.URL)
	if fetchErr != nil {
		return nil, types.NewErrorWithStatusCode(
			fmt.Errorf("failed to fetch ali tts audio: %v", fetchErr),
			types.ErrorCodeBadResponse,
			http.StatusBadGateway,
		)
	}
	defer audioResp.Body.Close()
	if audioResp.StatusCode != http.StatusOK {
		return nil, types.NewErrorWithStatusCode(
			fmt.Errorf("failed to fetch ali tts audio, oss status: %d", audioResp.StatusCode),
			types.ErrorCodeBadResponse,
			http.StatusBadGateway,
		)
	}
	audioData, readAudioErr := io.ReadAll(audioResp.Body)
	if readAudioErr != nil {
		return nil, types.NewErrorWithStatusCode(
			errors.New("failed to read ali tts audio bytes"),
			types.ErrorCodeReadResponseBodyFailed,
			http.StatusBadGateway,
		)
	}

	contentType := audioContentTypeFromURL(native.Output.Audio.URL)
	c.Data(http.StatusOK, contentType, audioData)

	charCount := info.GetEstimatePromptTokens()
	if v, exists := c.Get(aliTTSContextKeyCharCount); exists {
		if n, ok := v.(int); ok && n > 0 {
			charCount = n
		}
	}
	usage = &dto.Usage{
		PromptTokens:     charCount,
		CompletionTokens: 0,
		TotalTokens:      charCount,
	}
	return usage, nil
}
