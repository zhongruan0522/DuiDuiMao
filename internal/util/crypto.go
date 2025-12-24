package util

import "encoding/base64"

// DoubleEncode 双重Base64加密
func DoubleEncode(data string) string {
	// 第一次编码
	first := base64.StdEncoding.EncodeToString([]byte(data))
	// 第二次编码
	second := base64.StdEncoding.EncodeToString([]byte(first))
	return second
}

// DoubleDecode 双重Base64解密
func DoubleDecode(encoded string) (string, error) {
	// 第一次解码
	first, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	// 第二次解码
	second, err := base64.StdEncoding.DecodeString(string(first))
	if err != nil {
		return "", err
	}
	return string(second), nil
}
