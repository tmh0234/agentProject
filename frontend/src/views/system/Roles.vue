<template>
  <div class="page-shell">
    <div class="toolbar">
      <div class="filter-row">
        <el-input v-model="query.keyword" clearable placeholder="搜索角色名称或编码" style="width: 240px" />
        <el-button :icon="Search" type="primary" @click="loadData">查询</el-button>
      </div>
      <el-button :icon="Plus" type="primary" @click="openCreate">新增角色</el-button>
    </div>

    <el-table v-loading="loading" class="panel" :data="rows" row-key="id">
      <el-table-column label="角色名称" prop="name" min-width="150" />
      <el-table-column label="角色编码" prop="code" min-width="160" />
      <el-table-column label="描述" prop="description" min-width="220" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }: { row: Role }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? "启用" : "禁用" }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" fixed="right" width="230">
        <template #default="{ row }: { row: Role }">
          <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button link type="primary" @click="openPermission(row)">授权</el-button>
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

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑角色' : '新增角色'" width="520px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="92px">
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="角色编码" prop="code">
          <el-input v-model="form.code" placeholder="请输入角色编码" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" placeholder="请输入描述" type="textarea" />
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

    <el-dialog v-model="permissionVisible" title="角色授权" width="720px">
      <div class="permission-grid">
        <section>
          <h3>菜单权限</h3>
          <el-tree
            ref="menuTreeRef"
            show-checkbox
            node-key="id"
            :data="menus"
            :props="{ label: 'title', children: 'children' }"
          />
        </section>
        <section>
          <h3>API 权限</h3>
          <el-checkbox-group v-model="selectedApiIds" class="api-list">
            <el-checkbox v-for="api in apiPermissions" :key="api.id" :label="api.id">
              {{ api.method }} {{ api.path }}
            </el-checkbox>
          </el-checkbox-group>
        </section>
      </div>
      <template #footer>
        <el-button @click="permissionVisible = false">取消</el-button>
        <el-button :loading="saving" type="primary" @click="savePermission">保存授权</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { Plus, Search } from "@element-plus/icons-vue";
import type { ElTree, FormInstance, FormRules } from "element-plus";
import { ElMessage, ElMessageBox } from "element-plus";
import { onMounted, reactive, ref } from "vue";
import {
  createRoleApi,
  deleteRoleApi,
  getApiPermissionsApi,
  getMenusTreeApi,
  getRolePermissionsApi,
  getRolesApi,
  updateRoleApi,
  updateRolePermissionsApi
} from "@/api/system";
import type { PageQuery } from "@/types/common";
import type { ApiPermission, Menu, Role, RoleForm } from "@/types/system";

const loading = ref(false);
const saving = ref(false);
const rows = ref<Role[]>([]);
const total = ref(0);
const dialogVisible = ref(false);
const permissionVisible = ref(false);
const editingId = ref<number | null>(null);
const permissionRoleId = ref<number | null>(null);
const formRef = ref<FormInstance>();
const menuTreeRef = ref<InstanceType<typeof ElTree>>();
const menus = ref<Menu[]>([]);
const apiPermissions = ref<ApiPermission[]>([]);
const selectedApiIds = ref<number[]>([]);

const query = reactive<PageQuery>({
  page: 1,
  pageSize: 20,
  keyword: ""
});

const form = reactive<RoleForm>({
  name: "",
  code: "",
  status: 1,
  description: ""
});

const rules: FormRules<RoleForm> = {
  name: [{ required: true, message: "请输入角色名称", trigger: "blur" }],
  code: [{ required: true, message: "请输入角色编码", trigger: "blur" }]
};

async function loadData(): Promise<void> {
  loading.value = true;
  try {
    const data = await getRolesApi(query);
    rows.value = data.items;
    total.value = data.total;
  } finally {
    loading.value = false;
  }
}

function resetForm(): void {
  form.name = "";
  form.code = "";
  form.status = 1;
  form.description = "";
}

function openCreate(): void {
  editingId.value = null;
  resetForm();
  dialogVisible.value = true;
}

function openEdit(row: Role): void {
  editingId.value = row.id;
  form.name = row.name;
  form.code = row.code;
  form.status = row.status;
  form.description = row.description;
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
      await updateRoleApi(editingId.value, form);
    } else {
      await createRoleApi(form);
    }
    ElMessage.success("保存成功");
    dialogVisible.value = false;
    await loadData();
  } finally {
    saving.value = false;
  }
}

async function openPermission(row: Role): Promise<void> {
  permissionRoleId.value = row.id;
  permissionVisible.value = true;
  const [menuData, apiData, current] = await Promise.all([
    getMenusTreeApi(),
    getApiPermissionsApi(),
    getRolePermissionsApi(row.id)
  ]);
  menus.value = menuData;
  apiPermissions.value = apiData;
  selectedApiIds.value = current.apiPermissionIds;
  window.setTimeout(() => {
    menuTreeRef.value?.setCheckedKeys(current.menuIds, false);
  });
}

async function savePermission(): Promise<void> {
  if (!permissionRoleId.value) {
    return;
  }

  saving.value = true;
  try {
    await updateRolePermissionsApi(permissionRoleId.value, {
      menuIds: menuTreeRef.value?.getCheckedKeys(false).map(Number) || [],
      apiPermissionIds: selectedApiIds.value
    });
    ElMessage.success("授权已保存");
    permissionVisible.value = false;
  } finally {
    saving.value = false;
  }
}

async function removeRow(row: Role): Promise<void> {
  await ElMessageBox.confirm(`确认删除角色 ${row.name}？`, "删除确认", { type: "warning" });
  await deleteRoleApi(row.id);
  ElMessage.success("删除成功");
  await loadData();
}

onMounted(loadData);
</script>

<style scoped>
.pagination-row {
  display: flex;
  justify-content: flex-end;
  padding-top: 12px;
}

.permission-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 18px;
}

.permission-grid h3 {
  margin: 0 0 10px;
  font-size: 15px;
}

.api-list {
  display: grid;
  max-height: 360px;
  overflow: auto;
  gap: 8px;
}
</style>
