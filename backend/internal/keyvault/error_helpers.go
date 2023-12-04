package kv

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

var (
	ErrAzKeyVaultItemNotFound = errors.New("az key vault key not found")
)

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
