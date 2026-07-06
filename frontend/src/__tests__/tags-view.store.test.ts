import { createPinia, setActivePinia } from "pinia";
import { beforeEach, describe, expect, it } from "vitest";
import { useTagsViewStore } from "@/stores/tagsView";

describe("tags view store", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it("closes current tag and returns adjacent route", () => {
    const tags = useTagsViewStore();
    tags.addTag({ path: "/dashboard", title: "仪表盘", name: "Dashboard" });
    tags.addTag({ path: "/system/admin-users", title: "管理员管理", name: "AdminUsers" });
    tags.addTag({ path: "/product/products", title: "商品列表", name: "Products" });

    const next = tags.closeCurrent("/system/admin-users");

    expect(tags.visitedTags.map((tag) => tag.path)).toEqual(["/dashboard", "/product/products"]);
    expect(next?.path).toBe("/product/products");
  });

  it("keeps current tag when closing other tags", () => {
    const tags = useTagsViewStore();
    tags.addTag({ path: "/dashboard", title: "仪表盘", name: "Dashboard" });
    tags.addTag({ path: "/system/roles", title: "角色管理", name: "Roles" });
    tags.addTag({ path: "/product/brands", title: "品牌管理", name: "Brands" });

    tags.closeOthers("/system/roles");

    expect(tags.visitedTags).toEqual([{ path: "/system/roles", title: "角色管理", name: "Roles" }]);
    expect(localStorage.getItem("mall_admin_tags")).toContain("角色管理");
  });
});
