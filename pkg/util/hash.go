package util

import (
	"crypto/sha256"
	"fmt"
)

func Sha256(text string) string {
	encode := sha256.Sum256([]byte(text))
	return fmt.Sprintf("%x", encode)
}
