package state

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

type AccountMeta struct {
	AccountID string `yaml:"account_id"`
	Enabled   bool   `yaml:"enabled"`
	Label     string `yaml:"label,omitempty"`
	Note      string `yaml:"note,omitempty"`
}

type AccountRepository struct {
	path string
	mu   sync.Mutex
}

func NewAccountRepository(path string) *AccountRepository {
	return &AccountRepository{path: strings.TrimSpace(path)}
}

func (r *AccountRepository) LoadAll() ([]AccountMeta, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var accounts []AccountMeta
	if err := readYAMLFile(r.path, &accounts); err != nil {
		return nil, err
	}
	if err := validateAccounts(accounts); err != nil {
		return nil, err
	}
	if len(accounts) == 0 {
		return []AccountMeta{}, nil
	}
	return accounts, nil
}

func (r *AccountRepository) SaveAll(accounts []AccountMeta) error {
	if err := validateAccounts(accounts); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	return writeYAMLAtomic(r.path, accounts)
}

func validateAccounts(accounts []AccountMeta) error {
	seen := make(map[string]struct{}, len(accounts))
	for i, account := range accounts {
		accountID := strings.TrimSpace(account.AccountID)
		if accountID == "" {
			return fmt.Errorf("account metadata[%d]: account_id is required", i)
		}
		if _, ok := seen[accountID]; ok {
			return fmt.Errorf("duplicate account_id %q", accountID)
		}
		seen[accountID] = struct{}{}
	}
	return nil
}

func readYAMLFile(path string, out any) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return fmt.Errorf("state path is empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read state file: %w", err)
	}
	if len(data) == 0 {
		return nil
	}
	if err := yaml.Unmarshal(data, out); err != nil {
		return fmt.Errorf("decode state file: %w", err)
	}
	return nil
}

func writeYAMLAtomic(path string, value any) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return fmt.Errorf("state path is empty")
	}

	data, err := yaml.Marshal(value)
	if err != nil {
		return fmt.Errorf("encode state file: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create state dir: %w", err)
	}

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o600); err != nil {
		return fmt.Errorf("write state temp file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(path)
		if err2 := os.Rename(tmpPath, path); err2 != nil {
			_ = os.Remove(tmpPath)
			return fmt.Errorf("replace state file: %w", err2)
		}
	}
	return nil
}
