package state

import (
	"fmt"
	"strings"
	"sync"
)

type InstanceMeta struct {
	InstanceID    string `yaml:"instance_id"`
	AccountID     string `yaml:"account_id"`
	ServerAddress string `yaml:"server_address"`
	Enabled       bool   `yaml:"enabled"`
}

type InstanceRepository struct {
	path string
	mu   sync.Mutex
}

func NewInstanceRepository(path string) *InstanceRepository {
	return &InstanceRepository{path: strings.TrimSpace(path)}
}

func (r *InstanceRepository) LoadAll() ([]InstanceMeta, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var instances []InstanceMeta
	if err := readYAMLFile(r.path, &instances); err != nil {
		return nil, err
	}
	if err := validateInstances(instances); err != nil {
		return nil, err
	}
	if len(instances) == 0 {
		return []InstanceMeta{}, nil
	}
	return instances, nil
}

func (r *InstanceRepository) SaveAll(instances []InstanceMeta) error {
	if err := validateInstances(instances); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	return writeYAMLAtomic(r.path, instances)
}

func validateInstances(instances []InstanceMeta) error {
	seen := make(map[string]struct{}, len(instances))
	for i, instance := range instances {
		instanceID := strings.TrimSpace(instance.InstanceID)
		if instanceID == "" {
			return fmt.Errorf("instance metadata[%d]: instance_id is required", i)
		}
		if _, ok := seen[instanceID]; ok {
			return fmt.Errorf("duplicate instance_id %q", instanceID)
		}
		seen[instanceID] = struct{}{}

		if strings.TrimSpace(instance.AccountID) == "" {
			return fmt.Errorf("instance metadata[%d]: account_id is required", i)
		}
		if strings.TrimSpace(instance.ServerAddress) == "" {
			return fmt.Errorf("instance metadata[%d]: server_address is required", i)
		}
	}
	return nil
}
