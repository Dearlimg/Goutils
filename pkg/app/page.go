package app

import (
	"Goutils/pkg/utils"
	"net/http"
)

//分页处理

type Page struct {
	DefaultPageSize int32
	MaxPageSize     int32
	PageKey         string //URL 中 page 关键字
	PageSizeKey     string //URL 中 pagesize 关键字
}

// InitPage 初始化默认页数大小和最大页数限制以及查询的关键字
func InitPage(defaultPageSize, maxPageSize int32, pageKey, pageSizeKey string) *Page {
	return &Page{
		DefaultPageSize: defaultPageSize,
		MaxPageSize:     maxPageSize,
		PageKey:         pageKey,
		PageSizeKey:     pageSizeKey,
	}
}

// GetPageSizeAndOffset 从请求中获取页尺寸和偏移值
func (p *Page) GetPageSizeAndOffset(r *http.Request) (pageSize, offset int32) {
	page := utils.StrTo(r.FormValue(p.PageKey)).MustInt32()
	if page <= 0 {
		page = 1
	}
	pageSize = utils.StrTo(r.FormValue(p.PageSizeKey)).MustInt32()
	if pageSize <= 0 {
		pageSize = p.DefaultPageSize
	}
	if pageSize > p.MaxPageSize {
		pageSize = p.MaxPageSize
	}
	offset = (page - 1) * pageSize
	return
}

// CulOffset 计算偏移量
func (p *Page) CulOffset(page, pageSize int32) (offset int32) {
	return (page - 1) * pageSize
}
