package xenon

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
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
