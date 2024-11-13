package encryptutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
)

const EncryptionKey = "976e888b6a111a0f8097f4de42b3b133" // Make sure this is exactly 32 bytes

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
