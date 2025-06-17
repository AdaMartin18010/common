# 5. 密码学基础 (Cryptography)

## 5.1 密码学基础理论

### 5.1.1 密码学定义与分类

密码学是研究信息安全的科学，包括加密、解密、认证、完整性验证等技术。

**形式化定义**：
设 $M$ 为明文空间，$C$ 为密文空间，$K$ 为密钥空间。
加密函数 $E: M \times K \rightarrow C$ 定义为：
$$E(m,k) = c$$
解密函数 $D: C \times K \rightarrow M$ 定义为：
$$D(c,k) = m$$

**密码学分类**：
1. **对称密码学**：加密和解密使用相同密钥
2. **非对称密码学**：加密和解密使用不同密钥
3. **哈希函数**：单向函数，不可逆
4. **数字签名**：提供认证和不可否认性

### 5.1.2 密码学安全模型

**安全目标**：
1. **机密性**：确保信息不被未授权访问
2. **完整性**：确保信息不被篡改
3. **认证性**：确保信息发送者身份
4. **不可否认性**：确保发送者不能否认发送

**攻击模型**：
- **唯密文攻击**：攻击者只能获得密文
- **已知明文攻击**：攻击者知道部分明文和对应密文
- **选择明文攻击**：攻击者可以选择明文获得密文
- **选择密文攻击**：攻击者可以选择密文获得明文

## 5.2 对称密码学

### 5.2.1 对称加密算法

**常见算法**：
1. **AES (Advanced Encryption Standard)**：高级加密标准
2. **DES (Data Encryption Standard)**：数据加密标准
3. **3DES (Triple DES)**：三重DES
4. **ChaCha20**：流密码算法

**Go语言对称加密实现**：

```go
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "fmt"
    "io"
)

// SymmetricCrypto 对称加密
type SymmetricCrypto struct{}

// NewSymmetricCrypto 创建对称加密
func NewSymmetricCrypto() *SymmetricCrypto {
    return &SymmetricCrypto{}
}

// GenerateKey 生成密钥
func (s *SymmetricCrypto) GenerateKey(keySize int) ([]byte, error) {
    key := make([]byte, keySize)
    _, err := io.ReadFull(rand.Reader, key)
    if err != nil {
        return nil, fmt.Errorf("failed to generate key: %v", err)
    }
    return key, nil
}

// AESEncrypt AES加密
func (s *SymmetricCrypto) AESEncrypt(plaintext, key []byte) ([]byte, error) {
    // 创建AES密码块
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("failed to create cipher: %v", err)
    }

    // 生成随机IV
    iv := make([]byte, aes.BlockSize)
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, fmt.Errorf("failed to generate IV: %v", err)
    }

    // 创建GCM模式
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM: %v", err)
    }

    // 加密
    ciphertext := gcm.Seal(iv, iv, plaintext, nil)
    return ciphertext, nil
}

// AESDecrypt AES解密
func (s *SymmetricCrypto) AESDecrypt(ciphertext, key []byte) ([]byte, error) {
    // 创建AES密码块
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("failed to create cipher: %v", err)
    }

    // 创建GCM模式
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM: %v", err)
    }

    // 提取IV
    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    // 解密
    plaintext, err := gcm.Open(nil, iv, ciphertext, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt: %v", err)
    }

    return plaintext, nil
}

// ChaCha20Encrypt ChaCha20加密
func (s *SymmetricCrypto) ChaCha20Encrypt(plaintext, key, nonce []byte) ([]byte, error) {
    // 创建ChaCha20密码
    cipher, err := cipher.NewChaCha20(key, nonce)
    if err != nil {
        return nil, fmt.Errorf("failed to create ChaCha20: %v", err)
    }

    // 加密
    ciphertext := make([]byte, len(plaintext))
    cipher.XORKeyStream(ciphertext, plaintext)

    return ciphertext, nil
}

// ChaCha20Decrypt ChaCha20解密
func (s *SymmetricCrypto) ChaCha20Decrypt(ciphertext, key, nonce []byte) ([]byte, error) {
    // ChaCha20是对称的，解密就是加密
    return s.ChaCha20Encrypt(ciphertext, key, nonce)
}

// StreamCipher 流密码
type StreamCipher struct {
    key   []byte
    nonce []byte
}

// NewStreamCipher 创建流密码
func NewStreamCipher(key, nonce []byte) *StreamCipher {
    return &StreamCipher{
        key:   key,
        nonce: nonce,
    }
}

// Encrypt 流密码加密
func (s *StreamCipher) Encrypt(plaintext []byte) ([]byte, error) {
    cipher, err := cipher.NewChaCha20(s.key, s.nonce)
    if err != nil {
        return nil, fmt.Errorf("failed to create stream cipher: %v", err)
    }

    ciphertext := make([]byte, len(plaintext))
    cipher.XORKeyStream(ciphertext, plaintext)

    return ciphertext, nil
}

// Decrypt 流密码解密
func (s *StreamCipher) Decrypt(ciphertext []byte) ([]byte, error) {
    return s.Encrypt(ciphertext)
}
```

### 5.2.2 对称加密模式

**加密模式**：
1. **ECB (Electronic Codebook)**：电子密码本模式
2. **CBC (Cipher Block Chaining)**：密码块链接模式
3. **CFB (Cipher Feedback)**：密码反馈模式
4. **OFB (Output Feedback)**：输出反馈模式
5. **CTR (Counter)**：计数器模式
6. **GCM (Galois/Counter Mode)**：伽罗瓦/计数器模式

**Go语言加密模式实现**：

```go
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "fmt"
    "io"
)

// EncryptionModes 加密模式
type EncryptionModes struct{}

// NewEncryptionModes 创建加密模式
func NewEncryptionModes() *EncryptionModes {
    return &EncryptionModes{}
}

// CBCEncrypt CBC模式加密
func (e *EncryptionModes) CBCEncrypt(plaintext, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("failed to create cipher: %v", err)
    }

    // 填充
    plaintext = e.pad(plaintext)

    // 生成IV
    iv := make([]byte, aes.BlockSize)
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, fmt.Errorf("failed to generate IV: %v", err)
    }

    // 创建CBC模式
    mode := cipher.NewCBCEncrypter(block, iv)

    // 加密
    ciphertext := make([]byte, len(plaintext))
    mode.CryptBlocks(ciphertext, plaintext)

    // 组合IV和密文
    result := append(iv, ciphertext...)
    return result, nil
}

// CBCDecrypt CBC模式解密
func (e *EncryptionModes) CBCDecrypt(ciphertext, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("failed to create cipher: %v", err)
    }

    // 分离IV和密文
    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    // 创建CBC模式
    mode := cipher.NewCBCDecrypter(block, iv)

    // 解密
    plaintext := make([]byte, len(ciphertext))
    mode.CryptBlocks(plaintext, ciphertext)

    // 去除填充
    plaintext = e.unpad(plaintext)
    return plaintext, nil
}

// CTRMode CTR模式
func (e *EncryptionModes) CTRMode(plaintext, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("failed to create cipher: %v", err)
    }

    // 生成随机数
    nonce := make([]byte, aes.BlockSize)
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, fmt.Errorf("failed to generate nonce: %v", err)
    }

    // 创建CTR模式
    stream := cipher.NewCTR(block, nonce)

    // 加密
    ciphertext := make([]byte, len(plaintext))
    stream.XORKeyStream(ciphertext, plaintext)

    // 组合随机数和密文
    result := append(nonce, ciphertext...)
    return result, nil
}

// pad PKCS7填充
func (e *EncryptionModes) pad(data []byte) []byte {
    padding := aes.BlockSize - len(data)%aes.BlockSize
    padtext := make([]byte, padding)
    for i := range padtext {
        padtext[i] = byte(padding)
    }
    return append(data, padtext...)
}

// unpad PKCS7去填充
func (e *EncryptionModes) unpad(data []byte) []byte {
    length := len(data)
    if length == 0 {
        return data
    }
    padding := int(data[length-1])
    if padding > length {
        return data
    }
    return data[:length-padding]
}
```

## 5.3 非对称密码学

### 5.3.1 非对称加密算法

**常见算法**：
1. **RSA**：基于大整数分解问题
2. **ECC (Elliptic Curve Cryptography)**：基于椭圆曲线
3. **DSA (Digital Signature Algorithm)**：数字签名算法
4. **ElGamal**：基于离散对数问题

**Go语言非对称加密实现**：

```go
package main

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/pem"
    "fmt"
)

// AsymmetricCrypto 非对称加密
type AsymmetricCrypto struct{}

// NewAsymmetricCrypto 创建非对称加密
func NewAsymmetricCrypto() *AsymmetricCrypto {
    return &AsymmetricCrypto{}
}

// GenerateRSAKeyPair 生成RSA密钥对
func (a *AsymmetricCrypto) GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
    privateKey, err := rsa.GenerateKey(rand.Reader, bits)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to generate RSA key: %v", err)
    }

    return privateKey, &privateKey.PublicKey, nil
}

// RSAEncrypt RSA加密
func (a *AsymmetricCrypto) RSAEncrypt(plaintext []byte, publicKey *rsa.PublicKey) ([]byte, error) {
    // 使用OAEP填充
    ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, plaintext, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt: %v", err)
    }
    return ciphertext, nil
}

// RSADecrypt RSA解密
func (a *AsymmetricCrypto) RSADecrypt(ciphertext []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
    // 使用OAEP填充
    plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt: %v", err)
    }
    return plaintext, nil
}

// RSASign RSA签名
func (a *AsymmetricCrypto) RSASign(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
    // 计算哈希
    hash := sha256.Sum256(data)

    // 签名
    signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
    if err != nil {
        return nil, fmt.Errorf("failed to sign: %v", err)
    }
    return signature, nil
}

// RSAVerify RSA验证
func (a *AsymmetricCrypto) RSAVerify(data, signature []byte, publicKey *rsa.PublicKey) error {
    // 计算哈希
    hash := sha256.Sum256(data)

    // 验证签名
    err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature)
    if err != nil {
        return fmt.Errorf("signature verification failed: %v", err)
    }
    return nil
}

// ExportPrivateKey 导出私钥
func (a *AsymmetricCrypto) ExportPrivateKey(privateKey *rsa.PrivateKey) (string, error) {
    // 编码私钥
    privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
    privateKeyPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: privateKeyBytes,
    })
    return string(privateKeyPEM), nil
}

// ExportPublicKey 导出公钥
func (a *AsymmetricCrypto) ExportPublicKey(publicKey *rsa.PublicKey) (string, error) {
    // 编码公钥
    publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
    publicKeyPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PUBLIC KEY",
        Bytes: publicKeyBytes,
    })
    return string(publicKeyPEM), nil
}

// ImportPrivateKey 导入私钥
func (a *AsymmetricCrypto) ImportPrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
    block, _ := pem.Decode([]byte(privateKeyPEM))
    if block == nil {
        return nil, fmt.Errorf("failed to decode PEM block")
    }

    privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        return nil, fmt.Errorf("failed to parse private key: %v", err)
    }

    return privateKey, nil
}

// ImportPublicKey 导入公钥
func (a *AsymmetricCrypto) ImportPublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
    block, _ := pem.Decode([]byte(publicKeyPEM))
    if block == nil {
        return nil, fmt.Errorf("failed to decode PEM block")
    }

    publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
    if err != nil {
        return nil, fmt.Errorf("failed to parse public key: %v", err)
    }

    return publicKey, nil
}
```

### 5.3.2 椭圆曲线密码学

**椭圆曲线优势**：
1. **密钥长度短**：相同安全级别下密钥更短
2. **计算效率高**：运算速度更快
3. **内存占用少**：存储空间更小

**Go语言椭圆曲线实现**：

```go
package main

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/sha256"
    "crypto/x509"
    "encoding/pem"
    "fmt"
    "math/big"
)

// ECCrypto 椭圆曲线密码学
type ECCrypto struct{}

// NewECCrypto 创建椭圆曲线密码学
func NewECCrypto() *ECCrypto {
    return &ECCrypto{}
}

// GenerateECKeyPair 生成椭圆曲线密钥对
func (e *ECCrypto) GenerateECKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
    privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to generate EC key: %v", err)
    }

    return privateKey, &privateKey.PublicKey, nil
}

// ECSign 椭圆曲线签名
func (e *ECCrypto) ECSign(data []byte, privateKey *ecdsa.PrivateKey) (*big.Int, *big.Int, error) {
    // 计算哈希
    hash := sha256.Sum256(data)

    // 签名
    r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
    if err != nil {
        return nil, nil, fmt.Errorf("failed to sign: %v", err)
    }

    return r, s, nil
}

// ECVerify 椭圆曲线验证
func (e *ECCrypto) ECVerify(data []byte, r, s *big.Int, publicKey *ecdsa.PublicKey) bool {
    // 计算哈希
    hash := sha256.Sum256(data)

    // 验证签名
    return ecdsa.Verify(publicKey, hash[:], r, s)
}

// ExportECPrivateKey 导出椭圆曲线私钥
func (e *ECCrypto) ExportECPrivateKey(privateKey *ecdsa.PrivateKey) (string, error) {
    // 编码私钥
    privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
    if err != nil {
        return "", fmt.Errorf("failed to marshal private key: %v", err)
    }

    privateKeyPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "EC PRIVATE KEY",
        Bytes: privateKeyBytes,
    })
    return string(privateKeyPEM), nil
}

// ExportECPublicKey 导出椭圆曲线公钥
func (e *ECCrypto) ExportECPublicKey(publicKey *ecdsa.PublicKey) (string, error) {
    // 编码公钥
    publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
    if err != nil {
        return "", fmt.Errorf("failed to marshal public key: %v", err)
    }

    publicKeyPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "PUBLIC KEY",
        Bytes: publicKeyBytes,
    })
    return string(publicKeyPEM), nil
}
```

## 5.4 哈希函数

### 5.4.1 哈希函数特性

**哈希函数要求**：
1. **确定性**：相同输入总是产生相同输出
2. **快速计算**：计算速度要快
3. **抗碰撞性**：难以找到两个不同输入产生相同输出
4. **雪崩效应**：输入微小变化导致输出巨大变化

**常见哈希函数**：
1. **MD5**：128位哈希（已不安全）
2. **SHA-1**：160位哈希（已不安全）
3. **SHA-256**：256位哈希
4. **SHA-3**：新一代哈希函数

**Go语言哈希函数实现**：

```go
package main

import (
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "crypto/sha512"
    "encoding/hex"
    "fmt"
    "golang.org/x/crypto/sha3"
)

// HashFunctions 哈希函数
type HashFunctions struct{}

// NewHashFunctions 创建哈希函数
func NewHashFunctions() *HashFunctions {
    return &HashFunctions{}
}

// MD5Hash MD5哈希
func (h *HashFunctions) MD5Hash(data []byte) string {
    hash := md5.Sum(data)
    return hex.EncodeToString(hash[:])
}

// SHA1Hash SHA1哈希
func (h *HashFunctions) SHA1Hash(data []byte) string {
    hash := sha1.Sum(data)
    return hex.EncodeToString(hash[:])
}

// SHA256Hash SHA256哈希
func (h *HashFunctions) SHA256Hash(data []byte) string {
    hash := sha256.Sum256(data)
    return hex.EncodeToString(hash[:])
}

// SHA512Hash SHA512哈希
func (h *HashFunctions) SHA512Hash(data []byte) string {
    hash := sha512.Sum512(data)
    return hex.EncodeToString(hash[:])
}

// SHA3Hash SHA3哈希
func (h *HashFunctions) SHA3Hash(data []byte) string {
    hash := sha3.Sum256(data)
    return hex.EncodeToString(hash[:])
}

// KeccakHash Keccak哈希
func (h *HashFunctions) KeccakHash(data []byte) string {
    hash := sha3.NewLegacyKeccak256()
    hash.Write(data)
    return hex.EncodeToString(hash.Sum(nil))
}

// HMAC HMAC哈希
func (h *HashFunctions) HMAC(data, key []byte) string {
    hmac := hmac.New(sha256.New, key)
    hmac.Write(data)
    return hex.EncodeToString(hmac.Sum(nil))
}

// PBKDF2 密码派生函数
func (h *HashFunctions) PBKDF2(password, salt []byte, iterations int) string {
    derivedKey := pbkdf2.Key(password, salt, iterations, 32, sha256.New)
    return hex.EncodeToString(derivedKey)
}

// Argon2 Argon2密码哈希
func (h *HashFunctions) Argon2(password, salt []byte) string {
    hash := argon2.IDKey(password, salt, 1, 64*1024, 4, 32)
    return hex.EncodeToString(hash)
}
```

### 5.4.2 密码哈希

**密码哈希要求**：
1. **抗彩虹表攻击**：使用盐值
2. **抗暴力破解**：使用足够多的迭代次数
3. **抗硬件攻击**：使用内存密集型算法

**Go语言密码哈希实现**：

```go
package main

import (
    "crypto/rand"
    "fmt"
    "golang.org/x/crypto/bcrypt"
    "golang.org/x/crypto/scrypt"
)

// PasswordHashing 密码哈希
type PasswordHashing struct{}

// NewPasswordHashing 创建密码哈希
func NewPasswordHashing() *PasswordHashing {
    return &PasswordHashing{}
}

// GenerateSalt 生成盐值
func (p *PasswordHashing) GenerateSalt(length int) ([]byte, error) {
    salt := make([]byte, length)
    _, err := rand.Read(salt)
    if err != nil {
        return nil, fmt.Errorf("failed to generate salt: %v", err)
    }
    return salt, nil
}

// BcryptHash Bcrypt哈希
func (p *PasswordHashing) BcryptHash(password []byte, cost int) (string, error) {
    hash, err := bcrypt.GenerateFromPassword(password, cost)
    if err != nil {
        return "", fmt.Errorf("failed to hash password: %v", err)
    }
    return string(hash), nil
}

// BcryptVerify Bcrypt验证
func (p *PasswordHashing) BcryptVerify(password []byte, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), password)
    return err == nil
}

// ScryptHash Scrypt哈希
func (p *PasswordHashing) ScryptHash(password, salt []byte) (string, error) {
    hash, err := scrypt.Key(password, salt, 16384, 8, 1, 32)
    if err != nil {
        return "", fmt.Errorf("failed to hash password: %v", err)
    }
    return fmt.Sprintf("%x", hash), nil
}

// ScryptVerify Scrypt验证
func (p *PasswordHashing) ScryptVerify(password, salt []byte, hash string) bool {
    computedHash, err := p.ScryptHash(password, salt)
    if err != nil {
        return false
    }
    return computedHash == hash
}
```

## 5.5 数字签名

### 5.5.1 数字签名原理

**数字签名特性**：
1. **认证性**：证明签名者身份
2. **完整性**：确保数据未被篡改
3. **不可否认性**：签名者不能否认签名

**数字签名过程**：
1. **签名**：使用私钥对消息哈希进行签名
2. **验证**：使用公钥验证签名

**Go语言数字签名实现**：

```go
package main

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/pem"
    "fmt"
)

// DigitalSignature 数字签名
type DigitalSignature struct{}

// NewDigitalSignature 创建数字签名
func NewDigitalSignature() *DigitalSignature {
    return &DigitalSignature{}
}

// RSASign RSA数字签名
func (d *DigitalSignature) RSASign(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
    // 计算哈希
    hash := sha256.Sum256(data)

    // 签名
    signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
    if err != nil {
        return nil, fmt.Errorf("failed to sign: %v", err)
    }
    return signature, nil
}

// RSAVerify RSA签名验证
func (d *DigitalSignature) RSAVerify(data, signature []byte, publicKey *rsa.PublicKey) error {
    // 计算哈希
    hash := sha256.Sum256(data)

    // 验证签名
    err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature)
    if err != nil {
        return fmt.Errorf("signature verification failed: %v", err)
    }
    return nil
}

// ECDSASign ECDSA数字签名
func (d *DigitalSignature) ECDSASign(data []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
    // 计算哈希
    hash := sha256.Sum256(data)

    // 签名
    r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
    if err != nil {
        return nil, fmt.Errorf("failed to sign: %v", err)
    }

    // 编码签名
    signature := append(r.Bytes(), s.Bytes()...)
    return signature, nil
}

// ECDSAVerify ECDSA签名验证
func (d *DigitalSignature) ECDSAVerify(data, signature []byte, publicKey *ecdsa.PublicKey) error {
    // 计算哈希
    hash := sha256.Sum256(data)

    // 解码签名
    sigLen := len(signature)
    if sigLen%2 != 0 {
        return fmt.Errorf("invalid signature length")
    }

    r := new(big.Int).SetBytes(signature[:sigLen/2])
    s := new(big.Int).SetBytes(signature[sigLen/2:])

    // 验证签名
    if !ecdsa.Verify(publicKey, hash[:], r, s) {
        return fmt.Errorf("signature verification failed")
    }
    return nil
}
```

## 5.6 密钥管理

### 5.6.1 密钥生命周期

**密钥生命周期**：
1. **生成**：安全生成密钥
2. **分发**：安全分发密钥
3. **存储**：安全存储密钥
4. **使用**：安全使用密钥
5. **更新**：定期更新密钥
6. **销毁**：安全销毁密钥

### 5.6.2 密钥派生

**密钥派生函数**：
1. **PBKDF2**：基于密码的密钥派生函数
2. **Scrypt**：内存密集型密钥派生函数
3. **Argon2**：新一代密钥派生函数

## 5.7 密码学应用

### 5.7.1 安全通信

**安全通信协议**：
1. **TLS/SSL**：传输层安全协议
2. **SSH**：安全外壳协议
3. **IPsec**：IP安全协议

### 5.7.2 数字证书

**数字证书内容**：
1. **公钥**：证书持有者的公钥
2. **身份信息**：证书持有者的身份信息
3. **有效期**：证书的有效期
4. **签名**：CA对证书的签名

---

**相关链接**：
- [1.1 安全基础理论](../01-Security-Foundations.md#1.1-安全基础理论)
- [2.1 网络安全基础理论](../02-Network-Security.md#2.1-网络安全基础理论)
- [3.1 应用安全基础](../03-Application-Security.md#3.1-应用安全基础)

**下一步**：
- [6. 安全运营](../06-Security-Operations.md)
- [7. 威胁情报](../07-Threat-Intelligence.md)
- [8. 安全合规](../08-Security-Compliance.md) 