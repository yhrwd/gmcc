package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/crypto/scrypt"

	"gmcc/internal/logx"
	"gmcc/internal/session"
	"gmcc/internal/webtypes"
)

var (
	ErrInvalidPassword  = errors.New("invalid password")
	ErrTokenNotFound    = errors.New("token not found")
	ErrDecryptFailed    = errors.New("decryption failed")
	ErrInvalidTokenData = errors.New("invalid token data")
)

// Vault Token加密存储管理器
type Vault struct {
	config    webtypes.TokenVaultConfig
	keyGetter KeyGetter
	mu        sync.RWMutex
}

// KeyGetter 密钥获取函数类型
type KeyGetter func(password string) ([]byte, error)

// NewVault 创建Token Vault
func NewVault(config webtypes.TokenVaultConfig, keyGetter KeyGetter) *Vault {
	// 确保存储目录存在
	if err := os.MkdirAll(config.StoragePath, 0700); err != nil {
		logx.Warnf("创建Token存储目录失败: %v", err)
	}

	return &Vault{
		config:    config,
		keyGetter: keyGetter,
	}
}

// Store 加密并存储Token
func (v *Vault) Store(playerID string, token *session.TokenCache, password string) error {
	// 获取派生密钥
	key, err := v.keyGetter(password)
	if err != nil {
		return fmt.Errorf("key derivation failed: %w", err)
	}

	// 序列化Token数据
	tokenData, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	// 生成随机盐值
	salt := make([]byte, v.config.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	// 使用scrypt派生加密密钥
	// key已经是从机器指纹+密码派生的，这里我们再用scrypt处理一下
	// 实际上我们可以直接用key作为AES密钥，但为了安全起见，我们再用scrypt派生一次
	derivedKey, err := scrypt.Key(key, salt, v.config.ScryptN, v.config.ScryptR, v.config.ScryptP, 32)
	if err != nil {
		return fmt.Errorf("scrypt key derivation failed: %w", err)
	}

	// 生成随机nonce
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("failed to generate nonce: %w", err)
	}

	// AES-256-GCM加密
	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	ciphertext := aead.Seal(nil, nonce, tokenData, nil)

	// 构建加密Token结构
	encrypted := webtypes.EncryptedToken{
		Version:    1,
		Algorithm:  "aes-256-gcm",
		KDF:        "scrypt",
		ScryptN:    v.config.ScryptN,
		ScryptR:    v.config.ScryptR,
		ScryptP:    v.config.ScryptP,
		Salt:       salt,
		Nonce:      nonce,
		Ciphertext: ciphertext,
		PlayerID:   playerID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 写入文件
	filename := v.getFilePath(playerID)
	data, err := json.MarshalIndent(encrypted, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal encrypted token: %w", err)
	}

	// 使用互斥锁保护写操作
	v.mu.Lock()
	defer v.mu.Unlock()

	if err := os.WriteFile(filename, data, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	logx.Debugf("Token已加密存储: player_id=%s, file=%s", playerID, filename)
	return nil
}

// Retrieve 解密并读取Token
func (v *Vault) Retrieve(playerID string, password string) (*session.TokenCache, error) {
	filename := v.getFilePath(playerID)

	v.mu.RLock()
	defer v.mu.RUnlock()

	// 读取文件
	data, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrTokenNotFound
		}
		return nil, fmt.Errorf("failed to read token file: %w", err)
	}

	// 解析加密Token
	var encrypted webtypes.EncryptedToken
	if err := json.Unmarshal(data, &encrypted); err != nil {
		return nil, fmt.Errorf("failed to unmarshal encrypted token: %w", err)
	}

	// 获取派生密钥
	key, err := v.keyGetter(password)
	if err != nil {
		return nil, fmt.Errorf("key derivation failed: %w", err)
	}

	// 使用scrypt派生加密密钥
	derivedKey, err := scrypt.Key(key, encrypted.Salt, encrypted.ScryptN, encrypted.ScryptR, encrypted.ScryptP, 32)
	if err != nil {
		return nil, fmt.Errorf("scrypt key derivation failed: %w", err)
	}

	// AES-256-GCM解密
	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	plaintext, err := aead.Open(nil, encrypted.Nonce, encrypted.Ciphertext, nil)
	if err != nil {
		// 解密失败，可能是密码错误
		return nil, ErrDecryptFailed
	}

	// 解析Token数据
	var token session.TokenCache
	if err := json.Unmarshal(plaintext, &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token data: %w", err)
	}

	logx.Debugf("Token已解密读取: player_id=%s", playerID)
	return &token, nil
}

// Delete 删除存储的Token
func (v *Vault) Delete(playerID string) error {
	filename := v.getFilePath(playerID)

	v.mu.Lock()
	defer v.mu.Unlock()

	if err := os.Remove(filename); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrTokenNotFound
		}
		return fmt.Errorf("failed to delete token file: %w", err)
	}

	logx.Debugf("Token已删除: player_id=%s", playerID)
	return nil
}

// Exists 检查Token是否存在
func (v *Vault) Exists(playerID string) bool {
	filename := v.getFilePath(playerID)
	_, err := os.Stat(filename)
	return err == nil
}

// List 列出所有存储的Token
func (v *Vault) List() ([]string, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	entries, err := os.ReadDir(v.config.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read token directory: %w", err)
	}

	var playerIDs []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if filepath.Ext(name) == ".enc" {
			playerID := name[:len(name)-4]
			playerIDs = append(playerIDs, playerID)
		}
	}

	return playerIDs, nil
}

// getFilePath 获取Token文件路径
func (v *Vault) getFilePath(playerID string) string {
	// 对playerID进行base64编码以安全地作为文件名
	encodedID := base64.URLEncoding.EncodeToString([]byte(playerID))
	return filepath.Join(v.config.StoragePath, encodedID+".enc")
}

// ClearMemory 清零敏感数据（工具函数）
func ClearMemory(data []byte) {
	for i := range data {
		data[i] = 0
	}
}
