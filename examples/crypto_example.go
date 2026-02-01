// Package main demonstrates the usage of the crypto utility functions.
package main

import (
	"crypto/rand"
	"fmt"
	"log"

	"gmcc/pkg/crypto"
)

func main() {
	fmt.Println("Crypto Utility Functions 示例")

	// 示例1: AES-GCM 加密和解密
	fmt.Println("\n=== 示例1: AES-GCM 加密和解密 ===")
	// 创建一个32字节的密钥（256位）
	key := make([]byte, 32)
	// 用随机数据填充密钥
	if _, err := rand.Read(key); err != nil {
		log.Fatal("生成随机密钥失败:", err)
	}

	// 要加密的明文
	plaintext := []byte("Hello, World! This is a secret message.")

	// 加密
	ciphertext, err := crypto.AESGCMEncrypt(key, plaintext)
	if err != nil {
		log.Fatal("AES加密失败:", err)
	}
	fmt.Printf("AES加密成功，密文长度: %d 字节\n", len(ciphertext))

	// 解密
	decrypted, err := crypto.AESGCMDecrypt(key, ciphertext)
	if err != nil {
		log.Fatal("AES解密失败:", err)
	}
	fmt.Printf("AES解密成功: %s\n", string(decrypted))

	// 示例2: RSA 密钥生成和加解密
	fmt.Println("\n=== 示例2: RSA 密钥生成和加解密 ===")
	// 生成RSA密钥对（使用较小的密钥长度以提高速度）
	privateKey, publicKey, err := crypto.GenerateRSAKeyPair(2048)
	if err != nil {
		log.Fatal("RSA密钥对生成失败:", err)
	}
	fmt.Println("RSA密钥对生成成功")

	// 使用RSA公钥加密
	rsaPlaintext := []byte("This is a secret message for RSA encryption.")
	rsaCiphertext, err := crypto.RSAEncrypt(publicKey, rsaPlaintext)
	if err != nil {
		log.Fatal("RSA加密失败:", err)
	}
	fmt.Printf("RSA加密成功，密文长度: %d 字节\n", len(rsaCiphertext))

	// 使用RSA私钥解密
	rsaDecrypted, err := crypto.RSADecrypt(privateKey, rsaCiphertext)
	if err != nil {
		log.Fatal("RSA解密失败:", err)
	}
	fmt.Printf("RSA解密成功: %s\n", string(rsaDecrypted))

	// 示例3: RSA密钥的PEM编码和解码
	fmt.Println("\n=== 示例3: RSA密钥的PEM编码和解码 ===")
	// 将私钥编码为PEM格式
	privateKeyPEM, err := crypto.EncodePrivateKeyToPEM(privateKey)
	if err != nil {
		log.Fatal("私钥PEM编码失败:", err)
	}
	fmt.Printf("私钥PEM编码成功，长度: %d 字节\n", len(privateKeyPEM))

	// 将公钥编码为PEM格式
	publicKeyPEM, err := crypto.EncodePublicKeyToPEM(publicKey)
	if err != nil {
		log.Fatal("公钥PEM编码失败:", err)
	}
	fmt.Printf("公钥PEM编码成功，长度: %d 字节\n", len(publicKeyPEM))

	// 从PEM格式解码私钥
	decodedPrivateKey, err := crypto.DecodePrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		log.Fatal("私钥PEM解码失败:", err)
	}
	fmt.Println("私钥PEM解码成功")

	// 从PEM格式解码公钥
	decodedPublicKey, err := crypto.DecodePublicKeyFromPEM(publicKeyPEM)
	if err != nil {
		log.Fatal("公钥PEM解码失败:", err)
	}
	fmt.Println("公钥PEM解码成功")

	// 验证解码后的密钥是否与原始密钥相同
	if decodedPrivateKey.N.Cmp(privateKey.N) == 0 && decodedPrivateKey.E == privateKey.E {
		fmt.Println("解码的私钥与原始私钥匹配")
	}

	if decodedPublicKey.N.Cmp(publicKey.N) == 0 && decodedPublicKey.E == publicKey.E {
		fmt.Println("解码的公钥与原始公钥匹配")
	}

	fmt.Println("\n所有示例完成!")
}
