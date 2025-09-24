package model

import (
	"gorm.io/gorm"
	"time"
)

// ChangeType 变更类型枚举
const (
	ChangeTypeCreate = "create" // 创建
	ChangeTypeUpdate = "update" // 更新
	ChangeTypeDelete = "delete" // 删除
	ChangeTypeSync   = "sync"   // 同步
)

// ChangeSource 变更来源枚举
const (
	ChangeSourceManual    = "manual"    // 手动变更
	ChangeSourceAPI       = "api"       // API变更
	ChangeSourceSync      = "sync"      // 自动同步
	ChangeSourceImport    = "import"    // 批量导入
	ChangeSourceScheduled = "scheduled" // 定时任务
)

// 1. 资源变更历史记录表
type ResourceHistory struct {
	gorm.Model
	// 关联信息
	ResourceID   uint   `json:"resource_id" gorm:"index;not null;comment:'资源ID'"`
	ResourceUUID string `json:"resource_uuid" gorm:"type:varchar(100);index;not null;comment:'资源UUID'"`
	
	// 变更信息
	ChangeType   string    `json:"change_type" gorm:"type:varchar(20);not null;index;comment:'变更类型'"`
	ChangeSource string    `json:"change_source" gorm:"type:varchar(50);not null;index;comment:'变更来源'"`
	ChangeTime   time.Time `json:"change_time" gorm:"not null;index;comment:'变更时间'"`
	
	// 操作人信息
	OperatorID   string `json:"operator_id" gorm:"type:varchar(100);index;comment:'操作人ID'"`
	OperatorName string `json:"operator_name" gorm:"type:varchar(100);comment:'操作人姓名'"`
	OperatorIP   string `json:"operator_ip" gorm:"type:varchar(50);comment:'操作人IP'"`
	
	// 变更内容
	BeforeData   JSONMap `json:"before_data" gorm:"type:jsonb;comment:'变更前数据快照'"`
	AfterData    JSONMap `json:"after_data" gorm:"type:jsonb;comment:'变更后数据快照'"`
	ChangedFields JSONMap `json:"changed_fields" gorm:"type:jsonb;comment:'变更字段详情'"`
	
	// 附加信息
	ChangeReason string `json:"change_reason" gorm:"type:text;comment:'变更原因'"`
	Comment      string `json:"comment" gorm:"type:text;comment:'变更说明'"`
	
	// 版本信息
	Version      int64 `json:"version" gorm:"not null;index;comment:'版本号'"`
	
	// 关联
	Resource     Resource `json:"resource" gorm:"foreignKey:ResourceID"`
}

func (m *ResourceHistory) TableName() string {
	return "cmdb_resource_history"
}

// 2. 服务变更历史记录表
type ServiceHistory struct {
	gorm.Model
	// 关联信息
	ServiceID   uint   `json:"service_id" gorm:"index;not null;comment:'服务ID'"`
	ServiceUUID string `json:"service_uuid" gorm:"type:varchar(100);index;not null;comment:'服务UUID'"`
	
	// 变更信息
	ChangeType   string    `json:"change_type" gorm:"type:varchar(20);not null;index;comment:'变更类型'"`
	ChangeSource string    `json:"change_source" gorm:"type:varchar(50);not null;index;comment:'变更来源'"`
	ChangeTime   time.Time `json:"change_time" gorm:"not null;index;comment:'变更时间'"`
	
	// 操作人信息
	OperatorID   string `json:"operator_id" gorm:"type:varchar(100);index;comment:'操作人ID'"`
	OperatorName string `json:"operator_name" gorm:"type:varchar(100);comment:'操作人姓名'"`
	OperatorIP   string `json:"operator_ip" gorm:"type:varchar(50);comment:'操作人IP'"`
	
	// 变更内容
	BeforeData   JSONMap `json:"before_data" gorm:"type:jsonb;comment:'变更前数据快照'"`
	AfterData    JSONMap `json:"after_data" gorm:"type:jsonb;comment:'变更后数据快照'"`
	ChangedFields JSONMap `json:"changed_fields" gorm:"type:jsonb;comment:'变更字段详情'"`
	
	// 附加信息
	ChangeReason string `json:"change_reason" gorm:"type:text;comment:'变更原因'"`
	Comment      string `json:"comment" gorm:"type:text;comment:'变更说明'"`
	
	// 版本信息
	Version      int64 `json:"version" gorm:"not null;index;comment:'版本号'"`
	
	// 关联
	Service      Service `json:"service" gorm:"foreignKey:ServiceID"`
}

func (m *ServiceHistory) TableName() string {
	return "cmdb_service_history"
}

// 3. 业务变更历史记录表
type BusinessHistory struct {
	gorm.Model
	// 关联信息
	BusinessID   uint   `json:"business_id" gorm:"index;not null;comment:'业务ID'"`
	BusinessUUID string `json:"business_uuid" gorm:"type:varchar(100);index;not null;comment:'业务UUID'"`
	
	// 变更信息
	ChangeType   string    `json:"change_type" gorm:"type:varchar(20);not null;index;comment:'变更类型'"`
	ChangeSource string    `json:"change_source" gorm:"type:varchar(50);not null;index;comment:'变更来源'"`
	ChangeTime   time.Time `json:"change_time" gorm:"not null;index;comment:'变更时间'"`
	
	// 操作人信息
	OperatorID   string `json:"operator_id" gorm:"type:varchar(100);index;comment:'操作人ID'"`
	OperatorName string `json:"operator_name" gorm:"type:varchar(100);comment:'操作人姓名'"`
	OperatorIP   string `json:"operator_ip" gorm:"type:varchar(50);comment:'操作人IP'"`
	
	// 变更内容
	BeforeData   JSONMap `json:"before_data" gorm:"type:jsonb;comment:'变更前数据快照'"`
	AfterData    JSONMap `json:"after_data" gorm:"type:jsonb;comment:'变更后数据快照'"`
	ChangedFields JSONMap `json:"changed_fields" gorm:"type:jsonb;comment:'变更字段详情'"`
	
	// 附加信息
	ChangeReason string `json:"change_reason" gorm:"type:text;comment:'变更原因'"`
	Comment      string `json:"comment" gorm:"type:text;comment:'变更说明'"`
	
	// 版本信息
	Version      int64 `json:"version" gorm:"not null;index;comment:'版本号'"`
	
	// 关联
	Business     Business `json:"business" gorm:"foreignKey:BusinessID"`
}

func (m *BusinessHistory) TableName() string {
	return "cmdb_business_history"
}

// 4. 关系变更历史记录表
type RelationHistory struct {
	gorm.Model
	// 关联信息
	RelationID   uint   `json:"relation_id" gorm:"index;comment:'关系ID(删除时为空)'"`
	SourceID     uint   `json:"source_id" gorm:"index;not null;comment:'源资源ID'"`
	TargetID     uint   `json:"target_id" gorm:"index;not null;comment:'目标资源ID'"`
	RelationType string `json:"relation_type" gorm:"type:varchar(50);not null;index;comment:'关系类型'"`
	
	// 变更信息
	ChangeType   string    `json:"change_type" gorm:"type:varchar(20);not null;index;comment:'变更类型'"`
	ChangeSource string    `json:"change_source" gorm:"type:varchar(50);not null;index;comment:'变更来源'"`
	ChangeTime   time.Time `json:"change_time" gorm:"not null;index;comment:'变更时间'"`
	
	// 操作人信息
	OperatorID   string `json:"operator_id" gorm:"type:varchar(100);index;comment:'操作人ID'"`
	OperatorName string `json:"operator_name" gorm:"type:varchar(100);comment:'操作人姓名'"`
	OperatorIP   string `json:"operator_ip" gorm:"type:varchar(50);comment:'操作人IP'"`
	
	// 变更内容
	BeforeData   JSONMap `json:"before_data" gorm:"type:jsonb;comment:'变更前数据快照'"`
	AfterData    JSONMap `json:"after_data" gorm:"type:jsonb;comment:'变更后数据快照'"`
	ChangedFields JSONMap `json:"changed_fields" gorm:"type:jsonb;comment:'变更字段详情'"`
	
	// 附加信息
	ChangeReason string `json:"change_reason" gorm:"type:text;comment:'变更原因'"`
	Comment      string `json:"comment" gorm:"type:text;comment:'变更说明'"`
	
	// 版本信息
	Version      int64 `json:"version" gorm:"not null;index;comment:'版本号'"`
	
	// 关联
	Relation     *ResourceRelation `json:"relation" gorm:"foreignKey:RelationID"`
	Source       Resource          `json:"source" gorm:"foreignKey:SourceID"`
	Target       Resource          `json:"target" gorm:"foreignKey:TargetID"`
}

func (m *RelationHistory) TableName() string {
	return "cmdb_relation_history"
}

// 5. 快照管理表 (用于定期备份整个资源状态)
type ResourceSnapshot struct {
	gorm.Model
	// 快照信息
	SnapshotID   string    `json:"snapshot_id" gorm:"type:varchar(100);uniqueIndex;not null;comment:'快照唯一标识'"`
	SnapshotTime time.Time `json:"snapshot_time" gorm:"not null;index;comment:'快照时间'"`
	SnapshotType string    `json:"snapshot_type" gorm:"type:varchar(50);not null;index;comment:'快照类型(full/incremental)'"`
	
	// 快照范围
	TenantID     string `json:"tenant_id" gorm:"type:varchar(100);index;comment:'租户ID(为空表示全局快照)'"`
	BusinessID   string `json:"business_id" gorm:"type:varchar(100);index;comment:'业务ID(为空表示租户级快照)'"`
	
	// 快照数据统计
	ResourceCount  int `json:"resource_count" gorm:"type:int;not null;comment:'资源数量'"`
	ServiceCount   int `json:"service_count" gorm:"type:int;not null;comment:'服务数量'"`
	BusinessCount  int `json:"business_count" gorm:"type:int;not null;comment:'业务数量'"`
	RelationCount  int `json:"relation_count" gorm:"type:int;not null;comment:'关系数量'"`
	
	// 快照状态
	Status       string `json:"status" gorm:"type:varchar(50);not null;index;comment:'快照状态(creating/completed/failed)'"`
	
	// 存储信息
	StoragePath  string `json:"storage_path" gorm:"type:varchar(500);comment:'快照存储路径'"`
	FileSize     int64  `json:"file_size" gorm:"type:bigint;comment:'快照文件大小(字节)'"`
	Checksum     string `json:"checksum" gorm:"type:varchar(100);comment:'快照文件校验和'"`
	
	// 创建信息
	CreatorID    string `json:"creator_id" gorm:"type:varchar(100);comment:'创建人ID'"`
	CreatorName  string `json:"creator_name" gorm:"type:varchar(100);comment:'创建人姓名'"`
	
	// 过期信息
	ExpiresAt    *time.Time `json:"expires_at" gorm:"index;comment:'快照过期时间'"`
	
	Description  string `json:"description" gorm:"type:text;comment:'快照描述'"`
}

func (m *ResourceSnapshot) TableName() string {
	return "cmdb_resource_snapshots"
}

// 6. 变更审计配置表
type AuditConfig struct {
	gorm.Model
	// 配置信息
	ResourceType string `json:"resource_type" gorm:"type:varchar(50);not null;index;comment:'资源类型'"`
	TenantID     string `json:"tenant_id" gorm:"type:varchar(100);index;comment:'租户ID(为空表示全局配置)'"`
	
	// 审计设置
	EnableAudit  bool `json:"enable_audit" gorm:"default:true;comment:'是否启用审计'"`
	RetentionDays int `json:"retention_days" gorm:"type:int;default:365;comment:'审计日志保留天数'"`
	
	// 字段监控配置
	MonitoredFields JSONMap `json:"monitored_fields" gorm:"type:jsonb;comment:'需要监控的字段配置'"`
	
	// 通知配置
	NotificationConfig JSONMap `json:"notification_config" gorm:"type:jsonb;comment:'变更通知配置'"`
	
	IsActive     bool   `json:"is_active" gorm:"default:true;comment:'是否启用'"`
	Description  string `json:"description" gorm:"type:text;comment:'配置描述'"`
}

func (m *AuditConfig) TableName() string {
	return "cmdb_audit_configs"
}

// 7. 数据同步日志表 (记录外部数据源同步状态)
type SyncLog struct {
	gorm.Model
	// 同步信息
	SyncID       string    `json:"sync_id" gorm:"type:varchar(100);uniqueIndex;not null;comment:'同步任务唯一标识'"`
	SyncTime     time.Time `json:"sync_time" gorm:"not null;index;comment:'同步时间'"`
	SyncType     string    `json:"sync_type" gorm:"type:varchar(50);not null;index;comment:'同步类型(full/incremental)'"`
	
	// 数据源信息
	DataSource   string `json:"data_source" gorm:"type:varchar(100);not null;index;comment:'数据源标识'"`
	Provider     string `json:"provider" gorm:"type:varchar(50);not null;index;comment:'云提供商'"`
	Region       string `json:"region" gorm:"type:varchar(100);index;comment:'区域'"`
	
	// 同步范围
	TenantID     string `json:"tenant_id" gorm:"type:varchar(100);index;comment:'租户ID'"`
	ResourceTypes []string `json:"resource_types" gorm:"type:json;comment:'同步的资源类型列表'"`
	
	// 同步状态
	Status       string `json:"status" gorm:"type:varchar(50);not null;index;comment:'同步状态(running/completed/failed/partial)'"`
	
	// 同步统计
	TotalCount   int `json:"total_count" gorm:"type:int;not null;comment:'总数量'"`
	SuccessCount int `json:"success_count" gorm:"type:int;not null;comment:'成功数量'"`
	FailedCount  int `json:"failed_count" gorm:"type:int;not null;comment:'失败数量'"`
	SkippedCount int `json:"skipped_count" gorm:"type:int;not null;comment:'跳过数量'"`
	
	// 时间统计
	StartTime    time.Time  `json:"start_time" gorm:"not null;comment:'开始时间'"`
	EndTime      *time.Time `json:"end_time" gorm:"comment:'结束时间'"`
	Duration     int64      `json:"duration" gorm:"type:bigint;comment:'同步耗时(毫秒)'"`
	
	// 错误信息
	ErrorMessage string `json:"error_message" gorm:"type:text;comment:'错误信息'"`
	ErrorDetails JSONMap `json:"error_details" gorm:"type:jsonb;comment:'详细错误信息'"`
	
	// 同步详情
	SyncDetails  JSONMap `json:"sync_details" gorm:"type:jsonb;comment:'同步详情统计'"`
	
	Description  string `json:"description" gorm:"type:text;comment:'同步描述'"`
}

func (m *SyncLog) TableName() string {
	return "cmdb_sync_logs"
}
