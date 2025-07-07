# File Encryption and Decryption in Go

```go
package main

import (
	"fmt"
	"github.com/gabrielluizsf/file-encryptor/encryptor"
	"github.com/gabrielluizsf/file-encryptor/input"
)

func main() {
	// Get user input
	userInput, err := input.User(input.StdReader(), input.ReadPassword)
	if err != nil {
		fmt.Println(err)
		return
	}
	var (
		operation, secret     = userInput.Operation, userInput.Secret
		inputPath, outputPath = userInput.Path, userInput.OutputPath
	)
	// Create a new FileCrypto 'instance'
	fileCrypto, err := encryptor.NewFileCrypto(secret)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Perform the operation based on user input
	switch operation {
	case input.Encrypt:
		if err := fileCrypto.Encrypt(inputPath, outputPath); err != nil {
			fmt.Println("Error encrypting file:", err)
			return
		}
		fmt.Println("\nFile encrypted successfully")
	case input.Decrypt:
		if err := fileCrypto.Decrypt(inputPath, outputPath); err != nil {
			fmt.Println("Error decrypting file:", err)
			return
		}
		fmt.Println("\nFile decrypted successfully")
	}
}

```
