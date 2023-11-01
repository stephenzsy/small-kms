package cloudkey

import "errors"

var (
	errInvalidKeyType = errors.New("invalid key type")
	errInvalidKey     = errors.New("invalid key")
	errInvalidAlg     = errors.New("invalid algorithm")
)

func exportErr(err error) error {
	return err
}

var (
	ErrInvalidKeyType   = exportErr(errInvalidKeyType)
	ErrInvalidKey       = exportErr(errInvalidKey)
	ErrInvalidAlgorithm = exportErr(errInvalidAlg)
)
