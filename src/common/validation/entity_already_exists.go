package validation

import (
	"errors"
)

type (
	// EntityAlreadyExists é responsável emitir erros de entidade já existente
	EntityAlreadyExists struct {
		Err error
	}
)

// EntityAlreadyExistsError cria um novo EntityAlreadyExists
func EntityAlreadyExistsError(err string) *EntityAlreadyExists {
	return &EntityAlreadyExists{
		Err: errors.New(err),
	}
}

// Error implementa a interface error para que seja possível utilizar o EntityAlreadyExists como um erro
func (e *EntityAlreadyExists) Error() string {
	return e.Err.Error()
}
