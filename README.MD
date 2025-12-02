    protected_area_backend/
    ├── config/             # 存放数据库配置信息 (config.yaml 等)
    ├── internal/
    │   ├── model/      # 1. 定义数据结构
        ├── store/      # 2. 负责 SQL 查询 (GORM)
        ├── service/    # 3. 负责数据加工 (List 转 Map)
        ├── handler/    # 4. 负责 HTTP 请求响应
        └── router/     # 5. 路由配置
    ├── pkg/                # 公共工具包 (比如统一的响应格式)
    ├── main.go             # 项目入口
    └── go.mod              # 依赖管理