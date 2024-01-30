package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"movieexample.com/pkg/discovery"
)

type serviceNameType string
type instanceIDType string

type Registery struct {
	sync.RWMutex
	serviceAddrs map[serviceNameType]map[instanceIDType]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

// NewRegistry creates a new in-memory registry instance.
func NewRegistry() *Registery {
	return &Registery{
		serviceAddrs: make(map[serviceNameType]map[instanceIDType]*serviceInstance),
	}
}

// Register creates a service record in the registry
func (r *Registery) Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceNameType(serviceName)]; !ok {
		r.serviceAddrs[serviceNameType(serviceName)] = map[instanceIDType]*serviceInstance{}
	}
	r.serviceAddrs[serviceNameType(serviceName)][instanceIDType(instanceID)] = &serviceInstance{
		hostPort:   hostPort,
		lastActive: time.Now(),
	}
	return nil
}

// Deregister removes a service record from the registry
func (r *Registery) Deregister(ctx context.Context, instanceID string, serviceName string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceNameType(serviceName)]; !ok {
		return discovery.ErrNotFound
	}

	delete(r.serviceAddrs[serviceNameType(serviceName)], instanceIDType(instanceID))
	return nil
}

// HealthCheck marks a service instance as active
func (r *Registery) HealthCheck(instanceID string, serviceName string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceNameType(serviceName)]; !ok {
		return errors.New("service not registered yet")
	}

	if _, ok := r.serviceAddrs[serviceNameType(serviceName)][instanceIDType(instanceID)]; !ok {
		return errors.New("service instance not registered yet")
	}

	r.serviceAddrs[serviceNameType(serviceName)][instanceIDType(instanceID)].lastActive = time.Now()
	return nil
}

// Discover returns a list of service instances from the registry
func (r *Registery) Discover(ctx context.Context, serviceName string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()

	if len(r.serviceAddrs[serviceNameType(serviceName)]) == 0 {
		return nil, discovery.ErrNotFound
	}
	var res []string

	for _, v := range r.serviceAddrs[serviceNameType(serviceName)] {
		if time.Since(v.lastActive) > 5*time.Second {
			continue
		}

		res = append(res, v.hostPort)
	}
	return res, nil
}
