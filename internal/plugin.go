package internal

import (
	"fmt"

	"github.com/creasty/defaults"
	"github.com/mach-composer/mach-composer-plugin-helpers/helpers"
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"
	"github.com/mach-composer/mach-composer-plugin-sdk/schema"
	"github.com/mitchellh/mapstructure"
)

type Plugin struct {
	environment string
	provider    string
	siteConfigs map[string]SiteConfig
}

func NewAlgoliaPlugin() schema.MachComposerPlugin {
	state := &Plugin{
		provider:    "0.5.7",
		siteConfigs: map[string]SiteConfig{},
	}

	return plugin.NewPlugin(&schema.PluginSchema{
		Identifier: "algolia",

		Configure: state.Configure,
		IsEnabled: state.IsEnabled,

		GetValidationSchema: state.GetValidationSchema,

		// Config
		SetSiteConfig: state.SetSiteConfig,

		// Renders
		RenderTerraformProviders: state.TerraformRenderProviders,
		RenderTerraformResources: state.TerraformRenderResources,
		RenderTerraformComponent: state.RenderTerraformComponent,
	})
}

func (p *Plugin) Configure(environment string, provider string) error {
	p.environment = environment
	if provider != "" {
		p.provider = provider
	}
	return nil
}

func (p *Plugin) IsEnabled() bool {
	return true
}

func (p *Plugin) GetValidationSchema() (*schema.ValidationSchema, error) {
	result := getSchema()
	return result, nil
}

func (p *Plugin) SetSiteConfig(site string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}

	cfg := SiteConfig{}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}

	if err := defaults.Set(&cfg); err != nil {
		return err
	}
	p.siteConfigs[site] = cfg
	return nil
}

func (p *Plugin) TerraformRenderProviders(_ string) (string, error) {
	result := fmt.Sprintf(`
		algolia = {
			source  = "k-yomo/algolia"
			version = "%s"
		}`, helpers.VersionConstraint(p.provider))
	return result, nil
}

func (p *Plugin) TerraformRenderResources(site string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
	}

	templateContext := struct {
		ApiKey string
		AppId  string
	}{
		ApiKey: cfg.ApiKey,
		AppId:  cfg.AppId,
	}

	template := `
		provider "algolia" {
			{{ renderProperty "api_key" .ApiKey }}
			{{ renderProperty "app_id" .AppId }}
		}
	`
	return helpers.RenderGoTemplate(template, templateContext)
}

func (p *Plugin) RenderTerraformComponent(_ string, _ string) (*schema.ComponentSchema, error) {
	result := &schema.ComponentSchema{
		Providers: []string{
			"algolia = algolia",
		},
	}
	return result, nil
}

func (p *Plugin) getSiteConfig(site string) *SiteConfig {
	cfg, ok := p.siteConfigs[site]
	if !ok {
		return nil
	}
	return &cfg
}
