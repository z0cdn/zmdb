package model

import (
	"time"

	"gorm.io/gorm"
)

// 应用类型枚举
const (
	AppTypeDNSServer    = "dns_server"    // DNS服务
	AppTypeWebGateway   = "web_gateway"   // Web网关
	AppTypeCacheService = "cache_service" // 缓存服务
	AppTypeCDNEdge      = "cdn_edge"      // CDN边缘服务
	AppTypeLoadBalancer = "load_balancer" // 负载均衡
	AppTypeDatabase     = "database"      // 数据库
	AppTypeMessageQueue = "message_queue" // 消息队列
	AppTypeMonitoring   = "monitoring"    // 监控服务
	AppTypeLogging      = "logging"       // 日志服务
)

// 应用状态枚举
const (
	AppStatusRunning     = "running"     // 运行中
	AppStatusStopped     = "stopped"     // 已停止
	AppStatusStarting    = "starting"    // 启动中
	AppStatusStopping    = "stopping"    // 停止中
	AppStatusFailed      = "failed"      // 故障
	AppStatusMaintenance = "maintenance" // 维护中
	AppStatusUpgrading   = "upgrading"   // 升级中
)

// 1. 应用定义表 (应用类型的元数据)
type ApplicationType struct {
	gorm.Model
	// 基本信息
	TypeName    string `json:"type_name" gorm:"type:varchar(50);uniqueIndex;not null;comment:'应用类型名称'"`
	DisplayName string `json:"display_name" gorm:"type:varchar(100);not null;comment:'显示名称'"`
	Category    string `json:"category" gorm:"type:varchar(50);not null;index;comment:'应用分类'"`
	Version     string `json:"version" gorm:"type:varchar(50);comment:'应用版本'"`

	// 视觉信息
	Icon  string `json:"icon" gorm:"type:varchar(100);comment:'图标'"`
	Color string `json:"color" gorm:"type:varchar(20);comment:'颜色'"`

	// 资源要求
	ResourceRequirements JSONMap `json:"resource_requirements" gorm:"type:jsonb;comment:'资源要求(CPU/内存/存储等)'"`

	// 配置模板
	ConfigSchema  JSONMap `json:"config_schema" gorm:"type:jsonb;comment:'配置项schema定义'"`
	DefaultConfig JSONMap `json:"default_config" gorm:"type:jsonb;comment:'默认配置'"`

	// 监控和健康检查
	HealthCheckConfig JSONMap `json:"health_check_config" gorm:"type:jsonb;comment:'健康检查配置'"`
	MonitoringConfig  JSONMap `json:"monitoring_config" gorm:"type:jsonb;comment:'监控配置'"`

	// 部署信息
	DeploymentMethods []string `json:"deployment_methods" gorm:"type:json;comment:'支持的部署方式(docker/binary/k8s等)'"`
	SupportedOS       []string `json:"supported_os" gorm:"type:json;comment:'支持的操作系统'"`

	Description string `json:"description" gorm:"type:text;comment:'应用描述'"`
	IsActive    bool   `json:"is_active" gorm:"default:true;comment:'是否启用'"`

	// 关联
	Applications []Application `json:"applications" gorm:"foreignKey:TypeID;references:ID"`
}

func (m *ApplicationType) TableName() string {
	return "cmdb_application_types"
}

// 2. 应用实例表 (部署在资源上的应用实例)
type Application struct {
	gorm.Model
	// 基本信息
	AppID   string `json:"app_id" gorm:"type:varchar(100);uniqueIndex;not null;comment:'应用实例唯一标识'"`
	Name    string `json:"name" gorm:"type:varchar(200);not null;comment:'应用实例名称'"`
	TypeID  uint   `json:"type_id" gorm:"index;not null;comment:'应用类型ID'"`
	Version string `json:"version" gorm:"type:varchar(50);comment:'应用版本'"`
	Status  string `json:"status" gorm:"type:varchar(50);not null;index;comment:'应用状态'"`

	// 部署信息
	ResourceID     uint   `json:"resource_id" gorm:"index;not null;comment:'部署的资源ID'"`
	DeploymentType string `json:"deployment_type" gorm:"type:varchar(50);comment:'部署方式(docker/binary/k8s等)'"`
	WorkingDir     string `json:"working_dir" gorm:"type:varchar(500);comment:'工作目录'"`
	ExecutablePath string `json:"executable_path" gorm:"type:varchar(500);comment:'可执行文件路径'"`

	// 网络配置
	ListenPorts   JSONMap `json:"listen_ports" gorm:"type:jsonb;comment:'监听端口配置'"`
	NetworkConfig JSONMap `json:"network_config" gorm:"type:jsonb;comment:'网络配置'"`

	// 资源使用
	ResourceUsage  JSONMap `json:"resource_usage" gorm:"type:jsonb;comment:'资源使用情况'"`
	ResourceLimits JSONMap `json:"resource_limits" gorm:"type:jsonb;comment:'资源限制'"`

	// 环境和租户
	Environment string `json:"environment" gorm:"type:varchar(50);index;comment:'环境(prod/test/dev)'"`
	TenantID    string `json:"tenant_id" gorm:"type:varchar(100);index;comment:'租户ID'"`

	// 运行时信息
	ProcessID     int        `json:"process_id" gorm:"type:int;comment:'进程ID'"`
	StartTime     *time.Time `json:"start_time" gorm:"comment:'启动时间'"`
	LastHeartbeat *time.Time `json:"last_heartbeat" gorm:"comment:'最后心跳时间'"`

	// 健康状态
	HealthStatus    string     `json:"health_status" gorm:"type:varchar(50);comment:'健康状态'"`
	LastHealthCheck *time.Time `json:"last_health_check" gorm:"comment:'最后健康检查时间'"`

	Description string `json:"description" gorm:"type:text;comment:'应用描述'"`

	// 关联
	ApplicationType ApplicationType  `json:"application_type" gorm:"foreignKey:TypeID"`
	Resource        Resource         `json:"resource" gorm:"foreignKey:ResourceID"`
	Configurations  []Configuration  `json:"configurations" gorm:"foreignKey:ApplicationID;references:ID"`
	Tags            []ApplicationTag `json:"tags" gorm:"foreignKey:ApplicationID;references:ID"`
}

func (m *Application) TableName() string {
	return "cmdb_applications"
}

// 3. 配置实例表 (业务在应用上的配置)
type Configuration struct {
	gorm.Model
	// 基本信息
	ConfigID      string `json:"config_id" gorm:"type:varchar(100);uniqueIndex;not null;comment:'配置唯一标识'"`
	Name          string `json:"name" gorm:"type:varchar(200);not null;comment:'配置名称'"`
	ApplicationID uint   `json:"application_id" gorm:"index;not null;comment:'应用实例ID'"`

	// 业务关联
	BusinessID string `json:"business_id" gorm:"type:varchar(100);index;comment:'业务ID'"`
	ServiceID  string `json:"service_id" gorm:"type:varchar(100);index;comment:'服务ID'"`
	TenantID   string `json:"tenant_id" gorm:"type:varchar(100);index;comment:'租户ID'"`

	// 配置内容
	ConfigType   string  `json:"config_type" gorm:"type:varchar(50);not null;index;comment:'配置类型(service/business/tenant等)'"`
	ConfigData   JSONMap `json:"config_data" gorm:"type:jsonb;not null;comment:'配置数据'"`
	ConfigFormat string  `json:"config_format" gorm:"type:varchar(20);comment:'配置格式(json/yaml/ini等)'"`

	// 配置来源和管理
	Source     string `json:"source" gorm:"type:varchar(50);comment:'配置来源(manual/template/sync等)'"`
	TemplateID string `json:"template_id" gorm:"type:varchar(100);comment:'配置模板ID'"`

	// 优先级和分组
	Priority    int    `json:"priority" gorm:"type:int;default:1;comment:'配置优先级'"`
	ConfigGroup string `json:"config_group" gorm:"type:varchar(100);comment:'配置组'"`

	// 状态管理
	Status      string `json:"status" gorm:"type:varchar(50);not null;comment:'配置状态(active/inactive/pending)'"`
	IsEncrypted bool   `json:"is_encrypted" gorm:"default:false;comment:'是否加密'"`

	// 版本管理
	Version   string     `json:"version" gorm:"type:varchar(50);comment:'配置版本'"`
	CreatedBy string     `json:"created_by" gorm:"type:varchar(100);comment:'创建人'"`
	UpdatedBy string     `json:"updated_by" gorm:"type:varchar(100);comment:'更新人'"`
	AppliedAt *time.Time `json:"applied_at" gorm:"comment:'配置应用时间'"`

	// 生效条件
	EffectiveTime *time.Time `json:"effective_time" gorm:"comment:'生效时间'"`
	ExpireTime    *time.Time `json:"expire_time" gorm:"comment:'过期时间'"`

	Description string `json:"description" gorm:"type:text;comment:'配置描述'"`

	// 关联
	Application Application        `json:"application" gorm:"foreignKey:ApplicationID"`
	Tags        []ConfigurationTag `json:"tags" gorm:"foreignKey:ConfigurationID;references:ID"`
}

func (m *Configuration) TableName() string {
	return "cmdb_configurations"
}

// 4. 应用标签表
type ApplicationTag struct {
	gorm.Model
	ApplicationID uint   `json:"application_id" gorm:"index;not null;comment:'应用ID'"`
	Key           string `json:"key" gorm:"type:varchar(100);not null;index;comment:'标签键'"`
	Value         string `json:"value" gorm:"type:varchar(500);not null;index;comment:'标签值'"`

	Application Application `json:"application" gorm:"foreignKey:ApplicationID"`
}

func (m *ApplicationTag) TableName() string {
	return "cmdb_application_tags"
}

// 5. 配置标签表
type ConfigurationTag struct {
	gorm.Model
	ConfigurationID uint   `json:"configuration_id" gorm:"index;not null;comment:'配置ID'"`
	Key             string `json:"key" gorm:"type:varchar(100);not null;index;comment:'标签键'"`
	Value           string `json:"value" gorm:"type:varchar(500);not null;index;comment:'标签值'"`

	Configuration Configuration `json:"configuration" gorm:"foreignKey:ConfigurationID"`
}

func (m *ConfigurationTag) TableName() string {
	return "cmdb_configuration_tags"
}

// 6. 配置模板表 (配置的模板)
type ConfigurationTemplate struct {
	gorm.Model
	// 基本信息
	TemplateID string `json:"template_id" gorm:"type:varchar(100);uniqueIndex;not null;comment:'模板唯一标识'"`
	Name       string `json:"name" gorm:"type:varchar(200);not null;comment:'模板名称'"`
	AppTypeID  uint   `json:"app_type_id" gorm:"index;not null;comment:'应用类型ID'"`

	// 模板内容
	TemplateData JSONMap `json:"template_data" gorm:"type:jsonb;not null;comment:'模板数据'"`
	Variables    JSONMap `json:"variables" gorm:"type:jsonb;comment:'模板变量定义'"`

	// 适用范围
	BusinessTypes []string `json:"business_types" gorm:"type:json;comment:'适用的业务类型'"`
	Environments  []string `json:"environments" gorm:"type:json;comment:'适用的环境'"`

	// 模板属性
	Version   string `json:"version" gorm:"type:varchar(50);comment:'模板版本'"`
	Category  string `json:"category" gorm:"type:varchar(50);comment:'模板分类'"`
	IsDefault bool   `json:"is_default" gorm:"default:false;comment:'是否默认模板'"`
	IsActive  bool   `json:"is_active" gorm:"default:true;comment:'是否启用'"`

	// 创建信息
	CreatedBy string `json:"created_by" gorm:"type:varchar(100);comment:'创建人'"`
	UpdatedBy string `json:"updated_by" gorm:"type:varchar(100);comment:'更新人'"`

	Description string `json:"description" gorm:"type:text;comment:'模板描述'"`

	// 关联
	ApplicationType ApplicationType `json:"application_type" gorm:"foreignKey:AppTypeID"`
}

func (m *ConfigurationTemplate) TableName() string {
	return "cmdb_configuration_templates"
}

// 7. 应用依赖关系表
type ApplicationDependency struct {
	gorm.Model
	// 依赖关系
	SourceAppID    uint   `json:"source_app_id" gorm:"index;not null;comment:'源应用ID'"`
	TargetAppID    uint   `json:"target_app_id" gorm:"index;not null;comment:'目标应用ID'"`
	DependencyType string `json:"dependency_type" gorm:"type:varchar(50);not null;comment:'依赖类型'"`

	// 依赖详情
	DependencyConfig JSONMap `json:"dependency_config" gorm:"type:jsonb;comment:'依赖配置详情'"`
	IsRequired       bool    `json:"is_required" gorm:"default:true;comment:'是否必需依赖'"`

	// 健康检查
	HealthCheckEnabled bool    `json:"health_check_enabled" gorm:"default:true;comment:'是否启用健康检查'"`
	HealthCheckConfig  JSONMap `json:"health_check_config" gorm:"type:jsonb;comment:'健康检查配置'"`

	Description string `json:"description" gorm:"type:text;comment:'依赖描述'"`

	// 关联
	SourceApp Application `json:"source_app" gorm:"foreignKey:SourceAppID"`
	TargetApp Application `json:"target_app" gorm:"foreignKey:TargetAppID"`
}

func (m *ApplicationDependency) TableName() string {
	return "cmdb_application_dependencies"
}

// 依赖类型枚举
const (
	DependencyTypeNetwork  = "network"  // 网络依赖
	DependencyTypeService  = "service"  // 服务依赖
	DependencyTypeData     = "data"     // 数据依赖
	DependencyTypeConfig   = "config"   // 配置依赖
	DependencyTypeResource = "resource" // 资源依赖
)

// 8. 应用实例组表 (用于管理相同应用的多个实例)
type ApplicationGroup struct {
	gorm.Model
	// 基本信息
	GroupID string `json:"group_id" gorm:"type:varchar(100);uniqueIndex;not null;comment:'应用组唯一标识'"`
	Name    string `json:"name" gorm:"type:varchar(200);not null;comment:'应用组名称'"`
	TypeID  uint   `json:"type_id" gorm:"index;not null;comment:'应用类型ID'"`

	// 业务信息
	BusinessID  string `json:"business_id" gorm:"type:varchar(100);index;comment:'业务ID'"`
	ServiceID   string `json:"service_id" gorm:"type:varchar(100);index;comment:'服务ID'"`
	TenantID    string `json:"tenant_id" gorm:"type:varchar(100);index;comment:'租户ID'"`
	Environment string `json:"environment" gorm:"type:varchar(50);index;comment:'环境'"`

	// 组配置
	GroupConfig JSONMap `json:"group_config" gorm:"type:jsonb;comment:'组级配置'"`

	// 负载均衡配置
	LoadBalancerConfig JSONMap `json:"load_balancer_config" gorm:"type:jsonb;comment:'负载均衡配置'"`

	// 扩展配置
	AutoScaling  JSONMap `json:"auto_scaling" gorm:"type:jsonb;comment:'自动扩展配置'"`
	MinInstances int     `json:"min_instances" gorm:"type:int;default:1;comment:'最小实例数'"`
	MaxInstances int     `json:"max_instances" gorm:"type:int;default:10;comment:'最大实例数'"`

	Description string `json:"description" gorm:"type:text;comment:'应用组描述'"`

	// 关联
	ApplicationType ApplicationType `json:"application_type" gorm:"foreignKey:TypeID"`
	Applications    []Application   `json:"applications" gorm:"many2many:cmdb_application_group_members;"`
}

func (m *ApplicationGroup) TableName() string {
	return "cmdb_application_groups"
}

// 9. 应用组成员表
type ApplicationGroupMember struct {
	gorm.Model
	GroupID       uint   `json:"group_id" gorm:"index;not null;comment:'应用组ID'"`
	ApplicationID uint   `json:"application_id" gorm:"index;not null;comment:'应用ID'"`
	Role          string `json:"role" gorm:"type:varchar(50);comment:'在组中的角色'"`
	Weight        int    `json:"weight" gorm:"type:int;default:1;comment:'权重'"`
	IsActive      bool   `json:"is_active" gorm:"default:true;comment:'是否启用'"`

	// 关联
	Group       ApplicationGroup `json:"group" gorm:"foreignKey:GroupID"`
	Application Application      `json:"application" gorm:"foreignKey:ApplicationID"`
}

func (m *ApplicationGroupMember) TableName() string {
	return "cmdb_application_group_members"
}
