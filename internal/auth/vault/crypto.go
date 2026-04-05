package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/scrypt"

	"gmcc/internal/auth/session"
)

const (
	vaultVersion   = 1
	vaultAlgorithm = "aes-256-gcm"
	vaultKDF       = "scrypt"
	derivedKeyLen  = 32
)

var errInvalidVaultPayload = errors.New("invalid vault payload")

type encryptedRecord struct {
	Version    int       `json:"version"`
	Algorithm  string    `json:"algorithm"`
	KDF        string    `json:"kdf"`
	ScryptN    int       `json:"scrypt_n"`
	ScryptR    int       `json:"scrypt_r"`
	ScryptP    int       `json:"scrypt_p"`
	Salt       []byte    `json:"salt"`
	Nonce      []byte    `json:"nonce"`
	Ciphertext []byte    `json:"ciphertext"`
	AccountID  string    `json:"account_id"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func encodeRecord(cfg Config, record *session.AccountAuthRecord) ([]byte, error) {
	if record == nil {
		return nil, fmt.Errorf("account auth record is nil")
	}

	accountID := strings.TrimSpace(record.AccountID)
	if accountID == "" {
		return nil, fmt.Errorf("account auth record account id is empty")
	}

	plain, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("encode account auth record: %w", err)
	}

	salt := make([]byte, cfg.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("generate vault salt: %w", err)
	}

	key, err := deriveKey(cfg, salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create vault cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create vault gcm: %w", err)
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("generate vault nonce: %w", err)
	}

	encrypted := encryptedRecord{
		Version:    vaultVersion,
		Algorithm:  vaultAlgorithm,
		KDF:        vaultKDF,
		ScryptN:    cfg.ScryptN,
		ScryptR:    cfg.ScryptR,
		ScryptP:    cfg.ScryptP,
		Salt:       salt,
		Nonce:      nonce,
		Ciphertext: aead.Seal(nil, nonce, plain, []byte(accountID)),
		AccountID:  accountID,
		UpdatedAt:  record.UpdatedAt,
	}

	payload, err := json.MarshalIndent(encrypted, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode encrypted account auth record: %w", err)
	}

	return payload, nil
}

func decodeRecord(cfg Config, accountID string, payload []byte) (*session.AccountAuthRecord, error) {
	var encrypted encryptedRecord
	if err := json.Unmarshal(payload, &encrypted); err != nil {
		return nil, fmt.Errorf("decode encrypted account auth record: %w", err)
	}

	trimmedAccountID := strings.TrimSpace(accountID)
	storedAccountID := strings.TrimSpace(encrypted.AccountID)
	if storedAccountID == "" || storedAccountID != trimmedAccountID {
		return nil, errInvalidVaultPayload
	}
	if encrypted.Version != vaultVersion || encrypted.Algorithm != vaultAlgorithm || encrypted.KDF != vaultKDF {
		return nil, errInvalidVaultPayload
	}
	if len(encrypted.Salt) == 0 || len(encrypted.Nonce) == 0 || len(encrypted.Ciphertext) == 0 {
		return nil, errInvalidVaultPayload
	}

	key, err := deriveKey(Config{
		MasterKey: cfg.MasterKey,
		ScryptN:   encrypted.ScryptN,
		ScryptR:   encrypted.ScryptR,
		ScryptP:   encrypted.ScryptP,
	}, encrypted.Salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create vault cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create vault gcm: %w", err)
	}

	plain, err := aead.Open(nil, encrypted.Nonce, encrypted.Ciphertext, []byte(trimmedAccountID))
	if err != nil {
		return nil, fmt.Errorf("decrypt account auth record: %w", err)
	}

	record := session.NewAccountAuthRecord(trimmedAccountID)
	if err := json.Unmarshal(plain, record); err != nil {
		return nil, fmt.Errorf("decode account auth record: %w", err)
	}
	record.AccountID = trimmedAccountID

	return record, nil
}

func deriveKey(cfg Config, salt []byte) ([]byte, error) {
	if len(cfg.MasterKey) == 0 {
		return nil, fmt.Errorf("vault master key is empty")
	}
	if cfg.ScryptN <= 1 || cfg.ScryptR <= 0 || cfg.ScryptP <= 0 {
		return nil, fmt.Errorf("vault scrypt parameters are invalid")
	}

	key, err := scrypt.Key(cfg.MasterKey, salt, cfg.ScryptN, cfg.ScryptR, cfg.ScryptP, derivedKeyLen)
	if err != nil {
		return nil, fmt.Errorf("derive vault key: %w", err)
	}
	return key, nil
}
