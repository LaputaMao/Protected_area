package service

import (
	"ProtectedArea/internal/store"
	"fmt"
)

// --- 在文件顶部定义常量 ---
const (
	// ConstProtectedCount 这里的值我先随便写的，小鱼你需要改成项目实际的常量值
	ConstProtectedCount     = 1098        // 保护地个数
	ConstProtectedTotalArea = 88423024.14 // 保护地总面积
)

type NatureService interface {
	GetTrendAnalysis() (map[string]map[string]int64, error)

	GetYearlyOverview(year string) (map[string]interface{}, error)
	GetDamageAnalysisByBatch(year string) (map[string]map[string]interface{}, error)

	GetAdministrativeStats(year, scope, name string) (interface{}, error)
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

// GetYearlyOverview 1. 业务逻辑：获取年度概况（包含常量）
func (s *natureService) GetYearlyOverview(year string) (map[string]interface{}, error) {
	count, area, err := s.store.GetSummaryByYear(year)
	if err != nil {
		return nil, err
	}

	// 组装返回数据
	return map[string]interface{}{
		"year":                 year,
		"total_count":          count,                   // 当年图斑总数
		"total_area":           area,                    // 当年保护地面积总和
		"protected_count":      ConstProtectedCount,     // 常量：保护地个数
		"protected_total_area": ConstProtectedTotalArea, // 常量：保护地总面积
	}, nil
}

// GetDamageAnalysisByBatch 2. 业务逻辑：分批次统计资源损毁
func (s *natureService) GetDamageAnalysisByBatch(year string) (map[string]map[string]interface{}, error) {
	// 1. 先从数据库拿到原始的分组数据
	// 此时 rawStats 里的 PC 可能是 "202301", "2023-01", "A01" 等各种格式
	rawStats, err := s.store.GetDamageStatsByBatch(year)
	if err != nil {
		return nil, err
	}

	// 2. 初始化返回结构
	// 使用 map 暂存累加结果，防止多个原始 PC 对应同一个“第一批次”
	// 例如: "202301" 和 "2023-01" 都应该算进 "第一批次"
	countMap := make(map[string]int64)
	areaMap := make(map[string]float64)

	// 3. 遍历原始数据，进行清洗和聚合
	for _, item := range rawStats {
		// 调用辅助函数，解析出 "第一批次" 等标准名称
		batchName := s.getBatchNameFromPC(item.PC)

		// 累加数据
		countMap[batchName] += item.Count
		areaMap[batchName] += item.Area
	}

	// 4. 组装最终的 JSON 结构
	// 目标格式: {"资源损毁个数": {...}, "资源损毁面积": {...}}
	response := map[string]map[string]interface{}{
		"资源损毁个数": make(map[string]interface{}),
		"资源损毁面积": make(map[string]interface{}),
	}

	for name, count := range countMap {
		response["资源损毁个数"][name] = count
		response["资源损毁面积"][name] = areaMap[name] // 取出对应的面积
	}

	return response, nil
}

// getBatchNameFromPC 辅助函数：根据 PC 字段后两位判断批次
func (s *natureService) getBatchNameFromPC(pc string) string {
	// 防止空字符串或长度不足2位的脏数据
	if len(pc) < 2 {
		return "未知批次"
	}

	// 取最后两位
	suffix := pc[len(pc)-2:]

	switch suffix {
	case "01":
		return "第一批次"
	case "02":
		return "第二批次"
	case "03":
		return "第三批次"
	case "04":
		return "第四批次"
	default:
		// 如果后两位不是 01-04，归类为其他
		return "其他批次"
	}
}

func (s *natureService) GetAdministrativeStats(year, scope, name string) (interface{}, error) {
	// 1. 定义数据库字段映射
	// scope -> 对应的数据库字段名
	colMap := map[string]string{
		"province": "THSHENG",
		"city":     "THSHI",
		"county":   "THXIAN",
	}

	// 2. 校验 scope 是否合法
	currentCol, ok := colMap[scope]
	if !ok {
		return nil, fmt.Errorf("无效的查询范围(scope): %s", scope)
	}

	var groupCol string  // 最终我们要按哪一列分组
	var filterCol string // 我们要筛选哪一列

	// 3. 核心逻辑判断
	if name == "" {
		// 场景 A: 查当前层级的所有数据 (比如 scope=province, 查所有省)
		groupCol = currentCol
		filterCol = "" // 不筛选
	} else {
		// 场景 B: 查指定行政区的下级数据 (比如 scope=province, name=河北, 查河北下的市)

		// 边界检查: 县级没有下级
		if scope == "county" {
			return nil, fmt.Errorf("县级行政区无法查询下级详情")
		}

		filterCol = currentCol // 筛选当前层级 (WHERE THSHENG = '河北')

		// 确定下级分组列
		if scope == "province" {
			groupCol = colMap["city"] // 省 -> 市
		} else if scope == "city" {
			groupCol = colMap["county"] // 市 -> 县
		}
	}

	// 4. 调用 Store
	stats, err := s.store.GetRegionStats(year, groupCol, filterCol, name)
	if err != nil {
		return nil, err
	}

	// 5. 格式化返回数据
	// 返回一个 Map: {"河北省": {count: 10, area: 100}, "河南省": ...}
	response := make(map[string]map[string]interface{})
	for _, item := range stats {
		// 防止空名数据
		regionName := item.RegionName
		if regionName == "" {
			regionName = "未知区域"
		}

		response[regionName] = map[string]interface{}{
			"count": item.Count,
			"area":  item.Area,
		}
	}

	return response, nil
}
