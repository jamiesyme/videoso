import Vue from 'vue';
import Router from 'vue-router';
import Home from './views/Home.vue';

Vue.use(Router);

export default new Router({
	mode: 'history',
	base: process.env.BASE_URL,
	routes: [
		{
			path: '/',
			name: 'home',
			component: Home
		},
		{
			path: '/video/:videoId',
			name: 'viewer',
			component: () => import(/* webpackChunkName: "viewer" */ '@/views/Viewer')
		},
		{
			path: '/channel/:userName',
			name: 'channel',
			component: () => import(/* webpackChunkName: "channel" */ '@/views/Channel')
		},
		{
			path: '/search',
			name: 'search',
			component: () => import(/* webpackChunkName: "search" */ '@/views/Search')
		},
		{
			path: '/login',
			name: 'login',
			component: () => import(/* webpackChunkName: "login" */ '@/views/Login')
		},
		{
			path: '/admin',
			name: 'admin',
			redirect: '/admin/categories'
		},
		{
			path: '/admin/categories',
			name: 'admin-categories',
			component: () => import(/* webpackChunkName: "admin" */ '@/views/Admin'),
			props: {
				section: 'categories',
			},
		},
		{
			path: '/admin/users',
			name: 'admin-users',
			component: () => import(/* webpackChunkName: "admin" */ '@/views/Admin'),
			props: {
				section: 'users',
			},
		},
		{
			path: '/admin/videos',
			name: 'admin-videos',
			component: () => import(/* webpackChunkName: "admin" */ '@/views/Admin'),
			props: {
				section: 'videos',
			},
		},
		/*{
			path: '/about',
			name: 'about',
			// route level code-splitting
			// this generates a separate chunk (about.[hash].js) for this route
			// which is lazy-loaded when the route is visited.
			component: () => import(/* webpackChunkName: "about" * / './views/About.vue')
		},*/
	],
});
