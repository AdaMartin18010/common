package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// 计算字符串的SHA256哈希
func hashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// 计算文件的SHA256哈希
func hashFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func main() {
	str := "Hello, SHA256!"
	hash := hashString(str)
	fmt.Println("字符串哈希:", hash)

	// 文件哈希示例（假设有test.txt文件）
	fileHash, err := hashFile("test.txt")
	if err != nil {
		fmt.Println("文件哈希失败:", err)
	} else {
		fmt.Println("文件哈希:", fileHash)
	}
}
