import Vue from 'vue';
import BreakpointUtils from '@/utils/breakpoint';

function getCurrentBreakpoint () {
	const pointsByName = {
		small:  BreakpointUtils.isSmall(),
		medium: BreakpointUtils.isMedium(),
		large:  BreakpointUtils.isLarge(),
		xlarge: BreakpointUtils.isExtraLarge(),
	};
	return Object.keys(pointsByName).find(key => {
		return pointsByName[key];
	});
}

const breakpoint = Vue.observable({
	name: getCurrentBreakpoint(),
});

BreakpointUtils.addListener(() => {
	const current = getCurrentBreakpoint();
	if (breakpoint.name !== current) {
		breakpoint.name = current;
	}
});

export function install (Vue, options) {
	Vue.prototype.$breakpoint = breakpoint;
}

export default {
	install,
}
