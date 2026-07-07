# 商城管理后台设计规格

日期：2026-07-06

## 目标

从当前空仓库搭建一个标准工程版商城管理后台，包含权限管理员模块和商品模块。系统采用前后端分离架构，后端使用 Go Gin + GORM，前端使用 Vue 3 + TypeScript + Pinia + Vue Router + Element Plus + TailwindCSS，数据库使用 Docker Compose 管理的 MySQL 8。

第一版目标是可本机稳定运行、可联调、可继续扩展，不追求生产级数据库高可用。

## 范围

### 包含

- 登录与 JWT 鉴权。
- 标准 RBAC：管理员、角色、菜单、API 权限、角色授权。
- 商品模块：商品分类、品牌、商品 SPU、SKU、价格、库存、上下架。
- 管理后台主布局：左侧二级菜单、HeaderBar、TagsView、内容区。
- Docker Compose MySQL 部署、初始化 SQL、后端配置、前端代理。
- 基础单元测试和集成验证脚本或命令说明。

### 不包含

- 生产级数据库主从、高可用、自动备份和监控告警。
- 支付、订单、会员、营销活动、物流、售后。
- 复杂商品属性模板、批量导入、图库管理、库存流水。
- 多租户、部门数据权限、操作日志和登录日志。

## 工程结构

```text
.
├── docker-compose.yml
├── backend
│   ├── cmd/server
│   ├── config.yaml
│   ├── deploy/mysql/init
│   ├── internal
│   │   ├── config
│   │   ├── handler
│   │   ├── middleware
│   │   ├── model
│   │   ├── repository
│   │   ├── router
│   │   └── service
│   └── pkg
└── frontend
    ├── src
    │   ├── api
    │   ├── components
    │   ├── layout
    │   ├── router
    │   ├── stores
    │   ├── types
    │   └── views
    └── vite.config.ts
```

## 数据库部署

根目录提供 `docker-compose.yml`，启动 `mysql:8.0`。数据库通过 Docker volume 持久化，容器重启不会丢失数据。

默认配置：

- 数据库：`mall_admin`
- 用户：`mall_user`
- 字符集：`utf8mb4`
- 排序规则：`utf8mb4_unicode_ci`
- 初始化目录：`backend/deploy/mysql/init`

首次启动时，MySQL 自动执行初始化 SQL，创建表结构和种子数据。后端 `config.yaml` 固定连接本机映射端口的 MySQL。

## 数据模型

### 权限管理员

- `admin_users`：后台管理员账号，包含用户名、昵称、密码哈希、状态、最后登录时间。
- `roles`：角色，包含角色名、编码、状态、描述。
- `menus`：前端菜单和按钮权限，支持父子层级、路由路径、组件路径、图标、排序。
- `api_permissions`：后端 API 权限，包含方法、路径、权限编码、描述。
- `admin_user_roles`：管理员与角色多对多关系。
- `role_menus`：角色与菜单多对多关系。
- `role_api_permissions`：角色与 API 权限多对多关系。

### 商品

- `product_categories`：商品分类，支持父子层级、状态、排序。
- `brands`：品牌，包含名称、Logo、状态、排序。
- `products`：商品 SPU，包含商品名称、编码、分类、品牌、主图、售价范围、上下架状态、描述。
- `product_skus`：商品 SKU，包含规格 JSON、SKU 编码、价格、库存、状态。

金额字段使用整数分保存，避免浮点误差。软删除使用 GORM `DeletedAt`。

## 后端设计

后端必须遵循 `Handler -> Service -> Repository` 分层：

- Handler：参数绑定、validator 校验、统一响应。
- Service：业务规则、密码哈希、JWT、RBAC 权限判断、事务。
- Repository：GORM 查询、分页、持久化。

核心接口分组：

- `POST /api/v1/auth/login`
- `GET /api/v1/auth/profile`
- `GET /api/v1/auth/menus`
- `GET/POST/PUT/DELETE /api/v1/admin-users`
- `GET/POST/PUT/DELETE /api/v1/roles`
- `GET /api/v1/menus/tree`
- `PUT /api/v1/roles/:id/permissions`
- `GET/POST/PUT/DELETE /api/v1/product-categories`
- `GET/POST/PUT/DELETE /api/v1/brands`
- `GET/POST/PUT/DELETE /api/v1/products`
- `PATCH /api/v1/products/:id/status`

删除接口优先使用软删除。列表接口统一支持分页、关键词、状态筛选。

## 前端设计

前端必须使用 Vue 3 Composition API + TypeScript + Pinia + Vue Router + Element Plus + TailwindCSS。

布局要求：

- 左侧 Sidebar 使用 Element Plus `el-sub-menu` 实现二级菜单和手风琴模式。
- Sidebar 展开宽度 165px，折叠宽度 48px，容器用 `transition: width`。
- 右侧包含 HeaderBar、TagsView、router-view。
- TagsView 使用 Pinia 维护打开标签，支持切换、关闭当前、关闭其他，从 localStorage 恢复。
- 前端中文文案必须使用真实 UTF-8 中文字符，禁止 Unicode 转义序列。

页面：

- 登录页。
- 仪表盘占位页。
- 管理员管理：列表、筛选、新增、编辑、禁用、重置密码、分配角色。
- 角色管理：列表、编辑、授权菜单/API 权限。
- 菜单管理：树形列表。
- 商品分类：树形或表格列表、新增、编辑、启停。
- 品牌管理：列表、新增、编辑、启停。
- 商品列表：筛选、新增、编辑、SKU 行编辑、上下架。

## 错误处理

- 后端统一响应结构：`code`、`message`、`data`。
- 参数校验错误返回 400。
- 未登录返回 401。
- 无权限返回 403。
- 业务冲突返回 409，例如用户名重复、角色编码重复、SKU 编码重复。
- 前端统一拦截 401 并跳转登录页。

## 测试与验证

后端：

- Service 层测试登录、密码校验、RBAC 判断、商品创建和 SKU 校验。
- Handler 或集成测试覆盖核心接口。

前端：

- 类型检查和构建必须通过。
- 关键 store 行为测试：登录状态、TagsView、动态菜单。

集成验证：

1. `docker compose up -d mysql`
2. 启动后端并确认健康检查通过。
3. 启动前端并确认 `/api` 代理可用。
4. 使用默认账号 `admin / Admin@123456` 登录。
5. 验证菜单加载、管理员列表、角色列表、商品列表、新增商品、上下架。

## 实施流程

必须遵守仓库 `AGENTS.md`：

1. Planner 输出数据库 SQL、RESTful API 契约、前后端原子任务。
2. Frontend 子 Agent 只修改 `/frontend`。
3. Backend 子 Agent 只修改 `/backend`。
4. Reviewer 审查契约一致性、质量、测试覆盖和 Unicode 转义序列。
5. 主 Agent 只负责协调、根目录配置、集成验证，以及不归属子 Agent 的问题。

## 自审结果

- 无未定事项或待补充占位。
- 方案与用户确认的“方案 2、标准 RBAC、标准商品 + SKU”一致。
- 第一版范围明确排除了订单、支付、会员、营销和生产级数据库高可用。
- 前后端边界与 AGENTS.md 分工一致。
