package login

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"math/big"
	"strings"
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
func codeVerifierGenerator(codeVerifier *string) ([]byte, error) {
	maxLength := big.NewInt(85)
	n, err := rand.Int(rand.Reader, maxLength)
	if err != nil {
		return []byte{}, err
	}
	length := n.Int64() + 43
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return []byte{}, err
	}
	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}
	*codeVerifier = string(bytes)
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
