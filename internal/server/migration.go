package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	v1 "nunu-layout-admin/api/v1"
	"nunu-layout-admin/internal/model"
	"nunu-layout-admin/pkg/log"
	"nunu-layout-admin/pkg/sid"
	"os"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type MigrateServer struct {
	db  *gorm.DB
	log *log.Logger
	sid *sid.Sid
	e   *casbin.SyncedEnforcer
}

func NewMigrateServer(
	db *gorm.DB,
	log *log.Logger,
	sid *sid.Sid,
	e *casbin.SyncedEnforcer,
) *MigrateServer {
	return &MigrateServer{
		e:   e,
		db:  db,
		log: log,
		sid: sid,
	}
}
func (m *MigrateServer) Start(ctx context.Context) error {
	// 删除现有表
	m.db.Migrator().DropTable(
		&model.AdminUser{},
		&model.Menu{},
		&model.Role{},
		&model.Api{},
		// CMDB 核心表
		&model.Resource{},
		&model.ResourceTag{},
		&model.ResourceRelation{},
		&model.ResourceType{},
		&model.Service{},
		&model.ServiceResource{},
		&model.ServiceTag{},
		&model.Business{},
		&model.BusinessService{},
		&model.BusinessTag{},
		// CMDB 历史表
		&model.ResourceHistory{},
		&model.ServiceHistory{},
		&model.BusinessHistory{},
		&model.RelationHistory{},
		&model.ResourceSnapshot{},
		&model.AuditConfig{},
		&model.SyncLog{},
		// CMDB 缓存表
		&model.RelationCache{},
		&model.ResourceView{},
		&model.SearchIndex{},
		&model.QueryPerformance{},
		&model.CacheManager{},
		&model.DataStatistics{},
		// CMDB 应用层表
		&model.ApplicationType{},
		&model.Application{},
		&model.Configuration{},
		&model.ApplicationTag{},
		&model.ConfigurationTag{},
		&model.ConfigurationTemplate{},
		&model.ApplicationDependency{},
		&model.ApplicationGroup{},
		&model.ApplicationGroupMember{},
		// CMDB 关系表
		&model.UniversalRelation{},
		&model.UniversalRelationTag{},
		&model.MultiDimensionAssociation{},
		&model.RelationPath{},
		&model.DependencyGraph{},
		&model.RelationRule{},
	)

	// 创建新表
	if err := m.db.AutoMigrate(
		&model.AdminUser{},
		&model.Menu{},
		&model.Role{},
		&model.Api{},
		// CMDB 核心表
		&model.Resource{},
		&model.ResourceTag{},
		&model.ResourceRelation{},
		&model.ResourceType{},
		&model.Service{},
		&model.ServiceResource{},
		&model.ServiceTag{},
		&model.Business{},
		&model.BusinessService{},
		&model.BusinessTag{},
		// CMDB 历史表
		&model.ResourceHistory{},
		&model.ServiceHistory{},
		&model.BusinessHistory{},
		&model.RelationHistory{},
		&model.ResourceSnapshot{},
		&model.AuditConfig{},
		&model.SyncLog{},
		// CMDB 缓存表
		&model.RelationCache{},
		&model.ResourceView{},
		&model.SearchIndex{},
		&model.QueryPerformance{},
		&model.CacheManager{},
		&model.DataStatistics{},
		// CMDB 应用层表
		&model.ApplicationType{},
		&model.Application{},
		&model.Configuration{},
		&model.ApplicationTag{},
		&model.ConfigurationTag{},
		&model.ConfigurationTemplate{},
		&model.ApplicationDependency{},
		&model.ApplicationGroup{},
		&model.ApplicationGroupMember{},
		// CMDB 关系表
		&model.UniversalRelation{},
		&model.UniversalRelationTag{},
		&model.MultiDimensionAssociation{},
		&model.RelationPath{},
		&model.DependencyGraph{},
		&model.RelationRule{},
	); err != nil {
		m.log.Error("user migrate error", zap.Error(err))
		return err
	}
	err := m.initialAdminUser(ctx)
	if err != nil {
		m.log.Error("initialAdminUser error", zap.Error(err))
	}

	err = m.initialMenuData(ctx)
	if err != nil {
		m.log.Error("initialMenuData error", zap.Error(err))
	}

	err = m.initialApisData(ctx)
	if err != nil {
		m.log.Error("initialApisData error", zap.Error(err))
	}

	err = m.initialRBAC(ctx)
	if err != nil {
		m.log.Error("initialRBAC error", zap.Error(err))
	}

	err = m.initialCMDBData(ctx)
	if err != nil {
		m.log.Error("initialCMDBData error", zap.Error(err))
	}

	m.log.Info("AutoMigrate success")
	os.Exit(0)
	return nil
}
func (m *MigrateServer) Stop(ctx context.Context) error {
	m.log.Info("AutoMigrate stop")
	return nil
}
func (m *MigrateServer) initialAdminUser(ctx context.Context) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = m.db.Create(&model.AdminUser{
		Model:    gorm.Model{ID: 1},
		Username: "admin",
		Password: string(hashedPassword),
		Nickname: "Admin",
	}).Error
	return m.db.Create(&model.AdminUser{
		Model:    gorm.Model{ID: 2},
		Username: "user",
		Password: string(hashedPassword),
		Nickname: "运营人员",
	}).Error

}
func (m *MigrateServer) initialRBAC(ctx context.Context) error {
	roles := []model.Role{
		{Sid: model.AdminRole, Name: "超级管理员"},
		{Sid: "1000", Name: "运营人员"},
		{Sid: "1001", Name: "访客"},
	}
	if err := m.db.Create(&roles).Error; err != nil {
		return err
	}
	m.e.ClearPolicy()
	err := m.e.SavePolicy()
	if err != nil {
		m.log.Error("m.e.SavePolicy error", zap.Error(err))
		return err
	}
	_, err = m.e.AddRoleForUser(model.AdminUserID, model.AdminRole)
	if err != nil {
		m.log.Error("m.e.AddRoleForUser error", zap.Error(err))
		return err
	}
	menuList := make([]v1.MenuDataItem, 0)
	err = json.Unmarshal([]byte(menuData), &menuList)
	if err != nil {
		m.log.Error("json.Unmarshal error", zap.Error(err))
		return err
	}
	for _, item := range menuList {
		m.addPermissionForRole(model.AdminRole, model.MenuResourcePrefix+item.Path, "read")
	}
	apiList := make([]model.Api, 0)
	err = m.db.Find(&apiList).Error
	if err != nil {
		m.log.Error("m.db.Find(&apiList).Error error", zap.Error(err))
		return err
	}
	for _, api := range apiList {
		m.addPermissionForRole(model.AdminRole, model.ApiResourcePrefix+api.Path, api.Method)
	}

	// 添加运营人员权限
	_, err = m.e.AddRoleForUser("2", "1000")
	if err != nil {
		m.log.Error("m.e.AddRoleForUser error", zap.Error(err))
		return err
	}
	m.addPermissionForRole("1000", model.MenuResourcePrefix+"/profile/basic", "read")
	m.addPermissionForRole("1000", model.MenuResourcePrefix+"/profile/advanced", "read")
	m.addPermissionForRole("1000", model.MenuResourcePrefix+"/profile", "read")
	m.addPermissionForRole("1000", model.MenuResourcePrefix+"/dashboard", "read")
	m.addPermissionForRole("1000", model.MenuResourcePrefix+"/dashboard/workplace", "read")
	m.addPermissionForRole("1000", model.MenuResourcePrefix+"/dashboard/analysis", "read")
	m.addPermissionForRole("1000", model.MenuResourcePrefix+"/account/settings", "read")
	m.addPermissionForRole("1000", model.MenuResourcePrefix+"/account/center", "read")
	m.addPermissionForRole("1000", model.MenuResourcePrefix+"/account", "read")
	m.addPermissionForRole("1000", model.ApiResourcePrefix+"/v1/menus", http.MethodGet)
	m.addPermissionForRole("1000", model.ApiResourcePrefix+"/v1/admin/user", http.MethodGet)

	return nil
}
func (m *MigrateServer) addPermissionForRole(role, resource, action string) {
	_, err := m.e.AddPermissionForUser(role, resource, action)
	if err != nil {
		m.log.Sugar().Info("为角色 %s 添加权限 %s:%s 失败: %v", role, resource, action, err)
		return
	}
	fmt.Printf("为角色 %s 添加权限: %s %s\n", role, resource, action)
}
func (m *MigrateServer) initialApisData(ctx context.Context) error {
	initialApis := []model.Api{

		{Group: "基础API", Name: "获取用户菜单列表", Path: "/v1/menus", Method: http.MethodGet},
		{Group: "基础API", Name: "获取管理员信息", Path: "/v1/admin/user", Method: http.MethodGet},

		{Group: "菜单管理", Name: "获取管理菜单", Path: "/v1/admin/menus", Method: http.MethodGet},
		{Group: "菜单管理", Name: "创建菜单", Path: "/v1/admin/menu", Method: http.MethodPost},
		{Group: "菜单管理", Name: "更新菜单", Path: "/v1/admin/menu", Method: http.MethodPut},
		{Group: "菜单管理", Name: "删除菜单", Path: "/v1/admin/menu", Method: http.MethodDelete},

		{Group: "权限模块", Name: "获取用户权限", Path: "/v1/admin/user/permissions", Method: http.MethodGet},
		{Group: "权限模块", Name: "获取角色权限", Path: "/v1/admin/role/permissions", Method: http.MethodGet},
		{Group: "权限模块", Name: "更新角色权限", Path: "/v1/admin/role/permission", Method: http.MethodPut},
		{Group: "权限模块", Name: "获取角色列表", Path: "/v1/admin/roles", Method: http.MethodGet},
		{Group: "权限模块", Name: "创建角色", Path: "/v1/admin/role", Method: http.MethodPost},
		{Group: "权限模块", Name: "更新角色", Path: "/v1/admin/role", Method: http.MethodPut},
		{Group: "权限模块", Name: "删除角色", Path: "/v1/admin/role", Method: http.MethodDelete},

		{Group: "权限模块", Name: "获取管理员列表", Path: "/v1/admin/users", Method: http.MethodGet},
		{Group: "权限模块", Name: "更新管理员信息", Path: "/v1/admin/user", Method: http.MethodPut},
		{Group: "权限模块", Name: "创建管理员账号", Path: "/v1/admin/user", Method: http.MethodPost},
		{Group: "权限模块", Name: "删除管理员", Path: "/v1/admin/user", Method: http.MethodDelete},

		{Group: "权限模块", Name: "获取API列表", Path: "/v1/admin/apis", Method: http.MethodGet},
		{Group: "权限模块", Name: "创建API", Path: "/v1/admin/api", Method: http.MethodPost},
		{Group: "权限模块", Name: "更新API", Path: "/v1/admin/api", Method: http.MethodPut},
		{Group: "权限模块", Name: "删除API", Path: "/v1/admin/api", Method: http.MethodDelete},
	}

	return m.db.Create(&initialApis).Error
}
func (m *MigrateServer) initialMenuData(ctx context.Context) error {
	menuList := make([]v1.MenuDataItem, 0)
	err := json.Unmarshal([]byte(menuData), &menuList)
	if err != nil {
		m.log.Error("json.Unmarshal error", zap.Error(err))
		return err
	}
	menuListDb := make([]model.Menu, 0)
	for _, item := range menuList {
		menuListDb = append(menuListDb, model.Menu{
			Model: gorm.Model{
				ID: item.ID,
			},
			ParentID:   item.ParentID,
			Path:       item.Path,
			Title:      item.Title,
			Name:       item.Name,
			Component:  item.Component,
			Locale:     item.Locale,
			Weight:     item.Weight,
			Icon:       item.Icon,
			Redirect:   item.Redirect,
			URL:        item.URL,
			KeepAlive:  item.KeepAlive,
			HideInMenu: item.HideInMenu,
		})
	}
	return m.db.Create(&menuListDb).Error
}

// 初始化CMDB基础数据
func (m *MigrateServer) initialCMDBData(ctx context.Context) error {
	m.log.Info("开始初始化CMDB基础数据...")

	// 1. 初始化资源类型定义
	resourceTypes := []model.ResourceType{
		{
			TypeName:    model.ResourceTypeServer,
			DisplayName: "物理服务器",
			Category:    "infrastructure",
			Icon:        "server",
			Color:       "#1890ff",
			AttributeSchema: model.JSONMap{
				"properties": map[string]interface{}{
					"cpu_cores":    map[string]interface{}{"type": "integer", "description": "CPU核数"},
					"memory_gb":    map[string]interface{}{"type": "integer", "description": "内存大小(GB)"},
					"disk_gb":      map[string]interface{}{"type": "integer", "description": "磁盘大小(GB)"},
					"ip_address":   map[string]interface{}{"type": "string", "description": "IP地址"},
					"os":           map[string]interface{}{"type": "string", "description": "操作系统"},
					"manufacturer": map[string]interface{}{"type": "string", "description": "制造商"},
					"model":        map[string]interface{}{"type": "string", "description": "型号"},
				},
			},
			AllowedRelations: []string{model.RelationTypeContains, model.RelationTypeRunsOn, model.RelationTypeConnectsTo},
			Description:      "物理服务器资源类型",
			IsActive:         true,
		},
		{
			TypeName:    model.ResourceTypeCDNNode,
			DisplayName: "CDN节点",
			Category:    "cdn",
			Icon:        "cdn",
			Color:       "#52c41a",
			AttributeSchema: model.JSONMap{
				"properties": map[string]interface{}{
					"bandwidth_mbps":      map[string]interface{}{"type": "integer", "description": "带宽(Mbps)"},
					"cache_size_gb":       map[string]interface{}{"type": "integer", "description": "缓存大小(GB)"},
					"location":            map[string]interface{}{"type": "string", "description": "地理位置"},
					"provider":            map[string]interface{}{"type": "string", "description": "CDN提供商"},
					"node_type":           map[string]interface{}{"type": "string", "description": "节点类型"},
					"supported_protocols": map[string]interface{}{"type": "array", "description": "支持的协议"},
				},
			},
			AllowedRelations: []string{model.RelationTypeConnectsTo, model.RelationTypeProvides, model.RelationTypeBelongsTo},
			Description:      "CDN边缘节点资源类型",
			IsActive:         true,
		},
		{
			TypeName:    model.ResourceTypeK8sCluster,
			DisplayName: "K8s集群",
			Category:    "container",
			Icon:        "kubernetes",
			Color:       "#722ed1",
			AttributeSchema: model.JSONMap{
				"properties": map[string]interface{}{
					"version":        map[string]interface{}{"type": "string", "description": "Kubernetes版本"},
					"node_count":     map[string]interface{}{"type": "integer", "description": "节点数量"},
					"master_count":   map[string]interface{}{"type": "integer", "description": "主节点数量"},
					"network_plugin": map[string]interface{}{"type": "string", "description": "网络插件"},
					"storage_class":  map[string]interface{}{"type": "string", "description": "存储类"},
				},
			},
			AllowedRelations: []string{model.RelationTypeContains, model.RelationTypeManages, model.RelationTypeRunsOn},
			Description:      "Kubernetes集群资源类型",
			IsActive:         true,
		},
	}

	if err := m.db.Create(&resourceTypes).Error; err != nil {
		m.log.Error("创建资源类型失败", zap.Error(err))
		return err
	}

	// 2. 初始化审计配置
	auditConfigs := []model.AuditConfig{
		{
			ResourceType:  "*",
			TenantID:      "",
			EnableAudit:   true,
			RetentionDays: 365,
			MonitoredFields: model.JSONMap{
				"common":    []string{"name", "status", "tags", "attributes"},
				"sensitive": []string{"ip_address", "password", "secret"},
			},
			NotificationConfig: model.JSONMap{
				"email": map[string]interface{}{
					"enabled":    true,
					"recipients": []string{"admin@example.com"},
				},
				"webhook": map[string]interface{}{
					"enabled": false,
					"url":     "",
				},
			},
			IsActive:    true,
			Description: "全局审计配置",
		},
	}

	if err := m.db.Create(&auditConfigs).Error; err != nil {
		m.log.Error("创建审计配置失败", zap.Error(err))
		return err
	}

	// 3. 创建示例资源数据
	resources := []model.Resource{
		{
			ResourceID:  "server-001",
			Name:        "Web服务器-01",
			Type:        model.ResourceTypeServer,
			Status:      model.ResourceStatusActive,
			Provider:    model.ProviderSelfBuilt,
			Region:      "beijing",
			Zone:        "beijing-a",
			TenantID:    "tenant-001",
			BusinessID:  "web-service",
			Environment: "prod",
			Attributes: model.JSONMap{
				"cpu_cores":    8,
				"memory_gb":    32,
				"disk_gb":      500,
				"ip_address":   "192.168.1.10",
				"os":           "Ubuntu 20.04",
				"manufacturer": "Dell",
				"model":        "PowerEdge R440",
			},
			Description: "生产环境Web服务器",
		},
		{
			ResourceID:  "cdn-node-001",
			Name:        "北京CDN节点-01",
			Type:        model.ResourceTypeCDNNode,
			Status:      model.ResourceStatusActive,
			Provider:    model.ProviderAliyun,
			Region:      "beijing",
			Zone:        "beijing-a",
			TenantID:    "tenant-001",
			BusinessID:  "cdn-service",
			Environment: "prod",
			Attributes: model.JSONMap{
				"bandwidth_mbps":      1000,
				"cache_size_gb":       1024,
				"location":            "北京市朝阳区",
				"provider":            "阿里云CDN",
				"node_type":           "edge",
				"supported_protocols": []string{"HTTP", "HTTPS", "HTTP/2"},
			},
			Description: "北京地区CDN边缘节点",
		},
	}

	if err := m.db.Create(&resources).Error; err != nil {
		m.log.Error("创建示例资源失败", zap.Error(err))
		return err
	}

	// 4. 创建资源标签
	var serverResource, cdnResource model.Resource
	m.db.Where("resource_id = ?", "server-001").First(&serverResource)
	m.db.Where("resource_id = ?", "cdn-node-001").First(&cdnResource)

	tags := []model.ResourceTag{
		{ResourceID: serverResource.ID, Key: "environment", Value: "production"},
		{ResourceID: serverResource.ID, Key: "team", Value: "backend"},
		{ResourceID: serverResource.ID, Key: "criticality", Value: "high"},
		{ResourceID: cdnResource.ID, Key: "environment", Value: "production"},
		{ResourceID: cdnResource.ID, Key: "team", Value: "infrastructure"},
		{ResourceID: cdnResource.ID, Key: "region", Value: "north-china"},
	}

	if err := m.db.Create(&tags).Error; err != nil {
		m.log.Error("创建资源标签失败", zap.Error(err))
		return err
	}

	// 5. 创建服务定义
	services := []model.Service{
		{
			ServiceID:   "web-service-001",
			Name:        "Web应用服务",
			Type:        "web_application",
			Status:      "running",
			TenantID:    "tenant-001",
			BusinessID:  "web-service",
			Environment: "prod",
			Configuration: model.JSONMap{
				"port":          80,
				"protocol":      "HTTP",
				"load_balancer": true,
				"auto_scaling":  true,
			},
			Endpoints: model.JSONMap{
				"public":  "https://api.example.com",
				"private": "http://192.168.1.10:8080",
			},
			HealthStatus: "healthy",
			SLATarget:    99.9,
			Description:  "主要的Web应用服务",
		},
	}

	if err := m.db.Create(&services).Error; err != nil {
		m.log.Error("创建示例服务失败", zap.Error(err))
		return err
	}

	// 6. 创建业务定义
	businesses := []model.Business{
		{
			BusinessID:  "web-service",
			Name:        "Web服务业务线",
			Type:        "web_application",
			Status:      "active",
			TenantID:    "tenant-001",
			OwnerID:     "owner-001",
			TeamID:      "team-backend",
			Priority:    1,
			CostCenter:  "CC-001",
			Budget:      100000.00,
			Description: "公司主要的Web服务业务线",
		},
	}

	if err := m.db.Create(&businesses).Error; err != nil {
		m.log.Error("创建示例业务失败", zap.Error(err))
		return err
	}

	// 7. 创建应用类型定义
	applicationTypes := []model.ApplicationType{
		{
			TypeName:    model.AppTypeDNSServer,
			DisplayName: "DNS服务器",
			Category:    "network",
			Version:     "1.0.0",
			Icon:        "dns",
			Color:       "#13c2c2",
			ResourceRequirements: model.JSONMap{
				"min_cpu":    0.5,
				"min_memory": 512,
				"min_disk":   1024,
			},
			ConfigSchema: model.JSONMap{
				"properties": map[string]interface{}{
					"listen_port":      map[string]interface{}{"type": "integer", "default": 53},
					"upstream_dns":     map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
					"cache_size":       map[string]interface{}{"type": "integer", "default": 1000},
					"log_level":        map[string]interface{}{"type": "string", "default": "info"},
					"enable_recursion": map[string]interface{}{"type": "boolean", "default": true},
				},
			},
			DefaultConfig: model.JSONMap{
				"listen_port":      53,
				"upstream_dns":     []string{"8.8.8.8", "8.8.4.4"},
				"cache_size":       1000,
				"log_level":        "info",
				"enable_recursion": true,
			},
			HealthCheckConfig: model.JSONMap{
				"method":   "dns_query",
				"query":    "health.check",
				"timeout":  5,
				"interval": 30,
			},
			DeploymentMethods: []string{"binary", "docker"},
			SupportedOS:       []string{"linux", "windows"},
			Description:       "DNS服务器应用类型",
			IsActive:          true,
		},
		{
			TypeName:    model.AppTypeCacheService,
			DisplayName: "缓存服务",
			Category:    "cache",
			Version:     "1.0.0",
			Icon:        "cache",
			Color:       "#f5222d",
			ResourceRequirements: model.JSONMap{
				"min_cpu":    1.0,
				"min_memory": 1024,
				"min_disk":   2048,
			},
			ConfigSchema: model.JSONMap{
				"properties": map[string]interface{}{
					"port":        map[string]interface{}{"type": "integer", "default": 6379},
					"max_memory":  map[string]interface{}{"type": "string", "default": "1gb"},
					"persistence": map[string]interface{}{"type": "boolean", "default": true},
					"max_clients": map[string]interface{}{"type": "integer", "default": 10000},
					"timeout":     map[string]interface{}{"type": "integer", "default": 300},
				},
			},
			DefaultConfig: model.JSONMap{
				"port":        6379,
				"max_memory":  "1gb",
				"persistence": true,
				"max_clients": 10000,
				"timeout":     300,
			},
			HealthCheckConfig: model.JSONMap{
				"method":   "tcp_connect",
				"timeout":  3,
				"interval": 15,
			},
			DeploymentMethods: []string{"binary", "docker", "k8s"},
			SupportedOS:       []string{"linux", "windows", "macos"},
			Description:       "缓存服务应用类型",
			IsActive:          true,
		},
	}

	if err := m.db.Create(&applicationTypes).Error; err != nil {
		m.log.Error("创建应用类型失败", zap.Error(err))
		return err
	}

	// 8. 创建应用实例
	var dnsAppType, cacheAppType model.ApplicationType
	m.db.Where("type_name = ?", model.AppTypeDNSServer).First(&dnsAppType)
	m.db.Where("type_name = ?", model.AppTypeCacheService).First(&cacheAppType)

	applications := []model.Application{
		{
			AppID:          "dns-app-001",
			Name:           "DNS服务-01",
			TypeID:         dnsAppType.ID,
			Version:        "1.2.3",
			Status:         model.AppStatusRunning,
			ResourceID:     serverResource.ID,
			DeploymentType: "binary",
			WorkingDir:     "/opt/dns-server",
			ExecutablePath: "/opt/dns-server/bin/dns-server",
			ListenPorts: model.JSONMap{
				"dns": 53,
				"api": 8053,
			},
			NetworkConfig: model.JSONMap{
				"bind_ip":    "0.0.0.0",
				"interfaces": []string{"eth0"},
			},
			Environment:  "prod",
			TenantID:     "tenant-001",
			ProcessID:    12345,
			HealthStatus: "healthy",
			Description:  "生产环境DNS服务实例",
		},
		{
			AppID:          "cache-app-001",
			Name:           "缓存服务-01",
			TypeID:         cacheAppType.ID,
			Version:        "6.2.0",
			Status:         model.AppStatusRunning,
			ResourceID:     serverResource.ID,
			DeploymentType: "docker",
			WorkingDir:     "/data/cache",
			ListenPorts: model.JSONMap{
				"cache":   6379,
				"cluster": 16379,
			},
			ResourceUsage: model.JSONMap{
				"cpu_percent":    25.5,
				"memory_used_mb": 512,
				"disk_used_mb":   1024,
			},
			Environment:  "prod",
			TenantID:     "tenant-001",
			ProcessID:    12346,
			HealthStatus: "healthy",
			Description:  "生产环境缓存服务实例",
		},
	}

	if err := m.db.Create(&applications).Error; err != nil {
		m.log.Error("创建应用实例失败", zap.Error(err))
		return err
	}

	// 9. 创建配置实例
	var dnsApp, cacheApp model.Application
	m.db.Where("app_id = ?", "dns-app-001").First(&dnsApp)
	m.db.Where("app_id = ?", "cache-app-001").First(&cacheApp)

	configurations := []model.Configuration{
		{
			ConfigID:      "dns-config-web-001",
			Name:          "DNS-Web业务配置",
			ApplicationID: dnsApp.ID,
			BusinessID:    "web-service",
			ServiceID:     "web-service-001",
			TenantID:      "tenant-001",
			ConfigType:    "business",
			ConfigData: model.JSONMap{
				"zones": []map[string]interface{}{
					{
						"name": "example.com",
						"type": "master",
						"records": []map[string]interface{}{
							{"name": "@", "type": "A", "value": "192.168.1.10"},
							{"name": "www", "type": "A", "value": "192.168.1.10"},
							{"name": "api", "type": "A", "value": "192.168.1.10"},
						},
					},
				},
				"upstream_dns": []string{"8.8.8.8", "1.1.1.1"},
				"cache_size":   2000,
			},
			ConfigFormat: "json",
			Source:       "manual",
			Priority:     1,
			ConfigGroup:  "web-business",
			Status:       "active",
			Version:      "1.0.0",
			CreatedBy:    "admin",
			Description:  "Web业务的DNS配置",
		},
		{
			ConfigID:      "cache-config-web-001",
			Name:          "缓存-Web业务配置",
			ApplicationID: cacheApp.ID,
			BusinessID:    "web-service",
			ServiceID:     "web-service-001",
			TenantID:      "tenant-001",
			ConfigType:    "business",
			ConfigData: model.JSONMap{
				"databases": map[string]interface{}{
					"web_session": map[string]interface{}{
						"db":       0,
						"ttl":      3600,
						"max_size": "100mb",
					},
					"web_cache": map[string]interface{}{
						"db":       1,
						"ttl":      1800,
						"max_size": "500mb",
					},
				},
				"eviction_policy": "allkeys-lru",
				"max_memory":      "1gb",
			},
			ConfigFormat: "json",
			Source:       "template",
			Priority:     1,
			ConfigGroup:  "web-business",
			Status:       "active",
			Version:      "1.0.0",
			CreatedBy:    "admin",
			Description:  "Web业务的缓存配置",
		},
	}

	if err := m.db.Create(&configurations).Error; err != nil {
		m.log.Error("创建配置实例失败", zap.Error(err))
		return err
	}

	// 10. 创建多维度关联关系
	associations := []model.MultiDimensionAssociation{
		{
			AssocID:        "deploy-web-dns-001",
			AssocType:      model.AssocTypeDeployment,
			PrimaryType:    model.ObjectTypeResource,
			PrimaryID:      serverResource.ResourceID,
			SecondaryType:  model.ObjectTypeApplication,
			SecondaryID:    dnsApp.AppID,
			TertiaryType:   model.ObjectTypeConfiguration,
			TertiaryID:     "dns-config-web-001",
			QuaternaryType: model.ObjectTypeBusiness,
			QuaternaryID:   "web-service",
			AssocConfig: model.JSONMap{
				"deployment_mode": "single",
				"auto_start":      true,
				"health_check":    true,
			},
			Weight:      1.0,
			Priority:    1,
			Environment: "prod",
			TenantID:    "tenant-001",
			Status:      "active",
			IsActive:    true,
			Description: "Web业务DNS服务部署关联",
		},
		{
			AssocID:        "deploy-web-cache-001",
			AssocType:      model.AssocTypeDeployment,
			PrimaryType:    model.ObjectTypeResource,
			PrimaryID:      serverResource.ResourceID,
			SecondaryType:  model.ObjectTypeApplication,
			SecondaryID:    cacheApp.AppID,
			TertiaryType:   model.ObjectTypeConfiguration,
			TertiaryID:     "cache-config-web-001",
			QuaternaryType: model.ObjectTypeBusiness,
			QuaternaryID:   "web-service",
			AssocConfig: model.JSONMap{
				"deployment_mode": "single",
				"auto_start":      true,
				"health_check":    true,
				"backup_enabled":  true,
			},
			Weight:      1.0,
			Priority:    2,
			Environment: "prod",
			TenantID:    "tenant-001",
			Status:      "active",
			IsActive:    true,
			Description: "Web业务缓存服务部署关联",
		},
	}

	if err := m.db.Create(&associations).Error; err != nil {
		m.log.Error("创建多维度关联失败", zap.Error(err))
		return err
	}

	// 11. 创建通用关系
	universalRelations := []model.UniversalRelation{
		{
			RelationID:   "rel-resource-app-dns-001",
			SourceType:   model.ObjectTypeResource,
			SourceID:     serverResource.ResourceID,
			SourceName:   serverResource.Name,
			TargetType:   model.ObjectTypeApplication,
			TargetID:     dnsApp.AppID,
			TargetName:   dnsApp.Name,
			RelationType: model.UniversalRelationHosts,
			Direction:    "forward",
			Weight:       1.0,
			Priority:     1,
			Properties: model.JSONMap{
				"port_mapping": map[string]interface{}{
					"53":   "dns",
					"8053": "api",
				},
				"resource_allocation": map[string]interface{}{
					"cpu_cores": 1,
					"memory_mb": 512,
				},
			},
			Environment: "prod",
			TenantID:    "tenant-001",
			Status:      "active",
			IsActive:    true,
			CreatedBy:   "system",
			Description: "服务器托管DNS应用",
		},
		{
			RelationID:   "rel-app-config-dns-001",
			SourceType:   model.ObjectTypeApplication,
			SourceID:     dnsApp.AppID,
			SourceName:   dnsApp.Name,
			TargetType:   model.ObjectTypeConfiguration,
			TargetID:     "dns-config-web-001",
			TargetName:   "DNS-Web业务配置",
			RelationType: model.UniversalRelationAppliesTo,
			Direction:    "forward",
			Weight:       1.0,
			Priority:     1,
			Properties: model.JSONMap{
				"config_priority": 1,
				"reload_required": true,
			},
			Environment: "prod",
			TenantID:    "tenant-001",
			Status:      "active",
			IsActive:    true,
			CreatedBy:   "system",
			Description: "DNS配置应用于DNS应用",
		},
	}

	if err := m.db.Create(&universalRelations).Error; err != nil {
		m.log.Error("创建通用关系失败", zap.Error(err))
		return err
	}

	m.log.Info("CMDB基础数据初始化完成")
	return nil
}

var menuData = `[
 {
    "id": 18,
    "parentId": 15,
    "path": "/access/admin",
    "title": "管理员账号",
    "name": "accessAdmin",
    "component": "/access/admin",
    "locale": "menu.access.admin"
  },
  {
    "id": 2,
    "parentId": 0,
    "title": "分析页",
    "icon": "DashboardOutlined",
    "component": "/dashboard/analysis",
    "path": "/dashboard/analysis",
    "name": "DashboardAnalysis",
    "keepAlive": true,
    "locale": "menu.dashboard.analysis",
    "weight": 2
  },
  {
    "id": 1,
    "parentId": 0,
    "title": "仪表盘",
    "icon": "DashboardOutlined",
    "component": "RouteView",
    "redirect": "/dashboard/analysis",
    "path": "/dashboard",
    "name": "Dashboard",
    "locale": "menu.dashboard"
  },
  {
    "id": 3,
    "parentId": 0,
    "title": "表单页",
    "icon": "FormOutlined",
    "component": "RouteView",
    "redirect": "/form/basic",
    "path": "/form",
    "name": "Form",
    "locale": "menu.form"
  },
  {
    "id": 5,
    "parentId": 0,
    "title": "链接",
    "icon": "LinkOutlined",
    "component": "RouteView",
    "redirect": "/link/iframe",
    "path": "/link",
    "name": "Link",
    "locale": "menu.link"
  },
  {
    "id": 6,
    "parentId": 5,
    "title": "AntDesign",
    "url": "https://ant.design/",
    "component": "Iframe",
    "path": "/link/iframe",
    "name": "LinkIframe",
    "keepAlive": true,
    "locale": "menu.link.iframe"
  },
  {
    "id": 7,
    "parentId": 5,
    "title": "AntDesignVue",
    "url": "https://antdv.com/",
    "component": "Iframe",
    "path": "/link/antdv",
    "name": "LinkAntdv",
    "keepAlive": true,
    "locale": "menu.link.antdv"
  },
  {
    "id": 8,
    "parentId": 5,
    "path": "https://www.baidu.com",
    "name": "LinkExternal",
    "title": "跳转百度",
    "locale": "menu.link.external"
  },
  {
    "id": 9,
    "parentId": 0,
    "title": "菜单",
    "icon": "BarsOutlined",
    "component": "RouteView",
    "path": "/menu",
    "redirect": "/menu/menu1",
    "name": "Menu",
    "locale": "menu.menu"
  },
  {
    "id": 10,
    "parentId": 9,
    "title": "菜单1",
    "component": "/menu/menu1",
    "path": "/menu/menu1",
    "name": "MenuMenu11",
    "keepAlive": true,
    "locale": "menu.menu.menu1"
  },
  {
    "id": 11,
    "parentId": 9,
    "title": "菜单2",
    "component": "/menu/menu2",
    "path": "/menu/menu2",
    "keepAlive": true,
    "locale": "menu.menu.menu2"
  },
  {
    "id": 12,
    "parentId": 9,
    "path": "/menu/menu3",
    "redirect": "/menu/menu3/menu1",
    "title": "菜单1-1",
    "component": "RouteView",
    "locale": "menu.menu.menu3"
  },
  {
    "id": 13,
    "parentId": 12,
    "path": "/menu/menu3/menu1",
    "component": "/menu/menu-1-1/menu1",
    "title": "菜单1-1-1",
    "keepAlive": true,
    "locale": "menu.menu3.menu1"
  },
  {
    "id": 14,
    "parentId": 12,
    "path": "/menu/menu3/menu2",
    "component": "/menu/menu-1-1/menu2",
    "title": "菜单1-1-2",
    "keepAlive": true,
    "locale": "menu.menu3.menu2"
  },
  {
    "id": 15,
    "path": "/access",
    "component": "RouteView",
    "redirect": "/access/common",
    "title": "权限模块",
    "name": "Access",
    "parentId": 0,
    "icon": "ClusterOutlined",
    "locale": "menu.access",
    "weight": 1
  },
  {
    "id": 51,
    "parentId": 15,
    "path": "/access/role",	
    "title": "角色管理",
    "name": "AccessRoles",
    "component": "/access/role",
    "locale": "menu.access.roles"
  },
{
    "id": 52,
    "parentId": 15,
    "path": "/access/menu",	
    "title": "菜单管理",
    "name": "AccessMenu",
    "component": "/access/menu",
    "locale": "menu.access.menus"
  },
{
    "id": 53,
    "parentId": 15,
    "path": "/access/api",	
    "title": "API管理",
    "name": "AccessAPI",
    "component": "/access/api",
    "locale": "menu.access.api"
  },
  {
    "id": 19,
    "parentId": 0,
    "title": "异常页",
    "icon": "WarningOutlined",
    "component": "RouteView",
    "redirect": "/exception/403",
    "path": "/exception",
    "name": "Exception",
    "locale": "menu.exception"
  },
  {
    "id": 20,
    "parentId": 19,
    "path": "/exception/403",
    "title": "403",
    "name": "403",
    "component": "/exception/403",
    "locale": "menu.exception.not-permission"
  },
  {
    "id": 21,
    "parentId": 19,
    "path": "/exception/404",
    "title": "404",
    "name": "404",
    "component": "/exception/404",
    "locale": "menu.exception.not-find"
  },
  {
    "id": 22,
    "parentId": 19,
    "path": "/exception/500",
    "title": "500",
    "name": "500",
    "component": "/exception/500",
    "locale": "menu.exception.server-error"
  },
  {
    "id": 23,
    "parentId": 0,
    "title": "结果页",
    "icon": "CheckCircleOutlined",
    "component": "RouteView",
    "redirect": "/result/success",
    "path": "/result",
    "name": "Result",
    "locale": "menu.result"
  },
  {
    "id": 24,
    "parentId": 23,
    "path": "/result/success",
    "title": "成功页",
    "name": "ResultSuccess",
    "component": "/result/success",
    "locale": "menu.result.success"
  },
  {
    "id": 25,
    "parentId": 23,
    "path": "/result/fail",
    "title": "失败页",
    "name": "ResultFail",
    "component": "/result/fail",
    "locale": "menu.result.fail"
  },
  {
    "id": 26,
    "parentId": 0,
    "title": "列表页",
    "icon": "TableOutlined",
    "component": "RouteView",
    "redirect": "/list/card-list",
    "path": "/list",
    "name": "List",
    "locale": "menu.list"
  },
  {
    "id": 27,
    "parentId": 26,
    "path": "/list/card-list",
    "title": "卡片列表",
    "name": "ListCard",
    "component": "/list/card-list",
    "locale": "menu.list.card-list"
  },
  {
    "id": 28,
    "parentId": 0,
    "title": "详情页",
    "icon": "ProfileOutlined",
    "component": "RouteView",
    "redirect": "/profile/basic",
    "path": "/profile",
    "name": "Profile",
    "locale": "menu.profile"
  },
  {
    "id": 29,
    "parentId": 28,
    "path": "/profile/basic",
    "title": "基础详情页",
    "name": "ProfileBasic",
    "component": "/profile/basic/index",
    "locale": "menu.profile.basic"
  },
  {
    "id": 30,
    "parentId": 26,
    "path": "/list/search-list",
    "title": "搜索列表",
    "name": "SearchList",
    "component": "/list/search-list",
    "locale": "menu.list.search-list"
  },
  {
    "id": 31,
    "parentId": 30,
    "path": "/list/search-list/articles",
    "title": "搜索列表（文章）",
    "name": "SearchListArticles",
    "component": "/list/search-list/articles",
    "locale": "menu.list.search-list.articles"
  },
  {
    "id": 32,
    "parentId": 30,
    "path": "/list/search-list/projects",
    "title": "搜索列表（项目）",
    "name": "SearchListProjects",
    "component": "/list/search-list/projects",
    "locale": "menu.list.search-list.projects"
  },
  {
    "id": 33,
    "parentId": 30,
    "path": "/list/search-list/applications",
    "title": "搜索列表（应用）",
    "name": "SearchListApplications",
    "component": "/list/search-list/applications",
    "locale": "menu.list.search-list.applications"
  },
  {
    "id": 34,
    "parentId": 26,
    "path": "/list/basic-list",
    "title": "标准列表",
    "name": "BasicCard",
    "component": "/list/basic-list",
    "locale": "menu.list.basic-list"
  },
  {
    "id": 35,
    "parentId": 28,
    "path": "/profile/advanced",
    "title": "高级详细页",
    "name": "ProfileAdvanced",
    "component": "/profile/advanced/index",
    "locale": "menu.profile.advanced"
  },
  {
    "id": 4,
    "parentId": 3,
    "title": "基础表单",
    "component": "/form/basic-form/index",
    "path": "/form/basic-form",
    "name": "FormBasic",
    "keepAlive": false,
    "locale": "menu.form.basic-form"
  },
  {
    "id": 36,
    "parentId": 0,
    "title": "个人页",
    "icon": "UserOutlined",
    "component": "RouteView",
    "redirect": "/account/center",
    "path": "/account",
    "name": "Account",
    "locale": "menu.account"
  },
  {
    "id": 37,
    "parentId": 36,
    "path": "/account/center",
    "title": "个人中心",
    "name": "AccountCenter",
    "component": "/account/center",
    "locale": "menu.account.center"
  },
  {
    "id": 38,
    "parentId": 36,
    "path": "/account/settings",
    "title": "个人设置",
    "name": "AccountSettings",
    "component": "/account/settings",
    "locale": "menu.account.settings"
  },
  {
    "id": 39,
    "parentId": 3,
    "title": "分步表单",
    "component": "/form/step-form/index",
    "path": "/form/step-form",
    "name": "FormStep",
    "keepAlive": false,
    "locale": "menu.form.step-form"
  },
  {
    "id": 40,
    "parentId": 3,
    "title": "高级表单",
    "component": "/form/advanced-form/index",
    "path": "/form/advanced-form",
    "name": "FormAdvanced",
    "keepAlive": false,
    "locale": "menu.form.advanced-form"
  },
  {
    "id": 41,
    "parentId": 26,
    "path": "/list/table-list",
    "title": "查询表格",
    "name": "ConsultTable",
    "component": "/list/table-list",
    "locale": "menu.list.consult-table"
  },
  {
    "id": 42,
    "parentId": 1,
    "title": "监控页",
    "component": "/dashboard/monitor",
    "path": "/dashboard/monitor",
    "name": "DashboardMonitor",
    "keepAlive": true,
    "locale": "menu.dashboard.monitor"
  },
  {
    "id": 43,
    "parentId": 1,
    "title": "工作台",
    "component": "/dashboard/workplace",
    "path": "/dashboard/workplace",
    "name": "DashboardWorkplace",
    "keepAlive": true,
    "locale": "menu.dashboard.workplace"
  },
  {
    "id": 44,
    "parentId": 26,
    "path": "/list/crud-table",
    "title": "增删改查表格",
    "name": "CrudTable",
    "component": "/list/crud-table",
    "locale": "menu.list.crud-table"
  },
  {
    "id": 45,
    "parentId": 9,
    "path": "/menu/menu4",
    "redirect": "/menu/menu4/menu1",
    "title": "菜单2-1",
    "component": "RouteView",
    "locale": "menu.menu.menu4"
  },
  {
    "id": 46,
    "parentId": 45,
    "path": "/menu/menu4/menu1",
    "component": "/menu/menu-2-1/menu1",
    "title": "菜单2-1-1",
    "keepAlive": true,
    "locale": "menu.menu4.menu1"
  },
  {
    "id": 47,
    "parentId": 45,
    "path": "/menu/menu4/menu2",
    "component": "/menu/menu-2-1/menu2",
    "title": "菜单2-1-2",
    "keepAlive": true,
    "locale": "menu.menu4.menu2"
  }
]`
