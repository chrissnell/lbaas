package model

import (
	"github.com/coreos/go-etcd/etcd"
)

// Model contains the data model with the associated etcd Client
type Model struct {
	e *etcd.Client
}

// New creates a new data model with a new DB connection
func New(e *etcd.Client) *Model {
	m := &Model{
		e: e,
	}

	return m
}
