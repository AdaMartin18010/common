package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

// encryptAES 加密明文
func encryptAES(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext, nil
}

// decryptAES 解密密文
func decryptAES(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext, nil
}

func main() {
	key := []byte("1234567890abcdef") // 16字节密钥
	plaintext := []byte("Hello, Cybersecurity!")

	ciphertext, err := encryptAES(plaintext, key)
	if err != nil {
		fmt.Println("加密失败:", err)
		return
	}
	fmt.Println("密文:", hex.EncodeToString(ciphertext))

	decrypted, err := decryptAES(ciphertext, key)
	if err != nil {
		fmt.Println("解密失败:", err)
		return
	}
	fmt.Println("解密后:", string(decrypted))
}
