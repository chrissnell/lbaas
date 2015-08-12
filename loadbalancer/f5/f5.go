package f5

import (
	"github.com/chrissnell/lbaas/loadbalancer"
)

type F5 struct {
	ManagementIP   string
	ManagementPort uint8
	Username       string
	Password       string
}

func New() F5 {
	f := F5{}
	return f
}

func (f F5) CreateVIP(*loadbalancer.VIP) error {
	var err error
	return err
}

func (f F5) UpdateVIP(*loadbalancer.VIP) error {
	var err error
	return err
}

func (f F5) DeleteVIP(string) error {
	var err error
	return err
}

func (f F5) GetVIP(string) (*loadbalancer.VIP, error) {
	var err error
	var v *loadbalancer.VIP
	return v, err
}

func (f F5) AddPoolMember(*loadbalancer.PoolMember) error {
	var err error
	return err
}

func (f F5) DeletePoolMember(pm string) error {
	var err error
	return err
}

func (f F5) DeleteAllPoolMembers() error {
	var err error
	return err
}

func (f F5) GetPoolMembers() ([]*loadbalancer.PoolMember, error) {
	var err error
	var p []*loadbalancer.PoolMember
	return p, err
}
