import { defineStore } from "pinia";
import { ref } from "vue";
import type { TagView } from "@/types/common";
import { readStorage, writeStorage } from "@/utils/storage";

const TAGS_KEY = "mall_admin_tags";

function samePath(tag: TagView, path: string): boolean {
  return tag.path === path;
}

export const useTagsViewStore = defineStore("tagsView", () => {
  const visitedTags = ref<TagView[]>(readStorage<TagView[]>(TAGS_KEY) || []);

  function persist(): void {
    writeStorage(TAGS_KEY, visitedTags.value);
  }

  function addTag(tag: TagView): void {
    if (!visitedTags.value.some((item) => item.path === tag.path)) {
      visitedTags.value.push(tag);
      persist();
    }
  }

  function closeCurrent(path: string): TagView | null {
    if (visitedTags.value.length <= 1) {
      return visitedTags.value[0] || null;
    }

    const index = visitedTags.value.findIndex((tag) => samePath(tag, path));
    if (index < 0) {
      return visitedTags.value[visitedTags.value.length - 1] || null;
    }

    visitedTags.value.splice(index, 1);
    persist();
    return visitedTags.value[index] || visitedTags.value[index - 1] || null;
  }

  function closeOthers(path: string): void {
    const current = visitedTags.value.find((tag) => samePath(tag, path));
    visitedTags.value = current ? [current] : [];
    persist();
  }

  function clear(): void {
    visitedTags.value = [];
    persist();
  }

  return {
    visitedTags,
    addTag,
    closeCurrent,
    closeOthers,
    clear
  };
});
