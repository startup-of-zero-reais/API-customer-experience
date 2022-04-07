package validation

import (
	"encoding/json"
	"errors"
	"log"
)

type (
	// FieldValidator é responsável por validar campos
	FieldValidator struct {
		Errors []map[string]interface{} `json:"errors,omitempty"`
		Err    error                    `json:"error"`
	}
)

// NewFieldValidator cria um novo FieldValidator
func NewFieldValidator() *FieldValidator {
	return &FieldValidator{
		Err:    errors.New("erro de validação"),
		Errors: []map[string]interface{}{},
	}
}

// AddError adiciona um erro ao FieldValidator
func (f *FieldValidator) AddError(field string, message string) *FieldValidator {
	for i := range f.Errors {
		if f.Errors[i]["field"] == field {
			f.Errors[i]["errors"] = append(f.Errors[i]["errors"].([]string), message)
			return f
		}
	}

	fieldError := map[string]interface{}{
		"field":  field,
		"errors": []string{},
	}

	fieldError["errors"] = append(fieldError["errors"].([]string), message)
	f.Errors = append(f.Errors, fieldError)

	return f
}

// HasErrors verifica se o FieldValidator possui erros
func (f *FieldValidator) HasErrors() bool {
	return len(f.Errors) > 0
}

// Error implementa a interface error para que seja possível utilizar o FieldValidator como um erro
func (f *FieldValidator) Error() string {
	log.Println("Error:", f.Err)
	bytes, err := json.Marshal(f)
	if err != nil {
		return err.Error()
	}

	return string(bytes)
}

// MarshalJSON implementa a interface json.Marshaler para que seja possível marshalar o FieldValidator
func (f *FieldValidator) MarshalJSON() ([]byte, error) {
	mashalled := struct {
		Err    string      `json:"error"`
		Errors interface{} `json:"errors,omitempty"`
	}{
		Err:    f.Err.Error(),
		Errors: f.Errors,
	}

	return json.Marshal(mashalled)
}
