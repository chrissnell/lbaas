package controller

import (
	"github.com/chrissnell/lbaas/config"
	"github.com/chrissnell/lbaas/controller/restapi"
	"github.com/chrissnell/lbaas/loadbalancer"
	"github.com/chrissnell/lbaas/model"
	"github.com/gorilla/mux"
)

type Controller struct {
	c  config.Config
	m  *model.Model
	lb loadbalancer.LoadBalancer
	r  *restapi.RestAPI
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

	apiRouter.HandleFunc("/vip", a.r.CreateVIP).Methods("POST")
	apiRouter.HandleFunc("/vip/{vipid}", a.r.GetVIP).Methods("GET")
	apiRouter.HandleFunc("/vip/{vipid}", a.r.DeleteVIP).Methods("DELETE")
	apiRouter.HandleFunc("/vip/{vipid}", a.r.UpdateVIP).Methods("PUT")
	apiRouter.HandleFunc("/vips", a.r.GetAllVIPs).Methods("GET")

	apiRouter.HandleFunc("/pool", a.r.CreatePool).Methods("POST")
	apiRouter.HandleFunc("/pool", a.r.DeletePool).Methods("DELETE")
	apiRouter.HandleFunc("/pool/{poolid}/members", a.r.AddPoolMember).Methods("PUT")
	apiRouter.HandleFunc("/pool/{poolid}/members/{member}", a.r.DeletePoolMember).Methods("DELETE")
	apiRouter.HandleFunc("/pool/{poolid}/members", a.r.GetAllPoolMembers).Methods("GET")
	apiRouter.HandleFunc("/pool/{poolid}/members", a.r.DeleteAllPoolMembers).Methods("DELETE")

	return apiRouter
}

// APIRouter will create a new gorilla Router for handling all Kubernetes REST API calls for
// the Kubernetes RESTful CloudProvider interface
func (a *Controller) KubeRouter() *mux.Router {
	KubeRouter := mux.NewRouter().StrictSlash(true)

	// add some KubeRouter.HandleFunc()'s here

	return KubeRouter
}
