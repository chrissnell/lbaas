package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"

	"github.com/chrissnell/lbaas/config"
)

type Store struct {
	c        config.Config
	e        client.Client
	k        client.KeysAPI
	basePath string
}

// FetchVIP(string) (*VIP, error)
// StoreVIP(*VIP) error
// DeleteVIP(string) error
// FetchAllVIPs() ([]*VIP, error)

func (s *Store) New(c config.Config) *Store {
	// Open an etcd client
	var endpoints []string
	endpoints = append(endpoints, fmt.Sprintf("http://%v:%v", c.Etcd.Hostname, c.Etcd.Port))

	ec := client.Config{
		Endpoints: endpoints,
		Transport: client.DefaultTransport,
	}

	e, err := client.New(ec)
	if err != nil {
		log.Fatalln("Could not connect to etcd:", err)
	}

	// Create a new KeysAPI for etcd
	k := client.NewKeysAPI(e)

	ns := &Store{
		e:        e,
		c:        c,
		k:        k,
		basePath: fmt.Sprint(c.Etcd.BasePath, "/vips/"),
	}

	return ns
}

func (s *Store) keyPath(k string) string {
	return fmt.Sprint(s.basePath, k)
}

// GetVIP will fetch a VIP from etcd
func (s *Store) GetVIP(v string) (*VIP, error) {
	var dv *VIP

	er, err := s.SafeGet(s.keyPath(v), true, false)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error fetching /vips/%v : %v", v, err))
	}

	err = json.Unmarshal([]byte(er.Node.Value), &dv)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error decoding response from etcd:", err))
	}

	return dv, nil
}

func (s *Store) SetVIP(v *VIP) error {
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

	_, err = s.k.Set(ctx, s.keyPath(v.Name), string(vs), opt)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s *Store) DeleteVIP(v string) error {

	opt := &client.DeleteOptions{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.k.Delete(ctx, s.keyPath(v), opt)
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
	r, _ := regexp.Compile(`\.\./`)
	if r.MatchString(key) {
		return nil, errors.New(fmt.Sprint("Invalid key:", key))
	}

	return s.k.Get(ctx, key, opt)
}
