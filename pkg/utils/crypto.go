package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)


func MD5(inp string) string {
	if inp == "" {
		return ""
	}
	h := md5.New()
	h.Write([]byte(inp))
	return hex.EncodeToString(h.Sum(nil))
}

func SHA256(str, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
    mac.Write([]byte(str))
    return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}