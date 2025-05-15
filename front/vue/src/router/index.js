import { createRouter, createWebHistory } from "vue-router";
import LoginPage from "@/views/login";
import HomePage from "@/views/home";

const routes = [
  { path: "/" , redirect: "/login"},//默认路由
  { path: "/login", component: LoginPage },
  { path: "/home", component: HomePage }
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// 路由守卫，确保用户未登录时无法访问其他页面
// router.beforeEach((to, from, next) => {
//     const isLogin = localStorage.getItem("username");
//     if (isLogin == null && to.path !== "/login") {
//       next("/login");
//     } else if(isLogin != null && to.path === "/login") {
//       next("/main");
//     }else{
//       next();
//     }
//   });

export default router;
