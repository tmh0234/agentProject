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
  border-bottom: 1px solid #dfe6ef;
  background: #ffffff;
}

.user-entry {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  border: 0;
  background: transparent;
  color: #202833;
  cursor: pointer;
}
</style>
