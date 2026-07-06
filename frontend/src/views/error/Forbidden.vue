<template>
  <div class="error-page">
    <h1>403</h1>
    <p>当前账号没有访问该页面的权限。</p>
    <el-button type="primary" @click="$router.push('/dashboard')">返回仪表盘</el-button>
  </div>
</template>

<style scoped>
.error-page {
  display: grid;
  min-height: 100vh;
  place-items: center;
  align-content: center;
  gap: 12px;
  color: #202833;
}

h1 {
  margin: 0;
  font-size: 72px;
}

p {
  margin: 0;
  color: #6b7788;
}
</style>
