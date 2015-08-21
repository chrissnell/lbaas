package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/ugorji/go/codec"
)

type VIP struct {
	Name             string
	FrontendIP       string
	FrontendPort     uint8
	FrontendProtocol string
	PoolMembers      []PoolMember
}

// Marshal implements the json Encoder interface
func (v *VIP) Marshal() ([]byte, error) {
	jv, err := json.Marshal(&v)
	return jv, err
}

// Unmarshal implements the json Decoder interface
func (v *VIP) Unmarshal(jv string) error {
	err := json.Unmarshal([]byte(jv), &v)
	return err
}

// GetVIP will fetch a VIP from the database
func (m *Model) GetVIP(v string) (*VIP, error) {
	var db []byte
	var h codec.Handle = new(codec.JsonHandle)

	var dv *VIP

	enc := codec.NewEncoderBytes(&db, h)

	vresp, err := m.SafeGet(fmt.Sprint(m.c.Etcd.BasePath, "/vips/", v), true, false)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error fetching /vips/%v : %v\n", v, err))
	}

	vresp.CodecEncodeSelf(enc)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error encoding response from etcd:", err))
	}

	log.Println("Returned JSON:", vresp)
	return dv, nil
}
