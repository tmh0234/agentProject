<template>
  <div class="page-shell">
    <div class="toolbar">
      <div class="filter-row">
        <el-input v-model="query.keyword" clearable placeholder="搜索用户名或昵称" style="width: 220px" />
        <el-select v-model="query.status" clearable placeholder="状态" style="width: 120px">
          <el-option label="启用" :value="1" />
          <el-option label="禁用" :value="0" />
        </el-select>
        <el-button :icon="Search" type="primary" @click="loadData">查询</el-button>
      </div>
      <el-button :icon="Plus" type="primary" @click="openCreate">新增管理员</el-button>
    </div>

    <el-table v-loading="loading" class="panel" :data="rows" row-key="id">
      <el-table-column label="用户名" prop="username" min-width="140" />
      <el-table-column label="昵称" prop="nickname" min-width="140" />
      <el-table-column label="角色" min-width="160">
        <template #default="{ row }: { row: AdminUser }">
          {{ formatRoleNames(row) }}
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }: { row: AdminUser }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? "启用" : "禁用" }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="最后登录" prop="lastLoginAt" min-width="180" />
      <el-table-column label="操作" fixed="right" width="260">
        <template #default="{ row }: { row: AdminUser }">
          <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button link type="warning" @click="resetPassword(row)">重置密码</el-button>
          <el-button link type="danger" @click="removeRow(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-row">
      <el-pagination
        v-model:current-page="query.page"
        v-model:page-size="query.pageSize"
        layout="total, sizes, prev, pager, next"
        :total="total"
        @change="loadData"
      />
    </div>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑管理员' : '新增管理员'" width="520px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="92px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" :disabled="Boolean(editingId)" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="form.nickname" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item v-if="!editingId" label="密码" prop="password">
          <el-input v-model="form.password" placeholder="请输入密码" show-password type="password" />
        </el-form-item>
        <el-form-item label="角色" prop="roleIds">
          <el-select v-model="form.roleIds" multiple placeholder="请选择角色" style="width: 100%">
            <el-option v-for="role in roles" :key="role.id" :label="role.name" :value="role.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="form.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button :loading="saving" type="primary" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { Plus, Search } from "@element-plus/icons-vue";
import type { FormInstance, FormRules } from "element-plus";
import { ElMessage, ElMessageBox } from "element-plus";
import { onMounted, reactive, ref } from "vue";
import {
  createAdminUserApi,
  deleteAdminUserApi,
  getAdminUsersApi,
  getRolesApi,
  resetAdminPasswordApi,
  updateAdminUserApi
} from "@/api/system";
import type { PageQuery } from "@/types/common";
import type { AdminUser, AdminUserForm, Role } from "@/types/system";

const loading = ref(false);
const saving = ref(false);
const rows = ref<AdminUser[]>([]);
const roles = ref<Role[]>([]);
const total = ref(0);
const dialogVisible = ref(false);
const editingId = ref<number | null>(null);
const formRef = ref<FormInstance>();

const query = reactive<PageQuery>({
  page: 1,
  pageSize: 20,
  keyword: "",
  status: undefined
});

const form = reactive<AdminUserForm>({
  username: "",
  nickname: "",
  password: "",
  status: 1,
  roleIds: []
});

const rules: FormRules<AdminUserForm> = {
  username: [{ required: true, message: "请输入用户名", trigger: "blur" }],
  nickname: [{ required: true, message: "请输入昵称", trigger: "blur" }],
  password: [{ required: true, message: "请输入密码", trigger: "blur" }]
};

function resetForm(): void {
  form.username = "";
  form.nickname = "";
  form.password = "";
  form.status = 1;
  form.roleIds = [];
}

function adminRoleIds(row: AdminUser): number[] {
  return row.roleIds || row.roles?.map((role) => role.id) || [];
}

function formatRoleNames(row: AdminUser): string {
  const names = row.roles?.map((role) => role.name) || adminRoleIds(row).map((id) => roles.value.find((role) => role.id === id)?.name || `角色 ${id}`);
  return names.length ? names.join("、") : "未分配";
}

async function loadRoles(): Promise<void> {
  const data = await getRolesApi({ page: 1, pageSize: 200 });
  roles.value = data.items;
}

async function loadData(): Promise<void> {
  loading.value = true;
  try {
    const data = await getAdminUsersApi(query);
    rows.value = data.items;
    total.value = data.total;
  } finally {
    loading.value = false;
  }
}

function openCreate(): void {
  editingId.value = null;
  resetForm();
  dialogVisible.value = true;
}

function openEdit(row: AdminUser): void {
  editingId.value = row.id;
  form.username = row.username;
  form.nickname = row.nickname;
  form.password = "";
  form.status = row.status;
  form.roleIds = adminRoleIds(row);
  dialogVisible.value = true;
}

async function save(): Promise<void> {
  const valid = await formRef.value?.validate();
  if (!valid) {
    return;
  }

  saving.value = true;
  try {
    if (editingId.value) {
      const payload: AdminUserForm = { ...form };
      delete payload.password;
      await updateAdminUserApi(editingId.value, payload);
    } else {
      await createAdminUserApi(form);
    }
    ElMessage.success("保存成功");
    dialogVisible.value = false;
    await loadData();
  } finally {
    saving.value = false;
  }
}

async function resetPassword(row: AdminUser): Promise<void> {
  await ElMessageBox.confirm(`确认将 ${row.nickname} 的密码重置为 Admin@123456？`, "重置密码");
  await resetAdminPasswordApi(row.id, "Admin@123456");
  ElMessage.success("密码已重置");
}

async function removeRow(row: AdminUser): Promise<void> {
  await ElMessageBox.confirm(`确认删除管理员 ${row.nickname}？`, "删除确认", { type: "warning" });
  await deleteAdminUserApi(row.id);
  ElMessage.success("删除成功");
  await loadData();
}

onMounted(async () => {
  await Promise.all([loadRoles(), loadData()]);
});
</script>

<style scoped>
.pagination-row {
  display: flex;
  justify-content: flex-end;
  padding-top: 12px;
}
</style>
