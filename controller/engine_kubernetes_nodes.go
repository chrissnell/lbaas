package controller

import (
	"fmt"
	"log"
	"sync"
	"time"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/watch"

	"github.com/chrissnell/lbaas/model"
	"github.com/chrissnell/lbaas/util/log"
)

// Bryan's idea: goroutines for each node with a pub-sub channel (tv42's topic) to broadcast when nodes go away

type NodeChangeMessage struct {
	Type      watch.EventType
	UID       string
	Event     watch.Event
	NodeReady bool
	Hostname  string
}

type NodesEngine struct {
	sync.Mutex
	m              *model.Model
	w              watch.Interface
	activeNodes    map[string]string // node_UID -> node_IP
	NodeChangeChan chan NodeChangeMessage
}

func NewNodesEngine(m *model.Model) *NodesEngine {
	e := &NodesEngine{
		m:              m,
		activeNodes:    make(map[string]string),
		NodeChangeChan: make(chan NodeChangeMessage, 100),
	}

	// Start up the engine
	go e.start()

	return e
}

func (e *NodesEngine) start() {
	var err error
	e.w, err = e.m.K.NewKubeNodesWatcher(api.NamespaceDefault, nil)
	if err != nil {
		// This needs to reconnect...
		log.Println("Unable to get a Nodes watcher:", err)
	}

	ticker := time.NewTicker(time.Second * 5)

	for {
		select {
		case ev := <-e.w.ResultChan():

			if ev.Type == watch.Modified || ev.Type == watch.Added || ev.Type == watch.Deleted {

				msg := NodeChangeMessage{
					Type: ev.Type,
					UID:  fmt.Sprint(ev.Object.(*api.Node).UID),
					// Currently using the first address in the array...maybe we should send them all?
					Hostname: ev.Object.(*api.Node).Status.Addresses[0].Address,
					Event:    ev,
				}

				if ev.Object.(*api.Node).Status.Conditions[0].Status == api.ConditionTrue {
					msg.NodeReady = true
				} else {
					msg.NodeReady = false
				}

				// log.Printf("Sending NodeChangMessage: %+v\n", msg)

				e.NodeChangeChan <- msg
			}
		case <-ticker.C:
			logger.Log("The nodes engine is ticking...")

			nl, err := e.m.K.GetAllKubeNodes("")
			if err != nil {
				logger.Log(fmt.Sprintln("Could not get all nodes:", err))
			}

			for _, i := range nl.Items {
				log.Println("Node:", i.Name, i.ObjectMeta.UID)
			}

		}
	}

}

func (e *NodesEngine) addNode(uid, ip string) error {
	e.Lock()
	defer e.Unlock()

	return nil
}

func (e *NodesEngine) removeNode(uid string) error {
	e.Lock()
	defer e.Unlock()

	return nil
}
