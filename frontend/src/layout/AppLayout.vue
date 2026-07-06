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
  background: #f5f7fb;
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
}
</style>
