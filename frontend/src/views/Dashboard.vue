<template>
  <div class="page-shell dashboard-page">
    <div class="metric-row">
      <section v-for="item in metrics" :key="item.label" class="metric-card panel">
        <span>{{ item.label }}</span>
        <strong>{{ item.value }}</strong>
        <small>{{ item.hint }}</small>
      </section>
    </div>

    <section class="panel workbench">
      <div>
        <h2>今日工作台</h2>
        <p>关注管理员权限变更、商品上下架和库存异常，保持商城基础数据清晰可追踪。</p>
      </div>
      <el-timeline>
        <el-timeline-item timestamp="09:30">检查待发布商品 SKU 是否完整</el-timeline-item>
        <el-timeline-item timestamp="11:00">同步角色授权与菜单权限</el-timeline-item>
        <el-timeline-item timestamp="16:00">复核品牌和分类启停状态</el-timeline-item>
      </el-timeline>
    </section>
  </div>
</template>

<script setup lang="ts">
const metrics = [
  { label: "在线商品", value: "128", hint: "较昨日 +6" },
  { label: "启用品牌", value: "24", hint: "覆盖 7 个类目" },
  { label: "管理员", value: "12", hint: "3 个角色组" },
  { label: "待处理", value: "5", hint: "库存与权限提醒" }
];
</script>

<style scoped>
.dashboard-page {
  display: grid;
  gap: 16px;
}

.metric-row {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

.metric-card {
  padding: 18px;
}

.metric-card span,
.metric-card small {
  display: block;
  color: #6b7788;
}

.metric-card strong {
  display: block;
  margin: 8px 0;
  font-size: 30px;
  line-height: 1;
}

.workbench {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 26px;
  padding: 20px;
}

.workbench h2 {
  margin: 0 0 8px;
  font-size: 18px;
}

.workbench p {
  margin: 0;
  color: #6b7788;
  line-height: 1.7;
}
</style>
