package internal

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/e2b-dev/infra/tests/integration/internal/api"
	"github.com/e2b-dev/infra/tests/integration/internal/setup"
	"github.com/e2b-dev/infra/tests/integration/internal/utils"
)

func TestMain(m *testing.M) {
	log.Println("Setting up test environment")
	m.Run()
	log.Println("Environment set up")
}

// TestCacheTemplate starts a sandbox before all tests to cache the necessary files for the base template.
func TestCacheTemplate(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	c := setup.GetAPIClient()
	sbxTimeout := int32(60)
	sbx, err := c.PostSandboxesWithResponse(ctx, api.NewSandbox{
		TemplateID: setup.SandboxTemplateID,
		Timeout:    &sbxTimeout,
	}, setup.WithAPIKey())
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, sbx.StatusCode())

	t.Cleanup(func() {
		if sbx == nil {
			t.Logf("Error: %v", err)
		} else if sbx.JSON201 == nil {
			t.Logf("Response error: %d %v", sbx.StatusCode(), string(sbx.Body))
		} else {
			utils.TeardownSandbox(t, c, sbx.JSON201.SandboxID)
		}
	})
}
