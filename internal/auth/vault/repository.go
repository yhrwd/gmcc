package vault

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gmcc/internal/auth/session"
)

type Config struct {
	Dir       string
	MasterKey []byte
	ScryptN   int
	ScryptR   int
	ScryptP   int
	SaltLen   int
}

type Repository struct {
	dir   string
	cfg   Config
	locks sync.Map
}

func NewRepository(cfg Config) (*Repository, error) {
	cfg.Dir = strings.TrimSpace(cfg.Dir)
	if cfg.Dir == "" {
		return nil, fmt.Errorf("vault dir is empty")
	}
	if cfg.SaltLen <= 0 {
		return nil, fmt.Errorf("vault salt length is invalid")
	}
	if len(cfg.MasterKey) == 0 {
		return nil, fmt.Errorf("vault master key is empty")
	}
	if _, err := deriveKey(cfg, make([]byte, cfg.SaltLen)); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(cfg.Dir, 0o700); err != nil {
		return nil, fmt.Errorf("create vault dir: %w", err)
	}

	return &Repository{dir: cfg.Dir, cfg: cfg}, nil
}

func (r *Repository) PutAccount(record *session.AccountAuthRecord) error {
	if record == nil {
		return fmt.Errorf("account auth record is nil")
	}
	payload, err := encodeRecord(r.cfg, record)
	if err != nil {
		return err
	}
	return writeAtomic(r.pathForAccount(record.AccountID), payload, 0o600, r.lockForAccount(record.AccountID))
}

func (r *Repository) GetAccount(accountID string) (*session.AccountAuthRecord, error) {
	path := r.pathForAccount(accountID)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, session.ErrAccountNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("read account auth vault: %w", err)
	}

	record, err := decodeRecord(r.cfg, accountID, data)
	if err != nil {
		return nil, fmt.Errorf("decode account auth vault: %w", err)
	}
	return record, nil
}

func (r *Repository) DeleteAccount(accountID string) error {
	lock := r.lockForAccount(accountID)
	lock.Lock()
	defer lock.Unlock()

	if err := os.Remove(r.pathForAccount(accountID)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("delete account auth vault: %w", err)
	}
	return nil
}

func (r *Repository) pathForAccount(accountID string) string {
	return filepath.Join(r.dir, fileNameForAccount(accountID))
}

func (r *Repository) lockForAccount(accountID string) *sync.Mutex {
	key := fileNameForAccount(accountID)
	lock, _ := r.locks.LoadOrStore(key, &sync.Mutex{})
	return lock.(*sync.Mutex)
}

func fileNameForAccount(accountID string) string {
	trimmed := strings.TrimSpace(accountID)
	if trimmed == "" {
		trimmed = "default"
	}

	hash := sha256.Sum256([]byte(trimmed))
	base := sanitizeFileBase(trimmed)
	return base + "--" + hex.EncodeToString(hash[:8]) + ".vault"
}

func sanitizeFileBase(accountID string) string {
	var b strings.Builder
	b.Grow(len(accountID))
	for _, r := range accountID {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '_' || r == '-':
			b.WriteRune(r)
		default:
			b.WriteByte('-')
		}
	}

	clean := strings.Trim(b.String(), "-_")
	if clean == "" {
		return "account"
	}
	return clean
}

func writeAtomic(path string, data []byte, perm os.FileMode, lock *sync.Mutex) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create vault parent dir: %w", err)
	}

	lock.Lock()
	defer lock.Unlock()

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, perm); err != nil {
		return fmt.Errorf("write vault temp file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(path)
		if err2 := os.Rename(tmpPath, path); err2 != nil {
			_ = os.Remove(tmpPath)
			return fmt.Errorf("replace vault file: %w", err2)
		}
	}
	return nil
}
