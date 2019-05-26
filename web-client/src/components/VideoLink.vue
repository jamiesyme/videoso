<template>
	<router-link
		class="video-link"
		to="">
		<img :src="video.thumbnailUrl">
		<div class="video-title">{{ video.title }}</div>
	</router-link>
</template>

<script>
	export default {
		name: 'VideoLink',

		props: {
			video: Object,
		},

		data () {
			return {};
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
</style>
