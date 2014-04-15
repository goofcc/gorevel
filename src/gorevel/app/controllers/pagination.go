package controllers

import (
	"strconv"
)

const (
	PagesPerView = 11 //最多显示几个页码
	ItemsPerPage = 10 //每页几条记录
)

type Pagination struct {
	page      int //当前页码
	rows      int //记录总数
	url       string
	pageCount int //总页数
}

type PageNum struct {
	Num       int
	IsCurrent bool
	Url       string
}

func NewPagination(page, rows int, url string) *Pagination {
	return &Pagination{
		page: page,
		rows: rows,
		url:  url,
	}
}

func (p *Pagination) Pages() []PageNum {
	var result []PageNum

	p.pageCount = p.rows / ItemsPerPage
	if p.pageCount*ItemsPerPage < p.rows {
		p.pageCount += 1
	}
	if p.pageCount == 1 {
		return result
	}

	page := p.page
	page -= PagesPerView / 2
	if page < 0 {
		page = 0
	}

	count := page + PagesPerView
	if count > p.pageCount {
		count = p.pageCount
	}

	pageNum := 0
	for ; page < count; page++ {
		pageNum = page + 1
		result = append(result, PageNum{pageNum, page == p.page, p.url + strconv.Itoa(pageNum)})
	}

	return result
}
