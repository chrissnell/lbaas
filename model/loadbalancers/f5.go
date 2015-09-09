package f5

import (
	"errors"
	"github.com/chrissnell/lbaas/model"
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

func (f F5) CreateVIP(*model.VIP) error {
	var err error
	return err
}

func (f F5) UpdateVIP(*model.VIP) error {
	var err error
	return err
}

func (f F5) DeleteVIP(string) error {
	var err error
	return err
}

func (f F5) GetVIP(string) (*model.VIP, error) {
	var err error
	var v model.VIP
	v.Name = "My Test VIP"
	err = errors.New("This is broken")
	return &v, err
}

func (f F5) GetAllVIPs() ([]*model.VIP, error) {
	var err error
	var v []*model.VIP
	return v, err
}

func (f F5) AddPoolMembers([]*model.PoolMember) error {
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

func (f F5) GetPoolMembers() ([]*model.PoolMember, error) {
	var err error
	var p []*model.PoolMember
	return p, err
}
