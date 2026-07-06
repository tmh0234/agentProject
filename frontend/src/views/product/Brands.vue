<template>
  <div class="page-shell">
    <div class="toolbar">
      <div class="filter-row">
        <el-input v-model="query.keyword" clearable placeholder="搜索品牌名称" style="width: 220px" />
        <el-select v-model="query.status" clearable placeholder="状态" style="width: 120px">
          <el-option label="启用" :value="1" />
          <el-option label="停用" :value="0" />
        </el-select>
        <el-button :icon="Search" type="primary" @click="loadData">查询</el-button>
      </div>
      <el-button :icon="Plus" type="primary" @click="openCreate">新增品牌</el-button>
    </div>

    <el-table v-loading="loading" class="panel" :data="rows" row-key="id">
      <el-table-column label="Logo" width="90">
        <template #default="{ row }: { row: Brand }">
          <el-avatar shape="square" :src="row.logo">{{ row.name.slice(0, 1) }}</el-avatar>
        </template>
      </el-table-column>
      <el-table-column label="品牌名称" prop="name" min-width="180" />
      <el-table-column label="排序" prop="sort" width="100" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }: { row: Brand }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? "启用" : "停用" }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" fixed="right" width="180">
        <template #default="{ row }: { row: Brand }">
          <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
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

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑品牌' : '新增品牌'" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="88px">
        <el-form-item label="品牌名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入品牌名称" />
        </el-form-item>
        <el-form-item label="Logo" prop="logo">
          <el-input v-model="form.logo" placeholder="请输入 Logo 地址" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="form.sort" :min="0" />
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
import { createBrandApi, deleteBrandApi, getBrandsApi, updateBrandApi } from "@/api/product";
import type { PageQuery } from "@/types/common";
import type { Brand, BrandForm } from "@/types/product";

const loading = ref(false);
const saving = ref(false);
const rows = ref<Brand[]>([]);
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

const form = reactive<BrandForm>({
  name: "",
  logo: "",
  status: 1,
  sort: 0
});

const rules: FormRules<BrandForm> = {
  name: [{ required: true, message: "请输入品牌名称", trigger: "blur" }]
};

async function loadData(): Promise<void> {
  loading.value = true;
  try {
    const data = await getBrandsApi(query);
    rows.value = data.items;
    total.value = data.total;
  } finally {
    loading.value = false;
  }
}

function resetForm(): void {
  form.name = "";
  form.logo = "";
  form.status = 1;
  form.sort = 0;
}

function openCreate(): void {
  editingId.value = null;
  resetForm();
  dialogVisible.value = true;
}

function openEdit(row: Brand): void {
  editingId.value = row.id;
  form.name = row.name;
  form.logo = row.logo;
  form.status = row.status;
  form.sort = row.sort;
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
      await updateBrandApi(editingId.value, form);
    } else {
      await createBrandApi(form);
    }
    ElMessage.success("保存成功");
    dialogVisible.value = false;
    await loadData();
  } finally {
    saving.value = false;
  }
}

async function removeRow(row: Brand): Promise<void> {
  await ElMessageBox.confirm(`确认删除品牌 ${row.name}？`, "删除确认", { type: "warning" });
  await deleteBrandApi(row.id);
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
</style>
