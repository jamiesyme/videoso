import Vue from 'vue';
import App from './App.vue';
import router from './router';
import BreakpointPlugin from '@/plugins/breakpoint';

Vue.config.productionTip = false;

Vue.use(BreakpointPlugin);

const app = new Vue({
	router,
	render: h => h(App),
});
app.$mount('#app');
