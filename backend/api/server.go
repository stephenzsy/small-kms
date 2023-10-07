package api

import (
	ctx "context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/auth"
	certtemplate "github.com/stephenzsy/small-kms/backend/cert-template"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/profile"
)

type server struct {
	common.CommonConfig
	profileService      profile.ProfileService
	certTemplateService certtemplate.CertificateTemplateService
}

func (s *server) ServiceContext(c ctx.Context) common.ServiceContext {
	return common.WithClientProvider(c, s)
}

func wrapResponse[T interface{}](c *gin.Context, defaultStatus int, data T, err error) {
	switch {
	case err == nil:
		c.JSON(defaultStatus, data)
	case errors.Is(err, common.ErrStatusBadRequest):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, common.ErrStatusUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case errors.Is(err, common.ErrStatusForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, common.ErrStatusNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, common.ErrStatusConflict):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		log.Error().Err(err).Stack().Msg("internal error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}

// MsGraphDelegatedClient implements common.ClientProvider.
func (s *server) MsGraphDelegatedClient(c ctx.Context) (*msgraphsdkgo.GraphServiceClient, error) {
	if authIdentity, ok := auth.GetAuthIdentity(c); ok {
		if creds, err := authIdentity.GetOnBehalfOfTokenCredential(s, nil); err != nil {
			return nil, err
		} else {
			return msgraphsdkgo.NewGraphServiceClientWithCredentials(creds, nil)
		}
	}
	return nil, fmt.Errorf("%w: no auth header to authenticate to graph service", common.ErrStatusUnauthorized)
}

func NewServer() models.ServerInterface {
	commonConfig, err := common.NewCommonConfig()
	if err != nil {
		log.Panic().Err(err).Msg("failed to create common config")
	}
	s := server{
		CommonConfig:        &commonConfig,
		profileService:      profile.NewProfileService(),
		certTemplateService: certtemplate.NewCertificateTemplateService(),
	}
	return &s
}