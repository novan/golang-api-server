package schema

import (
	"fmt"
	"math"
)

// Pagination struct
type Pagination struct {
	CurrentPage int    `json:"current_page"`
	PageSize    int    `json:"page_size"`
	TotalPage   int    `json:"total_page"`
	TotalResult int    `json:"total_result"`
	Next        string `json:"next,omitempty"`
	Prev        string `json:"prev,omitempty"`
}

const DEFAULT_PAGESIZE = 25

func (p *Pagination) SetPageSize(size int) {
	p.PageSize = size
}

func (p *Pagination) GetPageSize() int {
	if p.PageSize == 0 {
		p.SetPageSize(DEFAULT_PAGESIZE)
	}
	return p.PageSize
}

func (p *Pagination) GetPageLimit() string {
	return fmt.Sprintf("LIMIT %d OFFSET %d", p.GetPageSize(), p.GetPageOffset())
}

func (p *Pagination) SetTotalPage(totalResult int) {
	totalAll := totalResult
	d := float64(totalAll) / float64(p.PageSize)
	totalPages := int(math.Ceil(d))

	if p.PageSize == 0 {
		p.PageSize = DEFAULT_PAGESIZE
	}

	p.TotalResult = totalResult
	p.TotalPage = totalPages
}

func (p *Pagination) GetTotalPage() int {
	return p.TotalPage
}

func (p *Pagination) SetCurrentPage(page int) {
	p.CurrentPage = page
}

func (p *Pagination) GetCurrentPage() int {
	if p.CurrentPage == 0 {
		p.SetCurrentPage(1)
	}
	return p.CurrentPage
}

func (p *Pagination) GetPageOffset() int {
	return (p.GetCurrentPage() - 1) * p.GetPageSize()
}

func (p *Pagination) GetFirstPage() int {
	if p.GetCurrentPage() > 1 {
		return (p.GetCurrentPage() * p.GetPageSize()) - p.GetPageSize()
	}
	return 0
}

func (p *Pagination) HasNextPage() bool {
	return (p.GetCurrentPage() * p.GetPageSize()) < p.TotalResult
}
