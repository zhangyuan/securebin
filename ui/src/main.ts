import { createApp } from "vue";
import { createRouter, createWebHashHistory } from "vue-router";

import "./style.css";
import App from "@/App.vue";
import Home from "@/components/Home.vue";
import View from "@/components/View.vue";

const routes = [
  { path: "/", component: Home, name: "home" },
  { path: "/messages/:id/:key", component: View, name: "view" },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

const app = createApp(App);
app.use(router);

app.mount("#app");
