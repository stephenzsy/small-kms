package agentconfigmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	agentclient "github.com/stephenzsy/small-kms/backend/agent/client/v2"
	agentutils "github.com/stephenzsy/small-kms/backend/agent/utils"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	keymodels "github.com/stephenzsy/small-kms/backend/models/key"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func enrollCert(c context.Context, cm ConfigManager, certPolicyID string) (*certmodels.Certificate, string, error) {

	var enrolledFileName string
	cert, _, err := agentutils.EnrollCertificate(c, cm.Client(), certPolicyID,
		func(cert *certmodels.Certificate) (*os.File, error) {
			enrolledFileName = cm.ConfigDir().Certs().File(fmt.Sprintf("%s.pem", cert.ID))
			return cm.ConfigDir().Certs().OpenFile(fmt.Sprintf("%s.pem", cert.ID), os.O_CREATE|os.O_WRONLY, 0400)
		}, false)
	return cert, enrolledFileName, err
}

func pullPublicJWK(c context.Context, cm ConfigManager, keyID string) (*keymodels.Key, error) {
	filename := cm.ConfigDir().JWKs().File(fmt.Sprintf("%s.json", keyID))
	if _, err := os.Stat(filename); err == nil {
		if content, err := os.ReadFile(filename); err == nil {
			var key keymodels.Key
			if err := json.Unmarshal(content, &key); err == nil {
				return &key, nil
			}
		}
	}
	resp, err := cm.Client().GetKeyWithResponse(c, models.NamespaceProviderServicePrincipal, "me", keyID, &agentclient.GetKeyParams{
		IncludeJwk: utils.ToPtr(true),
	})
	if err != nil {
		return nil, err
	} else if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	key := resp.JSON200
	if keyJson, err := json.Marshal(key); err != nil {
		return key, err
	} else if err := os.WriteFile(filename, keyJson, 0400); err != nil {
		return key, err
	}
	return key, nil
}
