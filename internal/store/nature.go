package store

import (
	"ProtectedArea/internal/model"
	"gorm.io/gorm"
)

// NatureStore 定义接口，方便后续扩展
type NatureStore interface {
	GetYearlyTrendStats() ([]model.StatResult, error)

	GetSummaryByYear(year string) (int64, float64, error)
	GetDamageStatsByBatch(year string) ([]model.BatchStatResult, error)
	// GetRegionStats
	// year: 年份
	// groupCol: 要分组统计的目标列 (比如 THSHI)
	// filterCol: 筛选条件的列名 (比如 THSHENG)，如果没有筛选则为空字符串
	// filterVal: 筛选条件的值 (比如 "河北省")
	GetRegionStats(year string, groupCol string, filterCol string, filterVal string) ([]model.RegionStatResult, error)
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

// GetSummaryByYear 1. 获取某年的总图斑数和总面积
func (s *natureStore) GetSummaryByYear(year string) (int64, float64, error) {
	var result struct {
		TotalCount int64
		TotalArea  float64
	}

	// SQL: SELECT count(*) as total_count, sum(BHMJ) as total_area FROM nature_data WHERE year = ?
	err := s.db.Model(&model.NatureData{}).
		Select("count(*) as total_count, sum(BHMJ) as total_area").
		Where("year = ?", year).
		Scan(&result).Error

	return result.TotalCount, result.TotalArea, err
}

// GetDamageStatsByBatch 2. 获取某年“资源损毁”的分批次统计
func (s *natureStore) GetDamageStatsByBatch(year string) ([]model.BatchStatResult, error) {
	var results []model.BatchStatResult

	// SQL: SELECT PC, count(*), sum(BHMJ) FROM nature_data WHERE year = ? AND BHDL = '资源损毁' GROUP BY PC
	err := s.db.Model(&model.NatureData{}).
		Select("PC, count(*) as count, sum(BHMJ) as area").
		Where("year = ? AND BHDL = ?", year, "资源损毁").
		Group("PC").
		Scan(&results).Error

	return results, err
}

func (s *natureStore) GetRegionStats(year string, groupCol string, filterCol string, filterVal string) ([]model.RegionStatResult, error) {
	var results []model.RegionStatResult

	// 构建基础查询
	// Select: 动态列名 as region_name, count, sum
	tx := s.db.Model(&model.NatureData{}).
		Select(groupCol+" as region_name, count(*) as count, sum(BHMJ) as area").
		Where("year = ?", year)

	// 如果有筛选条件 (比如查河北省下的市)，则追加 Where
	if filterCol != "" && filterVal != "" {
		tx = tx.Where(filterCol+" = ?", filterVal)
	}

	// 执行分组和查询
	err := tx.Group(groupCol).Scan(&results).Error

	return results, err
}
