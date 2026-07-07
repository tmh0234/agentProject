import type { Component } from "vue";
import { createRouter, createWebHistory, type RouteRecordRaw } from "vue-router";
import AppLayout from "@/layout/AppLayout.vue";
import { useAuthStore } from "@/stores/auth";
import { fallbackMenus, useMenuStore } from "@/stores/menu";
import { useTagsViewStore } from "@/stores/tagsView";
import type { Menu } from "@/types/system";
import { menuRouteName, resolveViewModulePath } from "./menuRoutes";

const pageModules = import.meta.glob<Component>("@/views/**/*.vue");

const fixedRoutes: RouteRecordRaw[] = [
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/Login.vue"),
    meta: { title: "登录", public: true }
  },
  {
    path: "/",
    component: AppLayout,
    redirect: "/dashboard",
    children: [
      {
        path: "dashboard",
        name: "Dashboard",
        component: () => import("@/views/Dashboard.vue"),
        meta: { title: "仪表盘" }
      }
    ]
  },
  {
    path: "/403",
    name: "Forbidden",
    component: () => import("@/views/error/Forbidden.vue"),
    meta: { title: "无权限", public: true }
  },
  {
    path: "/:pathMatch(.*)*",
    name: "NotFound",
    component: () => import("@/views/error/NotFound.vue"),
    meta: { title: "页面不存在", public: true }
  }
];

export const router = createRouter({
  history: createWebHistory(),
  routes: fixedRoutes
});

type RouteComponent = NonNullable<RouteRecordRaw["component"]>;

function resolveComponent(componentPath: string): RouteComponent {
  const modulePath = resolveViewModulePath(componentPath);
  return (pageModules[modulePath] || (() => import("@/views/error/NotFound.vue"))) as RouteComponent;
}

function toChildRoute(menu: Menu): RouteRecordRaw {
  return {
    path: menu.path,
    name: menuRouteName(menu),
    component: resolveComponent(menu.component),
    meta: {
      title: menu.title,
      icon: menu.icon
    }
  };
}

function flattenLeafMenus(menus: Menu[]): Menu[] {
  return menus.flatMap((menu) => (menu.children?.length ? flattenLeafMenus(menu.children) : [menu]));
}

export function registerDynamicRoutes(menus: Menu[]): void {
  flattenLeafMenus(menus).forEach((menu) => {
    const routeName = menuRouteName(menu);
    if (!router.hasRoute(routeName)) {
      router.addRoute({
        path: menu.path,
        component: AppLayout,
        meta: { title: menu.title, icon: menu.icon },
        children: [
          {
            ...toChildRoute(menu),
            path: ""
          }
        ]
      });
    }
  });
}

registerDynamicRoutes(fallbackMenus);

router.beforeEach(async (to) => {
  const auth = useAuthStore();
  const menu = useMenuStore();

  if (to.meta.public) {
    return true;
  }

  if (!auth.isLoggedIn) {
    return { path: "/login", query: { redirect: to.fullPath } };
  }

  await menu.loadMenus();
  registerDynamicRoutes(menu.menus);
  return true;
});

router.afterEach((to) => {
  if (!to.meta.public && to.name) {
    const tags = useTagsViewStore();
    tags.addTag({
      path: to.fullPath,
      title: String(to.meta.title || "未命名页面"),
      name: String(to.name)
    });
  }
});

export default router;
