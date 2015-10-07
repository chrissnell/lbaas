package controller

import (
	"github.com/emicklei/go-restful"

	"github.com/chrissnell/lbaas/config"
	"github.com/chrissnell/lbaas/controller/restapi"
	"github.com/chrissnell/lbaas/model"
)

type Controller struct {
	c  config.Config
	M  *model.Model
	R  *restapi.RestAPI
	WS *restful.WebService
	NE *NodesEngine
	LB *LBEngine
	SE *ServicesEngine
}

// New will create a new Controller
func NewController(config config.Config, m *model.Model) *Controller {

	c := &Controller{
		c: config,
		M: m,
	}

	// Initialize the REST API
	c.R = restapi.New(config, m)

	// Create and start the Nodes watcher engine
	c.NE = NewNodesEngine(m)

	// Create and start the Services watcher engine
	c.SE = NewServicesEngine(m)

	// Create and start the LB updater engine
	c.LB = NewLBEngine(m, c.NE.NodeChangeChan)

	// Start routing the API
	c.WS = c.APIRouter()

	return c
}

// APIRouter will create a new go-restful router for handling REST API calls
func (a *Controller) APIRouter() *restful.WebService {

	ws := new(restful.WebService)
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/vip").To(a.R.CreateVIP))
	ws.Route(ws.GET("/vip/{vipid}").To(a.R.GetVIP))
	ws.Route(ws.DELETE("/vip/{vipid}").To(a.R.DeleteVIP))
	ws.Route(ws.PUT("/vip/{vipid}").To(a.R.UpdateVIP))
	ws.Route(ws.PUT("/vips").To(a.R.GetAllVIPs))

	ws.Route(ws.POST("/pool").To(a.R.CreatePool))
	ws.Route(ws.DELETE("/pool").To(a.R.DeletePool))
	ws.Route(ws.PUT("/pool/{poolid}/members").To(a.R.AddPoolMembers))
	ws.Route(ws.DELETE("/pool/{poolid}/members/{member}").To(a.R.AddPoolMembers))
	ws.Route(ws.GET("/pool/{poolid}/members").To(a.R.GetAllPoolMembers))
	ws.Route(ws.DELETE("/pool/{poolid}/members").To(a.R.DeleteAllPoolMembers))

	return ws

}
