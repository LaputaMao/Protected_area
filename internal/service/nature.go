package service

import (
	"ProtectedArea/internal/store"
)

type NatureService interface {
	GetTrendAnalysis() (map[string]map[string]int64, error)
}

type natureService struct {
	store store.NatureStore
}

func NewNatureService(s store.NatureStore) NatureService {
	return &natureService{store: s}
}

// GetTrendAnalysis 处理业务逻辑：数据格式转换
func (s *natureService) GetTrendAnalysis() (map[string]map[string]int64, error) {
	// 1. 调用 Store 层获取原始数据
	rawStats, err := s.store.GetYearlyTrendStats()
	if err != nil {
		return nil, err
	}

	// 2. 组装业务数据结构
	// 目标格式: {"资源损毁": {"2020": 352}, "恢复治理": {"2020": 63}}
	response := map[string]map[string]int64{
		"资源损毁": make(map[string]int64),
		"恢复治理": make(map[string]int64),
	}

	for _, item := range rawStats {
		if subMap, ok := response[item.BHDL]; ok {
			subMap[item.Year] = item.Count
		}
	}

	return response, nil
}
