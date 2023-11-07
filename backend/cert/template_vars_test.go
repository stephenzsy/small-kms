package cert

import (
	"context"
	"testing"

	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stretchr/testify/assert"
)

func TestProcessTemplate(t *testing.T) {
	ctx := context.Background()

	templateStr := "Hello, {{member.graph.id}}!"
	graphObj := graphmodels.NewDirectoryObject()
	graphObj.SetId(stringPtr("123"))
	ctx = context.WithValue(ctx, groupMemberGraphObjectContextKey, graphObj)
	result, err := processTemplate(ctx, "test", templateStr)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, 123!", result)
}

func TestProcessTemplateMissingVar(t *testing.T) {
	ctx := context.Background()

	templateStr := "Hello, {{member.graph.id}}!"
	_, err := processTemplate(ctx, "test", templateStr)
	assert.Error(t, err)
}

func TestProcessTemplateBadTemplate(t *testing.T) {
	ctx := context.Background()

	templateStr := "Hello, {{member.graph.bad}}!"
	_, err := processTemplate(ctx, "test", templateStr)
	assert.Error(t, err)
}

func stringPtr(s string) *string {
	return &s
}
