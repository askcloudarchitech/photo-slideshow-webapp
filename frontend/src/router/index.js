import { createRouter, createWebHistory } from "vue-router";
import TVSlideshowView from "../views/TVSlideshowView.vue";
import Home from "../views/HomeView.vue";
import LoginView from "../views/LoginView.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "Home",
      component: Home,
    },
    {
      path: "/tv-slideshow",
      name: "TVSlideshow",
      component: TVSlideshowView,
    },
    {
      path: "/login",
      name: "Login",
      component: LoginView,
    },
  ],
});

router.beforeEach(async (to, from) => {
  const res = await fetch("/api/is-authenticated");
  if (res.status != 200) {
    if (to.name !== "Login") {
      return {
        path: "/login",
        query: {
          redirect_uri: encodeURIComponent(from.fullPath),
        },
      };
    }
  }
});

export default router;
