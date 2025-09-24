package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// JSONMap 用于存储 JSON 格式的扩展属性
type JSONMap map[string]interface{}

// Value 实现 driver.Valuer 接口
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan 实现 sql.Scanner 接口
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

// ResourceType 资源类型枚举
const (
	// 物理资源
	ResourceTypeServer    = "server"    // 物理服务器
	ResourceTypeVM        = "vm"        // 虚拟机
	ResourceTypeContainer = "container" // 容器
	ResourceTypeNetwork   = "network"   // 网络设备
	ResourceTypeStorage   = "storage"   // 存储设备

	// CDN相关资源
	ResourceTypeCDNNode      = "cdn_node"      // CDN节点
	ResourceTypePOPNode      = "pop_node"      // PoP节点
	ResourceTypeDNSServer    = "dns_server"    // DNS服务器
	ResourceTypeLoadBalancer = "load_balancer" // 负载均衡器

	// 云资源
	ResourceTypeCloudInstance = "cloud_instance" // 云实例
	ResourceTypeCloudDisk     = "cloud_disk"     // 云盘
	ResourceTypeCloudNetwork  = "cloud_network"  // 云网络

	// K8s资源
	ResourceTypeK8sCluster   = "k8s_cluster"   // K8s集群
	ResourceTypeK8sNode      = "k8s_node"      // K8s节点
	ResourceTypeK8sPod       = "k8s_pod"       // K8s Pod
	ResourceTypeK8sService   = "k8s_service"   // K8s Service
	ResourceTypeK8sNamespace = "k8s_namespace" // K8s命名空间
)

// ResourceStatus 资源状态枚举
const (
	ResourceStatusActive      = "active"      // 活跃
	ResourceStatusInactive    = "inactive"    // 非活跃
	ResourceStatusMaintenance = "maintenance" // 维护中
	ResourceStatusFault       = "fault"       // 故障
	ResourceStatusOffline     = "offline"     // 离线
	ResourceStatusTerminated  = "terminated"  // 已终止
)

// Provider 云提供商枚举
const (
	ProviderAWS       = "aws"
	ProviderAliyun    = "aliyun"
	ProviderTencent   = "tencent"
	ProviderBaidu     = "baidu"
	ProviderSelfBuilt = "self_built" // 自建
)

// 1. 资源层 - 核心资源表
type Resource struct {
	gorm.Model
	// 基础字段
	ResourceID string `json:"resource_id" gorm:"type:varchar(100);uniqueIndex;not null;comment:'资源唯一标识(UUID/Instance ID)'"`
	Name       string `json:"name" gorm:"type:varchar(200);not null;comment:'资源名称'"`
	Type       string `json:"type" gorm:"type:varchar(50);not null;index;comment:'资源类型'"`
	Status     string `json:"status" gorm:"type:varchar(50);not null;index;comment:'资源状态'"`

	// 位置信息
	Provider string `json:"provider" gorm:"type:varchar(50);index;comment:'云提供商/数据源'"`
	Region   string `json:"region" gorm:"type:varchar(100);index;comment:'区域'"`
	Zone     string `json:"zone" gorm:"type:varchar(100);index;comment:'可用区'"`

	// 业务信息
	TenantID    string `json:"tenant_id" gorm:"type:varchar(100);index;comment:'租户ID'"`
	BusinessID  string `json:"business_id" gorm:"type:varchar(100);index;comment:'业务ID'"`
	Environment string `json:"environment" gorm:"type:varchar(50);index;comment:'环境(prod/test/dev)'"`

	// 扩展属性 (支持异构资源)
	Attributes JSONMap `json:"attributes" gorm:"type:jsonb;comment:'扩展属性JSON格式'"`

	// 描述信息
	Description string `json:"description" gorm:"type:text;comment:'资源描述'"`

	// 时间戳
	LastSyncTime *time.Time `json:"last_sync_time" gorm:"comment:'最后同步时间'"`

	// 关联
	Tags      []ResourceTag      `json:"tags" gorm:"foreignKey:ResourceID;references:ID"`
	Relations []ResourceRelation `json:"relations" gorm:"foreignKey:SourceID;references:ID"`
}

func (m *Resource) TableName() string {
	return "cmdb_resources"
}

// 2. 资源标签表 (支持标签化管理)
type ResourceTag struct {
	gorm.Model
	ResourceID uint   `json:"resource_id" gorm:"index;not null;comment:'资源ID'"`
	Key        string `json:"key" gorm:"type:varchar(100);not null;index;comment:'标签键'"`
	Value      string `json:"value" gorm:"type:varchar(500);not null;index;comment:'标签值'"`

	// 联合唯一索引
	Resource Resource `json:"resource" gorm:"foreignKey:ResourceID"`
}

func (m *ResourceTag) TableName() string {
	return "cmdb_resource_tags"
}

// 3. 资源关系表 (图结构建模)
type ResourceRelation struct {
	gorm.Model
	SourceID     uint   `json:"source_id" gorm:"index;not null;comment:'源资源ID'"`
	TargetID     uint   `json:"target_id" gorm:"index;not null;comment:'目标资源ID'"`
	RelationType string `json:"relation_type" gorm:"type:varchar(50);not null;index;comment:'关系类型'"`
	Direction    string `json:"direction" gorm:"type:varchar(20);not null;default:'forward';comment:'关系方向(forward/backward/bidirectional)'"`
	Weight       int    `json:"weight" gorm:"type:int;default:1;comment:'关系权重'"`

	// 关系属性
	Properties  JSONMap `json:"properties" gorm:"type:jsonb;comment:'关系属性JSON格式'"`
	Description string  `json:"description" gorm:"type:text;comment:'关系描述'"`

	// 关联
	Source Resource `json:"source" gorm:"foreignKey:SourceID"`
	Target Resource `json:"target" gorm:"foreignKey:TargetID"`
}

func (m *ResourceRelation) TableName() string {
	return "cmdb_resource_relations"
}

// 关系类型常量
const (
	RelationTypeContains     = "contains"      // 包含关系 (集群包含节点)
	RelationTypeRunsOn       = "runs_on"       // 运行在 (容器运行在主机上)
	RelationTypeConnectsTo   = "connects_to"   // 连接到 (网络连接)
	RelationTypeDependsOn    = "depends_on"    // 依赖于
	RelationTypeProvides     = "provides"      // 提供服务
	RelationTypeManages      = "manages"       // 管理关系
	RelationTypeBelongsTo    = "belongs_to"    // 属于
	RelationTypeUses         = "uses"          // 使用
	RelationTypeBacksUp      = "backs_up"      // 备份
	RelationTypeLoadBalances = "load_balances" // 负载均衡
	RelationTypeHosts        = "hosts"         // 托管 (资源托管应用)
	RelationTypeConfigures   = "configures"    // 配置 (配置作用于应用)
	RelationTypeServes       = "serves"        // 服务于 (应用服务于业务)
)

// 4. 资源类型定义表 (元数据驱动)
type ResourceType struct {
	gorm.Model
	TypeName    string `json:"type_name" gorm:"type:varchar(50);uniqueIndex;not null;comment:'资源类型名称'"`
	DisplayName string `json:"display_name" gorm:"type:varchar(100);not null;comment:'显示名称'"`
	Category    string `json:"category" gorm:"type:varchar(50);not null;index;comment:'资源分类'"`
	Icon        string `json:"icon" gorm:"type:varchar(100);comment:'图标'"`
	Color       string `json:"color" gorm:"type:varchar(20);comment:'颜色'"`

	// 属性定义 (JSON Schema)
	AttributeSchema JSONMap `json:"attribute_schema" gorm:"type:jsonb;comment:'属性schema定义'"`

	// 允许的关系类型
	AllowedRelations []string `json:"allowed_relations" gorm:"type:json;comment:'允许的关系类型列表'"`

	Description string `json:"description" gorm:"type:text;comment:'类型描述'"`
	IsActive    bool   `json:"is_active" gorm:"default:true;comment:'是否启用'"`
}

func (m *ResourceType) TableName() string {
	return "cmdb_resource_types"
}

// 5. 服务层定义
type Service struct {
	gorm.Model
	ServiceID string `json:"service_id" gorm:"type:varchar(100);uniqueIndex;not null;comment:'服务唯一标识'"`
	Name      string `json:"name" gorm:"type:varchar(200);not null;comment:'服务名称'"`
	Type      string `json:"type" gorm:"type:varchar(50);not null;index;comment:'服务类型'"`
	Status    string `json:"status" gorm:"type:varchar(50);not null;index;comment:'服务状态'"`

	// 业务信息
	TenantID    string `json:"tenant_id" gorm:"type:varchar(100);index;comment:'租户ID'"`
	BusinessID  string `json:"business_id" gorm:"type:varchar(100);index;comment:'业务ID'"`
	Environment string `json:"environment" gorm:"type:varchar(50);index;comment:'环境'"`

	// 服务配置
	Configuration JSONMap `json:"configuration" gorm:"type:jsonb;comment:'服务配置'"`
	Endpoints     JSONMap `json:"endpoints" gorm:"type:jsonb;comment:'服务端点信息'"`

	// 监控和SLA
	HealthStatus string  `json:"health_status" gorm:"type:varchar(50);comment:'健康状态'"`
	SLATarget    float64 `json:"sla_target" gorm:"type:decimal(5,4);comment:'SLA目标'"`

	Description string `json:"description" gorm:"type:text;comment:'服务描述'"`

	// 关联资源
	ServiceResources []ServiceResource `json:"service_resources" gorm:"foreignKey:ServiceID;references:ID"`
	Tags             []ServiceTag      `json:"tags" gorm:"foreignKey:ServiceID;references:ID"`
}

func (m *Service) TableName() string {
	return "cmdb_services"
}

// 服务资源关联表
type ServiceResource struct {
	gorm.Model
	ServiceID  uint   `json:"service_id" gorm:"index;not null;comment:'服务ID'"`
	ResourceID uint   `json:"resource_id" gorm:"index;not null;comment:'资源ID'"`
	Role       string `json:"role" gorm:"type:varchar(50);comment:'资源在服务中的角色'"`
	Priority   int    `json:"priority" gorm:"type:int;default:1;comment:'优先级'"`

	Service  Service  `json:"service" gorm:"foreignKey:ServiceID"`
	Resource Resource `json:"resource" gorm:"foreignKey:ResourceID"`
}

func (m *ServiceResource) TableName() string {
	return "cmdb_service_resources"
}

// 服务标签表
type ServiceTag struct {
	gorm.Model
	ServiceID uint   `json:"service_id" gorm:"index;not null;comment:'服务ID'"`
	Key       string `json:"key" gorm:"type:varchar(100);not null;index;comment:'标签键'"`
	Value     string `json:"value" gorm:"type:varchar(500);not null;index;comment:'标签值'"`

	Service Service `json:"service" gorm:"foreignKey:ServiceID"`
}

func (m *ServiceTag) TableName() string {
	return "cmdb_service_tags"
}

// 6. 业务层定义
type Business struct {
	gorm.Model
	BusinessID string `json:"business_id" gorm:"type:varchar(100);uniqueIndex;not null;comment:'业务唯一标识'"`
	Name       string `json:"name" gorm:"type:varchar(200);not null;comment:'业务名称'"`
	Type       string `json:"type" gorm:"type:varchar(50);not null;index;comment:'业务类型'"`
	Status     string `json:"status" gorm:"type:varchar(50);not null;index;comment:'业务状态'"`

	// 所属信息
	TenantID string `json:"tenant_id" gorm:"type:varchar(100);index;comment:'租户ID'"`
	OwnerID  string `json:"owner_id" gorm:"type:varchar(100);index;comment:'业务负责人ID'"`
	TeamID   string `json:"team_id" gorm:"type:varchar(100);index;comment:'所属团队ID'"`

	// 业务属性
	Priority   int     `json:"priority" gorm:"type:int;default:1;comment:'业务优先级'"`
	CostCenter string  `json:"cost_center" gorm:"type:varchar(100);comment:'成本中心'"`
	Budget     float64 `json:"budget" gorm:"type:decimal(15,2);comment:'预算'"`

	Description string `json:"description" gorm:"type:text;comment:'业务描述'"`

	// 关联
	BusinessServices []BusinessService `json:"business_services" gorm:"foreignKey:BusinessID;references:ID"`
	Tags             []BusinessTag     `json:"tags" gorm:"foreignKey:BusinessID;references:ID"`
}

func (m *Business) TableName() string {
	return "cmdb_businesses"
}

// 业务服务关联表
type BusinessService struct {
	gorm.Model
	BusinessID  uint   `json:"business_id" gorm:"index;not null;comment:'业务ID'"`
	ServiceID   uint   `json:"service_id" gorm:"index;not null;comment:'服务ID'"`
	Role        string `json:"role" gorm:"type:varchar(50);comment:'服务在业务中的角色'"`
	Criticality string `json:"criticality" gorm:"type:varchar(50);comment:'重要性级别'"`

	Business Business `json:"business" gorm:"foreignKey:BusinessID"`
	Service  Service  `json:"service" gorm:"foreignKey:ServiceID"`
}

func (m *BusinessService) TableName() string {
	return "cmdb_business_services"
}

// 业务标签表
type BusinessTag struct {
	gorm.Model
	BusinessID uint   `json:"business_id" gorm:"index;not null;comment:'业务ID'"`
	Key        string `json:"key" gorm:"type:varchar(100);not null;index;comment:'标签键'"`
	Value      string `json:"value" gorm:"type:varchar(500);not null;index;comment:'标签值'"`

	Business Business `json:"business" gorm:"foreignKey:BusinessID"`
}

func (m *BusinessTag) TableName() string {
	return "cmdb_business_tags"
}
