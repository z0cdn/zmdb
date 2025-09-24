package model

import (
	"time"

	"gorm.io/gorm"
)

// 1. 关系缓存表 (用于加速图查询)
type RelationCache struct {
	gorm.Model
	// 源和目标资源
	SourceID uint `json:"source_id" gorm:"index;not null;comment:'源资源ID'"`
	TargetID uint `json:"target_id" gorm:"index;not null;comment:'目标资源ID'"`

	// 路径信息
	Path       string  `json:"path" gorm:"type:text;not null;comment:'关系路径(JSON格式)'"`
	PathLength int     `json:"path_length" gorm:"type:int;not null;index;comment:'路径长度'"`
	PathWeight float64 `json:"path_weight" gorm:"type:decimal(10,4);comment:'路径权重'"`

	// 关系类型链
	RelationChain string `json:"relation_chain" gorm:"type:varchar(500);not null;comment:'关系类型链'"`

	// 缓存元信息
	CacheTime   time.Time `json:"cache_time" gorm:"not null;index;comment:'缓存时间'"`
	ExpiresAt   time.Time `json:"expires_at" gorm:"not null;index;comment:'过期时间'"`
	AccessCount int       `json:"access_count" gorm:"type:int;default:0;comment:'访问次数'"`
	LastAccess  time.Time `json:"last_access" gorm:"comment:'最后访问时间'"`

	// 关联
	Source Resource `json:"source" gorm:"foreignKey:SourceID"`
	Target Resource `json:"target" gorm:"foreignKey:TargetID"`
}

func (m *RelationCache) TableName() string {
	return "cmdb_relation_cache"
}

// 2. 资源视图缓存表 (预计算的资源视图)
type ResourceView struct {
	gorm.Model
	// 视图信息
	ViewName string `json:"view_name" gorm:"type:varchar(100);not null;index;comment:'视图名称'"`
	ViewType string `json:"view_type" gorm:"type:varchar(50);not null;index;comment:'视图类型'"`

	// 查询条件
	QueryCondition JSONMap `json:"query_condition" gorm:"type:jsonb;not null;comment:'查询条件'"`
	QueryHash      string  `json:"query_hash" gorm:"type:varchar(100);uniqueIndex;not null;comment:'查询条件哈希'"`

	// 结果数据
	ResultData  JSONMap `json:"result_data" gorm:"type:jsonb;not null;comment:'查询结果数据'"`
	ResultCount int     `json:"result_count" gorm:"type:int;not null;comment:'结果数量'"`

	// 缓存信息
	CacheTime time.Time `json:"cache_time" gorm:"not null;index;comment:'缓存时间'"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index;comment:'过期时间'"`
	TTL       int       `json:"ttl" gorm:"type:int;not null;comment:'生存时间(秒)'"`

	// 访问统计
	AccessCount int       `json:"access_count" gorm:"type:int;default:0;comment:'访问次数'"`
	LastAccess  time.Time `json:"last_access" gorm:"comment:'最后访问时间'"`

	// 依赖信息
	DependentResources []string `json:"dependent_resources" gorm:"type:json;comment:'依赖的资源ID列表'"`

	IsActive    bool   `json:"is_active" gorm:"default:true;comment:'是否启用'"`
	Description string `json:"description" gorm:"type:text;comment:'视图描述'"`
}

func (m *ResourceView) TableName() string {
	return "cmdb_resource_views"
}

// 视图类型常量
const (
	ViewTypeTopology   = "topology"   // 拓扑视图
	ViewTypeList       = "list"       // 列表视图
	ViewTypeDashboard  = "dashboard"  // 仪表板视图
	ViewTypeReport     = "report"     // 报表视图
	ViewTypeStatistics = "statistics" // 统计视图
)

// 3. 搜索索引表 (用于全文搜索)
type SearchIndex struct {
	gorm.Model
	// 索引对象信息
	ObjectType string `json:"object_type" gorm:"type:varchar(50);not null;index;comment:'对象类型(resource/service/business)'"`
	ObjectID   uint   `json:"object_id" gorm:"index;not null;comment:'对象ID'"`
	ObjectUUID string `json:"object_uuid" gorm:"type:varchar(100);index;not null;comment:'对象UUID'"`

	// 搜索内容
	Title    string `json:"title" gorm:"type:varchar(500);not null;comment:'标题'"`
	Content  string `json:"content" gorm:"type:text;not null;comment:'搜索内容'"`
	Keywords string `json:"keywords" gorm:"type:text;comment:'关键词'"`
	Tags     string `json:"tags" gorm:"type:text;comment:'标签(逗号分隔)'"`

	// 分类信息
	Category    string `json:"category" gorm:"type:varchar(100);index;comment:'分类'"`
	SubCategory string `json:"sub_category" gorm:"type:varchar(100);index;comment:'子分类'"`

	// 权重和优先级
	Weight   float64 `json:"weight" gorm:"type:decimal(5,2);default:1.0;comment:'搜索权重'"`
	Priority int     `json:"priority" gorm:"type:int;default:1;comment:'优先级'"`

	// 索引元信息
	IndexTime  time.Time `json:"index_time" gorm:"not null;index;comment:'索引时间'"`
	LastUpdate time.Time `json:"last_update" gorm:"not null;comment:'最后更新时间'"`
	Version    int64     `json:"version" gorm:"not null;comment:'版本号'"`

	// 访问统计
	SearchCount int       `json:"search_count" gorm:"type:int;default:0;comment:'搜索次数'"`
	ClickCount  int       `json:"click_count" gorm:"type:int;default:0;comment:'点击次数'"`
	LastAccess  time.Time `json:"last_access" gorm:"comment:'最后访问时间'"`

	// 状态信息
	IsActive  bool `json:"is_active" gorm:"default:true;comment:'是否启用'"`
	IsDeleted bool `json:"is_deleted" gorm:"default:false;comment:'是否已删除'"`
}

func (m *SearchIndex) TableName() string {
	return "cmdb_search_index"
}

// 4. 查询性能监控表
type QueryPerformance struct {
	gorm.Model
	// 查询信息
	QueryType   string  `json:"query_type" gorm:"type:varchar(50);not null;index;comment:'查询类型'"`
	QuerySQL    string  `json:"query_sql" gorm:"type:text;not null;comment:'查询SQL'"`
	QueryHash   string  `json:"query_hash" gorm:"type:varchar(100);index;not null;comment:'查询哈希'"`
	QueryParams JSONMap `json:"query_params" gorm:"type:jsonb;comment:'查询参数'"`

	// 性能指标
	ExecutionTime float64 `json:"execution_time" gorm:"type:decimal(10,4);not null;comment:'执行时间(毫秒)'"`
	RowsScanned   int64   `json:"rows_scanned" gorm:"type:bigint;comment:'扫描行数'"`
	RowsReturned  int64   `json:"rows_returned" gorm:"type:bigint;comment:'返回行数'"`
	MemoryUsed    int64   `json:"memory_used" gorm:"type:bigint;comment:'内存使用量(字节)'"`

	// 时间信息
	QueryTime time.Time `json:"query_time" gorm:"not null;index;comment:'查询时间'"`
	Date      string    `json:"date" gorm:"type:varchar(10);index;not null;comment:'查询日期(YYYY-MM-DD)'"`
	Hour      int       `json:"hour" gorm:"type:int;index;comment:'查询小时(0-23)'"`

	// 用户信息
	UserID    string `json:"user_id" gorm:"type:varchar(100);index;comment:'用户ID'"`
	TenantID  string `json:"tenant_id" gorm:"type:varchar(100);index;comment:'租户ID'"`
	SessionID string `json:"session_id" gorm:"type:varchar(100);comment:'会话ID'"`
	ClientIP  string `json:"client_ip" gorm:"type:varchar(50);comment:'客户端IP'"`
	UserAgent string `json:"user_agent" gorm:"type:varchar(500);comment:'用户代理'"`

	// 状态信息
	Status       string `json:"status" gorm:"type:varchar(50);not null;comment:'查询状态(success/failed/timeout)'"`
	ErrorMessage string `json:"error_message" gorm:"type:text;comment:'错误信息'"`

	// 缓存信息
	CacheHit bool   `json:"cache_hit" gorm:"default:false;comment:'是否命中缓存'"`
	CacheKey string `json:"cache_key" gorm:"type:varchar(200);comment:'缓存键'"`
}

func (m *QueryPerformance) TableName() string {
	return "cmdb_query_performance"
}

// 5. 缓存管理表
type CacheManager struct {
	gorm.Model
	// 缓存信息
	CacheKey   string  `json:"cache_key" gorm:"type:varchar(200);uniqueIndex;not null;comment:'缓存键'"`
	CacheType  string  `json:"cache_type" gorm:"type:varchar(50);not null;index;comment:'缓存类型'"`
	CacheValue JSONMap `json:"cache_value" gorm:"type:jsonb;not null;comment:'缓存值'"`

	// 元信息
	ValueSize   int64     `json:"value_size" gorm:"type:bigint;comment:'缓存值大小(字节)'"`
	TTL         int       `json:"ttl" gorm:"type:int;not null;comment:'生存时间(秒)'"`
	CreatedTime time.Time `json:"created_time" gorm:"not null;comment:'创建时间'"`
	ExpiresAt   time.Time `json:"expires_at" gorm:"not null;index;comment:'过期时间'"`
	LastAccess  time.Time `json:"last_access" gorm:"comment:'最后访问时间'"`

	// 访问统计
	AccessCount int     `json:"access_count" gorm:"type:int;default:0;comment:'访问次数'"`
	HitRate     float64 `json:"hit_rate" gorm:"type:decimal(5,4);comment:'命中率'"`

	// 依赖和标签
	Tags         []string `json:"tags" gorm:"type:json;comment:'缓存标签'"`
	Dependencies []string `json:"dependencies" gorm:"type:json;comment:'依赖的资源'"`

	// 状态
	IsActive bool `json:"is_active" gorm:"default:true;comment:'是否启用'"`
	Priority int  `json:"priority" gorm:"type:int;default:1;comment:'优先级'"`

	Description string `json:"description" gorm:"type:text;comment:'缓存描述'"`
}

func (m *CacheManager) TableName() string {
	return "cmdb_cache_manager"
}

// 缓存类型常量
const (
	CacheTypeQuery     = "query"     // 查询缓存
	CacheTypeRelation  = "relation"  // 关系缓存
	CacheTypeView      = "view"      // 视图缓存
	CacheTypeSearch    = "search"    // 搜索缓存
	CacheTypeStatistic = "statistic" // 统计缓存
	CacheTypeReport    = "report"    // 报表缓存
)

// 6. 数据统计表 (预计算的统计数据)
type DataStatistics struct {
	gorm.Model
	// 统计维度
	StatType  string `json:"stat_type" gorm:"type:varchar(50);not null;index;comment:'统计类型'"`
	Dimension string `json:"dimension" gorm:"type:varchar(100);not null;index;comment:'统计维度'"`
	Period    string `json:"period" gorm:"type:varchar(20);not null;index;comment:'统计周期(day/week/month/quarter/year)'"`
	Date      string `json:"date" gorm:"type:varchar(20);index;not null;comment:'统计日期'"`

	// 统计范围
	TenantID   string `json:"tenant_id" gorm:"type:varchar(100);index;comment:'租户ID'"`
	BusinessID string `json:"business_id" gorm:"type:varchar(100);index;comment:'业务ID'"`
	Provider   string `json:"provider" gorm:"type:varchar(50);index;comment:'云提供商'"`
	Region     string `json:"region" gorm:"type:varchar(100);index;comment:'区域'"`

	// 统计数据
	TotalCount    int64 `json:"total_count" gorm:"type:bigint;not null;comment:'总数量'"`
	ActiveCount   int64 `json:"active_count" gorm:"type:bigint;comment:'活跃数量'"`
	InactiveCount int64 `json:"inactive_count" gorm:"type:bigint;comment:'非活跃数量'"`

	// 详细统计
	Statistics JSONMap `json:"statistics" gorm:"type:jsonb;comment:'详细统计数据'"`
	Trends     JSONMap `json:"trends" gorm:"type:jsonb;comment:'趋势数据'"`

	// 计算信息
	CalcTime     time.Time `json:"calc_time" gorm:"not null;comment:'计算时间'"`
	CalcDuration int64     `json:"calc_duration" gorm:"type:bigint;comment:'计算耗时(毫秒)'"`
	DataSource   string    `json:"data_source" gorm:"type:varchar(100);comment:'数据源'"`

	// 状态
	IsLatest bool  `json:"is_latest" gorm:"default:true;index;comment:'是否最新'"`
	Version  int64 `json:"version" gorm:"not null;comment:'版本号'"`

	Description string `json:"description" gorm:"type:text;comment:'统计描述'"`
}

func (m *DataStatistics) TableName() string {
	return "cmdb_data_statistics"
}

// 统计类型常量
const (
	StatTypeResource    = "resource"    // 资源统计
	StatTypeService     = "service"     // 服务统计
	StatTypeBusiness    = "business"    // 业务统计
	StatTypeRelation    = "relation"    // 关系统计
	StatTypePerformance = "performance" // 性能统计
	StatTypeUsage       = "usage"       // 使用统计
	StatTypeCost        = "cost"        // 成本统计
	StatTypeChange      = "change"      // 变更统计
)
