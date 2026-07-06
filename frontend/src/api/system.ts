import { request, requestPage } from "./http";
import type { PageQuery, PageResult } from "@/types/common";
import type {
  AdminUser,
  AdminUserForm,
  ApiPermission,
  Menu,
  Role,
  RoleForm,
  RolePermissions
} from "@/types/system";

export function getAdminUsersApi(params: PageQuery): Promise<PageResult<AdminUser>> {
  return requestPage<AdminUser>({ url: "/admin-users", method: "GET", params });
}

export function createAdminUserApi(data: AdminUserForm): Promise<AdminUser> {
  return request<AdminUser>({ url: "/admin-users", method: "POST", data });
}

export function updateAdminUserApi(id: number, data: AdminUserForm): Promise<AdminUser> {
  return request<AdminUser>({ url: `/admin-users/${id}`, method: "PUT", data });
}

export function resetAdminPasswordApi(id: number, password: string): Promise<void> {
  return request<void>({ url: `/admin-users/${id}/password`, method: "PATCH", data: { password } });
}

export function deleteAdminUserApi(id: number): Promise<void> {
  return request<void>({ url: `/admin-users/${id}`, method: "DELETE" });
}

export function getRolesApi(params: PageQuery): Promise<PageResult<Role>> {
  return requestPage<Role>({ url: "/roles", method: "GET", params });
}

export function createRoleApi(data: RoleForm): Promise<Role> {
  return request<Role>({ url: "/roles", method: "POST", data });
}

export function updateRoleApi(id: number, data: RoleForm): Promise<Role> {
  return request<Role>({ url: `/roles/${id}`, method: "PUT", data });
}

export function deleteRoleApi(id: number): Promise<void> {
  return request<void>({ url: `/roles/${id}`, method: "DELETE" });
}

export function getMenusTreeApi(): Promise<Menu[]> {
  return request<Menu[]>({ url: "/menus/tree", method: "GET" });
}

export function getApiPermissionsApi(): Promise<ApiPermission[]> {
  return request<ApiPermission[]>({ url: "/api-permissions", method: "GET" });
}

export function getRolePermissionsApi(id: number): Promise<RolePermissions> {
  return request<RolePermissions>({ url: `/roles/${id}/permissions`, method: "GET" });
}

export function updateRolePermissionsApi(id: number, data: RolePermissions): Promise<void> {
  return request<void>({ url: `/roles/${id}/permissions`, method: "PUT", data });
}
