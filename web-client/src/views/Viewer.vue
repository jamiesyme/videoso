<template>
	<div class="viewer-view">
		<div class="container" v-if="video">
			<section class="video-section">
				<VideoViewer :video="video"></VideoViewer>
				<hr>
			</section>
			<section class="related-section">
				<h2>Related Videos</h2>
				<VideoLink
					:video="video"
					:key="video.id"
					v-for="video in relatedVideos">
				</VideoLink>
				<hr>
			</section>
		</div>

		<div class="container" v-else>
			No video.
		</div>
	</div>
</template>

<script>
	import Content from '@/content';
	import ContentUtils from '@/utils/content';
	import VideoLink from '@/components/VideoLink';
	import VideoViewer from '@/components/VideoViewer';

	export default {
		components: {
			VideoLink,
			VideoViewer,
		},

		data () {
			return {
				content: null,
			};
		},

		async mounted () {
			// Load content (if hasn't already been loaded)
			if (Content.videos.length === 0) {
				await Content.load();
			}
			this.content = Content;
		},

		computed: {
			videoId () {
				return this.$route.params.videoId;
			},

			video () {
				if (!this.videoId) {
					return null;
				}
				if (!this.content) {
					return null;
				}

				const video = this.content.videos.find(vid => {
					return vid.id == this.videoId;
				});
				return Object.assign(
					ContentUtils.expandVideo(this.content, video),
					{
						publishedAt: new Date(video.publishedAt),
					},
				);
			},

			relatedVideos () {
				const dummyThumbUrl = 'https://dummyimage.com/200x112/000/fff';

				if (!this.videoId) {
					return null;
				}
				if (!this.content) {
					return null;
				}

				return this.content.videos.filter(vid => {
					const sameId = vid.id === this.video.id;
					const sameCat = vid.category === this.video.category.id;
					return !sameId && sameCat;
				}).map(vid => {
					return Object.assign(
						{
							thumbnailUrl: dummyThumbUrl,
						},
						ContentUtils.expandVideo(this.content, vid),
					);
				});
			},
		},
	}
</script>

<style lang="scss" scoped>
	.container {
		display: flex;
		padding: 0;
		margin: 4rem auto 6rem;
	}

	.video-section {
		flex: 1 1 75%;
	}

	.related-section {
		flex: 1 1 25%;
		padding: 0 0 0 3rem;

		h2 {
			font-size: 2.2rem;
			letter-spacing: -0.08rem;
			line-height: 1.35;
			margin-bottom: 1.5rem;
		}

		.video-link {
			margin-bottom: 2rem;
		}
	}
</style>
