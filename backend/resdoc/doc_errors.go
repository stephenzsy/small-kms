package resdoc

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

var (
	ErrAzCosmosDocNotFound = errors.New("az cosmos doc not found")
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
