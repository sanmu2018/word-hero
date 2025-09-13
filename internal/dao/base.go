package dao

import (
	"fmt"

	"gorm.io/gorm"
)

// Base 结构体提供基础数据库操作
type Base struct{}

// BaseList 分页查询参数结构体
type BaseList struct {
	PageNum  int    `form:"pageNum" json:"pageNum"`   // 页码，从1开始（可选）
	PageSize int    `form:"pageSize" json:"pageSize"` // 每页大小（可选）
	Sort     string `form:"sort" json:"sort"`         // 排序字段
}

// PaginationResult 分页结果结构体
type PaginationResult struct {
	List       interface{} `json:"list"`       // 数据列表
	Total      int64       `json:"total"`      // 总记录数
	PageNum    int         `json:"pageNum"`    // 当前页码
	PageSize   int         `json:"pageSize"`   // 每页大小
	TotalPages int         `json:"totalPages"` // 总页数
}

// ApplyPagination 应用分页和排序到数据库查询
func ApplyPagination(db *gorm.DB, baseList *BaseList) (*gorm.DB, error) {
	if baseList != nil {
		// 处理排序
		if baseList.Sort != "" {
			sort, err := NormalizeSorts(baseList.Sort)
			if err != nil {
				return nil, fmt.Errorf("invalid sort fields: %w", err)
			}
			if sort == "" {
				sort = "created_at desc"
			}
			if len(sort) > 0 {
				db = db.Order(sort)
			}
		}

		// 处理分页 - 只有当pageNum和pageSize都设置且大于0时才应用分页
		if baseList.PageNum > 0 && baseList.PageSize > 0 {
			offset := (baseList.PageNum - 1) * baseList.PageSize
			db = db.Offset(offset).Limit(baseList.PageSize)
		}
	}
	return db, nil
}

// PageList 保持向后兼容的分页函数
func PageList(db *gorm.DB, baseList *BaseList) (*gorm.DB, error) {
	return ApplyPagination(db, baseList)
}

// NormalizeSorts 规范化排序字段
func NormalizeSorts(str string) (string, error) {
	var sorts string
	for _, v := range str {
		if v >= 'A' && v <= 'Z' {
			sorts += "_" + string(v+32)
		} else if v == '|' {
			sorts += " "
		} else if v >= 'a' && v <= 'z' || v == '_' || v == '-' || v == ',' {
			sorts += string(v)
		} else {
			return "", fmt.Errorf("invalid sort fields")
		}
	}
	return sorts, nil
}

// CalculatePagination 计算分页信息
func CalculatePagination(total int64, pageNum, pageSize int) PaginationResult {
	totalPages := 0
	if pageSize > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	return PaginationResult{
		Total:      total,
		PageNum:    pageNum,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
