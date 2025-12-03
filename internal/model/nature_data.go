package model

// NatureData 对应数据库表 nature_data
type NatureData struct {
	TBBH      string  `gorm:"column:TBBH;primaryKey" json:"tbbh"` // 图斑编号 (设为主键)
	BHDL      string  `gorm:"column:BHDL" json:"bhdl"`            // 变化地类
	QLX       string  `gorm:"column:QLX" json:"qlx"`              // 前地类
	HLX       string  `gorm:"column:HLX" json:"hlx"`              // 后地类
	X         float64 `gorm:"column:X" json:"x"`                  // 经度
	Y         float64 `gorm:"column:Y" json:"y"`                  // 纬度
	BHMJ      float64 `gorm:"column:BHMJ" json:"bhmj"`            // 保护地图斑面积
	THBHDMC   string  `gorm:"column:THBHDMC" json:"thbhdmc"`      // 保护地名称
	BHDLX     string  `gorm:"column:BHDLX" json:"bhdlx"`          // 保护地类型
	PC        string  `gorm:"column:PC" json:"pc"`                // 批次
	BQSJ      string  `gorm:"column:BQSJ" json:"bqsj"`
	SQSJ      string  `gorm:"column:SQSJ" json:"sqsj"`
	THXDM     string  `gorm:"column:THXDM" json:"thxdm"`
	THSHENG   string  `gorm:"column:THSHENG" json:"thsheng"` // 省
	THSHI     string  `gorm:"column:THSHI" json:"thshi"`     // 市
	THXIAN    string  `gorm:"column:THXIAN" json:"thxian"`   // 县
	SFCXBH    int     `gorm:"column:SFCXBH" json:"sfcxbh"`
	SQTBBH    string  `gorm:"column:SQTBBH" json:"sqtbbh"`
	THBZ      string  `gorm:"column:THBZ" json:"thbz"`
	YBBHDMC   string  `gorm:"column:YBBHDMC" json:"ybbhdmc"`
	YBBHDLXBM string  `gorm:"column:YBBHDLXBM" json:"ybbhdlxbm"`
	YBSHENG   string  `gorm:"column:YBSHENG" json:"ybsheng"`
	YBSHI     string  `gorm:"column:YBSHI" json:"ybshi"`
	YBXIAN    string  `gorm:"column:YBXIAN" json:"ybxian"`
	YBXBM     string  `gorm:"column:YBXBM" json:"ybxbm"`
	YBBZ1     string  `gorm:"column:YBBZ1" json:"ybbz1"`
	YBBZ2     string  `gorm:"column:YBBZ2" json:"ybbz2"`
	Year      string  `gorm:"column:year" json:"year"` // 年份
}

// TableName 指定表名，防止 GORM 自动加 s
func (NatureData) TableName() string {
	return "nature_data"
}

// StatResult 用于接收数据库 Group By 查询出的原始结果
// 因为 GORM 聚合查询的结果往往不对应原始表结构，所以定义这个 DTO (Data Transfer Object)
type StatResult struct {
	Year  string `json:"year"`
	BHDL  string `json:"bhdl"`
	Count int64  `json:"count"`
}

// BatchStatResult --- 在 internal/model/nature_data.go 中补充一个 DTO 结构体 ---
// (请把这个加到 model 文件夹下的文件中，用于接收批次统计结果)
type BatchStatResult struct {
	PC    string  `json:"pc"`
	Count int64   `json:"count"`
	Area  float64 `json:"area"`
}

// RegionStatResult 用于接收行政区统计结果
type RegionStatResult struct {
	RegionName string  `json:"region_name"` // 可能是省名、市名或县名
	Count      int64   `json:"count"`
	Area       float64 `json:"area"`
}

// NatureQueryRequest 统一的查询参数结构体
type NatureQueryRequest struct {
	Year          string `form:"year" binding:"required"`  // 年份 (必选)
	Scope         string `form:"scope" binding:"required"` // 查询范围: province, city, county (必选)
	RegionName    string `form:"region_name"`              // 行政区名称 (可选)
	ProtectedType string `form:"protected_type"`           // 保护地类型 (可选)
	ChangeType    string `form:"change_type"`              // 变化地类 (可选)

	// 接口3 专用
	QLX string `form:"qlx"` // 前地类 (接口3必选)

	// 分页参数
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=10"`
}

// AlertQueryRequest 预警接口专用请求参数
type AlertQueryRequest struct {
	Year      string  `form:"year" binding:"required"`       // 年份
	AlertArea float64 `form:"alert_area" binding:"required"` // 预警面积阈值
	Page      int     `form:"page,default=1"`
	PageSize  int     `form:"page_size,default=10"`
}

// ProtectedAreaStat 接口1的返回结构
type ProtectedAreaStat struct {
	Name  string  `json:"name"`  // 保护地名称
	Count int64   `json:"count"` // 图斑个数
	Area  float64 `json:"area"`  // 面积
}

// SpotListItem 接口2专用：精简的图斑明细对象
// 只包含 TBBH, QLX, HLX, BHDL，这样 JSON 序列化时就只有这四个字段
type SpotListItem struct {
	TBBH string `json:"tbbh"`
	QLX  string `json:"qlx"`
	HLX  string `json:"hlx"`
	BHDL string `json:"bhdl"`
}

// TransitionStat 接口3的返回结构
type TransitionStat struct {
	HLX        string  `json:"hlx"`         // 后地类
	Count      int64   `json:"count"`       // 个数
	Area       float64 `json:"area"`        // 面积
	CountRatio float64 `json:"count_ratio"` // 个数占比 (%)
	AreaRatio  float64 `json:"area_ratio"`  // 面积占比 (%)
}

// AlertSpotItem 预警图斑返回项 (DTO)
type AlertSpotItem struct {
	THBHDMC string  `json:"thbhdmc"` // 保护地名称
	TBBH    string  `json:"tbbh"`    // 图斑编号
	BHMJ    float64 `json:"bhmj"`    // 图斑面积
	THSHENG string  `json:"thsheng"` // 所属省份
}
