package controller

import (
	"github.com/chrissnell/lbaas/config"
	"github.com/chrissnell/lbaas/loadbalancer"
	"github.com/chrissnell/lbaas/model"
	"github.com/gorilla/mux"
)

type Controller struct {
	c  config.Config
	m  *model.Model
	lb loadbalancer.LoadBalancer
}

// New will create a new Controller
func New(config config.Config, model *model.Model, lb loadbalancer.LoadBalancer) *Controller {
	a := &Controller{
		c:  config,
		m:  model,
		lb: lb,
	}
	return a
}

// APIRouter will create a new gorilla Router for handling all REST API calls
func (a *Controller) APIRouter() *mux.Router {
	apiRouter := mux.NewRouter().StrictSlash(true)

	// add some apiRouter.HandleFunc

	return apiRouter
}

// APIRouter will create a new gorilla Router for handling all Kubernetes REST API calls for
// the Kubernetes RESTful CloudProvider interface
func (a *Controller) KubeRouter() *mux.Router {
	KubeRouter := mux.NewRouter().StrictSlash(true)

	// add some KubeRouter.HandleFunc()'s here

	return KubeRouter
}
