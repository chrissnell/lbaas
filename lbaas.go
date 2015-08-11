package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"github.com/chrissnell/lbaas/config"
	"github.com/chrissnell/lbaas/controller"
	"github.com/chrissnell/lbaas/etcd"
	"github.com/chrissnell/lbaas/loadbalancer"
	"github.com/chrissnell/lbaas/model"
)

// Service contains our configuration and runtime objects
type Service struct {
	Config config.Config
	LB     *loadbalancer.LoadBalancer
	etcd   *etcd.Etcd
	api    *controller.Controller
	model  *model.Model
}

// New creates a new instance of Service with the given configuration file
func New(filename string) *Service {
	s := &Service{}

	// Read our server configuration
	filename, _ = filepath.Abs(filename)
	cfg, err := config.New(filename)
	if err != nil {
		log.Fatalln("Error reading config file.  Did you pass the -config flag?  Run with -h for help.\n", err)
	}
	s.Config = cfg

	// Open an etcd client
	s.etcd = etcd.New(s.Config.Etcd.Hostname, s.Config.Etcd.Port)
	defer s.etcd.Close()

	// Initialize the data model
	// s.model = model.New(s.db)

	// Initialize the Controller
	// s.api = controller.New(s.Config, s.model, s.notify, r)

	return s
}

// Listen will start the HTTP listeners for the API router.
func (s *Service) Listen() {

	// Set up our API endpoint router
	log.Fatal(http.ListenAndServe(s.Config.Service.APIListenAddr, s.api.APIRouter()))
}

// Close will shut down the service
func (s *Service) Close() {
	// s.db.Close()
}

func main() {
	cfgFile := flag.String("config", "config.yaml", "Path to config file (default: ./config.yaml)")
	flag.Parse()

	s := New(*cfgFile)
	defer s.Close()
	s.Listen()
}
