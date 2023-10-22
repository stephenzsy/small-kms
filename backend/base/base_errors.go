package base

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
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
	ErrAzCosmosDocNotFound = errors.New("az cosmos doc not found")

	ErrMsGraphResourceNotFound = errors.New("Request_ResourceNotFound")

	ErrResponseStatusBadRequest = httpResposneError(http.StatusBadRequest, "bad request")
	ErrResponseStatusForbidden  = httpResposneError(http.StatusForbidden, "forbidden")
)

func HandleMsGraphError(err error) error {
	if err == nil || errors.Is(err, ErrMsGraphResourceNotFound) {
		return err
	}
	errCode, _, ok := ExtractGraphODataErrorCode(err)
	if ok && errCode != nil && *errCode == "Request_ResourceNotFound" {
		return fmt.Errorf("%w:%w", ErrMsGraphResourceNotFound, err)
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
		return err
	}
}
