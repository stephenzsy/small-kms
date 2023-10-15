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
	fetchedConfig  *shared.AgentConfiguration
	pollShutdownCh chan struct{}
	processPending bool
}

// Name implements ConfigProcessor.
func (p *baseConfigProcessor) ConfigName() shared.AgentConfigName {
	return p.configName
}

func (p *baseConfigProcessor) loadFetchedConfig() error {
	configFilename := filepath.Join(p.configDir, string(p.configName), "config.json")
	configJson, err := os.ReadFile(configFilename)
	if err != nil {
		return err
	}
	config := shared.AgentConfiguration{}
	if err := json.Unmarshal(configJson, &config); err != nil {
		return err
	}
	p.fetchedConfig = &config
	return nil
}

func (p *baseConfigProcessor) saveNewFetchedVersion(config *shared.AgentConfiguration) error {
	version := config.Version
	versionedPathPart := fmt.Sprintf("%s.%s", p.configName, version)
	vdir := filepath.Join(p.versionedConfigDir, versionedPathPart)
	if err := os.MkdirAll(vdir, 0700); err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	vdirLink := filepath.Join(p.configDir, fmt.Sprintf("%s.next", p.configName))
	if _, err := os.Lstat(vdirLink); err == nil {
		if err = os.Remove(vdirLink); err != nil {
			return err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	linkTargetPath := filepath.Join(".", "versioned", versionedPathPart)
	log.Info().Msgf("symlink %s -> %s", vdirLink, linkTargetPath)
	if err := os.Symlink(linkTargetPath, vdirLink); err != nil {
		return err
	}
	configFilename := filepath.Join(vdir, "config.json")
	configJson, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(configFilename, configJson, 0600)
}

func (p *baseConfigProcessor) baseStart(c context.Context, scheduleToUpdate chan<- pollConfigMsg, exitCh chan<- error, tickerCh <-chan time.Time,
	onTicker func() *pollConfigMsg) {
	p.isActive = true
	for p.isActive {
		select {
		case <-p.pollShutdownCh:
			goto pollEnd
		case <-c.Done():
			exitCh <- fmt.Errorf("%w:processor:%s:%w", ErrForcedShutdown, p.ConfigName(), c.Err())
			return
		case <-tickerCh:
			if !p.processPending && p.isActive {
				p.processPending = true
				msg := onTicker()
				if msg != nil {
					scheduleToUpdate <- *msg
				}
			}
		}
	}
pollEnd:
	exitCh <- fmt.Errorf("%w:processor:%s:%w", ErrGracefullyShutdown, p.ConfigName(), c.Err())
}

func (p *baseConfigProcessor) baseShutdown() func() {
	p.isActive = false
	return func() {
		p.pollShutdownCh <- struct{}{}
	}
}

func (p *baseConfigProcessor) baseMarkProcessDone() func() {
	p.processPending = true
	return func() {
		p.processPending = false
	}
}
