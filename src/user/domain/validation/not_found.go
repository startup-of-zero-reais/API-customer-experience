package validation

import "errors"

type (
	// NotFound é responsável por emitir erros de entidade não encontrada
	NotFound struct {
		Err error
	}
)

// NewNotFound cria um novo NotFound
func NotFoundError(err string) *NotFound {
	return &NotFound{
		Err: errors.New(err),
	}
}

// Error implementa a interface error para que seja possível utilizar o NotFound como um erro
func (n *NotFound) Error() string {
	return n.Err.Error()
}
