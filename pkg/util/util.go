package util

import (
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	aimsv1 "github.com/dichque/grafana-operator/pkg/apis/grafana/v1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

var configPath = map[string]string{
	"grafana-dashboards": "dashboards.yaml",
	"kafka-dashboards":   "strimzi-kafka.json,strimzi-zookeeper.json,strimzi-kafka-exporter.json",
}

var configTmplPath = map[string]string{
	"grafana-config":      "grafana.ini",
	"grafana-datasources": "datasources.yaml",
}

func buildCMData(path string) map[string]string {
	m := make(map[string]string)

	for _, path := range strings.Split(path, ",") {
		data, err := ioutil.ReadFile(dashboardPath + path)
		if err != nil {
			klog.Errorf("Couldn't read file: %s %s", path, err)
		}
		m[path] = string(data)
	}

	return m
}

func buildCMDataFromTemplate(path string, cfg *grafanaConfig) map[string]string {
	m := make(map[string]string)

	tmpfile, err := ioutil.TempFile("/tmp", ".grafana-controller")
	if err != nil {
		klog.Errorf("Unable to create temporary file: %s : %s", path, err)
	}

	for _, path := range strings.Split(path, ",") {
		tmpl, err := template.ParseFiles(templatePath + path + ".tmpl")
		if err != nil {
			klog.Errorf("Unable to parse template file: %s: %s", path+".tmpl", err)
		}
		tmpl.Execute(tmpfile, cfg)
		tmpfile.Sync()
		tmpfile.Close()
		defer os.Remove(tmpfile.Name())

		data, err := ioutil.ReadFile(tmpfile.Name())
		if err != nil {
			klog.Errorf("Unable to generate config from template file: %s : %s", path, err)
		}
		m[path] = string(data)
	}

	return m
}

// CreateConfigMap returns configmaplist for loading to grafana deployment
func CreateConfigMap(grafana *aimsv1.Grafana, cmList *v1.ConfigMapList) *v1.ConfigMapList {

	gcfg := &grafanaConfig{
		AdminPassword: grafana.Spec.Password,
		AdminUser:     grafana.Spec.Username,
		PrometheusURL: grafana.Spec.PrometheusURL,
	}

	cmItems := []v1.ConfigMap{}

	for configName, path := range configPath {
		data := buildCMData(path)
		cm := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configName,
				Namespace: grafana.Namespace,
			},
			Data: data,
		}

		owner := metav1.NewControllerRef(
			grafana, aimsv1.SchemeGroupVersion.
				WithKind("Grafana"),
		)
		cm.ObjectMeta.OwnerReferences = append(cm.ObjectMeta.OwnerReferences, *owner)
		cmItems = append(cmItems, *cm)
	}

	for configTemplate, path := range configTmplPath {
		data := buildCMDataFromTemplate(path, gcfg)
		cm := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configTemplate,
				Namespace: grafana.Namespace,
			},
			Data: data,
		}

		owner := metav1.NewControllerRef(
			grafana, aimsv1.SchemeGroupVersion.
				WithKind("Grafana"),
		)
		cm.ObjectMeta.OwnerReferences = append(cm.ObjectMeta.OwnerReferences, *owner)
		cmItems = append(cmItems, *cm)
	}

	cmList.Items = cmItems
	return cmList

}

// Deployment creates grafana pod
func Deployment(grafana *aimsv1.Grafana) *appsv1.Deployment {
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      grafana.Name + "-deployment",
			Namespace: grafana.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "grafana"},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "grafana"}},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  grafana.Name + "-grafana",
							Image: grafana.Spec.Image,
							Ports: []v1.ContainerPort{{ContainerPort: 3000}},
							VolumeMounts: []v1.VolumeMount{
								{Name: "grafana-config", MountPath: "/etc/grafana"},
								{Name: "grafana-data", MountPath: "/var/lib/grafana"},
								{Name: "grafana-datasources", MountPath: "/etc/grafana/provisioning/datasources"},
								{Name: "grafana-dashboards", MountPath: "/etc/grafana/provisioning/dashboards"},
								{Name: "kafka-dashboards", MountPath: "/grafana-dashboard-definitions/0"},
							},
						},
					},
					ImagePullSecrets: []v1.LocalObjectReference{
						{
							Name: "intps-kafka-svc-pull-secret",
						},
					},
					Volumes: []v1.Volume{
						{
							Name: "grafana-config",
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: "grafana-config",
									},
								},
							},
						},
						{
							Name: "grafana-data",
							VolumeSource: v1.VolumeSource{
								EmptyDir: &v1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "grafana-datasources",
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: "grafana-datasources",
									},
								},
							},
						},
						{
							Name: "grafana-dashboards",
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: "grafana-dashboards",
									},
								},
							},
						},
						{
							Name: "kafka-dashboards",
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: "kafka-dashboards",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	owner := metav1.NewControllerRef(
		grafana, aimsv1.SchemeGroupVersion.
			WithKind("Grafana"),
	)
	deploy.ObjectMeta.OwnerReferences = append(deploy.ObjectMeta.OwnerReferences, *owner)

	return deploy
}
