package xenon

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
)

func EncodeMD5(unencrypted interface{}) (string, error) {
	switch t := unencrypted.(type) {
	case string:
		return fmt.Sprintf("%x", md5.Sum([]byte(unencrypted.(string)))), nil
	case multipart.File:
		file := unencrypted.(multipart.File)
		defer file.Seek(0, 0)
		hash := md5.New()
		_, err := io.Copy(hash, file)
		if err != nil {
			return "", err
		}
		return hex.EncodeToString(hash.Sum(nil)), nil
	default:
		return "", errors.New(fmt.Sprintf("md5 sum param data type error [%t]", t))
	}
}


func EncodeAesWithCommonKey(unencrypted string, commonKey string) (string, error) {
	block, err := aes.NewCipher([]byte(commonKey))
	if err != nil {
		return "", err
	}
	cipherText := make([]byte, aes.BlockSize+len(unencrypted))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cipher.NewCFBEncrypter(block, iv).XORKeyStream(cipherText[aes.BlockSize:],
		[]byte(unencrypted))
	return hex.EncodeToString(cipherText), nil

}
func DecodeAesWithCommonKey(encrypted string, commonKey string) (string, error) {
	cipherText, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher([]byte(commonKey))
	if err != nil {
		return "", err
	}
	if len(cipherText) < aes.BlockSize {
		return "", errors.New("cipherText too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	cipher.NewCFBDecrypter(block, iv).XORKeyStream(cipherText, cipherText)
	return string(cipherText), nil
}

func DecodeBase64(base64String string) ([]byte, error) {
	base64String = strings.TrimSpace(base64String)
	if strings.ContainsAny(base64String, "-_") {
		data, err := base64.RawURLEncoding.DecodeString(base64String)
		if err == nil {
			return data, nil
		}
		return base64.URLEncoding.DecodeString(base64String)
	}

	if len(base64String) % 4 != 0 {
		data, err := base64.RawStdEncoding.DecodeString(base64String)
		if err == nil {
			return data, nil
		}
		for len(base64String) % 4 != 0 {
			base64String += "="
		}
	}

	return base64.StdEncoding.DecodeString(base64String)
}