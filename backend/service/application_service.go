package service

import (
	"github.com/XR0118/1024Hackathon/backend/model"
)

type ApplicationService interface {
	Create(app *model.Application) error
	GetByID(id string) (*model.Application, error)
	List(repository, appType string) ([]*model.Application, error)
	Update(app *model.Application) error
	Delete(id string) error
	GetByRepository(repository string) ([]*model.Application, error)
}
