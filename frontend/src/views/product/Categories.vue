<template>
  <div class="page-shell">
    <div class="toolbar">
      <div class="filter-row">
        <el-button :icon="Refresh" @click="loadData">刷新</el-button>
      </div>
      <el-button :icon="Plus" type="primary" @click="openCreate">新增分类</el-button>
    </div>

    <el-table v-loading="loading" class="panel" :data="rows" row-key="id" default-expand-all>
      <el-table-column label="分类名称" prop="name" min-width="200" />
      <el-table-column label="排序" prop="sort" width="100" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }: { row: ProductCategory }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? "启用" : "停用" }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" fixed="right" width="180">
        <template #default="{ row }: { row: ProductCategory }">
          <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button link type="danger" @click="removeRow(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑分类' : '新增分类'" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
        <el-form-item label="上级分类">
          <el-tree-select
            v-model="form.parentId"
            check-strictly
            :data="parentOptions"
            :props="{ label: 'name', children: 'children', value: 'id' }"
            placeholder="请选择上级分类"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="分类名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入分类名称" />
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
import { Plus, Refresh } from "@element-plus/icons-vue";
import type { FormInstance, FormRules } from "element-plus";
import { ElMessage, ElMessageBox } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import { createCategoryApi, deleteCategoryApi, getCategoriesApi, updateCategoryApi } from "@/api/product";
import type { ProductCategory, ProductCategoryForm } from "@/types/product";

const loading = ref(false);
const saving = ref(false);
const rows = ref<ProductCategory[]>([]);
const dialogVisible = ref(false);
const editingId = ref<number | null>(null);
const formRef = ref<FormInstance>();

const form = reactive<ProductCategoryForm>({
  parentId: 0,
  name: "",
  status: 1,
  sort: 0
});

const rules: FormRules<ProductCategoryForm> = {
  name: [{ required: true, message: "请输入分类名称", trigger: "blur" }]
};

const parentOptions = computed<ProductCategory[]>(() => [{ id: 0, parentId: 0, name: "顶级分类", status: 1, sort: 0, createdAt: "", updatedAt: "", children: rows.value }]);

async function loadData(): Promise<void> {
  loading.value = true;
  try {
    rows.value = await getCategoriesApi();
  } finally {
    loading.value = false;
  }
}

function resetForm(): void {
  form.parentId = 0;
  form.name = "";
  form.status = 1;
  form.sort = 0;
}

function openCreate(): void {
  editingId.value = null;
  resetForm();
  dialogVisible.value = true;
}

function openEdit(row: ProductCategory): void {
  editingId.value = row.id;
  form.parentId = row.parentId;
  form.name = row.name;
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
      await updateCategoryApi(editingId.value, form);
    } else {
      await createCategoryApi(form);
    }
    ElMessage.success("保存成功");
    dialogVisible.value = false;
    await loadData();
  } finally {
    saving.value = false;
  }
}

async function removeRow(row: ProductCategory): Promise<void> {
  await ElMessageBox.confirm(`确认删除分类 ${row.name}？`, "删除确认", { type: "warning" });
  await deleteCategoryApi(row.id);
  ElMessage.success("删除成功");
  await loadData();
}

onMounted(loadData);
</script>
