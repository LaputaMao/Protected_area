    protected_area_backend/
    ├── config/             # 存放数据库配置信息 (config.yaml 等)
    ├── internal/
    │   ├── model/          # 存放对应数据库的结构体 (你的 nature_data 就在这)
    │   ├── dao/            # 数据访问层 (专门写 SQL 或 GORM 查询逻辑)
    │   ├── service/        # 业务逻辑层 (处理筛选、统计的具体算法)
    │   └── api/            # 接口层 (处理 HTTP 请求和响应)
    ├── router/             # 路由配置
    ├── pkg/                # 公共工具包 (比如统一的响应格式)
    ├── main.go             # 项目入口
    └── go.mod              # 依赖管理