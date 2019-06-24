import Vue from 'vue';
import App from '@/App.vue';
import router from '@/router';
import Content from '@/content';
import BreakpointPlugin from '@/plugins/breakpoint';

Vue.config.productionTip = false;

Vue.use(BreakpointPlugin);

Content.load().catch(err => {
	console.error('failed to load content:', err);
});

const app = new Vue({
	router,
	render: h => h(App),
});
app.$mount('#app');
