package input

import "golang.org/x/term"

func PasswordReader() ReadPSWDFn {
	return term.ReadPassword
}
