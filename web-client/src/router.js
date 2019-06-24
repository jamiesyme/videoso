import Vue from 'vue';
import Router from 'vue-router';
import Auth from '@/auth';
import HomeView from '@/views/Home.vue';

const views = {
	home:    HomeView,
	channel: () => import('@/views/Channel'),
	viewer:  () => import('@/views/Viewer'),
	search:  () => import('@/views/Search'),
	login:   () => import('@/views/Login'),
	admin:   () => import('@/views/Admin'),
};

function requireAuth (to, from, next) {
	if (Auth.loggedIn()) {
		next();
	} else {
		next('/login');
	}
}

Vue.use(Router);

export default new Router({
	mode: 'history',
	base: process.env.BASE_URL,
	routes: [
		{
			path: '/',
			component: views.home,
		},
		{
			path: '/video/:videoId',
			component: views.viewer,
		},
		{
			path: '/channel/:userName',
			component: views.channel,
		},
		{
			path: '/search',
			component: views.search,
		},
		{
			path: '/login',
			component: views.login,
			beforeEnter: (to, from, next) => {
				if (Auth.loggedIn()) {
					next('/admin');
				} else {
					next();
				}
			},
		},
		{
			path: '/admin',
			redirect: '/admin/categories',
		},
		{
			path: '/admin/categories',
			component: views.admin,
			props: {
				section: 'categories',
			},
			beforeEnter: requireAuth,
		},
		{
			path: '/admin/users',
			component: views.admin,
			props: {
				section: 'users',
			},
			beforeEnter: requireAuth,
		},
		{
			path: '/admin/videos',
			component: views.admin,
			props: {
				section: 'videos',
			},
			beforeEnter: requireAuth,
		},
	],
});
