package encryptutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"log"
)

func Encrypt(text string, key string) string {
    block, err := aes.NewCipher([]byte(key))
    if err != nil {
        log.Fatal(err)
    }
    plaintext := []byte(text)
    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        log.Fatal(err)
    }
    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
    return hex.EncodeToString(ciphertext)
}

func Decrypt(encryptedText string, key string) string {
    ciphertext, _ := hex.DecodeString(encryptedText)
    block, err := aes.NewCipher([]byte(key))
    if err != nil {
        log.Fatal(err)
    }
    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]
    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(ciphertext, ciphertext)
    return string(ciphertext)
}
func EncryptWithSalt(plainText, key string) (string, string, error) {
	// Generate a random IV (Initialization Vector)
	iv := make([]byte, aes.BlockSize)
	_, err := rand.Read(iv)
	if err != nil {
		return "", "", err
	}

	// Create AES cipher block
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", "", err
	}

	// Pad plaintext to be multiple of block size
	plainText = pad(plainText, aes.BlockSize)

	// Encrypt the data using AES and IV
	ciphertext := make([]byte, len(plainText))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, []byte(plainText))

	// Encode the ciphertext and iv to base64 strings for storage
	cipherTextBase64 := base64.StdEncoding.EncodeToString(ciphertext)
	ivBase64 := base64.StdEncoding.EncodeToString(iv)

	return cipherTextBase64, ivBase64, nil
}

func DecryptWithSalt(cipherTextBase64, ivBase64, key string) (string, error) {
	// Decode the ciphertext and IV from base64
	ciphertext, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", err
	}
	iv, err := base64.StdEncoding.DecodeString(ivBase64)
	if err != nil {
		return "", err
	}

	// Create AES cipher block
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Decrypt the data using AES and IV
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// Remove padding
	plainText := string(unpad(ciphertext, aes.BlockSize))
	return plainText, nil
}

// Padding functions (PKCS7)
func pad(data string, blockSize int) string {
	padding := blockSize - len(data)%blockSize
	padText := make([]byte, padding)
	for i := 0; i < padding; i++ {
		padText[i] = byte(padding)
	}
	return data + string(padText)
}

func unpad(data []byte, blockSize int) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}

// EncryptECB encrypts the given plaintext using AES in ECB mode
func EncryptECB(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Ensure plaintext is a multiple of AES block size by padding
	plaintextBytes := []byte(plaintext)
	padding := aes.BlockSize - len(plaintextBytes)%aes.BlockSize
	for i := 0; i < padding; i++ {
		plaintextBytes = append(plaintextBytes, byte(padding))
	}

	ciphertext := make([]byte, len(plaintextBytes))
	for start := 0; start < len(plaintextBytes); start += aes.BlockSize {
		block.Encrypt(ciphertext[start:start+aes.BlockSize], plaintextBytes[start:start+aes.BlockSize])
	}

	return hex.EncodeToString(ciphertext), nil
}

// DecryptECB decrypts the given hex-encoded ciphertext using AES in ECB mode
func DecryptECB(cipherHex string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", err
	}

	// Ensure ciphertext length is a multiple of the block size
	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("invalid ciphertext length")
	}

	plaintextBytes := make([]byte, len(ciphertext))
	for start := 0; start < len(ciphertext); start += aes.BlockSize {
		block.Decrypt(plaintextBytes[start:start+aes.BlockSize], ciphertext[start:start+aes.BlockSize])
	}

	// Remove padding
	padding := int(plaintextBytes[len(plaintextBytes)-1])
	plaintextBytes = plaintextBytes[:len(plaintextBytes)-padding]

	return string(plaintextBytes), nil
}