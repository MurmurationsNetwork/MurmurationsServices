package cryptoutil

import (
	"crypto/sha256"
	"encoding/hex"
)

func GetSHA256(s string) string {
	h := sha256.New()
	defer h.Reset()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
