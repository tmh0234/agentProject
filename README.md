# 商城管理后台

这是一个前后端分离的商城管理后台初版工程，当前分支为 `codex/init-admin-project`。

## 技术栈

- 后端：Go 1.22、Gin、GORM、JWT、bcrypt、MySQL 8。
- 前端：Vue 3、TypeScript、Vite、Pinia、Vue Router、Element Plus、TailwindCSS、Axios、Vitest。
- 数据库：Docker Compose 启动 MySQL 8，并通过 `backend/deploy/mysql/init/001_schema.sql` 初始化表结构和种子数据。

## 工程结构

```text
.
├── AGENTS.md
├── agents/
├── backend/
├── frontend/
├── docker-compose.yml
└── docs/superpowers/
```

## 本地启动

### 1. 切换分支

```bash
git checkout codex/init-admin-project
```

### 2. 启动数据库

项目默认使用根目录 `docker-compose.yml` 中的 MySQL 服务，后端 `backend/config.yaml` 已和该服务保持一致。

```bash
docker compose up -d mysql
```

### 3. 启动后端

```bash
cd backend
go mod download
go run ./cmd/server
```

健康检查：

```bash
curl http://localhost:8080/api/v1/health
```

### 4. 启动前端

```bash
cd frontend
npm install
npm run dev
```

默认前端地址：

```text
http://localhost:5173
```

## 验证命令

后端：

```bash
cd backend
gofmt -w .
go test ./...
```

前端：

```bash
cd frontend
npm test -- --run
npm run typecheck
npm run build
```

## Agent 协作规则

本仓库使用 `AGENTS.md` 和 `agents/*.toml` 约束多 Agent 协作：

- Planner 只负责规划和契约，不直接改代码。
- Frontend 只修改 `/frontend`。
- Backend 只修改 `/backend`。
- Reviewer 只做审查和报告，默认不直接修改业务代码。
- 主 Agent 负责根目录配置、集成验证和跨模块协调。

所有源码和配置必须使用 UTF-8，无 BOM。中文文案必须保存为真实中文字符，禁止把中文写成 Unicode 转义序列。
