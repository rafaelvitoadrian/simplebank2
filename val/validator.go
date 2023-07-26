package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValdiateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("Harus Berkaratker dari %d-%d", minLength, maxLength)
	}
	return nil
}

func ValdiateUsername(value string) error {
	if err := ValdiateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(value) {
		return fmt.Errorf("Hanya Boleh terdiri dari kata, angka dan garis bawah")
	}
	return nil
}

func ValidateFullName(value string) error {
	if err := ValdiateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidFullName(value) {
		return fmt.Errorf("Hanya Boleh terdiri dari kata atau spasi")
	}
	return nil
}

func ValidatePassword(value string) error {
	return ValdiateString(value, 6, 100)
}

func ValidateEmail(value string) error {
	if err := ValdiateString(value, 3, 100); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("Bukan Email Yang Valid")
	}

	return nil
}

func ValidateEmailId(value int64) error {
	if value <= 0 {
		return fmt.Errorf("Harus bernilai lebih dari 0")
	}

	return nil
}

func ValdiateSecretCode(value string) error {
	return ValdiateString(value, 32, 128)
}
