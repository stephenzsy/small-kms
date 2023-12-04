package cloudkey

import "errors"

var (
	errInvalidKeyType = errors.New("invalid key type")
	errInvalidKeySize = errors.New("invalid key size")
	errInvalidKey     = errors.New("invalid key")
	errInvalidAlg     = errors.New("invalid algorithm")
	errInvalidCurve   = errors.New("invalid curve")
)

func exportErr(err error) error {
	return err
}

var (
	ErrInvalidKeyType   = exportErr(errInvalidKeyType)
	ErrInvalidCurve     = exportErr(errInvalidCurve)
	ErrInvalidKey       = exportErr(errInvalidKey)
	ErrInvalidKeySize   = exportErr(errInvalidKeySize)
	ErrInvalidAlgorithm = exportErr(errInvalidAlg)
)
