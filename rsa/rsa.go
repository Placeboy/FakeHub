package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"
)

// 此文件中主要是用于加密解密的逻辑

// 可通过openssl产生
// openssl genrsa -out rsa_private_key.pem 1024
var privateKey = []byte(`  
-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQC6BZb8ya3cIrBli22TcHFFdZ5Sl5FHOcovGsO43lJTTrzPD1jS
UZ2dQEpWfRku9XzmSsq3x7tdmSPhKNRruERHfhpigX2kPGHGTS+HHCRDX2N6nGMF
dwW8jzdIS7dIAsX5vk8BPnXNTA+jr/lMx5WcenHoW1mAsTqq5pvKndwKnQIDAQAB
AoGAeNZgul0YP0OZap0j1P7Z1dENw4EJskbr+6VbNp/UwqEHLUo+3IB/7kJxB7XD
wildtQsonDF2mNp94Clxs3fDgcH2mkZFqJ7eh52Q+HIEch9hRzLrb03pfEdFgGXn
QkOJ5nWuNx2fu0ME2ATM5oZz6CHnz6KESmAnOCRK8/L5bwECQQDb+9a/8u1kTUCx
63LoXM5S7ub7T4YkQ8988J06bViT6x96X/bEKCf2kP6uUsOwjDptiRycv6Hj/ovW
XKttChbxAkEA2HpO3ybmQAUWuEt5qltcTH0yQo383650JorWZCecMQlobF0IPeNf
axgQDt/VrgdxdS5ZVSzziHTo2aB24LCmbQJAULvYUJHjNdB0UdfLUCPfROiQtOK2
pFCOsZfM3EiNHZxI7SyS7+Kc6AzGq0uMrhqIxvJvIcfirj4ZLA7OizIMwQJAHQLv
NQrSirvj2pkK2iDaUsnohXDf9d48ZLnwl4WTciLvoq4pH5osPH8CD+xBh8wpkWm/
wSGAFcaNOjU+GUizVQJAGiTuNkZfNZLiDAov9mCUt62r6pDrzoKW3eSjxPS3QI+/
fibFejzs5AqGnRPGfdxHCmLzPzhlggV5SiULjpslgg==
-----END RSA PRIVATE KEY-----
`)

// 先用openssl生成pkcs8公钥
//openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem
// 然后将pkcs8公钥转pkcs1公钥
//openssl rsa -pubin -in rsa_public_key.pem -RSAPublicKey_out
// pkcs1格式公钥
var publicKey = []byte(`  
-----BEGIN RSA PUBLIC KEY-----
MIGJAoGBALoFlvzJrdwisGWLbZNwcUV1nlKXkUc5yi8aw7jeUlNOvM8PWNJRnZ1A
SlZ9GS71fOZKyrfHu12ZI+Eo1Gu4REd+GmKBfaQ8YcZNL4ccJENfY3qcYwV3BbyP
N0hLt0gCxfm+TwE+dc1MD6Ov+UzHlZx6cehbWYCxOqrmm8qd3AqdAgMBAAE=
-----END RSA PUBLIC KEY-----   
`)

// pkcs8格式公钥
//-----BEGIN PUBLIC KEY-----
//MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC6BZb8ya3cIrBli22TcHFFdZ5S
//l5FHOcovGsO43lJTTrzPD1jSUZ2dQEpWfRku9XzmSsq3x7tdmSPhKNRruERHfhpi
//gX2kPGHGTS+HHCRDX2N6nGMFdwW8jzdIS7dIAsX5vk8BPnXNTA+jr/lMx5WcenHo
//W1mAsTqq5pvKndwKnQIDAQAB
//-----END PUBLIC KEY-----

func ParseFromPublicKey() (int, *big.Int, error){
	//解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return 0, big.NewInt(0), errors.New("public key error")
	}
	// 解析公钥
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return 0,big.NewInt(0), err
	}
	//fmt.Println(pub.N.BitLen())
	//buf := make([]byte, 1024)
	//pub.N.FillBytes(buf)
	//fmt.Println(buf)
	//str := pub.N.String()
	//fmt.Println(str)
	return pub.E, pub.N, nil
}

func ParseFromPrivateKey() (*big.Int, int, *big.Int, error){
	//解密pem格式的公钥
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return big.NewInt(0), 0, big.NewInt(0), errors.New("public key error")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		return big.NewInt(0), 0, big.NewInt(0), err
	}
	//fmt.Printf("N = %d, E = %d, D= %d\n", priv.N, priv.E, priv.D)
	return priv.D, priv.E, priv.N, nil
}

// 加密
func RsaEncrypt(origData []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	//pub := pubInterface.(*rsa.PublicKey)
	//fmt.Printf("N = %d, E = %d\n", pub.N, pub.E)
	//fmt.Printf("N = %d, E = %d\n", pub.N.Int64(), pub.E)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func RsaDecrypt(ciphertext []byte) ([]byte, error) {
	//解密
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("N = %d, E = %d, D= %d\n", priv.N, priv.E, priv.D)
	//fmt.Printf("N = %d, E = %d\n", pub.N.Int64(), pub.E)
	// 解密
	//rsa.DecryptPKCS1v15SessionKey()
	//return rsa.DecryptPKCS1v15(nil, priv, ciphertext)
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

func RsaDecryptNoPadding(ciphertext []byte) ([]byte, error) {
	//解密
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	c := new(big.Int).SetBytes(ciphertext)
	plainText := c.Exp(c, priv.D, priv.N).Bytes()
	// 解密
	return plainText, nil
}

