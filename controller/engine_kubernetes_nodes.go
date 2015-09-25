package controller

import (
	"fmt"
	"log"
	"sync"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/controller/framework"
	"k8s.io/kubernetes/pkg/watch"

	"github.com/chrissnell/lbaas/model"
)

// Bryan's idea: goroutines for each node with a pub-sub channel (tv42's topic) to broadcast when nodes go away

type NodeChangeMessage struct {
	UID       string
	Event     watch.Event
	EventType watch.EventType
	NodeReady bool
	Hostname  string
}

type NodesEngine struct {
	sync.Mutex
	m              *model.Model
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

	// keyFunc for endpoints and services.
	keyFunc := framework.DeletionHandlingMetaNamespaceKeyFunc

	<-e.m.KubeWorkQueueReady

	for {
		if e.m.K.NodeQueue != nil {
			item, _ := e.m.K.NodeQueue.Get()
			ev := item.(model.QueueEvent).Obj
			evtype := item.(model.QueueEvent).ObjType
			key, _ := keyFunc(ev)

			log.Printf("NODE Sync triggered by  %v\n", key)
			log.Printf("---->  [%v] UID: %v", ev.(*api.Node).Status.Addresses[0].Address, ev.(*api.Node).UID)
			log.Println("---->  Status:", ev.(*api.Node).Status.Conditions[0].Status)
			log.Println("---->  Message:", ev.(*api.Node).Status.Conditions[0].Message)
			log.Println("--- >  Reason:", ev.(*api.Node).Status.Conditions[0].Reason)
			log.Println("---->  Condition Type:", ev.(*api.Node).Status.Conditions[0].Type)

			msg := NodeChangeMessage{
				UID:       fmt.Sprint(ev.(*api.Node).UID),
				Hostname:  ev.(*api.Node).Status.Addresses[0].Address,
				Event:     ev.(watch.Event),
				EventType: evtype,
			}

			if ev.(*api.Node).Status.Conditions[0].Status == api.ConditionTrue {
				msg.NodeReady = true
			} else {
				msg.NodeReady = false
			}

			e.NodeChangeChan <- msg

			e.m.K.NodeQueue.Done(ev)
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
