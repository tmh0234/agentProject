import { request } from "./http";
import type { LoginRequest, LoginResponse, Menu, AdminUser } from "@/types/system";

export function loginApi(data: LoginRequest): Promise<LoginResponse> {
  return request<LoginResponse>({
    url: "/auth/login",
    method: "POST",
    data
  });
}

export function getProfileApi(): Promise<AdminUser> {
  return request<AdminUser>({
    url: "/auth/profile",
    method: "GET"
  });
}

export function getAuthMenusApi(): Promise<Menu[]> {
  return request<Menu[]>({
    url: "/auth/menus",
    method: "GET"
  });
}
