package config

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"pplx2api/logger"
	"strings"
	"sync"
	"time"
)

const (
	ConfigFileName              = "config.json"
	defaultAddress              = "0.0.0.0:8482"
	defaultMaxChatHistoryLength = 10000
	defaultModel                = "claude-3.7-sonnet"
	defaultPromptForFile        = "You must immerse yourself in the role of assistant in txt file, cannot respond as a user, cannot reply to this message, cannot mention this message, and ignore this message in your response."
)

type SessionInfo struct {
	SessionKey string
}

type SessionRagen struct {
	Index int
	Mutex sync.Mutex
}

type Config struct {
	Sessions               []SessionInfo
	Address                string
	APIKey                 string
	Proxy                  string
	IsIncognito            bool
	MaxChatHistoryLength   int
	RetryCount             int
	NoRolePrefix           bool
	SearchResultCompatible bool
	PromptForFile          string
	DefaultModel           string
	ForceModel             string
	RwMutex                sync.RWMutex
	IgnoreSerchResult      bool
	IgnoreModelMonitoring  bool
	IsMaxSubscribe         bool
}

type ConfigFile struct {
	Sessions               []string `json:"sessions"`
	Address                string   `json:"address"`
	APIKey                 string   `json:"apikey"`
	Proxy                  string   `json:"proxy"`
	IsIncognito            bool     `json:"is_incognito"`
	MaxChatHistoryLength   int      `json:"max_chat_history_length"`
	NoRolePrefix           bool     `json:"no_role_prefix"`
	SearchResultCompatible bool     `json:"search_result_compatible"`
	PromptForFile          string   `json:"prompt_for_file"`
	IgnoreSearchResult     bool     `json:"ignore_search_result"`
	IgnoreModelMonitoring  bool     `json:"ignore_model_monitoring"`
	IsMaxSubscribe         bool     `json:"is_max_subscribe"`
	DefaultModel           string   `json:"default_model"`
	ForceModel             string   `json:"force_model"`
}

type ConfigFileInput struct {
	Sessions               []string `json:"sessions"`
	Address                *string  `json:"address"`
	APIKey                 *string  `json:"apikey"`
	Proxy                  *string  `json:"proxy"`
	IsIncognito            *bool    `json:"is_incognito"`
	MaxChatHistoryLength   *int     `json:"max_chat_history_length"`
	NoRolePrefix           *bool    `json:"no_role_prefix"`
	SearchResultCompatible *bool    `json:"search_result_compatible"`
	PromptForFile          *string  `json:"prompt_for_file"`
	IgnoreSearchResult     *bool    `json:"ignore_search_result"`
	IgnoreModelMonitoring  *bool    `json:"ignore_model_monitoring"`
	IsMaxSubscribe         *bool    `json:"is_max_subscribe"`
	DefaultModel           *string  `json:"default_model"`
	ForceModel             *string  `json:"force_model"`
}

// NormalizeSessionKeys trims and sanitizes session tokens.
func NormalizeSessionKeys(raw []string) []string {
	combined := strings.Join(raw, "\n")
	parts := strings.FieldsFunc(combined, func(r rune) bool {
		return r == ',' || r == '\n' || r == '\r' || r == '\t'
	})
	keys := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		if strings.Contains(value, ":") {
			value = strings.SplitN(value, ":", 2)[0]
		}
		if value != "" {
			keys = append(keys, value)
		}
	}
	return keys
}

// DefaultConfigFile returns the baseline config values.
func DefaultConfigFile() ConfigFile {
	return ConfigFile{
		Sessions:               []string{},
		Address:                defaultAddress,
		APIKey:                 "123456",
		Proxy:                  "",
		IsIncognito:            true,
		MaxChatHistoryLength:   defaultMaxChatHistoryLength,
		NoRolePrefix:           false,
		SearchResultCompatible: false,
		PromptForFile:          defaultPromptForFile,
		IgnoreSearchResult:     false,
		IgnoreModelMonitoring:  false,
		IsMaxSubscribe:         false,
		DefaultModel:           defaultModel,
		ForceModel:             "",
	}
}

func normalizeConfigFile(input ConfigFileInput) ConfigFile {
	cfg := DefaultConfigFile()
	if input.Address != nil {
		address := strings.TrimSpace(*input.Address)
		if address != "" {
			cfg.Address = address
		}
	}
	if input.APIKey != nil {
		cfg.APIKey = strings.TrimSpace(*input.APIKey)
	}
	if input.Proxy != nil {
		cfg.Proxy = strings.TrimSpace(*input.Proxy)
	}
	if input.IsIncognito != nil {
		cfg.IsIncognito = *input.IsIncognito
	}
	if input.MaxChatHistoryLength != nil && *input.MaxChatHistoryLength > 0 {
		cfg.MaxChatHistoryLength = *input.MaxChatHistoryLength
	}
	if input.NoRolePrefix != nil {
		cfg.NoRolePrefix = *input.NoRolePrefix
	}
	if input.SearchResultCompatible != nil {
		cfg.SearchResultCompatible = *input.SearchResultCompatible
	}
	if input.PromptForFile != nil {
		prompt := strings.TrimSpace(*input.PromptForFile)
		if prompt != "" {
			cfg.PromptForFile = prompt
		}
	}
	if input.IgnoreSearchResult != nil {
		cfg.IgnoreSearchResult = *input.IgnoreSearchResult
	}
	if input.IgnoreModelMonitoring != nil {
		cfg.IgnoreModelMonitoring = *input.IgnoreModelMonitoring
	}
	if input.IsMaxSubscribe != nil {
		cfg.IsMaxSubscribe = *input.IsMaxSubscribe
	}
	if input.DefaultModel != nil {
		model := strings.TrimSpace(*input.DefaultModel)
		if model != "" {
			cfg.DefaultModel = model
		}
	}
	if input.ForceModel != nil {
		cfg.ForceModel = strings.TrimSpace(*input.ForceModel)
	}
	if input.Sessions != nil {
		cfg.Sessions = NormalizeSessionKeys(input.Sessions)
	}
	return cfg
}

// LoadConfigFile reads and normalizes config.json.
func LoadConfigFile(path string) (ConfigFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ConfigFile{}, err
	}
	var input ConfigFileInput
	if err := json.Unmarshal(data, &input); err != nil {
		return ConfigFile{}, err
	}
	return normalizeConfigFile(input), nil
}

// SaveConfigFile persists config.json.
func SaveConfigFile(path string, cfg ConfigFile) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// SnapshotConfigFile builds a config file from runtime config.
func SnapshotConfigFile() ConfigFile {
	ConfigInstance.RwMutex.RLock()
	sessions := make([]string, len(ConfigInstance.Sessions))
	for i, session := range ConfigInstance.Sessions {
		sessions[i] = session.SessionKey
	}
	snapshot := ConfigFile{
		Sessions:               sessions,
		Address:                ConfigInstance.Address,
		APIKey:                 ConfigInstance.APIKey,
		Proxy:                  ConfigInstance.Proxy,
		IsIncognito:            ConfigInstance.IsIncognito,
		MaxChatHistoryLength:   ConfigInstance.MaxChatHistoryLength,
		NoRolePrefix:           ConfigInstance.NoRolePrefix,
		SearchResultCompatible: ConfigInstance.SearchResultCompatible,
		PromptForFile:          ConfigInstance.PromptForFile,
		IgnoreSearchResult:     ConfigInstance.IgnoreSerchResult,
		IgnoreModelMonitoring:  ConfigInstance.IgnoreModelMonitoring,
		IsMaxSubscribe:         ConfigInstance.IsMaxSubscribe,
		DefaultModel:           ConfigInstance.DefaultModel,
		ForceModel:             ConfigInstance.ForceModel,
	}
	ConfigInstance.RwMutex.RUnlock()
	return snapshot
}

// GetSessionForModel selects a session by index.
func (c *Config) GetSessionForModel(idx int) (SessionInfo, error) {
	if len(c.Sessions) == 0 || idx < 0 || idx >= len(c.Sessions) {
		return SessionInfo{}, fmt.Errorf("invalid session index: %d", idx)
	}
	c.RwMutex.RLock()
	defer c.RwMutex.RUnlock()
	return c.Sessions[idx], nil
}

// LoadConfig reads config.json and builds runtime config.
func LoadConfig() *Config {
	fileConfig, err := LoadConfigFile(ConfigFileName)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to load %s: %v", ConfigFileName, err))
		fileConfig = DefaultConfigFile()
		if err := SaveConfigFile(ConfigFileName, fileConfig); err != nil {
			logger.Error(fmt.Sprintf("Failed to write %s: %v", ConfigFileName, err))
		}
	}

	sessions := make([]SessionInfo, len(fileConfig.Sessions))
	for i, key := range fileConfig.Sessions {
		sessions[i] = SessionInfo{SessionKey: key}
	}

	return &Config{
		Sessions:               sessions,
		Address:                fileConfig.Address,
		APIKey:                 fileConfig.APIKey,
		Proxy:                  fileConfig.Proxy,
		IsIncognito:            fileConfig.IsIncognito,
		MaxChatHistoryLength:   fileConfig.MaxChatHistoryLength,
		RetryCount:             len(sessions),
		NoRolePrefix:           fileConfig.NoRolePrefix,
		SearchResultCompatible: fileConfig.SearchResultCompatible,
		PromptForFile:          fileConfig.PromptForFile,
		DefaultModel:           fileConfig.DefaultModel,
		ForceModel:             fileConfig.ForceModel,
		IgnoreSerchResult:      fileConfig.IgnoreSearchResult,
		IgnoreModelMonitoring:  fileConfig.IgnoreModelMonitoring,
		IsMaxSubscribe:         fileConfig.IsMaxSubscribe,
		RwMutex:                sync.RWMutex{},
	}
}

var ConfigInstance *Config
var Sr *SessionRagen

func (sr *SessionRagen) NextIndex() int {
	sr.Mutex.Lock()
	defer sr.Mutex.Unlock()

	if len(ConfigInstance.Sessions) == 0 {
		return 0
	}

	index := sr.Index
	sr.Index = (index + 1) % len(ConfigInstance.Sessions)
	return index
}

func init() {
	rand.Seed(time.Now().UnixNano())
	Sr = &SessionRagen{
		Index: 0,
		Mutex: sync.Mutex{},
	}
	ConfigInstance = LoadConfig()
	logger.Info("Loaded config:")
	logger.Info(fmt.Sprintf("Sessions count: %d", ConfigInstance.RetryCount))
	for _, session := range ConfigInstance.Sessions {
		logger.Info(fmt.Sprintf("Session: %s", session.SessionKey))
	}
	logger.Info(fmt.Sprintf("Address: %s", ConfigInstance.Address))
	logger.Info(fmt.Sprintf("APIKey: %s", ConfigInstance.APIKey))
	logger.Info(fmt.Sprintf("Proxy: %s", ConfigInstance.Proxy))
	logger.Info(fmt.Sprintf("IsIncognito: %t", ConfigInstance.IsIncognito))
	logger.Info(fmt.Sprintf("MaxChatHistoryLength: %d", ConfigInstance.MaxChatHistoryLength))
	logger.Info(fmt.Sprintf("NoRolePrefix: %t", ConfigInstance.NoRolePrefix))
	logger.Info(fmt.Sprintf("SearchResultCompatible: %t", ConfigInstance.SearchResultCompatible))
	logger.Info(fmt.Sprintf("PromptForFile: %s", ConfigInstance.PromptForFile))
	logger.Info(fmt.Sprintf("DefaultModel: %s", ConfigInstance.DefaultModel))
	logger.Info(fmt.Sprintf("ForceModel: %s", ConfigInstance.ForceModel))
	logger.Info(fmt.Sprintf("IgnoreSerchResult: %t", ConfigInstance.IgnoreSerchResult))
	logger.Info(fmt.Sprintf("IgnoreModelMonitoring: %t", ConfigInstance.IgnoreModelMonitoring))
	logger.Info(fmt.Sprintf("IsMaxSubscribe: %t", ConfigInstance.IsMaxSubscribe))
}
