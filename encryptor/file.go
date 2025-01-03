package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"os"
)

// FileCrypto defines the structure for file encryption and decryption.
type FileCrypto struct {
	key []byte
}

// NewFileCrypto creates a new 'instance' of FileCrypto from a secret key.
func NewFileCrypto(secret string) (*FileCrypto, error) {
	if err := validateKey(secret); err != nil {
		return nil, err
	}
	key := deriveKey(secret)
	return &FileCrypto{key: key}, nil
}

// Encrypt encrypts the input file and writes to the output file.
func (fc *FileCrypto) Encrypt(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	block, err := aes.NewCipher(fc.key)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)                       // AES block size is 16 bytes
	if _, err := io.ReadFull(rand.Reader, iv); err != nil { // Generate random IV
		return err
	}

	// Write the IV to the output file
	if _, err := outputFile.Write(iv); err != nil {
		return err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	writer := &cipher.StreamWriter{S: stream, W: outputFile}

	// Encrypt the file content
	if _, err := io.Copy(writer, inputFile); err != nil {
		return err
	}

	return nil
}

// Decrypt decrypts the input file and writes to the output file.
func (fc *FileCrypto) Decrypt(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// Read the IV from the beginning of the file
	iv := make([]byte, aes.BlockSize) // AES block size is 16 bytes
	if _, err := io.ReadFull(inputFile, iv); err != nil {
		return err
	}

	block, err := aes.NewCipher(fc.key)
	if err != nil {
		return err
	}

	// Create the decrypter with the IV
	stream := cipher.NewCFBDecrypter(block, iv)
	reader := &cipher.StreamReader{S: stream, R: inputFile} // Read the rest of the encrypted file

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Copy the decrypted content to the output file
	if _, err := io.Copy(outputFile, reader); err != nil {
		return err
	}

	return nil
}

// ErrKeyTooShort is returned when the key is less than 11 characters long.
var ErrKeyTooShort = errors.New("the key must be at least 11 characters long")

// validateKey ensures the key is at least 11 characters long.
func validateKey(key string) error {
	if len(key) < 11 {
		return ErrKeyTooShort
	}
	return nil
}

// deriveKey generates a 32-byte key from the given secret.
func deriveKey(secret string) []byte {
	hash := sha256.Sum256([]byte(secret))
	return hash[:]
}
