package controller

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/chrissnell/lbaas/model"
)

// GeneratePoolMembers uses a deterministic algorithm to pick a set of pool members
// for a given VIP
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

	// For each node, we grab the first three digits of the UID and convert it to a uint16
	// and map this uint16 to the node's full UID and the node's IP
	for _, node := range nl.Items {
		nodeInt, _ := strconv.ParseUint(string(node.UID)[:3], 16, 32)
		nodeUIDMap[uint16(nodeInt)] = string(node.UID)
		nodeIPMap[uint16(nodeInt)] = node.Status.Addresses[0].Address
	}

	// Take the first 8 digits of the service UID and convert to uint16
	uid := string(ks.UID)
	svcSeed, err := strconv.ParseUint(uid[:8], 16, 32)
	if err != nil {
		log.Printf("Error parsing UID:", err)
		return nil, err
	}

	// If our Kubernetes cluster is smaller than our desired number
	// of LB pool members, we'll use what we have.  Otherwise, we'll
	// use what we want.
	if clusterSize < m.C.LoadBalancer.PoolMembersPerVIP {
		poolMemberCount = clusterSize
	} else {
		poolMemberCount = m.C.LoadBalancer.PoolMembersPerVIP
	}

	// Create a RNG and seed it with the uint that we grabbed from our service UID
	rand.Seed(int64(svcSeed))

	// Use our RNG to generate an int between 0 and 4095 and look for a match
	// in our node UID lookup map.  If a match is found, add that node to the
	// pool for this VIP.
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
