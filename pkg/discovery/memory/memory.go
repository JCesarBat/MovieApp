package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"movieexample.com/pkg/discovery"
)

type serviceName string
type InstanceID string

// Registry defines an in-mempry service registry.
type Registry struct {
	sync.RWMutex
	serviceAddrs map[serviceName]map[InstanceID]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

// NewRegistry creates a new in-memory service
// registry instance.
func New() *Registry {
	return &Registry{
		serviceAddrs: make(map[serviceName]map[InstanceID]*serviceInstance),
	}
}

// Register creates a service record in the registry
func (r *Registry) Register(_ context.Context, instanceId string,
	Name string, hostPort string) error {
	r.Lock()
	defer r.Unlock()
	sName := serviceName(Name)
	instID := InstanceID(instanceId)

	if _, ok := r.serviceAddrs[sName]; !ok {
		r.serviceAddrs[sName] = make(map[InstanceID]*serviceInstance)
	}
	r.serviceAddrs[sName][instID] = &serviceInstance{
		hostPort:   hostPort,
		lastActive: time.Now(),
	}
	return nil
}

// Deregister removes a service record from the
// registry.
func (r *Registry) Deregister(ctx context.Context, instanceId string,
	Name string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName(Name)]; !ok {
		return nil
	}
	delete(r.serviceAddrs[serviceName(Name)], InstanceID(instanceId))
	return nil

}

// ReportHealthyState is a push mechanism for
// reporting healty state to the registry.
func (r *Registry) ReportHealthyState(instanceId string,
	Name string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName(Name)]; !ok {
		return errors.New("service is not registered yet")
	}
	if _, ok := r.serviceAddrs[serviceName(Name)][InstanceID(instanceId)]; !ok {
		return errors.New("service instance is not registered yet")
	}
	r.serviceAddrs[serviceName(Name)][InstanceID(instanceId)].lastActive =
		time.Now()
	return nil

}

// ServiceAddresses  return a list of addresses of
// active instances of the given servic.
func (r *Registry) ServiceAddresses(ctx context.Context, name string) ([]string, error) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName(name)]; !ok {
		return nil, discovery.ErrNotFound
	}
	if len(r.serviceAddrs[serviceName(name)]) == 0 {
		return nil, discovery.ErrNotFound
	}
	var list []string
	for _, i := range r.serviceAddrs[serviceName(name)] {
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		list = append(list, i.hostPort)
	}

	return list, nil
}
