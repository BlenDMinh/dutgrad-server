package repositories

import (
	"fmt"

	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/databases/entities"
)

var DefaultPageSize = 20
var DefaultPreload = []string{}

type ICrudRepository[T entities.Entity, ID any] interface {
	GetAll(page int, pageSize int) ([]T, error)
	GetById(id ID) (*T, error)
	Create(model *T) (*T, error)
	Update(model *T) (*T, error)
	Delete(id ID) error
	GetByField(fieldName string, value interface{}) ([]T, error)
	DeleteByField(fieldName string, value interface{}) error
	Count() (int64, error)
	GetByFieldWithPagination(fieldName string, value interface{}, page int, pageSize int) ([]T, error)
}

type CrudRepository[T entities.Entity, ID any] struct {
	DefaultPageSize int
	Preload         []string
}

func NewCrudRepository[T entities.Entity, ID any]() *CrudRepository[T, ID] {
	return &CrudRepository[T, ID]{
		DefaultPageSize: DefaultPageSize,
		Preload:         DefaultPreload,
	}
}

func (c *CrudRepository[T, ID]) getModel() *T {
	return new(T)
}

func (c *CrudRepository[T, ID]) GetPaginated(page int, pageSize int) ([]T, Pagination, error) {
	pagination := NewPagination(page, pageSize, c.DefaultPageSize)

	db := databases.GetDB()
	entities := []T{}

	dbctx := pagination.ApplyPagination(db)
	if len(c.Preload) > 0 {
		for _, preload := range c.Preload {
			dbctx = dbctx.Preload(preload)
		}
	}
	dbctx = dbctx.Find(&entities)
	if dbctx.Error != nil {
		return nil, pagination, dbctx.Error
	}

	var count int64
	if err := db.Model(new(T)).Count(&count).Error; err != nil {
		return nil, pagination, err
	}

	pagination.Total = count

	return entities, pagination, nil
}

func (c *CrudRepository[T, ID]) GetAll(page int, pageSize int) ([]T, error) {
	entities, _, err := c.GetPaginated(page, pageSize)
	return entities, err
}

func (c *CrudRepository[T, ID]) GetById(id ID) (*T, error) {
	db := databases.GetDB()
	dbctx := db.Model(c.getModel())
	entity := new(T)
	if len(c.Preload) > 0 {
		for _, preload := range c.Preload {
			dbctx = dbctx.Preload(preload)
		}
	}
	dbctx = dbctx.First(entity, id)
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

func (c *CrudRepository[T, ID]) GetByField(fieldName string, value interface{}) ([]T, error) {
	var results []T
	db := databases.GetDB()
	dbctx := db.Model(c.getModel())
	if len(c.Preload) > 0 {
		for _, preload := range c.Preload {
			dbctx = dbctx.Preload(preload)
		}
	}
	err := dbctx.Where(fmt.Sprintf("%s = ?", fieldName), value).Find(&results).Error
	return results, err
}

func (c *CrudRepository[T, ID]) DeleteByField(fieldName string, value interface{}) error {
	db := databases.GetDB()
	return db.Where(fmt.Sprintf("%s = ?", fieldName), value).Delete(new(T)).Error
}

func (c *CrudRepository[T, ID]) Count() (int64, error) {
	db := databases.GetDB()
	var count int64
	err := db.Model(new(T)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *CrudRepository[T, ID]) GetByFieldPaginated(fieldName string, value interface{}, page int, pageSize int) ([]T, Pagination, error) {
	pagination := NewPagination(page, pageSize, c.DefaultPageSize)

	db := databases.GetDB()
	entities := []T{}

	dbctx := db.Model(c.getModel())

	if len(c.Preload) > 0 {
		for _, preload := range c.Preload {
			dbctx = dbctx.Preload(preload)
		}
	}

	whereClause := fmt.Sprintf("%s = ?", fieldName)
	dbctx = pagination.ApplyPagination(dbctx)
	dbctx = dbctx.Where(whereClause, value).Find(&entities)

	if dbctx.Error != nil {
		return nil, pagination, dbctx.Error
	}

	var count int64
	if err := db.Model(new(T)).Where(whereClause, value).Count(&count).Error; err != nil {
		return nil, pagination, err
	}

	pagination.Total = count

	return entities, pagination, nil
}

func (c *CrudRepository[T, ID]) GetByFieldWithPagination(fieldName string, value interface{}, page int, pageSize int) ([]T, error) {
	entities, _, err := c.GetByFieldPaginated(fieldName, value, page, pageSize)
	return entities, err
}
