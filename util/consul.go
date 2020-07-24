package util

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	consulapi "github.com/hashicorp/consul/api"
)

var ConsulClinet *consulapi.Client

var ServiceID string
var ServiceName string
var ServicePort int

func init() {
	config := consulapi.DefaultConfig()
	config.Address = "192.168.31.82:8500"

	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	ConsulClinet = client

	ServiceID = "userservice" + uuid.New().String()
}

func SetServiceNameAndPort(name string, port int) {
	ServiceName = name
	ServicePort = port
}

func RegisterService() {
	reg := consulapi.AgentServiceRegistration{}
	reg.ID = ServiceID
	reg.Name = ServiceName
	reg.Address = "192.168.31.82"

	reg.Port = ServicePort
	reg.Tags = []string{"primary", "v1"}

	check := consulapi.AgentServiceCheck{}
	check.Interval = "10s"
	check.Timeout = "5s"
	check.HTTP = fmt.Sprintf("http://%s:%d/health", reg.Address, ServicePort)

	reg.Check = &check

	err := ConsulClinet.Agent().ServiceRegister(&reg)
	if err != nil {
		log.Fatal(err)
	}
}

func DeregisterService() {
	ConsulClinet.Agent().ServiceDeregister(ServiceID)
}
