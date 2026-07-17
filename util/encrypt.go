package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

const (
	Key = "N4KM0haowQ2U87Nu"
	// Key         = "ABCDEFGHIJKLMNOP"
	CurrentMode = 2
	blocksize   = 16
)

// AesEncrypt aes加密用不同的模式
func AesEncrypt(data []byte) []byte {
	switch CurrentMode {
	case 1:
		return AesEncryptCBC(data)
	case 2:
		return AesEncryptECB(data)
	case 3:
		return AesEncryptCFB(data)
	default:
		return nil
	}
}
func AesDecrypt(data []byte) (decrypted []byte, err error) {
	switch CurrentMode {
	case 1:
		return AesDecryptCBC(data)
	case 2:
		return AesDecryptECB(data)
	case 3:
		return AesDecryptCFB(data)
	default:
		return nil, errors.New("unknow")
	}
}

// AesEncryptCBC aes加密算法，采用cbc模式
func AesEncryptCBC(data []byte) []byte {
	key := []byte(Key)
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	data = Pkcs5Padding(data, blockSize)                        // 补全码
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式
	encrypted := make([]byte, len(data))                        // 创建数组
	blockMode.CryptBlocks(encrypted, data)                      // 加密
	// hex.EncodeToString(encrypted)
	return []byte(base64.StdEncoding.EncodeToString(encrypted))
	// return encrypted
}

// AesDecryptCBC 解密
func AesDecryptCBC(data []byte) (decrypted []byte, err error) {
	data, err = base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		decrypted = nil
		return
	}
	key := []byte(Key)
	block, _ := aes.NewCipher(key)                              // 分组秘钥
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 加密模式
	decrypted = make([]byte, len(data))                         // 创建数组
	blockMode.CryptBlocks(decrypted, data)                      // 解密
	decrypted = pkcs5UnPadding(decrypted)                       // 去除补全码
	return
}

// AesDecrypt aes加密算法，采用ecb模式
func AesDecryptECB(origData []byte) (decrypted []byte, err error) {

	// key := []byte(Key)
	// tmp := ""
	// tmp, err = url.QueryUnescape(string(origData))
	// if err != nil {
	// 	return
	// }
	origData, err = base64.StdEncoding.DecodeString(string(origData))
	if err != nil {
		return
	}

	cipher, _ := aes.NewCipher(generateKey())
	decrypted = make([]byte, len(origData))

	for bs, be := 0, cipher.BlockSize(); bs < len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], origData[bs:be])
	}

	// trim := 0
	// if len(decrypted) > 0 {
	// 	trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	// }

	// decrypted = decrypted[:trim]
	decrypted = pkcs5UnPadding(decrypted)

	// decrypted, err = base64.StdEncoding.DecodeString(string(decrypted))
	// if err != nil {
	// 	decrypted = nil
	// 	return
	// }
	return
}

func AesDecryptECB1(origData []byte) (decrypted []byte, err error) {
	origData, err = base64.StdEncoding.DecodeString(string(origData))
	if err != nil {
		return
	}

	cipher, _ := aes.NewCipher(generateKey())
	decrypted = make([]byte, len(origData))

	for bs, be := 0, cipher.BlockSize(); bs < len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], origData[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	decrypted = decrypted[:trim]

	decrypted, err = base64.StdEncoding.DecodeString(string(decrypted))
	if err != nil {
		decrypted = nil
		return
	}
	return
}

func generateKey() (genKey []byte) {
	key := []byte(Key)
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

func AesEncryptECB(data []byte) []byte {
	key := []byte(Key)
	block, _ := aes.NewCipher(key)
	// data = []byte(base64.StdEncoding.EncodeToString(data))
	data = Pkcs5Padding(data, block.BlockSize())
	decrypted := make([]byte, len(data))
	size := block.BlockSize()

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Encrypt(decrypted[bs:be], data[bs:be])
	}

	return []byte(base64.StdEncoding.EncodeToString(decrypted))
}

func AesEncryptECB1(data []byte) []byte {
	key := []byte(Key)
	block, _ := aes.NewCipher(key)
	data = []byte(base64.StdEncoding.EncodeToString(data))
	// data = pkcs5Padding(data, block.BlockSize())
	data = noPadding(data)
	decrypted := make([]byte, len(data))
	size := block.BlockSize()

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Encrypt(decrypted[bs:be], data[bs:be])
	}

	return []byte(base64.StdEncoding.EncodeToString(decrypted))
}

// AesEncryptCFB aes加密cfb模式
func AesEncryptCFB(origData []byte) (encrypted []byte) {
	key := []byte(Key)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	encrypted = make([]byte, blocksize+len(origData))
	iv := encrypted[:blocksize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[blocksize:], origData)
	encrypted = []byte(base64.StdEncoding.EncodeToString(encrypted))
	return
}

func AesDecryptCFB(data []byte) (decrypted []byte, err error) {
	key := []byte(Key)
	block, _ := aes.NewCipher(key)
	if len(data) < blocksize {
		decrypted = nil
		err = errors.New("ciphertext too short")
		return
	}
	encrypted, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		decrypted = nil
		return
	}
	iv := encrypted[:blocksize]
	encrypted = encrypted[blocksize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding) //用0去填充
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}

func Pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// nopadding模式
func noPadding(src []byte) []byte {
	count := blocksize - len(src)%blocksize
	if len(src)%blocksize == 0 {
		return src
	} else {
		return append(src, bytes.Repeat([]byte{byte(0)}, count)...)
	}
}

// nopadding模式
func unNoPadding(src []byte) []byte {
	for i := len(src) - 1; ; i-- {
		if src[i] != 0 {
			return src[:i+1]
		}
	}
}
