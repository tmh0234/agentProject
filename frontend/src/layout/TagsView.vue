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
  height: 42px;
  flex: 0 0 auto;
  align-items: center;
  gap: 7px;
  overflow-x: auto;
  padding: 6px 12px;
  border-bottom: 1px solid var(--app-border);
  background: rgb(248 249 251 / 68%);
  box-shadow: 0 1px 0 rgb(255 255 255 / 86%) inset;
  backdrop-filter: blur(18px);
  scrollbar-width: thin;
}

.tag-button,
.more-button {
  display: inline-flex;
  height: 28px;
  align-items: center;
  gap: 6px;
  border: 1px solid transparent;
  border-radius: 7px;
  background: transparent;
  color: var(--app-text-muted);
  cursor: pointer;
  transition:
    background 0.18s ease,
    border-color 0.18s ease,
    color 0.18s ease,
    box-shadow 0.18s ease;
}

.tag-button {
  padding: 0 9px;
  font-size: 13px;
  font-weight: 600;
  white-space: nowrap;
}

.tag-button:hover,
.more-button:hover {
  border-color: var(--app-border);
  background: rgb(255 255 255 / 64%);
  color: #343941;
}

.tag-button.active {
  border-color: rgb(0 122 255 / 28%);
  background: rgb(255 255 255 / 86%);
  color: #0066d6;
  box-shadow:
    0 1px 0 rgb(255 255 255 / 92%) inset,
    0 8px 20px rgb(0 122 255 / 10%);
}

.close-icon {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  color: #8d949f;
  font-size: 11px;
  transition:
    background 0.18s ease,
    color 0.18s ease;
}

.tag-button.active .close-icon {
  color: #5f9ddd;
}

.close-icon:hover {
  background: rgb(29 29 31 / 8%);
  color: #343941;
}

.more-button {
  width: 30px;
  justify-content: center;
}
</style>
