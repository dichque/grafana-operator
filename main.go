package main

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"

	clientset "github.com/dichque/grafana-operator/pkg/client/clientset/versioned"
	ginformers "github.com/dichque/grafana-operator/pkg/client/informers/externalversions"
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	flag.StringVar(&kubeconfig, "kubeconfig", defaultKubeconfig(), "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")

	klog.InitFlags(nil)

	flag.Parse()

	cfg, err := rest.InClusterConfig()
	if err != nil {
		cfg, err = clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
		if err != nil {
			klog.Fatalf("Error building kubeconfig: %s", err.Error())
		}
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	grafanaClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building cnat clientset: %s", err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Minute*1)
	grafanaInformerFactory := ginformers.NewSharedInformerFactory(grafanaClient, time.Minute*1)

	controller := NewController(kubeClient, grafanaClient, grafanaInformerFactory.Aims().V1().Grafanas(), kubeInformerFactory.Apps().V1().Deployments())

	kubeInformerFactory.Start(wait.NeverStop)
	grafanaInformerFactory.Start(wait.NeverStop)

	if err = controller.Run(2, wait.NeverStop); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}

func defaultKubeconfig() string {
	fname := os.Getenv("KUBECONFIG")
	if fname != "" {
		return fname
	}
	home, err := os.UserHomeDir()
	if err != nil {
		klog.Warningf("failed to get home directory: %v", err)
		return ""
	}
	return filepath.Join(home, ".kube", "config")
}
