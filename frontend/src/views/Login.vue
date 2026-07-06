<template>
  <div class="login-page">
    <section class="login-panel">
      <div class="brand">
        <el-icon><Shop /></el-icon>
        <div>
          <h1>商城管理后台</h1>
          <p>权限、商品与库存统一管理</p>
        </div>
      </div>

      <el-form ref="formRef" :model="form" :rules="rules" label-position="top" @keyup.enter="submit">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" size="large" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" placeholder="请输入密码" size="large" show-password type="password" />
        </el-form-item>
        <el-button class="login-button" :loading="loading" size="large" type="primary" @click="submit">
          登录
        </el-button>
      </el-form>
    </section>
  </div>
</template>

<script setup lang="ts">
import { Shop } from "@element-plus/icons-vue";
import type { FormInstance, FormRules } from "element-plus";
import { ElMessage } from "element-plus";
import { reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import { useMenuStore } from "@/stores/menu";

const router = useRouter();
const route = useRoute();
const auth = useAuthStore();
const menu = useMenuStore();
const formRef = ref<FormInstance>();
const loading = ref(false);

const form = reactive({
  username: "admin",
  password: "Admin@123456"
});

const rules: FormRules<typeof form> = {
  username: [{ required: true, message: "请输入用户名", trigger: "blur" }],
  password: [{ required: true, message: "请输入密码", trigger: "blur" }]
};

async function submit(): Promise<void> {
  const valid = await formRef.value?.validate();
  if (!valid) {
    return;
  }

  loading.value = true;
  try {
    await auth.login(form);
    await menu.loadMenus();
    ElMessage.success("登录成功");
    const redirect = typeof route.query.redirect === "string" ? route.query.redirect : "/dashboard";
    await router.replace(redirect);
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.login-page {
  display: grid;
  min-height: 100vh;
  place-items: center;
  background:
    linear-gradient(135deg, rgb(64 158 255 / 14%), transparent 32%),
    linear-gradient(315deg, rgb(32 40 51 / 12%), transparent 36%),
    #f5f7fb;
}

.login-panel {
  width: 388px;
  border: 1px solid #dfe6ef;
  border-radius: 8px;
  background: #ffffff;
  padding: 30px;
  box-shadow: 0 18px 50px rgb(32 40 51 / 12%);
}

.brand {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 24px;
}

.brand .el-icon {
  display: grid;
  width: 46px;
  height: 46px;
  place-items: center;
  border-radius: 8px;
  background: #202833;
  color: #ffffff;
  font-size: 25px;
}

.brand h1 {
  margin: 0;
  font-size: 22px;
  letter-spacing: 0;
}

.brand p {
  margin: 4px 0 0;
  color: #6b7788;
}

.login-button {
  width: 100%;
}
</style>
