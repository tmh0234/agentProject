<template>
  <div class="admin-layout">
    <SidebarMenu />
    <section class="main-area">
      <HeaderBar />
      <TagsView />
      <main class="content-area">
        <router-view />
      </main>
    </section>
  </div>
</template>

<script setup lang="ts">
import HeaderBar from "./HeaderBar.vue";
import SidebarMenu from "./SidebarMenu.vue";
import TagsView from "./TagsView.vue";
</script>

<style scoped>
.admin-layout {
  display: flex;
  width: 100%;
  height: 100vh;
  overflow: hidden;
  background:
    linear-gradient(180deg, #fbfcfd 0%, #f4f5f7 44%, #eef1f5 100%);
}

.main-area {
  display: flex;
  min-width: 0;
  flex: 1;
  flex-direction: column;
}

.content-area {
  min-height: 0;
  flex: 1;
  overflow: auto;
  background:
    linear-gradient(90deg, rgb(255 255 255 / 42%), transparent 22%),
    linear-gradient(180deg, rgb(255 255 255 / 64%), transparent 190px),
    var(--app-bg);
}
</style>
