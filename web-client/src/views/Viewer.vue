<template>
	<div class="viewer-view">
		<div class="container" v-if="video">
			<section class="video-section">
				<div class="video-wrapper aspect-16-9">
					<video controls>
						<source :src="video.videoUrl">
					</video>
				</div>
				<div class="primary-info">
					<h1>{{ video.title }}</h1>
					<div class="view-count">
						{{ formatViewCount(video.viewCount) }}
					</div>
				</div>
				<div class="secondary-info">
					<div class="publish-info">
						Published by
						<router-link :to="`/channel/${video.author.name}`">
							{{ video.author.name }}
						</router-link>
						on
						{{ formatDate(video.publishedAt) }}
					</div>
					<div class="description" v-if="video.description">
						{{ video.description }}
					</div>
					<div class="description empty-description" v-else>
						No description.
					</div>
				</div>
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

	export default {
		components: {
			VideoLink,
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

			formatDate (date) {
				return date.toLocaleString('en-US', {
					month: 'short',
					day: 'numeric',
					year: 'numeric',
				});
			},

			formatViewCount (viewCount) {
				if (viewCount === 0) {
					return 'No views';
				}
				if (viewCount === 1) {
					return '1 view';
				}
				return `${viewCount} views`;
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

		.video-wrapper {
			background-color: black;
			position: relative;

			video {
				position: absolute;
				top: 0;
				left: 0;
				width: 100%;
				height: 100%;
				object-fit: contain;
			}
		}

		.primary-info {
			display: flex;
			margin: 2rem 0;
			padding: 0 1rem;

			h1 {
				flex: 1;
				font-size: 3.2rem;
				margin: 0;
			}

			.view-count {
				font-size: 2.0rem;
			}
		}

		.secondary-info {
			font-size: 1.4rem;
			padding: 0 1rem;

			.publish-info {
				a {
					font-weight: bold;
				}
			}

			.description {
				margin: 2rem 0 0;
			}

			.empty-description {
				font-style: italic;
			}
		}
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
			font-size: 1.4rem;
		}
	}
</style>
