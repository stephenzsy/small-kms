package common

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidCertificateFieldVarAttribute = errors.New("invalid certificate field variable attribute")
	ErrInvalidCertificateFieldVarSyntax    = errors.New("invalid certificate field variable syntax")
	ErrInvalidCertificateFieldVarSelector  = errors.New("invalid certificate field variable selector")
	ErrInvalidCertificateFieldVarName      = errors.New("invalid certificate field variable name")

	ErrCertificateFieldVarRequired = errors.New("certificate field variable is required")
)

type CertificateFieldVar struct {
	Required    bool
	Selector    CertificateFieldVarNamespaceSelectorToken
	Subselector string
	Name        string
}

func (v *CertificateFieldVar) FullVariableKey() string {
	var s string = string(v.Selector)
	if len(v.Subselector) > 0 {
		s += "." + v.Subselector
	}
	return s + "." + v.Name
}

func (v *CertificateFieldVar) parseIdentifier(s string) error {
	if token, rest, found := strings.Cut(s, "."); found {
		token = strings.TrimSpace(token)
		switch token {
		case string(CertFieldVarNsSelectorSelf),
			string(CertFieldVarNsSelectorCaller),
			string(CertFieldVarNsSelectorLinked),
			string(CertFieldVarNsSelectorRequest):
			v.Selector = CertificateFieldVarNamespaceSelectorToken(token)
		default:
			return fmt.Errorf("%w: %s", ErrInvalidCertificateFieldVarSelector, token)
		}
		s = rest
	} else {
		return fmt.Errorf("%w: no selector", ErrInvalidCertificateFieldVarSyntax)
	}
	if v.Selector == CertFieldVarNsSelectorLinked {
		if token, rest, found := strings.Cut(s, "."); found {
			token = strings.TrimSpace(token)
			v.Subselector = token
			s = rest
		} else {
			return fmt.Errorf("%w: no linked subselector", ErrInvalidCertificateFieldVarSyntax)
		}
	}
	// rest would be the name
	switch s {
	case string(CertFieldVarNameID),
		string(CertFieldVarNamePath),
		string(CertFieldVarNameURI),
		string(CertFieldVarNameRequestFQDN):
		v.Name = s

	default:
		return fmt.Errorf("%w: %s", ErrInvalidCertificateFieldVarName, s)
	}
	if v.Name == string(CertFieldVarNameRequestFQDN) && v.Selector != CertFieldVarNsSelectorRequest {
		return fmt.Errorf("%w: %s", ErrInvalidCertificateFieldVarName, s)
	}
	return nil
}

func (v *CertificateFieldVar) parseAttributes(s string) error {
	tokens := strings.Split(s, ",")
	for _, s := range tokens {
		s = strings.TrimSpace(s)

		switch s {
		case "":
			continue
		case string(CertFieldVarAttrRequired):
			v.Required = true
		default:
			return fmt.Errorf("%w: %s", ErrInvalidCertificateFieldVarAttribute, s)
		}
	}
	return nil
}

func ParseCertificateFieldVar(s string) (v CertificateFieldVar, err error) {
	s_identifier, s_attributes, _ := strings.Cut(s, ",")
	if err = v.parseIdentifier(s_identifier); err != nil {
		return
	}
	err = v.parseAttributes(s_attributes)
	return
}

func (v *CertificateFieldVar) Substitute(variableValues map[string]string) (string, error) {
	key := v.FullVariableKey()
	var value string
	if variableValues != nil {
		value = variableValues[key]
	}
	if len(value) == 0 && v.Required {
		return value, fmt.Errorf("%w: %s", ErrCertificateFieldVarRequired, key)
	}
	return value, nil
}
