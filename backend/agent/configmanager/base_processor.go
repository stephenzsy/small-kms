package cm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type baseConfigProcessor struct {
	*sharedConfig
	configName     shared.AgentConfigName
	isActive       bool
	pollShutdownCh chan struct{}
	processPending bool
	timer          *time.Timer
	configCtx      ConfigCtx
}

// Name implements ConfigProcessor.
func (p *baseConfigProcessor) ConfigName() shared.AgentConfigName {
	return p.configName
}

func loadFetchedConfig(filename string) (*shared.AgentConfiguration, error) {
	configJson, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config := shared.AgentConfiguration{}
	err = json.Unmarshal(configJson, &config)
	return &config, err
}

func (p *baseConfigProcessor) loadFetchedConfig(ctx context.Context, isPending bool) {
	var filename string
	if isPending {
		filename = filepath.Join(p.configDir, fmt.Sprintf("%s.next", string(p.configName)), "config.json")
	} else {
		filename = filepath.Join(p.configDir, string(p.configName), "config.json")
	}

	logger := log.Ctx(ctx).Info()
	if activeConfig, err := loadFetchedConfig(filename); err != nil {
		logger := logger.Err(err)
		if isPending {
			logger.Msgf("load fetched config - pending: %s", p.configName)
		} else {
			logger.Msgf("load fetched config - active: %s", p.configName)
		}
	} else {
		slotKey := fetchConfigActiveSlotKey
		if isPending {
			slotKey = fetchConfigPendingSlotKey
		}
		p.configCtx = withConfigSlot(p.configCtx, slotKey, configSlot[shared.AgentConfiguration]{
			config:  activeConfig,
			exp:     *activeConfig.NextRefreshAfter,
			version: activeConfig.Version,
		})
	}
}
func (p *baseConfigProcessor) persistVersionedConfig(ctx context.Context, config *shared.AgentConfiguration, overwriteIfExist bool) error {
	version := config.Version
	versionedPathPart := fmt.Sprintf("%s.%s", p.configName, version)
	versionedDir := filepath.Join(p.versionedConfigDir, versionedPathPart)
	if _, err := os.Stat(versionedDir); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err = os.MkdirAll(versionedDir, 0700); err != nil {
			return err
		}
	}
	configFilename := filepath.Join(versionedDir, "config.json")
	var persisted = false
	if _, err := os.Stat(configFilename); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	} else {
		persisted = true
	}
	if !persisted || overwriteIfExist {
		configJson, err := json.Marshal(config)
		if err != nil {
			return err
		}
		return os.WriteFile(configFilename, configJson, 0600)
	}
	return nil
}

func (p *baseConfigProcessor) symlinkFetchedConfig(ctx context.Context, version string, isPending bool) error {
	versionedPathPart := fmt.Sprintf("%s.%s", p.configName, version)
	var linkName string
	if isPending {
		linkName = filepath.Join(p.configDir, fmt.Sprintf("%s.%s", string(p.configName), "next"))
	} else {
		linkName = filepath.Join(p.configDir, string(p.configName))
	}
	relTarget := filepath.Join(".", "versioned", versionedPathPart)
	log.Ctx(ctx).Info().Msgf("symlink fetched config: %s -> %s", linkName, relTarget)
	return os.Symlink(relTarget, linkName)
}

func (p *baseConfigProcessor) baseStart(c context.Context, scheduleToUpdate chan<- pollConfigMsg,
	exitCh chan<- error,
	onTimer func() *pollConfigMsg) func() {
	p.timer = time.NewTimer(0)
	p.isActive = true
	for p.isActive {
		select {
		case <-p.pollShutdownCh:
			p.isActive = false
		case <-c.Done():
			exitCh <- fmt.Errorf("%w:processor:%s:%w", ErrForcedShutdown, p.ConfigName(), c.Err())
			return func() {}
		case <-p.timer.C:
			if !p.processPending && p.isActive {
				p.processPending = true
				msg := onTimer()
				if msg != nil {
					scheduleToUpdate <- *msg
				}
			}
		}
	}
	return func() {
		exitCh <- fmt.Errorf("%w:processor:%s:%w", ErrGracefullyShutdown, p.ConfigName(), c.Err())
	}
}

func (p *baseConfigProcessor) baseShutdown() func() {
	p.isActive = false
	p.timer.Stop()
	return func() {
		p.pollShutdownCh <- struct{}{}
	}
}

func (p *baseConfigProcessor) Shutdown() {
	cb := p.baseShutdown()
	cb()
}

func (p *baseConfigProcessor) baseMarkProcessDone() func(time.Duration) {
	p.processPending = true
	p.timer.Stop()
	return func(d time.Duration) {
		log.Info().Msgf("%s: next refresh after: %s", p.configName, d)
		p.processPending = false
		p.timer.Reset(d)
	}
}
