package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"math/big"
)

// ==================== 工具：从任意字符串生成 32 字节 AES 密钥 ====================
func deriveKey(secret string) []byte {
	sum := sha256.Sum256([]byte(secret))
	return sum[:]
}

// RandomScreatKey 生成一个用户可读、可自己输入的随机密钥
func RandomScreatKey() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	n := 24 // 密钥长度 24 字节

	buf := make([]byte, n)
	for i := range buf {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			// 出错就退回一个固定值，避免 panic
			log.Println("rand.Int err:", err)
			return "DefaultRandomSecretKey123"
		}
		buf[i] = letters[num.Int64()]
	}
	return string(buf)
}

// ==================== AES-GCM 加密 ====================
// Encrypt 用 AES-GCM 加密文本，secretKey 是用户输入的任意字符串
// 返回 Base64 编码后的 nonce + ciphertext
func Encrypt(plainText, secretKey string) string {
	// 1. 从用户密钥派生出固定长度 AES key
	key := deriveKey(secretKey)

	// 2. 创建 AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("Encrypt: NewCipher err:", err)
		return ""
	}

	// 3. 使用 GCM 模式
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("Encrypt: NewGCM err:", err)
		return ""
	}

	// 4. 生成随机 nonce
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println("Encrypt: ReadFull nonce err:", err)
		return ""
	}

	// 5. 加密
	ciphertext := aesgcm.Seal(nil, nonce, []byte(plainText), nil)

	// 6. 拼接 nonce + ciphertext，Base64 编码返回
	result := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(result)
}

// ==================== AES-GCM 解密 ====================
// Decrypt 用 AES-GCM 解密文本，secretKey 是用户输入的任意字符串
// cipherBase64 是 Encrypt 返回的 Base64 字符串
func Decrypt(cipherBase64, secretKey string) string {
	// 1. 从用户密钥派生出 AES key
	key := deriveKey(secretKey)

	// 2. Base64 解码，得到 nonce + ciphertext
	data, err := base64.StdEncoding.DecodeString(cipherBase64)
	if err != nil {
		log.Println("Decrypt: DecodeString err:", err)
		return ""
	}

	// 3. 创建 AES cipher 和 GCM
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("Decrypt: NewCipher err:", err)
		return ""
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("Decrypt: NewGCM err:", err)
		return ""
	}

	nonceSize := aesgcm.NonceSize()
	if len(data) < nonceSize {
		log.Println("Decrypt: data too short")
		return ""
	}

	// 4. 拆出 nonce 和真正的密文
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// 5. 解密
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Println("Decrypt: Open err:", err)
		return ""
	}

	return string(plaintext)
}
