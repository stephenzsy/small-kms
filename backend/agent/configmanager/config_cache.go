package configmanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type ConfigCache[T VersionedConfig] struct {
	configName    string
	configBaseDir string
	config        *T
}

func (cache *ConfigCache[T]) Load() (r T, ok bool, err error) {
	activeLink := filepath.Join(cache.configBaseDir, fmt.Sprintf("%s.active", cache.configName))
	if _, err = os.Stat(activeLink); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return r, false, nil
		}
		return r, false, err
	}
	configFilename := filepath.Join(activeLink, "config.json")
	if err = readJsonFile(configFilename, &r); err != nil {
		return r, false, err
	}
	cache.config = &r
	return r, true, nil
}

func (cache *ConfigCache[T]) SetPulledConfig(pulled T) error {
	old := cache.config
	cache.config = &pulled
	if old == nil || (*old).GetVersion() != pulled.GetVersion() {
		return cache.Persist(true)
	}
	return nil
}

func (cache *ConfigCache[T]) Persist(asActive bool) error {
	if cache.config == nil {
		return nil
	}
	versionedDir := filepath.Join(cache.configBaseDir, "versioned", fmt.Sprintf("%s.%s", cache.configName, (*cache.config).GetVersion()))
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
