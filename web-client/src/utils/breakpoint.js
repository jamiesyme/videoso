const querySmall  = window.matchMedia('(min-width: 600px)');
const queryMedium = window.matchMedia('(min-width: 900px)');
const queryLarge  = window.matchMedia('(min-width: 1200px)');

let listeners = [];

export function isSmall () {
	return !querySmall.matches;
}

export function isMedium () {
	return querySmall.matches && !queryMedium.matches;
}

export function isLarge () {
	return queryMedium.matches && !queryLarge.matches;
}

export function isExtraLarge () {
	return queryLarge.matches;
}

/**
 * @param {Function} func
 * @param {object}   [options]
 * @param {*}        [options.this]
 */
export function addListener (func, options) {
	let context;
	if (options && options.this) {
		context = options.this;
	}
	listeners.push({
		func,
		this: context,
	});
}

export function removeListener (func) {
	const i = listeners.findIndex(l => l.func === func);
	if (i >= 0) {
		listeners.splice(i, 1);
	}
}

function onChange () {
	for (const listener of listeners) {
		listener.func.call(listener.this);
	}
}

querySmall.addListener(onChange);
queryMedium.addListener(onChange);
queryLarge.addListener(onChange);

export default {
	isSmall,
	isMedium,
	isLarge,
	isExtraLarge,
	addListener,
	removeListener,
}
