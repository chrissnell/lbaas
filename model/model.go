package model

import (
	"encoding/json"
	"log"

	"github.com/chrissnell/lbaas/config"
)

type VIP struct {
	Name               string       `json:"name"`
	FrontEndIPClass    string       `json:"frontend_ip_class"`
	FrontendIP         string       `json:"frontend_ip"`
	FrontendIPUUID     string       `json:"frontend_ip_uuid"`
	FrontendPort       uint16       `json:"frontend_port"`
	FrontendProtocol   string       `json:"frontend_protocol"`
	FrontendProfiles   []string     `json:"frontend_profiles"`
	PoolMembers        []PoolMember `json:"pool_members"`
	PoolMemberProtocol string       `json:"pool_member_protocols"` // HTTP, HTTPS, UDP, FTP, etc
	KubeNamespace      string       `json:"kube_namespace"`
	KubeSvcName        string       `json:"kube_service_name"`
	KubeSvcPortName    string       `json:"kube_service_port_name"`
}

type LoadBalancer interface {
	CreateVIP(*VIP) error
	UpdateVIP(*VIP) error
	DeleteVIP(string) error
	GetVIP(string) (*VIP, error)
	GetAllVIPs() ([]*VIP, error)
	AddPoolMembers([]*PoolMember) error
	DeletePoolMember(string) error
	DeleteAllPoolMembers() error
	GetPoolMembers() ([]*PoolMember, error)
	ValidateProtocol(string) bool
}

// Model contains the data model with the associated etcd Client
type Model struct {
	c                  config.Config
	LB                 LoadBalancer
	S                  *Store
	K                  *Kube
	KubeWorkQueueReady chan struct{}
}

// New creates a new data model with a new DB connection and Kube API client
func New(lb LoadBalancer, c config.Config) *Model {

	KubeWorkQueueReady := make(chan struct{})

	s := &Store{}
	s = s.New(c)

	k := &Kube{}
	k, err := k.New(c, KubeWorkQueueReady)
	if err != nil {
		log.Fatalln("Error creating Kubernetes client:", err)
	}

	m := &Model{
		S:                  s,
		K:                  k,
		LB:                 lb.(LoadBalancer),
		c:                  c,
		KubeWorkQueueReady: KubeWorkQueueReady,
	}

	return m
}

// Marshal implements the json Encoder interface
func (v *VIP) Marshal() ([]byte, error) {
	jv, err := json.Marshal(&v)
	return jv, err
}

// Unmarshal implements the json Decoder interface
func (v *VIP) Unmarshal(jv string) error {
	err := json.Unmarshal([]byte(jv), &v)
	return err
}
