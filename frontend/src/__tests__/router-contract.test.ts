import { describe, expect, it } from "vitest";
import { menuRouteName, resolveViewModulePath } from "@/router/menuRoutes";
import type { Menu } from "@/types/system";

describe("menu route contract", () => {
  it("uses backend routeName as route identity", () => {
    const menu: Menu = {
      id: 1,
      parentId: 0,
      title: "管理员管理",
      routeName: "AdminUsers",
      path: "/system/admin-users",
      component: "views/system/AdminUsers.vue",
      icon: "User",
      permission: "system:admin-users:list",
      sort: 1,
      status: 1
    };

    expect(menuRouteName(menu)).toBe("AdminUsers");
  });

  it("maps backend views component path to glob module path", () => {
    expect(resolveViewModulePath("views/system/AdminUsers.vue")).toBe("/src/views/system/AdminUsers.vue");
    expect(resolveViewModulePath("system/Roles")).toBe("/src/views/system/Roles.vue");
  });

  it("registers dynamic leaf routes with backend route names", async () => {
    const { router } = await import("@/router");

    expect(router.hasRoute("AdminUsers")).toBe(true);
    expect(router.resolve({ name: "AdminUsers" }).path).toBe("/system/admin-users");
  });
});
