package cert

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessTemplate(t *testing.T) {
	ctx := context.Background()

	templateStr := "Hello, {{requester.graph.id}}!"

	ctx = context.WithValue(ctx, templateContextKeyRequesterGraph, &ResourceTemplateGraphVarData{
		ID: "123",
	})
	result, err := processTemplate(ctx, "test", templateStr)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, 123!", result)
}

func TestProcessTemplateMissingVar(t *testing.T) {
	ctx := context.Background()

	templateStr := "Hello, {{requester.graph.id}}!"
	_, err := processTemplate(ctx, "test", templateStr)
	assert.Error(t, err)
}

func TestProcessTemplateBadTemplate(t *testing.T) {
	ctx := context.Background()

	templateStr := "Hello, {{requester.graph.bad}}!"
	_, err := processTemplate(ctx, "test", templateStr)
	assert.Error(t, err)
}
