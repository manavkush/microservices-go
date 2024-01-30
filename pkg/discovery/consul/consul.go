package consul

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	consul "github.com/hashicorp/consul/api"
)

// Registery defines a consul based service registry
type Registery struct {
	client *consul.Client
}

// NewRegistery creates a new consul registry instance.
// Addr is the address where the consul agent is running
func NewRegistery(addr string) (*Registery, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Registery{client: client}, nil
}

// Register creates a service record in the registry
func (r *Registery) Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("invalid host:port format. Eg: localhost:8081")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	host := parts[0]

	err = r.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		Address: host,
		Port:    port,
		ID:      instanceID,
		Name:    serviceName,
		Check: &consul.AgentServiceCheck{
			CheckID: instanceID,
			TTL:     "5s",
		},
	})
	return err
}

// Deregister removes a service record from the registry
func (r *Registery) Deregister(ctx context.Context, instanceID string, _ string) error {
	err := r.client.Agent().ServiceDeregister(instanceID)
	return err
}

// HealthCheck is a push mechanism to update the health status of a service instance
func (r *Registery) HealthCheck(instanceID string, _ string) error {
	err := r.client.Agent().UpdateTTL(instanceID, "", "pass")
	return err
}

// Discover returns a list of addresses of active instances of the given service
func (r *Registery) Discover(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}
	var instances []string
	for _, entry := range entries {
		instances = append(instances, fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port))
	}
	return instances, nil
}
