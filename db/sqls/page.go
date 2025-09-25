// Package sqls -----------------------------
//
//	@file		: page.go
//	@author		: Carlos
//	@contact	: 534994749@qq.com
//	@time		: 2025/6/10 14:20
//
// -------------------------------------------
package sqls

type Paging struct {
	Page  int `form:"page"`  // 页码
	Limit int `form:"limit"` // 每页条数
}

func (p *Paging) Offset() int {
	offset := 0
	if p.Page > 0 {
		offset = (p.Page - 1) * p.Limit
	}
	return offset
}
