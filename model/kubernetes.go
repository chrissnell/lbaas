package model

import (
	"fmt"

	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/watch"

	"github.com/chrissnell/lbaas/config"
)

type Kube struct {
	c *client.Client
}

func (k *Kube) New(c config.Config) (*Kube, error) {
	k8s := &Kube{}

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

	kc, err := client.New(&conf)
	if err != nil {
		return nil, err
	}

	k8s.c = kc
	return k8s, nil
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

func (k *Kube) NewKubeServicesWatcher(namespace string, l labels.Selector) (watch.Interface, error) {
	if namespace == "" {
		namespace = api.NamespaceDefault
	}

	if l == nil {
		l = labels.Everything()
	}

	w, err := k.c.Services(namespace).Watch(l, fields.Everything(), "")
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (k *Kube) NewKubeNodesWatcher(namespace string, l labels.Selector) (watch.Interface, error) {
	if namespace == "" {
		namespace = api.NamespaceDefault
	}

	if l == nil {
		l = labels.Everything()
	}

	w, err := k.c.Nodes().Watch(l, fields.Everything(), "")

	if err != nil {
		return nil, err
	}
	return w, nil
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
