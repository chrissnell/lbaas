package model

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/cache"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/controller/framework"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/util"
	"k8s.io/kubernetes/pkg/util/workqueue"
	"k8s.io/kubernetes/pkg/watch"

	"github.com/chrissnell/lbaas/config"
)

type Kube struct {
	NodeQueue         *workqueue.Type
	ServiceQueue      *workqueue.Type
	NodeLister        cache.StoreToNodeLister
	ServiceLister     cache.StoreToServiceLister
	c                 *client.Client
	serviceController *framework.Controller
	nodeController    *framework.Controller
}

type QueueEvent struct {
	Obj     interface{}
	ObjType watch.EventType
}

func NewKubeClient(c config.Config, workQueueReady chan struct{}) (*Kube, error) {
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

	kube.NodeQueue = workqueue.New()
	kube.ServiceQueue = workqueue.New()

	// We will hardcode our namespace for now.
	namespace := api.NamespaceAll

	// Set up our enqueing function for node objects
	nodeEnqueueAsAdd := func(obj interface{}) {
		kube.NodeQueue.Add(QueueEvent{Obj: obj, ObjType: watch.Added})
	}

	nodeEnqueueAsDelete := func(obj interface{}) {
		kube.NodeQueue.Add(QueueEvent{Obj: obj, ObjType: watch.Deleted})
	}

	nodeEnqueueAsUpdate := func(obj interface{}) {
		kube.NodeQueue.Add(QueueEvent{Obj: obj, ObjType: watch.Modified})
	}

	// and one for service objects, too
	serviceEnqueueAsAdd := func(obj interface{}) {
		kube.ServiceQueue.Add(QueueEvent{Obj: obj, ObjType: watch.Added})
	}

	serviceEnqueueAsDelete := func(obj interface{}) {
		kube.ServiceQueue.Add(QueueEvent{Obj: obj, ObjType: watch.Deleted})
	}

	serviceEnqueueAsUpdate := func(obj interface{}) {
		kube.ServiceQueue.Add(QueueEvent{Obj: obj, ObjType: watch.Modified})
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

	kube.NodeLister.Store, kube.nodeController = framework.NewInformer(
		cache.NewListWatchFromClient(
			kube.c, "nodes", namespace, fields.Everything()),
		&api.Node{}, resyncPeriod, nodeEventHandlers)

	kube.ServiceLister.Store, kube.serviceController = framework.NewInformer(
		cache.NewListWatchFromClient(
			kube.c, "services", namespace, fields.Everything()),
		&api.Service{}, resyncPeriod, serviceEventHandlers)

	go kube.serviceController.Run(util.NeverStop)
	go kube.nodeController.Run(util.NeverStop)

	for !kube.serviceController.HasSynced() && !kube.nodeController.HasSynced() {
		log.Println("Waiting for serviceController and nodeController to sync...")
		time.Sleep(500 * time.Millisecond)
	}

	// Signal that the queue is ready
	close(workQueueReady)

	return kube, nil
}

// Gets a service by name, for a given namespace
func (k *Kube) GetKubeService(s string, namespace string) (*api.Service, error) {
	if namespace == "" {
		namespace = api.NamespaceAll
	}

	key := fmt.Sprint(namespace, "/", s)

	svc, exists, err := k.ServiceLister.Store.GetByKey(key)
	if !exists {
		return nil, fmt.Errorf("Service %v does not exist in namespace %v", s, namespace)
	}
	if err != nil {
		return nil, err
	}
	return svc.(*api.Service), nil
}

func (k *Kube) GetAllKubeServices(namespace string) (*[]api.Service, error) {
	var services []api.Service

	sl, err := k.ServiceLister.List()
	if err != nil {
		return nil, err
	}

	// Our ServiceLister returns all services in all namespaces.  If we weren't
	// passed a namespace, return slice of all services.
	if namespace == "" {
		return &sl.Items, nil
	}

	// Otherwise, go through the services and return only those in our namespace
	for _, sv := range sl.Items {
		if sv.Namespace == namespace {
			services = append(services, sv)
		}
	}

	return &services, nil
}

func (k *Kube) GetKubeServiceByUID(u string) (*api.Service, error) {
	svcs, err := k.GetAllKubeServices(api.NamespaceDefault)
	if err != nil {
		return nil, err
	}
	for _, s := range *svcs {
		if string(s.UID) == u {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("No service found with UID %v", u)
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

func (k *Kube) GetAllKubeNodes() (*api.NodeList, error) {
	nodes, err := k.NodeLister.List()
	if err != nil {
		return nil, err
	}
	return &nodes, nil
}
