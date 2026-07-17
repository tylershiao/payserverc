package sign

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/liangdas/mqant/log"
)

func SignCBC(dataStr string, key []byte, iv ...byte) string {
	data := []byte(dataStr)
	// 创建cipher.Block接口
	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}
	// 填充
	blockSize := block.BlockSize()
	if iv == nil {
		iv = make([]byte, blockSize)
	}
	data = pad(data, blockSize)
	// CBC模式加密
	ciphertext := make([]byte, len(data))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, data)

	return base64.StdEncoding.EncodeToString(ciphertext)
}

func SignCBCDecode(dataStr string, key []byte, iv ...byte) string {
	data, err := base64.StdEncoding.DecodeString(dataStr)
	if err != nil {
		log.Error("err:%v", err)
		return ""
	}
	//
	// 创建cipher.Block接口
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error("err:%v", err)
		return ""
	}
	// 填充
	blockSize := block.BlockSize()
	if iv == nil {
		iv = make([]byte, blockSize)
	}

	// CBC模式加密
	mode := cipher.NewCBCDecrypter(block, iv)
	ciphertext := make([]byte, len(data))
	mode.CryptBlocks(ciphertext, data)

	return string(ciphertext)
}

func SignCBCBase64(data string, okey string, aesIv string) string {
	// 将Base64编码的密钥解码
	key, err := base64.StdEncoding.DecodeString(okey)
	if err != nil {
		return ""
	}

	// 创建AES加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}

	// 将Base64编码的IV解码
	iv, err := base64.StdEncoding.DecodeString(aesIv)
	if err != nil {
		return ""
	}

	// 创建CBC加密模式实例
	mode := cipher.NewCBCEncrypter(block, iv)

	// 对数据进行填充操作
	paddedData := pad([]byte(data), block.BlockSize())

	// 加密数据
	encryptedData := make([]byte, len(paddedData))
	mode.CryptBlocks(encryptedData, paddedData)

	// 将加密结果进行Base64编码
	encodedResult := base64.StdEncoding.EncodeToString(encryptedData)

	return encodedResult
}

func pad(data []byte, blockSize int) []byte {
	padLen := blockSize - len(data)%blockSize
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(data, padding...)
}

func GetApiCheckSum(saltKey string, dataString string) string {
	saltKeyBytes := []byte(saltKey)
	saltKeyEncoded := base64.StdEncoding.EncodeToString(saltKeyBytes)
	hash := hmac.New(sha512.New, []byte(saltKeyEncoded))
	hash.Write([]byte(dataString))
	hashValue := hash.Sum(nil)
	hashValueHex := hex.EncodeToString(hashValue)
	hashValueHex = strings.ToUpper(hashValueHex)
	return hashValueHex
}

func SignCBCBase64Decode(data string, okey string, aesIv string) string {
	d, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Error("err:%v", err)
		return ""
	}
	key, err := base64.StdEncoding.DecodeString(okey)
	if err != nil {
		log.Error("err:%v", err)
		return ""
	}
	iv, err := base64.StdEncoding.DecodeString(aesIv)
	if err != nil {
		log.Error("err:%v", err)
		return ""
	}
	block, err := aes.NewCipher(key) // 分组秘钥
	if err != nil {
		log.Error("err:%v", err)
		return ""
	}
	blockMode := cipher.NewCBCDecrypter(block, iv) // 加密模式
	blockMode.CryptBlocks(d, d)                    // 解密
	d = pkcs5UnPadding(d)                          // 去除补全码

	return string(d)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
