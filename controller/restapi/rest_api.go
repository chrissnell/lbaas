package restapi

import (
	"errors"
	"fmt"
	"github.com/chrissnell/lbaas/config"
	"github.com/chrissnell/lbaas/model"
	"net/http"

	"github.com/emicklei/go-restful"
)

type RestAPI struct {
	c config.Config
	m *model.Model
}

func New(config config.Config, model *model.Model) *RestAPI {
	ra := &RestAPI{
		c: config,
		m: model,
	}

	return ra
}

func (ra *RestAPI) CreateVIP(req *restful.Request, resp *restful.Response) {
}

func (ra *RestAPI) DeleteVIP(req *restful.Request, resp *restful.Response) {
}

func (ra *RestAPI) UpdateVIP(req *restful.Request, resp *restful.Response) {
}

func (ra *RestAPI) GetVIP(req *restful.Request, resp *restful.Response) {

	vipid := req.PathParameter("vipid")

	res, err := ra.m.LB.GetVIP(vipid)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New(fmt.Sprint("VIP not found: ", vipid)))
	} else {
		resp.WriteEntity(res)
	}

}

func (ra *RestAPI) GetAllVIPs(req *restful.Request, resp *restful.Response) {
}

func (ra *RestAPI) CreatePool(req *restful.Request, resp *restful.Response) {
}

func (ra *RestAPI) DeletePool(req *restful.Request, resp *restful.Response) {
}

func (ra *RestAPI) AddPoolMembers(req *restful.Request, resp *restful.Response) {
}

func (ra *RestAPI) DeletePoolMember(req *restful.Request, resp *restful.Response) {
}

func (ra *RestAPI) GetAllPoolMembers(req *restful.Request, resp *restful.Response) {
}

func (ra *RestAPI) DeleteAllPoolMembers(req *restful.Request, resp *restful.Response) {
}
