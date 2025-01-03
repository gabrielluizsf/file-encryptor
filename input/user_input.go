package input

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// UserInput defines the structure for user input.
type UserInput struct {
	Operation, Path,
	OutputPath, Secret string
}

// ErrInvalidOperation is returned when the user enters an invalid operation.
var ErrInvalidOperation = errors.New("invalid operation. Please choose 'encrypt' or 'decrypt'")

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
	secret, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return UserInput{}, err
	}
	secretStr := string(secret)
	secretStr = strings.TrimSpace(secretStr)

	if operation != "encrypt" && operation != "decrypt" {
		return UserInput{}, ErrInvalidOperation
	}

	input := UserInput{
		Operation:  operation,
		Path:       inputPath,
		OutputPath: outputPath,
		Secret:     secretStr,
	}
	return input, nil
}
