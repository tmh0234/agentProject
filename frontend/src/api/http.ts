import axios, { AxiosError, type AxiosRequestConfig } from "axios";
import { ElMessage } from "element-plus";
import type { ApiResponse, PageResult } from "@/types/common";

const TOKEN_KEY = "mall_admin_token";

export const http = axios.create({
  baseURL: "/api/v1",
  timeout: 12000
});

http.interceptors.request.use((config) => {
  const token = localStorage.getItem(TOKEN_KEY);
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

http.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ApiResponse<unknown>>) => {
    if (error.response?.status === 401) {
      localStorage.removeItem(TOKEN_KEY);
      localStorage.removeItem("mall_admin_profile");
      const redirect = encodeURIComponent(window.location.pathname + window.location.search);
      window.location.href = `/login?redirect=${redirect}`;
      return Promise.reject(error);
    }

    const message = error.response?.data?.message || error.message || "网络请求失败";
    ElMessage.error(message);
    return Promise.reject(error);
  }
);

function unwrapResponse<T>(data: unknown): T {
  const body = data as ApiResponse<T>;
  if (body.code !== 0) {
    ElMessage.error(body.message || "请求失败");
    throw new Error(body.message || "请求失败");
  }
  return body.data;
}

export async function request<T>(config: AxiosRequestConfig): Promise<T> {
  const response = await http.request<ApiResponse<T>>(config);
  return unwrapResponse<T>(response.data);
}

export async function requestPage<T>(config: AxiosRequestConfig): Promise<PageResult<T>> {
  return request<PageResult<T>>(config);
}
