package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "io"
    "os"
)

// generateKey creates a 32-byte key for AES-256 encryption
func generateKey(password string) []byte {
    key := make([]byte, 32)
    copy(key, []byte(password))
    return key
}

// EncryptPassword encrypts a password using AES-256-GCM
func EncryptPassword(password string, secretKey []byte) (string, error) {
    block, err := aes.NewCipher(secretKey)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }

    ciphertext := gcm.Seal(nonce, nonce, []byte(password), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptPassword decrypts a password using AES-256-GCM
func DecryptPassword(encryptedPassword string, secretKey []byte) (string, error) {
    ciphertext, err := base64.StdEncoding.DecodeString(encryptedPassword)
    if err != nil {
        return "", err
    }

    block, err := aes.NewCipher(secretKey)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return "", errors.New("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}

// GetSecretKey retrieves or generates a secret key from environment variable
func GetSecretKey() []byte {
    secretKeyStr := os.Getenv("ENCRYPTION_SECRET")
    if secretKeyStr == "" {
        // Generate a default key (in production, use a more secure method)
        secretKeyStr = "your-very-secret-key-that-should-be-replaced"
    }
    return generateKey(secretKeyStr)
}