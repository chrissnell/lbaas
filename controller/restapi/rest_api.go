package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chrissnell/lbaas/config"
	"github.com/chrissnell/lbaas/model"

	"github.com/emicklei/go-restful"
	api "k8s.io/kubernetes/pkg/api"
)

type RestAPI struct {
	c config.Config
	m *model.Model
}

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
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
		WriteErrorJSON(resp, http.StatusNotAcceptable, fmt.Errorf("Error parsing VIP entity JSON: %v", err))
		return
	}

	// Validate the VIP fields that were provided
	valid, err := ra.validateVIPFields(v)
	if !valid {
		WriteErrorJSON(resp, http.StatusNotAcceptable, err)
		return
	}

	// Let's check to see if a VIP by this name already exists in the datbase
	_, err = ra.m.S.GetVIP(v.Name)
	if err == nil {
		WriteErrorJSON(resp, http.StatusNotAcceptable, fmt.Errorf("VIP %v already exists.  Please choose a new name.", v.Name))
		return
	}

	// Validate the Kubernetes service and port names that were provided
	valid, err = ra.m.K.VerifyKubeService(v)
	if !valid {
		WriteErrorJSON(resp, http.StatusNotAcceptable, err)
		return
	}

	// Check with cidrd to make sure this classname exists
	// Fetch an IP from cidrd and store UUID in VIP struct

	// Make sure a valid FE port was specified
	// Make sure a valid FE protocol was specified

	// Send some signals via channel

	// Everything looks good so we store the VIP in the database
	err = ra.m.S.SetVIP(v)
	if err != nil {
		WriteErrorJSON(resp, http.StatusInternalServerError, fmt.Errorf("Error storing VIP %v: %v", v.Name, err))
		return
	}

	// Return a 200 OK
	WriteSuccessJSON(resp, http.StatusOK, fmt.Sprintf("VIP %v created successfully.", v.Name))

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
	//	KubeSvcPortName    string       `json:"kube_service_port_name"`
	// }

}

func (ra *RestAPI) UpdateVIP(req *restful.Request, resp *restful.Response) {
	v := new(model.VIP)
	err := req.ReadEntity(&v)
	if err != nil {
		WriteErrorJSON(resp, http.StatusNotAcceptable, fmt.Errorf("Error parsing VIP entity: %v", err))
		return
	}

	// Validate the VIP fields that were provided
	valid, err := ra.validateVIPFields(v)
	if !valid {
		WriteErrorJSON(resp, http.StatusNotAcceptable, err)
		return
	}

	// First, we check to see if this VIP actually exists...
	_, err = ra.m.S.GetVIP(v.Name)
	if err != nil {
		WriteErrorJSON(resp, http.StatusNotFound, fmt.Errorf("VIP %v does not exist: %v", v.Name, err))
		return
	}

	// Validate the Kubernetes service and port names that were provided
	valid, err = ra.m.K.VerifyKubeService(v)
	if !valid {
		WriteErrorJSON(resp, http.StatusNotAcceptable, err)
		return
	}

	// The VIP is there, so we overwrite it with the VIP we were passed
	err = ra.m.S.SetVIP(v)
	if err != nil {
		WriteErrorJSON(resp, http.StatusInternalServerError, fmt.Errorf("Error storing VIP %v: %v", v.Name, err))
		return
	} else {

	}

	// Return a 200 OK
	WriteSuccessJSON(resp, http.StatusOK, fmt.Sprintf("VIP %v updated successfully.", v.Name))
	return
}

func (ra *RestAPI) DeleteVIP(req *restful.Request, resp *restful.Response) {
	vipid := req.PathParameter("vipid")

	// Make sure the VIP exists before we attempt to delete it.
	_, err := ra.m.S.GetVIP(vipid)
	if err != nil {
		WriteErrorJSON(resp, http.StatusNotFound, fmt.Errorf("VIP %v not found, cannot be deleted: %v", vipid, err))
		return
	}

	// Delete the VIP
	err = ra.m.S.DeleteVIP(vipid)
	if err != nil {
		WriteErrorJSON(resp, http.StatusInternalServerError, fmt.Errorf("VIP %v could not be deleted: %v", vipid, err))
		return
	} else {
		// Return a 200 OK
		WriteSuccessJSON(resp, http.StatusOK, fmt.Sprintf("VIP %v deleted successfully.", vipid))
	}

}

func (ra *RestAPI) GetVIP(req *restful.Request, resp *restful.Response) {
	vipid := req.PathParameter("vipid")

	// Get and respond with the VIP from our store or throw an error if
	// it doesn't exist
	vip, err := ra.m.S.GetVIP(vipid)
	if err != nil {
		WriteErrorJSON(resp, http.StatusNotFound, fmt.Errorf("VIP %v not found: %v", vipid, err))
		return
	}
	resp.PrettyPrint(false)
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

func (ra *RestAPI) validateVIPFields(v *model.VIP) (bool, error) {
	if v.Name == "" {
		return false, fmt.Errorf("VIP name cannot be empty.")
	}

	// Make sure frontend TCP/UDP port is valid
	if v.FrontendPort < 1 || v.FrontendPort > 65535 {
		return false, fmt.Errorf("Invalid VIP frontend port: %v", v.FrontendPort)
	}

	if v.KubeSvcName == "" || v.KubeSvcPortName == "" {
		return false, fmt.Errorf("VIP's Kubernetes service name and port name cannot be empty.")
	}

	// If a K8S namespace was not provided, just use the default namespace
	if v.KubeNamespace == "" {
		v.KubeNamespace = api.NamespaceDefault
	}

	return true, nil
}

func WriteErrorJSON(resp *restful.Response, respCode int, err error) {
	er := APIResponse{}

	fmt.Println("Error:", respCode, err.Error())

	er.Message = err.Error()

	content, _ := json.Marshal(er)

	resp.ResponseWriter.WriteHeader(respCode)
	resp.Write(content)
}

func WriteSuccessJSON(resp *restful.Response, respCode int, result string) {
	s := APIResponse{}

	s.Message = result
	s.Success = true

	content, _ := json.Marshal(s)

	resp.ResponseWriter.WriteHeader(respCode)
	resp.Write(content)
}
