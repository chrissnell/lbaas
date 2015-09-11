package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/github.com/ugorji/go/codec"
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"

	"github.com/chrissnell/lbaas/config"
)

type Store struct {
	c config.Config
	e client.Client
	k client.KeysAPI
}

// FetchVIP(string) (*VIP, error)
// StoreVIP(*VIP) error
// DeleteVIP(string) error
// FetchAllVIPs() ([]*VIP, error)

func (s *Store) New(e client.Client, c config.Config) *Store {
	// Create a new KeysAPI for etcd
	k := client.NewKeysAPI(e)

	ns := &Store{
		e: e,
		c: c,
		k: k,
	}

	return ns
}

// GetVIP will fetch a VIP from the etcd
func (s *Store) FetchVIP(v string) (*VIP, error) {
	var db []byte
	var h codec.Handle = new(codec.JsonHandle)

	var dv *VIP

	dec := codec.NewDecoderBytes(db, h)

	er, err := s.SafeGet(fmt.Sprint(s.c.Etcd.BasePath, "/vips/", v), true, false)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error fetching /vips/%v : %v\n", v, err))
	}

	//er.CodecEncodeSelf(enc)
	er.CodecDecodeSelf(dec)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error decoding response from etcd:", err))
	}

	log.Println("Returned JSON:", er)
	return dv, nil
}

func (s *Store) StoreVIP(v *VIP) error {
	if v.Name == "" {
		return fmt.Errorf("Cannot store a VIP if name is not set.")
	}

	vs, err := json.Marshal(v)
	if err != nil {
		return err
	}

	opt := &client.SetOptions{
		Dir: false,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = s.k.Set(ctx, fmt.Sprint(s.c.Etcd.BasePath, "/vips/", v.Name), string(vs), opt)
	if err != nil {
		return err
	} else {
		return nil
	}

}

func (s *Store) SafeGet(key string, sort, recursive bool) (*client.Response, error) {

	opt := &client.GetOptions{
		Recursive: recursive,
		Sort:      sort,
		Quorum:    true,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test for rude boys
	r, _ := regexp.Compile("../")
	if r.MatchString(key) {
		return nil, errors.New(fmt.Sprint("Invalid key:", key))
	}

	return s.k.Get(ctx, key, opt)
}
