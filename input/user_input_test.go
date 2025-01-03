package input

import (
	"fmt"
	"os"
	"testing"
)

type MockReadPswd struct {
	Key string
}

// Mock for the terminal password read function
func (m MockReadPswd) ReadPassword(fd int) ([]byte, error) {
	return []byte(m.Key), nil // Returning a fixed password for testing
}

func TestUser(t *testing.T) {
	originalStdin := os.Stdin
	originalReadPswd := readPassword

	defer func() {
		os.Stdin = originalStdin
		readPassword = originalReadPswd
	}()
	executeInput("secretkey")
	userInput, err := User()
	if err != ErrSecretTooShort {
		t.Fatal(err)
	}

	fmt.Println("Test user input:", userInput)

	if userInput.Operation != "encrypt" {
		t.Errorf("expected 'encrypt', got %s", userInput.Operation)
	}
	if userInput.Path != "file.txt" {
		t.Errorf("expected 'file.txt', got %s", userInput.Path)
	}
	if userInput.OutputPath != "file.enc" {
		t.Errorf("expected 'file.enc', got %s", userInput.OutputPath)
	}
	if userInput.Secret != "secretkey" {
		t.Errorf("expected 'secretkey', got %s", userInput.Secret)
	}
	
	executeInput("secretkey11")

	userInput, err = User()
	if err != nil {
		t.Fatal(err)
	}

}

func executeInput(key string) {
	r, w, _ := os.Pipe()

	go func() {
		fmt.Fprintln(w, "encrypt")
		fmt.Fprintln(w, "file.txt")
		fmt.Fprintln(w, "file.enc")
		w.Close()
	}()

	os.Stdin = r
	readPassword = MockReadPswd{key}.ReadPassword
}
