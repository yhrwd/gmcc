package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
)

// AESGCMEncrypt 使用AES-GCM算法加密数据
// 使用示例:
//
//	key := make([]byte, 32) // 256-bit key
//	rand.Read(key)
//	ciphertext, err := AESGCMEncrypt(key, []byte("Hello, World!"))
func AESGCMEncrypt(key, plaintext []byte) ([]byte, error) {
	// 创建AES密码块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// 创建随机nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// AESGCMDecrypt 使用AES-GCM算法解密数据
// 使用示例:
//
//	plaintext, err := AESGCMDecrypt(key, ciphertext)
func AESGCMDecrypt(key, ciphertext []byte) ([]byte, error) {
	// 创建AES密码块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// 检查密文长度
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// 分离nonce和密文
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 解密数据
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// GenerateRSAKeyPair 生成RSA密钥对
// 使用示例:
//
//	privateKey, publicKey, err := GenerateRSAKeyPair(2048)
func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate RSA key pair: %w", err)
	}

	// 验证密钥
	if err := privateKey.Validate(); err != nil {
		return nil, nil, fmt.Errorf("failed to validate private key: %w", err)
	}

	return privateKey, &privateKey.PublicKey, nil
}

// RSAEncrypt 使用RSA公钥加密数据
// 使用示例:
//
//	ciphertext, err := RSAEncrypt(publicKey, []byte("Hello, World!"))
func RSAEncrypt(publicKey *rsa.PublicKey, plaintext []byte) ([]byte, error) {
	// 使用OAEP填充方案加密
	label := []byte("")
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, plaintext, label)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt with RSA: %w", err)
	}

	return ciphertext, nil
}

// RSADecrypt 使用RSA私钥解密数据
// 使用示例:
//
//	plaintext, err := RSADecrypt(privateKey, ciphertext)
func RSADecrypt(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	// 使用OAEP填充方案解密
	label := []byte("")
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, label)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt with RSA: %w", err)
	}

	return plaintext, nil
}

// EncodePrivateKeyToPEM 将RSA私钥编码为PEM格式
// 使用示例:
//
//	pemBytes, err := EncodePrivateKeyToPEM(privateKey)
func EncodePrivateKeyToPEM(privateKey *rsa.PrivateKey) ([]byte, error) {
	// 创建PEM块
	privBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// 返回PEM编码的字节
	return pem.EncodeToMemory(privBlock), nil
}

// EncodePublicKeyToPEM 将RSA公钥编码为PEM格式
// 使用示例:
//
//	pemBytes, err := EncodePublicKeyToPEM(publicKey)
func EncodePublicKeyToPEM(publicKey *rsa.PublicKey) ([]byte, error) {
	// 创建PEM块
	pubBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}

	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}

	// 返回PEM编码的字节
	return pem.EncodeToMemory(pubBlock), nil
}

// DecodePrivateKeyFromPEM 从PEM格式解码RSA私钥
// 使用示例:
//
//	privateKey, err := DecodePrivateKeyFromPEM(pemBytes)
func DecodePrivateKeyFromPEM(pemBytes []byte) (*rsa.PrivateKey, error) {
	// 解码PEM块
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// 解析私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}

// DecodePublicKeyFromPEM 从PEM格式解码RSA公钥
// 使用示例:
//
//	publicKey, err := DecodePublicKeyFromPEM(pemBytes)
func DecodePublicKeyFromPEM(pemBytes []byte) (*rsa.PublicKey, error) {
	// 解码PEM块
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// 解析公钥
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	// 类型断言为RSA公钥
	publicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return publicKey, nil
}
