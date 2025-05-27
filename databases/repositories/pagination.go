package repositories

import "gorm.io/gorm"

type Pagination struct {
	Page     int
	PageSize int
	Total    int64
}

func NewPagination(page, pageSize int, defaultPageSize int) Pagination {
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if page <= 0 {
		page = 1
	}

	return Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p Pagination) HasNext() bool {
	return int64(p.Page*p.PageSize) < p.Total
}

func (p Pagination) HasPrev() bool {
	return p.Page > 1
}

func (p Pagination) TotalPages() int {
	if p.PageSize == 0 {
		return 0
	}

	totalPages := p.Total / int64(p.PageSize)
	if p.Total%int64(p.PageSize) != 0 {
		totalPages++
	}
	return int(totalPages)
}

func (p Pagination) ApplyPagination(tx *gorm.DB) *gorm.DB {
	return tx.Offset(p.Offset()).Limit(p.PageSize)
}
