package model

import (
	"fmt"
	"reflect"
	"time"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/cache"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/controller/framework"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/util"
	"k8s.io/kubernetes/pkg/util/workqueue"
	"k8s.io/kubernetes/pkg/watch"

	"github.com/chrissnell/lbaas/config"
)

type Kube struct {
	c                 *client.Client
	serviceController *framework.Controller
	nodeController    *framework.Controller
	NodeQueue         *workqueue.Type
	ServiceQueue      *workqueue.Type
	nodeLister        cache.StoreToNodeLister
	serviceLister     cache.StoreToServiceLister
}

type QueueEvent struct {
	Obj     interface{}
	ObjType watch.EventType
}

func (k *Kube) New(c config.Config, workQueueReady chan struct{}) (*Kube, error) {
	const resyncPeriod = 10 * time.Second
	var err error

	kube := &Kube{}

	conf := client.Config{
		Host: c.Kubernetes.APIendpoint,
	}

	// We only pass the username + password if they were set...
	if c.Kubernetes.Username != "" {
		conf.Username = c.Kubernetes.Username
	}

	if c.Kubernetes.Password != "" {
		conf.Password = c.Kubernetes.Password
	}

	kube.c, err = client.New(&conf)
	if err != nil {
		return nil, err
	}

	k.NodeQueue = workqueue.New()
	k.ServiceQueue = workqueue.New()

	// We will hardcode our namespace for now.
	namespace := api.NamespaceAll

	// Set up our enqueing function for node objects
	nodeEnqueueAsAdd := func(obj interface{}) {
		k.NodeQueue.Add(QueueEvent{Obj: obj, ObjType: watch.Added})
	}

	nodeEnqueueAsDelete := func(obj interface{}) {
		k.NodeQueue.Add(QueueEvent{Obj: obj, ObjType: watch.Deleted})
	}

	nodeEnqueueAsUpdate := func(obj interface{}) {
		k.NodeQueue.Add(QueueEvent{Obj: obj, ObjType: watch.Modified})
	}

	// and one for service objects, too
	serviceEnqueueAsAdd := func(obj interface{}) {
		k.ServiceQueue.Add(QueueEvent{Obj: obj, ObjType: watch.Added})
	}

	serviceEnqueueAsDelete := func(obj interface{}) {
		k.ServiceQueue.Add(QueueEvent{Obj: obj, ObjType: watch.Deleted})
	}

	serviceEnqueueAsUpdate := func(obj interface{}) {
		k.ServiceQueue.Add(QueueEvent{Obj: obj, ObjType: watch.Modified})
	}

	// Set up our event handlers.  These get called every time the cache client gets a new event from the API.
	nodeEventHandlers := framework.ResourceEventHandlerFuncs{
		AddFunc:    nodeEnqueueAsAdd,
		DeleteFunc: nodeEnqueueAsDelete,
		UpdateFunc: func(old, cur interface{}) {
			// We're only going to add updates to the queue when the node condition changes
			if old.(*api.Node).Status.Conditions[0].Status != cur.(*api.Node).Status.Conditions[0].Status {
				nodeEnqueueAsUpdate(cur)
			}
		},
	}

	serviceEventHandlers := framework.ResourceEventHandlerFuncs{
		AddFunc:    serviceEnqueueAsAdd,
		DeleteFunc: serviceEnqueueAsDelete,
		UpdateFunc: func(old, cur interface{}) {
			if !reflect.DeepEqual(old, cur) {
				serviceEnqueueAsUpdate(cur)
			}
		},
	}

	k.nodeLister.Store, k.nodeController = framework.NewInformer(
		cache.NewListWatchFromClient(
			kube.c, "nodes", namespace, fields.Everything()),
		&api.Node{}, resyncPeriod, nodeEventHandlers)

	k.serviceLister.Store, k.serviceController = framework.NewInformer(
		cache.NewListWatchFromClient(
			kube.c, "services", namespace, fields.Everything()),
		&api.Service{}, resyncPeriod, serviceEventHandlers)

	go k.serviceController.Run(util.NeverStop)
	go k.nodeController.Run(util.NeverStop)

	// Signal that the queue is ready
	close(workQueueReady)

	return kube, nil
}

// Gets a service by name, for a given namespace
func (k *Kube) GetKubeService(s string, namespace string) (*api.Service, error) {
	if namespace == "" {
		namespace = api.NamespaceDefault
	}

	svc, err := k.c.Services(namespace).Get(s)
	if err != nil {
		return nil, err
	}
	return svc, nil
}

func (k *Kube) GetAllKubeServices(namespace string) (*api.ServiceList, error) {
	if namespace == "" {
		namespace = api.NamespaceDefault
	}

	sl, err := k.c.Services(namespace).List(labels.Everything())
	if err != nil {
		return nil, err
	}

	return sl, nil
}

func (k *Kube) GetNodePortForServiceByPortName(s *api.Service, portName string) (int, error) {
	if portName == "" {
		return 0, fmt.Errorf("You must specifiy a port name.")
	}

	for _, p := range s.Spec.Ports {
		if p.Name == portName {
			return p.NodePort, nil
		}
	}

	return 0, fmt.Errorf("No port named %v for service %v were found.", portName, s.Name)
}

func (k *Kube) VerifyKubeService(v *VIP) (bool, error) {

	// Let's see if the service exists in Kubernetes...
	ks, err := k.GetKubeService(v.KubeSvcName, v.KubeNamespace)
	if err != nil {
		return false, fmt.Errorf("Kubernetes service name %v could not be found in namespace %v: %v", v.KubeSvcName, v.KubeNamespace, err)
	}

	// Let's make sure that the Kuberenetes service has a NodePort for the supplied port name
	np, err := k.GetNodePortForServiceByPortName(ks, v.KubeSvcPortName)
	if err != nil {
		return false, fmt.Errorf("Kubernetes service %v does not have a NodePort for port %v", v.KubeSvcName, np)
	}

	return true, nil
}

func (k *Kube) GetAllKubeNodes(namespace string) (*api.NodeList, error) {
	if namespace == "" {
		namespace = api.NamespaceDefault
	}

	nl, err := k.c.Nodes().List(labels.Everything(), fields.Everything())
	if err != nil {
		return nil, err
	}
	return nl, nil
}
