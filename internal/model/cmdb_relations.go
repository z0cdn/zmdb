package model

import (
	"gorm.io/gorm"
	"time"
)

// 对象类型枚举
const (
	ObjectTypeResource     = "resource"     // 资源
	ObjectTypeApplication  = "application"  // 应用
	ObjectTypeConfiguration = "configuration" // 配置
	ObjectTypeService      = "service"      // 服务
	ObjectTypeBusiness     = "business"     // 业务
	ObjectTypeProject      = "project"      // 项目
	ObjectTypeTenant       = "tenant"       // 租户
)

// 通用关系类型枚举 (扩展原有关系类型)
const (
	// 基础关系
	UniversalRelationContains   = "contains"   // 包含
	UniversalRelationBelongsTo  = "belongs_to" // 属于
	UniversalRelationDependsOn  = "depends_on" // 依赖
	UniversalRelationProvides   = "provides"   // 提供
	UniversalRelationConsumes   = "consumes"   // 消费
	
	// 部署关系
	UniversalRelationRunsOn     = "runs_on"     // 运行在
	UniversalRelationHosts      = "hosts"       // 托管
	UniversalRelationDeploys    = "deploys"     // 部署
	
	// 配置关系
	UniversalRelationConfigures = "configures" // 配置
	UniversalRelationAppliesTo  = "applies_to" // 应用于
	UniversalRelationInherits   = "inherits"   // 继承
	
	// 网络关系
	UniversalRelationConnectsTo = "connects_to" // 连接到
	UniversalRelationRoutes     = "routes"      // 路由到
	UniversalRelationExposes    = "exposes"     // 暴露
	
	// 业务关系
	UniversalRelationServes     = "serves"      // 服务于
	UniversalRelationManages    = "manages"     // 管理
	UniversalRelationOwns       = "owns"        // 拥有
	UniversalRelationUses       = "uses"        // 使用
	
	// 数据关系
	UniversalRelationBacksUp    = "backs_up"    // 备份
	UniversalRelationReplicates = "replicates"  // 复制
	UniversalRelationSyncs      = "syncs"       // 同步
)

// 1. 通用关系表 (支持多种对象类型之间的关系)
type UniversalRelation struct {
	gorm.Model
	// 关系标识
	RelationID   string `json:"relation_id" gorm:"type:varchar(100);uniqueIndex;not null;comment:'关系唯一标识'"`
	
	// 源对象
	SourceType   string `json:"source_type" gorm:"type:varchar(50);not null;index;comment:'源对象类型'"`
	SourceID     string `json:"source_id" gorm:"type:varchar(100);not null;index;comment:'源对象ID'"`
	SourceName   string `json:"source_name" gorm:"type:varchar(200);comment:'源对象名称(冗余字段)'"`
	
	// 目标对象
	TargetType   string `json:"target_type" gorm:"type:varchar(50);not null;index;comment:'目标对象类型'"`
	TargetID     string `json:"target_id" gorm:"type:varchar(100);not null;index;comment:'目标对象ID'"`
	TargetName   string `json:"target_name" gorm:"type:varchar(200);comment:'目标对象名称(冗余字段)'"`
	
	// 关系属性
	RelationType string  `json:"relation_type" gorm:"type:varchar(50);not null;index;comment:'关系类型'"`
	Direction    string  `json:"direction" gorm:"type:varchar(20);not null;default:'forward';comment:'关系方向'"`
	Weight       float64 `json:"weight" gorm:"type:decimal(10,4);default:1.0;comment:'关系权重'"`
	Priority     int     `json:"priority" gorm:"type:int;default:1;comment:'关系优先级'"`
	
	// 关系属性和配置
	Properties   JSONMap `json:"properties" gorm:"type:jsonb;comment:'关系属性'"`
	Constraints  JSONMap `json:"constraints" gorm:"type:jsonb;comment:'关系约束'"`
	
	// 生效条件
	EffectiveTime *time.Time `json:"effective_time" gorm:"comment:'生效时间'"`
	ExpireTime    *time.Time `json:"expire_time" gorm:"comment:'过期时间'"`
	
	// 环境和租户
	Environment  string `json:"environment" gorm:"type:varchar(50);index;comment:'环境'"`
	TenantID     string `json:"tenant_id" gorm:"type:varchar(100);index;comment:'租户ID'"`
	
	// 状态管理
	Status       string `json:"status" gorm:"type:varchar(50);not null;default:'active';comment:'关系状态'"`
	IsActive     bool   `json:"is_active" gorm:"default:true;comment:'是否启用'"`
	
	// 审计信息
	CreatedBy    string `json:"created_by" gorm:"type:varchar(100);comment:'创建人'"`
	UpdatedBy    string `json:"updated_by" gorm:"type:varchar(100);comment:'更新人'"`
	
	Description  string `json:"description" gorm:"type:text;comment:'关系描述'"`
	Tags         []UniversalRelationTag `json:"tags" gorm:"foreignKey:RelationID;references:ID"`
}

func (m *UniversalRelation) TableName() string {
	return "cmdb_universal_relations"
}

// 2. 通用关系标签表
type UniversalRelationTag struct {
	gorm.Model
	RelationID uint   `json:"relation_id" gorm:"index;not null;comment:'关系ID'"`
	Key        string `json:"key" gorm:"type:varchar(100);not null;index;comment:'标签键'"`
	Value      string `json:"value" gorm:"type:varchar(500);not null;index;comment:'标签值'"`
	
	Relation   UniversalRelation `json:"relation" gorm:"foreignKey:RelationID"`
}

func (m *UniversalRelationTag) TableName() string {
	return "cmdb_universal_relation_tags"
}

// 3. 多维度关联表 (处理复杂的多对多关系)
type MultiDimensionAssociation struct {
	gorm.Model
	// 关联标识
	AssocID      string `json:"assoc_id" gorm:"type:varchar(100);uniqueIndex;not null;comment:'关联唯一标识'"`
	AssocType    string `json:"assoc_type" gorm:"type:varchar(50);not null;index;comment:'关联类型'"`
	
	// 主维度 (比如：资源)
	PrimaryType  string `json:"primary_type" gorm:"type:varchar(50);not null;index;comment:'主对象类型'"`
	PrimaryID    string `json:"primary_id" gorm:"type:varchar(100);not null;index;comment:'主对象ID'"`
	
	// 次维度 (比如：应用)
	SecondaryType string `json:"secondary_type" gorm:"type:varchar(50);not null;index;comment:'次对象类型'"`
	SecondaryID   string `json:"secondary_id" gorm:"type:varchar(100);not null;index;comment:'次对象ID'"`
	
	// 第三维度 (比如：配置)
	TertiaryType  string `json:"tertiary_type" gorm:"type:varchar(50);index;comment:'第三对象类型'"`
	TertiaryID    string `json:"tertiary_id" gorm:"type:varchar(100);index;comment:'第三对象ID'"`
	
	// 第四维度 (比如：业务)
	QuaternaryType string `json:"quaternary_type" gorm:"type:varchar(50);index;comment:'第四对象类型'"`
	QuaternaryID   string `json:"quaternary_id" gorm:"type:varchar(100);index;comment:'第四对象ID'"`
	
	// 关联配置
	AssocConfig  JSONMap `json:"assoc_config" gorm:"type:jsonb;comment:'关联配置'"`
	
	// 权重和优先级
	Weight       float64 `json:"weight" gorm:"type:decimal(10,4);default:1.0;comment:'关联权重'"`
	Priority     int     `json:"priority" gorm:"type:int;default:1;comment:'关联优先级'"`
	
	// 环境和作用域
	Environment  string `json:"environment" gorm:"type:varchar(50);index;comment:'环境'"`
	TenantID     string `json:"tenant_id" gorm:"type:varchar(100);index;comment:'租户ID'"`
	Scope        string `json:"scope" gorm:"type:varchar(50);comment:'作用域'"`
	
	// 生效条件
	EffectiveTime *time.Time `json:"effective_time" gorm:"comment:'生效时间'"`
	ExpireTime    *time.Time `json:"expire_time" gorm:"comment:'过期时间'"`
	
	// 状态
	Status       string `json:"status" gorm:"type:varchar(50);not null;default:'active';comment:'关联状态'"`
	IsActive     bool   `json:"is_active" gorm:"default:true;comment:'是否启用'"`
	
	Description  string `json:"description" gorm:"type:text;comment:'关联描述'"`
}

func (m *MultiDimensionAssociation) TableName() string {
	return "cmdb_multi_dimension_associations"
}

// 关联类型枚举
const (
	AssocTypeDeployment    = "deployment"    // 部署关联 (资源→应用→配置→业务)
	AssocTypeService       = "service"       // 服务关联 (应用→配置→服务→业务)
	AssocTypeDependency    = "dependency"    // 依赖关联
	AssocTypeConfiguration = "configuration" // 配置关联
	AssocTypeMonitoring    = "monitoring"    // 监控关联
	AssocTypeBackup        = "backup"        // 备份关联
)

// 4. 关系路径缓存表 (用于加速复杂查询)
type RelationPath struct {
	gorm.Model
	// 路径标识
	PathID       string `json:"path_id" gorm:"type:varchar(100);uniqueIndex;not null;comment:'路径唯一标识'"`
	
	// 起始和终点
	StartType    string `json:"start_type" gorm:"type:varchar(50);not null;index;comment:'起始对象类型'"`
	StartID      string `json:"start_id" gorm:"type:varchar(100);not null;index;comment:'起始对象ID'"`
	EndType      string `json:"end_type" gorm:"type:varchar(50);not null;index;comment:'终点对象类型'"`
	EndID        string `json:"end_id" gorm:"type:varchar(100);not null;index;comment:'终点对象ID'"`
	
	// 路径信息
	PathType     string  `json:"path_type" gorm:"type:varchar(50);not null;comment:'路径类型'"`
	PathLength   int     `json:"path_length" gorm:"type:int;not null;comment:'路径长度'"`
	PathWeight   float64 `json:"path_weight" gorm:"type:decimal(10,4);comment:'路径权重'"`
	
	// 路径详情
	PathData     JSONMap `json:"path_data" gorm:"type:jsonb;not null;comment:'路径详情数据'"`
	RelationChain string `json:"relation_chain" gorm:"type:varchar(1000);comment:'关系链'"`
	
	// 缓存信息
	CacheTime    time.Time `json:"cache_time" gorm:"not null;comment:'缓存时间'"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"not null;index;comment:'过期时间'"`
	AccessCount  int       `json:"access_count" gorm:"type:int;default:0;comment:'访问次数'"`
	LastAccess   time.Time `json:"last_access" gorm:"comment:'最后访问时间'"`
	
	IsActive     bool   `json:"is_active" gorm:"default:true;comment:'是否启用'"`
}

func (m *RelationPath) TableName() string {
	return "cmdb_relation_paths"
}

// 5. 对象依赖图表 (用于存储预计算的依赖关系)
type DependencyGraph struct {
	gorm.Model
	// 图标识
	GraphID      string `json:"graph_id" gorm:"type:varchar(100);uniqueIndex;not null;comment:'依赖图唯一标识'"`
	GraphType    string `json:"graph_type" gorm:"type:varchar(50);not null;index;comment:'图类型'"`
	
	// 根节点
	RootType     string `json:"root_type" gorm:"type:varchar(50);not null;comment:'根节点类型'"`
	RootID       string `json:"root_id" gorm:"type:varchar(100);not null;comment:'根节点ID'"`
	
	// 图数据
	GraphData    JSONMap `json:"graph_data" gorm:"type:jsonb;not null;comment:'图数据'"`
	NodeCount    int     `json:"node_count" gorm:"type:int;not null;comment:'节点数量'"`
	EdgeCount    int     `json:"edge_count" gorm:"type:int;not null;comment:'边数量'"`
	MaxDepth     int     `json:"max_depth" gorm:"type:int;comment:'最大深度'"`
	
	// 计算信息
	CalcTime     time.Time `json:"calc_time" gorm:"not null;comment:'计算时间'"`
	CalcDuration int64     `json:"calc_duration" gorm:"type:bigint;comment:'计算耗时(毫秒)'"`
	Version      int64     `json:"version" gorm:"not null;comment:'版本号'"`
	
	// 环境和范围
	Environment  string `json:"environment" gorm:"type:varchar(50);index;comment:'环境'"`
	TenantID     string `json:"tenant_id" gorm:"type:varchar(100);index;comment:'租户ID'"`
	
	// 缓存控制
	ExpiresAt    time.Time `json:"expires_at" gorm:"not null;index;comment:'过期时间'"`
	AccessCount  int       `json:"access_count" gorm:"type:int;default:0;comment:'访问次数'"`
	LastAccess   time.Time `json:"last_access" gorm:"comment:'最后访问时间'"`
	
	IsActive     bool   `json:"is_active" gorm:"default:true;comment:'是否启用'"`
	Description  string `json:"description" gorm:"type:text;comment:'图描述'"`
}

func (m *DependencyGraph) TableName() string {
	return "cmdb_dependency_graphs"
}

// 图类型枚举
const (
	GraphTypeDependency    = "dependency"    // 依赖图
	GraphTypeTopology      = "topology"      // 拓扑图
	GraphTypeDataFlow      = "data_flow"     // 数据流图
	GraphTypeDeployment    = "deployment"    // 部署图
	GraphTypeNetworkFlow   = "network_flow"  // 网络流图
)

// 6. 关系规则表 (定义对象间允许的关系类型)
type RelationRule struct {
	gorm.Model
	// 规则标识
	RuleID       string `json:"rule_id" gorm:"type:varchar(100);uniqueIndex;not null;comment:'规则唯一标识'"`
	RuleName     string `json:"rule_name" gorm:"type:varchar(200);not null;comment:'规则名称'"`
	
	// 源和目标类型
	SourceType   string `json:"source_type" gorm:"type:varchar(50);not null;index;comment:'源对象类型'"`
	TargetType   string `json:"target_type" gorm:"type:varchar(50);not null;index;comment:'目标对象类型'"`
	
	// 允许的关系
	AllowedRelations []string `json:"allowed_relations" gorm:"type:json;not null;comment:'允许的关系类型列表'"`
	
	// 约束条件
	Constraints      JSONMap `json:"constraints" gorm:"type:jsonb;comment:'约束条件'"`
	
	// 规则配置
	IsRequired       bool    `json:"is_required" gorm:"default:false;comment:'是否必需'"`
	MaxConnections   int     `json:"max_connections" gorm:"type:int;default:0;comment:'最大连接数(0表示无限制)'"`
	MinConnections   int     `json:"min_connections" gorm:"type:int;default:0;comment:'最小连接数'"`
	
	// 验证配置
	ValidationRules  JSONMap `json:"validation_rules" gorm:"type:jsonb;comment:'验证规则'"`
	
	// 状态
	Priority         int    `json:"priority" gorm:"type:int;default:1;comment:'规则优先级'"`
	IsActive         bool   `json:"is_active" gorm:"default:true;comment:'是否启用'"`
	
	Description      string `json:"description" gorm:"type:text;comment:'规则描述'"`
}

func (m *RelationRule) TableName() string {
	return "cmdb_relation_rules"
}
