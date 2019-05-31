<template>
	<router-link
		class="video-link"
		:to="`/video/${video.id}`">
		<img :src="video.thumbnailUrl">
		<div class="video-title">{{ video.title }}</div>
		<div class="video-publisher">
			{{ video.author.name }}
		</div>
	</router-link>
</template>

<script>
	export default {
		name: 'VideoLink',

		props: {
			/**
			 * @typedef {object} Video
			 * @property {number} id
			 * @property {string} title
			 * @property {string} thumbnailUrl
			 * @property {object} author
			 * @property {number} author.id
			 * @property {string} author.name
			 */
			video: Object,
		},
	}
</script>

<style lang="scss" scoped>
	.video-link {
		display: block;
	}

	img {
		display: block;
		width: 100%;
		height: auto;
	}

	.video-title {
		font-weight: bold;
		margin: 0.4rem 0 0;
	}

	.video-publisher {
		color: #606c76;
	}
</style>
