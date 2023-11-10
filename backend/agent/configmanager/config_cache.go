package configmanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/stephenzsy/small-kms/backend/utils"
)

type ConfigCache[T VersionedConfig] struct {
	configName    string
	configBaseDir string
	config        utils.Nullable[T]
}

func (cache *ConfigCache[T]) Load() (bool, error) {
	activeLink := filepath.Join(cache.configBaseDir, fmt.Sprintf("%s.active", cache.configName))
	if _, err := os.Stat(activeLink); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	configFilename := filepath.Join(activeLink, "config.json")
	if err := readJsonFile(configFilename, &cache.config); err != nil {
		return false, err
	}
	return true, nil
}

func (cache *ConfigCache[T]) SetPulledConfig(pulled *T) error {
	oldVersion := cache.config.Value().GetVersion()
	cache.config.SetValue(*pulled)
	if oldVersion != cache.config.Value().GetVersion() {
		return cache.Persist(true)
	}
	return nil
}

func (cache *ConfigCache[T]) Persist(asActive bool) error {
	if !cache.config.HasValue() {
		return nil
	}
	versionedDir := filepath.Join(cache.configBaseDir, "versioned", fmt.Sprintf("%s.%s", cache.configName, cache.config.Value().GetVersion()))
	if _, err := os.Stat(versionedDir); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err := os.MkdirAll(versionedDir, 0750); err != nil {
			return err
		}
	}
	configFilename := filepath.Join(versionedDir, "config.json")
	if err := writeJsonFile(configFilename, cache.config); err != nil {
		return err
	}
	if asActive {
		activeLink := filepath.Join(cache.configBaseDir, fmt.Sprintf("%s.active", cache.configName))
		if _, err := os.Lstat(activeLink); err == nil {
			if err := os.Remove(activeLink); err != nil {
				return err
			}
		}
		if versionedRelPath, err := filepath.Rel(cache.configBaseDir, versionedDir); err != nil {
			return err
		} else if err := os.Symlink(versionedRelPath, activeLink); err != nil {
			return err
		}
	}
	return nil
}

func readJsonFile[T any](filename string, v T) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, v)
}

func writeJsonFile[T any](filename string, v T) error {
	content, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, content, 0640)
}

func NewConfigCache[T VersionedConfig](name, configBaseDir string) *ConfigCache[T] {
	return &ConfigCache[T]{
		configName:    name,
		configBaseDir: configBaseDir,
	}
}
