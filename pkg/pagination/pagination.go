package pagination

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/wangyaodream/gerty-goblog/pkg/config"
	"github.com/wangyaodream/gerty-goblog/pkg/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Page struct {
	URL    string
	Number int
}

type ViewData struct {
	HasPages   bool
	Next       Page
	HasNext    bool
	Prev       Page
	HasPrev    bool
	Current    Page
	TotalCount int64
	TotalPage  int
}

type Pagination struct {
	BaseURL string
	PerPage int
	Page    int
	Count   int64
	db      *gorm.DB
}

func New(r *http.Request, db *gorm.DB, baseURL string, PerPage int) *Pagination {
	if PerPage <= 0 {
		PerPage = config.GetInt("pagination.perpage")
	}

	p := &Pagination{
		db:      db,
		PerPage: PerPage,
		Page:    1,
		Count:   -1,
	}

	if strings.Contains(baseURL, "?") {
		p.BaseURL = baseURL + "&" + config.GetString("pagination.url_query") + "="
	} else {
		p.BaseURL = baseURL + "?" + config.GetString("pagination.url_query") + "="
	}

	// 设当前页
	p.SetPage(p.GetPageFromRequest(r))

	return p
}

func (p *Pagination) Paging() ViewData {
	return ViewData{
		HasPages: p.HasPages(),
		Next:     p.NewPage(p.NextPage()),
		HasNext:  p.HasNext(),

		Prev:    p.NewPage(p.PrevPage()),
		HasPrev: p.HasPrev(),

		Current:   p.NewPage(p.CurrentPage()),
		TotalPage: p.TotalPage(),

		TotalCount: p.Count,
	}
}

func (p Pagination) NewPage(page int) Page {
	return Page{
		Number: page,
		URL:    p.BaseURL + strconv.Itoa(page),
	}
}

func (p *Pagination) SetPage(page int) {
	if page <= 0 {
		page = 1
	}

	p.Page = page
}

func (p Pagination) CurrentPage() int {
	totalPage := p.TotalPage()
	if totalPage == 0 {
		return 0
	}

	if p.Page > totalPage {
		return totalPage
	}

	return p.Page
}

func (p Pagination) Results(data interface{}) error {
	var err error
	var offset int
	page := p.CurrentPage()
	if page == 0 {
		return err
	}

	if page > 1 {
		offset = (page - 1) * p.PerPage
	}

	return p.db.Preload(clause.Associations).Limit(p.PerPage).Offset(offset).Find(data).Error
}

func (p *Pagination) TotalCount() int64 {
	if p.Count == -1 {
		var count int64
		if err := p.db.Count(&count).Error; err != nil {
			return 0
		}
		p.Count = count
	}
	return p.Count
}

func (p *Pagination) HasPages() bool {
	n := p.TotalCount()
	return n > int64(p.PerPage)
}

func (p Pagination) HasNext() bool {
	totalPage := p.TotalPage()
	if totalPage == 0 {
		return false
	}

	page := p.CurrentPage()
	if page == 0 {
		return false
	}
	return page < totalPage
}

func (p Pagination) PrevPage() int {
	hasPrev := p.HasPrev()

	if !hasPrev {
		return 0
	}

	page := p.CurrentPage()
	if page == 0 {
		return 0
	}

	return page - 1
}

func (p Pagination) NextPage() int {
	hasNext := p.HasNext()
	if !hasNext {
		return 0
	}

	page := p.CurrentPage()

	if page == 0 {
		return 0
	}

	return page + 1
}

func (p Pagination) HasPrev() bool {
	page := p.CurrentPage()
	if page == 0 {
		return false
	}

	return page > 1
}

func (p Pagination) TotalPage() int {
	count := p.TotalCount()
	if count == 0 {
		return 0
	}

	nums := int64(math.Ceil(float64(count) / float64(p.PerPage)))
	if nums == 0 {
		nums = 1
	}

	return int(nums)
}

func (p Pagination) GetPageFromRequest(r *http.Request) int {
	page := r.URL.Query().Get(config.GetString("pagination.url_query"))
	if len(page) > 0 {
		pageInt := types.StringToInt(page)
		if pageInt <= 0 {
			return 1
		}
		return pageInt
	}
	return 0
}
