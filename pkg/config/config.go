package config

const (
	TemplatePath  string = "config/templates/"
	DashboardPath string = "config/templates/dashboards/"
)

type GrafanaConfig struct {
	AdminUser     string
	AdminPassword string
	PrometheusURL string
}
