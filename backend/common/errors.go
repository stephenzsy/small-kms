package common

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/gin-gonic/gin"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/rs/zerolog/log"
)

var (
	ErrStatus2xxCreated   = errors.New("created")            // 201
	ErrStatusBadRequest   = errors.New("invalid input")      // 400
	ErrStatusUnauthorized = errors.New("unauthorized")       // 401
	ErrStatusForbidden    = errors.New("forbidden")          // 403
	ErrStatusNotFound     = errors.New("resource not found") // 404
	ErrStatusConflict     = errors.New("resource conflict")  // 409
)

var (
	ErrMissingEnvVar = errors.New("missing env var")
)

// Deprecated: use WrapAzNotFoundErr instead
func IsAzNotFound(err error) bool {
	var respErr *azcore.ResponseError
	if errors.As(err, &respErr) {
		// Handle Error
		if respErr.StatusCode == http.StatusNotFound {
			return true
		}
	}
	return false
}

func IsAzCosmosNotFound(err error) (error, bool) {
	if err == nil {
		return nil, false
	}
	var respErr *azcore.ResponseError
	if errors.As(err, &respErr) && respErr.StatusCode == http.StatusNotFound {
		return respErr, true
	}
	return err, false
}

func WrapAzRsNotFoundErr(err error, resourceDescriptor string) error {
	if err == nil || errors.Is(err, ErrStatusNotFound) {
		return err
	}
	var respErr *azcore.ResponseError
	if errors.As(err, &respErr) && respErr.StatusCode == http.StatusNotFound {
		return fmt.Errorf("%w: az %s, %w", ErrStatusNotFound, resourceDescriptor, err)
	}
	return err
}

// Deprecated: use WrapMsGraphNotFoundErr instead
func IsGraphODataErrorNotFound(err error) bool {
	var odErr *odataerrors.ODataError
	if errors.As(err, &odErr) {
		if odErr.ResponseStatusCode == http.StatusNotFound {
			return true
		}
	}
	return false
}

func ExtractGraphODataErrorCode(err error) (errorCode *string, odErr *odataerrors.ODataError, ok bool) {
	ok = errors.As(err, &odErr)
	if ok {
		errorCode = odErr.GetErrorEscaped().GetCode()
	}
	return
}

func WrapMsGraphNotFoundErr(err error, resourceDescriptor string) error {
	if err == nil || errors.Is(err, ErrStatusNotFound) {
		return err
	}
	errCode, _, ok := ExtractGraphODataErrorCode(err)
	if ok && errCode != nil && *errCode == "Request_ResourceNotFound" {
		return fmt.Errorf("%w: graph %s, %w", ErrStatusNotFound, resourceDescriptor, err)
	}
	return err
}

func RespondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrStatusBadRequest):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	case errors.Is(err, ErrStatusUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case errors.Is(err, ErrStatusForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, ErrStatusNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, ErrStatusConflict):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		log.Error().Err(err).Stack().Msg("internal error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}
