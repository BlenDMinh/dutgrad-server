package services

import (
	"fmt"
	"reflect"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
)

type ICrudService[T entities.Entity, ID any] interface {
	GetAll(page int, pageSize int) ([]T, error)
	GetById(id ID) (*T, error)
	Create(model *T) (*T, error)
	Update(id ID, model *T) (*T, error)
	Upsert(id ID, model *T) (*T, error)
	Delete(id ID) error
}

type CrudService[T entities.Entity, ID any] struct {
	repo repositories.ICrudRepository[T, ID]
}

func NewCrudService[T entities.Entity, ID any](repo repositories.ICrudRepository[T, ID]) *CrudService[T, ID] {
	return &CrudService[T, ID]{
		repo: repo,
	}
}

func (s *CrudService[T, ID]) GetAll(page int, pageSize int) ([]T, error) {
	return s.repo.GetAll(page, pageSize)
}

func (s *CrudService[T, ID]) GetById(id ID) (*T, error) {
	return s.repo.GetById(id)
}

func (s *CrudService[T, ID]) Create(model *T) (*T, error) {
	return s.repo.Create(model)
}

func (s *CrudService[T, ID]) Update(id ID, model *T) (*T, error) {
	existing, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("%s entity record with id %v not found", reflect.TypeOf(new(T)).Elem().Name(), id)
	}
	reflect.ValueOf(model).Elem().FieldByName("CreatedAt").Set(reflect.ValueOf(existing).Elem().FieldByName("CreatedAt"))
	reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(id))

	return s.repo.Update(model)
}

func (s *CrudService[T, ID]) Delete(id ID) error {
	return s.repo.Delete(id)
}

func (s *CrudService[T, ID]) Upsert(id ID, model *T) (*T, error) {
	existing, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return s.repo.Create(model)
	}

	return s.repo.Update(model)
}
