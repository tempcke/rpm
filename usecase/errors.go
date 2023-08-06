package usecase

import (
	"errors"
)

var (
	ErrRepoNotSet = errors.New("use case repo is required")
	ErrRepo       = errors.New("error from repository")
)
