package model

import (
	"github.com/coreos/go-etcd/etcd"
)

// Model contains the data model with the associated etcd Client
type Model struct {
	e *etcd.Client
}

// New creates a new data model with the given DB connection handle
func New(e *etcd.Client) *Model {
	m := &Model{
		e: e,
	}

	return m
}
