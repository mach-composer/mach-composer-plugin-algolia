package internal

type SiteConfig struct {
	ApiKey string `mapstructure:"api_key"`
	AppId  string `mapstructure:"app_id"`
}
