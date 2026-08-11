import { createRouter, createWebHashHistory } from 'vue-router'
import { getToken } from '../services/auth.js'

import HomeView  from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{ path: '/',      redirect: '/chats'   },
		{ path: '/login', component: LoginView },
		{ path: '/chats', component: HomeView  },
	]
})

// Before each route, check if the user is logged in. If not, redirect to login page.
router.beforeEach((to, from, next) => {
	if (to.path !== '/login' && !getToken()) {
		next('/login')
	} else {
		next()
	}
})

export default router
