package cryptoutil

import (
	"crypto/sha256"
	"encoding/hex"
)

// ComputeSHA256 computes the SHA-256 hash of the given input string.
// If the input is a JSON string, please use jsonutil.Hash.
func ComputeSHA256(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}
