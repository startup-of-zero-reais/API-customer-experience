package validation

import "fmt"

func RequiredStringField(field, label string) error {
	if field == "" {
		return fmt.Errorf("o campo %s é obrigatório", label)
	}

	return nil
}

func RequiredIntField(field, label int) error {
	if field == 0 {
		return fmt.Errorf("o campo %d é obrigatório", label)
	}

	return nil
}
