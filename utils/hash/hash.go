package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}
