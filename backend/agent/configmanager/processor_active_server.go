package cm

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type activeServerProcessor struct {
	baseConfigProcessor
	configCtx     ConfigCtx[ActiveServerReadyConfig]
	serverReadyCh chan<- ActiveServerReadyConfig
	listener      string
	slotID        uint32
}

// Version implements ConfigProcessor.
func (*activeServerProcessor) Version() string {
	return "0.1"
}

// Name implements ConfigProcessor.
func (*activeServerProcessor) Name() shared.AgentConfigName {
	return shared.AgentConfigNameHeartbeat
}

type ActiveServerReadyConfig struct {
	ServerCertificateFile                         string   `json:"serverCertificateFile"`
	AuthorizedClientCertificateFingerprintsSHA384 [][]byte `json:"authorizedClientCertificateFingerprintsSHA384"`
}

func (p *activeServerProcessor) Process(ctx context.Context, task string) error {
	logger := log.Ctx(ctx)
	switch task {
	case TaskNameLoad:
		if err := p.configCtx.activeSlot.loadConfigFromFiles(p.configName, p.configDir, false); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return err
			}
		}
		if err := p.configCtx.pendingSlot.loadConfigFromFiles(p.configName, p.configDir, true); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return err
			}
		}
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

		latestSlot := p.configCtx.getLatestSlot()
		if latestSlot.version == nextConfig.Version {
			latestSlot.exp = *nextConfig.NextRefreshAfter
		} else {
			p.configCtx.pendingSlot.configFetched = *nextConfig
			p.configCtx.pendingSlot.version = nextConfig.Version
			p.configCtx.pendingSlot.exp = *nextConfig.NextRefreshAfter
			err = p.configCtx.pendingSlot.persistConfig(p.configDir, p.configName, true, true)
			if err != nil {
				// not critial to block next step
				logger.Error().Err(err).Msgf("failed persist versioned config: %s:%s", p.configName, nextConfig.Version)
			}
		}
		return nil
	case TaskNameActivate:
		pslot := &p.configCtx.pendingSlot
		if !pslot.hasValue() {
			return fmt.Errorf("no pending fetched config slot to process: %s", p.configName)
		}

		// activeServerCfg, err := pslot.configFetched.Config.AsAgentConfigurationAgentActiveServer()
		// if err != nil {
		// 	return err
		// }

		// fetch server certificate with private key
		// bundleFilename, err := p.processCertificate(ctx, *activeServerCfg.ServerCertificateId, true)
		// if err != nil {
		// 	return err
		// }

		readyConfig := ActiveServerReadyConfig{
			// ServerCertificateFile:                         bundleFilename,
			// AuthorizedClientCertificateFingerprintsSHA384: make([][]byte, 0, len(activeServerCfg.AuthorizedCertificateIds)),
		}
		// fetch authrorized client certificates
		// for _, clientCertId := range activeServerCfg.AuthorizedCertificateIds {
		// 	if rawCertFilename, err := p.processCertificate(ctx, clientCertId, false); err != nil {
		// 		return err
		// 	} else {
		// 		hash, err := getCertFingprintSHA384FromFile(rawCertFilename)
		// 		if err != nil {
		// 			return err
		// 		}
		// 		readyConfig.AuthorizedClientCertificateFingerprintsSHA384 =
		// 			append(readyConfig.AuthorizedClientCertificateFingerprintsSHA384, hash[:])
		// 	}
		// }

		p.configCtx.setActiveConfig(readyConfig, pslot)
		if err := p.configCtx.persistConfig(p.configDir, p.configName, true); err != nil {
			// not critical to block next step
			logger.Error().Err(err).Msgf("failed to persist config: %s", p.configName)
		}
		// not critical to block next step
		if err := p.configCtx.persistSymlinks(p.configDir, p.configName); err != nil {
			logger.Error().Err(err).Msgf("failed to persist symlinks: %s", p.configName)
		}
		return nil
	case TaskNameConfirm:
		params := shared.AgentConfigurationParameters{}
		// params.FromAgentConfigurationAgentActiveServer(shared.AgentConfigurationAgentActiveServer{
		// 	Name: shared.AgentConfigNameActiveServer,
		// 	Reply: &shared.AgentConfigurationAgentActiveServerReply{
		// 		Listener: p.listener,
		// 		SlotId:   p.slotID,
		// 		State:    shared.AgentConfigurationAgentActiveServerReplyStateUp,
		// 	},
		// })
		_, err := p.AgentClient().AgentCallbackWithResponse(ctx,
			shared.NamespaceKindServicePrincipal, meNamespaceIdIdentifier,
			shared.AgentConfigNameActiveServer,
			shared.AgentConfiguration{
				Config:  params,
				Version: p.configCtx.activeSlot.version,
			})
		logger.Info().Err(err).Msgf("Agent server callback with listener:%s,slot:%d,state:up", p.listener, p.slotID)
		return err
	default:
		return fmt.Errorf("unknown task: %s", task)
	}
}

// returns bytes of the first (chain) certificate in the chain
// func (p *activeServerProcessor) processCertificate(ctx context.Context, certID shared.Identifier, requirePrivateKey bool) (string, error) {
// 	bad := func(e error) (string, error) {
// 		return "", e
// 	}
// 	certDir := filepath.Join(p.configDir, "certs", string(p.configName), certID.String())
// 	needFetch := false
// 	certFilename := filepath.Join(certDir, "cert.pem")
// 	keyFilename := filepath.Join(certDir, "key.pem")
// 	bundleFilename := filepath.Join(certDir, "bundle.pem")

// 	var returnFilename = certFilename
// 	if requirePrivateKey {
// 		returnFilename = bundleFilename
// 	}

// 	if _, err := os.Stat(certFilename); err != nil {
// 		if !errors.Is(err, os.ErrNotExist) {
// 			return bad(err)
// 		}
// 		needFetch = true
// 	}
// 	if requirePrivateKey && !needFetch {
// 		if _, err := os.Stat(keyFilename); err != nil {
// 			if !errors.Is(err, os.ErrNotExist) {
// 				return bad(err)
// 			}
// 			needFetch = true
// 		} else if _, err := os.Stat(bundleFilename); err != nil {
// 			if !errors.Is(err, os.ErrNotExist) {
// 				return bad(err)
// 			}
// 			needFetch = true
// 		}
// 	}

// 	if needFetch {
// 		if _, err := os.Stat(certDir); err != nil {
// 			if !errors.Is(err, os.ErrNotExist) {
// 				return bad(err)
// 			}
// 			if err = os.MkdirAll(certDir, 0700); err != nil {
// 				return bad(err)
// 			}
// 		}
// if requirePrivateKey {
// certResp, err := p.AgentClient().GetCertificateWithResponse(ctx,
// 	shared.NamespaceKindServicePrincipal, meNamespaceIdIdentifier,
// 	certID, nil)
// if err != nil {
// 	return bad(err)
// }
// kid := azsecrets.ID(*certResp.JSON200.Jwk.KeyID)

// secretResp, err := p.AzSecretesClient().GetSecret(ctx, kid.Name(), kid.Version(), nil)
// if err != nil {
// 	return bad(err)
// }
// pemBytes := []byte(*secretResp.Value)
// block, rest := pem.Decode(pemBytes)
// if block.Type != "CERTIFICATE" {
// 	// this is the private key
// 	if err := os.WriteFile(keyFilename, pem.EncodeToMemory(block), 0400); err != nil {
// 		return bad(err)
// 	}
// 	if err := os.WriteFile(certFilename, rest, 0600); err != nil {
// 		return bad(err)
// 	}
// }
// if err := os.WriteFile(bundleFilename, pemBytes, 0400); err != nil {
// 	return bad(err)
// }
// } else {
// certResp, err := p.AgentClient().GetCertificateWithResponse(ctx,
// 	shared.NamespaceKindServicePrincipal, meNamespaceIdIdentifier,
// 	certID, &agentclient.GetCertificateParams{
// 		IncludeCertificate: utils.ToPtr(true),
// 	})
// if err != nil {
// 	return bad(err)
// }
// pemBytes := []byte(*certResp.JSON200.Pem)
// if err := os.WriteFile(certFilename, pemBytes, 0600); err != nil {
// 	return bad(err)
// }
// }
// 	}

// 	return returnFilename, nil
// }

func (p *activeServerProcessor) Start(c context.Context, scheduleToUpdate chan<- pollConfigMsg, shutdownNotifier common.LeafShutdownNotifier) {

	p.baseStart(c, scheduleToUpdate, func() *pollConfigMsg {
		if !p.attemptedLoad {
			p.attemptedLoad = true
			return &pollConfigMsg{
				name:      shared.AgentConfigNameActiveServer,
				task:      TaskNameLoad,
				processor: p,
			}
		}
		lslot := p.configCtx.getLatestSlot()
		if !lslot.hasValue() || lslot.exp.Before(time.Now()) {
			// no config at all, go for fetch
			return &pollConfigMsg{
				name:      shared.AgentConfigNameActiveServer,
				task:      TaskNameFetch,
				processor: p,
			}
		}
		pslot := &p.configCtx.pendingSlot
		if pslot.hasValue() {
			// we have a pending slot needs to be activated
			return &pollConfigMsg{
				name:      shared.AgentConfigNameActiveServer,
				task:      TaskNameActivate,
				processor: p,
			}
		}
		aslot := &p.configCtx.activeSlot
		if aslot.hasValue() {
			p.serverReadyCh <- aslot.configActive
			return &pollConfigMsg{
				name:      shared.AgentConfigNameActiveServer,
				task:      TaskNameConfirm,
				processor: p,
			}
		}
		return nil
	}, func() {
		// persist config before shutdown
		p.configCtx.persistConfig(p.configDir, p.configName, true)
		p.configCtx.persistSymlinks(p.configDir, p.configName)

		params := shared.AgentConfigurationParameters{}
		// params.FromAgentConfigurationAgentActiveServer(shared.AgentConfigurationAgentActiveServer{
		// 	Name: shared.AgentConfigNameActiveServer,
		// 	Reply: &shared.AgentConfigurationAgentActiveServerReply{
		// 		Listener: p.listener,
		// 		SlotId:   p.slotID,
		// 		State:    shared.AgentConfigurationAgentActiveServerReplyStateDown,
		// 	},
		// })
		_, err := p.AgentClient().AgentCallbackWithResponse(c,
			shared.NamespaceKindServicePrincipal, meNamespaceIdIdentifier,
			shared.AgentConfigNameActiveServer,
			shared.AgentConfiguration{
				Config:  params,
				Version: p.configCtx.activeSlot.version,
			})
		if err != nil {
			log.Ctx(c).Error().Err(err).Msgf("Agent server callback with listener:%s,slot:%d,state:down", p.listener, p.slotID)
		} else {
			log.Ctx(c).Info().Msgf("Agent server callback with listener:%s,slot:%d,state:down", p.listener, p.slotID)
		}
	}, shutdownNotifier)
}

var _ ConfigProcessor = (*activeServerProcessor)(nil)

const minRemoteDuration = 5 * time.Minute

func (p *activeServerProcessor) MarkProcessDone(taskName string, err error) {
	resetTimer := p.baseMarkProcessDone()

	switch taskName {
	case TaskNameLoad:
		resetTimer(3*time.Second, taskName)
		return
	case TaskNameFetch:
		if err == nil {
			// fetch succeeded, proceed if has pending slot
			if p.configCtx.pendingSlot.hasValue() {
				resetTimer(3*time.Second, taskName)
				return
			}
		} else {
			resetTimer(minRemoteDuration, taskName)
			return
		}
	case TaskNameActivate:
		if err == nil {
			// activate succeeded, proceed with execute
			resetTimer(3*time.Second, taskName)
		} else {
			resetTimer(minRemoteDuration, taskName)
		}
		return
	case TaskNameConfirm:
		if err != nil {
			resetTimer(minRemoteDuration, taskName)
			return
		}
	}

	resetTimer(p.configCtx.getWaitForNextRefresh(), taskName)
}

func newActiveServerProcessor(sc *sharedConfig, serverReadyCh chan<- ActiveServerReadyConfig, listener string, slotID uint32) *activeServerProcessor {
	return &activeServerProcessor{
		baseConfigProcessor: baseConfigProcessor{
			sharedConfig: sc,
			configName:   shared.AgentConfigNameActiveServer,
		},
		serverReadyCh: serverReadyCh,
		listener:      listener,
		slotID:        slotID,
	}
}
