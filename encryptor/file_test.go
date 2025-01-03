package encryptor

import (
	"github.com/i9si-sistemas/assert"
	"os"
	"testing"
)

// TestFileCrypto_NewFileCrypto tests the NewFileCrypto function.
func TestFileCrypto_NewFileCrypto(t *testing.T) {
	tests := []struct {
		secret   string
		expected error
	}{
		{"validsecret", nil},      // Valid key
		{"short", ErrKeyTooShort}, // Invalid key (less than 11 characters)
	}

	for _, test := range tests {
		t.Run(test.secret, func(t *testing.T) {
			_, err := NewFileCrypto(test.secret)
			assert.Equal(t, err, test.expected)
		})
	}
}

// TestFileCrypto_EncryptDecrypt tests encryption and decryption of files.
func TestFileCrypto_EncryptDecrypt(t *testing.T) {
	secret := "validsecret"
	fc, err := NewFileCrypto(secret)
	assert.NoError(t, err)

	// Create temporary directory
	tempDir := t.TempDir()
	inputFilePath := tempDir + "/input.txt"
	outputFilePath := tempDir + "/encrypted.bin"
	decryptedFilePath := tempDir + "/decrypted.txt"

	err = os.WriteFile(inputFilePath, []byte("This is a secret message."), 0644)
	assert.NoError(t, err)

	// Test encryption
	err = fc.Encrypt(inputFilePath, outputFilePath)
	assert.NoError(t, err)

	// Test decryption
	err = fc.Decrypt(outputFilePath, decryptedFilePath)
	assert.NoError(t, err)

	// Read and compare content
	originalContent, err := os.ReadFile(inputFilePath)
	assert.NoError(t, err)

	decryptedContent, err := os.ReadFile(decryptedFilePath)
	assert.NoError(t, err)

	assert.Equal(t, originalContent, decryptedContent)
}

// TestFileCrypto_EncryptDecrypt_InvalidKey tests the decryption with an invalid key.
func TestFileCrypto_EncryptDecrypt_InvalidKey(t *testing.T) {
	// Prepare encryption with a valid key
	secret := "validsecret"
	fc, err := NewFileCrypto(secret)
	assert.NoError(t, err)

	// Create temporary directory
	tempDir := t.TempDir()
	inputFilePath := tempDir + "/input.txt"
	outputFilePath := tempDir + "/encrypted.bin"
	fileBytes := []byte("This is a secret message.")
	err = os.WriteFile(inputFilePath, fileBytes, 0644)
	assert.NoError(t, err)

	err = fc.Encrypt(inputFilePath, outputFilePath)
	assert.NoError(t, err)

	// Prepare decryption with an invalid key
	invalidSecret := "invalidsecret"
	invalidFC, err := NewFileCrypto(invalidSecret)
	assert.NoError(t, err)

	// Try to decrypt with the invalid key
	decryptedFilePath := tempDir + "/decrypted_invalid.txt"
	err = invalidFC.Decrypt(outputFilePath, decryptedFilePath)
	assert.NoError(t, err)
	b, err := os.ReadFile(decryptedFilePath)
	assert.NoError(t, err)
	assert.NotEqual(t, fileBytes, b)
}

// TestFileCrypto_Encrypt_FileNotFound tests encryption with non-existent input file.
func TestFileCrypto_Encrypt_FileNotFound(t *testing.T) {
	secret := "validsecret"
	fc, err := NewFileCrypto(secret)
	assert.NoError(t, err)

	// Create temporary directory
	tempDir := t.TempDir()
	inputFilePath := tempDir + "/nonexistent.txt"
	outputFilePath := tempDir + "/encrypted.bin"

	// Try encrypting a non-existent file
	err = fc.Encrypt(inputFilePath, outputFilePath)
	assert.Error(t, err)
}

// TestFileCrypto_Decrypt_FileNotFound tests decryption with non-existent input file.
func TestFileCrypto_Decrypt_FileNotFound(t *testing.T) {
	secret := "validsecret"
	fc, err := NewFileCrypto(secret)
	assert.NoError(t, err)

	// Create temporary directory
	tempDir := t.TempDir()
	inputFilePath := tempDir + "/nonexistent_encrypted.bin"
	outputFilePath := tempDir + "/decrypted.txt"

	// Try decrypting a non-existent file
	err = fc.Decrypt(inputFilePath, outputFilePath)
	assert.Error(t, err)
}
