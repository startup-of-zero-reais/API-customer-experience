package domain

import "fmt"

func requiredStringField(field string) error {
	if field == "" {
		return fmt.Errorf("%s is required", field)
	}

	return nil
}

func requiredIntField(field int) error {
	if field == 0 {
		return fmt.Errorf("%s is required", field)
	}

	return nil
}
