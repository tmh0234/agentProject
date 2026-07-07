import type { StatusValue } from "./common";

export interface AdminUser {
  id: number;
  username: string;
  nickname: string;
  status: StatusValue;
  roles?: Role[];
  roleIds?: number[];
  lastLoginAt: string | null;
  createdAt: string;
  updatedAt: string;
}

export interface AdminUserForm {
  username: string;
  nickname: string;
  password?: string;
  status: StatusValue;
  roleIds: number[];
}

export interface Role {
  id: number;
  name: string;
  code: string;
  status: StatusValue;
  description: string;
  createdAt: string;
  updatedAt: string;
}

export interface RoleForm {
  name: string;
  code: string;
  status: StatusValue;
  description: string;
}

export interface Menu {
  id: number;
  parentId: number;
  title: string;
  routeName: string;
  path: string;
  component: string;
  icon: string;
  permission: string;
  sort: number;
  status: StatusValue;
  children?: Menu[];
}

export interface ApiPermission {
  id: number;
  method: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
  path: string;
  code: string;
  description: string;
}

export interface RolePermissions {
  menuIds: number[];
  apiPermissionIds: number[];
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  profile: AdminUser;
}
