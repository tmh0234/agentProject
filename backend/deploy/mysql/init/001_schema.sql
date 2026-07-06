CREATE DATABASE IF NOT EXISTS mall_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE mall_admin;
SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS admin_users (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  username VARCHAR(64) NOT NULL,
  nickname VARCHAR(64) NOT NULL,
  password_hash VARCHAR(128) NOT NULL,
  status TINYINT NOT NULL DEFAULT 1,
  last_login_at DATETIME NULL,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  deleted_at DATETIME(3) NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_admin_users_username (username),
  KEY idx_admin_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS roles (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL,
  code VARCHAR(64) NOT NULL,
  status TINYINT NOT NULL DEFAULT 1,
  description VARCHAR(255) NOT NULL DEFAULT '',
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  deleted_at DATETIME(3) NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_roles_code (code),
  KEY idx_roles_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS menus (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  parent_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  title VARCHAR(64) NOT NULL,
  route_name VARCHAR(64) NOT NULL DEFAULT '',
  path VARCHAR(128) NOT NULL DEFAULT '',
  component VARCHAR(128) NOT NULL DEFAULT '',
  icon VARCHAR(64) NOT NULL DEFAULT '',
  permission VARCHAR(128) NOT NULL DEFAULT '',
  sort INT NOT NULL DEFAULT 0,
  status TINYINT NOT NULL DEFAULT 1,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  deleted_at DATETIME(3) NULL,
  PRIMARY KEY (id),
  KEY idx_menus_parent_id (parent_id),
  KEY idx_menus_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS api_permissions (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  method VARCHAR(16) NOT NULL,
  path VARCHAR(255) NOT NULL,
  code VARCHAR(128) NOT NULL,
  description VARCHAR(255) NOT NULL DEFAULT '',
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  deleted_at DATETIME(3) NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_api_permissions_method_path (method, path),
  UNIQUE KEY uk_api_permissions_code (code),
  KEY idx_api_permissions_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS admin_user_roles (
  admin_user_id BIGINT UNSIGNED NOT NULL,
  role_id BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (admin_user_id, role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS role_menus (
  role_id BIGINT UNSIGNED NOT NULL,
  menu_id BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (role_id, menu_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS role_api_permissions (
  role_id BIGINT UNSIGNED NOT NULL,
  api_permission_id BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (role_id, api_permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS product_categories (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  parent_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  name VARCHAR(64) NOT NULL,
  status TINYINT NOT NULL DEFAULT 1,
  sort INT NOT NULL DEFAULT 0,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  deleted_at DATETIME(3) NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_product_categories_name (name),
  KEY idx_product_categories_parent_id (parent_id),
  KEY idx_product_categories_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS brands (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL,
  logo VARCHAR(255) NOT NULL DEFAULT '',
  status TINYINT NOT NULL DEFAULT 1,
  sort INT NOT NULL DEFAULT 0,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  deleted_at DATETIME(3) NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_brands_name (name),
  KEY idx_brands_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS products (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(128) NOT NULL,
  code VARCHAR(64) NOT NULL,
  category_id BIGINT UNSIGNED NOT NULL,
  brand_id BIGINT UNSIGNED NOT NULL,
  main_image VARCHAR(255) NOT NULL DEFAULT '',
  min_price_cents BIGINT NOT NULL DEFAULT 0,
  max_price_cents BIGINT NOT NULL DEFAULT 0,
  status TINYINT NOT NULL DEFAULT 0,
  description TEXT NULL,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  deleted_at DATETIME(3) NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_products_code (code),
  KEY idx_products_category_id (category_id),
  KEY idx_products_brand_id (brand_id),
  KEY idx_products_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS product_skus (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  product_id BIGINT UNSIGNED NOT NULL,
  code VARCHAR(64) NOT NULL,
  specs JSON NOT NULL,
  price_cents BIGINT NOT NULL,
  stock INT NOT NULL DEFAULT 0,
  status TINYINT NOT NULL DEFAULT 1,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  deleted_at DATETIME(3) NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_product_skus_code (code),
  KEY idx_product_skus_product_id (product_id),
  KEY idx_product_skus_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO roles (id, name, code, status, description, created_at, updated_at)
VALUES (1, '超级管理员', 'super_admin', 1, '拥有全部后台权限', NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE name = VALUES(name), status = VALUES(status), description = VALUES(description), updated_at = NOW(3);

INSERT INTO admin_users (id, username, nickname, password_hash, status, created_at, updated_at)
VALUES (1, 'admin', '超级管理员', '$2a$10$E.vZUD/NkdYZ9xxP8SAk/.4F0WGs1l3jPxwo4ra3J/a5Ti8vErqgu', 1, NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE nickname = VALUES(nickname), password_hash = VALUES(password_hash), status = VALUES(status), updated_at = NOW(3);

INSERT INTO admin_user_roles (admin_user_id, role_id) VALUES (1, 1)
ON DUPLICATE KEY UPDATE role_id = VALUES(role_id);

INSERT INTO menus (id, parent_id, title, route_name, path, component, icon, permission, sort, status, created_at, updated_at) VALUES
(1, 0, '仪表盘', 'Dashboard', '/dashboard', 'views/Dashboard.vue', 'Gauge', 'dashboard', 1, 1, NOW(3), NOW(3)),
(2, 0, '系统管理', 'System', '/system', 'layout', 'Settings', 'system', 10, 1, NOW(3), NOW(3)),
(3, 2, '管理员管理', 'AdminUsers', '/system/admin-users', 'views/system/AdminUsers.vue', 'UserCog', 'admin-user:list', 11, 1, NOW(3), NOW(3)),
(4, 2, '角色管理', 'Roles', '/system/roles', 'views/system/Roles.vue', 'ShieldCheck', 'role:list', 12, 1, NOW(3), NOW(3)),
(5, 2, '菜单管理', 'Menus', '/system/menus', 'views/system/Menus.vue', 'Menu', 'menu:list', 13, 1, NOW(3), NOW(3)),
(6, 0, '商品管理', 'Product', '/product', 'layout', 'Package', 'product', 20, 1, NOW(3), NOW(3)),
(7, 6, '商品分类', 'ProductCategories', '/product/categories', 'views/product/Categories.vue', 'ListTree', 'category:list', 21, 1, NOW(3), NOW(3)),
(8, 6, '品牌管理', 'Brands', '/product/brands', 'views/product/Brands.vue', 'Badge', 'brand:list', 22, 1, NOW(3), NOW(3)),
(9, 6, '商品列表', 'Products', '/product/products', 'views/product/Products.vue', 'Boxes', 'product:list', 23, 1, NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE parent_id = VALUES(parent_id), title = VALUES(title), route_name = VALUES(route_name), path = VALUES(path), component = VALUES(component), icon = VALUES(icon), permission = VALUES(permission), sort = VALUES(sort), status = VALUES(status), updated_at = NOW(3);

INSERT INTO api_permissions (id, method, path, code, description, created_at, updated_at) VALUES
(1, 'GET', '/api/v1/admin-users', 'admin-user:list', '管理员列表', NOW(3), NOW(3)),
(2, 'POST', '/api/v1/admin-users', 'admin-user:create', '创建管理员', NOW(3), NOW(3)),
(3, 'PUT', '/api/v1/admin-users/:id', 'admin-user:update', '更新管理员', NOW(3), NOW(3)),
(4, 'PATCH', '/api/v1/admin-users/:id/password', 'admin-user:password', '重置管理员密码', NOW(3), NOW(3)),
(5, 'DELETE', '/api/v1/admin-users/:id', 'admin-user:delete', '删除管理员', NOW(3), NOW(3)),
(6, 'GET', '/api/v1/roles', 'role:list', '角色列表', NOW(3), NOW(3)),
(7, 'POST', '/api/v1/roles', 'role:create', '创建角色', NOW(3), NOW(3)),
(8, 'PUT', '/api/v1/roles/:id', 'role:update', '更新角色', NOW(3), NOW(3)),
(9, 'DELETE', '/api/v1/roles/:id', 'role:delete', '删除角色', NOW(3), NOW(3)),
(10, 'GET', '/api/v1/menus/tree', 'menu:list', '菜单树', NOW(3), NOW(3)),
(11, 'GET', '/api/v1/api-permissions', 'api-permission:list', '接口权限列表', NOW(3), NOW(3)),
(12, 'GET', '/api/v1/roles/:id/permissions', 'role:permission-detail', '角色权限详情', NOW(3), NOW(3)),
(13, 'PUT', '/api/v1/roles/:id/permissions', 'role:permission-update', '更新角色权限', NOW(3), NOW(3)),
(14, 'GET', '/api/v1/product-categories', 'category:list', '商品分类列表', NOW(3), NOW(3)),
(15, 'POST', '/api/v1/product-categories', 'category:create', '创建商品分类', NOW(3), NOW(3)),
(16, 'PUT', '/api/v1/product-categories/:id', 'category:update', '更新商品分类', NOW(3), NOW(3)),
(17, 'DELETE', '/api/v1/product-categories/:id', 'category:delete', '删除商品分类', NOW(3), NOW(3)),
(18, 'GET', '/api/v1/brands', 'brand:list', '品牌列表', NOW(3), NOW(3)),
(19, 'POST', '/api/v1/brands', 'brand:create', '创建品牌', NOW(3), NOW(3)),
(20, 'PUT', '/api/v1/brands/:id', 'brand:update', '更新品牌', NOW(3), NOW(3)),
(21, 'DELETE', '/api/v1/brands/:id', 'brand:delete', '删除品牌', NOW(3), NOW(3)),
(22, 'GET', '/api/v1/products', 'product:list', '商品列表', NOW(3), NOW(3)),
(23, 'GET', '/api/v1/products/:id', 'product:detail', '商品详情', NOW(3), NOW(3)),
(24, 'POST', '/api/v1/products', 'product:create', '创建商品', NOW(3), NOW(3)),
(25, 'PUT', '/api/v1/products/:id', 'product:update', '更新商品', NOW(3), NOW(3)),
(26, 'DELETE', '/api/v1/products/:id', 'product:delete', '删除商品', NOW(3), NOW(3)),
(27, 'POST', '/api/v1/products/:id/skus', 'product-sku:create', '创建 SKU', NOW(3), NOW(3)),
(28, 'PUT', '/api/v1/products/:id/skus/:skuId', 'product-sku:update', '更新 SKU', NOW(3), NOW(3)),
(29, 'DELETE', '/api/v1/products/:id/skus/:skuId', 'product-sku:delete', '删除 SKU', NOW(3), NOW(3)),
(30, 'PATCH', '/api/v1/products/:id/status', 'product:status', '商品上下架', NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE method = VALUES(method), path = VALUES(path), code = VALUES(code), description = VALUES(description), updated_at = NOW(3);

INSERT INTO role_menus (role_id, menu_id)
SELECT 1, id FROM menus
ON DUPLICATE KEY UPDATE menu_id = VALUES(menu_id);

INSERT INTO role_api_permissions (role_id, api_permission_id)
SELECT 1, id FROM api_permissions
ON DUPLICATE KEY UPDATE api_permission_id = VALUES(api_permission_id);

INSERT INTO product_categories (id, parent_id, name, status, sort, created_at, updated_at) VALUES
(1, 0, '手机数码', 1, 1, NOW(3), NOW(3)),
(2, 0, '家用电器', 1, 2, NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE parent_id = VALUES(parent_id), name = VALUES(name), status = VALUES(status), sort = VALUES(sort), updated_at = NOW(3);

INSERT INTO brands (id, name, logo, status, sort, created_at, updated_at) VALUES
(1, 'Acme', '', 1, 1, NOW(3), NOW(3)),
(2, 'Northwind', '', 1, 2, NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE name = VALUES(name), logo = VALUES(logo), status = VALUES(status), sort = VALUES(sort), updated_at = NOW(3);
