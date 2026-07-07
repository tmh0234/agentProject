# 商城管理后台 Implementation Plan

> For agentic workers: use subagent-driven development or executing-plans to implement this plan task by task.

## Goal

Build a runnable standard mall admin system with Go Gin + GORM backend, Vue 3 admin frontend, MySQL 8 via Docker Compose, RBAC, and product SPU/SKU management.

## Architecture

The root directory owns integration-only configuration. `/backend` owns all Go server code, SQL initialization, API behavior, validation, tests, and backend startup docs. `/frontend` owns all Vue/TypeScript UI code, API clients, stores, routes, layout, pages, frontend tests, and build scripts.

## Tech Stack

- Backend: Go, Gin, GORM, MySQL driver, JWT, bcrypt.
- Frontend: Vue 3 Composition API, TypeScript, Vite, Pinia, Vue Router, Element Plus, TailwindCSS, Axios, Vitest.
- Database: MySQL 8 via Docker Compose.

## Source Documents

- Approved design: `docs/superpowers/specs/2026-07-06-mall-admin-design.md`
- Process rules: `AGENTS.md`
- Planner role rules: `agents/planner.toml`
- Frontend role rules: `agents/frontend.toml`
- Backend role rules: `agents/backend.toml`
- Reviewer role rules: `agents/reviewer.toml`

## API Contract Summary

All backend responses use:

```json
{
  "code": 0,
  "message": "成功",
  "data": {}
}
```

Paginated responses use:

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "items": [],
    "total": 0,
    "page": 1,
    "pageSize": 20
  }
}
```

Protected APIs require `Authorization: Bearer <token>`. Parameter errors return 400, missing login returns 401, missing permission returns 403, business conflicts return 409.

Required endpoints:

- `GET /api/v1/health`
- `POST /api/v1/auth/login`
- `GET /api/v1/auth/profile`
- `GET /api/v1/auth/menus`
- `GET /api/v1/admin-users`
- `POST /api/v1/admin-users`
- `PUT /api/v1/admin-users/:id`
- `PATCH /api/v1/admin-users/:id/password`
- `DELETE /api/v1/admin-users/:id`
- `GET /api/v1/roles`
- `POST /api/v1/roles`
- `PUT /api/v1/roles/:id`
- `DELETE /api/v1/roles/:id`
- `GET /api/v1/menus/tree`
- `GET /api/v1/api-permissions`
- `GET /api/v1/roles/:id/permissions`
- `PUT /api/v1/roles/:id/permissions`
- `GET /api/v1/product-categories`
- `POST /api/v1/product-categories`
- `PUT /api/v1/product-categories/:id`
- `DELETE /api/v1/product-categories/:id`
- `GET /api/v1/brands`
- `POST /api/v1/brands`
- `PUT /api/v1/brands/:id`
- `DELETE /api/v1/brands/:id`
- `GET /api/v1/products`
- `GET /api/v1/products/:id`
- `POST /api/v1/products`
- `PUT /api/v1/products/:id`
- `DELETE /api/v1/products/:id`
- `POST /api/v1/products/:id/skus`
- `PUT /api/v1/products/:id/skus/:skuId`
- `DELETE /api/v1/products/:id/skus/:skuId`
- `PATCH /api/v1/products/:id/status`

Route parameter naming convention:

- Product detail and nested SKU APIs use `:id` for product ID.
- Nested SKU APIs use `:skuId` for SKU ID.
- API permission seed data, backend routes, frontend clients, and documentation must use the same path shape.

## Root Integration Task

Ownership: Main Agent only.

Files:

- `docker-compose.yml`
- `README.md`
- `.gitignore`

Tasks:

- Add MySQL Compose service with MySQL 8, `mall_admin` database, `mall_user` account, utf8mb4, persistent volume, health check, and initialization mounted from `./backend/deploy/mysql/init`.
- Add README with local startup, verification commands, and configuration notes.
- Keep local private configuration out of version control.

## Backend Task

Ownership: Backend sub Agent only. Modify only `/backend`.

Files under `/backend`:

- `go.mod`, `go.sum`
- `config.example.yaml`
- `cmd/server/main.go`
- `deploy/mysql/init/001_schema.sql`
- `internal/config/config.go`
- `internal/database/database.go`
- `internal/model/*.go`
- `internal/repository/*.go`
- `internal/service/*.go`
- `internal/handler/*.go`
- `internal/middleware/*.go`
- `internal/router/router.go`
- `internal/response/response.go`
- `internal/testutil/*.go`
- focused `*_test.go` files for auth, RBAC, role permission transactions, product/SKU creation, and product status validation

Backend implementation requirements:

- Implement Gin router, config loading, GORM MySQL connection, models, repositories, services, handlers, middleware, SQL schema, and seed data.
- Use environment variables to override local config when needed.
- Keep API response fields consistent with this plan.
- Keep product nested SKU paths as `/products/:id/skus` and `/products/:id/skus/:skuId`.

Backend checks:

```bash
cd backend
gofmt -w .
go test ./...
```

## Frontend Task

Ownership: Frontend sub Agent only. Modify only `/frontend`.

Files under `/frontend`:

- `package.json`, `vite.config.ts`, `tsconfig*.json`, `tailwind.config.js`, `postcss.config.js`
- `index.html`
- `src/main.ts`, `src/App.vue`, `src/style.css`
- `src/api/*.ts`
- `src/types/*.ts`
- `src/stores/*.ts`
- `src/router/index.ts`
- `src/layout/*.vue`
- `src/views/Login.vue`, `src/views/Dashboard.vue`, `src/views/error/*.vue`
- `src/views/system/AdminUsers.vue`, `Roles.vue`, `Menus.vue`
- `src/views/product/Categories.vue`, `Brands.vue`, `Products.vue`
- focused Vitest tests for auth store, tags view store, money conversion, and SKU spec conversion

Frontend implementation requirements:

- Implement Vue 3 app, Element Plus admin layout, sidebar two-level menu, TagsView, auth/menu/router stores, typed API clients, login page, RBAC pages, and product pages.
- Keep product nested SKU frontend requests consistent with backend paths.
- Use real UTF-8 Chinese text in UI labels and messages.

Frontend checks:

```bash
cd frontend
npm test -- --run
npm run typecheck
npm run build
```

## Reviewer Task

Ownership: Reviewer sub Agent. Do not modify business code unless explicitly asked later.

Review checklist:

- Compare backend endpoints, SQL `api_permissions`, frontend API clients, and TypeScript types against this plan.
- Inspect password hashing, JWT validation, RBAC enforcement, transaction usage, product/SKU validation, frontend auth guard, and 401 handling.
- Scan for BOM and Chinese text escaping problems.
- Report P1/P2/P3 findings. P1 blocks release, P2 should be fixed unless repeated three times, P3 is backlog.

## Integration Verification

Run after Reviewer has no P1 and at most two P2 findings:

```bash
docker compose up -d mysql
cd backend
go test ./...
go run ./cmd/server
```

In another terminal:

```bash
cd frontend
npm test -- --run
npm run typecheck
npm run build
npm run dev
```

Manual/API flow:

1. `GET http://localhost:8080/api/v1/health` returns success.
2. Login with the seeded administrator account from the initialization SQL.
3. Fetch profile and menus with the returned token.
4. Confirm frontend login renders dashboard and left two-level menu.
5. Open 管理员管理、角色管理、商品分类、品牌管理、商品列表 pages through menu/tabs.
6. Create a product with two SKUs.
7. Publish and unpublish the product.

## Self-Review

- Spec coverage includes Docker Compose MySQL, backend RBAC/product APIs, frontend layout/pages/stores, review loop, and integration verification.
- Placeholder scan: no pending placeholder markers.
- Type consistency: API names and JSON field names use camelCase externally and snake_case only in SQL.
- Route consistency: nested SKU routes consistently use `/products/:id/skus` and `/products/:id/skus/:skuId`.
