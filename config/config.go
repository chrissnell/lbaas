package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Service      ServiceConfig      `yaml:"service"`
	LoadBalancer LoadBalancerConfig `yaml:"load_balancer"`
	Etcd         EtcdConfig         `yaml:"etcd"`
}

type ServiceConfig struct {
	APIListenAddr           string `yaml:"api_listen_address"`
	KubernetesAPIListenAddr string `yaml:"kubernetes_api_listen_address"`
}

type LoadBalancerConfig struct {
	Kind     string   `yaml:"kind"`
	F5Config F5Config `yaml:"f5,omitempty"`
}

type F5Config struct {
	IControlRESTBaseURL string `yaml:"iControl_REST_base_URL"`
	Username            string `yaml:"username"`
	Password            string `yaml:"password"`
}

type EtcdConfig struct {
	Hostname string `yaml:"hostname"`
	Port     string `yaml:"port"`
}

// New creates an new config object from the given filename.
func New(filename string) (Config, error) {
	cfgFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}
	c := Config{}
	err = yaml.Unmarshal(cfgFile, &c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}
