package dto

type BaseList struct {
	PageNum  int    `form:"pageNum" json:"pageNum"`   // 页码，从1开始（可选）
	PageSize int    `form:"pageSize" json:"pageSize"` // 每页大小（可选）
	Sort     string `form:"sort" json:"sort"`         // 排序字段
}
