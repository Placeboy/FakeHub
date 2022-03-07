package rsa

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"testing"
)

// 加密以后的长度为128
func TestRsaEncrypt(t *testing.T) {
	plaintext := "readTemperatursadasdasdasdae"
	ciphertext, _ := RsaEncrypt([]byte(plaintext))
	fmt.Println(ciphertext)
	fmt.Println(len(ciphertext))
}
// 以16进制形式打印数组
func PrintBytesArrayInHexadecimal(name string, b []byte) {
	length := len(b)
	fmt.Println("Array length = ", length)
	fmt.Printf("%s = {", name)
	for i := 0; i < length-1; i++ {
		fmt.Printf("0x%02X, ", b[i])
	}
	fmt.Printf("0x%02X}\n", b[length-1])
}

func TestParseFromPublicKey(t *testing.T) {
	length := 128
	e, n, _ := ParseFromPublicKey()
	fmt.Printf("e = %d, n = %d\n", e, n)
	//fmt.Printf("%d\n", n.Mod(n, big.NewInt(16)))
	buf := make([]byte, length)
	n.FillBytes(buf)
	//fmt.Println(buf)
	PrintBytesArrayInHexadecimal("Modulus", buf)
	temp := big.NewInt(int64(e))
	buf = temp.Bytes()
	PrintBytesArrayInHexadecimal("PublicExponent", buf)
}

func TestParseFromPrivateKey(t *testing.T) {
	length := 128
	d, e, n, _ := ParseFromPrivateKey()
	//fmt.Printf("e = %d, n = %d\n", e, n)
	//fmt.Printf("%d\n", n.Mod(n, big.NewInt(16)))
	buf_n := make([]byte, length)
	n.FillBytes(buf_n)
	//fmt.Println(buf)
	PrintBytesArrayInHexadecimal("Modulus", buf_n)
	buf_d := d.Bytes()
	//d.FillBytes(buf)
	PrintBytesArrayInHexadecimal("PrivateExponent", buf_d)
	temp := big.NewInt(int64(e))
	buf_e := temp.Bytes()
	PrintBytesArrayInHexadecimal("PublicExponent", buf_e)
}

//func TestBigInit(t *testing.T) {
//	n :=
//}

func TestRsaDecrypt(t *testing.T) {
	cryptograph := "dwVS3CAz3pULMJr/olxK3uWrnzltc9657mIED0wycZX66tx7z0eoirgvR+17pR0ppkuCUKfRdPsjGCTJnHvagiu+eZ9DMhU+ndlsBoOki5ZAJ0BiGex5EOWSenIUUVXbBjwgwrOqfAgc+YiUMTrlqtTCA8CPxRmK+4fxchaQQ64="
	decodeBytes, err := base64.StdEncoding.DecodeString(cryptograph)
	if err != nil {
		log.Fatalln(err)
	}
	origData, _ := RsaDecrypt(decodeBytes)
	fmt.Println(string(origData))
}

func TestRsaDecryptNoPadding(t *testing.T) {
	cryptograph := "G3Nvh8Ohlt+iepQ/V5OYp27oJFoxMO/wERCjsQvbOC3Y1E3QeXIMttTnvIsaxW6eS1autVFrsoa6rJtpN9mup+hVuIZ8+UcHaPPawQgOfHZSyNBdRF82BvbagZsrQO/S53KTL1mG5VdXuQ8vtGwwA3oErC12sDQG8WZcAicLIOs="
	cipherText, err := base64.StdEncoding.DecodeString(cryptograph)
	if err != nil {
		log.Fatalln(err)
	}
	origData, _ := RsaDecryptNoPadding(cipherText)
	plaintext := string(origData)
	fmt.Println(plaintext[:16])
}
