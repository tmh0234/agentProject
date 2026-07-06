<template>
  <aside class="sidebar-container" :class="{ collapsed: isCollapsed }">
    <div class="sidebar-logo">
      <el-icon class="logo-icon"><Shop /></el-icon>
      <span class="logo-text">商城后台</span>
      <button class="collapse-button" type="button" @click="isCollapsed = !isCollapsed">
        <el-icon>
          <Fold v-if="!isCollapsed" />
          <Expand v-else />
        </el-icon>
      </button>
    </div>

    <el-menu
      class="sidebar-menu"
      :default-active="route.path"
      :collapse="isCollapsed"
      :collapse-transition="false"
      router
      unique-opened
      background-color="#263445"
      text-color="#c7d2df"
      active-text-color="#ffffff"
    >
      <el-sub-menu v-for="group in menuGroups" :key="group.id" :index="group.path">
        <template #title>
          <el-icon><component :is="resolveIcon(group.icon)" /></el-icon>
          <span>{{ group.title }}</span>
        </template>
        <el-menu-item v-for="item in group.children" :key="item.id" :index="item.path">
          <el-icon><component :is="resolveIcon(item.icon)" /></el-icon>
          <template #title>{{ item.title }}</template>
        </el-menu-item>
      </el-sub-menu>
    </el-menu>
  </aside>
</template>

<script setup lang="ts">
import {
  Avatar,
  Box,
  Collection,
  Expand,
  Fold,
  Goods,
  Grid,
  Menu as MenuIcon,
  PriceTag,
  Setting,
  Shop,
  User
} from "@element-plus/icons-vue";
import { computed, ref, type Component } from "vue";
import { useRoute } from "vue-router";
import { fallbackMenus, useMenuStore } from "@/stores/menu";

const route = useRoute();
const menuStore = useMenuStore();
const isCollapsed = ref(false);

const iconMap: Record<string, Component> = {
  Avatar,
  Box,
  Collection,
  Goods,
  Grid,
  Menu: MenuIcon,
  PriceTag,
  Setting,
  User
};

const menuGroups = computed(() => (menuStore.menus.length ? menuStore.menus : fallbackMenus));

function resolveIcon(name: string): Component {
  return iconMap[name] || Grid;
}
</script>

<style scoped>
.sidebar-container {
  width: 165px;
  height: 100vh;
  flex: 0 0 auto;
  overflow: hidden;
  background: #263445;
  transition: width 0.2s ease;
}

.sidebar-container.collapsed {
  width: 48px;
}

.sidebar-logo {
  display: flex;
  height: 52px;
  align-items: center;
  gap: 8px;
  padding: 0 8px;
  border-bottom: 1px solid rgb(255 255 255 / 8%);
  color: #ffffff;
}

.logo-icon {
  flex: 0 0 28px;
  justify-content: center;
  font-size: 21px;
}

.logo-text {
  max-width: 78px;
  overflow: hidden;
  white-space: nowrap;
  font-size: 15px;
  font-weight: 700;
  opacity: 1;
  transition:
    opacity 0.2s ease,
    max-width 0.2s ease;
}

.collapsed .logo-text {
  max-width: 0;
  opacity: 0;
}

.collapse-button {
  display: grid;
  width: 24px;
  height: 24px;
  flex: 0 0 24px;
  place-items: center;
  border: 0;
  border-radius: 6px;
  background: rgb(255 255 255 / 10%);
  color: #dce7f3;
  cursor: pointer;
}

.collapsed .collapse-button {
  margin-left: -4px;
}

.sidebar-menu {
  width: 100%;
  border-right: 0;
}

:deep(.el-menu-item),
:deep(.el-sub-menu__title) {
  height: 44px;
  padding-left: 14px !important;
}

:deep(.el-menu-item:hover),
:deep(.el-sub-menu__title:hover) {
  background-color: #1f2d3d !important;
}

:deep(.el-menu-item.is-active) {
  background-color: #409eff !important;
}
</style>
