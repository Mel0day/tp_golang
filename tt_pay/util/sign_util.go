package util

import (
	"context"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"sort"
	"strings"
)

func BuildMd5WithRsa(signMap map[string]interface{}, privateKey string) (string, error) {
	signStr := GenSignStr(signMap)

	rsaSign, err := RsaSign(signStr, privateKey)
	if err != nil {
		Debug("BuildMd5WithRsa RsaSign sign str[%s] err[%s]\n", signStr, err)
		return "", err
	}

	return rsaSign, nil
}

func VerifyMd5WithRsa(signMap map[string]interface{}, sign, publicKey string) bool {
	var signStr string
	signStr = GenSignStr(signMap)

	return RsaVerify(signStr, sign, publicKey)
}

func RsaVerify(target, sign, publicKey string) bool {
	//
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil || block.Type != "PUBLIC KEY" {
		Debug("RsaVerify got invalid public key[%s]\n", publicKey)
		return false
	}

	// FIXME: ParsePKCS1PublicKey ???
	pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		Debug("RsaVerify ParsePKIXPublicKey publicKey[%s] err[%s]\n", publicKey, err)
		return false
	}

	pubRsaKey, ok := pubKeyInterface.(*rsa.PublicKey)
	if !ok {
		Debug("RsaVerify publicKey[%s] type is not *rsa.PublicKey\n", publicKey)
		return false
	}

	targetHash := md5.Sum([]byte(target))
	signByte, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		Debug("RsaVerify base64 decode sign[%s] err[%s]\n", sign, err)
		return false
	}

	err = rsa.VerifyPKCS1v15(pubRsaKey, crypto.MD5, targetHash[:], signByte)
	if err != nil {
		Debug("RsaVerify target[%s] sign[%s] public key[%s] failed\n", target, sign, publicKey)
		return false
	}

	return true
}

// 返回加密后base64数据
func RsaSign(target, privateKey string) (string, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		Debug("RsaSign got invalid private key[%s]\n", privateKey)
		Debug("pem decode private block is nil[%v]\n", block == nil)
		if block != nil {
		Debug("pem decode private block type[%s] not RSA PRIVATE KEY\n", block.Type)
		}
		return "", errors.New("decode private error")
	}

	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		Debug("RsaSign x509.ParsePKCS1PrivateKey privateKey[%s] err[%s]\n", privateKey, err)
		return "", err
	}

	targetHash := md5.Sum([]byte(target))

	signed, err := rsa.SignPKCS1v15(rand.Reader, priKey, crypto.MD5, targetHash[:])
	if err != nil {
		Debug("RsaSign target[%s] private key[%s] failed\n", target, privateKey)
		return "", err
	}

	base64Sign := base64.StdEncoding.EncodeToString(signed)
	Debug("RsaSign target[%s] privateKey[%s] to signed[%s]\n", target, privateKey, base64Sign)
	return base64Sign, nil
}

func BuildMd5WithSalt(signMap map[string]interface{}, salt string) string {
	signStr := GenSignStr(signMap)
	finStr := signStr + salt
	sign := fmt.Sprintf("%x", md5.Sum([]byte(finStr)))
	Debug("BuildMd5WithSalt get sign[%s] from sign map[%#v] salt[%s] sign str[%s]\n",
		sign, signMap, salt, finStr)
	return sign
}

func BuildMd5WithSaltCtx(ctx context.Context, signMap map[string]interface{}, salt string) string {
	signStr := GenSignStr(signMap)
	finStr := signStr + salt
	sign := fmt.Sprintf("%x", md5.Sum([]byte(finStr)))
	Debug( "Context: [%s], BuildMd5WithSalt get sign[%s] from sign map[%#v] salt[%s] sign str[%s]\n",
		ctx, sign, signMap, salt, finStr)
	return sign
}

func BuildSha1WithSalt(signMap map[string]interface{}, salt string) string {
	signStr := GenSignStr(signMap)
	finStr := signStr + salt
	sign := fmt.Sprintf("%x", sha1.Sum([]byte(finStr)))
	Debug("BuildSha1WithSalt get sign[%s] from sign map[%#v] salt[%s] sign str[%s]\n",
		sign, signMap, salt, finStr)
	return sign
}

// 去掉请求参数中的字节类型字段（如文件、字节流）、sign字段、值为空的字段
// 目前可认为仅有string，int等简单类型
func GenSignStr(signMap map[string]interface{}) string {
	signArr := make([]string, 0, len(signMap))
	for key, val := range signMap {
		if key == "sign" {
			continue
		}
		finalVal := fmt.Sprintf("%v", val)
		if finalVal == "" {
			continue
		}

		signArr = append(signArr, fmt.Sprintf("%v=%v", key, finalVal))
	}

	sort.Strings(signArr)
	ret := strings.Join(signArr, "&")
	Debug("GenSignStr get ret str[%s] from params[%#v]\n", ret, signMap)
	return ret
}

// 去掉请求参数中的字节类型字段（如文件、字节流）、sign字段、值为空的字段
// 目前可认为仅有string，int等简单类型
func GenSignStrWithSign(signMap map[string]interface{}) string {
	signArr := make([]string, 0, len(signMap))
	for key, val := range signMap {
		finalVal := fmt.Sprintf("%v", val)
		if finalVal == "" {
			continue
		}

		signArr = append(signArr, fmt.Sprintf("%v=%v", key, finalVal))
	}

	sort.Strings(signArr)
	ret := strings.Join(signArr, "&")
	Debug("GenSignStr get ret str[%s] from params[%#v]\n", ret, signMap)
	return ret
}

func GetIMapDefaultString(data map[string]interface{}, key, val string) (ret string) {
	ret = val
	intfaceVal, exist := data[key]
	if exist {
		realVal, ok := intfaceVal.(string)
		if ok {
			ret = realVal
		}
	}

	return ret
}

// // 获取trade_create需要验签的参数
// func GetTradeCreatSIMap(data map[string]string, params map[string]string) map[string]interface{} {
// //兼容内部case没有传version的情况
// version, ok := params["version"]
// if !ok {
// version = "1.0"
// }
// ret := make(map[string]interface{})
// signKeys, ok := consts.TradeCreateSignKeysVersionCtrl[version]
// if !ok {
// return ret
// }
// for key, val := range data {
// if _, ok := signKeys[key]; ok {
// ret[key] = val
// }
// }

// return ret
// }

// GetSIMap 获取签名参数
func GetSIMap(signKeys map[string]int, data map[string]string) map[string]interface{} {
	ret := make(map[string]interface{})
	for key, val := range data {
		if _, ok := signKeys[key]; ok {
			ret[key] = val
		}
	}

	return ret
}

func VerifySign(ctx context.Context, data map[string]interface{}, secret string) bool {
	reqSign := data["sign"]
	signType, ok := data["sign_type"]
	if !ok {
		Debug("Context: [%s], sign type not found\n", ctx)
		return false
	}

	delete(data, "sign")
	sign := ""
	switch signType {
	case "MD5":
		// sign = BuildMd5WithSalt(data, secret)
		sign = BuildMd5WithSaltCtx(ctx, data, secret)
	case "SHA":
		sign = BuildSha1WithSalt(data, secret)
	case "MD5withRSA":
		// FIXME: 获取公钥
		// 此处代码有问题
		rsaSign, error := BuildMd5WithRsa(data, secret) //TODO RSA private key
		if error != nil {
			Debug( "Context: [%s], build rsa sign error %v\n", ctx, error)
			return false
		}
		sign = rsaSign
	default:
		Debug("bad sign_type\n", ctx, reqSign, sign)
		return false
	}

	if sign != reqSign {
		Debug("Context: [%s], sign error req sign %s, computed sign %s\n", ctx, reqSign, sign)
		return false
	}

	return true
}

func BuildRsaSign(params map[string]interface{}, privateKey string, digest string) string {

	signStr := GenSignStrWithSign(params)
	if digest == "sha1" {
		signStr = fmt.Sprintf("%x", sha1.Sum([]byte(signStr)))
	}

	// TODO: check error
	signedStr, _ := RsaSignWithSha1(signStr, privateKey)

	return signedStr
}

func RsaSignWithSha1(target, privateKey string) (string, error) {
	Debug("signContent : [%s]\n", target)
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		Debug("RsaSign got invalid private key[%s]\n", privateKey)
		Debug("pem decode private block is nil[%v]\n", block == nil)
		if block != nil {
			Debug("pem decode private block type[%s] not RSA PRIVATE KEY\n", block.Type)
		}
		return "", errors.New("decode private error")
	}

	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		Debug("RsaSign x509.ParsePKCS1PrivateKey privateKey[%s] err[%s]\n", privateKey, err)
		return "", err
	}

	targetHash := sha1.Sum([]byte(target))

	signed, err := rsa.SignPKCS1v15(rand.Reader, priKey, crypto.SHA1, targetHash[:])
	if err != nil {
		Debug("RsaSign target[%s] private key[%s] failed\n", target, privateKey)
		return "", err
	}

	base64Sign := base64.StdEncoding.EncodeToString(signed)
	Debug("RsaSign target[%s] privateKey[%s] to signed[%s]\n", target, privateKey, base64Sign)
	return base64Sign, nil
}
