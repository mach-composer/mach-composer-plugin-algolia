package internal

import (
	"strings"
	"testing"

	"github.com/mach-composer/mach-composer-plugin-sdk/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func normalize(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), ""))
}

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

func TestRenderTerraformResourcesMultiApplication(t *testing.T) {
	p := NewAlgoliaPlugin()
	err := p.SetSiteConfig("test-site", map[string]any{
		"applications": []any{
			map[string]any{
				"name":    "lab_digital",
				"api_key": "lab-digital-key",
				"app_id":  "lab-digital-app",
			},
			map[string]any{
				"name":    "mach_composer",
				"api_key": "mach-composer-key",
				"app_id":  "mach-composer-app",
			},
		},
	})
	require.NoError(t, err)

	result, err := p.RenderTerraformResources("test-site")
	require.NoError(t, err)

	assert.Equal(t, normalize(`
		provider "algolia" {
			alias   = "lab_digital"
			api_key = "lab-digital-key"
			app_id  = "lab-digital-app"
		}

		provider "algolia" {
			alias   = "mach_composer"
			api_key = "mach-composer-key"
			app_id  = "mach-composer-app"
		}
	`), normalize(result))
}

func TestRenderTerraformResourcesMultiApplicationDefaultProvider(t *testing.T) {
	p := NewAlgoliaPlugin()
	err := p.SetSiteConfig("test-site", map[string]any{
		"applications": []any{
			map[string]any{
				"name":    "algolia",
				"api_key": "default-key",
				"app_id":  "default-app",
			},
		},
	})
	require.NoError(t, err)

	result, err := p.RenderTerraformResources("test-site")
	require.NoError(t, err)

	assert.Equal(t, normalize(`
		provider "algolia" {
			api_key = "default-key"
			app_id  = "default-app"
		}
	`), normalize(result))
}

func TestRenderTerraformComponentSingleApplication(t *testing.T) {
	p := NewAlgoliaPlugin()
	err := p.SetSiteConfig("test-site", map[string]any{
		"api_key": "foobar",
		"app_id":  "1234",
	})
	require.NoError(t, err)

	result, err := p.RenderTerraformComponent("test-site", "component")
	require.NoError(t, err)

	assert.Equal(t, []string{"algolia = algolia"}, result.Providers)
}

func TestRenderTerraformComponentMultiApplication(t *testing.T) {
	p := NewAlgoliaPlugin()
	err := p.SetSiteConfig("test-site", map[string]any{
		"applications": []any{
			map[string]any{
				"name":    "lab_digital",
				"api_key": "lab-digital-key",
				"app_id":  "lab-digital-app",
			},
			map[string]any{
				"name":    "mach_composer",
				"api_key": "mach-composer-key",
				"app_id":  "mach-composer-app",
			},
		},
	})
	require.NoError(t, err)

	result, err := p.RenderTerraformComponent("test-site", "lab_digital")
	require.NoError(t, err)

	assert.Equal(t, []string{"algolia = algolia.lab_digital"}, result.Providers)
}

func TestRenderTerraformComponentMultiApplicationDefaultProvider(t *testing.T) {
	p := NewAlgoliaPlugin()
	err := p.SetSiteConfig("test-site", map[string]any{
		"applications": []any{
			map[string]any{
				"name":    "algolia",
				"api_key": "default-key",
				"app_id":  "default-app",
			},
		},
	})
	require.NoError(t, err)

	result, err := p.RenderTerraformComponent("test-site", "algolia")
	require.NoError(t, err)

	assert.Empty(t, result.Providers)
}

func TestRenderTerraformComponentMultiApplicationNotFound(t *testing.T) {
	p := NewAlgoliaPlugin()
	err := p.SetSiteConfig("test-site", map[string]any{
		"applications": []any{
			map[string]any{
				"name":    "lab_digital",
				"api_key": "lab-digital-key",
				"app_id":  "lab-digital-app",
			},
		},
	})
	require.NoError(t, err)

	var customErr = &NoApplicationConfigError{}

	_, err = p.RenderTerraformComponent("test-site", "mach_composer")
	assert.ErrorAs(t, err, &customErr)
}

func TestSetSiteConfigEmpty(t *testing.T) {
	p := NewAlgoliaPlugin()

	err := p.SetSiteConfig("test-site", map[string]any{})
	assert.NoError(t, err)

	resources, err := p.RenderTerraformResources("test-site")
	require.NoError(t, err)
	assert.Empty(t, resources)
}

func TestSetSiteConfigPartialSingleApplication(t *testing.T) {
	p := NewAlgoliaPlugin()

	var customErr = &InvalidSiteConfigError{}

	err := p.SetSiteConfig("test-site", map[string]any{
		"api_key": "foobar",
	})
	assert.ErrorAs(t, err, &customErr)
}

func TestSetSiteConfigSingleAndMultiApplication(t *testing.T) {
	p := NewAlgoliaPlugin()

	var customErr = &InvalidSiteConfigError{}

	err := p.SetSiteConfig("test-site", map[string]any{
		"api_key": "foobar",
		"app_id":  "1234",
		"applications": []any{
			map[string]any{
				"name":    "lab_digital",
				"api_key": "lab-digital-key",
				"app_id":  "lab-digital-app",
			},
		},
	})
	assert.ErrorAs(t, err, &customErr)
}

func TestSetSiteConfigInvalidApplicationAlias(t *testing.T) {
	p := NewAlgoliaPlugin()

	err := p.SetSiteConfig("test-site", map[string]any{
		"applications": []any{
			map[string]any{
				"name":    "lab-digital",
				"api_key": "lab-digital-key",
				"app_id":  "lab-digital-app",
			},
		},
	})
	assert.Error(t, err)
}
