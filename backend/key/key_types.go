package key

import (
	"io"
	"strconv"

	"github.com/stephenzsy/small-kms/backend/base"
)

type (
	keyPolicyRefComposed struct {
		base.ResourceReference
		KeyPolicyRefFields
	}

	keyPolicyComposed struct {
		KeyPolicyRef
		KeyPolicyFields
	}

	keyComposed struct {
		base.ResourceReference
		KeySpec
		KeyFields
	}
)

func (ks *SigningKeySpec) WriteToDigest(w io.Writer) (s int, err error) {
	if ks == nil {
		return 0, nil
	}
	if ks.Alg != nil {
		if c, err := w.Write([]byte(*ks.Alg)); err != nil {
			return s + c, err
		} else {
			s += c
		}
	}
	if c, err := w.Write([]byte(ks.Kty)); err != nil {
		return s + c, err
	} else {
		s += c
	}
	switch ks.Kty {
	case "RSA":
		if ks.KeySize != nil {
			if c, err := w.Write([]byte(strconv.Itoa(int(*ks.KeySize)))); err != nil {
				return s + c, err
			} else {
				s += c
			}
		}
	case "EC":
		if ks.Crv != nil {
			if c, err := w.Write([]byte(*ks.Crv)); err != nil {
				return s + c, err
			} else {
				s += c
			}
		}
	}
	for _, op := range ks.KeyOperations {
		if c, err := w.Write([]byte(op)); err != nil {
			return s + c, err
		}
	}
	return s, nil
}

func (la *LifetimeAction) WriteToDigest(w io.Writer) (s int, err error) {
	if la == nil {
		return 0, nil
	}
	return la.Trigger.WriteToDigest(w)
}

func (lt *LifetimeTrigger) WriteToDigest(w io.Writer) (s int, err error) {
	if lt == nil {
		return 0, nil
	}
	if lt.TimeAfterCreate != nil {
		if c, err := w.Write([]byte("timeAfterCreate")); err != nil {
			return s + c, err
		} else {
			s += c
		}
		if c, err := w.Write(lt.TimeAfterCreate.Bytes()); err != nil {
			return s + c, err
		} else {
			s += c
		}
	}
	if lt.TimeBeforeExpiry != nil {
		if c, err := w.Write([]byte("timeBeforeExpiry")); err != nil {
			return s + c, err
		} else {
			s += c
		}
		if c, err := w.Write(lt.TimeBeforeExpiry.Bytes()); err != nil {
			return s + c, err
		} else {
			s += c
		}
	}
	if lt.PercentageAfterCreate != nil {
		if c, err := w.Write([]byte("percentageAfterCreate")); err != nil {
			return s + c, err
		} else {
			s += c
		}
		if c, err := w.Write([]byte(strconv.Itoa(int(*lt.PercentageAfterCreate)))); err != nil {
			return s + c, err
		} else {
			s += c
		}
	}
	return s, nil
}
