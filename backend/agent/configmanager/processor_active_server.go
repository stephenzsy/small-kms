package cm

import (
	"context"
	"encoding/pem"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type activeServerProcessor struct {
	baseConfigProcessor
	timer *time.Timer
}

type processorState struct {
	hasAttemptedLoad bool
	fetchExpire      time.Time
	fetchExtraExpire time.Time
}

// Version implements ConfigProcessor.
func (*activeServerProcessor) Version() string {
	return "0.1"
}

// Name implements ConfigProcessor.
func (*activeServerProcessor) Name() shared.AgentConfigName {
	return shared.AgentConfigNameHeartbeat
}

const errorVersion = "deadbeef"
const defuaultErrorBackoff = 5 * time.Minute

func (p *activeServerProcessor) Process(ctx context.Context, task string) error {
	logger := log.Ctx(ctx)
	switch task {
	case TaskNameLoad:
		p.loadFetchedConfig(ctx, false)
		p.loadFetchedConfig(ctx, true)
		return nil
	case TaskNameFetch:
		resp, err := p.AgentClient().GetAgentConfigurationWithResponse(ctx,
			shared.NamespaceKindServicePrincipal, meNamespaceIdIdentifier,
			shared.AgentConfigNameActiveServer,
			nil)
		if err != nil {
			return err
		}
		nextConfig := resp.JSON200
		logger.Info().Msgf("fetched config: %s:%s", p.configName, nextConfig.Version)
		err = p.persistVersionedConfig(ctx, nextConfig, false)
		if err != nil {
			// not critial
			logger.Error().Err(err).Msgf("failed persist versioned config: %s:%s", p.configName, nextConfig.Version)
		}
		if existingSlot, hasExistingSlot := getFetchConfigSlotPreferPending[shared.AgentConfiguration](p.configCtx); !hasExistingSlot || existingSlot.version != nextConfig.Version {
			// version changed, replace pending slot
			p.configCtx = withConfigSlot(p.configCtx, fetchConfigPendingSlotKey, configSlot[shared.AgentConfiguration]{
				config:  nextConfig,
				exp:     *nextConfig.NextRefreshAfter,
				version: nextConfig.Version,
			})
		} else {
			// update expiring time
			existingSlot.exp = *nextConfig.NextRefreshAfter
		}
		return nil
	case TaskNameReady:
		fetchedPendingSlot, ok := getConfigSlot[shared.AgentConfiguration](p.configCtx, fetchConfigPendingSlotKey)
		if !ok {
			// no pending slot
			log.Error().Msgf("no pending fetched config slot to process: %s", p.configName)
			return nil
		}
		activeServerCfg, err := fetchedPendingSlot.config.Config.AsAgentConfigurationAgentActiveServer()
		if err != nil {
			return err
		}
		// fetch certificate
		serverCertResp, err := p.AgentClient().GetCertificateWithResponse(ctx, shared.NamespaceKindServicePrincipal, meNamespaceIdIdentifier, *activeServerCfg.ServerCertificateId, nil)
		if err != nil {
			return err
		}
		if err := p.processKeyVaultSecret(ctx, serverCertResp.JSON200); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unknown task: %s", task)
	}
	return nil
}

func (p *activeServerProcessor) processKeyVaultSecret(ctx context.Context, certInfo *shared.CertificateInfo) error {
	kid := azsecrets.ID(*certInfo.Jwk.KeyID)
	certDir := filepath.Join(p.configDir, "certs", string(p.configName), kid.Name(), kid.Version())
	if _, err := os.Stat(certDir); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err = os.MkdirAll(certDir, 0700); err != nil {
			return err
		}
	}
	keyFilename := filepath.Join(certDir, "key.pem")
	certFilename := filepath.Join(certDir, "cert.pem")
	bundleFilename := filepath.Join(certDir, "bundle.pem")

	needFetch := false
	if _, err := os.Stat(keyFilename); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		needFetch = true
	} else if _, err := os.Stat(certFilename); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		needFetch = true
	} else if _, err := os.Stat(bundleFilename); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		needFetch = true
	}
	if needFetch {
		secretResp, err := p.AzSecretesClient().GetSecret(ctx, kid.Name(), kid.Version(), nil)
		if err != nil {
			return err
		}
		pemBytes := []byte(*secretResp.SecretBundle.Value)
		block, rest := pem.Decode(pemBytes)
		if block.Type != "CERTIFICATE" {
			// this is the private key
			if err := os.WriteFile(keyFilename, pem.EncodeToMemory(block), 0400); err != nil {
				return err
			}
			if err := os.WriteFile(certFilename, rest, 0600); err != nil {
				return err
			}
		}
		if err := os.WriteFile(bundleFilename, pemBytes, 0400); err != nil {
			return err
		}
	}
	return nil
}

func (p *activeServerProcessor) Start(c context.Context, scheduleToUpdate chan<- pollConfigMsg, exitCh chan<- error) {
	p.configCtx = context.Background()
	finally := p.baseStart(c, scheduleToUpdate, exitCh, func() *pollConfigMsg {
		if !configCtxHasAttemptedLoad(p.configCtx) {
			// attempt load if has not attempted
			p.configCtx = configCtxWithAttemptedLoad(p.configCtx)
			return &pollConfigMsg{
				name:      shared.AgentConfigNameActiveServer,
				task:      TaskNameLoad,
				processor: p,
			}
		}
		if fetchSlot, hasFetchSlot := getFetchConfigSlotPreferPending[shared.AgentConfiguration](p.configCtx); !hasFetchSlot || fetchSlot.exp.Before(time.Now()) {
			// no fetch slot or expired, go for fetch
			return &pollConfigMsg{
				name:      shared.AgentConfigNameActiveServer,
				task:      TaskNameFetch,
				processor: p,
			}
		}

		// TODO
		if _, hasReadyPendingSlot := getConfigSlot[shared.AgentConfiguration](p.configCtx, readyConfigPendingSlotKey); !hasReadyPendingSlot {
			return &pollConfigMsg{
				name:      shared.AgentConfigNameActiveServer,
				task:      TaskNameReady,
				processor: p,
			}
		}

		return nil
	})
	if slot, hasSlot := getConfigSlot[shared.AgentConfiguration](p.configCtx, fetchConfigPendingSlotKey); hasSlot {
		p.persistVersionedConfig(c, slot.config, true)
		p.symlinkFetchedConfig(c, slot.version, true)
	}
	finally()
}

var _ ConfigProcessor = (*activeServerProcessor)(nil)

const minRemoteDuration = 5 * time.Minute

func applyJitter(d time.Duration) time.Duration {
	if d <= 5*time.Minute {
		return minRemoteDuration + time.Duration(rand.Int63n(int64(5*time.Second)))
	}
	if d <= 1*time.Hour {
		return minRemoteDuration + time.Duration(rand.Int63n(int64(60*time.Second)))
	}
	return minRemoteDuration + time.Duration(rand.Int63n(int64(5*time.Minute)))
}

func (p *activeServerProcessor) MarkProcessDone(taskName string, err error) {
	resetTimer := p.baseMarkProcessDone()
	nextDuration := 5 * time.Minute

	if fetchSlot, hasFetchSlot := getFetchConfigSlotPreferPending[shared.AgentConfiguration](p.configCtx); !hasFetchSlot || fetchSlot.exp.Before(time.Now()) {
		// no fetch slot or expired, go for fetch
		resetTimer(5 * time.Second)
		return
	}

	// TODO
	if _, hasReadyPendingSlot := getConfigSlot[shared.AgentConfiguration](p.configCtx, readyConfigPendingSlotKey); !hasReadyPendingSlot {
		resetTimer(5 * time.Second)
		return
	}

	resetTimer(nextDuration)
}

func newActiveServerProcessor(sc *sharedConfig) *activeServerProcessor {
	return &activeServerProcessor{
		baseConfigProcessor: baseConfigProcessor{
			sharedConfig:   sc,
			configName:     shared.AgentConfigNameActiveServer,
			pollShutdownCh: make(chan struct{}, 1),
		},
	}
}
