package encrypt

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(password string) (string, error) {
	encryptString, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(encryptString), nil
}

func Verify(hashPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}

func MD5Encrypt(text string) string {
	hash := md5.New()                    // 创建MD5哈希对象
	hash.Write([]byte(text))             // 写入数据
	hashBytes := hash.Sum(nil)           // 计算哈希值
	return hex.EncodeToString(hashBytes) // 转为16进制字符串
}
