<template>
  <div class="page-shell dashboard-page">
    <div class="metric-row">
      <section v-for="item in metrics" :key="item.label" class="metric-card panel">
        <div class="metric-top">
          <span class="metric-icon" :class="item.tone">
            <el-icon><component :is="item.icon" /></el-icon>
          </span>
          <span class="metric-trend">{{ item.trend }}</span>
        </div>
        <span class="metric-label">{{ item.label }}</span>
        <strong>{{ item.value }}</strong>
        <small>{{ item.hint }}</small>
      </section>
    </div>

    <section class="panel workbench">
      <div class="workbench-summary">
        <span class="section-kicker">今日工作台</span>
        <h2>保持商城基础数据清晰可追踪</h2>
        <p>关注管理员权限变更、商品上下架和库存异常，优先处理会影响运营准确性的事项。</p>
        <div class="focus-grid">
          <span>权限</span>
          <span>商品</span>
          <span>库存</span>
        </div>
      </div>
      <div class="task-panel">
        <div v-for="task in tasks" :key="task.time" class="task-item">
          <span class="task-time">{{ task.time }}</span>
          <span class="task-dot" :class="task.tone"></span>
          <div>
            <strong>{{ task.title }}</strong>
            <small>{{ task.detail }}</small>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { Bell, CollectionTag, Goods, UserFilled } from "@element-plus/icons-vue";
import type { Component } from "vue";

interface MetricItem {
  label: string;
  value: string;
  hint: string;
  trend: string;
  tone: "blue" | "green" | "orange" | "graphite";
  icon: Component;
}

interface TaskItem {
  time: string;
  title: string;
  detail: string;
  tone: "blue" | "green" | "orange";
}

const metrics: MetricItem[] = [
  { label: "在线商品", value: "128", hint: "较昨日 +6", trend: "+4.9%", tone: "blue", icon: Goods },
  { label: "启用品牌", value: "24", hint: "覆盖 7 个类目", trend: "稳定", tone: "green", icon: CollectionTag },
  { label: "管理员", value: "12", hint: "3 个角色组", trend: "权限正常", tone: "graphite", icon: UserFilled },
  { label: "待处理", value: "5", hint: "库存与权限提醒", trend: "需关注", tone: "orange", icon: Bell }
];

const tasks: TaskItem[] = [
  { time: "09:30", title: "检查待发布商品 SKU 是否完整", detail: "避免上架后规格和库存不可选", tone: "blue" },
  { time: "11:00", title: "同步角色授权与菜单权限", detail: "确保管理员入口和 API 权限一致", tone: "green" },
  { time: "16:00", title: "复核品牌和分类启停状态", detail: "减少列表筛选和商品归类异常", tone: "orange" }
];
</script>

<style scoped>
.dashboard-page {
  display: grid;
  gap: 18px;
}

.metric-row {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
}

.metric-card {
  position: relative;
  min-height: 152px;
  padding: 16px;
}

.metric-card::after {
  position: absolute;
  right: 0;
  bottom: 0;
  width: 84px;
  height: 84px;
  background: linear-gradient(135deg, transparent 18%, rgb(255 255 255 / 64%));
  content: "";
  pointer-events: none;
}

.metric-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 18px;
}

.metric-icon {
  display: grid;
  width: 34px;
  height: 34px;
  place-items: center;
  border: 1px solid rgb(255 255 255 / 72%);
  border-radius: 8px;
  font-size: 18px;
  box-shadow:
    0 1px 0 rgb(255 255 255 / 88%) inset,
    0 10px 20px rgb(31 35 43 / 7%);
}

.metric-icon.blue {
  background: rgb(0 122 255 / 11%);
  color: #006edb;
}

.metric-icon.green {
  background: rgb(52 199 89 / 13%);
  color: #238a3b;
}

.metric-icon.orange {
  background: rgb(255 159 10 / 15%);
  color: #b36b00;
}

.metric-icon.graphite {
  background: rgb(47 51 58 / 10%);
  color: #3f444c;
}

.metric-trend {
  padding: 3px 7px;
  border: 1px solid rgb(29 29 31 / 8%);
  border-radius: 7px;
  background: rgb(255 255 255 / 58%);
  color: #626976;
  font-size: 12px;
  font-weight: 650;
}

.metric-label,
.metric-card small {
  display: block;
  color: var(--app-text-muted);
}

.metric-card strong {
  display: block;
  margin: 7px 0;
  color: #17191c;
  font-size: 34px;
  font-weight: 760;
  line-height: 1;
}

.metric-card small {
  font-weight: 600;
}

.workbench {
  display: grid;
  grid-template-columns: minmax(300px, 0.82fr) 1.18fr;
  gap: 22px;
  padding: 18px;
}

.workbench-summary {
  padding: 2px 4px 2px 2px;
}

.section-kicker {
  display: inline-flex;
  margin-bottom: 10px;
  color: #007aff;
  font-size: 13px;
  font-weight: 760;
}

.workbench h2 {
  max-width: 360px;
  margin: 0 0 10px;
  color: #17191c;
  font-size: 22px;
  font-weight: 760;
  line-height: 1.28;
}

.workbench p {
  margin: 0;
  color: var(--app-text-muted);
  line-height: 1.7;
}

.focus-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
  margin-top: 20px;
}

.focus-grid span {
  display: grid;
  min-height: 34px;
  place-items: center;
  border: 1px solid rgb(29 29 31 / 8%);
  border-radius: 7px;
  background: rgb(255 255 255 / 56%);
  color: #4f5661;
  font-size: 13px;
  font-weight: 700;
}

.task-panel {
  display: grid;
  gap: 10px;
}

.task-item {
  display: grid;
  grid-template-columns: 54px 14px minmax(0, 1fr);
  gap: 12px;
  align-items: start;
  padding: 12px;
  border: 1px solid rgb(29 29 31 / 8%);
  border-radius: 8px;
  background: rgb(255 255 255 / 58%);
  box-shadow: 0 1px 0 rgb(255 255 255 / 80%) inset;
}

.task-time {
  color: var(--app-text-subtle);
  font-size: 13px;
  font-weight: 700;
  line-height: 1.5;
}

.task-dot {
  width: 9px;
  height: 9px;
  margin-top: 6px;
  border-radius: 50%;
  box-shadow: 0 0 0 4px rgb(29 29 31 / 5%);
}

.task-dot.blue {
  background: var(--app-blue);
}

.task-dot.green {
  background: var(--app-green);
}

.task-dot.orange {
  background: var(--app-orange);
}

.task-item strong,
.task-item small {
  display: block;
}

.task-item strong {
  color: #25282d;
  font-size: 14px;
  font-weight: 720;
  line-height: 1.5;
}

.task-item small {
  margin-top: 3px;
  color: var(--app-text-muted);
  font-size: 12px;
  line-height: 1.55;
}
</style>
