package restapi

import (
	"fmt"
	"net/http"

	"github.com/chrissnell/lbaas/config"
	"github.com/chrissnell/lbaas/model"
	"github.com/chrissnell/lbaas/util/log"
	"github.com/chrissnell/lbaas/util/resterror"

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
	v := new(model.VIP)
	err := req.ReadEntity(&v)
	if err != nil {
		resterror.WriteErrorJSON(resp, http.StatusNotAcceptable, fmt.Errorf("Error parsing VIP entity: %v", err))
	}

	if v.Name == "" {
		resterror.WriteErrorJSON(resp, http.StatusNotAcceptable, fmt.Errorf("VIP name cannot be empty."))
	}

	// Let's check to see if a VIP by this name already exists in the datbase
	_, err = ra.m.S.FetchVIP(v.Name)
	if err == nil {
		resterror.WriteErrorJSON(resp, http.StatusNotFound, fmt.Errorf("VIP %v already exists.  Please choose a new name", v.Name))
		return
	}

	// Check with Kube to make sure the kube service name is valid
	// Check with cidrd to make sure this classname exists
	// Fetch an IP from cidrd and store UUID in VIP struct
	// Make sure a valid FE port was specified
	// Make sure a valid FE protocol was specified

	// Send some signals via channel

	err = ra.m.S.StoreVIP(v)
	if err != nil {
		resterror.WriteErrorJSON(resp, http.StatusInternalServerError, fmt.Errorf("Error storing VIP %v: %v", v.Name, err))

	}

	// type VIP struct {
	// 	Name               string       `json:"name"`
	// 	FrontEndIPClass    string       `json:"frontend_ip_class"`
	// 	FrontendIP         string       `json:"frontend_ip"`
	// 	FrontendIPUUID     string       `json:"frontend_ip_uuid"`
	// 	FrontendPort       uint8        `json:"frontend_port"`
	// 	FrontendProtocol   string       `json:"frontend_protocol"`
	// 	FrontendProfiles   []string     `json:"frontend_profiles"`
	// 	PoolMembers        []PoolMember `json:"pool_members"`
	// 	PoolMemberProtocol string       `json:"pool_member_protocols"` // HTTP, HTTPS, UDP, FTP, etc
	// 	KubeSvcName        string       `json:"kube_service_name"`
	// }

}

func (ra *RestAPI) DeleteVIP(req *restful.Request, resp *restful.Response) {
	vipid := req.PathParameter("vipid")

	logger.LogFn("model.LB.GetVIP", vipid)
	logger.Log("controller.LBEngine.DeleteVIPChan <- VIP struct...")
}

func (ra *RestAPI) UpdateVIP(req *restful.Request, resp *restful.Response) {
	v := new(model.VIP)
	err := req.ReadEntity(&v)
	if err != nil {
		resterror.WriteErrorJSON(resp, http.StatusNotAcceptable, fmt.Errorf("Error parsing VIP entity: %v", err))
		return
	}

	// First, we check to see if this VIP actually exists...
	_, err = ra.m.S.FetchVIP(v.Name)
	if err != nil {
		resterror.WriteErrorJSON(resp, http.StatusNotFound, fmt.Errorf("VIP %v does not exist: %v", v.Name, err))
		return
	}

	// The VIP is there, so we overwrite it with the VIP we were passed
	err = ra.m.S.StoreVIP(v)
	if err != nil {
		resterror.WriteErrorJSON(resp, http.StatusInternalServerError, fmt.Errorf("Error storing VIP %v: %v", v.Name, err))
		return
	}

	// Everything appears to have worked.
	return

}

func (ra *RestAPI) GetVIP(req *restful.Request, resp *restful.Response) {
	vipid := req.PathParameter("vipid")

	// Get and respond with the VIP from our model or throw an error if
	// it doesn't exist
	vip, err := ra.m.LB.GetVIP(vipid)
	if err != nil {
		resterror.WriteErrorJSON(resp, http.StatusNotFound, fmt.Errorf("VIP id %v not found", vipid))
		return
	}

	resp.WriteHeader(http.StatusOK)
	resp.WriteEntity(vip)

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
