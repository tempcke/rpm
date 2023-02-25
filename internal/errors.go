package internal

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInternal       = errors.New("internal error")
	ErrEntityInvalid  = errors.New("entity state invalid")
	ErrEntityNotFound = errors.New("entity not found")
)

func MakeErr(errType error, msg string) error {
	return fmt.Errorf("%w: %s", errType, msg)
}

type Errors []error

func NewErrors(errorList ...error) Errors {
	return append(Errors{}, errorList...)
}
func (e Errors) Append(errorList ...error) Errors {
	return append(append(Errors{}, e...), errorList...)
}
func (e Errors) Error() string {
	var msgs = make([]string, len(e))
	j := len(e) - 1
	for i, err := range e {
		msgs[j-i] = err.Error()
	}
	return strings.Join(msgs, ": ")
}

func (e Errors) Is(target error) bool {
	for _, err := range e {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}
