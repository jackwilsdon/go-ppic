package ppic

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// hashString hashes the provided string into an integer.
func hashString(s string) int64 {
	m := sha256.New()

	// Write our string to the SHA256 hash calculator.
	fmt.Fprint(m, s)

	// Convert the first 8 bytes into a number.
	return int64(binary.BigEndian.Uint64(m.Sum(nil)))
}
