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
	input, err := input.User()
	if err != nil {
		fmt.Println(err)
		return
	}
	var (
		operation, secret     = input.Operation, input.Secret
		inputPath, outputPath = input.Path, input.OutputPath
	)
	// Create a new FileCrypto 'instance'
	fileCrypto, err := encryptor.NewFileCrypto(secret)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Perform the operation based on user input
	switch operation {
	case "encrypt":
		if err := fileCrypto.Encrypt(inputPath, outputPath); err != nil {
			fmt.Println("Error encrypting file:", err)
			return
		}
		fmt.Println("\nFile encrypted successfully")
	case "decrypt":
		if err := fileCrypto.Decrypt(inputPath, outputPath); err != nil {
			fmt.Println("Error decrypting file:", err)
			return
		}
		fmt.Println("\nFile decrypted successfully")
	}
}


```
