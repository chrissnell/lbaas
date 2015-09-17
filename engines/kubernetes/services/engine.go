package servicesengine

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/chrissnell/lbaas/model"
	"github.com/chrissnell/lbaas/util/log"
)

type ServiceStatus int

const (
	serviceAdded ServiceStatus = iota
	serviceDeleted
)

type ServiceChangeMessage struct {
	UID    string
	Action ServiceStatus
}

type Engine struct {
	sync.Mutex
	m                 *model.Model
	activeNodes       map[string]string // node_UID -> node_IP
	ServiceChangeChan chan ServiceChangeMessage
}

func New(m *model.Model) *Engine {
	e := &Engine{
		m:                 m,
		activeNodes:       make(map[string]string),
		ServiceChangeChan: make(chan ServiceChangeMessage, 100),
	}

	// Start up the engine
	go e.start()

	return e
}

func (e *Engine) start() {
	ticker := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-ticker.C:
			logger.Log("The services engine is ticking...")

			sl, err := e.m.K.GetAllKubeServices("")
			if err != nil {
				logger.Log(fmt.Sprintln("Could not get all services:", err))
			}

			for _, i := range sl.Items {
				log.Println("Service:", i.Name, i.ObjectMeta.UID)
			}

			msg := ServiceChangeMessage{
				UID:    "12345",
				Action: serviceAdded,
			}
			e.ServiceChangeChan <- msg
		}
	}

}

func (e *Engine) addService(uid, ip string) error {
	e.Lock()
	defer e.Unlock()

	return nil
}

func (e *Engine) deleteService(uid string) error {
	e.Lock()
	defer e.Unlock()

	return nil
}
