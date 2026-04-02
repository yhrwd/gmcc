package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"gmcc/internal/logx"
	"gmcc/internal/webtypes"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrTokenExpired    = errors.New("token expired")
	ErrTokenInvalid    = errors.New("token invalid")
	ErrNoPasswords     = errors.New("no passwords configured")
)

// Manager 认证管理器
type Manager struct {
	config    webtypes.AuthConfig
	secretKey []byte
}

// NewManager 创建认证管理器
func NewManager(config webtypes.AuthConfig) (*Manager, error) {
	if len(config.Passwords) == 0 {
		return nil, ErrNoPasswords
	}

	// 生成JWT密钥（可以基于配置或其他方式）
	secretKey := []byte(fmt.Sprintf("gmcc-jwt-%d", time.Now().UnixNano()))

	return &Manager{
		config:    config,
		secretKey: secretKey,
	}, nil
}

// VerifyPassword 验证密码
// 返回匹配的密码ID和错误
func (m *Manager) VerifyPassword(password string) (string, error) {
	for _, entry := range m.config.Passwords {
		if !entry.Enabled {
			continue
		}

		err := bcrypt.CompareHashAndPassword([]byte(entry.Hash), []byte(password))
		if err == nil {
			logx.Debugf("密码验证成功: password_id=%s", entry.ID)
			return entry.ID, nil
		}
	}

	return "", ErrInvalidPassword
}

// HashPassword 生成密码的bcrypt hash
func (m *Manager) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// GenerateToken 生成JWT Token
func (m *Manager) GenerateToken(passwordID string) (string, time.Time, error) {
	expiresAt := time.Now().Add(m.config.TokenExpiry)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"password_id": passwordID,
		"exp":         expiresAt.Unix(),
		"iat":         time.Now().Unix(),
	})

	tokenString, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// ValidateToken 验证JWT Token
func (m *Manager) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", ErrTokenExpired
		}
		return "", ErrTokenInvalid
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		passwordID, ok := claims["password_id"].(string)
		if !ok {
			return "", ErrTokenInvalid
		}
		return passwordID, nil
	}

	return "", ErrTokenInvalid
}

// GetPasswordEntry 获取密码条目
func (m *Manager) GetPasswordEntry(id string) (webtypes.PasswordEntry, bool) {
	for _, entry := range m.config.Passwords {
		if entry.ID == id {
			return entry, true
		}
	}
	return webtypes.PasswordEntry{}, false
}

// HashPassword 生成密码哈希（工具函数）
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}
