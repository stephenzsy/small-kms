package graph

import (
	"errors"
	"fmt"

	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

var (
	ErrMsGraphResourceNotFound = errors.New("Request_ResourceNotFound")
)

func HandleMsGraphError(err error) error {
	if err == nil || errors.Is(err, ErrMsGraphResourceNotFound) {
		return err
	}
	errCode, _, ok := extractGraphODataErrorCode(err)
	if ok && errCode != nil && *errCode == "Request_ResourceNotFound" {
		return fmt.Errorf("%w:%w", ErrMsGraphResourceNotFound, err)
	}
	return err
}

func extractGraphODataErrorCode(err error) (errorCode *string, odErr *odataerrors.ODataError, ok bool) {
	ok = errors.As(err, &odErr)
	if ok {
		errorCode = odErr.GetErrorEscaped().GetCode()
	}
	return
}
