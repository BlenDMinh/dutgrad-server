package services

import (
	"fmt"
	"reflect"
	"time"

	"github.com/BlenDMinh/dutgrad-server/databases/entities"
	"github.com/BlenDMinh/dutgrad-server/databases/repositories"
)

type ICrudService[T entities.Entity, ID any] interface {
	GetAll(page int, pageSize int) ([]T, error)
	GetById(id ID) (*T, error)
	Create(model *T) (*T, error)
	Update(model *T) (*T, error)
	UpdateByID(id ID, model *T) (*T, error)
	PatchByID(id ID, patchData *T) (*T, error)
	Upsert(id ID, model *T) (*T, error)
	Delete(id ID) error
	GetByField(fieldName string, value interface{}) ([]T, error)
	DeleteByField(fieldName string, value interface{}) error
	Count() (int64, error)
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

func (s *CrudService[T, ID]) Update(model *T) (*T, error) {
	reflect.ValueOf(model).Elem().FieldByName("UpdatedAt").Set(reflect.ValueOf(time.Now()))
	return s.repo.Update(model)
}

func (s *CrudService[T, ID]) UpdateByID(id ID, model *T) (*T, error) {
	existing, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("%s entity record with id %v not found", reflect.TypeOf(new(T)).Elem().Name(), id)
	}
	reflect.ValueOf(model).Elem().FieldByName("UpdatedAt").Set(reflect.ValueOf(time.Now()))
	reflect.ValueOf(model).Elem().FieldByName("CreatedAt").Set(reflect.ValueOf(existing).Elem().FieldByName("CreatedAt"))
	reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(id))

	return s.repo.Update(model)
}

func (s *CrudService[T, ID]) PatchByID(id ID, patchData *T) (*T, error) {
	existing, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("%s entity record with id %v not found", reflect.TypeOf(new(T)).Elem().Name(), id)
	}

	updatedModel := *existing

	patchValue := reflect.ValueOf(patchData).Elem()
	existingValue := reflect.ValueOf(&updatedModel).Elem()
	modelType := patchValue.Type()

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		if field.Name == "ID" || field.Name == "CreatedAt" {
			continue
		}

		patchFieldValue := patchValue.Field(i)

		if field.Type.Kind() == reflect.Bool || !isZeroValue(patchFieldValue) {
			existingValue.Field(i).Set(patchFieldValue)
		}
	}

	reflect.ValueOf(patchData).Elem().FieldByName("UpdatedAt").Set(reflect.ValueOf(time.Now()))
	return s.repo.Update(&updatedModel)
}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Array, reflect.Slice, reflect.Map:
		return v.Len() == 0
	case reflect.Struct:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	default:
		return false
	}
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

func (s *CrudService[T, ID]) GetByField(fieldName string, value interface{}) ([]T, error) {
	return s.repo.GetByField(fieldName, value)
}

func (s *CrudService[T, ID]) DeleteByField(fieldName string, value interface{}) error {
	return s.repo.DeleteByField(fieldName, value)
}

func (s *CrudService[T, ID]) Count() (int64, error) {
	return s.repo.Count()
}
