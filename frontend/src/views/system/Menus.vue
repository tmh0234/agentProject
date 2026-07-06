<template>
  <div class="page-shell">
    <div class="toolbar">
      <div class="filter-row">
        <el-button :icon="Refresh" @click="loadData">刷新菜单</el-button>
      </div>
    </div>

    <el-table v-loading="loading" class="panel" :data="rows" row-key="id" default-expand-all>
      <el-table-column label="菜单名称" prop="title" min-width="180" />
      <el-table-column label="路由名称" prop="name" width="150" />
      <el-table-column label="路径" prop="path" min-width="180" />
      <el-table-column label="组件" prop="component" min-width="180" />
      <el-table-column label="权限标识" prop="permission" min-width="200" />
      <el-table-column label="排序" prop="sort" width="90" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }: { row: Menu }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? "启用" : "禁用" }}</el-tag>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { Refresh } from "@element-plus/icons-vue";
import { onMounted, ref } from "vue";
import { getMenusTreeApi } from "@/api/system";
import { fallbackMenus } from "@/stores/menu";
import type { Menu } from "@/types/system";

const rows = ref<Menu[]>([]);
const loading = ref(false);

async function loadData(): Promise<void> {
  loading.value = true;
  try {
    rows.value = await getMenusTreeApi();
  } catch {
    rows.value = fallbackMenus;
  } finally {
    loading.value = false;
  }
}

onMounted(loadData);
</script>
