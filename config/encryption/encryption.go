package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"regexp"
)

var encryptionKey []byte

var encStringRegex *regexp.Regexp = regexp.MustCompile(`ENC\((?P<value>.*)\)`)

func SetEncryptionKey(key string) {
	hasher := md5.New()
	hasher.Write([]byte(key))
	encryptionKey = []byte(hex.EncodeToString(hasher.Sum(nil)))
	if len(encryptionKey) != 32 {
		panic("Encryption key is not 32bit!")
	}
}

func EncryptConfigString(input string) (string, error) {
	res, err := Encrypt([]byte(input))
	if err != nil {
		return "", err
	}
	return "ENC(" + string(res) + ")", nil
}

func DecryptConfigString(input string) (string, error) {
	if !IsEncrypted(input) {
		return input, nil
	}
	matches := encStringRegex.FindStringSubmatch(input)
	if len(matches) < 2 {
		return "", errors.New("something is very odd about the decription input")
	}
	res, err := Decrypt([]byte(matches[1]))
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func IsEncrypted(input string) bool {
	return encStringRegex.Match([]byte(input))
}

func Encrypt(input []byte) ([]byte, error) {
	cip, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(cip)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return []byte(base64.RawStdEncoding.EncodeToString(gcm.Seal(nonce, nonce, input, nil))), nil
}

func Decrypt(input []byte) ([]byte, error) {
	decodedInput, err := base64.RawStdEncoding.DecodeString(string(input))
	if err != nil {
		return nil, err
	}
	cip, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(cip)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(decodedInput) < nonceSize {
		return nil, errors.New("given input for decryption is too small to be valid")
	}
	nonce, decodedInput := decodedInput[:nonceSize], decodedInput[nonceSize:]
	ciphertext, err := gcm.Open(nil, nonce, decodedInput, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}
