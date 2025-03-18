package repositories

import (
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
)

var DefaultPageSize = 20

type ICrudRepository[T entities.Entity, ID any] interface {
	GetAll(page int, pageSize int) ([]T, error)
	GetById(id ID) (*T, error)
	Create(model *T) (*T, error)
	Update(model *T) (*T, error)
	Delete(id ID) error
}

type CrudRepository[T entities.Entity, ID any] struct {
	DefaultPageSize int
}

func NewCrudRepository[T entities.Entity, ID any]() *CrudRepository[T, ID] {
	return &CrudRepository[T, ID]{
		DefaultPageSize: DefaultPageSize,
	}
}

func (c *CrudRepository[T, ID]) getModel() *T {
	return new(T)
}

func (c *CrudRepository[T, ID]) GetAll(page int, pageSize int) ([]T, error) {

	if pageSize <= 0 {
		pageSize = c.DefaultPageSize
	}
	if page <= 0 {
		page = 1
	}

	db := databases.GetDB()
	entities := []T{}
	offset := (page - 1) * pageSize
	dbctx := db.Limit(pageSize).Offset(offset).Find(&entities)

	if dbctx.Error != nil {
		return nil, dbctx.Error
	}

	return entities, nil
}

func (c *CrudRepository[T, ID]) GetById(id ID) (*T, error) {
	db := databases.GetDB()
	entity := new(T)
	dbctx := db.First(entity, id)
	if dbctx.Error != nil {
		return nil, dbctx.Error
	}

	return entity, nil
}

func (c *CrudRepository[T, ID]) Create(model *T) (*T, error) {
	db := databases.GetDB()
	dbctx := db.Create(model)
	if dbctx.Error != nil {
		return nil, dbctx.Error
	}

	return model, nil
}

func (c *CrudRepository[T, ID]) Update(model *T) (*T, error) {
	db := databases.GetDB()
	dbctx := db.Save(model)
	if dbctx.Error != nil {
		return nil, dbctx.Error
	}

	return model, nil
}

func (c *CrudRepository[T, ID]) Delete(id ID) error {
	db := databases.GetDB()
	model := c.getModel()
	dbctx := db.Delete(model, id)
	if dbctx.Error != nil {
		return dbctx.Error
	}

	return nil
}
