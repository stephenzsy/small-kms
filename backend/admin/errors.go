package admin

import (
	"errors"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

var notFoundError = errors.New("not found")

func wrapNotFoundError(err error) error {
	var respErr *azcore.ResponseError
	if errors.As(err, &respErr) {
		// Handle Error
		if respErr.StatusCode == http.StatusNotFound {
			return notFoundError
		}
	}
	return err
}
