<template>
  <div class="tags-view">
    <button
      v-for="tag in tags.visitedTags"
      :key="tag.path"
      class="tag-button"
      :class="{ active: tag.path === route.fullPath || tag.path === route.path }"
      type="button"
      @click="router.push(tag.path)"
    >
      <span>{{ tag.title }}</span>
      <el-icon v-if="tags.visitedTags.length > 1" class="close-icon" @click.stop="closeTag(tag.path)">
        <Close />
      </el-icon>
    </button>
    <el-dropdown v-if="tags.visitedTags.length" trigger="click" @command="handleCommand">
      <button class="more-button" type="button">
        <el-icon><MoreFilled /></el-icon>
      </button>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="close-current">关闭当前</el-dropdown-item>
          <el-dropdown-item command="close-others">关闭其他</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script setup lang="ts">
import { Close, MoreFilled } from "@element-plus/icons-vue";
import { useRoute, useRouter } from "vue-router";
import { useTagsViewStore } from "@/stores/tagsView";

const tags = useTagsViewStore();
const route = useRoute();
const router = useRouter();

function closeTag(path: string): void {
  const next = tags.closeCurrent(path);
  if (path === route.fullPath || path === route.path) {
    void router.push(next?.path || "/dashboard");
  }
}

function handleCommand(command: string | number | object): void {
  if (command === "close-current") {
    closeTag(route.fullPath);
  }

  if (command === "close-others") {
    tags.closeOthers(route.fullPath);
  }
}
</script>

<style scoped>
.tags-view {
  display: flex;
  height: 38px;
  flex: 0 0 auto;
  align-items: center;
  gap: 6px;
  overflow-x: auto;
  padding: 0 10px;
  border-bottom: 1px solid #dfe6ef;
  background: #ffffff;
}

.tag-button,
.more-button {
  display: inline-flex;
  height: 26px;
  align-items: center;
  gap: 6px;
  border: 1px solid #dfe6ef;
  border-radius: 6px;
  background: #ffffff;
  color: #526173;
  cursor: pointer;
}

.tag-button {
  padding: 0 8px;
  white-space: nowrap;
}

.tag-button.active {
  border-color: #409eff;
  background: #ecf5ff;
  color: #1677c8;
}

.close-icon {
  font-size: 12px;
}

.more-button {
  width: 28px;
  justify-content: center;
}
</style>
