package validator

// Validator defines an interface for data validation.
// The Validate method checks the provided data and returns an error if it's invalid.
type Validator interface {
	// Validate validates the provided data.
	// Returns an error if the data is invalid, otherwise returns nil.
	Validate(data []byte) error
}

type SecretValidator struct {
	Err error
}

func (v SecretValidator) Validate(data []byte) error {
	if len(data) < 11 {
		return v.Err
	}
	return nil
}

type KeyValidator struct {
	SecretValidator
	Err error
}

func Secret(err error) Validator {
	return SecretValidator{
		Err: err,
	}
}

func Key(err error) Validator {
	return KeyValidator{
		Err: err,
	}
}
