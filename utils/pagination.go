package pagination

const globalDefaultPerPage = 20

type Pagination struct {
	Page       int64 `query:"page" form:"page" json:"page" description:"目前頁面"`
	PerPage    int64 `query:"perPage" form:"perPage" json:"perPage" description:"每頁顯示多少筆"`
	TotalCount int64 `query:"totalCount" form:"totalCount" json:"totalCount" description:"總筆數"`
	TotalPage  int64 `query:"totalPage" form:"totalPage" json:"totalPage" description:"總頁數"`
}

func (p *Pagination) CheckOrSetDefault(params ...int64) *Pagination {
	var defaultPerPage int64
	if len(params) >= 1 {
		defaultPerPage = params[0]
	}

	if defaultPerPage <= 0 {
		defaultPerPage = globalDefaultPerPage
	}

	if p.Page == 0 {
		p.Page = 1
	}
	if p.PerPage == 0 {
		p.PerPage = defaultPerPage
	}
	return p
}

func (p *Pagination) LimitAndOffset() (uint64, uint64) {
	return uint64(p.PerPage), uint64(p.Offset())
}

func (p *Pagination) Offset() int64 {
	if p.Page <= 0 {
		return 0
	}
	return (p.Page - 1) * p.PerPage
}
