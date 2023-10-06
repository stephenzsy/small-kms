package certtemplate

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	TemplateVarNameDeviceURI              string = "device.uri"
	TemplateVarNameDeviceAltURI           string = "device.altUri"
	TemplateVarNameApplicationURI         string = "application.uri"
	TemplateVarNameApplicationAltURI      string = "application.altUri"
	TemplateVarNameServicePrincipalID     string = "servicePrincipal.id"
	TemplateVarNameServicePrincipalURI    string = "servicePrincipal.uri"
	TemplateVarNameServicePrincipalAltURI string = "servicePrincipal.altUri"
	TemplateVarNameGroupID                string = "group.id"
	TemplateVarNameGroupURI               string = "group.uri"
)

var (
	ErrTemplateInvalidSyntax = errors.New("template has invalid syntax")
)

type ResourceTemplateVarData struct {
	ID     string `json:"id"`
	URI    string `json:"uri"`
	AltURI string `json:"altUri"`
}

type TemplateVarData struct {
	Device           ResourceTemplateVarData `json:"device"`
	Application      ResourceTemplateVarData `json:"application"`
	ServicePrincipal ResourceTemplateVarData `json:"servicePrincipal"`
	Group            ResourceTemplateVarData `json:"group"`
}

var allowedTemplateVars = map[string]string{
	"device.uri":              ".Device.URI",
	"device.altUri":           ".Device.AltURI",
	"application.uri":         ".Application.URI",
	"application.altUri":      ".Application.AltURI",
	"servicePrincipal.id":     ".ServicePrincipal.ID",
	"servicePrincipal.uri":    ".ServicePrincipal.URI",
	"servicePrincipal.altUri": ".ServicePrincipal.AltURI",
	"group.id":                ".Group.ID",
	"group.uri":               ".Group.URI",
}

var varRegex = regexp.MustCompile(`\{\{([a-zA-Z0-9\.\-_]+)\}\}`)

func isSegmentValid(s string) (string, bool) {
	return s, !(strings.Contains(s, "{{") || strings.Contains(s, "}}"))
}

func preprocessTemplate(s string) (transformed string, hasTemplate bool, err error) {
	allMatches := varRegex.FindAllStringIndex(s, -1)
	sb := strings.Builder{}
	sbInd := 0
	for _, match := range allMatches {
		hasTemplate = true
		if preSegment, ok := isSegmentValid(s[sbInd:match[0]]); ok {
			sb.WriteString(preSegment)
		} else {
			return s, false, fmt.Errorf("%w: unmatched '{{' or '}}'", ErrTemplateInvalidSyntax)
		}

		matchedInner := varRegex.FindStringSubmatch(s)[1]
		if t, ok := allowedTemplateVars[matchedInner]; ok {
			sb.WriteString("{{ ")
			sb.WriteString(t)
			sb.WriteString(" }}")
		} else {
			return s, false, fmt.Errorf("%w: invalid template variable '%s'", ErrTemplateInvalidSyntax, matchedInner)
		}
		sbInd = match[1]
	}
	if sbInd < len(s) {
		if lastSegment, ok := isSegmentValid(s[sbInd:]); ok {
			sb.WriteString(lastSegment)
		} else {
			return s, false, fmt.Errorf("%w: unmatched '{{' or '}}'", ErrTemplateInvalidSyntax)
		}
	}

	return sb.String(), hasTemplate, nil
}
