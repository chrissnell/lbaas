package model

import (
	"fmt"

	api "k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"

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
	svc, err := k.c.Services(namespace).Get(s)
	if err != nil {
		return nil, err
	}
	return svc, nil
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
