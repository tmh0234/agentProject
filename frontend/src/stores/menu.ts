import { defineStore } from "pinia";
import { ref } from "vue";
import { getAuthMenusApi } from "@/api/auth";
import type { Menu } from "@/types/system";

export const fallbackMenus: Menu[] = [
  {
    id: 1,
    parentId: 0,
    title: "系统管理",
    routeName: "System",
    path: "/system",
    component: "Layout",
    icon: "Setting",
    permission: "",
    sort: 1,
    status: 1,
    children: [
      {
        id: 11,
        parentId: 1,
        title: "管理员管理",
        routeName: "AdminUsers",
        path: "/system/admin-users",
        component: "views/system/AdminUsers.vue",
        icon: "User",
        permission: "system:admin-users:list",
        sort: 1,
        status: 1
      },
      {
        id: 12,
        parentId: 1,
        title: "角色管理",
        routeName: "Roles",
        path: "/system/roles",
        component: "views/system/Roles.vue",
        icon: "Avatar",
        permission: "system:roles:list",
        sort: 2,
        status: 1
      },
      {
        id: 13,
        parentId: 1,
        title: "菜单管理",
        routeName: "Menus",
        path: "/system/menus",
        component: "views/system/Menus.vue",
        icon: "Menu",
        permission: "system:menus:list",
        sort: 3,
        status: 1
      }
    ]
  },
  {
    id: 2,
    parentId: 0,
    title: "商品管理",
    routeName: "Product",
    path: "/product",
    component: "Layout",
    icon: "Goods",
    permission: "",
    sort: 2,
    status: 1,
    children: [
      {
        id: 21,
        parentId: 2,
        title: "商品分类",
        routeName: "Categories",
        path: "/product/categories",
        component: "views/product/Categories.vue",
        icon: "Collection",
        permission: "product:categories:list",
        sort: 1,
        status: 1
      },
      {
        id: 22,
        parentId: 2,
        title: "品牌管理",
        routeName: "Brands",
        path: "/product/brands",
        component: "views/product/Brands.vue",
        icon: "PriceTag",
        permission: "product:brands:list",
        sort: 2,
        status: 1
      },
      {
        id: 23,
        parentId: 2,
        title: "商品列表",
        routeName: "Products",
        path: "/product/products",
        component: "views/product/Products.vue",
        icon: "Box",
        permission: "product:products:list",
        sort: 3,
        status: 1
      }
    ]
  }
];

export const useMenuStore = defineStore("menu", () => {
  const menus = ref<Menu[]>([]);
  const loaded = ref(false);

  async function loadMenus(): Promise<void> {
    if (loaded.value) {
      return;
    }

    try {
      const data = await getAuthMenusApi();
      menus.value = data.length ? data : fallbackMenus;
    } catch {
      menus.value = fallbackMenus;
    }
    loaded.value = true;
  }

  function reset(): void {
    menus.value = [];
    loaded.value = false;
  }

  return {
    menus,
    loaded,
    loadMenus,
    reset
  };
});
