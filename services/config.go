package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"portus/constants"
	"portus/models"
	"portus/repository"
	"portus/utils"
	"strings"
	"sync"

	"github.com/knadh/koanf/parsers/dotenv"
	kjson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// ConfigService provides methods to interact with configuration
type ConfigService interface {
	InitConfig(ctx context.Context) error
	GetConfig() *models.Configuration
	SaveConfig(cfg models.Configuration) error
	GetFileConfig() *models.Configuration
	SaveFileConfig(cfg models.Configuration) error
	ResetFileConfig(ctx context.Context) error
}

type configService struct {
	configRepo repository.ConfigRepository
	config     *models.Configuration
	configLock sync.RWMutex
	k          *koanf.Koanf
	configPath string
}

// NewConfigService creates a new configuration service
func NewConfigService(configRepo repository.ConfigRepository) ConfigService {
	return &configService{
		configRepo: configRepo,
		configPath: "config/app.config.json",
	}
}

// InitConfig initializes the configuration
func (s *configService) InitConfig(ctx context.Context) error {
	logger := utils.LoggerFromContext(ctx)
	logger.Info().Msg("Initilizing Config")
	s.k = koanf.New(".")

	// 1. Load defaults
	if err := s.k.Load(confmap.Provider(constants.DefaultConfig, "."), nil); err != nil {
		return fmt.Errorf("error loading defaults: %w", err)
	}

	// Ensure config directory exists
	if err := s.configRepo.EnsureConfigDir(); err != nil {
		return err
	}

	// 2. Load app.config.json
	f := file.Provider(s.configPath)

	if err := s.k.Load(f, kjson.Parser()); err != nil {
		// Create default config if file doesn't exist
		if os.IsNotExist(err) {
			defaultConfig := &models.Configuration{}
			if err := s.k.Unmarshal("", defaultConfig); err != nil {
				return fmt.Errorf("error unmarshaling default config: %w", err)
			}
			if err := s.configRepo.WriteConfigFile(defaultConfig); err != nil {
				return fmt.Errorf("error saving default config: %w", err)
			}
		} else {
			return fmt.Errorf("error loading config file: %w", err)
		}
	}

	// Set up file watcher
	s.configRepo.WatchConfigFile(func() {
		s.configLock.Lock()
		defer s.configLock.Unlock()

		// Create a new instance and reload
		s.k = koanf.New(".")

		// Reload in the correct order
		s.k.Load(confmap.Provider(constants.DefaultConfig, "."), nil)
		s.k.Load(file.Provider(s.configPath), kjson.Parser())
		s.k.Load(file.Provider(".env"), dotenv.Parser())
		s.k.Load(env.Provider("PORTUS_", ".", s.envKeyReplacer), nil)

		// Update the config struct
		if err := s.k.Unmarshal("", s.config); err != nil {
			fmt.Printf("error unmarshaling config: %v\n", err)
			return
		}

		fmt.Println("Configuration reloaded due to file change")
	})

	// 3. Load environment variables
	if err := s.k.Load(env.Provider("PORTUS_", ".", s.envKeyReplacer), nil); err != nil {
		return fmt.Errorf("error loading environment variables: %w", err)
	}

	// Load the final config
	s.configLock.Lock()
	defer s.configLock.Unlock()

	s.config = &models.Configuration{}
	if err := s.k.UnmarshalWithConf("", s.config, koanf.UnmarshalConf{
		Tag: "json",
	}); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}

	return nil
}

// Helper method for environment variable key conversion
func (s *configService) envKeyReplacer(key string) string {
	return strings.ReplaceAll(
		strings.ToLower(
			strings.TrimPrefix(key, "PORTUS_")),
		"_",
		".",
	)
}

// GetConfig returns the current configuration
func (s *configService) GetConfig() *models.Configuration {
	s.configLock.RLock()
	defer s.configLock.RUnlock()
	return s.config
}

// SaveConfig saves and updates the configuration
func (s *configService) SaveConfig(cfg models.Configuration) error {
	s.configLock.Lock()
	defer s.configLock.Unlock()

	// Convert config struct to map
	jsonBytes, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	var configMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &configMap); err != nil {
		return fmt.Errorf("error unmarshaling to map: %w", err)
	}

	// Load the new config
	if err := s.k.Load(confmap.Provider(configMap, "."), nil); err != nil {
		return fmt.Errorf("error loading new config: %w", err)
	}

	// Save to file
	if err := s.configRepo.WriteConfigFile(&cfg); err != nil {
		return fmt.Errorf("error saving config: %w", err)
	}

	s.config = &cfg
	return nil
}

// GetFileConfig returns only the file-based configuration
func (s *configService) GetFileConfig() *models.Configuration {
	cfg, err := s.configRepo.ReadConfigFile()
	if err != nil {
		return nil
	}
	return cfg
}

// SaveFileConfig saves the configuration to file only
func (s *configService) SaveFileConfig(cfg models.Configuration) error {
	return s.configRepo.WriteConfigFile(&cfg)
}

// ResetFileConfig resets config file to defaults
func (s *configService) ResetFileConfig(ctx context.Context) error {
	// Create default config
	k := koanf.New(".")
	if err := k.Load(confmap.Provider(constants.DefaultConfig, "."), nil); err != nil {
		return fmt.Errorf("error loading defaults: %w", err)
	}

	defaultConfig := &models.Configuration{}
	if err := k.Unmarshal("", defaultConfig); err != nil {
		return fmt.Errorf("error unmarshaling default config: %w", err)
	}

	// Save defaults to file
	if err := s.configRepo.WriteConfigFile(defaultConfig); err != nil {
		return fmt.Errorf("error writing default config: %w", err)
	}

	// Reload the main configuration
	return s.InitConfig(ctx)
}
