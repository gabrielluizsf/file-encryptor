package input

import (
	"bufio"
	"errors"
	"os"
	"testing"

	"github.com/i9si-sistemas/assert"
	"github.com/i9si-sistemas/stringx"
	"golang.org/x/term"
)

func TestStdReader(t *testing.T) {
	assert.Equal(t, StdReader()(), bufio.NewReader(os.Stdin))
}

func TestReadPassword(t *testing.T) {
	assert.Equal(t, ReadPassword, term.ReadPassword)
}

func TestOperations(t *testing.T) {
	encryptOp, decryptOp := "encrypt", "decrypt"
	assert.Equal(t, Encrypt, encryptOp)
	assert.Equal(t, Encrypt.String(), encryptOp)
	assert.Equal(t, Decrypt, decryptOp)
	assert.Equal(t, Decrypt.String(), decryptOp)
}

// the max of ReadString method calls is 3
var errInvalidCalledCount = errors.New("invalid called count")

func TestUser(t *testing.T) {
	tests := []struct {
		operation      string
		inputFilePath  string
		outputFilePath string
		secret         string
		expectedErr    error
	}{
		{
			operation:      "encrypt",
			inputFilePath:  "input.txt",
			outputFilePath: "output.bin",
			secret:         "secret",
			expectedErr:    ErrSecretTooShort,
		},
		{
			operation:      "decrypt",
			inputFilePath:  "encrypted.bin",
			outputFilePath: "decrypted.txt",
			secret:         "12345678912",
			expectedErr:    nil,
		},
		{
			operation:      "invalid",
			inputFilePath:  "invalid.txt",
			outputFilePath: "invalid.bin",
			secret:         "12345678912",
			expectedErr:    ErrInvalidOperation,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.operation, func(t *testing.T) {
			var fd int
			r := testStringReader(testCase.operation, testCase.inputFilePath, testCase.outputFilePath)
			input, err := User(
				func() InputReader {
					return r
				},
				testReadPSWD(testCase.secret, &fd),
			)
			assert.Equal(t, err, testCase.expectedErr)
			assert.Equal(t, r.byteSeparator, '\n')
			assert.Equal(t, fd, os.Stdin.Fd())
			if testCase.operation == "invalid"  {
				assert.Zero(t, input.Operation)
				return
			}
			assert.Equal(t, input.Operation, testCase.operation)
		})
	}
}

func testReadPSWD(pswd string, fdPtr *int) ReadPSWDFn {
	return func(fd int) ([]byte, error) {
		*fdPtr = fd
		return []byte(pswd), nil
	}
}

func testStringReader(
	op string,
	ifp, ofp string,
) *testTermStringReader {
	return &testTermStringReader{
		operation:      op,
		inputFilePath:  ifp,
		outputFilePath: ofp,
	}
}

type testTermStringReader struct {
	byteSeparator  byte
	operation      string
	inputFilePath  string
	outputFilePath string
	calledCount    int
}

func (r *testTermStringReader) ReadString(sep byte) (op string, err error) {
	r.byteSeparator = sep
	r.calledCount++
	switch r.calledCount {
	case 1:
		return r.operation, nil
	case 2:
		return r.inputFilePath, nil
	case 3:
		return r.outputFilePath, nil
	}
	return stringx.Empty.String(), errInvalidCalledCount
}
