package controller

import (
	"fmt"
	"log"
	"sync"

	"github.com/chrissnell/lbaas/model"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/controller/framework"
	"k8s.io/kubernetes/pkg/watch"
)

type ServiceStatus int

type ServiceChangeMessage struct {
	UID         string
	Event       *api.Service
	EventType   watch.EventType
	ServiceName string
}

type ServicesEngine struct {
	sync.Mutex
	m                 *model.Model
	activeNodes       map[string]string // node_UID -> node_IP
	ServiceChangeChan chan ServiceChangeMessage
}

func NewServicesEngine(m *model.Model) *ServicesEngine {
	e := &ServicesEngine{
		m:                 m,
		activeNodes:       make(map[string]string),
		ServiceChangeChan: make(chan ServiceChangeMessage, 100),
	}

	// Start up the engine
	go e.start()

	webby, _ := m.S.GetVIP("webby")
	pm, err := GeneratePoolMembers(m, webby)
	if err != nil {
		log.Println("Error generating pool members:", err)
	} else {
		for ud, ip := range pm {
			log.Printf("----> Pool member: %v [%v]", ud, ip)
		}
	}

	return e
}

func (e *ServicesEngine) start() {

	for {

		// keyFunc for endpoints and services.
		keyFunc := framework.DeletionHandlingMetaNamespaceKeyFunc

		for {
			item, _ := e.m.K.ServiceQueue.Get()
			ev := item.(model.QueueEvent).Obj
			evtype := item.(model.QueueEvent).ObjType
			key, _ := keyFunc(ev)

			log.Printf("SERVICE Sync triggered by  %v\n", key)
			log.Printf("---->  [%v] UID: %v", ev.(*api.Service).Name, ev.(*api.Service).UID)
			for _, p := range ev.(*api.Service).Spec.Ports {
				log.Printf("----> [%v] NodePort: %v\n", p.Name, p.Name)
			}
			log.Println("---->  Condition Type:", evtype)

			msg := ServiceChangeMessage{
				UID:         fmt.Sprint(ev.(*api.Service).UID),
				ServiceName: ev.(*api.Service).Name,
				Event:       ev.(*api.Service),
				EventType:   evtype,
			}

			e.ServiceChangeChan <- msg

			e.m.K.ServiceQueue.Done(ev)
		}

	}

}

func (e *ServicesEngine) addService(uid, ip string) error {
	e.Lock()
	defer e.Unlock()

	return nil
}

func (e *ServicesEngine) deleteService(uid string) error {
	e.Lock()
	defer e.Unlock()

	return nil
}
