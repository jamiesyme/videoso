export function expandVideo (content, video) {
	return Object.assign({}, video, {
		category: content.categories.find(cat => {
			return cat.id === video.category;
		}),
		author: content.users.find(user => {
			return user.id === video.author;
		}),
	});
}

export default {
	expandVideo,
}
