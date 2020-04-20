package config

const (
	TemplatePath  string = "config/templates/"
	DashboardPath string = "config/templates/dashboards/"
)

type grafanaConfig struct {
	AdminUser     string
	AdminPassword string
	PrometheusURL string
}
