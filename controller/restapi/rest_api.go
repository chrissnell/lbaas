package restapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/chrissnell/lbaas/model"
	"github.com/gorilla/mux"
)

type RestAPI struct {
}

func (ra *RestAPI) CreateVIP(w http.ResponseWriter, r *http.Request) {
}

func (ra *RestAPI) DeleteVIP(w http.ResponseWriter, r *http.Request) {
}

func (ra *RestAPI) UpdateVIP(w http.ResponseWriter, r *http.Request) {
}

func (ra *RestAPI) GetVIP(w http.ResponseWriter, r *http.Request) {
	var res model.VIP

	vars := mux.Vars(r)
	vipid := vars["vipid"]

	res, err := model.GetVIP(vipid)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		res.Name = fmt.Sprint("VIP ", vipid, " does not exist")

		log.Printf("struct: %+v\n", res)
		json.NewEncoder(w).Encode(res)
	}

	log.Println("GOT", vipid)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	log.Printf("struct: %+v\n", res)
	json.NewEncoder(w).Encode(res)
}

func (ra *RestAPI) GetAllVIPs(w http.ResponseWriter, r *http.Request) {
}

func (ra *RestAPI) CreatePool(w http.ResponseWriter, r *http.Request) {
}

func (ra *RestAPI) DeletePool(w http.ResponseWriter, r *http.Request) {
}

func (ra *RestAPI) AddPoolMembers(w http.ResponseWriter, r *http.Request) {
}

func (ra *RestAPI) DeletePoolMember(w http.ResponseWriter, r *http.Request) {
}

func (ra *RestAPI) GetAllPoolMembers(w http.ResponseWriter, r *http.Request) {
}

func (ra *RestAPI) DeleteAllPoolMembers(w http.ResponseWriter, r *http.Request) {
}
