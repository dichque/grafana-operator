package util

const templatePath string = "/Users/jaganaga/WA/golang/src/github.com/dichque/grafana-operator/config/templates/"
const dashboardPath string = "/Users/jaganaga/WA/golang/src/github.com/dichque/grafana-operator/config/templates/dashboards/"

type grafanaConfig struct {
	AdminUser     string
	AdminPassword string
	PrometheusURL string
}
