<template>
	<div class="search-view">
		<div class="container">
			<h1>Showing search results for: {{ query }}</h1>
			<div class="video-list">
				<VideoLink
					:video="video"
					:key="video.id"
					v-for="video in videos">
				</VideoLink>
			</div>
		</div>
	</div>
</template>

<script>
	import Content from '@/content';
	import ContentUtils from '@/utils/content';
	import VideoLink from '@/components/VideoLink';

	export default {
		components: {
			VideoLink,
		},

		data () {
			return {
				videos: [],
			};
		},

		mounted () {
			this.refreshContent();
		},

		methods: {
			async refreshContent () {
				const dummyThumbUrl = 'https://dummyimage.com/200x112/000/fff';

				// Load content (if hasn't already been loaded)
				if (Content.videos.length === 0) {
					await Content.load();
				}

				// Fill videos from content
				this.videos = Content.videos.map(vid => {
					return Object.assign(
						{
							thumbnailUrl: dummyThumbUrl,
						},
						ContentUtils.expandVideo(Content, vid),
						{
							publishedAt: new Date(vid.publishedAt),
						},
					);
				});
			},
		},

		computed: {
			query () {
				return this.$route.query.q;
			},
		},
	}
</script>

<style lang="scss" scoped>
	.container {
		padding: 0;
		margin: 4rem auto 6rem;
	}

	h1 {
		font-size: 2.4rem;
		line-height: 1.3;
		margin: 3rem 0;
	}

	.video-list {
		display: grid;
		grid-template-columns: 1fr 1fr 1fr 1fr;
		grid-column-gap: 0.4rem;
		grid-row-gap: 2.4rem;
	}
</style>
