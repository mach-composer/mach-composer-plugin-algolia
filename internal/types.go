package internal

type SiteConfig struct {
	ApiKey       string              `mapstructure:"api_key"`
	AppId        string              `mapstructure:"app_id"`
	Applications []ApplicationConfig `mapstructure:"applications"`
}

type ApplicationConfig struct {
	Name   string `mapstructure:"name"`
	ApiKey string `mapstructure:"api_key"`
	AppId  string `mapstructure:"app_id"`
}

func (s *SiteConfig) IsMultiApplication() bool {
	return len(s.Applications) > 0
}

func (s *SiteConfig) GetApplicationConfig(name string) *ApplicationConfig {
	for _, application := range s.Applications {
		if application.Name == name {
			return &application
		}
	}
	return nil
}
