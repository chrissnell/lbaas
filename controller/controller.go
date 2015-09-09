package controller

import (
	"github.com/chrissnell/lbaas/config"
	"github.com/chrissnell/lbaas/controller/restapi"
	"github.com/chrissnell/lbaas/model"
	"github.com/emicklei/go-restful"
)

type Controller struct {
	c  config.Config
	m  *model.Model
	R  *restapi.RestAPI
	WS *restful.WebService
}

// New will create a new Controller
func New(config config.Config, model *model.Model) *Controller {
	a := &Controller{
		c: config,
		m: model,
	}

	// Initialize the REST API
	a.R = restapi.New(config, model)

	a.WS = a.APIRouter()

	return a
}

// APIRouter will create a new go-restful router for handling REST API calls
func (a *Controller) APIRouter() *restful.WebService {

	ws := new(restful.WebService)
	ws.Consumes(restful.MIME_JSON)

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
