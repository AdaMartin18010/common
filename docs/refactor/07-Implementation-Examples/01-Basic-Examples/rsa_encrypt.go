package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// 生成RSA密钥对
func generateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return privKey, &privKey.PublicKey, nil
}

// 使用公钥加密
func encryptRSA(plaintext []byte, pub *rsa.PublicKey) ([]byte, error) {
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, plaintext, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// 使用私钥解密
func decryptRSA(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func main() {
	// 生成密钥对
	priv, pub, err := generateKeyPair(2048)
	if err != nil {
		fmt.Println("密钥生成失败:", err)
		return
	}

	message := []byte("Hello, RSA Encryption!")

	// 加密
	ciphertext, err := encryptRSA(message, pub)
	if err != nil {
		fmt.Println("加密失败:", err)
		return
	}
	fmt.Println("密文(Base64):", base64.StdEncoding.EncodeToString(ciphertext))

	// 解密
	decrypted, err := decryptRSA(ciphertext, priv)
	if err != nil {
		fmt.Println("解密失败:", err)
		return
	}
	fmt.Println("解密后:", string(decrypted))
}
