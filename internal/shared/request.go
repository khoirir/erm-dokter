package shared

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"erm-dokter/internal/shared/apperror"
)

const (
	DefaultMaxBodyBytes = 1024 * 1024
	MaxLogPayloadBytes  = 500
)

var reSensitiveData = regexp.MustCompile(`(?i)("(?:password|pin|token|secret|api_?key)"\s*:\s*)("[^"]*"?|[^,\}\]]*)`)

func MaskSensitiveData(raw string) string {
	return reSensitiveData.ReplaceAllString(raw, `${1}"***"`)
}

func DecodeJSON(w http.ResponseWriter, r *http.Request, target any, customClientMsg ...string) bool {
	return DecodeJSONWithLimit(w, r, target, DefaultMaxBodyBytes, customClientMsg...)
}
func DecodeJSONWithLimit(w http.ResponseWriter, r *http.Request, target any, maxBytes int64, customClientMsg ...string) bool {
	clientMsg := "Format data yang dikirim tidak valid"
	if len(customClientMsg) > 0 && strings.TrimSpace(customClientMsg[0]) != "" {
		clientMsg = customClientMsg[0]
	}

	bodyBytes, err := io.ReadAll(io.LimitReader(r.Body, maxBytes))
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError(
			clientMsg,
			fmt.Sprintf("Gagal membaca body request: %v", err),
		))
		return false
	}

	trimmed := strings.TrimSpace(string(bodyBytes))
	if len(trimmed) == 0 {
		apperror.HandleError(w, apperror.NewBusinessError(
			clientMsg,
			"Payload request kosong (0 bytes)",
		))
		return false
	}

	if err := json.Unmarshal(bodyBytes, target); err != nil {
		masked := MaskSensitiveData(trimmed)
		if len(masked) > MaxLogPayloadBytes {
			masked = masked[:MaxLogPayloadBytes] + "... (dipotong)"
		}

		apperror.HandleError(w, apperror.NewBusinessError(
			clientMsg,
			fmt.Sprintf("Gagal decode payload JSON: %v (payload: %s)", err, masked),
		))
		return false
	}

	return true
}
