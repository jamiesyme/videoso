import Vue from 'vue';
import CacheUtils from '@/utils/cache';

const dataUrl = process.env.VUE_APP_CONTENT_URL;

const state = Vue.observable({
	categories: [],
	users: [],
	videos: [],
	load: null,
	save: null,
});

function cacheContent (content) {
	const expiresAt = new Date;
	expiresAt.setTime(expiresAt.getTime() + 60 * 60 * 1000);
	const cacheObj = {
		content,
		expiresAt,
	};
	const cacheStr = JSON.stringify(cacheObj);
	CacheUtils.setItem(dataUrl, cacheStr);
}

function getCachedContent () {
	if (CacheUtils.hasItem(dataUrl)) {
		const cacheStr = CacheUtils.getItem(dataUrl);
		const cacheObj = JSON.parse(cacheStr);
		if ((new Date).toISOString() < cacheObj.expiresAt) {
			return cacheObj.content;
		} else {
			CacheUtils.removeItem(dataUrl);
		}
	}
	return null;
}

async function getFreshContent () {
	const res = await fetch(dataUrl);
	if (!res.ok) {
		console.err(res);
		throw new Error('failed to load content');
	}
	return await res.json();
}

async function load (options) {
	let content;
	if (!options || !options.fresh) {
		content = getCachedContent();
	}
	if (!content) {
		content = await getFreshContent();
		cacheContent(content);
	}

	state.categories = content.categories || [];
	state.users      = content.users || [];
	state.videos     = content.videos || [];
}

async function save () {
	CacheUtils.removeItem(dataUrl);
	const res = await fetch(dataUrl, {
		method: 'PUT',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify({
			categories: state.categories,
			users:      state.users,
			videos:     state.videos,
		}),
	});
	if (!res.ok) {
		console.err(res);
		throw new Error('failed to save content');
	}
}

state.load = load;
state.save = save;

export default state;
