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
