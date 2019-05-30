import CacheUtils from '@/utils/cache';

class Content {
	constructor (dataUrl) {
		this.categories = [];
		this.users = [];
		this.videos = [];
		this.dataUrl = dataUrl;
	}

	async load (options) {
		let content;
		if (!options || !options.fresh) {
			content = this.getCachedContent();
		}
		if (!content) {
			content = await this.getFreshContent();
			this.cacheContent(content);
		}

		this.categories = content.categories || [];
		this.users = content.users || [];
		this.videos = content.videos || [];
	}

	async save () {
		CacheUtils.removeItem(this.dataUrl);
		const res = await fetch(this.dataUrl, {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({
				categories: this.categories,
				users: this.users,
				videos: this.videos,
			}),
		});
		if (!res.ok) {
			console.err(res);
			throw new Error('failed to save content');
		}
	}

	getCachedContent () {
		if (CacheUtils.hasItem(this.dataUrl)) {
			const cacheStr = CacheUtils.getItem(this.dataUrl);
			const cacheObj = JSON.parse(cacheStr);
			if ((new Date).toISOString() < cacheObj.expiresAt) {
				return cacheObj.content;
			} else {
				CacheUtils.removeItem(this.dataUrl);
			}
		}
		return null;
	}

	cacheContent (content) {
		const expiresAt = new Date;
		expiresAt.setTime(expiresAt.getTime() + 60 * 60 * 1000);
		const cacheObj = { content, expiresAt };
		const cacheStr = JSON.stringify(cacheObj);
		CacheUtils.setItem(this.dataUrl, cacheStr);
	}

	async getFreshContent () {
		const res = await fetch(this.dataUrl);
		if (!res.ok) {
			console.err(res);
			throw new Error('failed to load content');
		}
		return await res.json();
	}
}

export default new Content(process.env.VUE_APP_CONTENT_URL);
