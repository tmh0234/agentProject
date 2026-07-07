<template>
  <div class="page-shell">
    <div class="toolbar">
      <div class="filter-row">
        <el-input v-model="query.keyword" clearable placeholder="搜索商品名称或编码" style="width: 240px" />
        <el-select v-model="query.status" clearable placeholder="上下架" style="width: 120px">
          <el-option label="上架" :value="1" />
          <el-option label="下架" :value="0" />
        </el-select>
        <el-button :icon="Search" type="primary" @click="loadData">查询</el-button>
      </div>
      <el-button :icon="Plus" type="primary" @click="openCreate">新增商品</el-button>
    </div>

    <el-table v-loading="loading" class="panel" :data="rows" row-key="id">
      <el-table-column label="商品" min-width="260">
        <template #default="{ row }: { row: Product }">
          <div class="product-cell">
            <el-avatar shape="square" :src="row.mainImage">{{ row.name.slice(0, 1) }}</el-avatar>
            <div>
              <strong>{{ row.name }}</strong>
              <span>{{ row.code }}</span>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="价格区间" width="150">
        <template #default="{ row }: { row: Product }">
          ￥{{ centsToYuan(row.minPriceCents) }} - ￥{{ centsToYuan(row.maxPriceCents) }}
        </template>
      </el-table-column>
      <el-table-column label="SKU 数" width="100">
        <template #default="{ row }: { row: Product }">
          {{ row.skus.length }}
        </template>
      </el-table-column>
      <el-table-column label="状态" width="120">
        <template #default="{ row }: { row: Product }">
          <el-switch
            :model-value="row.status"
            :active-value="1"
            :inactive-value="0"
            @change="(value: string | number | boolean) => changeStatus(row, Number(value))"
          />
        </template>
      </el-table-column>
      <el-table-column label="操作" fixed="right" width="190">
        <template #default="{ row }: { row: Product }">
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

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑商品' : '新增商品'" width="920px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="92px">
        <div class="form-grid">
          <el-form-item label="商品名称" prop="name">
            <el-input v-model="form.name" placeholder="请输入商品名称" />
          </el-form-item>
          <el-form-item label="商品编码" prop="code">
            <el-input v-model="form.code" placeholder="请输入商品编码" />
          </el-form-item>
          <el-form-item label="商品分类" prop="categoryId">
            <el-tree-select
              v-model="form.categoryId"
              check-strictly
              :data="categories"
              :props="{ label: 'name', children: 'children', value: 'id' }"
              placeholder="请选择分类"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item label="品牌" prop="brandId">
            <el-select v-model="form.brandId" placeholder="请选择品牌" style="width: 100%">
              <el-option v-for="brand in brands" :key="brand.id" :label="brand.name" :value="brand.id" />
            </el-select>
          </el-form-item>
          <el-form-item label="主图" prop="mainImage" class="span-2">
            <el-input v-model="form.mainImage" placeholder="请输入主图地址" />
          </el-form-item>
          <el-form-item label="商品描述" prop="description" class="span-2">
            <el-input v-model="form.description" placeholder="请输入商品描述" type="textarea" />
          </el-form-item>
          <el-form-item label="状态" prop="status">
            <el-switch v-model="form.status" :active-value="1" :inactive-value="0" />
          </el-form-item>
        </div>
      </el-form>

      <section class="sku-section">
        <div class="sku-header">
          <h3>SKU 行编辑</h3>
          <el-button :icon="Plus" @click="addSku">新增 SKU</el-button>
        </div>
        <el-table :data="form.skus" row-key="localId" border>
          <el-table-column label="SKU 编码" min-width="150">
            <template #default="{ row }: { row: ProductSkuDraft }">
              <el-input v-model="row.code" placeholder="SKU 编码" />
            </template>
          </el-table-column>
          <el-table-column label="规格" min-width="260">
            <template #default="{ row }: { row: ProductSkuDraft }">
              <div class="spec-editor">
                <div v-for="(pair, index) in row.specPairs" :key="index" class="spec-row">
                  <el-input v-model="pair.key" placeholder="规格名" />
                  <el-input v-model="pair.value" placeholder="规格值" />
                  <el-button :icon="Delete" circle @click="row.specPairs.splice(index, 1)" />
                </div>
                <el-button link type="primary" @click="row.specPairs.push({ key: '', value: '' })">添加规格</el-button>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="售价（元）" width="150">
            <template #default="{ row }: { row: ProductSkuDraft }">
              <el-input v-model="row.priceYuan" placeholder="0.00" />
            </template>
          </el-table-column>
          <el-table-column label="库存" width="130">
            <template #default="{ row }: { row: ProductSkuDraft }">
              <el-input-number v-model="row.stock" :min="0" controls-position="right" />
            </template>
          </el-table-column>
          <el-table-column label="状态" width="100">
            <template #default="{ row }: { row: ProductSkuDraft }">
              <el-switch v-model="row.status" :active-value="1" :inactive-value="0" />
            </template>
          </el-table-column>
          <el-table-column label="操作" width="90">
            <template #default="{ $index }: { $index: number }">
              <el-button link type="danger" @click="form.skus.splice($index, 1)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </section>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button :loading="saving" type="primary" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { Delete, Plus, Search } from "@element-plus/icons-vue";
import type { FormInstance, FormRules } from "element-plus";
import { ElMessage, ElMessageBox } from "element-plus";
import { onMounted, reactive, ref } from "vue";
import {
  createProductApi,
  deleteProductApi,
  getBrandsApi,
  getCategoriesApi,
  getProductsApi,
  updateProductApi,
  updateProductStatusApi
} from "@/api/product";
import type { PageQuery } from "@/types/common";
import type { Brand, Product, ProductCategory } from "@/types/product";
import {
  centsToYuan,
  createEmptySkuDraft,
  productDraftToPayload,
  productToDraft,
  type ProductDraft,
  type ProductSkuDraft
} from "@/utils/product";

const loading = ref(false);
const saving = ref(false);
const rows = ref<Product[]>([]);
const total = ref(0);
const dialogVisible = ref(false);
const editingId = ref<number | null>(null);
const formRef = ref<FormInstance>();
const categories = ref<ProductCategory[]>([]);
const brands = ref<Brand[]>([]);

const query = reactive<PageQuery>({
  page: 1,
  pageSize: 20,
  keyword: "",
  status: undefined
});

const form = reactive<ProductDraft>({
  name: "",
  code: "",
  categoryId: null,
  brandId: null,
  mainImage: "",
  status: 0,
  description: "",
  skus: []
});

const rules: FormRules<ProductDraft> = {
  name: [{ required: true, message: "请输入商品名称", trigger: "blur" }],
  code: [{ required: true, message: "请输入商品编码", trigger: "blur" }],
  categoryId: [{ required: true, message: "请选择商品分类", trigger: "change" }],
  brandId: [{ required: true, message: "请选择品牌", trigger: "change" }]
};

function createLocalId(): string {
  return `${Date.now()}-${Math.random().toString(16).slice(2)}`;
}

function resetForm(): void {
  form.name = "";
  form.code = "";
  form.categoryId = null;
  form.brandId = null;
  form.mainImage = "";
  form.status = 0;
  form.description = "";
  form.skus = [];
  addSku();
}

function addSku(): void {
  form.skus.push(createEmptySkuDraft(createLocalId()));
}

async function loadOptions(): Promise<void> {
  const [categoryData, brandData] = await Promise.all([
    getCategoriesApi({ status: 1 }),
    getBrandsApi({ page: 1, pageSize: 500, status: 1 })
  ]);
  categories.value = categoryData;
  brands.value = brandData.items;
}

async function loadData(): Promise<void> {
  loading.value = true;
  try {
    const data = await getProductsApi(query);
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

function openEdit(row: Product): void {
  editingId.value = row.id;
  const draft = productToDraft(row, createLocalId);
  form.name = draft.name;
  form.code = draft.code;
  form.categoryId = draft.categoryId;
  form.brandId = draft.brandId;
  form.mainImage = draft.mainImage;
  form.status = draft.status;
  form.description = draft.description;
  form.skus = draft.skus;
  if (!form.skus.length) {
    addSku();
  }
  dialogVisible.value = true;
}

async function save(): Promise<void> {
  const valid = await formRef.value?.validate();
  if (!valid) {
    return;
  }

  if (!form.skus.length) {
    ElMessage.warning("请至少添加一个 SKU");
    return;
  }

  saving.value = true;
  try {
    const payload = productDraftToPayload(form);
    if (editingId.value) {
      await updateProductApi(editingId.value, payload);
    } else {
      await createProductApi(payload);
    }
    ElMessage.success("保存成功");
    dialogVisible.value = false;
    await loadData();
  } finally {
    saving.value = false;
  }
}

async function changeStatus(row: Product, status: number): Promise<void> {
  await updateProductStatusApi(row.id, status);
  row.status = status === 1 ? 1 : 0;
  ElMessage.success("状态已更新");
}

async function removeRow(row: Product): Promise<void> {
  await ElMessageBox.confirm(`确认删除商品 ${row.name}？`, "删除确认", { type: "warning" });
  await deleteProductApi(row.id);
  ElMessage.success("删除成功");
  await loadData();
}

onMounted(async () => {
  await Promise.all([loadOptions(), loadData()]);
});
</script>

<style scoped>
.product-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.product-cell strong,
.product-cell span {
  display: block;
}

.product-cell span {
  margin-top: 4px;
  color: #6b7788;
  font-size: 12px;
}

.pagination-row {
  display: flex;
  justify-content: flex-end;
  padding-top: 12px;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0 14px;
}

.span-2 {
  grid-column: span 2;
}

.sku-section {
  margin-top: 4px;
}

.sku-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.sku-header h3 {
  margin: 0;
  font-size: 15px;
}

.spec-editor {
  display: grid;
  gap: 6px;
}

.spec-row {
  display: grid;
  grid-template-columns: 1fr 1fr 32px;
  gap: 6px;
}
</style>
