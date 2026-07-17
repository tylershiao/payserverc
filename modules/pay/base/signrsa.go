package base

import (
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"math/big"

	"github.com/liangdas/mqant/log"
	"github.com/wenzhenxi/gorsa"
)

// 使用SHA256withRSA算法对数据进行签名
func SignWithSHA256(data, privateKey []byte) (string, error) {
	// 解析私钥
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return "", errors.New("invalid private key")
	}
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 计算数据的SHA256哈希值
	hash := sha256.Sum256(data)

	// 对哈希值进行签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, priKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}

	// 将签名结果进行Base64编码
	return base64.StdEncoding.EncodeToString(signature), nil
}

// 使用SHA256withRSA算法对数据进行验签
func VerifyWithSHA256(data, publicKey []byte, signature string) (bool, error) {
	// 解析公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return false, errors.New("invalid public key")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}

	// 计算数据的SHA256哈希值
	hash := sha256.Sum256(data)

	// 对签名进行Base64解码
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}

	// 进行验签操作
	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, hash[:], signatureBytes)
	if err == nil {
		return true, nil
	} else {
		return false, err
	}
}

// 使用SHA512withRSA算法对数据进行签名
func SignWithSHA512(data []byte, privateKey string) (string, error) {
	// 解析私钥
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", errors.New("invalid private key")
	}
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 计算数据的SHA256哈希值
	hash := sha512.Sum512(data)

	// 对哈希值进行签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, priKey, crypto.SHA512, hash[:])
	if err != nil {
		return "", err
	}

	// 将签名结果进行Base64编码
	return base64.StdEncoding.EncodeToString(signature), nil
}

// 使用SHA512withRSA算法对数据进行验签
func VerifyWithSHA512(data []byte, signature string, publicKey string) (bool, error) {
	// 解析公钥
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return false, errors.New("invalid public key")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}

	// 计算数据的SHA512哈希值
	hash := sha512.Sum512(data)

	// 对签名进行Base64解码
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}

	// 进行验签操作
	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA512, hash[:], signatureBytes)
	if err == nil {
		return true, nil
	} else {
		return false, err
	}
}

// 私钥签名
func RsaSignSha1(data []byte, privateKey string) ([]byte, error) {
	h := sha1.New()
	h.Write(data)
	hashed := h.Sum(nil)
	//获取私钥
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("private key error")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA1, hashed)
}

// 公钥验证
func RsaSignVer(data []byte, signature, publicKey string) error {
	sign, _ := base64.StdEncoding.DecodeString(signature)
	hashed := sha1.Sum(data)
	block, _ := pem.Decode([]byte(publicKey))
	// 解析公钥
	pubInterface, _ := x509.ParsePKIXPublicKey(block.Bytes)

	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//验证签名
	return rsa.VerifyPKCS1v15(pub, crypto.SHA1, hashed[:], sign)
}

func RSA_public_decrypt(pubKey *rsa.PublicKey, data []byte) []byte {
	c := new(big.Int)
	m := new(big.Int)
	m.SetBytes(data)
	e := big.NewInt(int64(pubKey.E))
	c.Exp(m, e, pubKey.N)
	out := c.Bytes()
	skip := 0
	for i := 2; i < len(out); i++ {
		if i+1 >= len(out) {
			break
		}
		if out[i] == 0xff && out[i+1] == 0 {
			skip = i + 2
			break
		}
	}
	return out[skip:]
}

// 公钥解密
func RsaDecode(signature, publicKey string) error {
	block, _ := pem.Decode([]byte(publicKey))
	public, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Error("加载公钥错误 %+v\n", err)
		return err
	}
	decodeBytes, _ := base64.StdEncoding.DecodeString(signature)
	var result string
	// 因为加密解密用rsa有长度限值，因此这个地方就看具体限值的多少了
	for len(decodeBytes) > 0 {
		decodePart := decodeBytes[:128]
		plain := RSA_public_decrypt(public.(*rsa.PublicKey), decodePart)
		result += string(plain)
		decodeBytes = decodeBytes[128:]
	}
	// log.Debug("%+v", result)
	if len(result) > 0 {
		return nil
	}
	return errors.New("check sign fail")
}

// 私钥加密
func SignRSA(data string, privateKey string) (string, error) {
	gRsa := gorsa.RSASecurity{}
	gRsa.SetPrivateKey(privateKey)

	rsaData, err := gRsa.PriKeyENCTYPT([]byte(data))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(rsaData), nil
}

func SignRSADecode(data string, publicKey string) (string, error) {
	dataBs, err := base64.RawURLEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	gRsa := gorsa.RSASecurity{}
	if err := gRsa.SetPublicKey(publicKey); err != nil {
		return "", err
	}

	rsaData, err := gRsa.PubKeyDECRYPT(dataBs)
	if err != nil {
		return "", err
	}

	return string(rsaData), nil
}

// 使用SHA256withRSA算法对数据进行签名
func SignWithHmacSHA256(data []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))

	// 写入数据
	h.Write(data)

	// 获取 HMAC 的字节序列
	hashBytes := h.Sum(nil)

	// 转换为十六进制字符串
	hashString := hex.EncodeToString(hashBytes)
	return hashString
}

func SignWithHmacSHA256Base64(data, key string) string {
	// 创建 HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))

	// 计算签名
	sign := mac.Sum(nil)

	// 转 Base64
	return base64.StdEncoding.EncodeToString(sign)
}

func SignWithHmacSHA1Base64(data, key string) string {
	// 创建 HMAC-SHA1
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(data))

	// 计算签名
	sign := mac.Sum(nil)

	// 转 Base64
	return base64.StdEncoding.EncodeToString(sign)
}

func RsaEncrypt(plaintext string, base64PublicKey []byte) (string, error) {
	// keyPem := "-----BEGIN PUBLIC KEY-----\n" + insertLineBreaks(base64PublicKey, 64) + "\n-----END PUBLIC KEY-----"
	block, _ := pem.Decode(base64PublicKey)
	if block == nil {
		return "", errors.New("failed to parse PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("not RSA public key")
	}

	encryptedBytes, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPub, []byte(plaintext))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

// func insertLineBreaks(s string, every int) string {
// 	var b strings.Builder
// 	for i := 0; i < len(s); i += every {
// 		end := i + every
// 		if end > len(s) {
// 			end = len(s)
// 		}
// 		b.WriteString(s[i:end] + "\n")
// 	}
// 	return b.String()
// }
