package model

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/chrissnell/lbaas/config"

	"github.com/coreos/go-etcd/etcd"
)

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
	e  *etcd.Client
	lb LoadBalancer
	c  config.Config
}

// New creates a new data model with a new DB connection
func New(e *etcd.Client, lb LoadBalancer, c config.Config) *Model {
	m := &Model{
		e:  e,
		lb: lb,
		c:  c,
	}

	return m
}

func (m *Model) SafeGet(key string, sort, recursive bool) (*etcd.Response, error) {
	// Test for rude boys
	r, _ := regexp.Compile("../")
	if r.MatchString(key) {
		return nil, errors.New(fmt.Sprint("Invalid key:", key))
	}

	return m.e.Get(key, sort, recursive)
}
