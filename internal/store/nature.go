package store

import (
	"ProtectedArea/internal/model"
	"gorm.io/gorm"
)

// NatureStore 定义接口，方便后续扩展
type NatureStore interface {
	GetYearlyTrendStats() ([]model.StatResult, error)
}

// natureStore 结构体实现接口
type natureStore struct {
	db *gorm.DB
}

// NewNatureStore 构造函数
func NewNatureStore(db *gorm.DB) NatureStore {
	return &natureStore{db: db}
}

// GetYearlyTrendStats 执行具体的 SQL 统计查询
func (s *natureStore) GetYearlyTrendStats() ([]model.StatResult, error) {
	var results []model.StatResult

	// SQL: SELECT year, BHDL, count(*) FROM nature_data WHERE ... GROUP BY ...
	err := s.db.Model(&model.NatureData{}).
		Select("year, BHDL, count(*) as count").
		Where("BHDL IN ?", []string{"资源损毁", "恢复治理"}).
		Group("year, BHDL").
		Scan(&results).Error

	return results, err
}
