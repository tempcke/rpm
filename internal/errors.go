package internal

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotImplemented = knownErr("not implemented")
	ErrBadRequest     = knownErr("bad request")
	ErrInternal       = knownErr("internal error")
	ErrEntityInvalid  = knownErr("entity state invalid")
	ErrEntityNotFound = knownErr("entity not found")
)

type Errors []error

func MakeErr(errType error, msg string) error {
	return fmt.Errorf("%w: %s", errType, msg)
}
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
func (e Errors) ErrorOrNil() error {
	if len(e) > 0 {
		return e
	}
	return nil
}

type knownErr string

func IsKnownErr(err error) bool {
	var kErr knownErr
	return errors.As(err, &kErr)
}
func (e knownErr) Error() string {
	return string(e)
}
