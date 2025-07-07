package input

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/gabrielluizsf/file-encryptor/encryptor"
	"github.com/gabrielluizsf/file-encryptor/validator"
	"github.com/i9si-sistemas/stringx"
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

// ReadPSWDFn defines the function type for reading a password.
type ReadPSWDFn func(fd int) ([]byte, error)

var (
	// ErrInvalidOperation is returned when the user enters an invalid operation.
	ErrInvalidOperation = errors.New("invalid operation. Please choose 'encrypt' or 'decrypt'")

	// StdReader is the standard input reader.
	StdReader = bufio.NewReader(os.Stdin)

	// ReadPassword is the function for reading a password.
	ReadPassword ReadPSWDFn = PasswordReader()
)

// InputReader defines the interface for reading user input.
type InputReader interface {
	// ReadString reads a string from the input.
	ReadString(sep byte) (string, error)
}

// User prompts the user to enter the operation, input file path, output file path, and secret key.
func User(r InputReader, readPSWD ReadPSWDFn) (UserInput, error) {
	emptyInput := UserInput{}
	trimSpace := func(s string) stringx.String {
		return stringx.String(s).Trim(stringx.Space.String())
	}
	getInput := func() string {
		s, _ := r.ReadString('\n')
		return s
	}

	fmt.Print("Enter the operation (encrypt/decrypt): ")
	operation := getInput()
	operation = trimSpace(operation).String()

	fmt.Print("Enter the input file path: ")
	inputPath := getInput()
	inputPath = trimSpace(inputPath).String()

	fmt.Print("Enter the output file path: ")
	outputPath := getInput()
	outputPath = trimSpace(outputPath).String()

	// Secret key input with hidden characters
	fmt.Print("Enter the secret key (at least 11 characters): ")
	secret, err := readPSWD(int(os.Stdin.Fd()))
	if err != nil {
		return emptyInput, err
	}
	secretStr := string(secret)
	secretStr = trimSpace(secretStr).String()

	selectedOperation, err := validateOperation(operation)
	if err != nil {
		return emptyInput, err
	}

	input := UserInput{
		Operation:  selectedOperation,
		Path:       inputPath,
		OutputPath: outputPath,
		Secret:     secretStr,
	}
	return input, validateSecret(input.Secret)
}

// ErrSecretTooShort is returned when the secret key is less than 11 characters long.
var ErrSecretTooShort = encryptor.ErrKeyTooShort

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
		return cryptoOperation(stringx.Empty), ErrInvalidOperation
	}
}
