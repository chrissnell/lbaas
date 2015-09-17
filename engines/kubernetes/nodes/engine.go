package nodesengine

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/chrissnell/lbaas/model"
	"github.com/chrissnell/lbaas/util/log"
)

const (
	nodeAdded NodeStatus = iota
	nodeDeleted
)

// Bryan's idea: goroutines for each node with a pub-sub channel (tv42's topic) to broadcast when nodes go away

type NodeStatus int

type NodeChangeMessage struct {
	UID    string
	Action NodeStatus
}

type Engine struct {
	sync.Mutex
	m              *model.Model
	activeNodes    map[string]string // node_UID -> node_IP
	NodeChangeChan chan NodeChangeMessage
}

func New(m *model.Model) *Engine {
	e := &Engine{
		m:              m,
		activeNodes:    make(map[string]string),
		NodeChangeChan: make(chan NodeChangeMessage, 100),
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
			logger.Log("The nodes engine is ticking...")

			nl, err := e.m.K.GetAllKubeNodes("")
			if err != nil {
				logger.Log(fmt.Sprintln("Could not get all nodes:", err))
			}

			for _, i := range nl.Items {
				log.Println("Node:", i.Name, i.ObjectMeta.UID)
			}

			msg := NodeChangeMessage{
				UID:    "12345",
				Action: nodeAdded,
			}
			e.NodeChangeChan <- msg
		}
	}

}

func (e *Engine) addNode(uid, ip string) error {
	e.Lock()
	defer e.Unlock()

	return nil
}

func (e *Engine) removeNode(uid string) error {
	e.Lock()
	defer e.Unlock()

	return nil
}

func (e *Engine) RefreshAllServices() error {
	e.Lock()
	defer e.Unlock()

	return nil
}
