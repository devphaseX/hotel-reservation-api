package utils

import "github.com/gofiber/fiber/v2"

const (
	page    = 1
	perPage = 10
)

type PaginateQuery struct {
	Page    *int64 `json:"page"`
	PerPage *int64 `json:"perPage"`
}

func (p *PaginateQuery) SetDefault() {
	if p.Page == nil || *p.Page < 1 {
		var _page int64 = page
		p.Page = &_page
	}

	if p.PerPage == nil || *p.PerPage < 1 {
		var _perPage int64 = perPage
		p.PerPage = &_perPage
	}
}

func (p *PaginateQuery) ApplyPagination() (offset, limit int64) {
	p.SetDefault()
	offset = (*p.Page - 1) * *p.PerPage
	limit = (*p.PerPage)
	return
}

func NewPaginate(fiber *fiber.Ctx) (*PaginateQuery, error) {
	p := PaginateQuery{}

	if err := fiber.QueryParser(&p); err != nil {
		return nil, err
	}

	p.SetDefault()

	return &p, nil
}

type ResourceResp struct {
	Data    any `json:"data"`
	Results int `json:"results"`
	Page    int `json:"page"`
}

func NewResourceResp(data any, res, page int) ResourceResp {
	return ResourceResp{
		Data:    data,
		Results: res,
		Page:    page,
	}
}
