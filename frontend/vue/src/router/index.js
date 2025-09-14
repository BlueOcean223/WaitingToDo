import { createRouter, createWebHistory } from "vue-router";
import LoginPage from "@/views/login";
import HomePage from "@/views/home";
import FriendsPage from "@/views/friend";
import NoticePage from "@/views/notice";
import TeamPage from "@/views/team";
import PersonalPage from "@/views/profiles";
import FriendDetail from "@/views/friend/FriendDetail.vue";
import NotFound from "@/views/404"


const routes = [
  { path: "/" , component: HomePage},//默认路由
  { path: "/login", component: LoginPage },
  { path: "/friend", component: FriendsPage },
  { path: "/notice", component: NoticePage },
  { path: "/team", component: TeamPage },
  { path: "/profile", component: PersonalPage },
  { path: "/friendDetail/:id", component: FriendDetail },
  // 捕获所有未匹配的路由
  { path: '/:pathMatch(.*)', component: NotFound }
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// 路由守卫，确保用户未登录时无法访问其他页面
router.beforeEach((to, from, next) => {
  const isLogin = localStorage.getItem("user");
  if (isLogin == null && to.path !== "/login") {
    next("/login");
  } else if(isLogin != null && to.path === "/login") {
    next("/");
  }else{
    next();
  }
});

export default router;
