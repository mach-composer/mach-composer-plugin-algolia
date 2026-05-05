package internal

import (
	"fmt"
	"regexp"

	"github.com/mach-composer/mach-composer-plugin-helpers/helpers"
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"
	"github.com/mach-composer/mach-composer-plugin-sdk/schema"
	"github.com/mitchellh/mapstructure"
)

var terraformAliasPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

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

	hasSingleApplication := cfg.ApiKey != "" || cfg.AppId != ""
	hasMultiApplication := cfg.IsMultiApplication()

	if hasSingleApplication && hasMultiApplication {
		return NewInvalidSiteConfigError("site %s cannot have both api_key/app_id and applications set", site)
	}

	if hasSingleApplication && (cfg.ApiKey == "" || cfg.AppId == "") {
		return NewInvalidSiteConfigError("site %s must have both api_key and app_id set", site)
	}

	for _, application := range cfg.Applications {
		if application.Name == "" || application.ApiKey == "" || application.AppId == "" {
			return NewInvalidSiteConfigError("site %s applications must have name, api_key, and app_id set", site)
		}
		if !isValidTerraformAlias(application.Name) {
			return NewInvalidSiteConfigError("site %s application %s must be a valid Terraform provider alias", site, application.Name)
		}
	}

	if p.siteConfigs == nil {
		p.siteConfigs = map[string]SiteConfig{}
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

	template := `
		{{- if .IsMultiApplication }}
		{{- range .Applications }}
		provider "algolia" {
			{{if ne .Name "algolia" }}{{ renderProperty "alias" .Name }}{{ end }}
			{{ renderProperty "api_key" .ApiKey }}
			{{ renderProperty "app_id" .AppId }}
		}
		{{- end }}
		{{- else }}
		provider "algolia" {
			{{ renderProperty "api_key" .ApiKey }}
			{{ renderProperty "app_id" .AppId }}
		}
		{{- end }}
	`
	return helpers.RenderGoTemplate(template, cfg)
}

func (p *Plugin) RenderTerraformComponent(site string, component string) (*schema.ComponentSchema, error) {
	cfg := p.getSiteConfig(site)
	if cfg != nil && cfg.IsMultiApplication() {
		application := cfg.GetApplicationConfig(component)
		if application == nil {
			return nil, NewNoApplicationConfigError("application %s not found in site %s. An application must exist with the same name as the component", component, site)
		}

		if component == "algolia" {
			return &schema.ComponentSchema{}, nil
		}

		return &schema.ComponentSchema{
			Providers: []string{
				fmt.Sprintf("algolia = algolia.%s", component),
			},
		}, nil
	}

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

func isValidTerraformAlias(value string) bool {
	return terraformAliasPattern.MatchString(value)
}
