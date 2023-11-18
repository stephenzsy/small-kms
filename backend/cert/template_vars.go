package cert

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/rs/zerolog/log"
)

func ProcessTemplate(c context.Context, templateName, templateStr string) (string, error) {
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

var (
	ErrTemplateInvalidSyntax = errors.New("template has invalid syntax")
)

func lookupContextForTemplateVars(c context.Context) (d TemplateVarData) {
	if dirObj, ok := c.Value(groupMemberGraphObjectContextKey).(graphmodels.DirectoryObjectable); ok && dirObj != nil {
		if d.Member == nil {
			d.Member = &ResourceTemplateVarData{}
		}
		d.Member.Graph = &ResourceTemplateGraphVarData{
			ID: dirObj.GetId(),
		}
	}
	if dirObj, ok := c.Value(selfGraphObjectContextKey).(graphmodels.DirectoryObjectable); ok && dirObj != nil {
		if d.My == nil {
			d.My = &ResourceTemplateVarData{}
		}
		d.My.Graph = &ResourceTemplateGraphVarData{
			ID: dirObj.GetId(),
		}
	}
	return d
}

type ResourceTemplateGraphVarData struct {
	ID *string `json:"id,omitempty"`
}

type ResourceTemplateVarData struct {
	Graph *ResourceTemplateGraphVarData `json:"graph,omitempty"`
}

type TemplateVarData struct {
	Member *ResourceTemplateVarData `json:"member,omitempty"`
	My     *ResourceTemplateVarData `json:"my,omitempty"`
}

var allowedTemplateVars = map[string]string{
	"member.graph.id": ".Member.Graph.ID",
	"my.graph.id":     ".My.Graph.ID",
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

	return strings.TrimSpace(sb.String()), hasTemplate, nil
}
