package model

import (
	"encoding/json"
)

type PoolMember struct {
	IP       string
	Port     uint8
	Protocol string
	// HealthCheck HealthCheck  - Not Yet Implemented
}

// Marshal implements the json Encoder interface
func (pm *PoolMember) Marshal() ([]byte, error) {
	jpm, err := json.Marshal(&pm)
	return jpm, err
}

// Unmarshal implements the json Decoder interface
func (pm *PoolMember) Unmarshal(jpm string) error {
	err := json.Unmarshal([]byte(jpm), &pm)
	return err
}
