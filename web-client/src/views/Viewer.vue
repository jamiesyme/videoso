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
				video: null,
				relatedVideos: [],
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
				const rawVideo = Content.videos.find(vid => {
					return vid.id == this.$route.params.videoId;
				});
				if (rawVideo) {
					this.video = Object.assign(
						ContentUtils.expandVideo(Content, rawVideo),
						{
							publishedAt: new Date(rawVideo.publishedAt),
						},
					);

					this.relatedVideos = Content.videos.filter(vid => {
						const sameId = vid.id === rawVideo.id;
						const sameCat = vid.category === rawVideo.category;
						return !sameId && sameCat;
					}).map(vid => {
						return Object.assign(
							{
								thumbnailUrl: dummyThumbUrl,
							},
							ContentUtils.expandVideo(Content, vid),
						);
					});

				} else {
					this.video = null;
					this.relatedVideos = [];
				}
			},
		},

		watch: {
			'$route.params': function (newParams, oldParams) {
				if (newParams.videoId !== oldParams.videoId) {
					this.refreshContent();
				}
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
	}
</style>
