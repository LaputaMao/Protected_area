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

	GetProtectedAreaStats(req model.NatureQueryRequest) ([]model.ProtectedAreaStat, int64, error)
	GetSpotList(req model.NatureQueryRequest) ([]model.SpotListItem, int64, error)
	GetTransitionStats(req model.NatureQueryRequest) ([]model.TransitionStat, error)
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

// buildCommonQuery 构建公共的筛选条件
func (s *natureStore) buildCommonQuery(req model.NatureQueryRequest) *gorm.DB {
	tx := s.db.Model(&model.NatureData{}).Where("year = ?", req.Year)

	// 动态处理行政区范围
	if req.RegionName != "" {
		switch req.Scope {
		case "province":
			tx = tx.Where("THSHENG = ?", req.RegionName)
		case "city":
			tx = tx.Where("THSHI = ?", req.RegionName)
		case "county":
			tx = tx.Where("THXIAN = ?", req.RegionName)
		}
	}
	// 可选筛选
	if req.ProtectedType != "" {
		tx = tx.Where("BHDLX = ?", req.ProtectedType)
	}
	if req.ChangeType != "" {
		tx = tx.Where("BHDL = ?", req.ChangeType)
	}

	return tx
}

// GetProtectedAreaStats 接口1: 按保护地分组统计 (带分页)
func (s *natureStore) GetProtectedAreaStats(req model.NatureQueryRequest) ([]model.ProtectedAreaStat, int64, error) {
	var results []model.ProtectedAreaStat
	var total int64

	// 1. 复用筛选条件
	query := s.buildCommonQuery(req)

	// 2. 计算总组数 (用于前端分页显示总页数)
	// 注意：这里统计的是 DISTINCT THBHDMC 的数量
	if err := query.Distinct("THBHDMC").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 3. 执行分组查询 + 分页
	offset := (req.Page - 1) * req.PageSize
	err := query.Select("THBHDMC as name, count(*) as count, sum(BHMJ) as area").
		Group("THBHDMC").
		Limit(req.PageSize).Offset(offset).
		Scan(&results).Error

	return results, total, err
}

// GetSpotList 接口2: 获取图斑明细列表 (带分页)
func (s *natureStore) GetSpotList(req model.NatureQueryRequest) ([]model.SpotListItem, int64, error) {
	var results []model.SpotListItem // <--- 换成新的精简结构体
	var total int64

	query := s.buildCommonQuery(req)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize

	// GORM 会自动把查询到的字段映射到 SpotListItem 的同名字段上
	err := query.Select("TBBH, QLX, HLX, BHDL").
		Limit(req.PageSize).Offset(offset).
		Scan(&results).Error

	return results, total, err
}

// GetTransitionStats 接口3: 前地类 -> 后地类 流向统计 (不分页，计算占比)
func (s *natureStore) GetTransitionStats(req model.NatureQueryRequest) ([]model.TransitionStat, error) {
	var results []model.TransitionStat

	query := s.buildCommonQuery(req)

	// 额外增加前地类筛选
	if req.QLX != "" {
		query = query.Where("QLX = ?", req.QLX)
	}

	// 按后地类分组统计
	err := query.Select("HLX, count(*) as count, sum(BHMJ) as area").
		Group("HLX").
		Scan(&results).Error

	return results, err
}
