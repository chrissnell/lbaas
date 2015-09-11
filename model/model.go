package model

import (
	"encoding/json"

	"github.com/coreos/etcd/client"

	"github.com/chrissnell/lbaas/config"
)

type VIP struct {
	Name               string       `json:"name"`
	FrontEndIPClass    string       `json:"frontend_ip_class"`
	FrontendIP         string       `json:"frontend_ip"`
	FrontendIPUUID     string       `json:"frontend_ip_uuid"`
	FrontendPort       uint8        `json:"frontend_port"`
	FrontendProtocol   string       `json:"frontend_protocol"`
	FrontendProfiles   []string     `json:"frontend_profiles"`
	PoolMembers        []PoolMember `json:"pool_members"`
	PoolMemberProtocol string       `json:"pool_member_protocols"` // HTTP, HTTPS, UDP, FTP, etc
	KubeSvcName        string       `json:"kube_service_name"`
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
}

// Model contains the data model with the associated etcd Client
type Model struct {
	c  config.Config
	LB LoadBalancer
	S  *Store
}

// New creates a new data model with a new DB connection
func New(e client.Client, lb LoadBalancer, c config.Config) *Model {

	s := &Store{}
	s = s.New(e, c)

	m := &Model{
		S:  s,
		LB: lb.(LoadBalancer),
		c:  c,
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
