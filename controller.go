package main

import (
	"fmt"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"

	corev1informer "k8s.io/client-go/informers/core/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	corev1lister "k8s.io/client-go/listers/core/v1"

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

	cmLister corev1lister.ConfigMapLister
	cmSynced cache.InformerSynced

	workqueue workqueue.RateLimitingInterface
	recorder  record.EventRecorder
}

// NewController implementation for Grafana resources
func NewController(
	kubeClientset kubernetes.Interface,
	grafanaClientset clientset.Interface,
	ginformer ginformers.GrafanaInformer,
	cmInformer corev1informer.ConfigMapInformer) *Controller {

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
		cmLister:         cmInformer.Lister(),
		cmSynced:         cmInformer.Informer().HasSynced,
		workqueue:        workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Grafana"),
		recorder:         recorder,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when At resources change
	ginformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueGrafana,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueGrafana(new)
		},
	})
	// Set up an event handler for when Pod resources change
	cmInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueCM,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueCM(new)
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
	if ok := cache.WaitForCacheSync(stopCh, c.cmSynced); !ok {
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

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool

		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		if when, err := c.syncHandler(key); err != nil {
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		} else if when != time.Duration(0) {
			c.workqueue.AddAfter(key, when)
		} else {
			c.workqueue.Forget(obj)
		}
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) syncHandler(key string) (time.Duration, error) {
	klog.Infof("=== Reconciling At %s", key)

	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return time.Duration(0), nil
	}

	// Get the At resource with this namespace/name
	original, err := c.gLister.Grafanas(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("grafana '%s' in work queue no longer exists", key))
			return time.Duration(0), nil
		}

		return time.Duration(0), err
	}

	// Clone because the original object is owned by the lister.
	instance := original.DeepCopy()

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

// enqueueCM a configmap and checks that the owner reference points to an Grafana object. It then
// enqueues this Grafana object.
func (c *Controller) enqueueCM(obj interface{}) {
	var cm *corev1.ConfigMap
	var ok bool
	if cm, ok = obj.(*corev1.ConfigMap); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding configmap, invalid type"))
			return
		}
		cm, ok = tombstone.Obj.(*corev1.ConfigMap)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding configmap tombstone, invalid type"))
			return
		}
		klog.V(4).Infof("Recovered deleted configmap '%s' from tombstone", cm.GetName())
	}
	if ownerRef := metav1.GetControllerOf(cm); ownerRef != nil {
		if ownerRef.Kind != "Grafana" {
			return
		}

		grafana, err := c.gLister.Grafanas(cm.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			klog.V(4).Infof("ignoring orphaned cm '%s' of At '%s'", cm.GetSelfLink(), ownerRef.Name)
			return
		}

		klog.Infof("enqueuing At %s/%s because configmap changed", grafana.Namespace, grafana.Name)
		c.enqueueGrafana(grafana)
	}
}
