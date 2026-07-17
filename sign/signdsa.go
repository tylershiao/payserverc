package sign

import (
	"crypto/dsa"
	"crypto/rand"
	"crypto/sha1"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
)

// 签名结构体
type dsaSignature struct {
	R, S *big.Int
}

// 公钥结构体（ASN.1 编码）
type dsaPublicKey struct {
	P, Q, G, Y *big.Int
}

// 私钥结构体（ASN.1 编码）
type dsaPrivateKey struct {
	Public dsaPublicKey
	X      *big.Int
}

// GenerateKeyPairBase64 生成 DSA 密钥对（返回 Base64 编码）
func GenerateKeyPairBase64() (string, string, error) {
	var params dsa.Parameters
	if err := dsa.GenerateParameters(&params, rand.Reader, dsa.L1024N160); err != nil {
		return "", "", err
	}

	priv := new(dsa.PrivateKey)
	priv.Parameters = params
	if err := dsa.GenerateKey(priv, rand.Reader); err != nil {
		return "", "", err
	}

	pubEnc := dsaPublicKey{
		P: priv.P,
		Q: priv.Q,
		G: priv.G,
		Y: priv.Y,
	}
	privEnc := dsaPrivateKey{
		Public: pubEnc,
		X:      priv.X,
	}

	pubBytes, _ := asn1.Marshal(pubEnc)
	privBytes, _ := asn1.Marshal(privEnc)

	pubB64 := base64.StdEncoding.EncodeToString(pubBytes)
	privB64 := base64.StdEncoding.EncodeToString(privBytes)
	return pubB64, privB64, nil
}

// 解码 Base64 公钥
// func decodePublicKey(pubB64 string) (*dsa.PublicKey, error) {
// 	data, err := base64.StdEncoding.DecodeString(pubB64)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var pk dsaPublicKey
// 	if _, err := asn1.Unmarshal(data, &pk); err != nil {
// 		return nil, err
// 	}

// 	pub := &dsa.PublicKey{}
// 	pub.Y = pk.Y
// 	pub.Parameters.P = pk.P
// 	pub.Parameters.Q = pk.Q
// 	pub.Parameters.G = pk.G
// 	return pub, nil
// }

// 解码 Base64 私钥
// func decodePrivateKey(privB64 string) (*dsa.PrivateKey, error) {
// 	data, err := base64.StdEncoding.DecodeString(privB64)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var sk dsaPrivateKey
// 	if _, err := asn1.Unmarshal(data, &sk); err != nil {
// 		return nil, err
// 	}
// 	priv := &dsa.PrivateKey{}
// 	priv.X = sk.X
// 	priv.PublicKey.Y = sk.Public.Y
// 	priv.PublicKey.Parameters.P = sk.Public.P
// 	priv.PublicKey.Parameters.Q = sk.Public.Q
// 	priv.PublicKey.Parameters.G = sk.Public.G
// 	return priv, nil
// 	// key, err := x509.ParsePKCS8PrivateKey(data)
// 	// if err != nil {
// 	// 	panic("ParsePKCS8PrivateKey error: " + err.Error())
// 	// }

// 	// // 3. 断言为DSA私钥类型
// 	// dsaPrivKey, ok := key.(*dsa.PrivateKey)
// 	// if !ok {
// 	// 	panic("Not a DSA private key")
// 	// }
// 	// return dsaPrivKey, nil
// }

// SignDsa 使用 Base64 私钥对消息签名
func SignDsa(message, privB64 string) (string, error) {
	priv, err := decodePKCS8DSAPrivateKey(privB64)
	if err != nil {
		return "", err
	}
	hash := sha1.Sum([]byte(message))
	r, s, err := dsa.Sign(rand.Reader, priv, hash[:])
	if err != nil {
		return "", err
	}
	sigBytes, _ := asn1.Marshal(dsaSignature{r, s})
	return base64.StdEncoding.EncodeToString(sigBytes), nil
}

// VerifyDsa 使用 Base64 公钥验证签名
func VerifyDsa(message, sigB64, pubB64 string) (bool, error) {
	pub, err := decodeDSAPublicKey(pubB64)
	if err != nil {
		return false, err
	}
	sigBytes, err := base64.StdEncoding.DecodeString(sigB64)
	if err != nil {
		return false, err
	}
	var sig dsaSignature
	if _, err := asn1.Unmarshal(sigBytes, &sig); err != nil {
		return false, err
	}
	hash := sha1.Sum([]byte(message))
	return dsa.Verify(pub, hash[:], sig.R, sig.S), nil
}

// ExampleMain 测试入口
func ExampleMain() {
	message := "288889a170964134871008Dsr5c1729585712879payout"

	pub, priv, err := GenerateKeyPairBase64()
	if err != nil {
		panic(err)
	}
	fmt.Println("Public Key (Base64):", pub)
	fmt.Println("Private Key (Base64):", priv)

	sig, err := SignDsa(message, priv)
	if err != nil {
		panic(err)
	}
	fmt.Println("Signature (Base64):", sig)

	valid, err := VerifyDsa(message, sig, pub)
	if err != nil {
		panic(err)
	}
	fmt.Println("Signature is valid:", valid)
}

// PKCS#8 最外层结构
type pkcs8PrivateKey struct {
	Version    int
	Algo       pkix.AlgorithmIdentifier
	PrivateKey []byte // OCTET STRING
}

// DSA 算法参数 Dss-Parms ::= SEQUENCE { p, q, g }
type dssParms struct {
	P, Q, G *big.Int
}

func decodePKCS8DSAPrivateKey(b64 string) (*dsa.PrivateKey, error) {
	// 1) Base64 解码 DER
	der, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}

	// 2) 解析外层 PKCS#8
	var p8 pkcs8PrivateKey
	if _, err := asn1.Unmarshal(der, &p8); err != nil {
		return nil, fmt.Errorf("unmarshal PKCS#8: %w", err)
	}
	// 确保是 DSA OID
	if !p8.Algo.Algorithm.Equal(asn1.ObjectIdentifier{1, 2, 840, 10040, 4, 1}) {
		return nil, errors.New("unexpected algorithm OID, not DSA")
	}

	// 3) 从 Algo.Parameters 解出 P, Q, G
	var params dssParms
	if _, err := asn1.Unmarshal(p8.Algo.Parameters.FullBytes, &params); err != nil {
		return nil, fmt.Errorf("unmarshal DSA parameters: %w", err)
	}

	// 4) 从 PrivateKey OCTET STRING 解出 X（私钥整数）
	var x *big.Int
	if _, err := asn1.Unmarshal(p8.PrivateKey, &x); err != nil {
		return nil, fmt.Errorf("unmarshal DSA private integer: %w", err)
	}

	// 5) 计算 Y = G^X mod P
	y := new(big.Int).Exp(params.G, x, params.P)

	// 6) 构造 dsa.PrivateKey
	priv := &dsa.PrivateKey{
		PublicKey: dsa.PublicKey{
			Parameters: dsa.Parameters{
				P: params.P,
				Q: params.Q,
				G: params.G,
			},
			Y: y,
		},
		X: x,
	}
	return priv, nil
}

// SPKI 结构
type subjectPublicKeyInfo struct {
	Algo      pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
}

func decodeDSAPublicKey(b64 string) (*dsa.PublicKey, error) {
	// 1. 解码 base64
	der, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("base64 decode: %w", err)
	}

	// 2. 解码 SPKI 外层
	var spki subjectPublicKeyInfo
	if _, err := asn1.Unmarshal(der, &spki); err != nil {
		return nil, fmt.Errorf("unmarshal SPKI: %w", err)
	}

	// 3. 验证算法 OID 是 DSA
	if !spki.Algo.Algorithm.Equal(asn1.ObjectIdentifier{1, 2, 840, 10040, 4, 1}) {
		return nil, fmt.Errorf("unexpected algorithm OID (not DSA)")
	}

	// 4. 解析 DSA 参数
	var params dssParms
	if _, err := asn1.Unmarshal(spki.Algo.Parameters.FullBytes, &params); err != nil {
		return nil, fmt.Errorf("unmarshal DSA parameters: %w", err)
	}

	// 5. 解析 BIT STRING 中的 Y
	var y *big.Int
	if _, err := asn1.Unmarshal(spki.PublicKey.RightAlign(), &y); err != nil {
		return nil, fmt.Errorf("unmarshal public key Y: %w", err)
	}

	// 6. 构造 dsa.PublicKey
	pub := &dsa.PublicKey{
		Parameters: dsa.Parameters{
			P: params.P,
			Q: params.Q,
			G: params.G,
		},
		Y: y,
	}
	return pub, nil
}
