package cert

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
)

type ResourceTemplateGraphVarData struct {
	ID string `json:"id,omitempty"`
}

type ResourceTemplateVarData struct {
	Graph *ResourceTemplateGraphVarData `json:"graph,omitempty"`
}

type TemplateVarData struct {
	Namespace *ResourceTemplateVarData `json:"namespace,omitempty"`
	Requester *ResourceTemplateVarData `json:"requester,omitempty"`
}

type templateContextKey string

const (
	templateContextKeyRequesterGraph templateContextKey = "requester.graph"
	templateContextKeyNamespaceGraph templateContextKey = "namespace.graph"
)

func processSubjectTemplate(c context.Context, subject certmodels.CertificateSubject) (certmodels.CertificateSubject, error) {
	var err error
	subject.CommonName, err = processTemplate(c, "subject.CommonName", subject.CommonName)
	if err != nil {
		return subject, err
	}
	return subject, nil
}

func processTemplate(c context.Context, templateName, templateStr string) (string, error) {
	bad := func(e error) (string, error) {
		return templateStr, e
	}
	preprocessed, hasTemplate, err := preprocessTemplate(templateStr)
	if err != nil {
		return bad(err)
	}
	if !hasTemplate {
		return templateStr, nil
	}
	tmpl, err := template.New(templateName).Option("missingkey=error").Parse(preprocessed)
	if err != nil {
		// template parse failed, something is wrong probably with preprocess
		log.Ctx(c).Error().Err(err).Str("originalTemplate", templateStr).Str("preprocessedTemplate", preprocessed).Msg("template parse failed")
		return bad(err)
	}
	sb := strings.Builder{}
	if err := tmpl.Execute(&sb, lookupContextForTemplateVars(c)); err != nil {
		return bad(err)
	}
	return sb.String(), nil
}

func lookupContextForTemplateVars(c context.Context) (d TemplateVarData) {
	if data, ok := c.Value(templateContextKeyRequesterGraph).(*ResourceTemplateGraphVarData); ok && data != nil {
		if d.Requester == nil {
			d.Requester = &ResourceTemplateVarData{}
		}
		d.Requester.Graph = data
	}
	if data, ok := c.Value(templateContextKeyNamespaceGraph).(*ResourceTemplateGraphVarData); ok && data != nil {
		if d.Namespace == nil {
			d.Namespace = &ResourceTemplateVarData{}
		}
		d.Namespace.Graph = data
	}
	return d
}

var allowedTemplateVars = map[string]string{
	"requester.graph.id": ".Requester.Graph.ID",
	"namespace.graph.id": ".Namespace.Graph.ID",
}

var varRegex = regexp.MustCompile(`\{\{([a-zA-Z0-9\.\-_]+)\}\}`)

func isSegmentValid(s string) (string, bool) {
	return s, !(strings.Contains(s, "{{") || strings.Contains(s, "}}"))
}

var (
	ErrTemplateInvalidSyntax = errors.New("template has invalid syntax")
)

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

	return strings.TrimSpace(sb.String()), hasTemplate, nil
}
