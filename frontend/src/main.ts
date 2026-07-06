import "element-plus/dist/index.css";
import "./style.css";

import { createPinia } from "pinia";
import ElementPlus from "element-plus";
import { createApp } from "vue";
import * as ElementPlusIconsVue from "@element-plus/icons-vue";
import App from "./App.vue";
import router from "./router";

const app = createApp(App);

Object.entries(ElementPlusIconsVue).forEach(([key, component]) => {
  app.component(key, component);
});

app.use(createPinia());
app.use(ElementPlus);
app.use(router);
app.mount("#app");
