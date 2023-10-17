package cm

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type configSlot[TR any] struct {
	configFetched shared.AgentConfiguration
	configActive  TR
	exp           time.Time
	version       string
}

type ConfigCtx[TR any] struct {
	pendingSlot configSlot[TR]
	activeSlot  configSlot[TR]
}

const (
	activeSymlinkSuffix   = ".active"
	pendingSymlinkSuffix  = ".pending"
	activeConfigFilename  = "config.active.json"
	fetchedConfigFilename = "config.json"
)

func (s *configSlot[TR]) loadConfigFromFiles(configName shared.AgentConfigName, configDir string, isPending bool) error {
	dirSymlinkSuffix := activeSymlinkSuffix
	if isPending {
		dirSymlinkSuffix = pendingSymlinkSuffix
	}

	fnFetched := filepath.Join(configDir, fmt.Sprintf("%s%s", string(configName), dirSymlinkSuffix), fetchedConfigFilename)
	cbFetched, err := os.ReadFile(fnFetched) // config bytes
	if err != nil {
		return err
	}
	err = json.Unmarshal(cbFetched, &s.configFetched)
	if err != nil {
		return err
	}
	s.exp = *s.configFetched.NextRefreshAfter
	s.version = s.configFetched.Version

	if !isPending {
		fnActive := filepath.Join(configDir, fmt.Sprintf("%s%s", string(configName), dirSymlinkSuffix), activeConfigFilename)
		cbActive, err := os.ReadFile(fnActive) // config bytes
		if err != nil {
			return err
		}
		err = json.Unmarshal(cbActive, &s.configActive)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *configSlot[TR]) hasValue() bool {
	return s.version != ""
}

func (s *configSlot[TR]) persistConfig(configDir string, configName shared.AgentConfigName, isPending bool, overwriteIfExist bool) error {
	if s.version == "" {
		return nil
	}
	versionedPathPart := fmt.Sprintf("%s.%s", configName, s.version)
	versionedDir := filepath.Join(configDir, "versioned", versionedPathPart)
	if _, err := os.Stat(versionedDir); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err = os.MkdirAll(versionedDir, 0700); err != nil {
			return err
		}
	}
	fnFetched := filepath.Join(versionedDir, fetchedConfigFilename)
	if overwriteIfExist || !utils.FileExists(fnFetched) {
		if b, err := json.Marshal(s.configFetched); err != nil {
			return err
		} else if err = os.WriteFile(fnFetched, b, 0600); err != nil {
			return err
		}
	}
	if !isPending {
		fnActive := filepath.Join(versionedDir, activeConfigFilename)
		if overwriteIfExist || !utils.FileExists(fnActive) {
			if b, err := json.Marshal(s.configActive); err != nil {
				return err
			} else if err = os.WriteFile(fnActive, b, 0600); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *configSlot[TR]) persistSymlink(configDir string, configName shared.AgentConfigName, linkSuffix string) error {
	linkName := filepath.Join(configDir, fmt.Sprintf("%s%s", configName, linkSuffix))
	// remove symlink
	if utils.SymlinkExists(linkName) {
		log.Debug().Msgf("remove symlink: %s, %s", linkName, s.version)
		if err := os.Remove(linkName); err != nil {
			return err
		}
	}
	if s.version != "" {
		// create symlink
		versionedRelDir := filepath.Join(".", "versioned", fmt.Sprintf("%s.%s", configName, s.version))
		log.Debug().Msgf("create symlink: %s -> %s", linkName, versionedRelDir)
		return os.Symlink(versionedRelDir, linkName)
	}
	return nil
}

func (c *ConfigCtx[TR]) persistConfig(
	configDir string,
	configName shared.AgentConfigName,
	overwriteIfExist bool) error {
	if err := c.activeSlot.persistConfig(configDir, configName, false, overwriteIfExist); err != nil {
		return err
	}
	if err := c.pendingSlot.persistConfig(configDir, configName, true, overwriteIfExist); err != nil {
		return err
	}
	return nil
}

func (c *ConfigCtx[TR]) persistSymlinks(
	configDir string,
	configName shared.AgentConfigName) error {
	if err := c.activeSlot.persistSymlink(configDir, configName, activeSymlinkSuffix); err != nil {
		return err
	}
	if err := c.pendingSlot.persistSymlink(configDir, configName, pendingSymlinkSuffix); err != nil {
		return err
	}
	return nil
}

func (c *ConfigCtx[TR]) getLatestSlot() *configSlot[TR] {
	if c.pendingSlot.version != "" {
		return &c.pendingSlot
	}
	return &c.activeSlot
}

func (c *ConfigCtx[TR]) setActiveConfig(config TR, basePendingSlot *configSlot[TR]) {
	c.activeSlot.configActive = config
	c.activeSlot.version = basePendingSlot.version
	c.activeSlot.exp = basePendingSlot.exp
	c.activeSlot.configFetched = basePendingSlot.configFetched
	if basePendingSlot.version == c.pendingSlot.version {
		// reset pending slot
		c.pendingSlot = configSlot[TR]{}
	}
}

func (c *ConfigCtx[TR]) getWaitForNextRefresh() time.Duration {
	hasSetDuration := false
	var duration time.Duration
	if c.pendingSlot.hasValue() {
		duration = time.Until(c.pendingSlot.exp)
		hasSetDuration = true
	}
	if c.activeSlot.hasValue() {
		if !hasSetDuration {
			duration = time.Until(c.activeSlot.exp)
			hasSetDuration = true
		} else if d := time.Until(c.activeSlot.exp); d < duration {
			duration = d
		}
	}
	if duration < minRemoteDuration {
		duration = minRemoteDuration
	}
	return duration
}
