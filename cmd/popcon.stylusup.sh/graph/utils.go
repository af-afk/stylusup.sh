package graph

import (
	"encoding/hex"
	"crypto/hmac"
	"crypto/sha256"
)

func hashIpAddr(ip, arch, lang, os string) string {
	m := hmac.New(sha256.New, []byte(arch + lang + os))
	m.Write([]byte(ip))
	return hex.EncodeToString(m.Sum(nil))
}
