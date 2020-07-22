package util

import (
	"log"

	consulapi "github.com/hashicorp/consul/api"
)

var ConsulClinet *consulapi.Client

func init() {
	config := consulapi.DefaultConfig()
	config.Address = "192.168.31.82:8500"

	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	ConsulClinet = client
}

func RegisterService() {
	reg := consulapi.AgentServiceRegistration{}
	reg.ID = "userservice1"
	reg.Name = "userservice"
	reg.Address = "192.168.31.82"
	reg.Port = 9050
	reg.Tags = []string{"primary", "v1"}

	check := consulapi.AgentServiceCheck{}
	check.Interval = "10s"
	check.Timeout = "5s"
	check.HTTP = "http://192.168.31.82:9050/health"

	reg.Check = &check

	err := ConsulClinet.Agent().ServiceRegister(&reg)
	if err != nil {
		log.Fatal(err)
	}
}

func DeregisterService() {
	ConsulClinet.Agent().ServiceDeregister("userservice1")
}
