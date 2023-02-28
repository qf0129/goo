package goo

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var ErrTokenExpired = errors.New("TokenExpired")

// 哈希密码
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// 校验密码
func VerifyPassword(psd string, psdHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(psdHash), []byte(psd))
	return err == nil
}

// 创建令牌
func CreateToken(body string) (string, error) {
	tokenByte := make([]byte, len(body)+8) // 前8位放时间戳
	binary.BigEndian.PutUint64(tokenByte[:8], uint64(time.Now().Unix()))
	for i, v := range body {
		tokenByte[i+8] = byte(v)
	}

	token, err := DesEncrypt(tokenByte, []byte(Config.SecretKey))
	return hex.EncodeToString(token), err
}

// 解析令牌
func ParseToken(token string) (string, error) {
	text, err := hex.DecodeString(token)
	if err != nil {
		return "", err
	}
	tokenStr, err := DesDecrypt([]byte(text), []byte(Config.SecretKey))
	if err != nil {
		return "", err
	}

	tokenTs := int64(binary.BigEndian.Uint64(tokenStr[0:8]))
	nowTs := time.Now().Unix()

	if nowTs > tokenTs+int64(Config.TokenExpiredTime) {
		return "", ErrTokenExpired
	}

	body := string(tokenStr[8:])
	return body, nil
}

// Des算法加密
func DesEncrypt(origData []byte, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())

	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// Des算法解密
func DesDecrypt(crypted []byte, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
