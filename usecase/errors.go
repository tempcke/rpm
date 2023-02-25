package usecase

import (
	"errors"
)

var (
	ErrRepoNotSet = errors.New("use case repo is required")
)
