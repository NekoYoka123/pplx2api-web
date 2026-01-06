package service

import (
	"errors"
	"net/http"
	"strings"

	"pplx2api/config"

	"github.com/gin-gonic/gin"
)

type AdminConfigResponse struct {
	Address                string   `json:"address"`
	APIKeySet              bool     `json:"api_key_set"`
	APIKeyHint             string   `json:"api_key_hint"`
	Proxy                  string   `json:"proxy"`
	IsIncognito            bool     `json:"is_incognito"`
	MaxChatHistoryLength   int      `json:"max_chat_history_length"`
	RetryCount             int      `json:"retry_count"`
	NoRolePrefix           bool     `json:"no_role_prefix"`
	SearchResultCompatible bool     `json:"search_result_compatible"`
	PromptForFile          string   `json:"prompt_for_file"`
	IgnoreSearchResult     bool     `json:"ignore_search_result"`
	IgnoreModelMonitoring  bool     `json:"ignore_model_monitoring"`
	IsMaxSubscribe         bool     `json:"is_max_subscribe"`
	RejectModelMismatch    bool     `json:"reject_model_mismatch"`
	DefaultModel           string   `json:"default_model"`
	ForceModel             string   `json:"force_model"`
	Sessions               []string `json:"sessions"`
	SessionsCount          int      `json:"sessions_count"`
}

type AdminConfigUpdateRequest struct {
	APIKey                 *string  `json:"apikey,omitempty"`
	Proxy                  *string  `json:"proxy,omitempty"`
	IsIncognito            *bool    `json:"is_incognito,omitempty"`
	MaxChatHistoryLength   *int     `json:"max_chat_history_length,omitempty"`
	NoRolePrefix           *bool    `json:"no_role_prefix,omitempty"`
	SearchResultCompatible *bool    `json:"search_result_compatible,omitempty"`
	PromptForFile          *string  `json:"prompt_for_file,omitempty"`
	IgnoreSearchResult     *bool    `json:"ignore_search_result,omitempty"`
	IgnoreModelMonitoring  *bool    `json:"ignore_model_monitoring,omitempty"`
	IsMaxSubscribe         *bool    `json:"is_max_subscribe,omitempty"`
	DefaultModel           *string  `json:"default_model,omitempty"`
	ForceModel             *string  `json:"force_model,omitempty"`
	Sessions               *[]string `json:"sessions,omitempty"`
	RejectModelMismatch    *bool    `json:"reject_model_mismatch,omitempty"`
}

type AdminConfigUpdateResponse struct {
	Status  string   `json:"status"`
	Changed []string `json:"changed"`
}

func AdminConfigGetHandler(c *gin.Context) {
	config.ConfigInstance.RwMutex.RLock()
	sessions := make([]string, 0, len(config.ConfigInstance.Sessions))
	for _, session := range config.ConfigInstance.Sessions {
		sessions = append(sessions, session.SessionKey)
	}
	retryCount := config.ConfigInstance.RetryCount
	config.ConfigInstance.RwMutex.RUnlock()

	resp := AdminConfigResponse{
		Address:                config.ConfigInstance.Address,
		APIKeySet:              config.ConfigInstance.APIKey != "",
		APIKeyHint:             maskAPIKey(config.ConfigInstance.APIKey),
		Proxy:                  config.ConfigInstance.Proxy,
		IsIncognito:            config.ConfigInstance.IsIncognito,
		MaxChatHistoryLength:   config.ConfigInstance.MaxChatHistoryLength,
		RetryCount:             retryCount,
		NoRolePrefix:           config.ConfigInstance.NoRolePrefix,
		SearchResultCompatible: config.ConfigInstance.SearchResultCompatible,
		PromptForFile:          config.ConfigInstance.PromptForFile,
		IgnoreSearchResult:     config.ConfigInstance.IgnoreSerchResult,
		IgnoreModelMonitoring:  config.ConfigInstance.IgnoreModelMonitoring,
		IsMaxSubscribe:         config.ConfigInstance.IsMaxSubscribe,
		RejectModelMismatch:    config.ConfigInstance.RejectModelMismatch,
		DefaultModel:           config.ConfigInstance.DefaultModel,
		ForceModel:             config.ConfigInstance.ForceModel,
		Sessions:               sessions,
		SessionsCount:          len(sessions),
	}

	c.JSON(http.StatusOK, resp)
}

func AdminConfigUpdateHandler(c *gin.Context) {
	var req AdminConfigUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	changed := []string{}

	if req.APIKey != nil {
		apiKey := strings.TrimSpace(*req.APIKey)
		if apiKey == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "API key cannot be empty"})
			return
		}
		config.ConfigInstance.APIKey = apiKey
		changed = append(changed, "apikey")
	}

	if req.Proxy != nil {
		config.ConfigInstance.Proxy = strings.TrimSpace(*req.Proxy)
		changed = append(changed, "proxy")
	}

	if req.IsIncognito != nil {
		config.ConfigInstance.IsIncognito = *req.IsIncognito
		changed = append(changed, "is_incognito")
	}

	if req.MaxChatHistoryLength != nil {
		if *req.MaxChatHistoryLength < 1 {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "max_chat_history_length must be >= 1"})
			return
		}
		config.ConfigInstance.MaxChatHistoryLength = *req.MaxChatHistoryLength
		changed = append(changed, "max_chat_history_length")
	}

	if req.NoRolePrefix != nil {
		config.ConfigInstance.NoRolePrefix = *req.NoRolePrefix
		changed = append(changed, "no_role_prefix")
	}

	if req.SearchResultCompatible != nil {
		config.ConfigInstance.SearchResultCompatible = *req.SearchResultCompatible
		changed = append(changed, "search_result_compatible")
	}

	if req.PromptForFile != nil {
		config.ConfigInstance.PromptForFile = *req.PromptForFile
		changed = append(changed, "prompt_for_file")
	}

	if req.IgnoreSearchResult != nil {
		config.ConfigInstance.IgnoreSerchResult = *req.IgnoreSearchResult
		changed = append(changed, "ignore_search_result")
	}

	if req.IgnoreModelMonitoring != nil {
		config.ConfigInstance.IgnoreModelMonitoring = *req.IgnoreModelMonitoring
		changed = append(changed, "ignore_model_monitoring")
	}

	if req.RejectModelMismatch != nil {
		config.ConfigInstance.RejectModelMismatch = *req.RejectModelMismatch
		changed = append(changed, "reject_model_mismatch")
	}

	if req.IsMaxSubscribe != nil {
		config.ConfigInstance.IsMaxSubscribe = *req.IsMaxSubscribe
		changed = append(changed, "is_max_subscribe")
	}

	if req.DefaultModel != nil {
		defaultModel := strings.TrimSpace(*req.DefaultModel)
		if defaultModel == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "default_model cannot be empty"})
			return
		}
		config.ConfigInstance.DefaultModel = defaultModel
		changed = append(changed, "default_model")
	}

	if req.ForceModel != nil {
		forceModel := strings.TrimSpace(*req.ForceModel)
		config.ConfigInstance.ForceModel = forceModel
		changed = append(changed, "force_model")
	}

	if req.Sessions != nil {
		if err := updateSessions(*req.Sessions); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		changed = append(changed, "sessions")
	}

	if len(changed) > 0 {
		if err := config.SaveConfigFile(config.ConfigFileName, config.SnapshotConfigFile()); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to persist config"})
			return
		}
	}

	c.JSON(http.StatusOK, AdminConfigUpdateResponse{
		Status:  "ok",
		Changed: changed,
	})
}

func updateSessions(raw []string) error {
	keys := config.NormalizeSessionKeys(raw)
	if len(keys) == 0 {
		return errors.New("sessions cannot be empty")
	}

	sessions := make([]config.SessionInfo, len(keys))
	for i, key := range keys {
		sessions[i] = config.SessionInfo{SessionKey: key}
	}

	config.ConfigInstance.RwMutex.Lock()
	config.ConfigInstance.Sessions = sessions
	config.ConfigInstance.RetryCount = len(sessions)
	config.ConfigInstance.RwMutex.Unlock()

	config.Sr.Mutex.Lock()
	config.Sr.Index = 0
	config.Sr.Mutex.Unlock()

	return nil
}

func maskAPIKey(key string) string {
	if key == "" {
		return ""
	}
	if len(key) <= 4 {
		return "****"
	}
	return "****" + key[len(key)-4:]
}
