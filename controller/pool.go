package controller

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/chrissnell/lbaas/model"
)

// GeneratePoolMembers returns a list of pool members for a given Service
func GeneratePoolMembers(m *model.Model, v *model.VIP) (map[string]string, error) {
	var matched bool
	var iters, poolMemberCount uint16

	poolMembers := make(map[string]string)
	nodeUIDMap := make(map[uint16]string)
	nodeIPMap := make(map[uint16]string)

	ks, err := m.K.GetKubeService(v.KubeSvcName, v.KubeNamespace)
	if err != nil {
		return nil, err
	}

	nl, err := m.K.GetAllKubeNodes()

	clusterSize := uint16(len(nl.Items))

	for _, node := range nl.Items {
		nodeInt, _ := strconv.ParseUint(string(node.UID)[:3], 16, 32)
		nodeUIDMap[uint16(nodeInt)] = string(node.UID)
		nodeIPMap[uint16(nodeInt)] = node.Status.Addresses[0].Address
	}

	uid := string(ks.UID)
	svcSeed, err := strconv.ParseUint(uid[:8], 16, 32)
	if err != nil {
		log.Printf("Error parsing UID:", err)
		return nil, err
	}

	if clusterSize < m.C.LoadBalancer.PoolMembersPerVIP {
		poolMemberCount = clusterSize
	} else {
		poolMemberCount = m.C.LoadBalancer.PoolMembersPerVIP
	}

	rand.Seed(int64(svcSeed))

	start := time.Now()
	var i uint16
	for i = 0; i < poolMemberCount; i++ {
		matched = false
		iters = 0
		for !matched {
			iters++
			rn := uint16(rand.Int31n(4095))
			_, presentInPool := poolMembers[nodeUIDMap[rn]]
			if nodeUIDMap[rn] != "" && !presentInPool {
				log.Printf("[iters: %v] Found a node: %v\n", iters, nodeUIDMap[rn])
				poolMembers[nodeUIDMap[rn]] = nodeIPMap[rn]
				matched = true
			}
		}
	}
	log.Printf("Time to find %v pool members: %v\n", poolMemberCount, time.Since(start))

	return poolMembers, nil
}
