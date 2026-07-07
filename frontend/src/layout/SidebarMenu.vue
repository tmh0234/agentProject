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
  border-right: 1px solid rgb(255 255 255 / 14%);
  background:
    linear-gradient(180deg, rgb(68 72 80 / 82%), rgb(45 49 56 / 88%)),
    #343840;
  box-shadow:
    1px 0 0 rgb(29 29 31 / 10%),
    18px 0 44px rgb(31 35 43 / 7%);
  transition: width 0.22s ease;
  backdrop-filter: blur(22px);
}

.sidebar-container.collapsed {
  width: 48px;
}

.sidebar-logo {
  position: relative;
  display: flex;
  height: 52px;
  align-items: center;
  gap: 8px;
  padding: 0 8px;
  border-bottom: 1px solid rgb(255 255 255 / 12%);
  color: rgb(255 255 255 / 94%);
}

.logo-icon {
  display: grid;
  width: 28px;
  height: 28px;
  flex: 0 0 28px;
  place-items: center;
  justify-content: center;
  border: 1px solid rgb(255 255 255 / 14%);
  border-radius: 8px;
  background: rgb(255 255 255 / 13%);
  color: #f7f8fa;
  font-size: 17px;
  box-shadow: 0 1px 0 rgb(255 255 255 / 16%) inset;
}

.logo-text {
  max-width: 78px;
  overflow: hidden;
  white-space: nowrap;
  font-size: 14px;
  font-weight: 700;
  letter-spacing: 0;
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
  border: 1px solid rgb(255 255 255 / 12%);
  border-radius: 6px;
  background: rgb(255 255 255 / 10%);
  color: rgb(255 255 255 / 78%);
  cursor: pointer;
  transition:
    background 0.18s ease,
    color 0.18s ease,
    transform 0.18s ease;
}

.collapse-button:hover {
  background: rgb(255 255 255 / 17%);
  color: #ffffff;
}

.collapsed .collapse-button {
  position: absolute;
  right: 3px;
  bottom: 6px;
  width: 18px;
  height: 18px;
  font-size: 11px;
}

.sidebar-menu {
  width: 100%;
  border-right: 0;
  padding: 8px 6px;
  background: transparent;
}

:deep(.el-menu-item),
:deep(.el-sub-menu__title) {
  position: relative;
  height: 38px;
  margin: 2px 0;
  padding-left: 9px !important;
  border-radius: 7px;
  color: rgb(255 255 255 / 74%);
  line-height: 38px;
  transition:
    background 0.18s ease,
    color 0.18s ease,
    box-shadow 0.18s ease;
}

:deep(.el-sub-menu .el-menu-item) {
  min-width: 0;
  padding-left: 28px !important;
  font-size: 13px;
}

:deep(.el-menu-item .el-icon),
:deep(.el-sub-menu__title .el-icon) {
  width: 18px;
  margin-right: 8px;
  color: rgb(255 255 255 / 58%);
  font-size: 16px;
}

:deep(.el-sub-menu__title span),
:deep(.el-menu-item .el-menu-tooltip__trigger span),
:deep(.el-menu-item span) {
  overflow: hidden;
  text-overflow: ellipsis;
}

:deep(.el-menu-item:hover),
:deep(.el-sub-menu__title:hover) {
  background: rgb(255 255 255 / 11%) !important;
  color: #ffffff;
}

:deep(.el-menu-item.is-active) {
  background: rgb(255 255 255 / 20%) !important;
  color: #ffffff;
  box-shadow:
    0 1px 0 rgb(255 255 255 / 16%) inset,
    0 8px 22px rgb(0 0 0 / 13%);
}

:deep(.el-menu-item.is-active::before) {
  position: absolute;
  left: 5px;
  width: 3px;
  height: 16px;
  border-radius: 3px;
  background: #7cc8ff;
  content: "";
}

:deep(.el-menu-item.is-active .el-icon),
:deep(.el-menu-item:hover .el-icon),
:deep(.el-sub-menu__title:hover .el-icon) {
  color: #ffffff;
}

:deep(.el-sub-menu .el-menu) {
  background: transparent;
}

:deep(.el-sub-menu__icon-arrow) {
  right: 9px;
  color: rgb(255 255 255 / 48%);
}

.collapsed .sidebar-logo {
  justify-content: center;
  padding: 0 6px;
}

.collapsed .logo-icon {
  width: 30px;
  height: 30px;
  flex-basis: 30px;
}

.collapsed .sidebar-menu {
  padding-inline: 5px;
}

.collapsed :deep(.el-menu-item),
.collapsed :deep(.el-sub-menu__title) {
  justify-content: center;
  padding-left: 0 !important;
}

.collapsed :deep(.el-menu-item .el-icon),
.collapsed :deep(.el-sub-menu__title .el-icon) {
  margin-right: 0;
}

.collapsed :deep(.el-menu-item.is-active::before) {
  left: 3px;
}
</style>
