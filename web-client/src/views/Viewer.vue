<template>
	<div class="viewer-view">
		<div :class="containerClasses" v-if="video">
			<section class="video-section">
				<VideoViewer
					:video="video"
					:extraWidth="videoViewerExtraWidth">
				</VideoViewer>
				<hr>
			</section>
			<section class="related-section">
				<h2>Related Videos</h2>
				<div :class="relatedVideoListClasses">
					<VideoLink
						:video="video"
						:extraWidth="relatedVideoLinkExtraWidth"
						:key="video.id"
						v-for="video in relatedVideos">
					</VideoLink>
				</div>
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

			containerClasses () {
				const bp = this.$breakpoint.name;
				const multi = bp !== 'small' && bp !== 'medium';
				return {
					'container': true,
					'container-multi-column': multi,
				};
			},

			videoViewerExtraWidth () {
				const bp = this.$breakpoint.name;
				const extra = bp === 'small' || bp === 'medium';
				return extra ? '4rem' : null;
			},

			relatedVideoListClasses () {
				const bp = this.$breakpoint.name;
				const multi = bp === 'medium';
				return {
					'video-list': true,
					'video-list-multi': multi,
				};
			},

			relatedVideoLinkExtraWidth () {
				const bp = this.$breakpoint.name;
				const extra = bp === 'small';
				return extra ? '4rem' : null;
			},
		},
	}
</script>

<style lang="scss" scoped>
	.viewer-view {
		margin: 4rem 0 6rem;
	}

	.container.container-multi-column {
		display: grid;
		grid-template-columns: 3fr 1fr;
		grid-column-gap: 3rem;
	}

	.related-section {
		h2 {
			font-size: 2.2rem;
			letter-spacing: -0.08rem;
			line-height: 1.35;
			margin-bottom: 1.5rem;
		}

		.video-list {
			display: grid;
			grid-template-columns: 1fr;
			grid-row-gap: 2.4rem;

			&.video-list-multi {
				grid-column-gap: 0.4rem;
				grid-template-columns: 1fr 1fr;
			}
		}
	}
</style>
