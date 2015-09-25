package controller

import (
	_ "fmt"
	"log"
	"sync"

	_ "k8s.io/kubernetes/pkg/api"

	"github.com/chrissnell/lbaas/model"
)

type LBEngine struct {
	sync.Mutex
	m              *model.Model
	activeNodes    map[string]string // node_UID -> node_IP
	nodeChangeChan chan NodeChangeMessage
}

func NewLBEngine(m *model.Model, nu chan NodeChangeMessage) *LBEngine {
	e := &LBEngine{
		m:              m,
		activeNodes:    make(map[string]string),
		nodeChangeChan: nu,
	}

	// Start up the engine
	go e.start()

	return e
}

func (e *LBEngine) start() {
	for {
		select {
		case ev := <-e.nodeChangeChan:
			// log.Printf("Received event: %+v\n", ev)

			if _, ok := e.activeNodes[ev.UID]; ok {
				// This node is already in activeNodes

				// Delete from activeNodes if the node is not ready
				if !ev.NodeReady {
					log.Println("Received a", ev.EventType, "message for", ev.UID)
					log.Println("   ---> DELETING from activeNodes")
					delete(e.activeNodes, ev.UID)
				} else {
					// If the node is ready and it's a MODIFY even, update the IP.
					log.Println("Received a", ev.EventType, "message for", ev.UID)
					log.Println("   ---> UPDATING activeNodes")
					e.activeNodes[ev.UID] = ev.Hostname
				}
			} else {
				// This node is not in activeNodes
				log.Println("Received a", ev.EventType, "message for", ev.UID)
				log.Println("   ---> ADDING to activeNodes")
				// So we add it to activeNodes
				e.activeNodes[ev.UID] = ev.Hostname
			}

			// log.Printf("Recieved event: [%v] [%v] %v\n", ev.Type, ev.Object.(*api.Node).UID, ev.Object.(*api.Node).Status.Conditions[0].Status)

			// if ev.Type == watch.Modified || ev.Type == watch.Added || ev.Type == watch.Deleted {

			// 	msg := NodeChangeMessage{
			// 		Type: ev.Type,
			// 		UID:  fmt.Sprint(ev.Object.(*api.Node).UID),
			// 		// Currently using the first address in the array...maybe we should send them all?
			// 		Hostname: ev.Object.(*api.Node).Status.Addresses[0].Address,
			// 		Event:    ev,
			// 	}

			// 	if ev.Object.(*api.Node).Status.Conditions[0].Status == api.ConditionTrue {
			// 		msg.NodeReady = true
			// 	} else {
			// 		msg.NodeReady = false
			// 	}

			// 	log.Printf("Sending NodeChangMessage: %+v\n", msg)

			// 	e.NodeChangeChan <- msg
			// }
		}
	}

}
