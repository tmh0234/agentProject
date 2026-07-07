<template>
  <header class="header-bar">
    <el-breadcrumb separator="/">
      <el-breadcrumb-item>商城管理后台</el-breadcrumb-item>
      <el-breadcrumb-item>{{ routeTitle }}</el-breadcrumb-item>
    </el-breadcrumb>

    <el-dropdown @command="handleCommand">
      <button class="user-entry" type="button">
        <el-avatar :size="28">{{ userInitial }}</el-avatar>
        <span>{{ auth.profile?.nickname || "管理员" }}</span>
        <el-icon><ArrowDown /></el-icon>
      </button>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="profile">个人资料</el-dropdown-item>
          <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </header>
</template>

<script setup lang="ts">
import { ArrowDown } from "@element-plus/icons-vue";
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import { useMenuStore } from "@/stores/menu";

const route = useRoute();
const router = useRouter();
const auth = useAuthStore();
const menu = useMenuStore();

const routeTitle = computed(() => String(route.meta.title || "工作台"));
const userInitial = computed(() => (auth.profile?.nickname || auth.profile?.username || "管").slice(0, 1));

function handleCommand(command: string | number | object): void {
  if (command === "logout") {
    auth.logout();
    menu.reset();
    void router.push("/login");
  }
}
</script>

<style scoped>
.header-bar {
  display: flex;
  height: 52px;
  flex: 0 0 auto;
  align-items: center;
  justify-content: space-between;
  padding: 0 18px;
  border-bottom: 1px solid var(--app-border);
  background: rgb(255 255 255 / 72%);
  box-shadow: 0 1px 0 rgb(255 255 255 / 84%) inset;
  backdrop-filter: blur(20px);
}

:deep(.el-breadcrumb) {
  font-size: 13px;
}

:deep(.el-breadcrumb__inner) {
  color: var(--app-text-subtle);
  font-weight: 500;
}

:deep(.el-breadcrumb__item:last-child .el-breadcrumb__inner) {
  color: #3d424a;
  font-weight: 650;
}

:deep(.el-breadcrumb__separator) {
  color: #c4c8cf;
  font-weight: 400;
}

.user-entry {
  display: inline-flex;
  height: 34px;
  align-items: center;
  gap: 8px;
  padding: 2px 8px 2px 3px;
  border: 1px solid var(--app-border);
  border-radius: 8px;
  background: rgb(255 255 255 / 72%);
  color: var(--app-text);
  cursor: pointer;
  box-shadow:
    0 1px 0 rgb(255 255 255 / 90%) inset,
    0 8px 22px rgb(31 35 43 / 5%);
  transition:
    border-color 0.18s ease,
    background 0.18s ease,
    box-shadow 0.18s ease;
}

.user-entry:hover {
  border-color: rgb(29 29 31 / 18%);
  background: rgb(255 255 255 / 92%);
  box-shadow:
    0 1px 0 rgb(255 255 255 / 95%) inset,
    0 10px 24px rgb(31 35 43 / 8%);
}

.user-entry span {
  max-width: 120px;
  overflow: hidden;
  color: #343941;
  font-size: 13px;
  font-weight: 650;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-entry :deep(.el-avatar) {
  background: linear-gradient(180deg, #f2f4f7, #d9dee7);
  color: #4a5059;
  font-size: 12px;
  font-weight: 800;
  box-shadow: 0 1px 0 rgb(255 255 255 / 92%) inset;
}
</style>
