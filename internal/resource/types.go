package resource

import (
	"errors"
	"strings"

	"gmcc/internal/auth/session"
	"gmcc/internal/state"
)

type ReasonCode string

const (
	ReasonAccountNotFound    ReasonCode = "account_not_found"
	ReasonAccountDisabled    ReasonCode = "account_disabled"
	ReasonAccountNotLoggedIn ReasonCode = "account_not_logged_in"
	ReasonAccountAuthInvalid ReasonCode = "account_auth_invalid"
)

var (
	ErrAccountIDRequired     = errors.New("account id is required")
	ErrInstanceIDRequired    = errors.New("instance id is required")
	ErrServerAddressRequired = errors.New("server address is required")
	ErrInstanceAlreadyExists = errors.New("instance already exists")
	ErrInstanceNotFound      = errors.New("instance not found")
	ErrAccountAlreadyExists  = errors.New("account already exists")
	ErrAccountNotFound       = errors.New("account not found")
	ErrAccountInUse          = errors.New("account in use")
	ErrAccountDisabled       = errors.New("account disabled")
	ErrAccountNotLoggedIn    = errors.New("account not logged in")
	ErrAccountAuthInvalid    = errors.New("account auth invalid")
)

type CreateAccountInput struct {
	AccountID string
	Enabled   bool
	Label     string
	Note      string
}

type CreateInstanceInput struct {
	InstanceID    string
	AccountID     string
	ServerAddress string
	Enabled       bool
}

type SkippedInstance struct {
	InstanceID string
	AccountID  string
	Reason     ReasonCode
}

type RestoreResourcesResult struct {
	RestoredInstances []state.InstanceMeta
	SkippedInstances  []SkippedInstance
	RestoredCount     int
	SkippedCount      int
}

type AccountRecord struct {
	Meta state.AccountMeta
}

type AccountMetadataRepository interface {
	LoadAll() ([]state.AccountMeta, error)
	SaveAll(accounts []state.AccountMeta) error
}

type InstanceMetadataRepository interface {
	LoadAll() ([]state.InstanceMeta, error)
	SaveAll(instances []state.InstanceMeta) error
}

type AuthStatusReader interface {
	GetAccountAuthStatus(accountID string) (session.AccountAuthStatus, error)
	Clear(accountID string) error
}

func normalizeCreateInstanceInput(in CreateInstanceInput) CreateInstanceInput {
	in.InstanceID = strings.TrimSpace(in.InstanceID)
	in.AccountID = strings.TrimSpace(in.AccountID)
	in.ServerAddress = strings.TrimSpace(in.ServerAddress)
	return in
}

func reasonForAuthStatus(status session.AccountAuthStatus) ReasonCode {
	switch status {
	case session.AccountAuthStatusAuthInvalid:
		return ReasonAccountAuthInvalid
	case session.AccountAuthStatusLoggedIn:
		return ""
	default:
		return ReasonAccountNotLoggedIn
	}
}

func errorForReason(reason ReasonCode) error {
	switch reason {
	case ReasonAccountNotFound:
		return ErrAccountNotFound
	case ReasonAccountDisabled:
		return ErrAccountDisabled
	case ReasonAccountAuthInvalid:
		return ErrAccountAuthInvalid
	case ReasonAccountNotLoggedIn:
		return ErrAccountNotLoggedIn
	default:
		return nil
	}
}
