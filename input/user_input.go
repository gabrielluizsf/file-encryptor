package input

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gabrielluizsf/file-encryptor/validator"
	"golang.org/x/term"
)

// UserInput defines the structure for user input.
type UserInput struct {
	Operation cryptoOperation
	Path,
	OutputPath,
	Secret string
}

type cryptoOperation string

func (c cryptoOperation) String() string {
	return string(c)
}

const (
	Encrypt cryptoOperation = "encrypt"
	Decrypt cryptoOperation = "decrypt"
)

// ErrInvalidOperation is returned when the user enters an invalid operation.
var (
	ErrInvalidOperation = errors.New("invalid operation. Please choose 'encrypt' or 'decrypt'")
	readPassword        = term.ReadPassword
)

// User prompts the user to enter the operation, input file path, output file path, and secret key.
func User() (UserInput, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the operation (encrypt/decrypt): ")
	operation, _ := reader.ReadString('\n')
	operation = strings.TrimSpace(operation)

	fmt.Print("Enter the input file path: ")
	inputPath, _ := reader.ReadString('\n')
	inputPath = strings.TrimSpace(inputPath)

	fmt.Print("Enter the output file path: ")
	outputPath, _ := reader.ReadString('\n')
	outputPath = strings.TrimSpace(outputPath)

	// Secret key input with hidden characters
	fmt.Print("Enter the secret key (at least 11 characters): ")
	secret, err := readPassword(int(os.Stdin.Fd()))
	if err != nil {
		return UserInput{}, err
	}
	secretStr := string(secret)
	secretStr = strings.TrimSpace(secretStr)

	selectedOperation, err := validateOperation(operation)
	if err != nil {
		return UserInput{}, err
	}

	input := UserInput{
		Operation:  selectedOperation,
		Path:       inputPath,
		OutputPath: outputPath,
		Secret:     secretStr,
	}
	err = validateSecret(input.Secret)
	return input, err
}

// ErrSecretTooShort is returned when the secret key is less than 11 characters long.
var ErrSecretTooShort = errors.New("the secret key must be at least 11 characters long")

// validateKey ensures the key is at least 11 characters long.
func validateSecret(key string) error {
	return validator.Secret(ErrSecretTooShort).Validate([]byte(key))
}

func validateOperation(operation string) (cryptoOperation, error) {
	switch operation {
	case Encrypt.String():
		return Encrypt, nil
	case Decrypt.String():
		return Decrypt, nil
	default:
		return "", ErrInvalidOperation
	}
}
