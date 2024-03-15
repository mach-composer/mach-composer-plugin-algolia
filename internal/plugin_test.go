package internal

import (
	"testing"

	"github.com/mach-composer/mach-composer-plugin-sdk/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderTerraformResources(t *testing.T) {
	tests := []struct {
		name   string
		create func() schema.MachComposerPlugin
	}{
		{
			name: "Render",
			create: func() schema.MachComposerPlugin {
				p := NewAlgoliaPlugin()
				err := p.SetSiteConfig("test-site", map[string]any{
					"api_key": "foobar",
					"app_id":  "1234",
				})
				assert.NoError(t, err)
				return p
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plugin := tt.create()
			result, err := plugin.RenderTerraformResources("test-site")
			require.NoError(t, err)

			assert.Contains(t, result, `api_key = "foobar"`)
			assert.Contains(t, result, `app_id = "1234"`)
		})
	}
}
