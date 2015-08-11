package etcd

import (
	"fmt"

	"github.com/coreos/go-etcd/etcd"
)

type Etcd struct {
	Client *etcd.Client
}

func New(hostname, port string) *Etcd {
	e := &Etcd{}
	e.Open(hostname, port)
	return e
}

func (e *Etcd) Open(hostname, port string) {
	var eh []string
	eh = append(eh, fmt.Sprintf("%v:%v", hostname, port))
	e.Client = etcd.NewClient(eh)
	return
}

func (e *Etcd) Close() {
	e.Client.Close()
	return
}
