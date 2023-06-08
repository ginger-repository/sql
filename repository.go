package sql

import (
	"github.com/ginger-core/errors"
	"github.com/ginger-core/repository"
)

type Repository interface {
	repository.Transational
	Initialize() errors.Error
}
