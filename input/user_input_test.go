package input

import (
	"fmt"
	"os"
	"testing"
)

// Mock for the terminal password read function
func mockReadPassword(fd int) ([]byte, error) {
	return []byte("secretkey"), nil // Returning a fixed password for testing
}

func TestUser(t *testing.T) {
	r, w, _ := os.Pipe()

	go func() {
		fmt.Fprintln(w, "encrypt")
		fmt.Fprintln(w, "file.txt")
		fmt.Fprintln(w, "file.enc")
		w.Close()
	}()

	originalStdin := os.Stdin
	originalReadPswd := readPassword

	defer func() {
		os.Stdin = originalStdin
		readPassword = originalReadPswd
	}()

	os.Stdin = r
	readPassword = mockReadPassword

	userInput, err := User()

	if err != nil {
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
}
