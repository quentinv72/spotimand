package login

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"strings"
	"time"
)

const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_.~"

var codeVerifier string

// codeChallenger returns a base64url code challenger using SHA256 hash.
func codeChallenger() (string, error) {
	codeVerifierBytes, err := codeVerifierGenerator(&codeVerifier) // Possible test by mocking the codeVerifierGenerator
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(codeVerifierBytes)
	return base64URL(hash[:]), nil
}

// codeVerifierGenerator generates a code verifier of length between 43 and 127.
// The code verifier is made up of the characters in the chars constant.
// The function updates the codeVerifier variable.
func codeVerifierGenerator(s *string) ([]byte, error) {
	rand.Seed(time.Now().UnixNano())
	length := rand.Intn(85) + 43
	bytes := make([]byte, length)
	if _, err := crand.Read(bytes); err != nil {
		return []byte{}, err
	}
	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}
	*s = string(bytes)
	return bytes, nil
}

// Return base64url encoding of code challenger
func base64URL(hash []byte) string {
	encoding := base64.StdEncoding.EncodeToString(hash)
	encoding = strings.Replace(encoding, "+", "-", -1) // 62nd char of encoding
	encoding = strings.Replace(encoding, "/", "_", -1) // 63rd char of encoding
	encoding = strings.Replace(encoding, "=", "", -1)  // Remove any trailing '='s
	return encoding
}
