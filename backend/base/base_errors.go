package base

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/labstack/echo/v4"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/rs/zerolog/log"
)

type HttpResponseError struct {
	error
	StatusCode int
}

func httpResposneError(statusCode int, message string) error {
	return &HttpResponseError{
		error:      errors.New(message),
		StatusCode: statusCode,
	}
}

var (
	ErrAzCosmosDocNotFound    = errors.New("az cosmos doc not found")
	ErrAzKeyVaultItemNotFound = errors.New("az key vault key not found")

	ErrResponseStatusBadRequest      = httpResposneError(http.StatusBadRequest, "bad request")
	ErrResposneStatusUnauthorized    = httpResposneError(http.StatusUnauthorized, "unauthorized")
	ErrResponseStatusForbidden       = httpResposneError(http.StatusForbidden, "forbidden")
	ErrResponseStatusNotFound        = httpResposneError(http.StatusNotFound, "not found")
	ErrResponseStatusTooManyRequests = httpResposneError(http.StatusTooManyRequests, "conflict")
)

func HandleAzCosmosError(err error) error {
	if err == nil || errors.Is(err, ErrAzCosmosDocNotFound) {
		return err
	}
	var respError *azcore.ResponseError
	if errors.As(err, &respError) {
		if respError.StatusCode == http.StatusNotFound {
			return fmt.Errorf("%w:%w", ErrAzCosmosDocNotFound, err)
		}
	}
	return err
}

func HandleAzKeyVaultError(err error) error {
	if err == nil || errors.Is(err, ErrAzKeyVaultItemNotFound) {
		return err
	}
	var respError *azcore.ResponseError
	if errors.As(err, &respError) {
		if respError.StatusCode == http.StatusNotFound {
			return fmt.Errorf("%w:%w", ErrAzKeyVaultItemNotFound, err)
		}
	}
	return err
}

func ExtractGraphODataErrorCode(err error) (errorCode *string, odErr *odataerrors.ODataError, ok bool) {
	ok = errors.As(err, &odErr)
	if ok {
		errorCode = odErr.GetErrorEscaped().GetCode()
	}
	return
}

func HandleResponseError(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return err
		}
		var respErr *HttpResponseError
		if errors.As(err, &respErr) {
			return c.JSON(respErr.StatusCode, map[string]string{"message": err.Error()})
		}
		log.Error().Err(err).Msg("internal error")
		return err
	}
}
