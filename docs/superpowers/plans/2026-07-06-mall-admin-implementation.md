# 商城管理后台 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a runnable standard mall admin system with Go Gin + GORM backend, Vue 3 admin frontend, MySQL 8 via Docker Compose, RBAC, and product SPU/SKU management.

**Architecture:** The root directory owns integration-only configuration. `/backend` owns all Go server code, SQL initialization, API behavior, validation, tests, and backend startup docs. `/frontend` owns all Vue/TypeScript UI code, API clients, stores, routes, layout, pages, frontend tests, and build scripts.

**Tech Stack:** Backend: Go, Gin, GORM, MySQL driver, validator, JWT, bcrypt. Frontend: Vue 3 Composition API, TypeScript, Vite, Pinia, Vue Router, Element Plus, TailwindCSS, Axios, Vitest. Database: MySQL 8 via Docker Compose.

---

## Source Documents

- Approved design: `docs/superpowers/specs/2026-07-06-mall-admin-design.md`
- Process rules: `AGENTS.md`
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
- `POST /api/v1/products/:productId/skus`
- `PUT /api/v1/products/:productId/skus/:skuId`
- `DELETE /api/v1/products/:productId/skus/:skuId`
- `PATCH /api/v1/products/:id/status`
- `GET /api/v1/health`

Default login for integration verification:

```text
admin / Admin@123456
```

## Root Integration Task

**Files:**

- Create: `docker-compose.yml`

- [ ] **Step 1: Add MySQL Compose service**

Create a root `docker-compose.yml` with MySQL 8, `mall_admin` database, `mall_user` account, password `mall_password`, `utf8mb4`, a persistent volume, health check, and initialization mounted from `./backend/deploy/mysql/init`.

- [ ] **Step 2: Verify Compose syntax**

Run:

```bash
docker compose config
```

Expected: the command prints normalized YAML and exits successfully.

## Backend Task

**Ownership:** Backend sub Agent only. Modify only `/backend`.

**Files to create under `/backend`:**

- `go.mod`, `go.sum`
- `config.yaml`
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

- [ ] **Step 1: Write failing tests for core service behavior**

Add tests first for password login, JWT/profile flow, RBAC allow/deny behavior, transactional role permission update, product create with SKU price range sync, SKU uniqueness, and product publish validation requiring at least one enabled SKU.

- [ ] **Step 2: Run tests and confirm red**

Run:

```bash
cd backend
go test ./...
```

Expected: tests fail because implementation does not exist yet.

- [ ] **Step 3: Implement backend**

Implement Gin router, config loading, GORM MySQL connection, models, repositories, services, handlers, middleware, SQL schema, seed data with a real bcrypt hash for `Admin@123456`, and all required endpoints.

- [ ] **Step 4: Run backend tests and formatting**

Run:

```bash
cd backend
gofmt -w .
go test ./...
```

Expected: all tests pass.

- [ ] **Step 5: Backend self-review**

Check that only `/backend` changed, Chinese text is UTF-8 characters, no Unicode escape sequences were introduced, and all API fields match this plan.

## Frontend Task

**Ownership:** Frontend sub Agent only. Modify only `/frontend`.

**Files to create under `/frontend`:**

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

- [ ] **Step 1: Write failing tests for store/helper behavior**

Add tests first for auth token persistence, TagsView close-current and close-other behavior, yuan-to-cents conversion, cents-to-yuan display conversion, and SKU spec key-value conversion.

- [ ] **Step 2: Run tests and confirm red**

Run:

```bash
cd frontend
npm test -- --run
```

Expected: tests fail because implementation does not exist yet.

- [ ] **Step 3: Implement frontend**

Implement Vue 3 app, Element Plus admin layout, sidebar two-level menu, TagsView, auth/menu/router stores, typed API clients, login page, RBAC pages, and product pages. Use real UTF-8 Chinese text in UI labels and messages.

- [ ] **Step 4: Run frontend checks**

Run:

```bash
cd frontend
npm test -- --run
npm run typecheck
npm run build
```

Expected: all checks pass.

- [ ] **Step 5: Frontend self-review**

Check that only `/frontend` changed, no `any` type was introduced, no Unicode escape sequences were introduced, and all API fields match this plan.

## Reviewer Task

**Ownership:** Reviewer sub Agent. Do not modify business code unless explicitly asked later.

- [ ] **Step 1: Review contract consistency**

Compare `/backend` endpoint response fields and `/frontend/src/types` plus API clients against this plan.

- [ ] **Step 2: Review quality and tests**

Inspect security-sensitive behavior: password hashing, JWT validation, RBAC enforcement, transaction usage, product/SKU validation, frontend auth guard, and 401 handling.

- [ ] **Step 3: Unicode escape scan**

Run:

```bash
rg -F "\\u" frontend backend docs -n
```

Expected: no source file contains Unicode escape sequences for Chinese UI/business text. Mentions in process docs are acceptable only if they are explicitly describing the forbidden pattern.

- [ ] **Step 4: Report severity**

Report P1/P2/P3 findings. P1 blocks release, P2 should be fixed unless repeated three times, P3 is backlog.

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
2. Login with `admin / Admin@123456`.
3. Fetch profile and menus with the returned token.
4. Confirm frontend login renders dashboard and left two-level menu.
5. Open 관리자管理, 角色管理, 商品分类, 品牌管理, 商品列表 pages through menu/tabs.
6. Create a product with two SKUs.
7. Publish and unpublish the product.

## Self-Review

- Spec coverage: includes Docker Compose MySQL, backend RBAC/product APIs, frontend layout/pages/stores, review loop, and integration verification.
- Placeholder scan: no pending placeholder markers.
- Type consistency: API names and JSON field names use camelCase externally and snake_case only in SQL.
