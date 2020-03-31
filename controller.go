package main

import (
	"fmt"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsv1informer "k8s.io/client-go/informers/apps/v1"
	appsv1lister "k8s.io/client-go/listers/apps/v1"

	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	aimsv1 "github.com/dichque/grafana-operator/pkg/apis/grafana/v1"
	clientset "github.com/dichque/grafana-operator/pkg/client/clientset/versioned"
	"github.com/dichque/grafana-operator/pkg/client/clientset/versioned/scheme"
	gscheme "github.com/dichque/grafana-operator/pkg/client/clientset/versioned/scheme"
	ginformers "github.com/dichque/grafana-operator/pkg/client/informers/externalversions/grafana/v1"
	glisters "github.com/dichque/grafana-operator/pkg/client/listers/grafana/v1"
)

const controllerName = "grafana-controller"

// Controller struct for Grafana resources
type Controller struct {
	kubeClientset    kubernetes.Interface
	grafanaClientset clientset.Interface

	gLister glisters.GrafanaLister
	gSynced cache.InformerSynced

	deploymentLister appsv1lister.DeploymentLister
	deploymentSynced cache.InformerSynced

	workqueue workqueue.RateLimitingInterface
	recorder  record.EventRecorder
}

// NewController implementation for Grafana resources
func NewController(
	kubeClientset kubernetes.Interface,
	grafanaClientset clientset.Interface,
	ginformer ginformers.GrafanaInformer,
	deploymentInformer appsv1informer.DeploymentInformer) *Controller {

	utilruntime.Must(gscheme.AddToScheme(scheme.Scheme))
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerName})

	controller := &Controller{
		kubeClientset:    kubeClientset,
		grafanaClientset: grafanaClientset,
		gLister:          ginformer.Lister(),
		gSynced:          ginformer.Informer().HasSynced,
		deploymentLister: deploymentInformer.Lister(),
		deploymentSynced: deploymentInformer.Informer().HasSynced,
		workqueue:        workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Grafana"),
		recorder:         recorder,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when Grafana resources change
	ginformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueGrafana,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueGrafana(new)
		},
	})
	// Set up an event handler for when Deployment resources change
	deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueDeployment,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueDeployment(new)
		},
	})
	return controller

}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting grafana controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.gSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}
	if ok := cache.WaitForCacheSync(stopCh, c.deploymentSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	// Launch two workers to process At resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	defer c.workqueue.Done(obj)

	var key string
	var ok bool

	if key, ok = obj.(string); !ok {
		c.workqueue.Forget(obj)
		utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
		return true
	}

	err := c.syncHandler(key)

	if err == nil {
		c.workqueue.Forget(key)
		klog.Info("Successfully processed")
	} else if c.workqueue.NumRequeues(key) < 3 {
		c.workqueue.AddRateLimited(key)
		klog.Info("Re-processing the queue")
	} else {
		c.workqueue.Forget(key)
		klog.Error("Max retries reached")
	}

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) syncHandler(key string) error {
	klog.Infof("=== Reconciling At %s", key)

	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the At resource with this namespace/name
	original, err := c.gLister.Grafanas(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("grafana '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	// Clone because the original object is owned by the lister.
	instance := original.DeepCopy()

	if instance.Status.Phase == "" {
		instance.Status.Phase = aimsv1.PhasePending
	}

	return err
}

// enqueueGrafana takes a Grafana resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Grafana.
func (c *Controller) enqueueGrafana(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

// enqueue a  deployment and checks that the owner reference points to an Grafana object. It then
// enqueues this Grafana object.
func (c *Controller) enqueueDeployment(obj interface{}) {
	var deploy *appsv1.Deployment
	var ok bool
	if deploy, ok = obj.(*appsv1.Deployment); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding deployment, invalid type"))
			return
		}
		deploy, ok = tombstone.Obj.(*appsv1.Deployment)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding deployment tombstone, invalid type"))
			return
		}
		klog.V(4).Infof("Recovered deleted deployment '%s' from tombstone", deploy.GetName())
	}
	if ownerRef := metav1.GetControllerOf(deploy); ownerRef != nil {
		if ownerRef.Kind != "Grafana" {
			return
		}

		grafana, err := c.gLister.Grafanas(deploy.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			klog.V(4).Infof("ignoring orphaned deploy '%s' of At '%s'", deploy.GetSelfLink(), ownerRef.Name)
			return
		}

		klog.Infof("enqueuing Grafana %s/%s because deployment changed", grafana.Namespace, grafana.Name)
		c.enqueueGrafana(grafana)
	}
}
