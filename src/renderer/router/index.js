import { createRouter, createMemoryHistory } from "vue-router";
import LoginPage from "@/views/login/index.vue";
import HomePage from "@/views/home/index.vue";
import FriendsPage from "@/views/friend/index.vue";
import NoticePage from "@/views/notice/index.vue";
import TeamPage from "@/views/team/index.vue";
import PersonalPage from "@/views/profiles/index.vue";
import FriendDetail from "@/views/friend/FriendDetail.vue";
import NotFound from "@/views/404/index.vue"


const routes = [
  { path: "/" , redirect: "/login"},//默认路由
  { path: "/login", component: LoginPage },
  { path: "/home", component: HomePage },
  { path: "/friend", component: FriendsPage },
  { path: "/notice", component: NoticePage },
  { path: "/team", component: TeamPage },
  { path: "/profile", component: PersonalPage },
  { path: "/friendDetail/:id", component: FriendDetail },
  // 捕获所有未匹配的路由
  { path: '/:pathMatch(.*)', component: NotFound }
];

const router = createRouter({
  history: createMemoryHistory(),
  routes,
});

// 路由守卫，确保用户未登录时无法访问其他页面
router.beforeEach((to, from, next) => {
  const isLogin = localStorage.getItem("user");
  if (isLogin == null && to.path !== "/login") {
    next("/login");
  } else if(isLogin != null && to.path === "/login") {
    next("/home");
  }else{
    next();
  }
});

export default router;
