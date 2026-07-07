import { createPinia, setActivePinia } from "pinia";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { useAuthStore } from "@/stores/auth";
import type { AdminUser } from "@/types/system";

vi.mock("@/api/auth", () => ({
  loginApi: vi.fn(),
  getProfileApi: vi.fn()
}));

const { loginApi } = await import("@/api/auth");
const mockedLoginApi = vi.mocked(loginApi);

describe("auth store", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it("persists token and profile to localStorage", () => {
    const auth = useAuthStore();
    const profile: AdminUser = {
      id: 1,
      username: "admin",
      nickname: "超级管理员",
      status: 1,
      roles: [{ id: 1, name: "超级管理员", code: "super_admin", status: 1, description: "", createdAt: "", updatedAt: "" }],
      lastLoginAt: "2026-07-06T10:00:00+08:00",
      createdAt: "2026-07-06T10:00:00+08:00",
      updatedAt: "2026-07-06T10:00:00+08:00"
    };

    auth.setSession("jwt-token", profile);

    expect(auth.token).toBe("jwt-token");
    expect(auth.profile?.nickname).toBe("超级管理员");
    expect(localStorage.getItem("mall_admin_token")).toBe("jwt-token");

    setActivePinia(createPinia());
    const restored = useAuthStore();

    expect(restored.token).toBe("jwt-token");
    expect(restored.profile?.username).toBe("admin");
  });

  it("stores login response profile from backend contract", async () => {
    const auth = useAuthStore();
    mockedLoginApi.mockResolvedValueOnce({
      token: "jwt-token",
      profile: {
        id: 1,
        username: "admin",
        nickname: "超级管理员",
        status: 1,
        roles: [],
        lastLoginAt: null,
        createdAt: "2026-07-06T10:00:00+08:00",
        updatedAt: "2026-07-06T10:00:00+08:00"
      }
    });

    await auth.login({ username: "admin", password: "Admin@123456" });

    expect(auth.token).toBe("jwt-token");
    expect(auth.profile?.username).toBe("admin");
    expect(localStorage.getItem("mall_admin_profile")).toContain("超级管理员");
  });
});
