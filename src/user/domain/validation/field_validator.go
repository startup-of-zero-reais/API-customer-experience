package validation

import (
	"encoding/json"
	"errors"
	"log"
)

type (
	FieldValidator struct {
		Errors []map[string]interface{} `json:"errors,omitempty"`
		Err    error                    `json:"error"`
	}
)

func NewFieldValidator() *FieldValidator {
	return &FieldValidator{
		Err:    errors.New("erro de validação"),
		Errors: []map[string]interface{}{},
	}
}

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

func (f *FieldValidator) HasErrors() bool {
	return len(f.Errors) > 0
}

func (f *FieldValidator) Error() string {
	log.Println("Error:", f.Err)
	bytes, err := json.Marshal(f)
	if err != nil {
		return err.Error()
	}

	return string(bytes)
}

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
