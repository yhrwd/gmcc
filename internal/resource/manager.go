package resource

import (
	"fmt"
	"strings"

	"gmcc/internal/auth/session"
	"gmcc/internal/state"
)

type Manager struct {
	accounts  AccountMetadataRepository
	instances InstanceMetadataRepository
	auth      AuthStatusReader
}

func NewManager(accounts AccountMetadataRepository, instances InstanceMetadataRepository, auth AuthStatusReader) *Manager {
	return &Manager{accounts: accounts, instances: instances, auth: auth}
}

func (m *Manager) ListAccounts() ([]AccountRecord, error) {
	accounts, err := m.accounts.LoadAll()
	if err != nil {
		return nil, fmt.Errorf("load accounts: %w", err)
	}

	records := make([]AccountRecord, 0, len(accounts))
	for _, account := range accounts {
		records = append(records, AccountRecord{Meta: account})
	}
	return records, nil
}

func (m *Manager) GetAccount(accountID string) (AccountRecord, error) {
	accountID = normalizeID(accountID)
	if accountID == "" {
		return AccountRecord{}, ErrAccountIDRequired
	}

	accounts, err := m.accounts.LoadAll()
	if err != nil {
		return AccountRecord{}, fmt.Errorf("load accounts: %w", err)
	}

	account, ok := findAccount(accounts, accountID)
	if !ok {
		return AccountRecord{}, ErrAccountNotFound
	}
	return AccountRecord{Meta: account}, nil
}

func (m *Manager) CreateAccount(in CreateAccountInput) (state.AccountMeta, error) {
	in.AccountID = normalizeID(in.AccountID)
	if in.AccountID == "" {
		return state.AccountMeta{}, ErrAccountIDRequired
	}

	accounts, err := m.accounts.LoadAll()
	if err != nil {
		return state.AccountMeta{}, fmt.Errorf("load accounts: %w", err)
	}
	for _, existing := range accounts {
		if existing.AccountID == in.AccountID {
			return state.AccountMeta{}, ErrAccountAlreadyExists
		}
	}

	saved := state.AccountMeta{
		AccountID: in.AccountID,
		Enabled:   in.Enabled,
		Label:     strings.TrimSpace(in.Label),
		Note:      strings.TrimSpace(in.Note),
	}
	accounts = append(accounts, saved)
	if err := m.accounts.SaveAll(accounts); err != nil {
		return state.AccountMeta{}, fmt.Errorf("save accounts: %w", err)
	}

	return saved, nil
}

func (m *Manager) CreateInstance(in CreateInstanceInput) (state.InstanceMeta, error) {
	in = normalizeCreateInstanceInput(in)
	if in.InstanceID == "" {
		return state.InstanceMeta{}, ErrInstanceIDRequired
	}
	if in.AccountID == "" {
		return state.InstanceMeta{}, ErrAccountIDRequired
	}
	if in.ServerAddress == "" {
		return state.InstanceMeta{}, ErrServerAddressRequired
	}

	accounts, err := m.accounts.LoadAll()
	if err != nil {
		return state.InstanceMeta{}, fmt.Errorf("load accounts: %w", err)
	}
	instances, err := m.instances.LoadAll()
	if err != nil {
		return state.InstanceMeta{}, fmt.Errorf("load instances: %w", err)
	}

	for _, existing := range instances {
		if existing.InstanceID == in.InstanceID {
			return state.InstanceMeta{}, ErrInstanceAlreadyExists
		}
	}

	if err := m.validateAccount(accounts, in.AccountID); err != nil {
		return state.InstanceMeta{}, err
	}

	saved := state.InstanceMeta{
		InstanceID:    in.InstanceID,
		AccountID:     in.AccountID,
		ServerAddress: in.ServerAddress,
		Enabled:       in.Enabled,
	}
	instances = append(instances, saved)
	if err := m.instances.SaveAll(instances); err != nil {
		return state.InstanceMeta{}, fmt.Errorf("save instances: %w", err)
	}

	return saved, nil
}

func (m *Manager) RestoreResources() (RestoreResourcesResult, error) {
	accounts, err := m.accounts.LoadAll()
	if err != nil {
		return RestoreResourcesResult{}, fmt.Errorf("load accounts: %w", err)
	}
	instances, err := m.instances.LoadAll()
	if err != nil {
		return RestoreResourcesResult{}, fmt.Errorf("load instances: %w", err)
	}

	accountByID := make(map[string]state.AccountMeta, len(accounts))
	for _, account := range accounts {
		accountByID[account.AccountID] = account
	}

	result := RestoreResourcesResult{
		RestoredInstances: make([]state.InstanceMeta, 0, len(instances)),
		SkippedInstances:  make([]SkippedInstance, 0),
	}

	for _, instance := range instances {
		reason, err := m.restoreReasonForInstance(accountByID, instance)
		if err != nil {
			return RestoreResourcesResult{}, err
		}
		if reason != "" {
			result.SkippedInstances = append(result.SkippedInstances, SkippedInstance{
				InstanceID: instance.InstanceID,
				AccountID:  instance.AccountID,
				Reason:     reason,
			})
			continue
		}

		result.RestoredInstances = append(result.RestoredInstances, instance)
	}

	result.RestoredCount = len(result.RestoredInstances)
	result.SkippedCount = len(result.SkippedInstances)
	return result, nil
}

func (m *Manager) DeleteAccount(accountID string) error {
	accountID = normalizeID(accountID)
	if accountID == "" {
		return ErrAccountIDRequired
	}

	accounts, err := m.accounts.LoadAll()
	if err != nil {
		return fmt.Errorf("load accounts: %w", err)
	}
	instances, err := m.instances.LoadAll()
	if err != nil {
		return fmt.Errorf("load instances: %w", err)
	}

	index := -1
	for i, account := range accounts {
		if account.AccountID == accountID {
			index = i
			break
		}
	}
	if index == -1 {
		return ErrAccountNotFound
	}

	for _, instance := range instances {
		if instance.AccountID == accountID {
			return ErrAccountInUse
		}
	}

	if err := m.auth.Clear(accountID); err != nil {
		return fmt.Errorf("clear account auth state: %w", err)
	}
	accounts = append(accounts[:index], accounts[index+1:]...)
	if err := m.accounts.SaveAll(accounts); err != nil {
		return fmt.Errorf("save accounts: %w", err)
	}
	return nil
}

func (m *Manager) DeleteInstance(instanceID string) error {
	instanceID = normalizeID(instanceID)
	if instanceID == "" {
		return ErrInstanceIDRequired
	}

	instances, err := m.instances.LoadAll()
	if err != nil {
		return fmt.Errorf("load instances: %w", err)
	}

	index := -1
	for i, instance := range instances {
		if instance.InstanceID == instanceID {
			index = i
			break
		}
	}
	if index == -1 {
		return ErrInstanceNotFound
	}

	instances = append(instances[:index], instances[index+1:]...)
	if err := m.instances.SaveAll(instances); err != nil {
		return fmt.Errorf("save instances: %w", err)
	}
	return nil
}

func (m *Manager) validateAccount(accounts []state.AccountMeta, accountID string) error {
	account, ok := findAccount(accounts, accountID)
	if !ok {
		return ErrAccountNotFound
	}
	if !account.Enabled {
		return ErrAccountDisabled
	}

	status, err := m.auth.GetAccountAuthStatus(accountID)
	if err != nil {
		if err == session.ErrAccountNotFound {
			return ErrAccountNotLoggedIn
		}
		return fmt.Errorf("get account auth status: %w", err)
	}

	if reason := reasonForAuthStatus(status); reason != "" {
		return errorForReason(reason)
	}
	return nil
}

func (m *Manager) restoreReasonForInstance(accounts map[string]state.AccountMeta, instance state.InstanceMeta) (ReasonCode, error) {
	account, ok := accounts[instance.AccountID]
	if !ok {
		return ReasonAccountNotFound, nil
	}
	if !account.Enabled {
		return ReasonAccountDisabled, nil
	}

	status, err := m.auth.GetAccountAuthStatus(instance.AccountID)
	if err != nil {
		if err == session.ErrAccountNotFound {
			return ReasonAccountNotLoggedIn, nil
		}
		return "", fmt.Errorf("get account auth status for instance %q: %w", instance.InstanceID, err)
	}

	return reasonForAuthStatus(status), nil
}

func findAccount(accounts []state.AccountMeta, accountID string) (state.AccountMeta, bool) {
	for _, account := range accounts {
		if account.AccountID == accountID {
			return account, true
		}
	}
	return state.AccountMeta{}, false
}

func normalizeID(value string) string {
	return strings.TrimSpace(value)
}
